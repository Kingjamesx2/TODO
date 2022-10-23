//filename: internal/data/todo.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"todo.jamesfaber.net/internal/validator"
)

type Todo struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Task      string    `json:"task"`
	Version   int32     `json:"version"`
}

func ValidateTodo(v *validator.Validator, todo *Todo) {
	// Use the check() method to execute our validation checks
	v.Check(todo.Name != "", "name", "must be provided")
	v.Check(len(todo.Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(todo.Task != "", "task", "must be provided")
	v.Check(len(todo.Task) <= 200, "task", "must not be more than 200 bytes long")
}

// Define a todo list model which wraps a sql.DB connection pool
type TodoModel struct {
	DB *sql.DB
}

// Insert() allows us to create a new todo task
func (m TodoModel) Insert(todo *Todo) error {
	query := `
	INSERT INTO todo (name, task)
	VALUES ($1, $2)
	RETURNING id, created_at, version
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	// Collect the data fields into a slice
	args := []interface{}{todo.Name, todo.Task}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt, &todo.Version)
}

// GET() allows us to retrieve a specific todo item
func (m TodoModel) Get(id int64) (*Todo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Create query
	query := `
		SELECT id, created_at, name, task, version
		FROM todo
		WHERE id = $1
	`
	// Declare a Todo variable to hold the return data
	var todo Todo
	// Execute Query using the QueryRow
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.Name,
		&todo.Task,
		&todo.Version,
	)
	// Handle any errors
	if err != nil {
		// Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Success
	return &todo, nil
}

// Update() allows us to edit/alter a todo item in the list
func (m TodoModel) Update(todo *Todo) error {
	query := `
		UPDATE todo 
		set name = $1, task = $2,  
		version = version + 1
		WHERE id = $3
		AND version = $4
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	args := []interface{}{
		todo.Name,
		todo.Task,
		todo.ID,
		todo.Version,
	}
	// Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete() removes a specific Task
func (m TodoModel) Delete(id int64) error {
	// Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	// Create the delete query
	query := `
		DELETE FROM todo
		WHERE id = $1
	`
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	// Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Check how many rows were affected by the delete operation. We
	// call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// the GetAll() method returns a list of all the Todo sorted by id
func (m TodoModel) GetAll(name string, task string, filters Filters) ([]*Todo, Metadata, error) {
	//construct the query to return all todo
	//make query into formated string to be able to sort by field and asc or dec dynaimicaly
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(),id, created_at, name, task, version
		FROM todo
		WHERE (to_tsvector('simple',name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple',task) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortOrder())

	//create a 3 second timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//execute the query
	args := []interface{}{name, task, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	//close the result set
	defer rows.Close()
	//store total records
	totalRecords := 0
	//intialize an empty slice to hold the Todo data
	todos := []*Todo{}
	//iterate over the rows in the result set
	for rows.Next() {
		var todo Todo
		//scan the values from the row into the Todo struct
		err := rows.Scan(
			&totalRecords,
			&todo.ID,
			&todo.CreatedAt,
			&todo.Name,
			&todo.Task,
			&todo.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		//add the todo to our slice
		todos = append(todos, &todo)
	}
	//check if any errors occured while proccessing the result set
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	//return the result set. the slice of todos
	return todos, metadata, nil
}

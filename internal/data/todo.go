//filename: internal/data/todo.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"todo.jamesfaber.net/internal/validator"
)

type Todo struct {
	ID        int64     `json: "id"`
	CreatedAt time.Time `json: "created_at"`
	Name      string    `json: "name"`
	Task      string    `json: "task"`
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

// insert() allows us to crete a new todo task
func (m TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO (name, task)
		VALUE ($1, $2)
		RETURNING id, created_at, version
	`
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	//colllect the data fields into a slice
	args := []interface{}{
		todo.Name,
		todo.Task,
	}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt, &todo.Version)
}

// Get() allows us to retrieve a specific task
func (m TodoModel) Get(id int64) (*Todo, error) {
	// Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Create the query
	query := `
		SELECT  id, created_at, name, task, version
		FROM todo
		WHERE id = $1
	`
	// Declare a todo variable to hold the returned data
	var todo Todo

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	// Execute the query using QueryRow()
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

// Update() allows us to edi/alter a specific tool
func (m TodoModel) Update(todo *Todo) error {
	// Create a query
	query := `
		UPDATE todo
		SET name = $1, task = $2, version = version + 1
		WHERE id = $3
		AND version = $4
		RETURNING version
	`
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	args := []interface{}{
		todo.Name,
		todo.Task,
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

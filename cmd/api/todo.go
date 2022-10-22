//FIlename: cmd/api/todo.go

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"todo.jamesfaber.net/internal/data"
	"todo.jamesfaber.net/internal/validator"
)

// createTodoInfoHandler for the POST /v1/todoInfo" endpoints
func (app *application) createTodoInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Our target decode destination
	var input struct {
		Name string `json: "name"`
		Task string `json: "task"`
	}
	// Initialize a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//Copy the values from the input struct to a new Todo struct
	todo := &data.Todo{
		Name: input.Name,
		Task: input.Task,
	}

	//Initialize a new Validator instance
	v := validator.New()

	//check the map to determine if there were any validation errors
	if data.ValidateTodo(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Create a todo
	err = app.models.Todo.Insert(todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	//Create a location header for the newly created resource/todo
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/todoInfo/%d", todo.ID))

	//Write the JSON response with 201 - Created status code with the body
	//being the todo data and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// showTodoInfoHandler for the "GET /v1/todoInfo/" endpoint
func (app *application) showTodoInfoHandler(w http.ResponseWriter, r *http.Request) {
	//use the "ParamsFromCOntext()" function to get the reuest context as a slice
	params := httprouter.ParamsFromContext(r.Context())
	// GET the value of the "id" parameter
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	//fetch the specific Task
	todo, err := app.models.Todo.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateTodoInfoHandler(w http.ResponseWriter, r *http.Request) {
	// This method does a partial replacement
	// Get the id for the todo that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the orginal record from the database
	todo, err := app.models.Todo.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Create an input struct to hold data read in from the client
	// We update input struct to use pointers because pointers have a
	// default value of nil
	// If a field remains nil then we know that the client did not update it
	var input struct {
		Name *string `json:"name"`
		Task *string `json:"task"`
	}

	// Initialize a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check for updates
	if input.Name != nil {
		todo.Name = *input.Name
	}
	if input.Task != nil {
		todo.Task = *input.Task
	}

	// Perform validation on the updated Todo. If validation fails, then
	// we send a 422 - Unprocessable Entity respose to the client
	// Initialize a new Validator instance
	v := validator.New()

	// Check the map to determine if there were any validation errors
	if data.ValidateTodo(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updated Todo record to the Update() method
	err = app.models.Todo.Update(todo)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTodoInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id for the todo that needs deleting
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the Todo from the database. Send a 404 Not Found status code to the
	// client if there is no matching record
	err = app.models.Todo.Delete(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return 200 Status OK to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Todo Task successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// The listTodoInfosHandler() allows the client to see a listing of Todo task based on a set of criteria
func (app *application) listTodoInfoHandler(w http.ResponseWriter, r *http.Request) {
	//Create an input struct to hold our query parameters
	var input struct {
		Name string
		Task string
		data.Filters
	}

	//Initialize a validator
	v := validator.New()
	//Get the URL values map
	qs := r.URL.Query()
	//Use the helper methods to extract the values
	input.Name = app.readString(qs, "name", "")
	input.Task = app.readString(qs, "Task", "")
	//Get the page information
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	//Get the sort information
	input.Filters.Sort = app.readString(qs, "sort", "id")
	//Specify the allowed sort values

	//CHECK THIS
	input.Filters.SortList = []string{"id", "name", "task", "-id", "-name", "-task"}

	//Checking for validation errors
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Get a listing of all Todo Info
	todo, metadata, err := app.models.Todo.GetAll(input.Name, input.Task, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	//Send JSON responce containing all the todo info
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

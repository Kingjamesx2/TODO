//FIlename: cmd/pi/todo.go

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

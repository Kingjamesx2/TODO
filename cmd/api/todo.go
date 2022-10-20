//FIlename: cmd/pi/todo.go

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	//Display the request
	fmt.Fprintf(w, "%+v\n", input)
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
	//create a new instance of the todo struct containing the ID we extracted from our URL and some sample data
	//display the id
	todo := data.Todo{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Hilary To DO list",
		Task:      "bake cake for birthday",
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encounteed a problem and could not process your reuest", http.StatusInternalServerError)
	}

}

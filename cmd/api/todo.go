//FIlename: cmd/pi/todo.go

package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// createTodoInfoHandler for the POST /v1/todoInfo" endpoints
func (app *application) createTodoInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "create a todo list")
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
	//display the id
	fmt.Fprintf(w, "show the details for school %d\n", id)

}

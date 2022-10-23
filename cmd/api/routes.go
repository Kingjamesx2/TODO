//FILename: cmd/api/routes

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	//Create a new  httprouter ruter instance
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/todoInfo/", app.listTodoInfoHandler)

	router.HandlerFunc(http.MethodPost, "/v1/todoInfo", app.createTodoInfoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/todoInfo/:id", app.showTodoInfoHandler)
	router.HandlerFunc(http.MethodPut, "/v1/todoInfo/:id", app.updateTodoInfoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/todoInfo/:id", app.deleteTodoInfoHandler)

	return router
}

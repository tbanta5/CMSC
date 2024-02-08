package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) Routes() *httprouter.Router {
	router := httprouter.New()

	// // Declare a session manager middleware for certain routes.
	// dynamic := alice.New(app.sessionManager.LoadAndSave)
	// Server crude html page.
	router.Handler(http.MethodGet, "/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.index)))
	router.Handler(http.MethodGet, "/v1/coffee", app.sessionManager.LoadAndSave(http.HandlerFunc(app.listCoffees)))
	router.Handler(http.MethodGet, "/v1/coffee/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.getDescription)))

	// Liveness is used by kubernetes
	router.HandlerFunc(http.MethodGet, "/v1/liveness", app.liveness)

	return router
}

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
	router.Handler(http.MethodGet, "/coffee", app.sessionManager.LoadAndSave(http.HandlerFunc(app.coffees)))
	router.Handler(http.MethodGet, "/coffee/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.coffeeDesc)))

	// Liveness is used by kubernetes
	router.HandlerFunc(http.MethodGet, "/liveness", app.liveness)

	return router
}

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) Routes() *httprouter.Router {
	router := httprouter.New()

	// Server crude html page.
	router.Handler(http.MethodGet, "/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.index)))
	router.Handler(http.MethodGet, "/coffee", app.sessionManager.LoadAndSave(http.HandlerFunc(app.coffees)))
	router.Handler(http.MethodGet, "/coffee/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.coffeeDesc)))

	// Liveness is used by kubernetes
	router.HandlerFunc(http.MethodGet, "/liveness", app.liveness)

	return router
}

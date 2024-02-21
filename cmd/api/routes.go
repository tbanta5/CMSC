package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) Routes() *httprouter.Router {
	router := httprouter.New()

	// Server crude html page.
	router.Handler(http.MethodGet, "/", http.HandlerFunc(app.index))
	router.Handler(http.MethodGet, "/coffee", app.sessionManager.LoadAndSave(http.HandlerFunc(app.coffees)))
	router.Handler(http.MethodGet, "/coffee/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.coffeeDetails)))
	// Add middleware to the below endpoints to validate authToken
	router.Handler(http.MethodPost, "/coffee", http.HandlerFunc(app.newCoffee))

	router.Handler(http.MethodPost, "/cart/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.addCoffee)))
	router.Handler(http.MethodDelete, "/cart/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.removeCoffee)))
	router.Handler(http.MethodGet, "/cart", app.sessionManager.LoadAndSave(http.HandlerFunc(app.shoppingCart)))

	router.Handler(http.MethodPost, "/auth", app.sessionManager.LoadAndSave(http.HandlerFunc(app.auth)))

	// Liveness is used by kubernetes
	router.HandlerFunc(http.MethodGet, "/liveness", app.liveness)

	return router
}

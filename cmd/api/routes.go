package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) Routes() *httprouter.Router {
	router := httprouter.New()

	// Server crude html page.
	router.Handler(http.MethodGet, "/", http.HandlerFunc(app.index))
	router.Handler(http.MethodGet, "/coffee", app.sessionManager.LoadAndSave(app.loadCoffeeList(http.HandlerFunc(app.coffees))))
	router.Handler(http.MethodGet, "/coffee/:id", app.sessionManager.LoadAndSave(app.loadCoffeeList(http.HandlerFunc(app.coffeeDetails))))
	// Add middleware to the below endpoints to validate authToken
	router.Handler(http.MethodPost, "/coffee", app.authenticate(http.HandlerFunc(app.newCoffee)))
	router.Handler(http.MethodDelete, "/coffee/:id", app.authenticate(http.HandlerFunc(app.deleteCoffee)))
	router.Handler(http.MethodPatch, "/coffee/:id", app.authenticate(http.HandlerFunc(app.updateCoffee)))

	router.Handler(http.MethodPost, "/cart/:id", app.sessionManager.LoadAndSave(app.loadCoffeeList(http.HandlerFunc(app.addCoffee))))
	router.Handler(http.MethodDelete, "/cart/:id", app.sessionManager.LoadAndSave(app.loadCoffeeList(http.HandlerFunc(app.removeCoffee))))
	router.Handler(http.MethodGet, "/cart", app.sessionManager.LoadAndSave(app.loadCoffeeList(http.HandlerFunc(app.shoppingCart))))

	// Generate a new authentication token for admin usage
	router.Handler(http.MethodPost, "/auth", http.HandlerFunc(app.createAuthToken))

	// Liveness is used by kubernetes
	router.HandlerFunc(http.MethodGet, "/liveness", app.liveness)

	return router
}

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
	router.Handler(http.MethodGet, "/coffee/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.coffeeDetails)))

	router.Handler(http.MethodPost, "/cart/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.addCoffee)))
	router.Handler(http.MethodPost, "/admin/coffee", app.sessionManager.LoadAndSave(http.HandlerFunc(app.adminAddCoffee)))
	router.Handler(http.MethodGet, "/cart", app.sessionManager.LoadAndSave(http.HandlerFunc(app.shoppingCart)))
	// Liveness is used by kubernetes
	router.HandlerFunc(http.MethodGet, "/liveness", app.liveness)

	return router
}

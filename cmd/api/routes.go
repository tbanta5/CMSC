package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) Routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.index)

	router.HandlerFunc(http.MethodGet, "/v1/liveness", app.liveness)
	router.HandlerFunc(http.MethodGet, "/v1/readiness", app.readiness)

	return router
}

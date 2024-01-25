package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) Routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/liveness", app.liveness)
	return router
}

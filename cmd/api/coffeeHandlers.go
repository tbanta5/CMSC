package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"cmsc.group2.coffee-api/internal/dataModels"
	"github.com/julienschmidt/httprouter"
)

func (app *application) coffees(w http.ResponseWriter, r *http.Request) {

}

// curl -XGET localhost:8585/coffee/1 - understand 1 is the id here.
func (app *application) coffeeDetails(w http.ResponseWriter, r *http.Request) {
	// Get the parameters from the url context
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Error("parsing id param", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Define a timeout for the database retrieval
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Pass our timeout context and app.db connection
	// info to coffee method.
	coffee, err := dataModels.CoffeeDetails(ctx, app.db, id)
	if err != nil {
		app.logger.Error("dataModels.CoffeDetails", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	// Format response to json
	js, err := json.Marshal(coffee)
	if err != nil {
		app.logger.Error("marshal json", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	// Write response to http.ResponseWriter
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

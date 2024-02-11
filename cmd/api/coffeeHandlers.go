package main

import (
	"encoding/json"
	"net/http"

	"cmsc.group2.coffee-api/internal/dataModels"
)

func (app *application) coffees(w http.ResponseWriter, r *http.Request) {

}

// NOT REALLY FUNCTIONAL JUST AN EXAMPLE
// curl -XGET localhost:8585/coffee/1 - understand 1 is the id here.
func (app *application) coffeeDesc(w http.ResponseWriter, r *http.Request) {
	coffee := dataModels.Coffee{
		Name:        "Latte",
		Description: "Frothy coffee drink",
		Price:       3.85,
	}

	js, err := json.Marshal(coffee)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

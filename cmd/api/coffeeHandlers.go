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

// coffees returns a list of coffees to the user. The function
// is optimized to store the coffee list retrieved from db
// in the session data - thereby reducing db calls.
func (app *application) coffees(w http.ResponseWriter, r *http.Request) {
	// First check the session data to see if we already have the product list
	// Note, the type assertion here will return a fully nil object
	// So we cannot proliferate coffeeList []Coffee{} outside of this declaration
	// Hence a little code duplication below.
	coffeeList, ok := app.sessionManager.Get(r.Context(), "coffeeList").([]dataModels.Coffee)
	if !ok {
		// If coffeeList doesn't exists yet, we call the database.
		// Define a timeout for the database retrieval
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Call the database
		coffeeList, err := dataModels.CoffeeList(ctx, app.db)
		if err != nil {
			app.logger.Error("dataModels.CoffeList", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		// Cache the coffee list in the session token
		// this is to avoid subsequent database queries
		// if list becomes too large - this is a bad optimization.
		app.sessionManager.Put(r.Context(), "coffeeList", coffeeList)
		coffees := make([]dataModels.Coffee, len(coffeeList))
		for i, v := range coffeeList {
			coffees[i] = dataModels.Coffee{
				ID:          v.ID,
				Name:        v.Name,
				Description: v.Description,
				Price:       v.Price,
			}
		}
		// Format response to json
		js, err := json.Marshal(coffees)
		if err != nil {
			app.logger.Error("marshal json", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		// Write response to http.ResponseWriter
		js = append(js, '\n')
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	} // End initial DB Call for list of products

	coffees := make([]dataModels.Coffee, len(coffeeList))
	for i, v := range coffeeList {
		coffees[i] = dataModels.Coffee{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
		}
	}

	// Format response to json
	js, err := json.Marshal(coffees)
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

// coffeeDetails will return the specified coffee product
// with additional details such as caffeine and calories.
// This function attempts to avoid unnecessary db calls
// by first checking the session for list of products.
func (app *application) coffeeDetails(w http.ResponseWriter, r *http.Request) {
	// Get the parameters from the request url context, ie ":id"
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Error("parsing id param", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Pull coffeeList from the session data.
	// We must type assert to type Coffee
	coffeeList, ok := app.sessionManager.Get(r.Context(), "coffeeList").([]dataModels.Coffee)
	if !ok {
		app.logger.Error("Session doesn't contain coffeeList")
		http.Error(w, "Coffee products not available", http.StatusBadRequest)
		return
	}

	coffeeDesc := dataModels.Coffee{}
	for _, coffee := range coffeeList {
		if coffee.ID == id {
			coffeeDesc = coffee
			break
		}
	}

	js, err := json.Marshal(coffeeDesc)
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

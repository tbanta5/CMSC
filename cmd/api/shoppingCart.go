package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cmsc.group2.coffee-api/internal/dataModels"
	"github.com/julienschmidt/httprouter"
)

func (app *application) shoppingCart(w http.ResponseWriter, r *http.Request) {
	shoppingCart, ok := app.sessionManager.Get(r.Context(), "shoppingCart").([]dataModels.Coffee)
	if !ok {
		shoppingCart = []dataModels.Coffee{}
	}
	// Here we need to calculate the items price in shopping cart and return a price
	js, err := json.Marshal(shoppingCart)
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

func (app *application) addCoffee(w http.ResponseWriter, r *http.Request) {
	// Get the parameters from the request url context, ie ":id"
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Error("parsing id param", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Ensure id is valid, then add to user's selection
	coffeeList, ok := app.sessionManager.Get(r.Context(), "coffeeList").([]dataModels.Coffee)
	if !ok {
		app.logger.Error("Session doesn't contain coffeeList")
		http.Error(w, "Coffee products not available", http.StatusBadRequest)
		return
	}

	coffee := dataModels.Coffee{}
	for _, c := range coffeeList {
		if c.ID == id {
			coffee = c
			break
		}
	}

	if coffee.Name == "" {
		app.logger.Error("User selected value does not exist", "id", strconv.Itoa(id))
		http.Error(w, "Invalid Selection", http.StatusBadRequest)
		return
	}

	// Retrieve current session data
	shoppingCart, ok := app.sessionManager.Get(r.Context(), "shoppingCart").([]dataModels.Coffee)
	if !ok {
		// If nothing in the cart, go ahead and create it.
		shoppingCart = []dataModels.Coffee{}
	}
	// Add coffee product to shopping cart
	shoppingCart = append(shoppingCart, coffee)
	app.sessionManager.Put(r.Context(), "shoppingCart", shoppingCart)

	msg := map[string]string{"success": "cart updated"}
	js, err := json.Marshal(msg)
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

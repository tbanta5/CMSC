package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cmsc.group2.coffee-api/internal/dataModels"
	"github.com/julienschmidt/httprouter"
)

// Define a new struct to hold the items and the total
type userCart struct {
	Items []dataModels.Coffee `json:"items"`
	Total float64             `json:"total"`
}

// shoppingCart lists a users shopping cart
// this list is tied to a sessionCookie.
func (app *application) shoppingCart(w http.ResponseWriter, r *http.Request) {
	// Retrieve the current shopping cart from the session
	shoppingCart, ok := app.sessionManager.Get(r.Context(), "shoppingCart").([]dataModels.Coffee)
	if !ok {
		shoppingCart = []dataModels.Coffee{}
	}

	// Calculate the total price of the items in the shopping cart
	var total float64
	for _, item := range shoppingCart {
		total += item.Price // Assuming that the Price field is a float64
	}

	// Create a userCart struct to hold the items and the total
	cart := userCart{
		Items: shoppingCart,
		Total: total,
	}

	// Marshal the userCart struct into JSON
	js, err := json.Marshal(cart)
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

// addCoffee adds a coffee to a users shopping cart
// The user cart is tied to a session cookie.
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
		app.logger.Error("User selected value does not exist", fmt.Errorf("id: %d", id))
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

	msg := fmt.Sprintf("Coffee %d added", id)
	success := map[string]string{"success": msg}
	js, err := json.Marshal(success)
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

// removeCoffee removes coffee from users shopping cart
// The user cart is tied to a sessionCookie
func (app *application) removeCoffee(w http.ResponseWriter, r *http.Request) {
	// Get the coffee ID from the URL parameter
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Error("parsing coffee ID", err)
		http.Error(w, "Invalid coffee ID", http.StatusBadRequest)
		return
	}

	// Retrieve the current shopping cart from the session
	shoppingCart, ok := app.sessionManager.Get(r.Context(), "shoppingCart").([]dataModels.Coffee)
	if !ok {
		app.logger.Error("Session doesn't contain shoppingCart")
		http.Error(w, "Shopping cart not found", http.StatusBadRequest)
		return
	}

	// Find and remove the coffee from the shopping cart
	updatedCart := []dataModels.Coffee{}
	count := 1 // Use a counter to ensure only one item is deleted at a time
	for _, coffee := range shoppingCart {
		if coffee.ID == id && count == 1 {
			count = 0
		} else {
			updatedCart = append(updatedCart, coffee)
		}
	}

	// If the length of the cart is the same after the removal attempt,
	// the item was not found
	if len(updatedCart) == len(shoppingCart) {
		app.logger.Error("Coffee ID not found in cart", fmt.Errorf("id: %d", id))
		http.Error(w, "Coffee ID not found in cart", http.StatusBadRequest)
		return
	}

	// Save the updated cart back to the session
	app.sessionManager.Put(r.Context(), "shoppingCart", updatedCart)

	// Write a successful response
	msg := fmt.Sprintf("Deleted coffee %d", id)
	success := map[string]string{"success": msg}
	js, err := json.Marshal(success)
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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cmsc.group2.coffee-api/internal/auth"
	"cmsc.group2.coffee-api/internal/dataModels"
	"cmsc.group2.coffee-api/internal/validation"
	"github.com/julienschmidt/httprouter"
)

func (app *application) coffees(w http.ResponseWriter, r *http.Request) {
	// Define a timeout for the database retrieval
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	coffeeList, err := dataModels.CoffeeList(ctx, app.db)
	if err != nil {
		app.logger.Error("dataModels.CoffeList", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(coffeeList)
	// Cache the coffee list in the session token
	// this is to avoid subsequent database queries
	// if list becomes too large - this is a bad optimization.
	app.sessionManager.Put(r.Context(), "coffee_list", coffeeList)

	coffees := []dataModels.Coffee{}
	for _, v := range coffeeList {
		coffee := dataModels.Coffee{}
		// Give only the basic coffee overview for this call
		coffee.ID = v.ID
		coffee.Name = v.Name
		coffee.Description = v.Description
		coffee.Price = v.Price
		coffees = append(coffees, coffee)
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

// curl -XGET localhost:8585/coffee/1 - understand 1 is the id here.
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
	coffeeList, ok := app.sessionManager.Get(r.Context(), "coffee_list").([]dataModels.Coffee)
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

func (app *application) adminAddCoffee(w http.ResponseWriter, r *http.Request) {
	// We expect a content-type of application/json for this request
	if r.Header.Get("Content-Type") != "application/json" {
		app.logger.Error("Content-Type is not application/json")
		http.Error(w, "Content-Type header is not application/json", http.StatusBadRequest)
		return
	}

	var newCoffee dataModels.Coffee
	err := json.NewDecoder(r.Body).Decode(&newCoffee)
	if err != nil {
		app.logger.Error("decode json", err)
		http.Error(w, "Error decoding JSON body", http.StatusBadRequest)
		return
	}

	// Call the ValidateCoffee function from the validation package
	err = validation.ValidateCoffee(&newCoffee)
	if err != nil {
		// Handle validation errors
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Define a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Insert the coffee into the database
	err = newCoffee.AddCoffee(ctx, app.db)
	if err != nil {
		app.logger.Error("inserting coffee", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Respond to the client with the ID of the new coffee
	w.Header().Set("Location", fmt.Sprintf("/coffee/%d", newCoffee.ID))
	w.WriteHeader(http.StatusCreated)
	// Optionally return the new coffee object in the response
	js, _ := json.Marshal(newCoffee)
	w.Write(js)
}

func (app *application) adminLogin(w http.ResponseWriter, r *http.Request) {
	// Decode the request body to get admin credentials
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if valid, err := auth.AuthenticateAdminCredentials(app.datastore, creds.Username, creds.Password); !valid {
		if err != nil {
			// Log the error for internal tracking
			app.logger.Error("authenticateAdminCredentials failed", err)
		}
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// If authentication is successful, generate a JWT token
	token, err := auth.GenerateAdminToken(creds.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the token in a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(theError)
}

package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"cmsc.group2.coffee-api/internal/dataModels"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure that we indicate authorization may vary
		w.Header().Add("Vary", "Authorization")

		// Returns "" empty string if nothing is found.
		authorizationHeader := r.Header.Get("Authorization")

		// If no auth header, set user as anonymous
		if authorizationHeader == "" {
			w.Header().Set("WWW-Authenticate", "Bearer")
			message := "invalid or missing auth token"
			http.Error(w, message, http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.Header().Set("WWW-Authenticate", "Bearer")
			message := "invalid or missing auth token"
			http.Error(w, message, http.StatusUnauthorized)
			return
		}

		token := headerParts[1]

		valid, err := dataModels.IsValidAdmin(app.db, token)
		if err != nil {
			app.logger.Error("Error with token validation", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		if !valid {
			w.Header().Set("WWW-Authenticate", "Bearer")
			message := "invalid or missing auth token"
			http.Error(w, message, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) loadCoffeeList(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure coffeeList is in user session from the start.
		_, ok := app.sessionManager.Get(r.Context(), "coffeeList").([]dataModels.Coffee)
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
			app.sessionManager.Put(r.Context(), "coffeeList", coffeeList)

		}
		next.ServeHTTP(w, r)
	})
}

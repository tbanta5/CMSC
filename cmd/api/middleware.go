package main

import (
	"net/http"
	"strings"

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

		if valid := dataModels.IsValidAdmin(token); !valid {
			w.Header().Set("WWW-Authenticate", "Bearer")
			message := "invalid or missing auth token"
			http.Error(w, message, http.StatusUnauthorized)
			return
		}
	})
}

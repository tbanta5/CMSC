package main

import (
	"encoding/json"
	"net/http"

	"cmsc.group2.coffee-api/internal/auth"
)

func (app *application) auth(w http.ResponseWriter, r *http.Request) {
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

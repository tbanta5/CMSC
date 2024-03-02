package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cmsc.group2.coffee-api/internal/dataModels"
	"cmsc.group2.coffee-api/internal/validation"
)

func (app *application) createAuthToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Error("decode json", err)
		http.Error(w, "Error decoding JSON body", http.StatusBadRequest)
		return
	}
	v := validation.New()

	dataModels.ValidateEmail(v, input.Email)
	dataModels.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.logger.Error("validation failed")
		http.Error(w, "Authentication failure", http.StatusUnauthorized)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := dataModels.GetByEmail(ctx, app.db, input.Email)

	if err != nil {
		app.logger.Error("User account not found", err)
		http.Error(w, "Authentication failure", http.StatusUnauthorized)
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.logger.Error("Account mismatch", err)
		http.Error(w, "Authentication failure", http.StatusUnauthorized)
		return
	}

	if !match {
		http.Error(w, "Authentication failure", http.StatusUnauthorized)
	}

	// If password is correct, we issue the token
	token, err := dataModels.NewToken(app.db, user.ID, 4*time.Hour, dataModels.ScopeAuthentication)
	if err != nil {
		app.logger.Error("Error from New Token creation", err)
		http.Error(w, "error processing", http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(token)
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

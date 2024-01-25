package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) liveness(w http.ResponseWriter, r *http.Request) {
	msg := map[string]string{
		"status":      "alive",
		"environment": app.config.env,
		"version":     app.config.version,
	}
	js, err := json.Marshal(msg)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

package main

import (
	"html/template"
	"net/http"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	// Initialize a slice containing the paths to the two files.
	// It's important // to note that the file containing our base template must be the *first*
	// file in the slice.
	files := []string{
		"./ui/html/base.html",
		"./ui/html/index.html"}
	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

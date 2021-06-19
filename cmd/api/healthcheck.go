package main

import (
	"net/http"
)

// Declare a handler which writes a plain-text response.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a fixed-format JSON response from a string.
	// js := `{"status":"available", "environment": %q, "version": %q}`
	// js = fmt.Sprintf(js, app.config.env, version)

	// Create a map which holds the information that we want to send in the response.
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	// Pass the map to the json.Marshal() function.
	// This returns a []byte slice containing the encoded JSON.
	// If there was an error, we log it and send the client a generic error message.
	// js, err := json.Marshal(data)
	// if err != nil {
	// 	app.logger.Println(err)
	// 	http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	// 	return
	// }

	// Append a newline to the JSON. This is just a small nicety to make it easier to view in terminal application.
	// js = append(js, '\n')

	// Set the "Content-Type: application/json" header on the response.
	// If you forget to this, Go will default to sending
	// a "Content-Type: text/plain; charset=utf-8" header instead.
	// w.Header().Set("Content-Type", "application/json")

	// Write the JSON as the HTTP response body.
	// w.Write([]byte(js))
}

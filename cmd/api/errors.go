package main

import (
	"fmt"
	"net/http"
)

// The logError() method is a generic helper for logging an error message.
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// For sending JSON-formatted error messages to the client with a given status code.
// Note that we are using an interface{} type for the message parameter, rather than just a string type, as
// this gives us more flexibility over the values that we can include in the response.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	// If this happens to return an error then log it, and fall back to sending the client
	// an empty response with 500 Internal Server Error status code.
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// Used when our application encounters an unexpected problem at runtime.
// It logs the detailed error message, then send a 500 to the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request."
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// Used to send a 405 Method Not Allowed status code and JSON response to the client.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

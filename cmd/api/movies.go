package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/able8/greenlight/internal/data"
)

// For the "POST /v1/movies" endpoint.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be
	// in the HTTP request body. This struct will be our target decode destination.
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use the Decode() method to decode the body contents into the input struct.
	// Note that when we call Decode() we pass a pointer to the inpurt struct as the target decode destination.
	// If there was an error during decoding, we also use our generic errorResponse() helper to
	// send the client a 400 Bad Request response containing the error message.
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Dump the contents of the input struct in a HTTP response.
	fmt.Fprintf(w, "%v\n", input)
}

// For the "GET /v1/movies/:id" endpoints.
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil || id < 1 {
		// http.NotFound(w, r)

		// Use the new notFoundResponse() helper
		app.notFoundResponse(w, r)
		return
	}

	// Create a new instance of the Movie struct.
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	// Create an envelope{"movie": movie} instance.
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		// app.logger.Println(err)
		// http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		app.serverErrorResponse(w, r, err)
	}
}

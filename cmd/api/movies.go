package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/able8/greenlight/internal/data"
	"github.com/able8/greenlight/internal/validator"
)

// For the "POST /v1/movies" endpoint.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be
	// in the HTTP request body. This struct will be our target decode destination.
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"` // Make this field a data.Runtime type.
		Genres  []string     `json:"genres"`
	}

	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use the Decode() method to decode the body contents into the input struct.
	// Note that when we call Decode() we pass a pointer to the inpurt struct as the target decode destination.
	// If there was an error during decoding, we also use our generic errorResponse() helper to
	// send the client a 400 Bad Request response containing the error message.
	// err := json.NewDecoder(r.Body).Decode(&input)

	// Use the new readJSON() helper to decode the request body into the input struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		// app.errorResponse(w, r, http.StatusBadRequest, err.Error())

		// Use the new badRequestResponse() helper.
		app.badRequestResponse(w, r, err)
		return
	}

	// Initialize a new Validator instance.
	v := validator.New()

	// Use the check() method to execute our validation check.
	// This will add the provided key and error message to the errors map if the check fails.
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 characters long")
	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(input.Runtime != 0, "runtime", "must be provided")
	v.Check(input.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")

	// Note that we're using the Unique helper in the line below to check
	// that all values in the input.Genres slice are unique.
	v.Check(validator.Unique(input.Genres), "genres", "must not contain deplicate values")

	// User the Valid() method to see if any of the checks failed.
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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

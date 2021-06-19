package main

import (
	"net/http"

	"github.com/julienschemidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := &httprouter.New()

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using th HandlerFunc() method.
	router.HandleFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandleFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandleFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	// Return the httprouter instance
	return router
}

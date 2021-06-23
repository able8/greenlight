package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function which will always be run in the event of
		// a panic as Go unwinds the stack.
		defer func() {
			// User the builtin recover function to check if there has been
			// a panic or not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header
				// on the response. This acts as a trigger to make Go's HTTP
				// server automatically close the current connection after a sponse has been sent.
				w.Header().Set("Connection", "close")

				// The value returned by recover() has the type interface{},
				// so we use fmt.Errorf() to normalize it into an error.
				// This will log the error using our custom Logger type at the
				// ERROR level and send the client a 500 Internal Server Error response.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// Initialize a new rate limiter which allows an average of
	// 2 requests per second, with a maximum of 4 requests in a single 'burst'.
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call limiter.Allow() to see if the request is permitted, and if it's not,
		// then we call the rateLimitExceededResponse() helper to return a
		// 429 Too Many Requests response (we will create this helper in a minute).
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

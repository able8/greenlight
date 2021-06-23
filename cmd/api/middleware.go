package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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

	// Define a client struct to hold the rate limiter and last seen time for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare a mutex and a map to hold the client's IP addresses and rate limiters.
	var (
		mu sync.Mutex
		// clients = make(map[string]*rate.Limiter)
		// Update the map so the values are pointers to a client struct.
		clients = make(map[string]*client)
	)

	// Launch a background goroutine which removes old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock the mutex to prevent any rate limiter checks from happening while the cleanup is taking place.
			mu.Lock()

			// Loop through all clients. If they haven't been seen within the last three minutes,
			// delete the corresponding entry from the map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			// Inportantly, unlock the mutex when the cleanup is complete.
			mu.Unlock()
		}
	}()

	// Initialize a new rate limiter which allows an average of
	// 2 requests per second, with a maximum of 4 requests in a single 'burst'.
	// limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only carry out the check if rate limiting is enabled.
		if app.config.limiter.enabled {
			// Extract the client's IP address from the request.
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			// Lock the mutex to prevent tis code from being executed concurrently.
			mu.Lock()

			// Check to see if the IP address already exists in the map. If it doesn't,
			// then initialize a new rate limiter adn add the IP address and limiter to the map.
			if _, found := clients[ip]; !found {
				// clients[ip] = rate.NewLimiter(2, 4)

				// Create and add a new client struct to the map if it doesn't already exist.
				// clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}

				// Use the requests-per-second and burst values from the config struct.
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
			}
			// update the last seen time for the client.
			clients[ip].lastSeen = time.Now()

			// Call the Allow() method on the rate limiter for the current IP address.
			// If the request isn't allowed, unlock the mutex and send a 429 response, just like befor.
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			// Very important, unlock the mutex before calling the next handler in the chain.
			// Notice that we DON'T use defer to unlock the mutex, as that
			// would mean that the mutex isn't unlocker until all the handlers donwstream of this
			// middleware have also returned.
			mu.Unlock()

			// Call limiter.Allow() to see if the request is permitted, and if it's not,
			// then we call the rateLimitExceededResponse() helper to return a
			// 429 Too Many Requests response (we will create this helper in a minute).
			// if !limiter.Allow() {
			// 	app.rateLimitExceededResponse(w, r)
			// 	return
			// }

		}
		next.ServeHTTP(w, r)
	})
}

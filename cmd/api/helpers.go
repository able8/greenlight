package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Define an envelope type.
type envelope map[string]interface{}

// Retrieve the "id" URL parameter from the current request context, then
// convert it to an integer and return it. If the operation isn't successful,
// return 0 and an error.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// Define a writeJSON() helper for sending response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	// js, err := json.Marshal(data)

	// Use the json.MarshalIndent() function so that whitespace is added to the encoded JSON.
	// Here we use no line prefix ("") and tab indents ("\t") for each element.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// Note that it's OK if the provided header map is nil.
	// Go doesn't throw an error if you try to range over (or generally, read from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(js))

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Decode the request body into the target destination
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {

		// If there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Use the errors.As() function to check whether the error has the type.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formatted JSON (at character %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formatted JSON")

		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// A json.InvalidUnmarshalError error will be returned if we pass a non-nil pointer
		// to Decode(). We catch this and painc, rather then returning an error to our handler.

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

			// For anything else, return the error message as-is.
		default:
			return err
		}

	}

	return nil
}

package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/able8/greenlight/internal/data"
	"github.com/able8/greenlight/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the data from the request body into a new User struct. Notice also that we set
	// the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default.
	// But setting this explicitly helps to make out intentions clear to anyone reading the code.
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate the user struct and return the error messages to the client if any of the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Launch a goroutine which runs an anonymous function that sends the welcome email.
	go func() {

		// Run a deferred function which uses recover() to catch any panic, and
		// log an error message instead of terminating the application.
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		// panic("test a panic")
		err = app.mailer.Send(user.Email, "user_welcome.tmpl.html", user)
		if err != nil {
			// app.serverErrorResponse(w, r, err)
			// Importantly, if there is an error sending the email then we use the
			// app.logger.PrintError() helper to manage it, instead of the
			// app.serverErrorResponse() helper like before.
			app.logger.PrintError(err, nil)
			return
		}
		app.logger.PrintInfo("Send email successfully", map[string]string{
			"email": user.Email,
		})
	}()
	// // Call the Send() method on our Mailer, passing in the user's email address,
	// // name of the template file, and the User stuct containing the new user's data.
	// err = app.mailer.Send(user.Email, "user_welcome.tmpl.html", user)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }

	// Write a JSON response containing the user data along with a 201 Created status code.
	// err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)

	// Note that we also change this to send the client a 202 Accepted status code.
	// This code indicates that the request has been accepted for processing, but
	// the processing has not been completed.
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

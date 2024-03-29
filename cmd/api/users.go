package main

import (
	"errors"
	"net/http"
	"time"

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

	// Add the "movie:read" permission for the new user.
	err = app.models.Permissions.AddForUser(user.ID, "movies:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Launch a goroutine which runs an anonymous function that sends the welcome email.
	// go func() {

	// 	// Run a deferred function which uses recover() to catch any panic, and
	// 	// log an error message instead of terminating the application.
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			app.logger.PrintError(fmt.Errorf("%s", err), nil)
	// 		}
	// 	}()

	// 	// panic("test a panic")
	// 	err = app.mailer.Send(user.Email, "user_welcome.tmpl.html", user)
	// 	if err != nil {
	// 		// app.serverErrorResponse(w, r, err)
	// 		// Importantly, if there is an error sending the email then we use the
	// 		// app.logger.PrintError() helper to manage it, instead of the
	// 		// app.serverErrorResponse() helper like before.
	// 		app.logger.PrintError(err, nil)
	// 		return
	// 	}
	// 	app.logger.PrintInfo("Send email successfully", map[string]string{
	// 		"email": user.Email,
	// 	})
	// }()

	// After the user record has been created in the databasee, generate a new activation.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Use the background helper to execute an anonymous function that sends the welcome email.
	app.background(func() {
		// As there are now nultiple pieces of data that we want to pass to our email
		// templates, we create a map to act as  a holding structure for the data.
		// This contains the plaintext versino of the activation token for the user,
		// along with their ID.
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		// Send the welcome email, passing in the map above as dynamic data.
		err = app.mailer.Send(user.Email, "user_welcome.tmpl.html", data)
		if err != nil {
			app.logger.PrintError(err, nil)
			return
		}

		app.logger.PrintInfo("Send email successfully", map[string]string{
			"email": user.Email,
		})
	})

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

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body.
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the plaintext token provided by the client.
	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve the details of the user associated with token using the GetForToken() method
	// If no matching recorrd is found, then we let the client know that the token they provided is not valid.
	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Update the user's activation status.
	user.Activated = true

	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// If everything went successfully, then we delete all activation tokens for the user.
	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel.
// We'll add other models to this, like a UserModel and PermissionModel.
type Models struct {
	Movies      MovieModel
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
}

// For ease of use, we also add a New() method which returns a Models struct
// containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
		Users:       UserModel{DB: db}, // Initialize a new UserModel instance.
		Tokens:      TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}

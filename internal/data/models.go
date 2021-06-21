package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// Create a Models struct which wraps the MovieModel.
// We'll add other models to this, like a UserModel and PermissionModel.
type Models struct {
	Movies MovieModel
}

// For ease of use, we also add a New() method which returns a Models struct
// containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}

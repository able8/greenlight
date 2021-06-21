package data

import (
	"database/sql"
	"time"

	"github.com/able8/greenlight/internal/validator"
	"github.com/lib/pq"
)

// Annotate the Movie struct with struct tags to control how the keys appear in the JSON-encoded output.

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Timestamp for when the movie is added to our database
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"` // Movie runtime in minutes
	Genres    []string  `json:"genres,omitempty"`  // Slice of genres for the movie, romance, comedy, etc.
	Version   int32     `json:"version"`           // The version number will be incremented each time the information is updated.
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}

// Insert a new record in the movies table.
// The Insert() method accepts a pointer to a movie struct, which should contain the data for the new record.
func (m MovieModel) Insert(movie *Movie) error {
	// Define the SQL query for inserting a new record in the movies table and
	// returning the system-generated data.
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
	`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately  next to our SQL
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the
	// system-generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Get a specific record from the movies table.
func (m MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

// Update a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {
	return nil
}

// Delete a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	return nil
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	// Use the check() method to execute our validation check.
	// This will add the provided key and error message to the errors map if the check fails.
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 characters long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")

	// Note that we're using the Unique helper in the line below to check
	// that all values in the movie.Genres slice are unique.
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain deplicate values")
}

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately  next to our SQL
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the
	// system-generated id, created_at and version values into the movie struct.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Get a specific record from the movies table.
func (m MovieModel) Get(id int64) (*Movie, error) {
	// The PostgresSQL bigserial type that we're using for the movie ID starts
	// auto-incrementing at 1 by default, so we know that no movies will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shotcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Define the SQL query for retrieving the movie data.
	// Update the query to return pg_sleep(10) as the first value.
	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1
	`

	// Declare a Movie struct to hold the data returned by the query.
	var movie Movie

	// Use the context.WithTimeout() function to create a context.Context which carries a 3-second timeout deadline.
	// Note that we're using the empty context.Background() as the parent context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Importantly, use defer to make sure that we cancel the context before the Get() method returns
	defer cancel()

	// Scan the response data into the fileds of the Movie struct.
	// Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	// err := m.DB.QueryRow(query, id).Scan(

	// User the QueryRowContext() method to execute the query,
	// passing in the context with the deadline as the first argument.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		// Imprtantly, update the Scan() parameters so that the
		// pg_sleep(10) return value is scanned into a byte slice.
		// &[]byte{},
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We chack for this and return our custom ErrRecordNotFound error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Otherwise, return a pointer to the Movie struct.
	return &movie, nil
}

// Update a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {
	// Declare the SQL query for updating the record and returning the new version number.
	query := `
		UPDATE movies
		SET title=$1, year=$2, runtime=$3, genres=$4, version=version+1
		WHERE id=$5 AND version=$6
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version, // Add the expected version.
	}

	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version into the movie struct.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM movies
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// The Exec() method returns a sql.Result object.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get
	// the number of rows affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows are affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it.
	// In that case we return an ErrRecordNotFound error.

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

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

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	// Create a new GetAll() method whch returns a slice of movies.

	// Add an ORDER BY clause and interpolate the sort column and direction. Importantly notice that we also include a secondary sort on the movie ID to ensure a consistent ordering.

	// Update the SQL query to include the window function which
	// counts the total (filtered) records.
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// As our SQL query now has quite a few placeholder parameters, let's collect the values
	// for placeholders in a slice. Notice here how we call the limit() and offset() methods on the Filters struct to get the appropriate values for
	// the LIMIT and OFFSET clause.
	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}
	// This returns a sql.Rows resultset containing the result.
	// Pass the title and genres as the placholder. parameters values.
	// rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(genres))

	// And the pass the args slice to QueryContext() as a variadic parameter.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		// Update this to return an empty Metadata struct
		return nil, Metadata{}, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed before GetAll() returns.
	defer rows.Close()

	totalRecords := 0
	// Initialize an empty slice to hold the movie data.
	movies := []*Movie{}

	for rows.Next() {
		// Initialize an empty Movie struct to hold the dat for an individual movie.
		var movie Movie

		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		// fmt.Printf("%+v\n", movie)
		// Add the Movie struct to the slice.
		movies = append(movies, &movie)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and
	// pagination parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of movies.
	return movies, metadata, nil
}

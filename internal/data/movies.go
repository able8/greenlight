package data

import "time"

type Movie struct {
	ID        int64
	CreatedAt time.Time // Timestamp for when the movie is added to our database
	Title     string
	Year      int32
	Runtime   int32    // Movie runtime in minutes
	Genres    []string // Slice of genres for the movie, romance, comedy, etc.
	Version   int32    // The version number will be incremented each time the information is updated.
}

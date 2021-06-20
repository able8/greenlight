package data

import (
	"fmt"
	"strconv"
)

// Declare a custom Runtime type.
type Runtime int32

// Implement a MarshalJSON() method on the Runtime type so that it
// statisfies the json.Marshaler interface.
func (r Runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format.
	jsonValue := fmt.Sprintf("%d mins", r)

	// Wrap it in double quotes. It needs to be surrounded by
	// double quotes in order to be a valid JSON string.
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

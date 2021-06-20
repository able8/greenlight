package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Define an error that out UnmarshalJSON() method can return if we are unable to
// parse or convert the JSON string successfully.
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

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

// Note that because UnmarshalJSON() needs to modify the receiver (our Runtime type), we must
// use a pointer receiver for this to work correctly. Otherwise, we will only be modifying
// a copy (which is then discarded when this method is returned).
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// We expect that the incoming JSON value will be a string in the format at
	// "<runtime> mins", and the first thing we need to do is remove the surrounding
	// double quotes from the string.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))

	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedJSONValue, " ")

	// Sanity check the parts of the string to make sure it was in the exected format.
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Otherwise, parse the string containing the number into adn int32.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Convert the int32 to a Runtime type and assign this to the receiver.
	// Note that we use the * operator to deference the receiver (which is a pointer to
	// a Runtime type) in order to set the underlying value of the pointer.
	*r = Runtime(i)

	return nil
}

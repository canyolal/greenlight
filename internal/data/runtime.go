package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// Custom runtime type has underlying type int32
type Runtime int32

// MarshalJSON() method for runtime type which satisfies Marshaler interface so that
// runtime can return in "r mins" format in JSON.
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonVal := fmt.Sprintf("%d mins", r)

	// Use the strconv.Quote() function on the string to wrap it in double quotes. It
	// needs to be surrounded by double quotes in order to be a valid *JSON string*.
	quotedJsonVal := strconv.Quote(jsonVal)

	return []byte(quotedJsonVal), nil

}

// UnmarshalJSON method that satisfies Unmarshaler interface which custom checks for
// Runtime field in JSON.
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

	// We expect that the incoming JSON value will be a string in the format
	// "<runtime> mins", and the first thing we need to do is remove the surrounding
	// double-quotes from this string. If we can't unquote it, then we return the
	// ErrInvalidRuntimeFormat error.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedJSONValue, " ")

	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidRuntimeFormat error again.
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Otherwise, parse the string containing the number into an int32. Again, if this
	// fails return the ErrInvalidRuntimeFormat error.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Convert the int32 to a Runtime type and assign this to the receiver. Note that we use
	// use the * operator to deference the receiver (which is a pointer to a Runtime
	// type) in order to set the underlying value of the pointer.
	*r = Runtime(i)

	return nil
}

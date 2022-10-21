package data

import (
	"fmt"
	"strconv"
)

// Custom runtime type has underlying type int32
type runtime int32

// MarshalJSON() method for runtime type which satisfies Marshaler interface so that
// runtime can return in "r mins" format in JSON.
func (r *runtime) MarshalJSON() ([]byte, error) {
	jsonVal := fmt.Sprintf("%d mins", r)

	// Use the strconv.Quote() function on the string to wrap it in double quotes. It
	// needs to be surrounded by double quotes in order to be a valid *JSON string*.
	quotedJsonVal := strconv.Quote(jsonVal)

	return []byte(quotedJsonVal), nil

}

package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrJSONInvalid is a named application error.
var ErrJSONInvalid = errors.New("invalid JSON")

func unmarshalJSON(d string) (data interface{}, err error) {
	if err = json.Unmarshal([]byte(d), &data); err != nil {
		// Wrap encoding/json error with an error identifiable on the application level.
		return nil, WrapWithMessage(err, ErrJSONInvalid, "invalid JSON (%q)", d)
	}
	return data, nil
}

func ExampleWrapWithMessage() {
	_, err := unmarshalJSON(`{`)
	fmt.Println(err)

	appErr := Fetch(err, ErrJSONInvalid)
	fmt.Println(appErr == ErrJSONInvalid)
	fmt.Println(appErr)

	// Output:
	// invalid JSON ("{") : invalid JSON : unexpected end of JSON input
	// true
	// invalid JSON
}

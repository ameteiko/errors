package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrJSONUnmarshal is a named application error.
var ErrJSONUnmarshal = errors.New("json unmarshal error")

func unmarshalUser(ud string) (user interface{}, err error) {
	if err = json.Unmarshal([]byte(ud), &user); err != nil {
		// Wrap encoding/json error with an error identifiable on the application level.
		return nil, Wrap(err, ErrJSONUnmarshal)
	}
	return user, nil
}

func ExampleWrap() {
	_, err := unmarshalUser(`{`)
	fmt.Println(err)

	appErr := Fetch(err, ErrJSONUnmarshal)
	fmt.Println(appErr == ErrJSONUnmarshal)
	fmt.Println(appErr)

	// Output:
	// json unmarshal error : unexpected end of JSON input
	// true
	// json unmarshal error
}

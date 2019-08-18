package errors

import (
	"errors"
	"fmt"
)

var (
	ErrTooShort = errors.New("string is too short")
	ErrTooLong  = errors.New("string is too long")
)

func validateUsername(un string) error {
	if len(un) > 20 {
		return WithMessage(ErrTooLong, "username len exceeds the limit of 20 chars (%q)", un)
	}

	if len(un) < 3 {
		return WithMessage(ErrTooShort, "username len is less than 3 chars (%q)", un)
	}

	return nil
}

func ExampleWithMessage() {
	err := validateUsername("Le")
	fmt.Println(err)
	fmt.Println(Fetch(err, ErrTooShort))

	// Output:
	// username len is less than 3 chars ("Le") : string is too short
	// string is too short
}

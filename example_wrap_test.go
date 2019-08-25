package errors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoSpecialCharacters = errors.New("no special characters")
	ErrStringTooShort      = errors.New("string is too short")
	ErrStringTooLong       = errors.New("string is too long")
)

func validatePassword(pwd string) (err error) {
	if !strings.ContainsAny(pwd, `~!@#$%^&*()_+{}<>?`) {
		err = Wrap(err, ErrNoSpecialCharacters)
	}

	if len(pwd) < 3 {
		err = Wrap(err, ErrStringTooShort)
	}

	if len(pwd) > 12 {
		err = Wrap(err, ErrStringTooLong)
	}

	return err
}

func ExampleFetch() {
	err := validatePassword("")
	fmt.Println(err)

	tooShortErr := Fetch(err, ErrStringTooShort)
	fmt.Println(tooShortErr == ErrStringTooShort)

	tooLongErr := Fetch(err, ErrStringTooLong)
	fmt.Println(tooLongErr)

	noSpecialCharactersErr := Fetch(err, ErrNoSpecialCharacters)
	fmt.Println(noSpecialCharactersErr == ErrNoSpecialCharacters)

	// Output:
	// string is too short : no special characters
	// true
	// <nil>
	// true
}

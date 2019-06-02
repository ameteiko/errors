// Package errors provides the typed errors checking functionality.
// TODO: Write package comment here.
package errors

import (
	"fmt"
	"reflect"
)

// Wrap wraps variadic number of errors into an error chain.
// TODO: test appending another enchainer.
func Wrap(errs ...error) error {
	var (
		enchainer          *chain
		doesEnchainerExist bool
		enchainerIndex     int
	)

	// Find enchainer instance if exists.
	for i := range errs {
		enchainerIndex = len(errs) - i - 1
		err := errs[enchainerIndex]
		enchainer, doesEnchainerExist = err.(*chain)
		if doesEnchainerExist {
			break
		}
	}

	if !doesEnchainerExist {
		enchainer = newChain()
		for _, err := range errs {
			enchainer.append(err)
		}

		return enchainer
	}

	// prepend errors
	for i := range errs[:enchainerIndex] {
		err := errs[enchainerIndex-i-1]
		anotherEnchainer, isAnotherEnchainer := err.(*chain)
		if !isAnotherEnchainer {
			enchainer.prepend(err)
			continue
		}

		for _, innerEnchainerError := range anotherEnchainer.getErrors() {
			enchainer.prepend(innerEnchainerError)
		}
	}

	// append errors
	for _, err := range errs[enchainerIndex+1:] {
		enchainer.append(err)
	}

	return enchainer
}

func WrapWithMessage(err1, err2 error, format string, args ...interface{}) error {
	return nil
}

// Cause returns a first error that implements causer.
// causer is always must be a pointer to either a structure or interface.
func Cause(err error, causerObj interface{}) error {
	// Get causer type element.
	causerType := reflect.TypeOf(causerObj)
	if nil == causerType || reflect.Ptr != causerType.Kind() {
		return nil
	}

	causer := causerType.Elem()
	chainer, ok := err.(*chain)
	if !ok {
		// The error object is not a Chainer instance.
		errType := reflect.TypeOf(err)
		if doesErrorMatch(errType, causerType, causer) {
			return err
		}

		return nil
	}

	for _, e := range chainer.getErrors() {
		errType := reflect.TypeOf(e)
		if doesErrorMatch(errType, causerType, causer) {
			return e
		}
	}

	return nil
}

// WithMessage returns an error wrapped with message.
func WithMessage(err error, format string, args ...interface{}) error {
	if format != "" {
		return Wrap(err, fmt.Errorf(format, args...))
	}

	if chain, ok := err.(*chain); ok {
		return chain
	}

	return Wrap(err)
}

// doesErrorMatch returns true if err object implements or is assignable to causer.
func doesErrorMatch(err reflect.Type, causerType reflect.Type, causerElem reflect.Type) bool {
	doesImplement := reflect.Interface == causerElem.Kind() && err.Implements(causerElem)
	isTypeOf := reflect.Struct == causerElem.Kind() && (err.AssignableTo(causerType) || err.AssignableTo(causerElem))

	if doesImplement || isTypeOf {
		return true
	}

	return false
}

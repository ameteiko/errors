package errors

import (
	"errors"
	"fmt"
	"reflect"
)

//
// Wrap wraps errors list into an error chain.
// TODO: if there are two enchainers, raise an error.
// TODO: test appending another enchainer.
//
func Wrap(errs ...error) Chainer {
	var (
		enchainer          Chainer
		doesEnchainerExist bool
		enchainerIndex     int
	)

	// Find enchainer instance if exists.
	for i := range errs {
		enchainerIndex = len(errs) - i - 1
		err := errs[enchainerIndex]
		enchainer, doesEnchainerExist = err.(Chainer)
		if doesEnchainerExist {
			break
		}
	}

	if !doesEnchainerExist {
		enchainer = NewChain()
		for _, err := range errs {
			enchainer.Append(err)
		}

		return enchainer
	}

	// Prepend errors
	for i := range errs[:enchainerIndex] {
		err := errs[enchainerIndex-i-1]
		anotherEnchainer, isAnotherEnchainer := err.(Chainer)
		if !isAnotherEnchainer {
			enchainer.Prepend(err)
			continue
		}

		for _, innerEnchainerError := range anotherEnchainer.GetErrors() {
			enchainer.Prepend(innerEnchainerError)
		}
	}

	// Append errors
	for _, err := range errs[enchainerIndex+1:] {
		enchainer.Append(err)
	}

	return enchainer
}

//
// Cause returns a first error that implements causer.
// causer is always must be a pointer to either a structure or interface.
//
func Cause(err error, causerObj interface{}) error {
	// Get causer type element.
	causerType := reflect.TypeOf(causerObj)
	if nil == causerType || reflect.Ptr != causerType.Kind() {
		return nil
	}

	causer := causerType.Elem()
	chainer, ok := err.(Chainer)
	if !ok {
		// The error object is not a Chainer instance.
		errType := reflect.TypeOf(err)
		if doesErrorMatch(errType, causerType, causer) {
			return err
		} else {
			return nil
		}
	}

	for _, e := range chainer.GetErrors() {
		errType := reflect.TypeOf(e)
		if doesErrorMatch(errType, causerType, causer) {
			return e
		}
	}

	return nil
}

//
// WithMessage returns an error wrapped with message.
//
func WithMessage(err error, format string, args ...interface{}) error {
	if "" != format {
		return Wrap(err, errors.New(fmt.Sprintf(format, args...)))
	}

	if chain, ok := err.(Chainer); ok {
		return chain
	}

	return Wrap(err)
}

//
// doesErrorMatch returns true if err object implements or is assignable to causer.
//
func doesErrorMatch(err reflect.Type, causerType reflect.Type, causerElem reflect.Type) bool {
	doesImplement := reflect.Interface == causerElem.Kind() && err.Implements(causerElem)
	isTypeOf := reflect.Struct == causerElem.Kind() && (err.AssignableTo(causerType) || err.AssignableTo(causerElem))

	if doesImplement || isTypeOf {
		return true
	}

	return false
}

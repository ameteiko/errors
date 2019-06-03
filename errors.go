// Package errors provides the typed errors checking functionality.
// TODO: Write package comment here.
package errors

import (
	"fmt"
	"reflect"
)

// Wrap wraps variadic number of errors into an error errChain.
func Wrap(errs ...error) error {
	var (
		enchainer          *errChain
		doesEnchainerExist bool
		enchainerIndex     int
	)

	// Find enchainer instance if exists.
	for i := range errs {
		enchainerIndex = len(errs) - i - 1
		err := errs[enchainerIndex]
		enchainer, doesEnchainerExist = err.(*errChain)
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
		anotherEnchainer, isAnotherEnchainer := err.(*errChain)
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

// To returns a first error that implements causer.
// causer is always must be a pointer to either a structure or interface.

// chainErr is supposed to be a errChain instance, but it can also be any error.
// targetErr is the error to be searched over the err. It can be: pointer to the error type, interface, or error value.
func To(chainedErr error, targetErr interface{}) error {
	if chainedErr == nil || targetErr == nil {
		return nil
	}

	chain, ok := chainedErr.(*errChain)
	if !ok {
		// The chainedErr is not an errChain instance. Check if it matches targetErr itself.
		if errorMatches(chainedErr, targetErr) {
			return chainedErr
		}
		return nil
	}

	for _, e := range chain.getErrors() {
		if errorMatches(e, targetErr) {
			return e
		}
	}

	return nil
}

// Tos returns the list of matching errors.
func Tos(chainedErr error, targetErr interface{}) (errs []error) {
	if chainedErr == nil || targetErr == nil {
		return nil
	}

	chain, ok := chainedErr.(*errChain)
	if !ok {
		// The chainedErr is not an errChain instance. Check if it matches targetErr itself.
		if errorMatches(chainedErr, targetErr) {
			return []error{chainedErr}
		}
		return nil
	}

	for _, e := range chain.getErrors() {
		if errorMatches(e, targetErr) {
			errs = append(errs, e)
		}
	}

	return errs
}

// WithMessage returns an error wrapped with message.
func WithMessage(err error, format string, args ...interface{}) error {
	if format != "" {
		return Wrap(err, fmt.Errorf(format, args...))
	}

	if chain, ok := err.(*errChain); ok {
		return chain
	}

	return Wrap(err)
}

// WrapWithMessage wraps two errors with message.
func WrapWithMessage(err1, err2 error, format string, args ...interface{}) error {
	return WithMessage(Wrap(err1, err2), format, args...)
}

// errorMatches returns true if targetErr matches sourceErr.
// target parameter options:
//     - errorInstance
//     - (*errorInterface)(nil)
//     - (*customErrorStruct)(nil)
// TODO: apply some optimization here not to calculate targetType, targetElem for each iteration.
func errorMatches(sourceErr error, target interface{}) bool {
	// Invariant verification.
	if sourceErr == nil || target == nil {
		return false
	}

	// Check if target is of type error.
	if targetErr, ok := target.(error); ok && sourceErr == targetErr {
		return true
	}

	// target must be a nil pointer to either interface or an error struct.
	targetType := reflect.TypeOf(target)
	targetElem := targetType.Elem()
	sourceType := reflect.TypeOf(sourceErr)
	doesImplement := reflect.Interface == targetElem.Kind() && sourceType.Implements(targetElem)
	isTypeOf := reflect.Struct == targetElem.Kind() &&
		(sourceType.AssignableTo(targetType) || sourceType.AssignableTo(targetElem))
	if doesImplement || isTypeOf {
		return true
	}

	return false
}

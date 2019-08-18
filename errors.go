// Package errors provides error context wrappers to for handling errors in Go.
//
// The principle of handling errors is to either handle the error when it happens or to delegate its handling to upper
// layer via providing current context. The purpose of this package is to provide a way to attach some additional
// context to the error.
//
// The package's core is the internal type errors.queue which is exposed to clients as error interface. All
// application error handling is supposed to operate with errors using Wrap(), WithMessage() or WrapWithMessage()
// providing the existing error and attaching some additional context to it. Each of these methods either create a new
// error queue instance if it doesn't exist or uses existing one.
//
// Once the error reaches to layer to handle the error, one of functions Fetch(), FetchByType() should be used to
// examine the error on having one of contexts.
package errors

import (
	"errors"
	"fmt"
	"reflect"
)

// New returns error object wrapped with a stacktrace.
func New(msg string) error {
	return Wrap(errors.New(msg))
}

// isErrNil returns true if error is nil.
func isErrNil(err error) bool {
	val := reflect.ValueOf(err)
	if err == nil || (val.Kind() == reflect.Ptr && val.IsNil()) {
		return true
	}
	return false
}

// Wrap wraps variadic number of errors into an error queue providing context about the error.
//
// This function wraps several errors into internal queue object that satisfies the error interface and is inspectable
// by Fetch(), FetchByType() functions.
// It is used to provide extended local context to the existing error object.
// If several queue instances are passed, then the they will be disintegrated and all the errors from them will be
// injected into the last found queue object.
// TODO: think on removing the side-effects of changing the queue objects.
func Wrap(errs ...error) error {
	var (
		q    *queue
		qIdx int
		ok   bool
	)

	// Ignore all nil errors from the list.
	ee := make([]error, 0, len(errs))
	for i := 0; i < len(errs); i++ {
		if !isErrNil(errs[i]) {
			ee = append(ee, errs[i])
		}
	}
	errs = ee

	if len(errs) == 0 {
		return nil
	}

	// Find an existing queue instance if exists. Iterating in reverse order.
	for i := range errs {
		qIdx = len(errs) - i - 1
		err := errs[qIdx]
		if q, ok = err.(*queue); ok {
			break
		}
	}

	// There is no error queue, so create one and inject all errors into it.
	if q == nil {
		q = newQueue()
		q.errs = errs
		return q
	}

	// Prepend errors before qIdx in reverse order.
	for i := range errs[:qIdx] {
		err := errs[qIdx-i-1]
		q.prepend(err)
	}

	// Append errors after the queue instance in direct order.
	q.errs = append(q.errs, errs[qIdx+1:]...)

	return q
}

// WithMessage returns an error wrapped with message.
// This function is used to attach some context in a form of formatted to existing error.
func WithMessage(err error, format string, args ...interface{}) error {
	if format == "" || isErrNil(err) {
		return err
	}

	return Wrap(err, fmt.Errorf(format, args...))
}

// WrapWithMessage wraps two errors with message.
// It is used when an external call returns a 3rd-party error, that needs to be wrapped not only into an application
// error but with additional context too.
func WrapWithMessage(err1, err2 error, format string, args ...interface{}) error {
	return WithMessage(Wrap(err1, err2), format, args...)
}

// validateErrors returns flags whether errors are valid (meaning not nil) and match each other.
func validateErrors(sourceErr, targetErr error) (valid bool, matched bool) {
	isSourceNil := isErrNil(sourceErr)
	isTargetNil := isErrNil(targetErr)
	if isSourceNil && isTargetNil {
		return false, true
	}
	if isSourceNil || isTargetNil {
		return false, false
	}

	if sourceErr == targetErr || sourceErr.Error() == targetErr.Error() {
		return true, true
	}

	return true, false
}

// Fetch returns targetErr from the err queue.
// Provided that qErr is an error queue, this function iterates over all queue errors and returns the first matched one.
func Fetch(qErr, targetErr error) error {
	valid, matched := validateErrors(qErr, targetErr)
	if !valid {
		return nil
	} else if matched {
		return targetErr
	}

	q, ok := qErr.(*queue)
	if !ok {
		return nil
	}

	for _, err := range q.getErrors() {
		if _, matched := validateErrors(err, targetErr); matched {
			return targetErr
		}
	}

	return nil
}

// FetchByType returns a first error from the error queue that implements or assignable to targetErr.
// targetErr must be a pointer to either a structure or interface.
//
// qErr is supposed to be a queue instance, but it can also be any error object.
// targetErr is the error to be searched over the err. It can be: pointer to the error type, interface, or error value.
func FetchByType(qErr error, targetErr interface{}) error {
	if isErrNil(qErr) {
		return nil
	}

	if targetErr == nil {
		return nil
	}

	q, ok := qErr.(*queue)
	if !ok {
		// The qErr is not a queue instance but error, so check if it matches target error.
		if errorMatches(qErr, targetErr) {
			return qErr
		}
		return nil
	}

	for _, e := range q.getErrors() {
		if errorMatches(e, targetErr) {
			return e
		}
	}

	return nil
}

// errorMatches returns true if targetErr matches sourceErr.
// target parameter options:
//     - errorInstance
//     - (*errorInterface)(nil)
//     - (*customErrorStruct)(nil)
// TODO: apply some optimization here not to calculate targetType, targetElem for each iteration.
// nolint:gocyclo
func errorMatches(sourceErr error, target interface{}) bool {
	sourceType, _, sourceConversionErr := getTypeElem(sourceErr)
	targetType, targetElem, targetConversionErr := getTypeElem(target)
	if sourceConversionErr != nil || targetConversionErr != nil {
		return false
	}

	doesImplement := reflect.Interface == targetElem.Kind() && sourceType.Implements(targetElem)
	isTypeOf := reflect.Struct == targetElem.Kind() &&
		(sourceType.AssignableTo(targetType) || sourceType.AssignableTo(targetElem))
	if doesImplement || isTypeOf {
		return true
	}

	return false
}

func getTypeElem(obj interface{}) (tp reflect.Type, el reflect.Type, err error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, nil, errors.New("obj is nil")
	}
	objTypeKind := objType.Kind()
	if objTypeKind != reflect.Struct && objTypeKind != reflect.Ptr && objTypeKind != reflect.Interface {
		return nil, nil, errors.New("must be valid obj")
	}

	if objTypeKind == reflect.Struct {
		return reflect.New(objType).Type(), objType, nil
	}

	return objType, objType.Elem(), nil
}

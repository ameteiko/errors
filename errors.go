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
	"fmt"
	"reflect"
)

// New returns error object wrapped with a stacktrace.
func New(msg string) error {
	return Wrap(fmt.Errorf(msg))
}

// Wrap wraps variadic number of errors into an error queue providing context about the error.
//
// This function wraps several errors into internal queue object that satisfies the error interface and is inspectable
// by Fetch(), FetchByType() functions.
// It is used to provide extended local context to the existing error object.
// If several queue instances are passed, then the they will be disintegrated and all the errors from them will be
// injected into the last found queue object.
func Wrap(errs ...error) error {
	var ee []error
	for _, err := range errs {
		if isErrNil(err) {
			continue
		}

		if q, ok := err.(*queue); ok {
			ee = append(ee, q.errs...)
			continue
		}

		ee = append(ee, err)
	}

	if len(ee) == 0 {
		return nil
	}

	return &queue{errs: ee}
}

// WithMessage returns an error wrapped with message.
// This function is used to attach some context in a form of formatted text to an existing error.
func WithMessage(err error, format string, args ...interface{}) error {
	if isErrNil(err) {
		return nil
	}

	if format == "" {
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

// Fetch returns targetErr from the err queue.
// Provided that qErr is an error queue, this function iterates over all queue errors and returns the first matched one.
func Fetch(qErr, targetErr error) error {
	if isErrNil(qErr) || isErrNil(targetErr) {
		return nil
	}

	q, ok := qErr.(*queue)
	if !ok {
		if compareErrs(qErr, targetErr) {
			return targetErr
		}
		return nil
	}

	for _, err := range q.getErrors() {
		if compareErrs(err, targetErr) {
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
// TODO: think on errors that are created by values, not by by references.
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

// isErrNil returns true if error object is nil.
func isErrNil(err error) bool {
	if err == nil {
		return true
	}
	if val := reflect.ValueOf(err); val.Kind() == reflect.Ptr && val.IsNil() {
		return true
	}

	return false
}

// compareErrs returns true if errors are the same.
// Method Contract: sourceErr and targetErr are not nil.
func compareErrs(sourceErr, targetErr error) bool {
	return sourceErr.Error() == targetErr.Error()
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
		return nil, nil, fmt.Errorf("obj is nil")
	}
	objTypeKind := objType.Kind()
	if objTypeKind != reflect.Struct && objTypeKind != reflect.Ptr && objTypeKind != reflect.Interface {
		return nil, nil, fmt.Errorf("must be valid obj")
	}

	if objTypeKind == reflect.Struct {
		return reflect.New(objType).Type(), objType, nil
	}

	return objType, objType.Elem(), nil
}

// Package errors provides error context wrappers to for handling errors in Go.
//
// The principle of handling errors is to either handle the error when it happens or to delegate its handling to upper
// layer via providing current context. The purpose of this package is to provide a way to attach some additional
// context to the error.
//
// The package's core is the internal type errors.queue which is returned to clients as error interface. All
// application error handling is supposed to operate with errors using Wrap(), WithMessage() or WrapWithMessage()
// providing the existing error and attaching some additional context to it. Each of those functions creates a new
// error queue instance.
//
// Once the error reaches to layer to handle the error, one of functions Fetch(), FetchByType() or FetchAllByType()
// should be used to examine the error on having one of contexts.
package errors

import (
	"fmt"
	"reflect"
)

// New returns an error.
// This method is a replacement for built-in errors.New function.
func New(format string, args ...interface{}) error {
	if format == "" {
		return nil
	}

	return fmt.Errorf(format, args...)
}

// Wrap wraps variadic number of errors into an error queue providing additional error context.
//
// This function wraps several errors into internal queue object that conforms to the error interface and is inspectable
// by Fetch(), FetchByType(), FetchAllByType() functions. It is used to provide extended local context to the existing
// error object. If any of errors is a queue instance, then Wrap() fetches it's error list and puts to the newly created
// queue instance.
func Wrap(errs ...error) error {
	var (
		q                = newQueue()
		foundQs          int
		foundQStacktrace fmt.Formatter
	)
	for _, err := range errs {
		if isErrNil(err) {
			continue
		}

		errQ, ok := err.(*queue)
		if !ok {
			q.errs = append(q.errs, err)
			continue
		}

		q.errs = append(q.errs, errQ.errs...)
		if foundQs == 0 {
			foundQStacktrace = errQ.stacktrace
		}
		foundQs++
	}

	if len(q.errs) == 0 {
		return nil
	}

	if foundQs == 1 {
		q.stacktrace = foundQStacktrace
	}

	return q
}

// WithMessage returns an error wrapped with message.
// This function is used to attach some context in a form of formatted text to an existing error.
func WithMessage(err error, format string, args ...interface{}) error {
	if isErrNil(err) {
		return nil
	}

	return Wrap(err, New(format, args...))
}

// WrapWithMessage wraps two errors with message.
// It is used when an external call returns a 3rd-party error, that needs to be wrapped not only into an application
// error but with additional context too.
func WrapWithMessage(err1, err2 error, format string, args ...interface{}) error {
	if isErrNil(err1) && isErrNil(err2) {
		return nil
	}

	return Wrap(err1, err2, New(format, args...))
}

// Fetch returns targetErr from the err queue.
// Provided that qErr is an error queue, this function iterates over all queue errors and returns the first matched one.
func Fetch(qErr, targetErr error) error {
	if isErrNil(qErr) || isErrNil(targetErr) {
		return nil
	}

	var errs []error
	if q, ok := qErr.(*queue); ok {
		errs = q.getErrors()
	} else {
		errs = []error{qErr}
	}

	for _, err := range errs {
		if compareErrs(err, targetErr) {
			return targetErr
		}
	}

	return nil
}

// FetchByType returns a first error from the error queue that implements or assignable to targetErr.
// targetErr must be a pointer to either a structure or interface.
// qErr is supposed to be a queue instance, but it can also be any error object.
func FetchByType(qErr error, targetErr interface{}) error {
	errs := fetchAllByType(qErr, targetErr, true)
	if errs == nil {
		return nil
	}

	return errs[0]
}

// FetchAllByType returns all matched errors from the error queue that implement or are assignable to targetErr.
func FetchAllByType(qErr error, targetErr interface{}) []error {
	return fetchAllByType(qErr, targetErr, false)
}

func fetchAllByType(qErr error, targetErr interface{}, returnFirst bool) (errs []error) {
	targetType, targetElem, err := getTypeElem(targetErr)
	if isErrNil(qErr) || err != nil || targetType.Kind() != reflect.Ptr {
		return nil
	}

	var qErrs []error
	if q, ok := qErr.(*queue); ok {
		qErrs = q.getErrors()
	} else {
		qErrs = []error{qErr}
	}

	for _, e := range qErrs {
		if !errorMatches(e, targetType, targetElem) {
			continue
		}

		if returnFirst {
			return []error{e}
		}

		errs = append(errs, e)
	}

	if len(errs) > 0 {
		return errs
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
// Method contract: sourceErr and targetErr are not nil.
func compareErrs(sourceErr, targetErr error) bool {
	return sourceErr.Error() == targetErr.Error()
}

// errorMatches returns true if targetErr matches sourceErr.
// target parameter options:
//     - errorInstance
//     - (*errorInterface)(nil)
//     - (*customErrorStruct)(nil)
func errorMatches(sourceErr error, targetType reflect.Type, targetElem reflect.Type) bool {
	sourceType, _, err := getTypeElem(sourceErr)
	if err != nil {
		return false
	}

	doesImplement := reflect.Interface == targetElem.Kind() && sourceType.Implements(targetElem)
	isTypeOf := reflect.Struct == targetElem.Kind() &&
		(sourceType.AssignableTo(targetType) || sourceType.AssignableTo(targetElem))

	return doesImplement || isTypeOf
}

func getTypeElem(obj interface{}) (tp reflect.Type, el reflect.Type, err error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, nil, New("obj is nil")
	}
	objTypeKind := objType.Kind()
	if objTypeKind != reflect.Struct && objTypeKind != reflect.Ptr && objTypeKind != reflect.Interface {
		return nil, nil, New("obj must be one of: struct, pointer, interface")
	}

	if objTypeKind == reflect.Struct {
		return reflect.New(objType).Type(), objType, nil
	}

	return objType, objType.Elem(), nil
}

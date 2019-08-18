package errors

import (
	"errors"
	"fmt"
	"testing"
)

// typedError and typedErrorInterface represent a custom error type for a test suite.
type (
	typedError          struct{ msg string }
	typedErrorPtr       struct{ msg string }
	typedErrorInterface interface{ Message() string }
	customInterface     interface{ Send() string }
)

func (e typedError) Error() string     { return e.msg }
func (e typedError) Message() string   { return e.msg }
func (e *typedErrorPtr) Error() string { return e.msg }

func TestValidateErrors(t *testing.T) {
	var (
		err1               = errors.New("1")
		err2               = errors.New("2")
		nilTypedErr        *typedErrorPtr
		anotherNilTypedErr *typedErrorPtr
	)

	tcs := []struct {
		name                 string
		sourceErr, targetErr error
		valid, matched       bool
	}{
		{
			name:      "ForBothNilErrors",
			sourceErr: nil, targetErr: nil,
			valid: false, matched: true,
		},
		{
			name:      "ForNilSourceError",
			sourceErr: nil, targetErr: err1,
			valid: false, matched: false,
		},
		{
			name:      "ForNilTargetError",
			sourceErr: err1, targetErr: nil,
			valid: false, matched: false,
		},
		{
			name:      "ForMatchingErrors",
			sourceErr: err1, targetErr: err1,
			valid: true, matched: true,
		},
		{
			name:      "ForNotMatchingErrors",
			sourceErr: err2, targetErr: err1,
			valid: true, matched: false,
		},
		{
			name:      "ForAnErrorAndAMatchingErrorQueueAsSource",
			sourceErr: &queue{errs: []error{err1}}, targetErr: err1,
			valid: true, matched: true,
		},
		{
			name:      "ForAnErrorAndAMatchingErrorQueueAsTarget",
			sourceErr: err1, targetErr: &queue{errs: []error{err1}},
			valid: true, matched: true,
		},
		{
			name:      "ForAnErrorAndANotMatchingErrorQueueAsSource",
			sourceErr: &queue{errs: []error{err1}}, targetErr: err2,
			valid: true, matched: false,
		},
		{
			name:      "ForANilErrorByReference",
			sourceErr: nilTypedErr, targetErr: anotherNilTypedErr,
			valid: false, matched: true,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			valid, match := validateErrors(tc.sourceErr, tc.targetErr)

			if tc.valid != valid {
				t.Errorf("validateErrors(%v, %v).valid != %v, got %v", tc.sourceErr, tc.targetErr, tc.valid, valid)
			}

			if tc.matched != match {
				t.Errorf("validateErrors(%v, %v).match != %v, got %v", tc.sourceErr, tc.targetErr, tc.matched, match)
			}
		})
	}
}

func TestFetch(t *testing.T) {
	var (
		err1 = errors.New("1")
		err2 = errors.New("2")
		err3 = errors.New("3")
	)

	tcs := []struct {
		name            string
		qErr, targetErr error
		fetchedErr      error
	}{
		{
			name: "ForBothNils",
			qErr: nil, targetErr: nil,
			fetchedErr: nil,
		},
		{
			name: "ForNilErrorQueue",
			qErr: nil, targetErr: err1,
			fetchedErr: nil,
		},
		{
			name: "ForNilTargetError",
			qErr: err1, targetErr: nil,
			fetchedErr: nil,
		},
		{
			name: "ForMatchingErrors",
			qErr: err1, targetErr: err1,
			fetchedErr: err1,
		},
		{
			name: "ForNotMatchingErrors",
			qErr: err2, targetErr: err1,
			fetchedErr: nil,
		},
		{
			name: "ForAnEmptyErrorQueue",
			qErr: &queue{}, targetErr: err1,
			fetchedErr: nil,
		},
		{
			name: "ForAnErrorQueueWithNotMatchingErrors",
			qErr: &queue{errs: []error{err1, err2}}, targetErr: err3,
			fetchedErr: nil,
		},
		{
			name: "ForAnErrorQueueWithError",
			qErr: &queue{errs: []error{err1, err2}}, targetErr: err2,
			fetchedErr: err2,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			fetchedErr := Fetch(tc.qErr, tc.targetErr)

			if fetchedErr == nil && tc.fetchedErr == nil {
				return
			}

			if fetchedErr == nil || tc.fetchedErr == nil {
				t.Errorf("Fetch(%v, %v) != %v, got %v", tc.qErr, tc.targetErr, tc.fetchedErr, fetchedErr)
			}

			if tc.fetchedErr.Error() != fetchedErr.Error() {
				t.Errorf("Fetch(%v, %v) != %v, got %v", tc.qErr, tc.targetErr, tc.fetchedErr, fetchedErr)
			}
		})
	}
}

func TestErrorMatches(t *testing.T) {
	var (
		err1 = errors.New("1")
	)

	tcs := []struct {
		name      string
		sourceErr error
		targetErr interface{}
		matched   bool
	}{
		{
			name:      "ForTwoSameErrors",
			sourceErr: err1, targetErr: err1,
			matched: true,
		},
		{
			name:      "ForDifferentErrorsWithSameValue",
			sourceErr: errors.New("error"), targetErr: errors.New("error"),
			matched: true,
		},
		{
			name:      "ForCustomInterfacePointer",
			sourceErr: err1, targetErr: (*customInterface)(nil),
			matched: false,
		},
		{
			name:      "ForErrorInterfacePointer",
			sourceErr: err1, targetErr: (*error)(nil),
			matched: true,
		},
		{
			name:      "ForACustomErrorAndErrorInterfacePointer",
			sourceErr: typedError{msg: "error"}, targetErr: (*error)(nil),
			matched: true,
		},
		{
			name:      "ForACustomErrorAndCustomInterfacePointer",
			sourceErr: typedError{msg: "error"}, targetErr: (*typedErrorInterface)(nil),
			matched: true,
		},
		{
			name:      "ForACustomErrorAndCustomInterface",
			sourceErr: typedError{msg: "error"}, targetErr: (typedErrorInterface)(nil),
			matched: false,
		},
		{
			name:      "ForACustomErrorAndAPointerToErrorStruct",
			sourceErr: typedError{msg: "error"}, targetErr: (*typedError)(nil),
			matched: true,
		},
		{
			name:      "ForACustomErrorPointerAndAPointerToErrorStruct",
			sourceErr: &typedError{msg: "error"}, targetErr: (*typedError)(nil),
			matched: true,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			matched := errorMatches(tc.sourceErr, tc.targetErr)

			if tc.matched != matched {
				t.Errorf("errorMatches(%v, %v) != %v", tc.sourceErr, tc.targetErr, tc.matched)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	var (
		err1  = errors.New("1")
		err2  = errors.New("2")
		err31 = errors.New("31")
		err32 = errors.New("32")
		err4  = errors.New("4")
		err5  = errors.New("4")
	)

	tcs := []struct {
		name string
		errs []error
		q    *queue
	}{
		{
			name: "ForANilError",
			errs: []error{nil},
			q:    nil,
		},
		{
			name: "ForSeveralNilErrors",
			errs: []error{nil, nil},
			q:    nil,
		},
		{
			name: "ForAnError",
			errs: []error{err1},
			q:    &queue{errs: []error{err1}},
		},
		{
			name: "ForASingleErrorQueue",
			errs: []error{&queue{errs: []error{err1}}},
			q:    &queue{errs: []error{err1}},
		},
		{
			name: "ForSeveralErrors",
			errs: []error{nil, err1, nil, err2, nil},
			q:    &queue{errs: []error{err1, err2}},
		},
		{
			name: "ForSeveralErrorsAndAQueue",
			errs: []error{nil, err1, nil, err2, &queue{errs: []error{err31, err32}}, err4, nil},
			q:    &queue{errs: []error{err1, err2, err31, err32, err4}},
		},
		{
			name: "ForAQueueThatGoesLast",
			errs: []error{nil, err1, nil, err2, &queue{errs: []error{err31, err32}}},
			q:    &queue{errs: []error{err1, err2, err31, err32}},
		},
		{
			name: "ForAQueueThatGoesFirst",
			errs: []error{&queue{errs: []error{err31, err32}}, err4},
			q:    &queue{errs: []error{err31, err32, err4}},
		},
		{
			name: "ForErrorsAndSeveralErrorQueuesInBetween",
			errs: []error{err1, err2, &queue{errs: []error{err31, err32}}, &queue{errs: []error{err4}}, err5},
			q:    &queue{errs: []error{err1, err2, err31, err32, err4, err5}},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			res := Wrap(tc.errs...)
			if res == nil && tc.q == nil {
				return
			}

			q, ok := res.(*queue)
			if !ok {
				t.Errorf("Wrap(%v) returned not an errors.queue instance", tc.errs)
			}

			if len(q.getErrors()) != len(tc.q.getErrors()) {
				t.Errorf("Wrap(%v) must contain %d errors, got %d", tc.errs, len(tc.q.getErrors()), len(q.getErrors()))
			}

			if q.Error() != tc.q.Error() {
				t.Errorf("Wrap(%v) error message mismatch, %q != %q", tc.errs, tc.q.Error(), q.Error())
			}
		})
	}
}

func TestWithMessage(t *testing.T) {
	var (
		err1 = errors.New("1")
	)

	tcs := []struct {
		name   string
		err    error
		format string
		args   []interface{}
		res    error
	}{
		{
			name:   "ForANilAndAnEmptyMessage",
			err:    nil,
			format: "",
			res:    nil,
		},
		{
			name:   "ForAnEmptyMessage",
			err:    err1,
			format: "",
			res:    err1,
		},
		{
			name:   "ForAnEmptyMessage",
			err:    err1,
			format: "",
			res:    err1,
		},
		{
			name:   "ForAnErrorAndAMessage",
			err:    typedError{"custom error"},
			format: "",
			res:    typedError{"custom error"},
		},
		{
			name:   "ForAnErrorAndAMessage",
			err:    typedError{"custom error"},
			format: "error message",
			res:    &queue{errs: []error{typedError{"custom error"}, errors.New("error message")}},
		},
		{
			name:   "ForAMessageAndAnErrorQueue",
			err:    &queue{errs: []error{err1}},
			format: "error message",
			res:    &queue{errs: []error{err1, errors.New("error message")}},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			res := WithMessage(tc.err, tc.format, tc.args...)
			if res == nil && tc.res == nil {
				return
			}

			if res == nil || tc.res == nil {
				t.Errorf("WithMessage(%v, %s) error message mismatch, %q != %q", tc.err, tc.format, tc.res.Error(), res)
			}

			if res.Error() != tc.res.Error() {
				t.Errorf("WithMessage(%v, %s) error message mismatch, %q != %q", tc.err, tc.format, tc.res.Error(), res.Error())
			}
		})
	}
}

func TestWrapWithMessage(t *testing.T) {
	var (
		err1 = errors.New("1")
		err2 = errors.New("2")
	)

	tcs := []struct {
		name       string
		err1, err2 error
		format     string
		args       []interface{}
		res        error
	}{
		{
			name: "ForANilsAndAnEmptyMessage",
			err1: nil, err2: nil,
			format: "",
			res:    nil,
		},
		{
			name: "ForAnEmptyMessage",
			err1: err1, err2: err2,
			format: "",
			res:    &queue{errs: []error{err1, err2}},
		},
		{
			name: "ForAMessage",
			err1: err1, err2: err2,
			format: "message",
			res:    &queue{errs: []error{err1, err2, errors.New("message")}},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			res := WrapWithMessage(tc.err1, tc.err2, tc.format, tc.args...)
			if res == nil && tc.res == nil {
				return
			}

			if res == nil || tc.res == nil {
				t.Errorf(
					"WrapWithMessage(%v, %v, %s) error message mismatch, %q != %q",
					tc.err1, tc.err2, tc.format, tc.res.Error(), res,
				)
			}

			if res.Error() != tc.res.Error() {
				t.Errorf(
					"WrapWithMessage(%v, %v, %s) error message mismatch, %q != %q",
					tc.err1, tc.err2, tc.format, tc.res.Error(), res.Error(),
				)
			}
		})
	}
}

func TestIsErrNil(t *testing.T) {
	var (
		zeroErrorValue typedError
		zeroErrorPtr   *typedError
	)
	tcs := []struct {
		name string
		err  error
		res  bool
	}{
		{
			name: "ForANilError",
			err:  nil,
			res:  true,
		},
		{
			name: "ForANotNilError",
			err:  errors.New("some error"),
			res:  false,
		},
		{
			name: "ForAnEmptyValueError",
			err:  zeroErrorValue,
			res:  false,
		},
		{
			name: "ForAnEmptyPointerError",
			err:  zeroErrorPtr,
			res:  true,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			res := isErrNil(tc.err)
			if res != tc.res {
				t.Errorf("isErrNil(%v) != %v", tc.err, tc.res)
			}
		})
	}
}

func TestFetchByType(t *testing.T) {
	var (
		err1 = errors.New("1")
		err2 = errors.New("2")
		err3 = errors.New("3")
	)

	tcs := []struct {
		name   string
		source error
		target interface{}
		res    error
	}{
		{
			name:   "ForNilParameters",
			source: nil,
			target: nil,
			res:    nil,
		},
		{
			name:   "ForANilErrorQueue",
			source: nil,
			target: err1,
			res:    nil,
		},
		{
			name:   "ForANilTarget",
			source: err1,
			target: nil,
			res:    nil,
		},
		{
			name:   "ForExistingErrorType",
			source: &queue{errs: []error{err1, typedError{"2"}, err3}},
			target: (*typedError)(nil),
			res:    typedError{"2"},
		},
		{
			name:   "ForAnInterface",
			source: &queue{errs: []error{err1, typedError{"2"}, err3}},
			target: (typedErrorInterface)(nil),
			res:    nil,
		},
		{
			name:   "ForAnInterfacePointer",
			source: &queue{errs: []error{err1, typedError{"2"}, err3}},
			target: (*typedErrorInterface)(nil),
			res:    typedError{"2"},
		},
		{
			name:   "ForSeveralCustomErrors",
			source: &queue{errs: []error{err1, typedError{"2"}, err3, typedError{"4"}}},
			target: (*typedError)(nil),
			res:    typedError{"4"},
		},
		{
			name:   "ForAnErrorThatIsNotInErrorQueue",
			source: &queue{errs: []error{err1, err2}},
			target: (*typedError)(nil),
			res:    nil,
		},
		{
			name:   "ForAnErrorSelectedByTypePointer",
			source: typedError{"2"},
			target: (*typedError)(nil),
			res:    typedError{"2"},
		},
		{
			name:   "ForAnErrorSelectedByInterfacePointer",
			source: typedError{"2"},
			target: (*typedErrorInterface)(nil),
			res:    typedError{"2"},
		},
		{
			name:   "ForAnErrorSelectedByNotMatchingInterfacePointer",
			source: typedError{"2"},
			target: (*fmt.Formatter)(nil),
			res:    nil,
		},
		{
			name:   "ForAnErrorSelectedByInterface",
			source: typedError{"2"},
			target: (typedErrorInterface)(nil),
			res:    nil,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			res := FetchByType(tc.source, tc.target)

			if res == nil && tc.res == nil {
				return
			}

			if res == nil || tc.res == nil {
				t.Errorf("FetchByType(%v, %v) != %v, got %v", tc.source, tc.target, tc.res, res)
			}

			if res.Error() != tc.res.Error() {
				t.Errorf("FetchByType(%v, %v) != %v, got %v", tc.source, tc.target, tc.res, res)
			}
		})
	}
}

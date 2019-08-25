package errors

import (
	"fmt"
	"reflect"
	"testing"
)

type (
	customError struct{ msg string }
	//typedErrorPtr       struct{ msg string }
	customErrorInterface interface{ Message() string }
	customInterface      interface{ Send() string }
)

func (e customError) Error() string   { return e.msg }
func (e customError) Message() string { return e.msg }

//func (e *typedErrorPtr) Error() string { return e.msg }

func TestWrapNotToCauseSideEffects(t *testing.T) {
	var (
		err1    = New("1")
		err2    = New("2")
		err3    = New("3")
		q21     = &queue{errs: []error{err1, err2}}
		q21copy = &queue{errs: []error{err1, err2}}
		q3      = &queue{errs: []error{err3}}
		q3copy  = &queue{errs: []error{err3}}
	)

	q := Wrap(q21, q3)

	if !reflect.DeepEqual(q21, q21copy) {
		t.Errorf("Wrap(q1, q2) mustn't cause side effects. q1 was changed: %v != %v", q21, q21copy)
	}
	if !reflect.DeepEqual(q3, q3copy) {
		t.Errorf("Wrap(q1, q2) mustn't cause side effects. q1 was changed: %v != %v", q3, q3copy)
	}
	if q == q21 {
		t.Errorf("Wrap(q1, q2) must return new queue instance. q1 returned")
	}
	if q == q3 {
		t.Errorf("Wrap(q1, q2) must return new queue instance. q2 returned")
	}
}

func TestFetch(t *testing.T) {
	var (
		err1 = New("1")
		err2 = New("2")
		err3 = New("3")
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
			name: "ForAnErrorQueueWithMatchingError",
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

func TestWrap(t *testing.T) {
	var (
		err1  = New("1")
		err2  = New("2")
		err31 = New("31")
		err32 = New("32")
		err4  = New("4")
		err5  = New("4")
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
			errs: []error{newQueue(err1)},
			q:    newQueue(err1),
		},
		{
			name: "ForSeveralErrors",
			errs: []error{nil, err1, nil, err2, nil},
			q:    newQueue(err1, err2),
		},
		{
			name: "ForSeveralErrorsAndAQueue",
			errs: []error{nil, err1, nil, err2, newQueue(err31, err32), err4, nil},
			q:    newQueue(err1, err2, err31, err32, err4),
		},
		{
			name: "ForAQueueThatGoesLast",
			errs: []error{nil, err1, nil, err2, newQueue(err31, err32)},
			q:    newQueue(err1, err2, err31, err32),
		},
		{
			name: "ForAQueueThatGoesFirst",
			errs: []error{newQueue(err31, err32), err4},
			q:    newQueue(err31, err32, err4),
		},
		{
			name: "ForErrorsAndSeveralErrorQueuesInBetween",
			errs: []error{err1, err2, newQueue(err31, err32), newQueue(err4), err5},
			q:    newQueue(err1, err2, err31, err32, err4, err5),
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
		err1 = New("1")
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
			err:    customError{"custom error"},
			format: "",
			res:    customError{"custom error"},
		},
		{
			name:   "ForAnErrorAndAMessage",
			err:    customError{"custom error"},
			format: "error message",
			res:    &queue{errs: []error{customError{"custom error"}, New("error message")}},
		},
		{
			name:   "ForAMessageAndAnErrorQueue",
			err:    &queue{errs: []error{err1}},
			format: "error message",
			res:    &queue{errs: []error{err1, New("error message")}},
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
		err1 = New("1")
		err2 = New("2")
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
			res:    &queue{errs: []error{err1, err2, New("message")}},
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
		zeroErrorValue customError
		zeroErrorPtr   *customError
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
			err:  New("some error"),
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
		err1 = New("1")
		err2 = New("2")
		err3 = New("3")
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
			source: &queue{errs: []error{err1, customError{"2"}, err3}},
			target: (*customError)(nil),
			res:    customError{"2"},
		},
		{
			name:   "ForAnInterface",
			source: &queue{errs: []error{err1, customError{"2"}, err3}},
			target: (customErrorInterface)(nil),
			res:    nil,
		},
		{
			name:   "ForAnInterfacePointer",
			source: &queue{errs: []error{err1, customError{"2"}, err3}},
			target: (*customErrorInterface)(nil),
			res:    customError{"2"},
		},
		{
			name:   "ForSeveralCustomErrors",
			source: &queue{errs: []error{err1, customError{"2"}, err3, customError{"4"}}},
			target: (*customError)(nil),
			res:    customError{"4"},
		},
		{
			name:   "ForAnErrorThatIsNotInErrorQueue",
			source: &queue{errs: []error{err1, err2}},
			target: (*customError)(nil),
			res:    nil,
		},
		{
			name:   "ForAnErrorSelectedByTypePointer",
			source: customError{"2"},
			target: (*customError)(nil),
			res:    customError{"2"},
		},
		{
			name:   "ForAnErrorSelectedByInterfacePointer",
			source: customError{"2"},
			target: (*customErrorInterface)(nil),
			res:    customError{"2"},
		},
		{
			name:   "ForAnErrorSelectedByNotMatchingInterfacePointer",
			source: customError{"2"},
			target: (*fmt.Formatter)(nil),
			res:    nil,
		},
		{
			name:   "ForAnErrorSelectedByInterface",
			source: customError{"2"},
			target: (customErrorInterface)(nil),
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

func TestCompareErrors(t *testing.T) {
	var (
		err1 = New("1")
		err2 = New("2")
	)

	tcs := []struct {
		name                 string
		sourceErr, targetErr error
		matched              bool
	}{
		{
			name:      "ForMatchingErrors",
			sourceErr: err1, targetErr: err1,
			matched: true,
		},
		{
			name:      "ForNotMatchingErrors",
			sourceErr: err2, targetErr: err1,
			matched: false,
		},
		{
			name:      "ForAnErrorAndAMatchingErrorQueueAsSource",
			sourceErr: &queue{errs: []error{err1}}, targetErr: err1,
			matched: true,
		},
		{
			name:      "ForAnErrorAndAMatchingErrorQueueAsTarget",
			sourceErr: err1, targetErr: &queue{errs: []error{err1}},
			matched: true,
		},
		{
			name:      "ForAnErrorAndANotMatchingErrorQueueAsSource",
			sourceErr: &queue{errs: []error{err1}}, targetErr: err2,
			matched: false,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			match := compareErrs(tc.sourceErr, tc.targetErr)

			if tc.matched != match {
				t.Errorf("compareErrs(%v, %v) != %v", tc.sourceErr, tc.targetErr, tc.matched)
			}
		})
	}
}

func TestErrorMatches(t *testing.T) {
	var (
		err1 = New("1")
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
			sourceErr: New("error"), targetErr: New("error"),
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
			sourceErr: customError{msg: "error"}, targetErr: (*error)(nil),
			matched: true,
		},
		{
			name:      "ForACustomErrorAndCustomInterfacePointer",
			sourceErr: customError{msg: "error"}, targetErr: (*customErrorInterface)(nil),
			matched: true,
		},
		{
			name:      "ForACustomErrorAndAPointerToErrorStruct",
			sourceErr: customError{msg: "error"}, targetErr: (*customError)(nil),
			matched: true,
		},
		{
			name:      "ForACustomErrorPointerAndAPointerToErrorStruct",
			sourceErr: &customError{msg: "error"}, targetErr: (*customError)(nil),
			matched: true,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			targetType, targetElem, err := getTypeElem(tc.targetErr)
			if err != nil {
				t.Errorf("target parameter cannot be nil")
			}
			matched := errorMatches(tc.sourceErr, targetType, targetElem)

			if tc.matched != matched {
				t.Errorf("errorMatches(%v, %v) != %v", tc.sourceErr, tc.targetErr, tc.matched)
			}
		})
	}
}

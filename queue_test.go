package errors

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	tcs := []struct {
		name string
		errs []error
		msg  string
	}{
		{
			name: "ForAnEmptyList",
			errs: []error{},
			msg:  "",
		},
		{
			name: "ForASingleError",
			errs: []error{New("err1")},
			msg:  "err1",
		},
		{
			name: "ForSeveralErrors",
			errs: []error{New("err1"), New("err2"), New("err3")},
			msg:  "err3 : err2 : err1",
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			msg := newQueue(tc.errs...).Error()

			if msg != tc.msg {
				t.Errorf("newQueue(%v).Error(), %q != %q", tc.errs, tc.msg, msg)
			}
		})
	}
}

// Format() :: for a %+v flag :: returns an error with stacktrace.
func TestFormatForVFlag(t *testing.T) {
	var (
		stacktraceMsg  = "error stacktrace trace"
		err1           = New("1")
		err2           = New("2")
		expectedOutput = fmt.Sprintf("%s : %s\n%s", err2, err1, stacktraceMsg)
		q              = newQueue(err1, err2)
	)
	q.stacktrace = &formatterStub{stacktraceMsg}

	output := fmt.Sprintf("%+v", q)

	if expectedOutput != output {
		t.Errorf(`error output (%%+v format) %q != %q`, expectedOutput, output)
	}
}

// formatterStub is a stub for the stacktrace instance to test error formatting.
type formatterStub struct{ msg string }

func (fs *formatterStub) Format(s fmt.State, _ rune) { _, _ = s.Write([]byte(fs.msg)) }

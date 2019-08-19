package errors

import (
	"fmt"
	"testing"
)

// Format() :: for a %+v flag :: returns an error with stacktrace.
func TestFormatForVFlag(t *testing.T) {
	stackMessage := "error stacktrace trace"
	q := queue{errs: []error{fmt.Errorf("1"), fmt.Errorf("2")}}
	q.stacktrace = &formatterStub{stackMessage}
	expectedO := "2 : 1\n" + stackMessage

	o := fmt.Sprintf("%+v", &q)

	if expectedO != o {
		t.Errorf(`error output (%%+v format) %q != %q`, expectedO, o)
	}
}

// formatterStub is a stub for the stacktrace instance to test error formatting.
type formatterStub struct{ msg string }

func (fs *formatterStub) Format(state fmt.State, _ rune) { _, _ = state.Write([]byte(fs.msg)) }

package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestQueueAppendPrepend(t *testing.T) {
	err1 := fmt.Errorf("1")
	err2 := fmt.Errorf("2")
	err4 := fmt.Errorf("4")
	err5 := fmt.Errorf("5")
	q3231 := newQueue()
	q3231.errs = []error{fmt.Errorf("31"), fmt.Errorf("32")}

	tcs := []struct {
		name   string
		Q      func() *queue
		len    int
		output string
	}{
		{
			name: "AppendForQueueWithTwoErrors",
			Q: func() *queue {
				q := newQueue()
				q.append(err1)
				q.append(err2)
				q.append(q3231)
				q.append(err4)
				return q
			},
			len: 5, output: "4 : 32 : 31 : 2 : 1",
		},
		{
			name: "PrependForAQueueWithTwoErrors",
			Q: func() *queue {
				q := newQueue()
				q.prepend(err5)
				q.prepend(err4)
				q.prepend(q3231)
				q.prepend(err2)
				return q
			},
			len: 5, output: "5 : 4 : 32 : 31 : 2",
		},
		{
			name: "PrependForASingleError",
			Q: func() *queue {
				q := newQueue()
				q.prepend(err1)
				return q
			},
			len: 1, output: "1",
		},
		{
			name: "ErrorsWithEmptyMessagesReturnEmptyPlaceholders",
			Q: func() *queue {
				q := newQueue()
				q.append(errors.New(""))
				q.append(err1)
				q.append(errors.New(""))
				q.append(err2)
				return q
			},
			len: 4, output: "2 :  : 1 : ",
		},
		{
			name: "ErrorsWithEmptyMessagesReturnEmptyPlaceholders",
			Q: func() *queue {
				q := newQueue()
				q.append(errors.New(""))
				q.append(errors.New(""))
				return q
			},
			len: 2, output: " : ",
		},
		{
			name: "EmptyErrorQueue",
			Q: func() *queue {
				q := newQueue()
				return q
			},
			len: 0, output: "",
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			q := tc.Q()

			output := q.Error()
			outputS := fmt.Sprintf("%s", q)
			outputV := fmt.Sprintf("%v", q)

			if tc.len != len(q.getErrors()) {
				t.Errorf(`error queue must contain %d errors, got %d`, tc.len, len(q.getErrors()))
			}
			if tc.output != output {
				t.Errorf(`error output %q != %q`, tc.output, q.Error())
			}
			if tc.output != outputS {
				t.Errorf(`error output (%%s format) %q != %q`, tc.output, outputS)
			}
			if tc.output != outputV {
				t.Errorf(`error output (%%v format) %q != %q`, tc.output, outputV)
			}
		})
	}
}

// Format() :: for a %+v flag :: returns an error with stacktrace.
func TestFormatForVFlag(t *testing.T) {
	stackMessage := "error stacktrace trace"
	q := newQueue()
	q.stacktrace = &formatterStub{stackMessage}
	q.append(fmt.Errorf("1"))
	q.append(fmt.Errorf("2"))
	expectedO := "2 : 1\n" + stackMessage

	o := fmt.Sprintf("%+v", q)

	if expectedO != o {
		t.Errorf(`error output (%%+v format) %q != %q`, expectedO, o)
	}
}

// formatterStub is a stub for the stacktrace instance to test error formatting.
type formatterStub struct{ msg string }

func (fs *formatterStub) Format(state fmt.State, _ rune) { _, _ = state.Write([]byte(fs.msg)) }

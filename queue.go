package errors

import (
	"bytes"
	"fmt"
)

// errMsgSeparator joins error messages in form of "outer error : inner error".
const errMsgSeparator = " : "

// queue object queues application errors into an ordered collection.
// All errors are stored in LIFO order, that's why getErrors() reverses the list.
type queue struct {
	errs       []error       // The Double-Ended Queue with errors.
	stacktrace fmt.Formatter // Stacktrace at the moment of creation.
}

// newQueue returns a new queue instance with a stacktrace data at the moment of invocation.
// Contract: all errors from the errs list are not nil.
func newQueue(errs ...error) *queue {
	return &queue{errs: errs, stacktrace: newStacktrace()}
}

// Error returns an error message.
func (q *queue) Error() (errMsg string) {
	buf := new(bytes.Buffer)
	for _, err := range q.getErrors() {
		if buf.Len() != 0 {
			buf.WriteString(errMsgSeparator)
		}
		buf.WriteString(err.Error())
	}

	return buf.String()
}

// Format formats an error message for the queue object.
// %+v additionally prints out an error stacktrace.
func (q *queue) Format(st fmt.State, verb rune) {
	_, _ = st.Write([]byte(q.Error()))
	if verb == 'v' && st.Flag('+') {
		_, _ = st.Write([]byte("\n"))
		q.stacktrace.Format(st, verb)
	}
}

// getErrors returns the errors in reverse order.
func (q *queue) getErrors() []error {
	errsLen := len(q.errs)
	errs := make([]error, errsLen)
	for i := range q.errs {
		errs[i] = q.errs[errsLen-i-1]
	}

	return errs
}

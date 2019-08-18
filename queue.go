package errors

import (
	"bytes"
	"fmt"

	"github.com/ameteiko/errors/stacktrace"
)

// errMsgSeparator joins error messages in form of "outer error" : "inner error".
const errMsgSeparator = " : "

// queue queues application errors into an ordered collection.
// All errors are stored LIFO order, that's why getErrors() reverses the list.
type queue struct {
	errs       []error       // The Double-Ended Queue with errors.
	stacktrace fmt.Formatter // Stacktrace at the moment of creation.
}

// newQueue returns a new queue instance with invocation stacktrace data.
func newQueue() *queue {
	return &queue{stacktrace: stacktrace.New()}
}

// Error returns a composite error message.
func (q *queue) Error() (errMsg string) {
	if len(q.errs) == 0 {
		return ""
	}

	buf := new(bytes.Buffer)
	for i, err := range q.getErrors() {
		if i != 0 {
			buf.WriteString(errMsgSeparator)
		}
		buf.WriteString(err.Error())
	}

	return buf.String()
}

// Format formats a queue error message.
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

// prepend adds an error to the front of the queue.
// If prepended error is yet another queue object, then extract all its errors and prepend them sequentially.
func (q *queue) prepend(err error) {
	var errsToPrepend []error
	if errQ, ok := err.(*queue); ok {
		// Use errQ.errs, not errQ.getErrors() to retrieve the errors in reversed order.
		// It's crucial to keep reversed order, because errors are prepended at once with copy().
		errsToPrepend = errQ.errs
	} else {
		errsToPrepend = []error{err}
	}

	errsNum := len(errsToPrepend)
	errs := make([]error, len(q.errs)+errsNum)
	copy(errs[:errsNum], errsToPrepend)
	copy(errs[errsNum:], q.errs)
	q.errs = errs
}

// append adds an error to the back of the queue.
// If appended error is yet another queue object, then extract all its errors and append them sequentially.
func (q *queue) append(err error) {
	if errQ, ok := err.(*queue); ok {
		q.errs = append(q.errs, errQ.errs...)
		return
	}

	q.errs = append(q.errs, err)
}

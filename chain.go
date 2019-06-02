package errors

import (
	"fmt"

	"github.com/ameteiko/errors/stack"
)

const (
	errorMessageSeparator = " : " // a separator value used to join error messages. (outer error : inner error)
)

// chain is the aggregate for the application errors chained in the order they arose.
type chain struct {
	errs  []error      // The collection of the errors.
	stack *stack.Stack // The stacktrace at the chain initialization moment.
}

// newChain returns a new chain instance.
func newChain() *chain {
	return &chain{stack: stack.Get()}
}

// Format formats the error message for the error chain.
// On %+v it additionally prints the error stacktrace.
func (c *chain) Format(state fmt.State, verb rune) {
	fmt.Fprint(state, c.Error())
	if verb == 'v' && state.Flag('+') {
		fmt.Fprintln(state)
		c.stack.Format(state, verb)
	}
}

// Error returns an error message for the error chain.
func (c *chain) Error() (errMsg string) {
	if len(c.errs) == 0 {
		return ""
	}

	for i, err := range c.getErrors() {
		if i == 0 {
			errMsg = err.Error()
			continue
		}
		errMsg += errorMessageSeparator + err.Error()
	}

	return errMsg
}

// getErrors returns the error chain in reverse order.
func (c *chain) getErrors() []error {
	errsLen := len(c.errs)
	errs := make([]error, errsLen)
	for i := range c.errs {
		errs[i] = c.errs[errsLen-i-1]
	}

	return errs
}

// prepend adds an error to the bottom of the chain.
func (c *chain) prepend(err error) {
	errs := make([]error, len(c.errs)+1)
	errs[0] = err
	copy(errs[1:], c.errs)
	c.errs = errs
}

// append adds an error to the top of the chain.
func (c *chain) append(err error) {
	c.errs = append(c.errs, err)
}

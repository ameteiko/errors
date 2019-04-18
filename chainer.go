package errors

type Enchainer interface {
}

//
// Chainer interface wraps error objects into a manageable chain.
//
type Chainer interface {
	error

	//
	// Append appends an error to the head of the chain.
	//
	Append(error) *Chain

	//
	// Prepend prepends an error to the bottom of the chain.
	//
	Prepend(error) *Chain

	//
	// GetErrors returns the error chain.
	//
	GetErrors() []error

	//
	// WithMessage returns a chain with attached formatted message.
	//
	WithMessage(format string, args ...interface{}) error
}

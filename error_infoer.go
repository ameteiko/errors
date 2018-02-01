package errors


//
// ErrorInfoer is an Error interface.
//
type ErrorInfoer interface {
	error

	//
	// WithMessage returns a new error instance with extra info.
	//
	WithMessage(format string, args ...interface{}) error
}
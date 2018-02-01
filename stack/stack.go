package stack

import (
	"fmt"
	"runtime"
)

//
// General constants.
//
const (
	StacktraceDepth = 16
)

//
// Stack wraps a stack of program counters.
//
type Stack []uintptr

//
// Format prints the stack trace.
//
func (s *Stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		fallthrough
	default:
		for _, pc := range *s {
			var (
				function = "unknown"
				file     string
				line     int
			)
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				function = fn.Name()
				file, line = fn.FileLine(pc)
			}
			fmt.Fprintf(st, "\t%s\t↪\t%s:%d\n", function, file, line)
		}
	}
}

//
// Get returns current invoker stack.
//
func Get() *Stack {
	pcs := make([]uintptr, StacktraceDepth)
	n := runtime.Callers(3, pcs[:])
	var st Stack = pcs[0:n]

	return &st
}

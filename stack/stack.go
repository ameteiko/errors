package stack

import (
	"fmt"
	"go/build"
	"runtime"
	"strings"
)

// General constants.
const (
	StacktraceDepth = 16
)

// Stack wraps a stack of program counters.
type Stack []uintptr

// Format prints the stack trace.
func (s *Stack) Format(st fmt.State, _ rune) {
	for _, pc := range *s {
		var (
			function = "unknown"
			file     string
			line     int
		)
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			function = sanitizeFunctionName(fn.Name())
			file, line = fn.FileLine(pc)
			file = sanitizeFilename(file)
		}

		fmt.Fprintf(st, " @ %s:%s:%d\n", file, function, line)
	}
}

// Get returns current invoker stack.
func Get() *Stack {
	pcs := make([]uintptr, StacktraceDepth)
	n := runtime.Callers(3, pcs)
	var st Stack = pcs[0:n]

	return &st
}

// sanitizeFunctionName cleans up a function name from the module full name info.
//
// The function name comes with full module path like github.com/ameteiko/errors/errors.Wrap, so it's
// necessary to strip all the information that comes before the last dot.
func sanitizeFunctionName(n string) string {
	lastDotIndex := strings.LastIndex(n, ".")

	return n[lastDotIndex+1:]
}

// sanitizeFilename cleans up a file name from the GOROOT and GOPATH prefixes.
//
// The filename comes with full module path like /Users/ameteiko/Projects/go/src/github.com/ameteiko/errors/errors.go:,
// so it's necessary to strip all the information that is related to GOPATH or GOROOT.
func sanitizeFilename(filename string) string {
	pathPrefixes := []string{
		build.Default.GOPATH + "/src",
		runtime.GOROOT() + "/src",
	}

	for _, pathPrefix := range pathPrefixes {
		if strings.HasPrefix(filename, pathPrefix) {
			return filename[len(pathPrefix)+1:]
		}
	}

	return filename
}

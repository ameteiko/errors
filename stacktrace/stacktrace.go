package stacktrace

import (
	"bytes"
	"fmt"
	"go/build"
	"runtime"
	"strconv"
	"strings"
)

const (
	stacktraceDepth = 16
)

var (
	// nolint:gochecknoglobals
	gopath = build.Default.GOPATH + "/src/"
	// nolint:gochecknoglobals
	goroot = runtime.GOROOT() + "/src/"
)

// stacktrace wraps a stacktrace of program counters.
type Stacktrace []uintptr

// New returns a stacktrace.
func New() Stacktrace {
	s := make(Stacktrace, stacktraceDepth)
	// Skipping 3 runtime callers:
	//   0 - runtime.Callers()
	//   1 - stacktrace.New()
	//   2 - errors.newQueue()
	//   3 - errors.Wrap()
	n := runtime.Callers(3, s)

	return s[0:n]
}

// Format prints the stacktrace.
func (s Stacktrace) Format(st fmt.State, _ rune) {
	b := new(bytes.Buffer)
	ff := runtime.CallersFrames([]uintptr(s))
	var mainProcessed bool
	for {
		f, more := ff.Next()
		file := sanitizeFilename(f.File)
		line := strconv.Itoa(f.Line)
		fun := sanitizeFuncName(f.Function)

		// Don't print anything beyond main.main.
		if mainProcessed && !strings.HasPrefix(fun, "main.") {
			break
		}
		if !mainProcessed && fun == "main.main()" {
			mainProcessed = true
		}

		b.WriteString("\t")
		b.WriteString(file)
		b.WriteString(":")
		b.WriteString(line)
		b.WriteString(" ")
		b.WriteString(fun)
		b.WriteString("\n")

		if !more {
			break
		}
	}

	// nolint
	_, _ = b.WriteTo(st)
}

// sanitizeFuncName trims fully qualified module path from the function name.
// Transforms:
//     github.com/ameteiko/errors/stacktrace.New -> stacktrace.New()
func sanitizeFuncName(n string) string {
	if n == "" {
		return "unknown"
	}
	lastSlashIndex := strings.LastIndex(n, "/")

	return n[lastSlashIndex+1:] + "()"
}

// sanitizeFilename trims GOROOT and GOPATH prefixes from the fully qualified file name.
// Transforms:
//     /Users/ameteiko/Projects/go/src/github.com/ameteiko/errors/errors.go -> github.com/ameteiko/errors/errors.go
func sanitizeFilename(n string) string {
	n = strings.TrimPrefix(n, gopath)
	n = strings.TrimPrefix(n, goroot)

	return n
}

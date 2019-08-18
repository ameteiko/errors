package stacktrace

import (
	"testing"
)

func TestSanitizeFuncName(t *testing.T) {
	tcs := []struct {
		desc string
		name string
		want string
	}{
		{`for an empty name`, "", "unknown"},
		{`for a name`, "main", "main()"},
		{`for a name with package`, "stacktrace.New", "stacktrace.New()"},
		{`for a fully qualified name`, "github.com/ameteiko/errors/stacktrace.New", "stacktrace.New()"},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(tc.desc, func(t *testing.T) {
			n := sanitizeFuncName(tc.name)
			if n != tc.want {
				t.Errorf("sanitizeFuncName(%q)=%q, got %q", tc.name, tc.want, n)
			}
		})
	}
}

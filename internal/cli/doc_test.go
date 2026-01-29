package cli

import "testing"

type stubArgs []string

func (a stubArgs) Get(n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}

func (a stubArgs) First() string {
	return a.Get(0)
}

func (a stubArgs) Tail() []string {
	if len(a) >= 2 {
		tail := make([]string, len(a)-1)
		copy(tail, a[1:])
		return tail
	}
	return []string{}
}

func (a stubArgs) Len() int {
	return len(a)
}

func (a stubArgs) Present() bool {
	return len(a) != 0
}

func (a stubArgs) Slice() []string {
	out := make([]string, len(a))
	copy(out, a)
	return out
}

func TestParseDocArgs(t *testing.T) {
	tests := []struct { //nolint:govet // fieldalignment
		name      string
		args      []string
		wantPkg   string
		wantPat   string
		wantError bool
	}{
		{
			name:    "pattern only",
			args:    []string{"Range"},
			wantPkg: "",
			wantPat: "Range",
		},
		{
			name:    "package and pattern",
			args:    []string{"builtin/core", "Range"},
			wantPkg: "builtin/core",
			wantPat: "Range",
		},
		{
			name:    "package dot pattern",
			args:    []string{"builtin/core", ".", "Range"},
			wantPkg: "builtin/core",
			wantPat: "Range",
		},
		{
			name:      "too few arguments",
			args:      []string{},
			wantError: true,
		},
		{
			name:      "too many arguments",
			args:      []string{"a", "b", "c", "d"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, pat, err := parseDocArgs(stubArgs(tt.args))
			if (err != nil) != tt.wantError {
				t.Fatalf("expected error=%v got %v", tt.wantError, err)
			}
			if err != nil {
				return
			}
			if pkg != tt.wantPkg {
				t.Fatalf("expected pkg %q got %q", tt.wantPkg, pkg)
			}
			if pat != tt.wantPat {
				t.Fatalf("expected pattern %q got %q", tt.wantPat, pat)
			}
		})
	}
}

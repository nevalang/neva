package view

import (
	"testing"

	"github.com/nevalang/neva/pkg/core"
)

func TestModuleID(t *testing.T) {
	t.Parallel()

	if got := moduleID(core.ModuleRef{Path: "@"}); got != "module/@" {
		t.Fatalf("moduleID() = %q, want %q", got, "module/@")
	}

	if got := moduleID(core.ModuleRef{Path: "std", Version: "1.2.3"}); got != "module/std@1.2.3" {
		t.Fatalf("moduleID() = %q, want %q", got, "module/std@1.2.3")
	}
}

func TestPathIDs(t *testing.T) {
	t.Parallel()

	loc := core.Location{
		ModRef:   core.ModuleRef{Path: "github.com/acme/mod", Version: "v1.0.0"},
		Package:  "pkg",
		Filename: "main",
	}

	if got := packageID(loc.ModRef, loc.Package); got != "module/github.com_acme_mod@v1.0.0/package/pkg" {
		t.Fatalf("packageID() = %q", got)
	}

	if got := fileID(loc); got != "module/github.com_acme_mod@v1.0.0/package/pkg/file/main" {
		t.Fatalf("fileID() = %q", got)
	}

	if got := componentID(loc, "Main", 0); got != "module/github.com_acme_mod@v1.0.0/package/pkg/file/main/component/Main@0" {
		t.Fatalf("componentID() = %q", got)
	}

	if got := nodeID(componentID(loc, "Main", 0), "echo"); got != "module/github.com_acme_mod@v1.0.0/package/pkg/file/main/component/Main@0/node/echo" {
		t.Fatalf("nodeID() = %q", got)
	}
}

func TestSanitizeSegment(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"":        "_",
		"  ":      "__",
		"a/b":     "a_b",
		"a:b":     "a_b",
		"a b":     "a_b",
		"a\tb\n#": "a_b__",
	}

	for in, want := range cases {
		if got := sanitizeSegment(in); got != want {
			t.Fatalf("sanitizeSegment(%q) = %q, want %q", in, got, want)
		}
	}
}

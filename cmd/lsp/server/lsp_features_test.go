package server

import (
	"testing"

	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
)

func TestGeneralCompletionsIncludesEntitiesAndNodeNames(t *testing.T) {
	t.Parallel()

	moduleRef := core.ModuleRef{Path: "@"}
	build := &src.Build{
		Modules: map[core.ModuleRef]src.Module{
			moduleRef: {
				Packages: map[string]src.Package{
					"main": {
						"file1": {
							Entities: map[string]src.Entity{
								"FooComponent": {Kind: src.ComponentEntity},
								"FooType":      {Kind: src.TypeEntity},
								"FooConst":     {Kind: src.ConstEntity},
								"FooInterface": {Kind: src.InterfaceEntity},
							},
						},
					},
				},
			},
		},
	}

	fileCtx := &fileContext{
		moduleRef:   moduleRef,
		packageName: "main",
		file: src.File{
			Imports: map[string]src.Import{
				"fmt": {Package: "fmt"},
			},
		},
	}
	compCtx := &componentContext{
		component: src.Component{
			Nodes: map[string]src.Node{
				"parser": {},
				"writer": {},
			},
		},
	}

	items := (&Server{}).generalCompletions(build, fileCtx, compCtx)
	labels := map[string]struct{}{}
	for _, item := range items {
		labels[item.Label] = struct{}{}
	}

	for _, expectedLabel := range []string{
		"FooComponent",
		"FooType",
		"FooConst",
		"FooInterface",
		"parser",
		"writer",
	} {
		if _, ok := labels[expectedLabel]; !ok {
			t.Fatalf("missing completion label %q", expectedLabel)
		}
	}
}

func TestParseCodeLensDataValidation(t *testing.T) {
	t.Parallel()

	assertParsed := func(testName string, raw any, wantValid bool) {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			_, ok := parseCodeLensData(raw)
			if ok != wantValid {
				t.Fatalf("parseCodeLensData() valid=%v, want %v", ok, wantValid)
			}
		})
	}

	assertParsed("references_kind_is_valid", map[string]any{
		"uri":  "file:///tmp/main.neva",
		"name": "Main",
		"kind": "references",
	}, true)
	assertParsed("implementations_kind_is_valid", map[string]any{
		"uri":  "file:///tmp/main.neva",
		"name": "Main",
		"kind": "implementations",
	}, true)
	assertParsed("unknown_kind_is_rejected", map[string]any{
		"uri":  "file:///tmp/main.neva",
		"name": "Main",
		"kind": "unknown",
	}, false)
	assertParsed("missing_name_is_rejected", map[string]any{
		"uri":  "file:///tmp/main.neva",
		"kind": "references",
	}, false)
}

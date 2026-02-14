package server

import (
	"testing"

	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TestGeneralCompletionsIncludesEntitiesAndNodeNames verifies MVP completion coverage:
// package entities and current-component node names must both be included.
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
	itemsByLabel := map[string]protocol.CompletionItem{}
	for _, item := range items {
		itemsByLabel[item.Label] = item
	}

	for _, expectedLabel := range []string{
		"FooComponent",
		"FooType",
		"FooConst",
		"FooInterface",
		"parser",
		"writer",
	} {
		if _, ok := itemsByLabel[expectedLabel]; !ok {
			t.Fatalf("missing completion label %q", expectedLabel)
		}
	}

	assertCompletionKind(t, itemsByLabel, "FooComponent", protocol.CompletionItemKindFunction)
	assertCompletionKind(t, itemsByLabel, "FooType", protocol.CompletionItemKindClass)
	assertCompletionKind(t, itemsByLabel, "FooConst", protocol.CompletionItemKindConstant)
	assertCompletionKind(t, itemsByLabel, "FooInterface", protocol.CompletionItemKindInterface)
	assertCompletionKind(t, itemsByLabel, "parser", protocol.CompletionItemKindVariable)
}

// TestParseCodeLensDataValidation checks strict validation of the generic LSP CodeLens.Data payload.
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

// TestImplementationLocationsForInterface verifies MVP interface-implementation discovery for CodeLens.
func TestImplementationLocationsForInterface(t *testing.T) {
	t.Parallel()

	ifaceMeta := core.Meta{
		Start: core.Position{Line: 1, Column: 0},
		Stop:  core.Position{Line: 1, Column: 8},
	}
	componentMeta := core.Meta{
		Start: core.Position{Line: 10, Column: 0},
		Stop:  core.Position{Line: 10, Column: 8},
	}

	interfaceEntity := src.Entity{
		Kind:      src.InterfaceEntity,
		Interface: src.Interface{IO: src.IO{In: map[string]src.Port{}, Out: map[string]src.Port{}}, Meta: ifaceMeta},
	}
	componentEntity := src.Entity{
		Kind: src.ComponentEntity,
		Component: []src.Component{
			{
				Interface: src.Interface{IO: src.IO{In: map[string]src.Port{}, Out: map[string]src.Port{}}},
				Meta:      componentMeta,
			},
		},
	}

	moduleRef := core.ModuleRef{Path: "@"}
	build := &src.Build{
		Modules: map[core.ModuleRef]src.Module{
			moduleRef: {
				Packages: map[string]src.Package{
					"main": {
						"iface_file": {
							Entities: map[string]src.Entity{
								"Greeter": interfaceEntity,
							},
						},
						"component_file": {
							Entities: map[string]src.Entity{
								"HelloGreeter": componentEntity,
							},
						},
					},
				},
			},
		},
	}

	server := &Server{workspacePath: "/tmp/workspace"}
	target := &resolvedEntity{
		moduleRef:   moduleRef,
		packageName: "main",
		name:        "Greeter",
		filePath:    "/tmp/workspace/main/iface_file.neva",
		entity:      interfaceEntity,
	}

	locations := server.implementationLocationsForEntity(build, target)
	if len(locations) != 1 {
		t.Fatalf("implementationLocationsForEntity() count=%d, want 1", len(locations))
	}
}

func assertCompletionKind(
	t *testing.T,
	itemsByLabel map[string]protocol.CompletionItem,
	label string,
	wantKind protocol.CompletionItemKind,
) {
	t.Helper()
	item, ok := itemsByLabel[label]
	if !ok {
		t.Fatalf("missing completion label %q", label)
	}
	if item.Kind == nil {
		t.Fatalf("completion %q kind is nil", label)
	}
	if *item.Kind != wantKind {
		t.Fatalf("completion %q kind=%v, want %v", label, *item.Kind, wantKind)
	}
}

package server

import (
	"sync"
	"testing"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	ts "github.com/nevalang/neva/pkg/typesystem"
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

// TestTextDocumentCodeLensEmitsExpectedKinds ensures interface declarations get two lenses
// (references + implementations), while other entities keep references only.
func TestTextDocumentCodeLensEmitsExpectedKinds(t *testing.T) {
	t.Parallel()

	server, docURI := buildTestLSPServerWithSingleFile()

	lenses, err := server.TextDocumentCodeLens(nil, &protocol.CodeLensParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("TextDocumentCodeLens() error = %v", err)
	}
	if len(lenses) != 4 {
		t.Fatalf("TextDocumentCodeLens() count=%d, want 4", len(lenses))
	}

	// We expect: Greeter references + implementations, HelloGreeter references, Answer references.
	got := map[string]map[codeLensKind]struct{}{}
	for _, lens := range lenses {
		parsedCodeLensData, ok := parseCodeLensData(lens.Data)
		if !ok {
			t.Fatalf("invalid code lens data payload: %#v", lens.Data)
		}
		if _, exists := got[parsedCodeLensData.Name]; !exists {
			got[parsedCodeLensData.Name] = map[codeLensKind]struct{}{}
		}
		got[parsedCodeLensData.Name][parsedCodeLensData.Kind] = struct{}{}
	}

	assertLensKinds(t, got, "Greeter", []codeLensKind{codeLensKindReferences, codeLensKindImplementations})
	assertLensKinds(t, got, "HelloGreeter", []codeLensKind{codeLensKindReferences})
	assertLensKinds(t, got, "Answer", []codeLensKind{codeLensKindReferences})
}

// TestCodeLensResolveForInterfaceImplementation verifies resolved implementation lenses
// point to implementing components via show-references command payload.
func TestCodeLensResolveForInterfaceImplementation(t *testing.T) {
	t.Parallel()

	server, docURI := buildTestLSPServerWithSingleFile()
	lens := &protocol.CodeLens{
		Range: protocol.Range{Start: protocol.Position{Line: 0, Character: 0}},
		Data: codeLensData{
			URI:  docURI,
			Name: "Greeter",
			Kind: codeLensKindImplementations,
		},
	}

	resolvedLens, err := server.CodeLensResolve(nil, lens)
	if err != nil {
		t.Fatalf("CodeLensResolve() error = %v", err)
	}
	if resolvedLens.Command == nil {
		t.Fatalf("CodeLensResolve() command is nil")
	}
	if resolvedLens.Command.Title != "1 implementations" {
		t.Fatalf("CodeLensResolve() title=%q, want %q", resolvedLens.Command.Title, "1 implementations")
	}
}

// TestCodeLensResolveForInterfaceReferences verifies interface references include explicit refs
// plus component implementations in the MVP relationship model.
func TestCodeLensResolveForInterfaceReferences(t *testing.T) {
	t.Parallel()

	server, docURI := buildTestLSPServerWithSingleFile()
	lens := &protocol.CodeLens{
		Range: protocol.Range{Start: protocol.Position{Line: 0, Character: 0}},
		Data: codeLensData{
			URI:  docURI,
			Name: "Greeter",
			Kind: codeLensKindReferences,
		},
	}

	resolvedLens, err := server.CodeLensResolve(nil, lens)
	if err != nil {
		t.Fatalf("CodeLensResolve() error = %v", err)
	}
	if resolvedLens.Command == nil {
		t.Fatalf("CodeLensResolve() command is nil")
	}
	// One explicit const reference + one implementing component.
	if resolvedLens.Command.Title != "2 references" {
		t.Fatalf("CodeLensResolve() title=%q, want %q", resolvedLens.Command.Title, "2 references")
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

func assertLensKinds(
	t *testing.T,
	got map[string]map[codeLensKind]struct{},
	entityName string,
	want []codeLensKind,
) {
	t.Helper()
	kinds, ok := got[entityName]
	if !ok {
		t.Fatalf("missing entity lens set for %q", entityName)
	}
	if len(kinds) != len(want) {
		t.Fatalf("lens kinds for %q count=%d, want %d", entityName, len(kinds), len(want))
	}
	for _, expectedKind := range want {
		if _, ok := kinds[expectedKind]; !ok {
			t.Fatalf("missing lens kind %q for %q", expectedKind, entityName)
		}
	}
}

func buildTestLSPServerWithSingleFile() (*Server, string) {
	moduleRef := core.ModuleRef{Path: "@"}

	// Shared type expression for interface/component ports in this test fixture.
	portType := ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}}
	interfaceMeta := core.Meta{Start: core.Position{Line: 1, Column: 0}, Stop: core.Position{Line: 1, Column: 7}}
	componentMeta := core.Meta{Start: core.Position{Line: 3, Column: 0}, Stop: core.Position{Line: 3, Column: 12}}
	constMeta := core.Meta{Start: core.Position{Line: 5, Column: 0}, Stop: core.Position{Line: 5, Column: 6}}

	file := src.File{
		Entities: map[string]src.Entity{
			"Greeter": {
				Kind: src.InterfaceEntity,
				Interface: src.Interface{
					IO: src.IO{
						In:  map[string]src.Port{"in": {TypeExpr: portType}},
						Out: map[string]src.Port{"out": {TypeExpr: portType}},
					},
					Meta: interfaceMeta,
				},
			},
			"HelloGreeter": {
				Kind: src.ComponentEntity,
				Component: []src.Component{
					{
						Interface: src.Interface{
							IO: src.IO{
								In:  map[string]src.Port{"in": {TypeExpr: portType}},
								Out: map[string]src.Port{"out": {TypeExpr: portType}},
							},
						},
						Meta: componentMeta,
					},
				},
			},
			"Answer": {
				Kind: src.ConstEntity,
				Const: src.Const{
					Meta: constMeta,
					Value: src.ConstValue{
						Ref: &core.EntityRef{
							Name: "Greeter",
							Meta: core.Meta{Start: core.Position{Line: 5, Column: 12}, Stop: core.Position{Line: 5, Column: 19}},
						},
					},
				},
			},
		},
	}

	build := src.Build{
		Modules: map[core.ModuleRef]src.Module{
			moduleRef: {
				Packages: map[string]src.Package{
					"main": {
						"main": file,
					},
				},
			},
		},
	}

	server := &Server{
		workspacePath: "/tmp/workspace",
		indexMutex:    &sync.Mutex{},
	}
	server.setBuild(build)
	return server, pathToURI("/tmp/workspace/main/main.neva")
}

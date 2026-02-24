package analyzer

import (
	"testing"

	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	"github.com/nevalang/neva/pkg"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeConstBytesLiteralRejected(t *testing.T) {
	a := testAnalyzer(t)
	scope := testScope(t)
	value := "hello"

	_, err := a.analyzeConst(src.Const{
		TypeExpr: ts.Expr{
			Inst: &ts.InstExpr{
				Ref: core.EntityRef{Name: "bytes"},
			},
		},
		Value: src.ConstValue{
			Message: &src.MsgLiteral{
				Str: &value,
			},
		},
	}, scope)

	require.NotNil(t, err)
	require.Contains(t, err.Message, "Bytes constants are not supported")
}

func TestAnalyzeConstStringLiteralAllowed(t *testing.T) {
	a := testAnalyzer(t)
	scope := testScope(t)
	value := "hello"

	_, err := a.analyzeConst(src.Const{
		TypeExpr: ts.Expr{
			Inst: &ts.InstExpr{
				Ref: core.EntityRef{Name: "string"},
			},
		},
		Value: src.ConstValue{
			Message: &src.MsgLiteral{
				Str: &value,
			},
		},
	}, scope)

	require.Nil(t, err)
}

func testAnalyzer(t *testing.T) Analyzer {
	t.Helper()

	terminator := ts.Terminator{}
	checker := ts.MustNewSubtypeChecker(terminator)
	resolver := ts.MustNewResolver(ts.Validator{}, checker, terminator)

	return MustNew(resolver)
}

func testScope(t *testing.T) src.Scope {
	t.Helper()

	mainModRef := core.ModuleRef{Path: "example.com/main"}
	stdModRef := core.ModuleRef{Path: "std", Version: pkg.Version}

	builtinTypes := map[string]src.Entity{
		"any": {
			IsPublic: true,
			Kind:     src.TypeEntity,
			Type:     ts.Def{},
		},
		"bytes": {
			IsPublic: true,
			Kind:     src.TypeEntity,
			Type:     ts.Def{},
		},
		"string": {
			IsPublic: true,
			Kind:     src.TypeEntity,
			Type:     ts.Def{},
		},
	}

	return src.NewScope(src.Build{
		EntryModRef: mainModRef,
		Modules: map[core.ModuleRef]src.Module{
			mainModRef: {
				Packages: map[string]src.Package{
					"main": {
						"main.neva": {
							Imports:  map[string]src.Import{},
							Entities: map[string]src.Entity{},
						},
					},
				},
			},
			stdModRef: {
				Packages: map[string]src.Package{
					"builtin": {
						"types.neva": {
							Imports:  map[string]src.Import{},
							Entities: builtinTypes,
						},
					},
				},
			},
		},
	}, core.Location{
		ModRef:   mainModRef,
		Package:  "main",
		Filename: "main.neva",
	})
}

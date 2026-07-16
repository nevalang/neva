package analyzer

import (
	"testing"

	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	"github.com/nevalang/neva/pkg"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestInferUnionLiteralSenderTypeRejectsPayloadTypeMismatch(t *testing.T) {
	// A typed payload literal must be checked before the connection receiver erases its tag type.
	value := "bad"
	_, err := testAnalyzer(t).inferUnionLiteralSenderType(
		&src.UnionLiteral{
			EntityRef: core.EntityRef{Name: "U"},
			Tag:       "Int",
			Data: &src.ConstValue{
				Message: &src.MsgLiteral{Str: &value},
			},
		},
		nil,
		unionLiteralTestScope(t),
	)

	require.Error(t, err)
	require.Contains(t, err.Message, "Union literal payload type")
}

func unionLiteralTestScope(t *testing.T) src.Scope {
	t.Helper()

	mainModRef := core.ModuleRef{Path: "example.com/main"}
	stdModRef := core.ModuleRef{Path: "std", Version: pkg.Version}
	intExpr := ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}}
	unionExpr := ts.Expr{Lit: &ts.LitExpr{Union: map[string]*ts.Expr{"Int": &intExpr}}}

	return src.NewScope(src.Build{
		EntryModRef: mainModRef,
		Modules: map[core.ModuleRef]src.Module{
			mainModRef: {
				Packages: map[string]src.Package{
					"main": {
						"main.neva": {
							Imports: map[string]src.Import{},
							Entities: map[string]src.Entity{
								"U": {
									Kind: src.TypeEntity,
									Type: ts.Def{BodyExpr: &unionExpr},
								},
							},
						},
					},
				},
			},
			stdModRef: {
				Packages: map[string]src.Package{
					"builtin": {
						"types.neva": {
							Imports: map[string]src.Import{},
							Entities: map[string]src.Entity{
								"int":    {Kind: src.TypeEntity, Type: ts.Def{}},
								"string": {Kind: src.TypeEntity, Type: ts.Def{}},
							},
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

// We need unit tests for parser because it contains not only antlr grammar but also mapping logic.

package parser_test

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/stretchr/testify/require"
)

// TestParser_ParseFile_Directives checks only how compiler directives are parsed.
func TestParser_ParseFile_Directives(t *testing.T) {
	text := []byte(`
		components {
			#runtime_func(d1)
			C1() ()

			#runtime_func(d2)
			C2() () {
				nodes {
					#runtime_func_msg(d3)
					n1 C1

					#runtime_func_msg(d4)
					n2 C1
				}
			}
		}
	`)

	p := parser.New(false)

	got, err := p.ParseFile(text)
	require.NoError(t, err)

	d1 := got.Entities["C1"].Component.Directives[compiler.RuntimeFuncDirective][0]
	require.Equal(t, "d1", d1)

	c2 := got.Entities["C2"].Component

	d2 := c2.Directives[compiler.RuntimeFuncDirective][0]
	require.Equal(t, "d2", d2)

	d3 := c2.Nodes["n1"].Directives[compiler.RuntimeFuncMsgDirective][0]
	require.Equal(t, "d3", d3)

	d4 := c2.Nodes["n2"].Directives[compiler.RuntimeFuncMsgDirective][0]
	require.Equal(t, "d4", d4)
}

// FIXME this suit contains broken tests (they fail because of bugs in parser), but we can't write tests this way because of src.Meta field
// func TestParser_ParseFile(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		bb      []byte
// 		want    src.File
// 		wantErr error
// 	}{
// 		// === Use ===
// 		{
// 			name: "use_statement_with_dots",
// 			bb: []byte(`
// 				use {
// 					std/tmp
// 					github.com/nevalang/neva/pkg/typesystem
// 					some/really/deeply/nested/path/to/local/package/at/the/project
// 				}
// 			`),
// 			want: src.File{
// 				Imports: map[string]src.Import{
// 					"tmp":        {PkgName: "std/tmp"},
// 					"typesystem": {PkgName: "github.com/nevalang/neva/pkg/typesystem"},
// 					"project":    {PkgName: "some/really/deeply/nested/path/to/local/package/at/the/project"},
// 				},
// 				Entities: map[string]src.Entity{},
// 			},
// 		},
// 		{
// 			name: `use_statement_with_"in"_word`, // FIXME keywords collision
// 			bb: []byte(`
// 				use {
// 					package/in/the/project
// 				}
// 			`),
// 			want: src.File{
// 				Imports: map[string]src.Import{
// 					"project": {PkgName: "package/in/the/project"},
// 				},
// 				Entities: map[string]src.Entity{},
// 			},
// 		},
// 		{
// 			name: "inline comment",
// 			bb: []byte(`
// 				use { // inline comment
// 					pkg
// 				}
// 			`),
// 			want: src.File{
// 				Imports: map[string]src.Import{
// 					"pkg": {PkgName: "pkg"},
// 				},
// 				Entities: map[string]src.Entity{},
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "duplicated_imports",
// 			bb: []byte(`
// 				use {
// 					dupl
// 					path/with/parts
// 					withalias dupl
// 					withalias @/local/path/with/parts
// 				}
// 			`),
// 			want: src.File{
// 				Imports: map[string]src.Import{
// 					"dupl":      {PkgName: "dupl"},
// 					"parts":     {PkgName: "path/with/parts"},
// 					"withalias": {PkgName: "@/local/path/with/parts"},
// 				},
// 				Entities: map[string]src.Entity{},
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "several_use_statements",
// 			bb: []byte(`
// 				use {
// 					foo
// 				}
// 				use {
// 					bar
// 				}
// 			`),
// 			want: src.File{
// 				Imports: map[string]src.Import{
// 					"foo": {PkgName: "foo"},
// 					"bar": {PkgName: "bar"},
// 				},
// 				Entities: map[string]src.Entity{},
// 			},
// 			wantErr: nil,
// 		},
// 		// === Interfaces ===
// 		{
// 			name: "just_a_couple_of_simple_interfaces",
// 			bb: []byte(`
// 				interfaces {
// 					IReader(path string) (i int, e err)
// 					IWriter(path) (i int, anything)
// 				}
// 			`),
// 			want: src.File{
// 				Imports: map[string]src.Import{},
// 				Entities: map[string]src.Entity{
// 					"IReader": {
// 						IsPublic: false,
// 						Kind:     src.InterfaceEntity,
// 						Interface: src.Interface{
// 							IO: src.IO{
// 								In: map[string]src.Port{
// 									"path": {
// 										TypeExpr: ts.Expr{
// 											Inst: &ts.InstExpr{
// 												Ref: src.EntityRef{Name: "string"},
// 											},
// 										},
// 										IsArray: false,
// 									},
// 								},
// 								Out: map[string]src.Port{
// 									"i": {
// 										TypeExpr: ts.Expr{
// 											Inst: &ts.InstExpr{
// 												Ref: src.EntityRef{Name: "int"},
// 											},
// 										},
// 										IsArray: false,
// 									},
// 									"e": {
// 										TypeExpr: ts.Expr{
// 											Inst: &ts.InstExpr{
// 												Ref: src.EntityRef{Name: "err"},
// 											},
// 										},
// 										IsArray: false,
// 									},
// 								},
// 							},
// 						},
// 					},
// 					"IWriter": {
// 						IsPublic: false,
// 						Kind:     src.InterfaceEntity,
// 						Interface: src.Interface{
// 							IO: src.IO{
// 								In: map[string]src.Port{
// 									"path": {
// 										TypeExpr: ts.Expr{
// 											Inst: &ts.InstExpr{
// 												Ref: src.EntityRef{Name: "any"},
// 											},
// 										},
// 										IsArray: false,
// 									},
// 								},
// 								Out: map[string]src.Port{
// 									"i": {
// 										TypeExpr: ts.Expr{
// 											Inst: &ts.InstExpr{
// 												Ref: src.EntityRef{Name: "int"},
// 											},
// 										},
// 										IsArray: false,
// 									},
// 									"anything": {
// 										TypeExpr: ts.Expr{
// 											Inst: &ts.InstExpr{
// 												Ref: src.EntityRef{Name: "any"},
// 											},
// 										},
// 										IsArray: false,
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			wantErr: nil,
// 		},
// 		// === Types ===
// 		// Struct type literal expression
// 		{
// 			name: "empty struct",
// 			bb: []byte(`
// 				types {
// 					SomeStruct {}
// 				}
// 			`),
// 			want: src.File{
// 				Imports: map[string]src.Import{},
// 				Entities: map[string]src.Entity{
// 					"SomeStruct": {
// 						IsPublic: false,
// 						Kind:     src.TypeEntity,
// 						Type: ts.Def{
// 							Params: nil,
// 							BodyExpr: &ts.Expr{
// 								Lit: &ts.LitExpr{
// 									Struct: map[string]ts.Expr{},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "struct with one int field",
// 			bb: []byte(`
// 				types {
// 					SomeStruct {
// 						age int
// 					}
// 				}
// 			`),
// 			want: src.File{
// 				Imports: map[string]src.Import{},
// 				Entities: map[string]src.Entity{
// 					"SomeStruct": {
// 						IsPublic: false,
// 						Kind:     src.TypeEntity,
// 						Type: ts.Def{
// 							Params: nil,
// 							BodyExpr: &ts.Expr{
// 								Lit: &ts.LitExpr{
// 									Struct: map[string]ts.Expr{
// 										"age": {
// 											Inst: &ts.InstExpr{Ref: src.EntityRef{Name: "int"}},
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			wantErr: nil,
// 		},
// 		// === Const ===
// 		// Float const
// 		{
// 			name: "float const",
// 			bb: []byte(`
// 				const {
// 					pi float 3.14
// 				}
// 			`),
// 			want: src.File{
// 				Imports: map[string]src.Import{},
// 				Entities: map[string]src.Entity{
// 					"pi": {
// 						IsPublic: false,
// 						Kind:     src.ConstEntity,
// 						Const: src.Const{
// 							Value: &src.Msg{
// 								TypeExpr: ts.Expr{
// 									Inst: &ts.InstExpr{Ref: src.EntityRef{Name: "float"}},
// 								},
// 								Float: utils.Pointer(3.14),
// 							},
// 						},
// 					},
// 				},
// 			},
// 			wantErr: nil,
// 		},
// 		// Components
// 		{
// 			name: "empty_main_component",
// 			bb: []byte(`components {
// 				Main(enter) (exit) {}
// 			}`),
// 			want: src.File{
// 				Entities: map[string]src.Entity{
// 					"Main": {
// 						Kind: src.ComponentEntity,
// 						Component: src.Component{
// 							Interface: src.Interface{
// 								IO: src.IO{
// 									In: map[string]src.Port{
// 										"enter": {
// 											TypeExpr: ts.Expr{
// 												Inst: &ts.InstExpr{Ref: src.EntityRef{Name: "any"}},
// 											},
// 										},
// 									},
// 									Out: map[string]src.Port{
// 										"exit": {
// 											TypeExpr: ts.Expr{
// 												Inst: &ts.InstExpr{Ref: src.EntityRef{Name: "any"}},
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			wantErr: nil,
// 		},
// 	}

// 	p := parser.New(false)

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			got, err := p.ParseFile(tt.bb)
// 			require.Equal(t, tt.want, got)
// 			require.ErrorIs(t, err, tt.wantErr)
// 		})
// 	}
// }

// Package parser implements source code parsing.
package parser_test

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseFile(t *testing.T) {
	tests := []struct {
		name    string
		bb      []byte
		want    src.File
		wantErr error
	}{
		// === Use ===
		{
			name: "use_statement_with_dots",
			bb: []byte(`
				use {
					std/tmp
					github.com/nevalang/neva/pkg/typesystem
					some/really/deeply/nested/path/to/local/package/at/the/project
				}
			`),
			want: src.File{
				Imports: map[string]string{
					"tmp":        "std/tmp",
					"typesystem": "github.com/nevalang/neva/pkg/typesystem",
					"project":    "some/really/deeply/nested/path/to/local/package/at/the/project",
				},
				Entities: map[string]src.Entity{},
			},
		},
		{
			name: `use_statement_with_"in"_word`, // FIXME keywords collision
			bb: []byte(`
				use {
					package/in/the/project
				}
			`),
			want: src.File{
				Imports: map[string]string{
					"project": "package/in/the/project",
				},
				Entities: map[string]src.Entity{},
			},
		},
		{
			name: "inline comment",
			bb: []byte(`
				use { // inline comment
					pkg
				}
			`),
			want: src.File{
				Imports: map[string]string{
					"pkg": "pkg",
				},
				Entities: map[string]src.Entity{},
			},
			wantErr: nil,
		},
		{
			name: "duplicated_imports",
			bb: []byte(`
				use {
					dupl
					path/with/parts
					withalias dupl
					withalias @/local/path/with/parts
				}
			`),
			want: src.File{
				Imports: map[string]string{
					"dupl":      "dupl",
					"parts":     "path/with/parts",
					"withalias": "@/local/path/with/parts",
				},
				Entities: map[string]src.Entity{},
			},
			wantErr: nil,
		},
		{
			name: "several_use_statements",
			bb: []byte(`
				use {
					foo
				}
				use {
					bar
				}
			`),
			want: src.File{
				Imports: map[string]string{
					"foo": "foo",
					"bar": "bar",
				},
				Entities: map[string]src.Entity{},
			},
			wantErr: nil,
		},
		// === Interfaces ===
		{
			name: "just_a_couple_of_simple_interfaces",
			bb: []byte(`
				interfaces {
					IReader(path string) (i int, e err)
					IWriter(path) (i int, anything)
				}
			`),
			want: src.File{
				Imports: map[string]string{},
				Entities: map[string]src.Entity{
					"IReader": {
						Exported: false,
						Kind:     src.InterfaceEntity,
						Interface: src.Interface{
							IO: src.IO{
								In: map[string]src.Port{
									"path": {
										TypeExpr: ts.Expr{
											Inst: &ts.InstExpr{
												Ref: src.EntityRef{Name: "string"},
											},
										},
										IsArray: false,
									},
								},
								Out: map[string]src.Port{
									"i": {
										TypeExpr: ts.Expr{
											Inst: &ts.InstExpr{
												Ref: src.EntityRef{Name: "int"},
											},
										},
										IsArray: false,
									},
									"e": {
										TypeExpr: ts.Expr{
											Inst: &ts.InstExpr{
												Ref: src.EntityRef{Name: "err"},
											},
										},
										IsArray: false,
									},
								},
							},
						},
					},
					"IWriter": {
						Exported: false,
						Kind:     src.InterfaceEntity,
						Interface: src.Interface{
							IO: src.IO{
								In: map[string]src.Port{
									"path": {
										TypeExpr: ts.Expr{
											Inst: &ts.InstExpr{
												Ref: src.EntityRef{Name: "any"},
											},
										},
										IsArray: false,
									},
								},
								Out: map[string]src.Port{
									"i": {
										TypeExpr: ts.Expr{
											Inst: &ts.InstExpr{
												Ref: src.EntityRef{Name: "int"},
											},
										},
										IsArray: false,
									},
									"anything": {
										TypeExpr: ts.Expr{
											Inst: &ts.InstExpr{
												Ref: src.EntityRef{Name: "any"},
											},
										},
										IsArray: false,
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		// === Types ===
		// Struct type literal expression
		{
			name: "empty struct",
			bb: []byte(`
				types {
					SomeStruct {}
				}
			`),
			want: src.File{
				Imports: map[string]string{},
				Entities: map[string]src.Entity{
					"SomeStruct": {
						Exported: false,
						Kind:     src.TypeEntity,
						Type: ts.Def{
							Params: nil,
							BodyExpr: &ts.Expr{
								Lit: &ts.LitExpr{
									Rec: map[string]ts.Expr{},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "struct with one int field",
			bb: []byte(`
				types {
					SomeStruct {
						age int
					}
				}
			`),
			want: src.File{
				Imports: map[string]string{},
				Entities: map[string]src.Entity{
					"SomeStruct": {
						Exported: false,
						Kind:     src.TypeEntity,
						Type: ts.Def{
							Params: nil,
							BodyExpr: &ts.Expr{
								Lit: &ts.LitExpr{
									Rec: map[string]ts.Expr{
										"age": {
											Inst: &ts.InstExpr{Ref: src.EntityRef{Name: "int"}},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		// === Const ===
		// Float const
		{
			name: "float const",
			bb: []byte(`
				const {
					pi float 3.14
				}
			`),
			want: src.File{
				Imports: map[string]string{},
				Entities: map[string]src.Entity{
					"pi": {
						Exported: false,
						Kind:     src.ConstEntity,
						Const: src.Const{
							Value: &src.Msg{
								TypeExpr: ts.Expr{
									Inst: &ts.InstExpr{Ref: src.EntityRef{Name: "float"}},
								},
								Float: 3.14,
							},
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	p := parser.MustNew(false)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := p.ParseFile(context.Background(), tt.bb)
			require.Equal(t, tt.want, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

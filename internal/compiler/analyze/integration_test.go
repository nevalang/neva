//go:build integration
// +build integration

package analyze_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/emil14/neva/internal/compiler/analyze"
	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var h src.Helper

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	type testcase struct {
		name    string
		prog    src.Prog
		wantErr error
	}

	tests := []testcase{
		// {
		// 	name: "root pkg refers to type and component in another pkg",
		// 	prog: src.Prog{
		// 		Pkgs: map[string]src.Pkg{
		// 			"pkg2": {
		// 				Entities: map[string]src.Entity{
		// 					"t1": h.TypeEntity(
		// 						true,
		// 						h.Def( // type t1<a> = vec<a>
		// 							h.Inst("vec", h.Inst("a")),
		// 							h.ParamWithNoConstr("a"),
		// 						),
		// 					),
		// 					"c1": {
		// 						Exported: true,
		// 						Kind:     src.ComponentEntity,
		// 					},
		// 				},
		// 			},
		// 			"pkg1": {
		// 				Imports: h.Imports("pkg2"),
		// 				Entities: map[string]src.Entity{
		// 					"t1": h.TypeEntity(
		// 						true,
		// 						h.Def( // type t1 = pkg2.t1<int>
		// 							h.Inst("pkg2.t1", h.Inst("int")),
		// 						),
		// 					),
		// 					"c1": h.RootComponentEntity(map[string]src.Node{
		// 						"n1": h.ComponentNode("pkg1", "c1"),
		// 					}),
		// 				},
		// 				RootComponent: "c1",
		// 			},
		// 		},
		// 		RootPkg: "pkg1",
		// 	},
		// 	wantErr: nil,
		// },
		{
			name: "root pkg refers another pkg that refers another pkg via types",
			prog: src.Prog{
				Pkgs: map[string]src.Pkg{
					"pkg3": {
						Entities: map[string]src.Entity{
							"t1": h.TypeEntity(
								true,
								h.Def( // type t1<a> = vec<a>
									h.Inst("vec", h.Inst("a")),
									h.ParamWithNoConstr("a"),
								),
							),
						},
					},
					"pkg2": {
						Imports: h.Imports("pkg3"),
						Entities: map[string]src.Entity{
							"t1": h.TypeEntity(
								true,
								h.Def( // type t1<a> = t1<a>
									h.Inst("pkg3.t1", h.Inst("a")),
									h.ParamWithNoConstr("a"),
								),
							),
							"c1": {
								Exported: true,
								Kind:     src.ComponentEntity,
							},
						},
					},
					"pkg1": {
						Imports: h.Imports("pkg2"),
						Entities: map[string]src.Entity{
							"t1": h.TypeEntity(
								true,
								h.Def( // type t1 = pkg2.t1<int>
									h.Inst("pkg2.t1", h.Inst("int")),
								),
							),
							"c1": h.RootComponentEntity(map[string]src.Node{
								"n1": h.ComponentNode("pkg1", "c1"),
							}),
						},
						RootComponent: "c1",
					},
				},
				RootPkg: "pkg1",
			},
			wantErr: nil,
		},
		// {
		// 	name: "inassignable message",
		// 	prog: src.Prog{
		// 		Pkgs: map[string]src.Pkg{
		// 			"pkg1": {
		// 				Entities: map[string]src.Entity{
		// 					"m1": h.IntMsgEntity(true, 42),
		// 					"c1": {
		// 						Exported: true,
		// 						Kind:     src.ComponentEntity,
		// 					},
		// 				},
		// 			},
		// 			"pkg2": {
		// 				Imports: h.Imports("pkg1"),
		// 				Entities: map[string]src.Entity{
		// 					"m1": h.MsgWithRefEntity(true, &src.EntityRef{
		// 						Pkg:  "pkg1",
		// 						Name: "m1",
		// 					}),
		// 				},
		// 			},
		// 			"pkg3": {
		// 				Imports: h.Imports("pkg1", "pkg2"),
		// 				Entities: map[string]src.Entity{
		// 					"m1": h.IntVecMsgEntity(
		// 						true,
		// 						[]src.Msg{
		// 							{
		// 								Ref: &src.EntityRef{
		// 									Pkg:  "pkg1",
		// 									Name: "m1",
		// 								},
		// 							},
		// 							{
		// 								Ref: &src.EntityRef{
		// 									Pkg:  "pkg2",
		// 									Name: "m1",
		// 								},
		// 							},
		// 							{Value: h.IntMsgValue(43)},
		// 						},
		// 					),
		// 				},
		// 			},
		// 			"pkg4": {
		// 				Imports: h.Imports("pkg1", "pkg2", "pkg3"),
		// 				Entities: map[string]src.Entity{
		// 					"m1": h.IntVecMsgEntity(
		// 						true,
		// 						[]src.Msg{
		// 							{Value: h.IntMsgValue(44)},
		// 							{
		// 								Ref: &src.EntityRef{
		// 									Pkg:  "pkg3",
		// 									Name: "m1",
		// 								},
		// 							},
		// 						},
		// 					),
		// 					"c1": h.RootComponentEntity(map[string]src.Node{
		// 						"n1": h.ComponentNode("pkg1", "c1"),
		// 					}),
		// 				},
		// 				RootComponent: "c1",
		// 			},
		// 		},
		// 		RootPkg: "pkg4",
		// 	},
		// 	wantErr: analyze.ErrVecEl,
		// },
	}

	a := analyze.Analyzer{
		Resolver: ts.NewDefaultResolver(),
		Compator: ts.NewDefaultCompatChecker(),
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := a.Analyze(context.Background(), tt.prog)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

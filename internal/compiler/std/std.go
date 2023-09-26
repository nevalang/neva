// TODO maybe load std.neva here without analyzer?
package std

import (
	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/ts"
)

// New returns single-package std library with only runtime-function interfaces. TODO add HOCs.
func New() src.File { //nolint:funlen
	return src.File{
		Entities: map[string]src.Entity{
			"Read": {
				Exported: true,
				Kind:     src.ComponentEntity,
				Component: src.Component{
					Interface: src.Interface{
						Params: []ts.Param{},
						IO: src.IO{
							In: map[string]src.Port{
								"sig": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "any"},
									},
								},
							},
							Out: map[string]src.Port{
								"v": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "str"},
									},
								},
							},
						},
					},
				},
			},
			"Print": {
				Exported: true,
				Kind:     src.ComponentEntity,
				Component: src.Component{
					Interface: src.Interface{
						Params: []ts.Param{
							{Name: "t1"},
						},
						IO: src.IO{
							In: map[string]src.Port{
								"v": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "t1"},
									},
								},
							},
							Out: map[string]src.Port{
								"v": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "t1"},
									},
								},
							},
						},
					},
				},
			},
			"Lock": {
				Exported: true,
				Kind:     src.ComponentEntity,
				Component: src.Component{
					Interface: src.Interface{
						Params: []ts.Param{
							{Name: "t1"},
						},
						IO: src.IO{
							In: map[string]src.Port{
								"v": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "t1"},
									},
								},
								"sig": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "any"},
									},
								},
							},
							Out: map[string]src.Port{
								"v": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "t1"},
									},
								},
							},
						},
					},
				},
			},
			"Const": {
				Exported: true,
				Kind:     src.ComponentEntity,
				Component: src.Component{
					Interface: src.Interface{
						Params: []ts.Param{
							{Name: "t1"},
						},
						IO: src.IO{
							In: map[string]src.Port{},
							Out: map[string]src.Port{
								"v": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "t1"},
									},
								},
							},
						},
					},
				},
			},
			"Add": {
				Exported: true,
				Kind:     src.ComponentEntity,
				Component: src.Component{
					Interface: src.Interface{
						Params: []ts.Param{
							{
								Name: "t1",
								Constr: ts.Expr{
									Lit: ts.LitExpr{
										Union: []ts.Expr{
											{
												Inst: ts.InstExpr{Ref: "int"},
											},
											{
												Inst: ts.InstExpr{Ref: "float"},
											},
											{
												Inst: ts.InstExpr{Ref: "str"},
											},
										},
									},
								},
							},
						},
						IO: src.IO{
							In: map[string]src.Port{
								"a": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "t1"},
									},
								},
								"b": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "t1"},
									},
								},
							},
							Out: map[string]src.Port{
								"v": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "t1"},
									},
								},
							},
						},
					},
				},
			},
			"ParseNum": {
				Exported: true,
				Kind:     src.ComponentEntity,
				Component: src.Component{
					Interface: src.Interface{
						Params: []ts.Param{
							{
								Name: "t1",
								Constr: ts.Expr{
									Lit: ts.LitExpr{
										Union: []ts.Expr{
											{
												Inst: ts.InstExpr{Ref: "int"},
											},
											{
												Inst: ts.InstExpr{Ref: "float"},
											},
										},
									},
								},
							},
						},
						IO: src.IO{
							In: map[string]src.Port{
								"v": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "str"},
									},
								},
							},
							Out: map[string]src.Port{
								"v": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "t1"},
									},
								},
								"err": {
									Type: ts.Expr{
										Inst: ts.InstExpr{Ref: "str"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

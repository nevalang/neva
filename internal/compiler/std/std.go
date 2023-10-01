// TODO maybe load std.neva here without analyzer?
package std

import (
	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
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
						TypeParams: []ts.Param{},
						IO: src.IO{
							In: map[string]src.Port{
								"sig": {},
							},
							Out: map[string]src.Port{
								"v": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "str"},
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
						TypeParams: []ts.Param{
							{Name: "t1"},
						},
						IO: src.IO{
							In: map[string]src.Port{
								"v": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "t1"},
									},
								},
							},
							Out: map[string]src.Port{
								"v": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "t1"},
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
						TypeParams: []ts.Param{
							{Name: "t1"},
						},
						IO: src.IO{
							In: map[string]src.Port{
								"v": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "t1"},
									},
								},
								"sig": {},
							},
							Out: map[string]src.Port{
								"v": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "t1"},
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
						TypeParams: []ts.Param{
							{Name: "t1"},
						},
						IO: src.IO{
							In: map[string]src.Port{},
							Out: map[string]src.Port{
								"v": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "t1"},
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
						TypeParams: []ts.Param{
							{
								Name: "t1",
								Constr: &ts.Expr{
									Lit: &ts.LitExpr{
										Union: []ts.Expr{
											{
												Inst: &ts.InstExpr{Ref: "int"},
											},
											{
												Inst: &ts.InstExpr{Ref: "float"},
											},
											{
												Inst: &ts.InstExpr{Ref: "str"},
											},
										},
									},
								},
							},
						},
						IO: src.IO{
							In: map[string]src.Port{
								"a": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "t1"},
									},
								},
								"b": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "t1"},
									},
								},
							},
							Out: map[string]src.Port{
								"v": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "t1"},
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
						TypeParams: []ts.Param{
							{
								Name: "t1",
								Constr: &ts.Expr{
									Lit: &ts.LitExpr{
										Union: []ts.Expr{
											{
												Inst: &ts.InstExpr{Ref: "int"},
											},
											{
												Inst: &ts.InstExpr{Ref: "float"},
											},
										},
									},
								},
							},
						},
						IO: src.IO{
							In: map[string]src.Port{
								"v": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "str"},
									},
								},
							},
							Out: map[string]src.Port{
								"v": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "t1"},
									},
								},
								"err": {
									TypeExpr: ts.Expr{
										Inst: &ts.InstExpr{Ref: "str"},
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

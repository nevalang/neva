package std

import (
	"github.com/nevalang/neva/internal/src"
	"github.com/nevalang/neva/pkg/ts"
)

// New returns single-package std library with only runtime-function interfaces. TODO add HOCs.
func New() src.Package {
	return src.Package{
		Entities: map[string]src.Entity{
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
		},
	}
}

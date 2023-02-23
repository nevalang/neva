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

func TestDefaultResolver(t *testing.T) {
	t.Parallel()

	type testcase struct {
		name    string
		prog    src.Prog
		wantErr error
	}

	var h ts.Helper

	tests := []testcase{
		{
			name: "",
			prog: src.Prog{
				Pkgs: map[string]src.Pkg{
					"pkg_1": {
						Entities: map[string]src.Entity{
							"t1": {
								Exported: true,
								Kind:     src.TypeEntity,
								Type:     h.Def(h.Inst("vec", h.Inst("a")), h.ParamWithNoConstr("a")), // type t1<a> = vec<a>
							},
							"c1": {
								Exported: true,
								Kind:     src.ComponentEntity,
							},
						},
						RootComponent: "",
					},
					"pkg_2": {
						Imports: map[string]string{
							"pkg_1": "pkg_1",
						},
						Entities: map[string]src.Entity{
							"t1": {
								Exported: true,
								Kind:     src.TypeEntity,
								Type:     h.Def(h.Inst("pkg_1.t1", h.Inst("int"))), // type t1 = pkg_1.t1<int>
							},
							"c1": {
								Kind: src.ComponentEntity,
								Component: src.Component{
									TypeParams: []ts.Param{
										h.ParamWithNoConstr("t"),
									},
									IO: src.IO{
										In: map[string]src.Port{
											"sig": {
												Type: h.Inst("t"),
											},
										},
									},
									Nodes: map[string]src.Node{
										"n1": {
											Instance: src.Instance{
												Ref: src.EntityRef{
													Pkg:  "pkg_1",
													Name: "c1",
												},
											},
										},
									},
								},
							},
						},
						RootComponent: "c1",
					},
				},
				RootPkg: "pkg_2",
			},
			wantErr: nil,
		},
	}

	a := analyze.Analyzer{
		Resolver: ts.NewDefaultResolver(),
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

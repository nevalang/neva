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

func TestDefaultResolver(t *testing.T) {
	t.Parallel()

	type testcase struct {
		name    string
		prog    src.Prog
		wantErr error
	}

	tests := []testcase{
		{
			name: "",
			prog: src.Prog{
				Pkgs: map[string]src.Pkg{
					"pkg_2": {
						Entities: map[string]src.Entity{
							"t1": h.TypeEntity(
								true,
								h.Def( // type t1<a> = vec<a>
									h.Inst("vec", h.Inst("a")),
									h.ParamWithNoConstr("a"),
								),
							),
							"c1": {
								Exported: true,
								Kind:     src.ComponentEntity,
							},
						},
						RootComponent: "",
					},
					"pkg_1": {
						Imports: h.Imports("pkg_2"),
						Entities: map[string]src.Entity{
							"t1": h.TypeEntity(
								true,
								h.Def( // type t1 = pkg_2.t1<int>
									h.Inst("pkg_2.t1", h.Inst("int")),
								),
							),
							"c1": h.RootComponentEntity(map[string]src.Node{
								"n1": {
									Instance: src.Instance{
										Ref: src.EntityRef{
											Pkg:  "pkg_1",
											Name: "c1",
										},
									},
								},
							}),
						},
						RootComponent: "c1",
					},
				},
				RootPkg: "pkg_1",
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

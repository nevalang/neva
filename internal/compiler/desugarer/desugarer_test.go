package desugarer

import (
	"testing"

	"github.com/nevalang/neva/internal/utils"
	src "github.com/nevalang/neva/pkg/sourcecode"
	"github.com/nevalang/neva/pkg/typesystem"
	"github.com/stretchr/testify/require"
)

func TestDesugarer_Desugar(t *testing.T) {
	tests := []struct {
		name    string
		build   src.Build
		want    src.Build
		wantErr error
	}{
		{
			name: "inject_std_dep_and_const_node",
			build: src.Build{
				Modules: map[src.ModuleRef]src.Module{
					{}: {
						Manifest: src.ModuleManifest{WantCompilerVersion: "0.0.1"},
						Packages: map[string]src.Package{
							"main": {
								"file": src.File{
									Entities: map[string]src.Entity{
										"bar": { // const must be present so desugarer can figure out type args for Const node
											Kind: src.ConstEntity,
											Const: src.Const{
												Value: &src.Msg{
													TypeExpr: typesystem.Expr{
														Inst: &typesystem.InstExpr{Ref: src.EntityRef{Name: "int"}},
													},
													Int: utils.Pointer(42),
												},
											},
										},
										"Main": {
											Kind: src.ComponentEntity,
											Component: src.Component{
												Net: []src.Connection{
													{
														SenderSide: src.SenderConnectionSide{
															ConstRef: &src.EntityRef{Name: "bar"},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: src.Build{
				Modules: map[src.ModuleRef]src.Module{
					{}: {
						Manifest: src.ModuleManifest{
							WantCompilerVersion: "0.0.1",
							Deps: map[string]src.ModuleRef{
								"std": { // stdlib dep injected
									Path:    "std",
									Version: "0.0.1",
								},
							},
						},
						Packages: map[string]src.Package{
							"main": {
								"file": src.File{
									Entities: map[string]src.Entity{
										"bar": {
											Kind: src.ConstEntity,
											Const: src.Const{
												Value: &src.Msg{
													TypeExpr: typesystem.Expr{
														Inst: &typesystem.InstExpr{Ref: src.EntityRef{Name: "int"}},
													},
													Int: utils.Pointer(42),
												},
											},
										},
										"Main": {
											Kind: src.ComponentEntity,
											Component: src.Component{
												Nodes: map[string]src.Node{
													"bar": {
														Directives: map[src.Directive][]string{
															"runtime_func_msg": {"bar"},
														},
														EntityRef: src.EntityRef{
															Pkg:  "std/builtin",
															Name: "Const",
														},
														TypeArgs: []typesystem.Expr{
															{
																Inst: &typesystem.InstExpr{Ref: src.EntityRef{Name: "int"}},
															},
														},
													},
												},
												Net: []src.Connection{
													{
														SenderSide: src.SenderConnectionSide{
															PortAddr: &src.PortAddr{
																Node: "bar",
																Port: "out",
															},
														},
													},
												},
											},
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
	}

	d := Desugarer{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.Desugar(tt.build)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

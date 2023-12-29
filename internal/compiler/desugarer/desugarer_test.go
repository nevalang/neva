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
		// every module must have std module dependency and std/builtin import in every file
		// {
		// 	name: "std_mod_dep_and_builtin_import",
		// 	build: src.Build{
		// 		Modules: map[src.ModuleRef]src.Module{
		// 			{}: {
		// 				Manifest: src.ModuleManifest{
		// 					Deps: map[string]src.ModuleRef{}, // <-- no std mod dep
		// 				},
		// 				Packages: map[string]src.Package{
		// 					"pkgName": {
		// 						"fileName": src.File{
		// 							Imports: map[string]src.Import{}, // <-- no imports of std/builtin
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	want: src.Build{
		// 		Modules: map[src.ModuleRef]src.Module{
		// 			{}: {
		// 				Manifest: src.ModuleManifest{
		// 					Deps: map[string]src.ModuleRef{
		// 						"std": {Path: "std", Version: "0.0.1"}, // <-- std mod dep added
		// 					},
		// 				},
		// 				Packages: map[string]src.Package{
		// 					"pkgName": {
		// 						"fileName": src.File{
		// 							Imports: map[string]src.Import{
		// 								"builtin": { // <-- std/builtin import added
		// 									ModuleName: "std",
		// 									PkgName:    "builtin",
		// 								},
		// 							},
		// 							Entities: map[string]src.Entity{},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantErr: nil,
		// },
		{
			// every network with const ref must be desugared into special node and connections to it
			name: "const_node",
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
								"std": { // <-- stdlib mod dep added
									Path:    "std",
									Version: "0.0.1",
								},
							},
						},
						Packages: map[string]src.Package{
							"main": {
								"file": src.File{
									Imports: map[string]src.Import{
										"builtin": { // <-- std/builtin import added
											ModuleName: "std",
											PkgName:    "builtin",
										},
									},
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
													"bar": { // <-- const node added
														Directives: map[src.Directive][]string{
															"runtime_func_msg": {"bar"},
														},
														EntityRef: src.EntityRef{
															Pkg:  "builtin",
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
														SenderSide: src.SenderConnectionSide{ // <-- const ref conn replaced with normal one
															PortAddr: &src.PortAddr{
																Node: "bar",
																Port: "v",
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
		// {
		// 	name: "inject_void_nodes_and_connections",
		// 	build: src.Build{
		// 		Modules: map[src.ModuleRef]src.Module{
		// 			{}: {
		// 				Manifest: src.ModuleManifest{},
		// 				Packages: map[string]src.Package{
		// 					"main": {
		// 						"file": src.File{
		// 							Entities: map[string]src.Entity{
		// 								"Foo": {
		// 									Kind: src.ComponentEntity,
		// 									Component: src.Component{
		// 										Nodes: map[string]src.Node{
		// 											"bar": {EntityRef: src.EntityRef{Name: "Bar"}}, // <-- node with `x` outport
		// 										},
		// 										Net: []src.Connection{}, // <-- no bar.x usage
		// 									},
		// 								},
		// 								"Bar": {
		// 									Kind: src.ComponentEntity,
		// 									Component: src.Component{
		// 										Interface: src.Interface{
		// 											IO: src.IO{
		// 												Out: map[string]src.Port{
		// 													"x": {}, // <-- unused by Foo
		// 												},
		// 											},
		// 										},
		// 									},
		// 								},
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	want: src.Build{
		// 		Modules: map[src.ModuleRef]src.Module{
		// 			{}: {
		// 				Manifest: src.ModuleManifest{
		// 					Deps: map[string]src.ModuleRef{
		// 						"std": {
		// 							Path:    "std",
		// 							Version: "0.0.1",
		// 						},
		// 					},
		// 				},
		// 				Packages: map[string]src.Package{
		// 					"main": {
		// 						"file": src.File{
		// 							Entities: map[string]src.Entity{
		// 								"Foo": {
		// 									Kind: src.ComponentEntity,
		// 									Component: src.Component{
		// 										Nodes: map[string]src.Node{
		// 											"bar": {EntityRef: src.EntityRef{Name: "Bar"}}, // that one node
		// 											"void": {
		// 												EntityRef: src.EntityRef{
		// 													Name: "Void",
		// 													Pkg:  "builtin",
		// 												},
		// 											},
		// 										},
		// 										Net: []src.Connection{
		// 											{},
		// 										}, // <-- no bar.x usage
		// 									},
		// 								},
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantErr: nil,
		// },
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

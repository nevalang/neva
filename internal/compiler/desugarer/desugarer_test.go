package desugarer

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/utils"
	src "github.com/nevalang/neva/pkg/sourcecode"
	"github.com/nevalang/neva/pkg/typesystem"
	"github.com/stretchr/testify/require"
)

var p = parser.New(false)

func TestDesugarer_Desugar(t *testing.T) {
	tests := []struct {
		name    string
		build   src.Build
		want    src.Build
		wantErr bool
	}{
		// every module must have std module dependency and std/builtin import in every file
		{
			name: "std_mod_dep_and_builtin_import",
			build: src.Build{
				Modules: map[src.ModuleRef]src.Module{
					{}: {
						Manifest: src.ModuleManifest{
							Deps: map[string]src.ModuleRef{}, // <-- no std mod dep
						},
						Packages: map[string]src.Package{
							"pkgName": {
								"fileName": src.File{
									Imports: map[string]src.Import{}, // <-- no imports of std/builtin
								},
							},
						},
					},
				},
			},
			want: src.Build{
				Modules: map[src.ModuleRef]src.Module{
					{}: {
						Manifest: defaultManifest(), // <-- std mod dep
						Packages: map[string]src.Package{
							"pkgName": {
								"fileName": src.File{
									Imports:  defaultImports(), // <-- std/builtin import
									Entities: map[string]src.Entity{},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		// every network with const ref must be desugared into special node and connections to it
		{
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
													"__bar__": { // <-- const node added
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
																Node: "__bar__",
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
			wantErr: false,
		},
		// every unused outport must be connected to special void node
		{
			name: "void_node",
			build: src.Build{
				Modules: map[src.ModuleRef]src.Module{
					{}: {
						Manifest: src.ModuleManifest{},
						Packages: map[string]src.Package{
							"main": {
								"file": src.File{
									Entities: map[string]src.Entity{
										"Foo": {
											Kind: src.ComponentEntity,
											Component: src.Component{
												Nodes: map[string]src.Node{
													"bar": {EntityRef: src.EntityRef{Name: "Bar"}}, // <-- node with `x` outport
												},
												Net: []src.Connection{}, // <-- no bar.x usage
											},
										},
										"Bar": {
											Kind: src.ComponentEntity,
											Component: src.Component{
												Interface: src.Interface{
													IO: src.IO{
														Out: map[string]src.Port{
															"x": {}, // <-- unused by Foo
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
							Deps: map[string]src.ModuleRef{
								"std": {
									Path:    "std",
									Version: "0.0.1",
								},
							},
						},
						Packages: map[string]src.Package{
							"main": {
								"file": src.File{
									Imports: map[string]src.Import{
										"builtin": {
											ModuleName: "std",
											PkgName:    "builtin",
										},
									},
									Entities: map[string]src.Entity{
										"Foo": {
											Kind: src.ComponentEntity,
											Component: src.Component{
												Nodes: map[string]src.Node{
													"bar": {EntityRef: src.EntityRef{Name: "Bar"}}, // that one node
													"__void__": { // <-- new node
														EntityRef: src.EntityRef{
															Name: "Void",
															Pkg:  "builtin",
														},
													},
												},
												Net: []src.Connection{
													{ // <-- (bar.x -> void.void)
														SenderSide: src.SenderConnectionSide{
															PortAddr: &src.PortAddr{
																Node: "bar",
																Port: "x",
															},
														},
														ReceiverSides: []src.ReceiverConnectionSide{
															{
																PortAddr: src.PortAddr{
																	Node: "__void__",
																	Port: "v",
																},
															},
														},
													},
												},
											},
										},
										"Bar": {
											Kind: src.ComponentEntity,
											Component: src.Component{
												Interface: src.Interface{
													IO: src.IO{
														Out: map[string]src.Port{
															"x": {}, // <-- now used with void
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
			wantErr: false,
		},
	}

	d := Desugarer{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.Desugar(tt.build)
			require.Equal(t, err != nil, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func mustParseFile(s string) src.File {
	file, err := p.ParseFile([]byte(s))
	if err != nil {
		panic(err)
	}
	return file
}

func injectStd(file src.File) src.File {
	if file.Imports == nil {
		file.Imports = map[string]src.Import{}
	}
	file.Imports["builtin"] = src.Import{
		ModuleName: "std",
		PkgName:    "builtin",
	}
	return file
}

func defaultManifest() src.ModuleManifest {
	return src.ModuleManifest{
		Deps: map[string]src.ModuleRef{
			"std": {
				Path:    "std",
				Version: "0.0.1",
			},
		},
	}
}

func defaultImports() map[string]src.Import {
	return map[string]src.Import{
		"builtin": {
			ModuleName: "std",
			PkgName:    "builtin",
		},
	}
}

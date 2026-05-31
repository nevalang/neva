package analyzer

import (
	"strings"
	"testing"

	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	"github.com/nevalang/neva/pkg"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeNodeRequiresExplicitTypeArgs(t *testing.T) {
	build := analyzerBuildWithGenericNode(nil)

	_, err := testAnalyzer(t).Analyze(build, "main")
	require.NotNil(t, err)
	msg := err.Error()
	require.True(t,
		strings.Contains(msg, "count of arguments mismatch count of parameters, want 1 got 0") ||
			strings.Contains(msg, "Referenced node not found"),
		"unexpected analyzer error: %s", msg,
	)
}

func analyzerBuildWithGenericNode(nodeTypeArgs src.TypeArgs) src.Build {
	mainModRef := core.ModuleRef{Path: "example.com/main"}
	stdModRef := core.ModuleRef{Path: "std", Version: pkg.Version}

	anyExpr := ts.Expr{
		Inst: &ts.InstExpr{
			Ref: core.EntityRef{Name: "any"},
		},
	}
	genericExpr := ts.Expr{
		Inst: &ts.InstExpr{
			Ref: core.EntityRef{Name: "T"},
		},
	}

	return src.Build{
		EntryModRef: mainModRef,
		Modules: map[core.ModuleRef]src.Module{
			mainModRef: {
				Manifest: src.ModuleManifest{LanguageVersion: pkg.Version},
				Packages: map[string]src.Package{
					"main": {
						"main.neva": {
							Imports: map[string]src.Import{},
							Entities: map[string]src.Entity{
								"Main": {
									Kind:     src.ComponentEntity,
									IsPublic: false,
									Component: []src.Component{
										{
											Interface: src.Interface{
												IO: src.IO{
													In: map[string]src.Port{
														"start": {TypeExpr: anyExpr},
													},
													Out: map[string]src.Port{
														"stop": {TypeExpr: anyExpr},
													},
												},
											},
											Nodes: map[string]src.Node{
												"gen": {
													EntityRef: core.EntityRef{Name: "Generic"},
													TypeArgs:  nodeTypeArgs,
												},
											},
											Net: []src.Connection{
												{
													Senders: []src.ConnectionSender{{
														PortAddr: &src.PortAddr{Port: "start"},
													}},
													Receivers: []src.ConnectionReceiver{{
														PortAddr: &src.PortAddr{
															Node: "gen",
															Port: "in",
														},
													}},
												},
												{
													Senders: []src.ConnectionSender{{
														PortAddr: &src.PortAddr{
															Node: "gen",
															Port: "out",
														},
													}},
													Receivers: []src.ConnectionReceiver{{
														PortAddr: &src.PortAddr{Port: "stop"},
													}},
												},
											},
										},
									},
								},
								"Generic": {
									Kind: src.ComponentEntity,
									Component: []src.Component{
										{
											Interface: src.Interface{
												TypeParams: src.TypeParams{
													Params: []ts.Param{
														{
															Name:   "T",
															Constr: anyExpr,
														},
													},
												},
												IO: src.IO{
													In: map[string]src.Port{
														"in": {TypeExpr: genericExpr},
													},
													Out: map[string]src.Port{
														"out": {TypeExpr: genericExpr},
													},
												},
											},
											Net: []src.Connection{
												{
													Senders: []src.ConnectionSender{{
														PortAddr: &src.PortAddr{Port: "in"},
													}},
													Receivers: []src.ConnectionReceiver{{
														PortAddr: &src.PortAddr{Port: "out"},
													}},
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
			stdModRef: {
				Manifest: src.ModuleManifest{LanguageVersion: pkg.Version},
				Packages: map[string]src.Package{
					"builtin": {
						"types.neva": {
							Imports: map[string]src.Import{},
							Entities: map[string]src.Entity{
								"any": {
									Kind: src.TypeEntity,
									Type: ts.Def{},
								},
							},
						},
					},
				},
			},
		},
	}
}

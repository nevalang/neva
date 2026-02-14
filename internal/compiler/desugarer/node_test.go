package desugarer

import (
	"testing"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestDesugarNodes_preserves_err_guard_for_anon_di(t *testing.T) {
	modRef := core.ModuleRef{Path: "@"}

	location := core.Location{
		ModRef:   modRef,
		Package:  "pkg",
		Filename: "file",
	}

	tapComponent := src.Component{
		Nodes: map[string]src.Node{
			"handler": {
				EntityRef: core.EntityRef{Name: "IHandler"},
			},
		},
	}

	build := src.Build{
		EntryModRef: modRef,
		Modules: map[core.ModuleRef]src.Module{
			modRef: {
				Manifest: src.ModuleManifest{Deps: map[string]core.ModuleRef{}},
				Packages: map[string]src.Package{
					"pkg": {
						"file": {
							Entities: map[string]src.Entity{
								"Tap": {
									Kind:      src.ComponentEntity,
									Component: []src.Component{tapComponent},
								},
								"IHandler": {
									Kind:      src.InterfaceEntity,
									Interface: src.Interface{},
								},
							},
						},
					},
				},
			},
		},
	}

	scope := src.NewScope(build, location)
	overloadIdx := 1
	component := src.Component{
		Nodes: map[string]src.Node{
			"tap": {
				EntityRef: core.EntityRef{Name: "Tap"},
				ErrGuard:  true,
				DIArgs: map[string]src.Node{
					"": {EntityRef: core.EntityRef{Name: "HandlerImpl"}},
				},
				OverloadIndex: &overloadIdx,
				Meta:          core.Meta{Location: location},
			},
		},
	}

	d := Desugarer{}
	desugaredNodes, _, err := d.desugarNodes(component, scope, map[string]src.Entity{})
	require.NoError(t, err)

	node, ok := desugaredNodes["tap"]
	require.True(t, ok)
	require.True(t, node.ErrGuard)
	require.NotNil(t, node.OverloadIndex)
	require.Equal(t, overloadIdx, *node.OverloadIndex)

	_, hasNamed := node.DIArgs["handler"]
	require.True(t, hasNamed)
}

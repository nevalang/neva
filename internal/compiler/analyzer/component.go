package analyzer

import (
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

func (a Analyzer) analyzeComponent(
	component src.Component,
	scope src.Scope,
) (src.Component, *compiler.Error) {
	externArgs, hasExtern := component.Directives[compiler.ExternDirective]
	if hasExtern && len(externArgs) == 0 {
		return src.Component{}, &compiler.Error{
			Message: "Component that use #extern directive must provide at least one argument",
			Meta:    &component.Meta,
		}
	}

	resolvedInterface, err := a.analyzeInterface(
		component.Interface,
		scope,
		analyzeInterfaceParams{
			allowEmptyInports:  hasExtern,
			allowEmptyOutports: hasExtern,
		},
	)
	if err != nil {
		return src.Component{}, compiler.Error{
			Meta: &component.Meta,
		}.Wrap(err)
	}

	if hasExtern {
		if len(component.Nodes) != 0 || len(component.Net) != 0 {
			return src.Component{}, &compiler.Error{
				Message: "Component with nodes or network cannot use #extern directive",
				Meta:    &component.Meta,
			}
		}
		return component, nil
	}

	resolvedNodes, nodesIfaces, hasGuard, err := a.analyzeNodes(
		component.Interface,
		component.Nodes,
		scope,
	)
	if err != nil {
		return src.Component{}, compiler.Error{
			Meta: &component.Meta,
		}.Wrap(err)
	}

	if len(component.Net) == 0 {
		return src.Component{}, &compiler.Error{
			Message: "Component must have network",
			Meta:    &component.Meta,
		}
	}

	analyzedNet, err := a.analyzeNetwork(
		component.Net,
		resolvedInterface,
		hasGuard,
		resolvedNodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return src.Component{}, compiler.Error{
			Meta: &component.Meta,
		}.Wrap(err)
	}

	return src.Component{
		Interface: resolvedInterface,
		Nodes:     resolvedNodes,
		Net:       analyzedNet,
		Meta:      component.Meta,
	}, nil
}

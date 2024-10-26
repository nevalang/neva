package analyzer

import (
	"errors"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

func (a Analyzer) analyzeComponent(
	component src.Component,
	scope src.Scope,
) (src.Component, *compiler.Error) {
	runtimeFuncArgs, isRuntimeFunc := component.Directives[compiler.ExternDirective]

	if isRuntimeFunc && len(runtimeFuncArgs) == 0 {
		return src.Component{}, &compiler.Error{
			Err:      errors.New("Flow that use #extern directive must provide at least one argument"),
			Location: &scope.Location,
			Meta:     &component.Meta,
		}
	}

	if len(runtimeFuncArgs) > 1 {
		for _, runtimeFuncArg := range runtimeFuncArgs {
			parts := strings.Split(runtimeFuncArg, " ")
			if len(parts) != 2 {
				return src.Component{}, &compiler.Error{
					Err:      errors.New("Component that use #extern with more than one argument must provide arguments in a form of <type, flow_ref> pairs"),
					Location: &scope.Location,
					Meta:     &component.Meta,
				}
			}
		}
	}

	resolvedInterface, err := a.analyzeInterface(
		component.Interface,
		scope,
		analyzeInterfaceParams{
			allowEmptyInports:  isRuntimeFunc,
			allowEmptyOutports: isRuntimeFunc,
		},
	)
	if err != nil {
		return src.Component{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &component.Meta,
		}.Wrap(err)
	}

	if isRuntimeFunc {
		if len(component.Nodes) != 0 || len(component.Net) != 0 {
			return src.Component{}, &compiler.Error{
				Err:      errors.New("Flow with nodes or network cannot use #extern directive"),
				Location: &scope.Location,
				Meta:     &component.Meta,
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
			Location: &scope.Location,
			Meta:     &component.Meta,
		}.Wrap(err)
	}

	if len(component.Net) == 0 {
		return src.Component{}, &compiler.Error{
			Err:      errors.New("Flow must have network"),
			Location: &scope.Location,
			Meta:     &component.Meta,
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
			Location: &scope.Location,
			Meta:     &component.Meta,
		}.Wrap(err)
	}

	return src.Component{
		Interface: resolvedInterface,
		Nodes:     resolvedNodes,
		Net:       analyzedNet,
	}, nil
}

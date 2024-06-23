package analyzer

import (
	"errors"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

//nolint:lll
var (
	ErrNodeWrongEntity     = errors.New("Node can only refer to flows or interfaces")
	ErrNodeTypeArgsMissing = errors.New("Not enough type arguments")
	ErrNodeTypeArgsTooMuch = errors.New("Too much type arguments")
	ErrNonFlowNodeWithDI   = errors.New("Only flow node can have dependency injection")
	ErrUnusedNode          = errors.New("Unused node found")
	ErrUnusedNodeInport    = errors.New("Unused node inport")
	ErrUnusedNodeOutports  = errors.New("All node's outports are unused")
	ErrSenderIsEmpty       = errors.New("Sender in network must contain port address, constant reference or message literal")
	ErrReadSelfOut         = errors.New("Flow cannot read from self outport")
	ErrWriteSelfIn         = errors.New("Flow cannot write to self inport")
	ErrInportNotFound      = errors.New("Referenced inport not found in flow's interface")
	ErrOutportNotFound     = errors.New(
		"Referenced inport not found in flow's interface",
	)
	ErrNodeNotFound               = errors.New("Referenced node not found")
	ErrNormCompWithExtern         = errors.New("Flow with nodes or network cannot use #extern directive")
	ErrNormFlowWithoutNet         = errors.New("Flow must have network except it uses #extern directive")
	ErrNormNodeBind               = errors.New("Node can't use #bind if it isn't instantiated with the flow that use #extern")
	ErrInterfaceNodeBindDirective = errors.New("Interface node cannot use #bind directive")
	ErrExternNoArgs               = errors.New("Flow that use #extern directive must provide at least one argument")
	ErrBindDirectiveArgs          = errors.New("Node with #bind directive must provide exactly one argument")
	ErrExternOverloadingArg       = errors.New("Flow that use #extern with more than one argument must provide arguments in a form of <type, flow_ref> pairs")
	ErrExternOverloadingNodeArgs  = errors.New("Node instantiated with flow with #extern with > 1 argument, must have exactly one type-argument for overloading")
)

// Maybe start here
func (a Analyzer) analyzeComponent( //nolint:funlen
	component src.Component,
	scope src.Scope,
) (src.Component, *compiler.Error) {
	runtimeFuncArgs, isRuntimeFunc := component.Directives[compiler.ExternDirective]

	if isRuntimeFunc && len(runtimeFuncArgs) == 0 {
		return src.Component{}, &compiler.Error{
			Err:      ErrExternNoArgs,
			Location: &scope.Location,
			Meta:     &component.Meta,
		}
	}

	if len(runtimeFuncArgs) > 1 {
		for _, runtimeFuncArg := range runtimeFuncArgs {
			parts := strings.Split(runtimeFuncArg, " ")
			if len(parts) != 2 {
				return src.Component{}, &compiler.Error{
					Err:      ErrExternOverloadingArg,
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
				Err:      ErrNormCompWithExtern,
				Location: &scope.Location,
				Meta:     &component.Meta,
			}
		}
		return component, nil
	}

	resolvedNodes, nodesIfaces, hasGuard, err := a.analyzeFlowNodes(
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
			Err:      ErrNormFlowWithoutNet,
			Location: &scope.Location,
			Meta:     &component.Meta,
		}
	}

	analyzedNet, err := a.analyzeFlowNetwork(
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

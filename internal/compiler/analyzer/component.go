package analyzer

import (
	"errors"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

//nolint:lll
var (
	ErrNodeWrongEntity        = errors.New("Node can only refer to components or interfaces")
	ErrNodeTypeArgsMissing    = errors.New("Not enough type arguments")
	ErrNodeTypeArgsTooMuch    = errors.New("Too much type arguments")
	ErrNonComponentNodeWithDI = errors.New("Only component node can have dependency injection")
	ErrUnusedNode             = errors.New("Unused node found")
	ErrUnusedNodeInport       = errors.New("Unused node inport")
	ErrUnusedNodeOutports     = errors.New("All node's outports are unused")
	ErrSenderIsEmpty          = errors.New("Sender in network must contain port address, constant reference or message literal")
	ErrReadSelfOut            = errors.New("Component cannot read from self outport")
	ErrWriteSelfIn            = errors.New("Component cannot write to self inport")
	ErrInportNotFound         = errors.New("Referenced inport not found in component's interface")
	ErrOutportNotFound        = errors.New(
		"Referenced inport not found in component's interface",
	)
	ErrNodeNotFound               = errors.New("Referenced node not found")
	ErrNormCompWithExtern         = errors.New("Component with nodes or network cannot use #extern directive")
	ErrNormComponentWithoutNet    = errors.New("Component must have network except it uses #extern directive")
	ErrNormNodeBind               = errors.New("Node can't use #bind if it isn't instantiated with the component that use #extern")
	ErrInterfaceNodeBindDirective = errors.New("Interface node cannot use #bind directive")
	ErrExternNoArgs               = errors.New("Component that use #extern directive must provide at least one argument")
	ErrBindDirectiveArgs          = errors.New("Node with #bind directive must provide exactly one argument")
	ErrExternOverloadingArg       = errors.New("Component that use #extern with more than one argument must provide arguments in a form of <type, component_ref> pairs")
	ErrExternOverloadingNodeArgs  = errors.New("Node instantiated with component with #extern with > 1 argument, must have exactly one type-argument for overloading")
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

	resolvedNodes, nodesIfaces, err := a.analyzeComponentNodes(
		component.Interface.TypeParams,
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
			Err:      ErrNormComponentWithoutNet,
			Location: &scope.Location,
			Meta:     &component.Meta,
		}
	}

	analyzedNet, err := a.analyzeComponentNetwork(
		component.Net,
		resolvedInterface,
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

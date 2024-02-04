package analyzer

import (
	"errors"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

//nolint:lll
var (
	ErrNodeWrongEntity                = errors.New("Node can only refer to components or interfaces")
	ErrNodeTypeArgsMissing            = errors.New("Not enough type arguments")
	ErrNodeTypeArgsTooMuch            = errors.New("Too much type arguments")
	ErrNonComponentNodeWithDI         = errors.New("Only component node can have dependency injection")
	ErrUnusedNode                     = errors.New("Unused node found")
	ErrUnusedNodeInport               = errors.New("Unused node inport found")
	ErrUnusedNodeOutports             = errors.New("All node's outports are unused")
	ErrSenderIsEmpty                  = errors.New("Sender in network must contain either port address or constant reference")
	ErrReadSelfOut                    = errors.New("Component cannot read from self outport")
	ErrWriteSelfIn                    = errors.New("Component cannot write to self inport")
	ErrInportNotFound                 = errors.New("Referenced inport not found in component's interface")
	ErrNodeNotFound                   = errors.New("Referenced node not found")
	ErrPortNotFound                   = errors.New("Port not found")
	ErrNormCompWithRuntimeFunc        = errors.New("Component with nodes or network cannot use #runtime_func directive")
	ErrNormComponentWithoutNet        = errors.New("Component must have network except it uses #runtime_func directive")
	ErrNormNodeBindDirective          = errors.New("Node can't use #bind if it isn't instantiated with the component that use #runtime_func")
	ErrInterfaceNodeBindDirective     = errors.New("Interface node cannot use #bind directive")
	ErrRuntimeFuncZeroArgs            = errors.New("Component that use #runtime_func directive must provide at least one argument")
	ErrBindDirectiveArgs              = errors.New("Node with #bind directive must provide exactly one argument")
	ErrRuntimeFuncOverloadingArg      = errors.New("Component that use #runtime_func with more than one argument must provide arguments in a form of <type, component_ref> pairs")
	ErrRuntimeFuncOverloadingNodeArgs = errors.New("Node instantiated with component with #runtime_func with > 1 argument, must have exactly one type-argument for overloading")
)

func (a Analyzer) analyzeComponent( //nolint:funlen
	component src.Component,
	scope src.Scope,
) (src.Component, *compiler.Error) {
	runtimeFuncArgs, isRuntimeFunc := component.Directives[compiler.RuntimeFuncDirective]

	if isRuntimeFunc && len(runtimeFuncArgs) == 0 {
		return src.Component{}, &compiler.Error{
			Err:      ErrRuntimeFuncZeroArgs,
			Location: &scope.Location,
			Meta:     &component.Meta,
		}
	}

	if len(runtimeFuncArgs) > 1 {
		for _, runtimeFuncArg := range runtimeFuncArgs {
			parts := strings.Split(runtimeFuncArg, " ")
			if len(parts) != 2 {
				return src.Component{}, &compiler.Error{
					Err:      ErrRuntimeFuncOverloadingArg,
					Location: &scope.Location,
					Meta:     &component.Meta,
				}
			}
		}
	}

	analyzedInterface, err := a.analyzeInterface(component.Interface, scope, analyzeInterfaceParams{
		allowEmptyInports:  isRuntimeFunc,
		allowEmptyOutports: isRuntimeFunc,
	})
	if err != nil {
		return src.Component{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &component.Meta,
		}.Merge(err)
	}

	if isRuntimeFunc {
		if len(component.Nodes) != 0 || len(component.Net) != 0 {
			return src.Component{}, &compiler.Error{
				Err:      ErrNormCompWithRuntimeFunc,
				Location: &scope.Location,
				Meta:     &component.Meta,
			}
		}
		return component, nil
	}

	analyzedNodes, nodesIfaces, err := a.analyzeComponentNodes(component.Nodes, scope)
	if err != nil {
		return src.Component{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &component.Meta,
		}.Merge(err)
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
		analyzedInterface,
		analyzedNodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return src.Component{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &component.Meta,
		}.Merge(err)
	}

	return src.Component{
		Interface: analyzedInterface,
		Nodes:     analyzedNodes,
		Net:       analyzedNet,
	}, nil
}

package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	ts "github.com/nevalang/neva/pkg/types"
)

var (
	ErrComponentTypeParams = errors.New("component type parameters")
	ErrComponentIO         = errors.New("component io")
	ErrComponentNodes      = errors.New("nodes")
	ErrComponentNet        = errors.New("net")
	ErrNodeInstance        = errors.New("node instance")
	ErrStaticInports       = errors.New("static inports")
	ErrUnusedOutports      = errors.New("unused outports")
	ErrInterfaceDIArgs     = errors.New("node instance that refers to interface cannot have DI args")
)

func (a Analyzer) analyzeCmp(
	component compiler.Component,
	scope Scope,
) (
	compiler.Component,
	map[compiler.EntityRef]struct{},
	error,
) {
	resolvedTypeParams, usedByTypeParams, err := a.analyzeTypeParameters(component.TypeParams, scope)
	if err != nil {
		return compiler.Component{}, nil, errors.Join(ErrComponentTypeParams, err)
	}

	resolvedIO, usedByIO, err := a.analyzeIO(component.IO, scope, resolvedTypeParams)
	if err != nil {
		return compiler.Component{}, nil, errors.Join(ErrComponentIO, err)
	}

	resolvedNodes, usedByNodes, err := a.analyzeNodes(component.Nodes, scope)
	if err != nil {
		return compiler.Component{}, nil, errors.Join(ErrComponentNodes, err)
	}

	unusedOutports, err := a.analyzeNet(component.Net, scope)
	if err != nil {
		return compiler.Component{}, nil, errors.Join(ErrComponentNet, err)
	}

	if len(unusedOutports) > 0 {
		return compiler.Component{}, nil, fmt.Errorf("%w: %v", ErrUnusedOutports, unusedOutports)
	}

	return compiler.Component{
		TypeParams: resolvedTypeParams,
		IO:         resolvedIO,
		Nodes:      resolvedNodes,
		Net:        []compiler.Connection{},
	}, a.mergeUsed(usedByTypeParams, usedByIO, usedByNodes), nil
}

func (a Analyzer) analyzeNodes(
	nodes map[string]compiler.Node,
	scope Scope,
) (
	map[string]compiler.Node,
	map[compiler.EntityRef]struct{},
	error,
) {
	resolvedNodes := make(map[string]compiler.Node, len(nodes))
	used := map[compiler.EntityRef]struct{}{}

	for name, node := range nodes {
		resolvedInstance, usedByInstance, err := a.analyzeNodeInstance(node, scope)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", errors.Join(ErrNodeInstance, err), name)
		}

		usedByStaticInports, err := a.analyzeStaticInports(node, scope)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", errors.Join(ErrStaticInports, err), name)
		}

		resolvedNodes[name] = compiler.Node{
			Instance:      resolvedInstance,
			StaticInports: node.StaticInports,
		}

		used = a.mergeUsed(used, usedByInstance, usedByStaticInports)
	}

	return resolvedNodes, used, nil
}

// analyzeNodeInstance finds interface or component that node is reffering to
// and checks whether it's possible to instantiate it with the given arguments
func (a Analyzer) analyzeNodeInstance(
	instance compiler.Instance,
	scope Scope,
) (
	compiler.Instance,
	map[compiler.EntityRef]struct{},
	error,
) {
	var params []ts.Param

	interf, err := scope.getInterface(instance.Ref)
	if err == nil {
		if len(instance.ComponentDI) != 0 {
			return compiler.Instance{}, nil, ErrInterfaceDIArgs
		}
		params = interf.Params
	} else {
		component, err := scope.getComponent(instance.Ref)
		if err != nil {
			return compiler.Instance{}, nil, err
		}
		// TODO implement DI analysis
		// get component's DI params (list of interfaces and names)
		// make sure DIargs are compatible with DIparams
		// make sure every DIarg refers to component and not the interface (sure?)
		params = component.TypeParams
	}

	// compatCheckParams(params, instance.TypeArgs)

	fmt.Println(params)

	return compiler.Instance{}, nil, nil
}

func (a Analyzer) analyzeStaticInports(node compiler.Node, scope Scope) (map[compiler.EntityRef]struct{}, error) {
	return nil, nil
}

// analyzeNet returns set of unused outports. It makes sure that:
// All nodes are used; Every node's inport is used; All connections refers to existing ports and are; Type-safe.
// All IO nodes are used;
func (a Analyzer) analyzeNet(net []compiler.Connection, scope Scope) (map[compiler.ConnPortAddr]struct{}, error) {
	if len(net) == 0 {
		panic("")
	}

	return nil, nil
}

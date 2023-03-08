package analyze

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
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

func (a Analyzer) analyzeComponent(
	component src.Component,
	scope Scope,
) (
	src.Component,
	map[src.EntityRef]struct{},
	error,
) {
	resolvedTypeParams, usedByTypeParams, err := a.analyzeTypeParameters(component.TypeParams, scope)
	if err != nil {
		return src.Component{}, nil, errors.Join(ErrComponentTypeParams, err)
	}

	resolvedIO, usedByIO, err := a.analyzeIO(component.IO, scope, resolvedTypeParams)
	if err != nil {
		return src.Component{}, nil, errors.Join(ErrComponentIO, err)
	}

	resolvedNodes, usedByNodes, err := a.analyzeNodes(component.Nodes, scope)
	if err != nil {
		return src.Component{}, nil, errors.Join(ErrComponentNodes, err)
	}

	unusedOutports, err := a.analyzeComponentNetwork(component.Net, scope)
	if err != nil {
		return src.Component{}, nil, errors.Join(ErrComponentNet, err)
	}

	if len(unusedOutports) > 0 {
		return src.Component{}, nil, fmt.Errorf("%w: %v", ErrUnusedOutports, unusedOutports)
	}

	return src.Component{
		TypeParams: resolvedTypeParams,
		IO:         resolvedIO,
		Nodes:      resolvedNodes,
		Net:        []src.Connection{},
	}, a.mergeUsed(usedByTypeParams, usedByIO, usedByNodes), nil
}

func (a Analyzer) analyzeNodes(
	nodes map[string]src.Node,
	scope Scope,
) (
	map[string]src.Node,
	map[src.EntityRef]struct{},
	error,
) {
	resolvedNodes := make(map[string]src.Node, len(nodes))
	used := map[src.EntityRef]struct{}{}

	for name, node := range nodes {
		resolvedInstance, usedByInstance, err := a.analyzeNodeInstance(node.Instance, scope)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", errors.Join(ErrNodeInstance, err), name)
		}

		usedByStaticInports, err := a.analyzeStaticInports(node, scope)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", errors.Join(ErrStaticInports, err), name)
		}

		resolvedNodes[name] = src.Node{
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
	instance src.NodeInstance,
	scope Scope,
) (
	src.NodeInstance,
	map[src.EntityRef]struct{},
	error,
) {
	var params []ts.Param

	interf, err := scope.getInterface(instance.Ref)
	if err == nil {
		if len(instance.DIArgs) != 0 {
			return src.NodeInstance{}, nil, ErrInterfaceDIArgs
		}
		params = interf.Params
	} else {
		component, err := scope.getComponent(instance.Ref)
		if err != nil {
			return src.NodeInstance{}, nil, err
		}
		// TODO implement DI analysis
		// get component's DI params (list of interfaces and names)
		// make sure DIargs are compatible with DIparams
		// make sure every DIarg refers to component and not the interface (sure?)
		params = component.TypeParams
	}

	// compatCheckParams(params, instance.TypeArgs)

	fmt.Println(params)

	return src.NodeInstance{}, nil, nil
}

func (a Analyzer) analyzeStaticInports(node src.Node, scope Scope) (map[src.EntityRef]struct{}, error) {
	return nil, nil
}

// analyzeComponentNetwork returns set of unused outports. It makes sure that:
// All nodes are used; Every node's inport is used; All connections refers to existing ports and are; Type-safe.
// All IO nodes are used;
func (a Analyzer) analyzeComponentNetwork(net []src.Connection, scope Scope) (map[src.ConnectionPortRef]struct{}, error) {
	return nil, nil
}

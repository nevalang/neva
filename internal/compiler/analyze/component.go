package analyze

import (
	"errors"

	"github.com/emil14/neva/internal/compiler/src"
)

var (
	ErrComponentTypeParams = errors.New("component type parameters")
	ErrComponentIO         = errors.New("component io")
	ErrComponentNodes      = errors.New("nodes")
	ErrComponentNet        = errors.New("net")
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

	resolvedNodes, usedByNodes, err := a.analyzeComponentNodes(component.Nodes, scope)
	if err != nil {
		return src.Component{}, nil, errors.Join(ErrComponentNodes, err)
	}

	unusedOutports, err := a.analyzeComponentNetwork(component.Net, scope)
	if err != nil {
		return src.Component{}, nil, errors.Join(ErrComponentNet, err)
	}

	if len(unusedOutports) > 0 {
		panic(unusedOutports)
	}

	return src.Component{
		TypeParams: resolvedTypeParams,
		IO:         resolvedIO,
		Nodes:      resolvedNodes,
		Net:        []src.Connection{},
	}, a.mergeUsed(usedByTypeParams, usedByIO, usedByNodes), nil
}



func (a Analyzer) analyzeComponentNodes(
	nodes map[string]src.Node,
	scope Scope,
) (map[string]src.Node,
	map[src.EntityRef]struct{},
	error,
) {
	return nil, nil, nil
}

// analyzeComponentNetwork returns set of unused outports. It makes sure that:
// All nodes are used; Every node's inport is used; All connections refers to existing ports and are; Type-safe.
// All IO nodes are used;
func (a Analyzer) analyzeComponentNetwork(net []src.Connection, scope Scope) (map[src.ConnectionPortRef]struct{}, error) {
	return nil, nil
}

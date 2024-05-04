package desugarer

import (
	"errors"
	"maps"
	"slices"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

var ErrConstSenderEntityKind = errors.New(
	"Entity that is used as a const reference in component's network must be of kind constant",
)

type handleComponentResult struct {
	desugaredComponent src.Component         // desugared component to replace
	virtualEntities    map[string]src.Entity //nolint:lll // sometimes after desugaring component we need to insert some entities to the file
}

func (d Desugarer) handleComponent( //nolint:funlen
	component src.Component,
	scope src.Scope,
) (handleComponentResult, *compiler.Error) {
	if len(component.Net) == 0 && len(component.Nodes) == 0 {
		return handleComponentResult{desugaredComponent: component}, nil
	}

	virtualEntities := map[string]src.Entity{}

	desugaredNodes, virtConnsForNodes, err := d.handleNodes(
		component,
		scope,
		virtualEntities,
	)
	if err != nil {
		return handleComponentResult{}, err
	}

	netToDesugar := append(virtConnsForNodes, component.Net...)
	handleNetResult, err := d.handleNetwork(
		netToDesugar,
		desugaredNodes,
		scope,
	)
	if err != nil {
		return handleComponentResult{}, err
	}

	desugaredNetwork := slices.Clone(handleNetResult.desugaredConnections)

	// add virtual constants created by network handler to virtual entities
	for name, constant := range handleNetResult.virtualConstants {
		virtualEntities[name] = src.Entity{
			Kind:  src.ConstEntity,
			Const: constant,
		}
	}

	// merge real nodes with virtual ones created by network handler
	maps.Copy(desugaredNodes, handleNetResult.virtualNodes)

	// create virtual destructor nodes and connections to handle unused outports
	unusedOutports := d.findUnusedOutports(
		component,
		scope,
		handleNetResult.usedNodePorts,
	)
	if unusedOutports.len() != 0 {
		unusedOutportsResult := d.handleUnusedOutports(unusedOutports)
		desugaredNetwork = append(desugaredNetwork, unusedOutportsResult.virtualConnections...)
		desugaredNodes[unusedOutportsResult.voidNodeName] = unusedOutportsResult.voidNode
	}

	return handleComponentResult{
		desugaredComponent: src.Component{
			Directives: component.Directives,
			Interface:  component.Interface,
			Nodes:      desugaredNodes,
			Net:        desugaredNetwork,
			Meta:       component.Meta,
		},
		virtualEntities: virtualEntities,
	}, nil
}

func (d Desugarer) handleNodes(
	component src.Component,
	scope src.Scope,
	virtualEntities map[string]src.Entity,
) (map[string]src.Node, []src.Connection, *compiler.Error) {
	desugaredNodes := make(map[string]src.Node, len(component.Nodes))
	virtualConns := []src.Connection{}

	for nodeName, node := range component.Nodes {
		extraConns, err := d.handleNode(
			scope,
			node,
			desugaredNodes,
			nodeName,
			virtualEntities,
		)
		if err != nil {
			return nil, nil, err
		}

		virtualConns = append(virtualConns, extraConns...)
	}

	return desugaredNodes, virtualConns, nil
}

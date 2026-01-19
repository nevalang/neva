package desugarer

import (
	"maps"
	"slices"

	src "github.com/nevalang/neva/internal/compiler/ast"
)

type handleComponentResult struct {
	desugaredFlow   src.Component
	virtualEntities map[string]src.Entity
}

func (d *Desugarer) desugarComponent(
	component src.Component,
	scope src.Scope,
) (handleComponentResult, error) {
	if len(component.Net) == 0 && len(component.Nodes) == 0 {
		return handleComponentResult{desugaredFlow: component}, nil
	}

	virtualEntities := map[string]src.Entity{}

	desugaredNodes, virtualConnections, err := d.desugarNodes(
		component,
		scope,
		virtualEntities,
	)
	if err != nil {
		return handleComponentResult{}, err
	}

	connectionsToDesugar := append(virtualConnections, component.Net...)

	desugarNetResult, err := d.desugarNetwork(
		component.Interface,
		connectionsToDesugar,
		desugaredNodes,
		scope,
	)
	if err != nil {
		return handleComponentResult{}, err
	}

	desugaredNetwork := slices.Clone(desugarNetResult.desugaredConnections)

	// add virtual constants created by network handler to virtual entities
	for name, constant := range desugarNetResult.constsToInsert {
		virtualEntities[name] = src.Entity{
			Kind:  src.ConstEntity,
			Const: constant,
		}
	}

	// merge real nodes with virtual ones created by network handler
	maps.Copy(desugaredNodes, desugarNetResult.nodesToInsert)

	desugaredNetwork, err = d.insertErrFanInIfNeeded(
		desugaredNetwork,
		desugaredNodes,
		component.Meta,
	)
	if err != nil {
		return handleComponentResult{}, err
	}

	// create and connect Del nodes to handle unused outports
	unusedOutports := d.findUnusedOutports(
		component,
		scope,
		desugarNetResult.nodesPortsUsed,
	)
	if unusedOutports.len() != 0 {
		unusedOutportsResult := d.handleUnusedOutports(unusedOutports, component.Meta)
		desugaredNetwork = append(desugaredNetwork, unusedOutportsResult.virtualConnections...)
		desugaredNodes[unusedOutportsResult.voidNodeName] = unusedOutportsResult.delNode
	}

	return handleComponentResult{
		desugaredFlow: src.Component{
			Directives: component.Directives,
			Interface:  component.Interface,
			Nodes:      desugaredNodes,
			Net:        desugaredNetwork,
			Meta:       component.Meta,
		},
		virtualEntities: virtualEntities,
	}, nil
}

func (d *Desugarer) desugarNodes(
	component src.Component,
	scope src.Scope,
	virtualEntities map[string]src.Entity,
) (map[string]src.Node, []src.Connection, error) {
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

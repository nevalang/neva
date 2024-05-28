package desugarer

import (
	"errors"
	"maps"
	"slices"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

var ErrConstSenderEntityKind = errors.New(
	"Entity that is used as a const reference in flow's network must be of kind constant",
)

type handleFlowResult struct {
	desugaredFlow   src.Flow
	virtualEntities map[string]src.Entity
}

func (d Desugarer) handleFlow(
	flow src.Flow,
	scope src.Scope,
) (handleFlowResult, *compiler.Error) {
	if len(flow.Net) == 0 && len(flow.Nodes) == 0 {
		return handleFlowResult{desugaredFlow: flow}, nil
	}

	virtualEntities := map[string]src.Entity{}

	desugaredNodes, virtConnsForNodes, err := d.handleNodes(
		flow,
		scope,
		virtualEntities,
	)
	if err != nil {
		return handleFlowResult{}, err
	}

	netToDesugar := append(virtConnsForNodes, flow.Net...)
	handleNetResult, err := d.handleNetwork(
		netToDesugar,
		desugaredNodes,
		scope,
	)
	if err != nil {
		return handleFlowResult{}, err
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
		flow,
		scope,
		handleNetResult.usedNodePorts,
	)
	if unusedOutports.len() != 0 {
		unusedOutportsResult := d.handleUnusedOutports(unusedOutports)
		desugaredNetwork = append(desugaredNetwork, unusedOutportsResult.virtualConnections...)
		desugaredNodes[unusedOutportsResult.voidNodeName] = unusedOutportsResult.voidNode
	}

	return handleFlowResult{
		desugaredFlow: src.Flow{
			Directives: flow.Directives,
			Interface:  flow.Interface,
			Nodes:      desugaredNodes,
			Net:        desugaredNetwork,
			Meta:       flow.Meta,
		},
		virtualEntities: virtualEntities,
	}, nil
}

func (d Desugarer) handleNodes(
	flow src.Flow,
	scope src.Scope,
	virtualEntities map[string]src.Entity,
) (map[string]src.Node, []src.Connection, *compiler.Error) {
	desugaredNodes := make(map[string]src.Node, len(flow.Nodes))
	virtualConns := []src.Connection{}

	for nodeName, node := range flow.Nodes {
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

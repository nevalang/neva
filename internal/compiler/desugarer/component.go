package desugarer

import (
	"errors"
	"fmt"
	"maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
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
	// if it's native component, nothing to desugar
	if len(component.Net) == 0 && len(component.Nodes) == 0 {
		return handleComponentResult{desugaredComponent: component}, nil
	}

	// we are going to create some entities from scratch
	virtualEntities := map[string]src.Entity{}

	// handle nodes
	desugaredNodes := make(map[string]src.Node, len(component.Nodes))
	for nodeName, node := range component.Nodes {
		err := d.handleNode(scope, node, desugaredNodes, nodeName, virtualEntities)
		if err != nil {
			return handleComponentResult{}, err
		}
	}

	// handle network
	handleNetResult, err := d.handleNetwork(component.Net, desugaredNodes, scope)
	if err != nil {
		return handleComponentResult{}, err
	}

	// add virtual constants created by network handler to virtual entities
	for name, constant := range handleNetResult.virtualConstants {
		virtualEntities[name] = src.Entity{
			Kind:  src.ConstEntity,
			Const: constant,
		}
	}

	// create alias
	desugaredNetwork := handleNetResult.desugaredConnections

	// merge real nodes with virtual ones created by network handler
	maps.Copy(desugaredNodes, handleNetResult.virtualNodes)

	// create virtual destructor nodes and connections to handle unused outports
	unusedOutports := d.findUnusedOutports(component, scope, handleNetResult.usedNodePorts)
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

func (Desugarer) handleNode(
	scope src.Scope,
	node src.Node,
	desugaredNodes map[string]src.Node,
	nodeName string,
	virtualEntities map[string]src.Entity,
) *compiler.Error {
	entity, _, err := scope.Entity(node.EntityRef)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &scope.Location,
		}
	}

	if entity.Kind != src.ComponentEntity {
		desugaredNodes[nodeName] = node
		return nil
	}

	_, ok := entity.Component.Directives[compiler.AutoportsDirective]
	if !ok {
		desugaredNodes[nodeName] = node
		return nil
	}

	structFields := node.TypeArgs[0].Lit.Struct

	inports := make(map[string]src.Port, len(structFields))
	for fieldName, fieldTypeExpr := range structFields {
		inports[fieldName] = src.Port{
			TypeExpr: fieldTypeExpr,
		}
	}

	outports := map[string]src.Port{
		"msg": {
			TypeExpr: node.TypeArgs[0],
		},
	}

	localBuilderComponent := src.Component{
		Interface: src.Interface{
			IO: src.IO{In: inports, Out: outports},
		},
	}

	localBuilderName := fmt.Sprintf("struct_builder_%v", nodeName)

	virtualEntities[localBuilderName] = src.Entity{
		Kind:      src.ComponentEntity,
		Component: localBuilderComponent,
	}

	desugaredNodes[nodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "StructBuilder",
		},
		Directives: node.Directives,
		TypeArgs:   node.TypeArgs,
		Deps:       node.Deps,
		Meta:       node.Meta,
	}

	return nil
}

package desugarer

import (
	"errors"
	"fmt"
	"maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

var ErrConstSenderEntityKind = errors.New(
	"Entity that is used as a const reference in component's network must be of kind constant",
)

type desugarComponentResult struct {
	component        src.Component         // desugared component to replace
	entitiesToInsert map[string]src.Entity //nolint:lll // sometimes after desugaring component we need to insert some entities to the file
}

func (d Desugarer) desugarComponent( //nolint:funlen
	component src.Component,
	scope src.Scope,
) (desugarComponentResult, *compiler.Error) {
	if len(component.Net) == 0 && len(component.Nodes) == 0 {
		return desugarComponentResult{
			component: component,
		}, nil
	}

	entitiesToInsert := map[string]src.Entity{}

	desugaredNodes := make(map[string]src.Node, len(component.Nodes))
	for nodeName, node := range component.Nodes {
		entity, _, err := scope.Entity(node.EntityRef)
		if err != nil {
			return desugarComponentResult{}, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
			}
		}

		if entity.Kind != src.ComponentEntity {
			desugaredNodes[nodeName] = node
			continue
		}

		_, ok := entity.Component.Directives[compiler.AutoportsDirective]
		if !ok {
			desugaredNodes[nodeName] = node
			continue
		}

		structFields := node.TypeArgs[0].Lit.Struct // must be resolved struct after analysis stage

		inports := make(map[string]src.Port, len(structFields))
		for fieldName, fieldTypeExpr := range structFields {
			inports[fieldName] = src.Port{
				TypeExpr: fieldTypeExpr,
			}
		}

		outports := map[string]src.Port{
			"v": {
				TypeExpr: node.TypeArgs[0],
			},
		}

		// create local variation of the struct builder component with inports corresponding to struct fields
		localBuilderComponent := src.Component{
			Interface: src.Interface{
				IO: src.IO{In: inports, Out: outports}, // these ports gonna be used by irgen and then by runtime func
			},
		}

		localBuilderName := fmt.Sprintf("__struct_builder_%v__", nodeName)

		entitiesToInsert[localBuilderName] = src.Entity{
			Kind:      src.ComponentEntity,
			Component: localBuilderComponent,
		}

		// finally replace component ref for this current node with the ref to newly created local builder variation
		desugaredNodes[nodeName] = src.Node{
			EntityRef: src.EntityRef{
				Pkg:  "builtin",
				Name: "StructBuilder",
			},
			Directives: node.Directives,
			TypeArgs:   node.TypeArgs,
			Deps:       node.Deps,
			Meta:       node.Meta,
		}
	}

	handleConnsResult, err := d.handleConns(component.Net, desugaredNodes, scope)
	if err != nil {
		return desugarComponentResult{}, err
	}

	desugaredNet := handleConnsResult.desugaredConns
	maps.Copy(desugaredNodes, handleConnsResult.extraNodes)

	unusedOutports := d.findUnusedOutports(component, scope, handleConnsResult.usedNodePorts)
	if unusedOutports.len() != 0 {
		voidResult := d.getVoidNodeAndConns(unusedOutports)
		desugaredNet = append(desugaredNet, voidResult.voidConns...)
		desugaredNodes[voidResult.voidNodeName] = voidResult.voidNode
	}

	for name, constant := range handleConnsResult.extraConsts {
		entitiesToInsert[name] = src.Entity{
			Kind:  src.ConstEntity,
			Const: constant,
		}
	}

	return desugarComponentResult{
		component: src.Component{
			Directives: component.Directives,
			Interface:  component.Interface,
			Nodes:      desugaredNodes,
			Net:        desugaredNet,
			Meta:       component.Meta,
		},
		entitiesToInsert: entitiesToInsert,
	}, nil
}

type handleConnsResult struct {
	desugaredConns []src.Connection     // desugared network
	extraConsts    map[string]src.Const // constants that needs to be inserted in to make desugared network work
	extraNodes     map[string]src.Node  // nodes that needs to be inserted in to make desugared network work
	usedNodePorts  nodePortsMap         // ports that were used in processed network
}

func (d Desugarer) handleConns( //nolint:funlen
	conns []src.Connection,
	nodes map[string]src.Node,
	scope src.Scope,
) (handleConnsResult, *compiler.Error) {
	nodesToInsert := map[string]src.Node{}
	desugaredConns := make([]src.Connection, 0, len(conns))
	usedNodePorts := newNodePortsMap()
	constsToInsert := map[string]src.Const{}

	for _, conn := range conns {
		if conn.SenderSide.PortAddr != nil { // const sender are not interested, we 100% they're used (we handle that here)
			usedNodePorts.set(
				conn.SenderSide.PortAddr.Node,
				conn.SenderSide.PortAddr.Port,
			)
		}

		if conn.SenderSide.Const == nil &&
			len(conn.SenderSide.Selectors) == 0 &&
			len(conn.ReceiverSide.ThenConnections) == 0 {
			desugaredConns = append(desugaredConns, conn)
			continue
		}

		if len(conn.SenderSide.Selectors) != 0 {
			result, err := d.desugarStructSelectors(
				conn,
				nodes,
				scope,
			)
			if err != nil {
				return handleConnsResult{}, compiler.Error{
					Err:      errors.New("Cannot desugar struct selectors"),
					Location: &scope.Location,
					Meta:     &conn.Meta,
				}.Merge(err)
			}
			nodesToInsert[result.nodeToInsertName] = result.nodeToInsert
			constsToInsert[result.constToInsertName] = result.constToInsert
			conn = result.connToReplace
			desugaredConns = append(desugaredConns, result.connToInsert)
		}

		if conn.SenderSide.Const != nil { //nolint:nestif
			if conn.SenderSide.Const.Ref != nil {
				result, err := d.handleConstRefSender(conn, scope)
				if err != nil {
					return handleConnsResult{}, err
				}
				nodesToInsert[result.emitterNodeName] = result.emitterNode
				conn = result.desugaredConn
			} else if conn.SenderSide.Const.Value != nil {
				result, err := d.handleLiteralSender(conn)
				if err != nil {
					return handleConnsResult{}, err
				}
				constsToInsert[result.constName] = *conn.SenderSide.Const
				nodesToInsert[result.emitterNodeName] = result.emitterNode
				conn = result.desugaredConn
			}
		}

		desugaredConns = append(desugaredConns, conn)

		if len(conn.ReceiverSide.ThenConnections) == 0 {
			continue
		}

		result, err := d.handleThenConns(conn, nodes, scope)
		if err != nil {
			return handleConnsResult{}, err
		}

		// handleThenConns recursively calls this function so it returns the same structure
		maps.Copy(usedNodePorts.m, result.usedNodesPorts.m)
		maps.Copy(constsToInsert, result.extraConsts)
		maps.Copy(nodesToInsert, result.extraNodes)

		desugaredConns = append(desugaredConns, result.extraConns...)
	}

	return handleConnsResult{
		desugaredConns: desugaredConns,
		usedNodePorts:  usedNodePorts,
		extraConsts:    constsToInsert,
		extraNodes:     nodesToInsert,
	}, nil
}

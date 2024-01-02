package desugarer

import (
	"errors"
	"fmt"
	"maps"

	"github.com/nevalang/neva/internal/compiler"

	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var ErrConstSenderEntityKind = errors.New("Entity that is used as a const reference in component's network must be of kind constant") //nolint:lll

// desugarComponent replaces const ref in net with regular port addr and injects const node with directive.
func (d Desugarer) desugarComponent( //nolint:funlen
	component src.Component,
	scope src.Scope,
) (src.Component, *compiler.Error) {
	if len(component.Net) == 0 && len(component.Nodes) == 0 {
		return component, nil
	}

	desugaredNodes := maps.Clone(component.Nodes)
	if desugaredNodes == nil {
		desugaredNodes = map[string]src.Node{}
	}

	usedNodePorts := newNodePortsMap()
	desugaredNet := make([]src.Connection, 0, len(component.Net))
	for _, conn := range component.Net {
		if conn.SenderSide.ConstRef == nil {
			desugaredNet = append(desugaredNet, conn)
			if conn.SenderSide.PortAddr != nil {
				usedNodePorts.set(
					conn.SenderSide.PortAddr.Node,
					conn.SenderSide.PortAddr.Port,
				)
			}
			continue
		}

		constTypeExpr, err := d.getConstType(*conn.SenderSide.ConstRef, scope)
		if err != nil {
			return src.Component{}, compiler.Error{
				Err:      fmt.Errorf("Unable to get constant type by reference '%v'", *conn.SenderSide.ConstRef),
				Location: &scope.Location,
				Meta:     &conn.SenderSide.ConstRef.Meta,
			}.Merge(err)
		}

		constRefStr := conn.SenderSide.ConstRef.String()
		constNodeName := fmt.Sprintf("__%v__", constRefStr)

		desugaredNodes[constNodeName] = src.Node{
			Directives: map[src.Directive][]string{
				compiler.RuntimeFuncMsgDirective: {constRefStr},
			},
			EntityRef: src.EntityRef{
				Pkg:  "builtin",
				Name: "Const",
			},
			TypeArgs: []ts.Expr{constTypeExpr},
		}

		constNodeOutportAddr := src.PortAddr{
			Node: constNodeName,
			Port: "v",
		}

		desugaredNet = append(desugaredNet, src.Connection{
			SenderSide: src.SenderConnectionSide{
				PortAddr:  &constNodeOutportAddr,
				Selectors: conn.SenderSide.Selectors,
				Meta:      conn.SenderSide.Meta,
			},
			ReceiverSides: conn.ReceiverSides,
			Meta:          conn.Meta,
		})
	}

	// try to find unused nodes outports to sugar them with voids
	unusedOutports := newNodePortsMap()
	for nodeName, node := range component.Nodes {
		entity, _, err := scope.Entity(node.EntityRef)
		if err != nil {
			continue
		}
		if entity.Kind != src.InterfaceEntity && entity.Kind != src.ComponentEntity {
			continue
		}

		var io src.IO
		if entity.Kind == src.InterfaceEntity {
			io = entity.Interface.IO
		} else {
			io = entity.Component.Interface.IO
		}

		for outportName := range io.Out {
			ok := usedNodePorts.get(nodeName, outportName)
			if !ok {
				unusedOutports.set(nodeName, outportName)
			}
		}
	}

	if unusedOutports.len() == 0 { // no need to mess with voids
		return src.Component{
			Directives: component.Directives,
			Interface:  component.Interface,
			Nodes:      desugaredNodes,
			Net:        desugaredNet,
			Meta:       component.Meta,
		}, nil
	}

	// add void node and connections to it
	voidNodeName := "__void__"
	desugaredNodes[voidNodeName] = src.Node{
		EntityRef: src.EntityRef{
			Pkg:  "builtin",
			Name: "Void",
		},
	}
	receiverSides := []src.ReceiverConnectionSide{
		{PortAddr: src.PortAddr{Node: voidNodeName, Port: "v"}},
	}
	for nodeName, ports := range unusedOutports.m {
		for portName := range ports {
			desugaredNet = append(desugaredNet, src.Connection{
				SenderSide: src.SenderConnectionSide{
					PortAddr: &src.PortAddr{
						Node: nodeName,
						Port: portName,
					},
				},
				ReceiverSides: receiverSides,
				Meta:          src.Meta{},
			})
		}
	}

	return src.Component{
		Directives: component.Directives,
		Interface:  component.Interface,
		Nodes:      desugaredNodes,
		Net:        desugaredNet,
		Meta:       component.Meta,
	}, nil
}

// getConstType is needed to figure out type parameters for Const node
func (d Desugarer) getConstType(ref src.EntityRef, scope src.Scope) (ts.Expr, *compiler.Error) {
	entity, _, err := scope.Entity(ref)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &ref.Meta,
		}
	}

	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrConstSenderEntityKind, entity.Kind),
			Location: &scope.Location,
			Meta:     entity.Meta(),
		}
	}

	if entity.Const.Ref != nil {
		expr, err := d.getConstType(*entity.Const.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Location: &scope.Location,
				Meta:     entity.Meta(),
			}.Merge(err)
		}
		return expr, nil
	}

	return entity.Const.Value.TypeExpr, nil
}

type nodePortsMap struct {
	m map[string]map[string]struct{}
}

func (n nodePortsMap) set(node string, outport string) {
	if n.m[node] == nil {
		n.m[node] = map[string]struct{}{}
	}
	n.m[node][outport] = struct{}{}
}

func (n nodePortsMap) get(node, port string) bool {
	ports, ok := n.m[node]
	if !ok {
		return false
	}
	_, ok = ports[port]
	return ok
}

func (n nodePortsMap) len() int {
	return len(n.m)
}

func newNodePortsMap() nodePortsMap {
	return nodePortsMap{
		m: map[string]map[string]struct{}{},
	}
}

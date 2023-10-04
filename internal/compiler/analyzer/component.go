package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

func (a Analyzer) analyzeComponent(comp src.Component, scope Scope) (src.Component, error) {
	resolvedInterface, err := a.analyzeInterface(comp.Interface, scope)
	if err != nil {
		return src.Component{}, fmt.Errorf("analyze interface: %w", err)
	}

	resolvedNodes, err := a.analyzeComponentNodes(comp.Nodes, scope)
	if err != nil {
		return src.Component{}, fmt.Errorf("analyze component nodes: %w", err)
	}

	resolvedNet, err := a.analyzeComponentNet(comp.Net, resolvedInterface, resolvedNodes, scope)
	if err != nil {
		return src.Component{}, fmt.Errorf("analyze component network: %w", err)
	}

	return src.Component{
		Interface: resolvedInterface,
		Nodes:     resolvedNodes,
		Net:       resolvedNet,
	}, nil
}

func (a Analyzer) analyzeComponentNodes(nodes map[string]src.Node, scope Scope) (map[string]src.Node, error) {
	resolvedNodes := make(map[string]src.Node, len(nodes))
	for name, node := range nodes {
		resolvedNode, err := a.analyzeComponentNode(node, scope)
		if err != nil {
			return nil, fmt.Errorf("analyze node: %w", err)
		}
		resolvedNodes[name] = resolvedNode
	}
	return resolvedNodes, nil
}

var (
	ErrNodeWrongEntity           = fmt.Errorf("node entity is not a component or interface")
	ErrNodeTypeArgsCountMismatch = errors.New("node type args count mismatch")
	ErrNodeInterfaceDI           = errors.New("interface node cannot have dependency injection")
)

func (a Analyzer) analyzeComponentNode(node src.Node, scope Scope) (src.Node, error) {
	entity, _, err := scope.prog.Entity(node.EntityRef)
	if err != nil {
		return src.Node{}, fmt.Errorf("entity: %w", err)
	}

	if entity.Kind != src.ComponentEntity && entity.Kind != src.InterfaceEntity {
		return src.Node{}, fmt.Errorf("%w: %v", ErrNodeWrongEntity, entity.Kind)
	}

	var compInterface src.Interface
	if entity.Kind == src.ComponentEntity {
		compInterface = entity.Component.Interface
	} else {
		if node.ComponentDI != nil {
			return src.Node{}, ErrNodeInterfaceDI
		}
		compInterface = entity.Interface
	}

	if len(node.TypeArgs) != len(compInterface.TypeParams) {
		return src.Node{}, fmt.Errorf(
			"%w: want %v, got %v",
			ErrNodeTypeArgsCountMismatch, compInterface.TypeParams, node.TypeArgs,
		)
	}

	resolvedArgs, _, err := a.resolver.ResolveFrame(node.TypeArgs, compInterface.TypeParams, scope)
	if err != nil {
		return src.Node{}, fmt.Errorf("resolve args: %w", err)
	}

	if node.ComponentDI == nil {
		return src.Node{
			EntityRef: node.EntityRef,
			TypeArgs:  resolvedArgs,
		}, nil
	}

	resolvedComponentDI := make(map[string]src.Node, len(node.ComponentDI))
	for depName, depNode := range node.ComponentDI {
		resolvedDep, err := a.analyzeComponentNode(depNode, scope)
		if err != nil {
			return src.Node{}, fmt.Errorf("analyze dependency node: %w", err)
		}
		resolvedComponentDI[depName] = resolvedDep
	}

	return src.Node{
		EntityRef:   node.EntityRef,
		TypeArgs:    resolvedArgs,
		ComponentDI: resolvedComponentDI,
	}, nil
}

func (a Analyzer) analyzeComponentNet(
	net []src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	scope Scope,
) ([]src.Connection, error) {
	for _, conn := range net {
		senderType, err := a.getSenderType(conn.SenderSide, compInterface.IO.In, nodes, scope)
		if err != nil {
			return nil, fmt.Errorf("get sender type: %w", err)
		}
		for _, receiver := range conn.ReceiverSides {
			receiverType, err := a.getReceiverType(receiver, compInterface.IO.Out, nodes, scope)
			if err != nil {
				return nil, fmt.Errorf("get sender type: %w", err)
			}
			if err := a.resolver.IsSubtypeOf(senderType, receiverType, scope); err != nil {
				return nil, fmt.Errorf("is subtype of: %w", err)
			}
		}
	}
	return net, nil
}

func (a Analyzer) getReceiverType(
	receiverSide src.ReceiverConnectionSide,
	outports map[string]src.Port,
	nodes map[string]src.Node,
	scope Scope,
) (ts.Expr, error) {
	if receiverSide.PortAddr.Node == "in" {
		return ts.Expr{}, ErrWriteSelfIn
	}

	if receiverSide.PortAddr.Node == "out" {
		outport, ok := outports[receiverSide.PortAddr.Port]
		if !ok {
			return ts.Expr{}, ErrInportNotFound
		}
		return outport.TypeExpr, nil
	}

	nodeInportType, err := a.getNodeInportType(receiverSide.PortAddr, nodes, scope)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("get node inport type: %w", err)
	}

	return nodeInportType, nil
}

func (a Analyzer) getNodeInportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	scope Scope,
) (ts.Expr, error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, ErrNodeNotFound
	}
	component, _, _ := scope.prog.Entity(node.EntityRef) // nodes analyzed so we don't check error
	return a.getResolvedPortType(component.Interface.IO.In, component.Interface.TypeParams, portAddr, node, scope)
}

func (a Analyzer) getResolvedPortType(
	ports map[string]src.Port,
	params []ts.Param,
	portAddr src.PortAddr,
	node src.Node,
	scope Scope,
) (ts.Expr, error) {
	port, ok := ports[portAddr.Node]
	if !ok {
		return ts.Expr{}, ErrNodeOutportNotFound
	}

	_, frame, err := a.resolver.ResolveFrame(node.TypeArgs, params, scope)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("resolve args: %w", err)
	}

	resolvedOutportType, err := a.resolver.ResolveExprWithFrame(port.TypeExpr, frame, scope)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("resolve expr with frame: %w", err)
	}

	return resolvedOutportType, nil
}

var (
	ErrSenderConstRefEntityKind = errors.New("entity reference in network sender  points to not constant")
	ErrSenderEmpty              = errors.New("network sender must contain either const ref or port addr")
	ErrReadSelfOut              = errors.New("component cannot read from self outport")
	ErrWriteSelfIn              = errors.New("component cannot write to self inports")
	ErrInportNotFound           = errors.New("network references to inport that is not found in component's IO")
	ErrNodeNotFound             = errors.New("network references node that is not found in component")
	ErrNodeOutportNotFound      = errors.New("network references to not existing node's outport")
)

func (a Analyzer) getSenderType(
	senderSide src.SenderConnectionSide,
	inports map[string]src.Port,
	nodes map[string]src.Node,
	scope Scope,
) (ts.Expr, error) {
	if senderSide.ConstRef != nil {
		constTypeExpr, err := a.getConstType(*senderSide.ConstRef, scope)
		if err != nil {
			return ts.Expr{}, fmt.Errorf("get const type: %w", err)
		}
		return constTypeExpr, nil
	}

	if senderSide.PortAddr == nil {
		return ts.Expr{}, ErrSenderEmpty
	}
	if senderSide.PortAddr.Node == "out" {
		return ts.Expr{}, ErrReadSelfOut
	}

	if senderSide.PortAddr.Node == "in" {
		inport, ok := inports[senderSide.PortAddr.Port]
		if !ok {
			return ts.Expr{}, ErrInportNotFound
		}
		return inport.TypeExpr, nil
	}

	nodeOutportType, err := a.getNodeOutportType(*senderSide.PortAddr, nodes, scope)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("get node outport type: %w", err)
	}

	return nodeOutportType, nil
}

func (a Analyzer) getNodeOutportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	scope Scope,
) (ts.Expr, error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, ErrNodeNotFound
	}
	component, _, _ := scope.prog.Entity(node.EntityRef) // nodes analyzed so we don't check error
	return a.getResolvedPortType(component.Interface.IO.Out, component.Interface.TypeParams, portAddr, node, scope)
}

func (a Analyzer) getConstType(ref src.EntityRef, scope Scope) (ts.Expr, error) {
	entity, _, err := scope.prog.Entity(ref)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("prog entity: %w", err)
	}
	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, fmt.Errorf("%w: %v", ErrSenderConstRefEntityKind, entity.Kind)
	}
	if entity.Const.Ref != nil {
		return a.getConstType(*entity.Const.Ref, scope)
	}
	return entity.Const.Value.TypeExpr, nil
}

package analyzer

import (
	"errors"
	"fmt"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var (
	ErrNodeWrongEntity           = errors.New("Node can only refer to components or interfaces")
	ErrNodeTypeArgsCountMismatch = errors.New("Type arguments count between node and its referenced entity does not match")
	ErrNodeInterfaceDI           = errors.New("Interface node cannot have dependency injection")
	ErrUnusedNode                = errors.New("Unused node found")
	ErrUnusedNodeInport          = errors.New("Unused node inport found")
	ErrUnusedNodeOutport         = errors.New("Unused node outport found")
	ErrSenderConstRefEntityKind  = errors.New("Sender in network with entity reference can only refer to a constant")
	ErrSenderIsEmpty             = errors.New("Sender in network must either refer to some port address or constant")
	ErrReadSelfOut               = errors.New("Component cannot read from self outport")
	ErrWriteSelfIn               = errors.New("Component cannot write to self inport")
	ErrInportNotFound            = errors.New("Referenced inport not found in component's interface")
	ErrNodeNotFound              = errors.New("Referenced node not found")
	ErrNodePortNotFound          = errors.New("Referenced node port not found")
)

type analyzeComponentParams struct {
	iface analyzeInterfaceParams
}

func (a Analyzer) analyzeComponent(
	comp src.Component,
	scope src.Scope,
	params analyzeComponentParams,
) (src.Component, error) {
	resolvedInterface, err := a.analyzeInterface(comp.Interface, scope, params.iface)
	if err != nil {
		return src.Component{}, fmt.Errorf("analyze interface: %w", err)
	}

	resolvedNodes, nodesIfaces, err := a.analyzeComponentNodes(comp.Nodes, scope)
	if err != nil {
		return src.Component{}, fmt.Errorf("analyze component nodes: %w", err)
	}

	resolvedNet, err := a.analyzeComponentNetwork(comp.Net, resolvedInterface, resolvedNodes, nodesIfaces, scope)
	if err != nil {
		return src.Component{}, fmt.Errorf("analyze component network: %w", err)
	}

	return src.Component{
		Interface: resolvedInterface,
		Nodes:     resolvedNodes,
		Net:       resolvedNet,
	}, nil
}

func (a Analyzer) analyzeComponentNodes(
	nodes map[string]src.Node,
	scope src.Scope,
) (map[string]src.Node, map[string]src.Interface, error) {
	resolvedNodes := make(map[string]src.Node, len(nodes))
	nodesIfaces := make(map[string]src.Interface, len(nodes))
	for name, node := range nodes {
		resolvedNode, iface, err := a.analyzeComponentNode(node, scope)
		if err != nil {
			return nil, nil, fmt.Errorf("analyze node: %v: %w", name, err)
		}
		nodesIfaces[name] = iface
		resolvedNodes[name] = resolvedNode
	}
	return resolvedNodes, nodesIfaces, nil
}

func (a Analyzer) analyzeComponentNode(node src.Node, scope src.Scope) (src.Node, src.Interface, error) {
	entity, _, err := scope.Entity(node.EntityRef)
	if err != nil {
		return src.Node{}, src.Interface{}, fmt.Errorf("entity: %v: %w", node.EntityRef, err)
	}

	if entity.Kind != src.ComponentEntity && entity.Kind != src.InterfaceEntity {
		return src.Node{}, src.Interface{}, fmt.Errorf("%w: %v", ErrNodeWrongEntity, entity.Kind)
	}

	var iface src.Interface
	if entity.Kind == src.ComponentEntity {
		iface = entity.Component.Interface
	} else {
		if node.ComponentDI != nil {
			return src.Node{}, src.Interface{}, ErrNodeInterfaceDI
		}
		iface = entity.Interface
	}

	if len(node.TypeArgs) != len(iface.TypeParams.Params) {
		return src.Node{}, src.Interface{}, fmt.Errorf(
			"%w: want %v, got %v",
			ErrNodeTypeArgsCountMismatch, iface.TypeParams, node.TypeArgs,
		)
	}

	resolvedArgs, _, err := a.resolver.ResolveFrame(node.TypeArgs, iface.TypeParams.Params, scope)
	if err != nil {
		return src.Node{}, src.Interface{}, fmt.Errorf("resolve args: %w", err)
	}

	if node.ComponentDI == nil {
		return src.Node{
			EntityRef: node.EntityRef,
			TypeArgs:  resolvedArgs,
		}, iface, nil
	}

	resolvedComponentDI := make(map[string]src.Node, len(node.ComponentDI))
	for depName, depNode := range node.ComponentDI {
		resolvedDep, _, err := a.analyzeComponentNode(depNode, scope)
		if err != nil {
			return src.Node{}, src.Interface{}, fmt.Errorf("analyze dep node: %v: %w", depNode, err)
		}
		resolvedComponentDI[depName] = resolvedDep
	}

	return src.Node{
		EntityRef:   node.EntityRef,
		TypeArgs:    resolvedArgs,
		ComponentDI: resolvedComponentDI,
	}, iface, nil
}

func (a Analyzer) analyzeComponentNetwork(
	net []src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) ([]src.Connection, error) {
	nodesUsage := make(map[string]NodeNetUsage, len(nodes))

	for _, conn := range net {
		outportTypeExpr, err := a.getSenderType(conn.SenderSide, compInterface.IO.In, nodes, nodesIfaces, scope)
		if err != nil {
			return nil, fmt.Errorf("get sender type: %w", err)
		}

		if conn.SenderSide.PortAddr != nil { // mark node's outport as used if sender isn't const ref
			senderNodeName := conn.SenderSide.PortAddr.Node
			outportName := conn.SenderSide.PortAddr.Port
			if _, ok := nodesUsage[senderNodeName]; !ok {
				nodesUsage[senderNodeName] = NodeNetUsage{
					In:  map[string]struct{}{}, // we don't use nodeIfaces for make with len
					Out: map[string]struct{}{}, // because sender could be const or io node (in/out)
				}
			}
			nodesUsage[senderNodeName].Out[outportName] = struct{}{}
		}

		for _, receiver := range conn.ReceiverSides {
			inportTypeExpr, err := a.getReceiverType(receiver, compInterface.IO.Out, nodes, nodesIfaces, scope)
			if err != nil {
				return nil, fmt.Errorf("get sen der type: %w", err)
			}

			if err := a.resolver.IsSubtypeOf(outportTypeExpr, inportTypeExpr, scope); err != nil {
				return nil, fmt.Errorf("is subtype of: %v -> %v: %w", conn.SenderSide, receiver, err)
			}

			// mark node's inport as used
			receiverNodeName := receiver.PortAddr.Node
			inportName := receiver.PortAddr.Port
			if _, ok := nodesUsage[receiverNodeName]; !ok {
				nodesUsage[receiverNodeName] = NodeNetUsage{
					In:  map[string]struct{}{}, // we don't use nodeIfaces for the same reason
					Out: map[string]struct{}{}, // as with outports
				}
			}
			nodesUsage[receiverNodeName].In[inportName] = struct{}{}
		}
	}

	if err := a.checkNodeUsage(nodesIfaces, scope, nodesUsage); err != nil {
		return nil, fmt.Errorf("check unused outports: %w", err)
	}

	return net, nil
}

// NodeNetUsage represents how network uses node's ports.
type NodeNetUsage struct {
	In, Out map[string]struct{}
}

// checkNodeUsage returns err if some node or node's outport is unused to avoid deadlocks.
func (Analyzer) checkNodeUsage(
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
	nodesUsage map[string]NodeNetUsage,
) error {
	for nodeName, nodeIface := range nodesIfaces {
		nodeUsage, ok := nodesUsage[nodeName]
		if !ok {
			return fmt.Errorf("%w: %v", ErrUnusedNode, nodeName)
		}

		for inportName := range nodeIface.IO.In {
			if _, ok := nodeUsage.In[inportName]; !ok {
				return fmt.Errorf("%w: %v.in.%v", ErrUnusedNodeInport, nodeName, inportName)
			}
		}

		for outportName := range nodeIface.IO.Out {
			if _, ok := nodeUsage.Out[outportName]; !ok {
				return fmt.Errorf("%w: %v.out.%v", ErrUnusedNodeOutport, nodeName, outportName)
			}
		}
	}

	return nil
}

func (a Analyzer) getReceiverType(
	receiverSide src.ReceiverConnectionSide,
	outports map[string]src.Port,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
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
	scope src.Scope,
) (ts.Expr, error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, fmt.Errorf("%w: %v", ErrNodeNotFound, portAddr.Node)
	}
	entity, _, err := scope.Entity(node.EntityRef) // nodes analyzed so we don't check error
	if err != nil {
		panic("")
	}

	typ, err := a.getResolvedPortType(
		entity.Component.Interface.IO.In,
		entity.Component.Interface.TypeParams.Params,
		portAddr,
		node,
		scope,
	)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("get resolved inport type: %w", err)
	}

	return typ, nil
}

func (a Analyzer) getResolvedPortType(
	ports map[string]src.Port,
	params []ts.Param,
	portAddr src.PortAddr,
	node src.Node,
	scope src.Scope,
) (ts.Expr, error) {
	port, ok := ports[portAddr.Port]
	if !ok {
		return ts.Expr{}, fmt.Errorf("%w: %v", ErrNodePortNotFound, portAddr)
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

func (a Analyzer) getSenderType(
	senderSide src.SenderConnectionSide,
	inports map[string]src.Port,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, error) {
	if senderSide.ConstRef != nil {
		constTypeExpr, err := a.getConstType(*senderSide.ConstRef, scope)
		if err != nil {
			return ts.Expr{}, fmt.Errorf("get const type: %w", err)
		}
		return constTypeExpr, nil
	}

	if senderSide.PortAddr == nil {
		return ts.Expr{}, ErrSenderIsEmpty
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

	nodeOutportType, err := a.getNodeOutportType(*senderSide.PortAddr, nodes, nodesIfaces, scope)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("get node outport type: %w", err)
	}

	return nodeOutportType, nil
}

func (a Analyzer) getNodeOutportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, fmt.Errorf("%w: %v", ErrNodeNotFound, portAddr.Node)
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, fmt.Errorf("%w: %v", ErrNodeNotFound, portAddr.Node)
	}

	typ, err := a.getResolvedPortType(
		nodeIface.IO.Out,
		nodeIface.TypeParams.Params,
		portAddr,
		node,
		scope,
	)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("get resolved outport type: %w", err)
	}

	return typ, err
}

func (a Analyzer) getConstType(ref src.EntityRef, scope src.Scope) (ts.Expr, error) {
	entity, _, err := scope.Entity(ref)
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

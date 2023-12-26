package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var (
	ErrNodeWrongEntity             = errors.New("Node can only refer to components or interfaces")
	ErrNodeTypeArgsCountMismatch   = errors.New("Type arguments count between node and its referenced entity not matches")
	ErrNonComponentNodeWithDI      = errors.New("Only component node can have dependency injection")
	ErrUnusedNode                  = errors.New("Unused node found")
	ErrUnusedNodeInport            = errors.New("Unused node inport found")
	ErrUnusedNodeOutport           = errors.New("Unused node outport found")
	ErrSenderIsEmpty               = errors.New("Sender in network must refer to some port address")
	ErrReadSelfOut                 = errors.New("Component cannot read from self outport")
	ErrWriteSelfIn                 = errors.New("Component cannot write to self inport")
	ErrInportNotFound              = errors.New("Referenced inport not found in component's interface")
	ErrNodeNotFound                = errors.New("Referenced node not found")
	ErrNodePortNotFound            = errors.New("Referenced node port not found")
	ErrNormCompWithRuntimeFunc     = errors.New("Component with nodes or network cannot use #runtime_func directive")
	ErrNormComponentWithoutNet     = errors.New("Component must have network except it uses #runtime_func directive")
	ErrNormNodeRuntimeMsg          = errors.New("Node can't use #runtime_func_msg if it isn't instantiated with the component that use #runtime_func") //nolint:lll
	ErrInterfaceNodeWithRuntimeMsg = errors.New("Interface node cannot use #runtime_func_msg directive")
	ErrRuntimeFuncDirectiveArgs    = errors.New("Component that use #runtime_func directive must provide exactly one argument") //nolint:lll
	ErrRuntimeMsgArgs              = errors.New("Node with #runtime_func_msg directive must provide exactly one argument")
)

type analyzeComponentParams struct {
	iface analyzeInterfaceParams
}

func (a Analyzer) analyzeComponent(component src.Component, scope src.Scope) (src.Component, *Error) { //nolint:funlen
	_, isRuntimeFunc := component.Directives[compiler.RuntimeFuncDirective]

	if isRuntimeFunc && len(component.Directives[compiler.RuntimeFuncDirective]) != 1 {
		return src.Component{}, &Error{
			Err:      ErrRuntimeFuncDirectiveArgs,
			Location: &scope.Location,
			Meta:     &component.Meta,
		}
	}

	resolvedInterface, err := a.analyzeInterface(component.Interface, scope, analyzeInterfaceParams{
		allowEmptyInports:  isRuntimeFunc,
		allowEmptyOutports: isRuntimeFunc,
	})
	if err != nil {
		return src.Component{}, Error{
			Location: &scope.Location,
			Meta:     &component.Meta,
		}.Merge(err)
	}

	if isRuntimeFunc {
		if len(component.Nodes) != 0 || len(component.Net) != 0 {
			return src.Component{}, &Error{
				Err:      ErrNormCompWithRuntimeFunc,
				Location: &scope.Location,
				Meta:     &component.Meta,
			}
		}
		return component, nil
	}

	resolvedNodes, nodesIfaces, err := a.analyzeComponentNodes(component.Nodes, scope)
	if err != nil {
		return src.Component{}, Error{
			Location: &scope.Location,
			Meta:     &component.Meta,
		}.Merge(err)
	}

	if len(component.Net) == 0 {
		return src.Component{}, &Error{
			Err:      ErrNormComponentWithoutNet,
			Location: &scope.Location,
			Meta:     &component.Meta,
		}
	}

	resolvedNet, err := a.analyzeComponentNetwork(component.Net, resolvedInterface, resolvedNodes, nodesIfaces, scope)
	if err != nil {
		return src.Component{}, Error{
			Location: &scope.Location,
			Meta:     &component.Meta,
		}.Merge(err)
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
) (map[string]src.Node, map[string]src.Interface, *Error) {
	resolvedNodes := make(map[string]src.Node, len(nodes))
	nodesInterfaces := make(map[string]src.Interface, len(nodes))

	for nodeName, node := range nodes {
		resolvedNode, nodeInterface, err := a.analyzeComponentNode(node, scope)
		if err != nil {
			return nil, nil, Error{
				Err:      fmt.Errorf("Invalid node: %v", nodeName),
				Location: &scope.Location,
				Meta:     &node.Meta,
			}.Merge(err)
		}

		nodesInterfaces[nodeName] = nodeInterface
		resolvedNodes[nodeName] = resolvedNode
	}

	return resolvedNodes, nodesInterfaces, nil
}

//nolint:funlen
func (a Analyzer) analyzeComponentNode(node src.Node, scope src.Scope) (src.Node, src.Interface, *Error) {
	entity, location, err := scope.Entity(node.EntityRef)
	if err != nil {
		return src.Node{}, src.Interface{}, &Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &node.Meta,
		}
	}

	if entity.Kind != src.ComponentEntity && entity.Kind != src.InterfaceEntity {
		return src.Node{}, src.Interface{}, &Error{
			Err:      fmt.Errorf("%w: %v", ErrNodeWrongEntity, entity.Kind),
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	runtimeMsgArgs, hasRuntimeMsg := node.Directives[compiler.RuntimeFuncMsgDirective]
	if hasRuntimeMsg && len(runtimeMsgArgs) != 1 {
		return src.Node{}, src.Interface{}, &Error{
			Err:      ErrRuntimeMsgArgs,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	var iface src.Interface
	if entity.Kind == src.ComponentEntity {
		_, isRuntimeFunc := entity.Component.Directives[compiler.RuntimeFuncDirective]

		if hasRuntimeMsg && !isRuntimeFunc {
			return src.Node{}, src.Interface{}, &Error{
				Err:      ErrNormNodeRuntimeMsg,
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		iface = entity.Component.Interface
	} else {
		if hasRuntimeMsg {
			return src.Node{}, src.Interface{}, &Error{
				Err:      ErrInterfaceNodeWithRuntimeMsg,
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		if node.Deps != nil {
			return src.Node{}, src.Interface{}, &Error{
				Err:      ErrNonComponentNodeWithDI,
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		iface = entity.Interface
	}

	if len(node.TypeArgs) != len(iface.TypeParams.Params) {
		return src.Node{}, src.Interface{}, &Error{
			Err: fmt.Errorf(
				"%w: want %v, got %v",
				ErrNodeTypeArgsCountMismatch, iface.TypeParams, node.TypeArgs,
			),
			Location: &location,
			Meta:     &node.Meta,
		}
	}

	resolvedArgs, _, err := a.resolver.ResolveFrame(node.TypeArgs, iface.TypeParams.Params, scope)
	if err != nil {
		return src.Node{}, src.Interface{}, &Error{
			Err:      err,
			Location: &location,
			Meta:     &node.Meta,
		}
	}

	if node.Deps == nil {
		return src.Node{
			EntityRef: node.EntityRef,
			TypeArgs:  resolvedArgs,
		}, iface, nil
	}

	resolvedComponentDI := make(map[string]src.Node, len(node.Deps))
	for depName, depNode := range node.Deps {
		resolvedDep, _, err := a.analyzeComponentNode(depNode, scope)
		if err != nil {
			return src.Node{}, src.Interface{}, Error{
				Err:      fmt.Errorf("Invalid node dependency: node '%v'", depNode),
				Location: &location,
				Meta:     &depNode.Meta,
			}.Merge(err)
		}
		resolvedComponentDI[depName] = resolvedDep
	}

	return src.Node{
		EntityRef: node.EntityRef,
		TypeArgs:  resolvedArgs,
		Deps:      resolvedComponentDI,
	}, iface, nil
}

//nolint:funlen
func (a Analyzer) analyzeComponentNetwork(
	net []src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) ([]src.Connection, *Error) {
	nodesUsage := make(map[string]NodeNetUsage, len(nodes))

	for _, conn := range net {
		outportTypeExpr, err := a.getSenderType(conn.SenderSide, compInterface.IO.In, nodes, nodesIfaces, scope)
		if err != nil {
			return nil, Error{
				Location: &scope.Location,
				Meta:     &conn.SenderSide.Meta,
			}.Merge(err)
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
				return nil, Error{
					Err:      errors.New("Unable to get receiver type"),
					Location: &scope.Location,
					Meta:     &receiver.Meta,
				}.Merge(err)
			}

			if err := a.resolver.IsSubtypeOf(outportTypeExpr, inportTypeExpr, scope); err != nil {
				return nil, &Error{
					Err: fmt.Errorf(
						"Subtype checking failed: sender %v, receiver %v, error %w",
						conn.SenderSide, receiver, err,
					),
					Location: &scope.Location,
					Meta:     &conn.Meta,
				}
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
		return nil, Error{Location: &scope.Location}.Merge(err)
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
) *Error {
	for nodeName, nodeIface := range nodesIfaces {
		nodeUsage, ok := nodesUsage[nodeName]
		if !ok {
			return &Error{
				Err:      fmt.Errorf("%w: %v", ErrUnusedNode, nodeName),
				Location: &scope.Location,
			}
		}

		for inportName := range nodeIface.IO.In {
			if _, ok := nodeUsage.In[inportName]; !ok {
				meta := nodeIface.IO.In[inportName].Meta
				return &Error{
					Err:      fmt.Errorf("%w: node '%v', inport '%v'", ErrUnusedNodeInport, nodeName, inportName),
					Location: &scope.Location,
					Meta:     &meta,
				}
			}
		}

		for outportName := range nodeIface.IO.Out {
			if _, ok := nodeUsage.Out[outportName]; !ok {
				meta := nodeIface.IO.Out[outportName].Meta
				return &Error{
					Err:      fmt.Errorf("%w: %v.out.%v", ErrUnusedNodeOutport, nodeName, outportName),
					Location: &scope.Location,
					Meta:     &meta,
				}
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
) (ts.Expr, *Error) {
	if receiverSide.PortAddr.Node == "in" {
		return ts.Expr{}, &Error{
			Err:      ErrWriteSelfIn,
			Location: &scope.Location,
			Meta:     &receiverSide.PortAddr.Meta,
		}
	}

	if receiverSide.PortAddr.Node == "out" {
		outport, ok := outports[receiverSide.PortAddr.Port]
		if !ok {
			return ts.Expr{}, &Error{
				Err:      ErrInportNotFound,
				Location: &scope.Location,
				Meta:     &receiverSide.PortAddr.Meta,
			}
		}
		return outport.TypeExpr, nil
	}

	nodeInportType, err := a.getNodeInportType(receiverSide.PortAddr, nodes, scope)
	if err != nil {
		return ts.Expr{}, &Error{
			Err:      fmt.Errorf("get node inport type: %w", err),
			Location: &scope.Location,
			Meta:     &receiverSide.PortAddr.Meta,
		}
	}

	return nodeInportType, nil
}

func (a Analyzer) getNodeInportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	scope src.Scope,
) (ts.Expr, *Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, &Error{
			Err:      fmt.Errorf("%w: %v", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	entity, _, err := scope.Entity(node.EntityRef)
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
		return ts.Expr{}, &Error{
			Err:      fmt.Errorf("Unable to ger resolved port type: %w", err),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	return typ, nil
}

func (a Analyzer) getResolvedPortType(
	ports map[string]src.Port,
	params []ts.Param,
	portAddr src.PortAddr,
	node src.Node,
	scope src.Scope,
) (ts.Expr, *Error) {
	port, ok := ports[portAddr.Port]
	if !ok {
		return ts.Expr{}, &Error{
			Err:      fmt.Errorf("%w: %v", ErrNodePortNotFound, portAddr),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	_, frame, err := a.resolver.ResolveFrame(node.TypeArgs, params, scope)
	if err != nil {
		return ts.Expr{}, &Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &node.Meta,
		}
	}

	resolvedOutportType, err := a.resolver.ResolveExprWithFrame(port.TypeExpr, frame, scope)
	if err != nil {
		return ts.Expr{}, &Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &node.Meta,
		}
	}

	return resolvedOutportType, nil
}

func (a Analyzer) getSenderType(
	senderSide src.SenderConnectionSide,
	inports map[string]src.Port,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, *Error) {
	if senderSide.PortAddr == nil {
		return ts.Expr{}, &Error{
			Err:      ErrSenderIsEmpty,
			Location: &scope.Location,
			Meta:     &senderSide.Meta,
		}
	}

	if senderSide.PortAddr.Node == "out" {
		return ts.Expr{}, &Error{
			Err:      ErrReadSelfOut,
			Location: &scope.Location,
			Meta:     &senderSide.PortAddr.Meta,
		}
	}

	if senderSide.PortAddr.Node == "in" {
		inport, ok := inports[senderSide.PortAddr.Port]
		if !ok {
			return ts.Expr{}, &Error{
				Err:      ErrInportNotFound,
				Location: &scope.Location,
				Meta:     &senderSide.PortAddr.Meta,
			}
		}
		return inport.TypeExpr, nil
	}

	nodeOutportType, err := a.getNodeOutportType(*senderSide.PortAddr, nodes, nodesIfaces, scope)
	if err != nil {
		return ts.Expr{}, Error{
			Err:      ErrInportNotFound,
			Location: &scope.Location,
			Meta:     &senderSide.PortAddr.Meta,
		}.Merge(err)
	}

	return nodeOutportType, nil
}

func (a Analyzer) getNodeOutportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, *Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, &Error{
			Err:      fmt.Errorf("%w: %v", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, &Error{
			Err:      fmt.Errorf("%w: %v", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	typ, err := a.getResolvedPortType(
		nodeIface.IO.Out,
		nodeIface.TypeParams.Params,
		portAddr,
		node,
		scope,
	)
	if err != nil {
		return ts.Expr{}, &Error{
			Err:      fmt.Errorf("get resolved outport type: %w", err),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	return typ, err
}

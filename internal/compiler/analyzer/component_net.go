package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var ErrStructFieldNotFound = errors.New("Struct field not found")

func (a Analyzer) analyzeComponentNetwork(
	net []src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) ([]src.Connection, *compiler.Error) {
	nodesUsage := make(map[string]NodeNetUsage, len(nodes)) // we create it here because there's recursion down there

	if err := a.analyzeConnections(net, compInterface, nodes, nodesIfaces, nodesUsage, scope); err != nil {
		return nil, compiler.Error{Location: &scope.Location}.Merge(err)
	}

	if err := a.checkNodesPortsUsage(nodesIfaces, scope, nodesUsage); err != nil {
		return nil, compiler.Error{Location: &scope.Location}.Merge(err)
	}

	return net, nil
}

// analyzeConnections does two things:
// 1. Analyzes every connection and terminates with non-nil error if any of them is invalid.
// 2. Updates nodesUsage (we mutate it in-place instead of returning to avoid merging across recursive calls).
func (a Analyzer) analyzeConnections(
	net []src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	nodesUsage map[string]NodeNetUsage,
	scope src.Scope,
) *compiler.Error {
	for _, conn := range net {
		if err := a.analyzeNetConn(conn, compInterface, nodes, nodesIfaces, scope, nodesUsage); err != nil {
			return err
		}
	}
	return nil
}

func (a Analyzer) analyzeNetConn( //nolint:funlen
	conn src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
	nodesUsage map[string]NodeNetUsage,
) *compiler.Error {
	outportTypeExpr, err := a.getSenderType(conn.SenderSide, compInterface.IO.In, nodes, nodesIfaces, scope)
	if err != nil {
		return compiler.Error{
			Location: &scope.Location,
			Meta:     &conn.SenderSide.Meta,
		}.Merge(err)
	}

	// mark node's outport as used if sender isn't const ref
	if conn.SenderSide.PortAddr != nil {
		senderNodeName := conn.SenderSide.PortAddr.Node
		outportName := conn.SenderSide.PortAddr.Port
		if _, ok := nodesUsage[senderNodeName]; !ok {
			nodesUsage[senderNodeName] = NodeNetUsage{
				In:  map[string]struct{}{},
				Out: map[string]struct{}{},
			}
		}
		nodesUsage[senderNodeName].Out[outportName] = struct{}{}
	}

	if len(conn.SenderSide.Selectors) > 0 {
		lastFieldType, err := ts.GetStructFieldTypeByPath(outportTypeExpr, conn.SenderSide.Selectors)
		if err != nil {
			return &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}
		outportTypeExpr = lastFieldType
	}

	if len(conn.ReceiverSide.ThenConnections) == 0 && len(conn.ReceiverSide.Receivers) == 0 {
		if err != nil {
			return &compiler.Error{
				Err: errors.New(
					"Connection's receiver side cannot be empty, it must either have then connection or receivers",
				),
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}
	} else if len(conn.ReceiverSide.ThenConnections) != 0 && len(conn.ReceiverSide.Receivers) != 0 {
		if err != nil {
			return &compiler.Error{
				Err: errors.New(
					"Connection's receiver side must either have then connection or receivers, not both",
				),
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}
	}

	if conn.ReceiverSide.ThenConnections != nil {
		// note that we call analyzeConnections instead of analyzeComponentNetwork
		// because we only need to analyze connections and update nodesUsage
		// analyzeComponentNetwork OTOH will also validate nodesUsage by itself
		return a.analyzeConnections( // indirect recursion
			conn.ReceiverSide.ThenConnections,
			compInterface,
			nodes,
			nodesIfaces,
			nodesUsage,
			scope,
		)
	}

	for _, receiver := range conn.ReceiverSide.Receivers {
		inportTypeExpr, err := a.resolveReceiverType(
			receiver,
			compInterface.IO.Out,
			nodes,
			nodesIfaces,
			scope,
		)
		if err != nil {
			return compiler.Error{
				Err:      errors.New("Bad receiver"),
				Location: &scope.Location,
				Meta:     &receiver.Meta,
			}.Merge(err)
		}

		if err := a.resolver.IsSubtypeOf(outportTypeExpr, inportTypeExpr, scope); err != nil {
			return &compiler.Error{
				Err: fmt.Errorf(
					"Subtype checking failed: %v -> %v: %w",
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

	return nil
}

// NodeNetUsage represents how network uses node's ports.
type NodeNetUsage struct {
	In, Out map[string]struct{}
}

// checkNodesPortsUsage ensures that for every node out there we use all its inports and outports.
func (Analyzer) checkNodesPortsUsage(
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
	nodesUsage map[string]NodeNetUsage,
) *compiler.Error {
	for nodeName, nodeIface := range nodesIfaces {
		nodeUsage, ok := nodesUsage[nodeName]
		if !ok {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrUnusedNode, nodeName),
				Location: &scope.Location,
			}
		}

		for inportName := range nodeIface.IO.In {
			if _, ok := nodeUsage.In[inportName]; !ok {
				meta := nodeIface.IO.In[inportName].Meta
				return &compiler.Error{
					Err:      fmt.Errorf("%w: node '%v', inport '%v'", ErrUnusedNodeInport, nodeName, inportName),
					Location: &scope.Location,
					Meta:     &meta,
				}
			}
		}

		if len(nodeIface.IO.Out) == 0 { // such components exist in stdlib, user cannot create them
			continue
		}

		atLeastOneOutportIsUsed := false
		for outportName := range nodeIface.IO.Out {
			if _, ok := nodeUsage.Out[outportName]; ok {
				atLeastOneOutportIsUsed = true
			}
		}
		if !atLeastOneOutportIsUsed {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrUnusedNodeOutports, nodeName),
				Location: &scope.Location,
				Meta:     &nodeIface.Meta,
			}
		}
	}

	return nil
}

func (a Analyzer) resolveReceiverType(
	receiverSide src.ConnectionReceiver,
	outports map[string]src.Port,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	if receiverSide.PortAddr.Node == "in" {
		return ts.Expr{}, &compiler.Error{
			Err:      ErrWriteSelfIn,
			Location: &scope.Location,
			Meta:     &receiverSide.PortAddr.Meta,
		}
	}

	if receiverSide.PortAddr.Node == "out" {
		outport, ok := outports[receiverSide.PortAddr.Port]
		if !ok {
			return ts.Expr{}, &compiler.Error{
				Err:      ErrInportNotFound,
				Location: &scope.Location,
				Meta:     &receiverSide.PortAddr.Meta,
			}
		}
		return outport.TypeExpr, nil
	}

	nodeInportType, err := a.getNodeInportType(receiverSide.PortAddr, nodes, nodesIfaces, scope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &receiverSide.PortAddr.Meta,
		}
	}

	return nodeInportType, nil
}

func (a Analyzer) getNodeInportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("Node not found '%v'", portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	iface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w '%v'", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	// TODO optimize: we can resolve every node's interface just once before processing the network
	typ, aerr := a.getResolvedPortType(
		iface.IO.In,
		iface.TypeParams.Params,
		portAddr,
		node,
		scope,
	)
	if aerr != nil {
		return ts.Expr{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}.Merge(aerr)
	}

	return typ, nil
}

func (a Analyzer) getResolvedPortType(
	ports map[string]src.Port,
	params []ts.Param,
	portAddr src.PortAddr,
	node src.Node,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	port, ok := ports[portAddr.Port]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w '%v'", ErrPortNotFound, portAddr),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	_, frame, err := a.resolver.ResolveFrame(node.TypeArgs, params, scope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &node.Meta,
		}
	}

	resolvedOutportType, err := a.resolver.ResolveExprWithFrame(port.TypeExpr, frame, scope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &node.Meta,
		}
	}

	return resolvedOutportType, nil
}

func (a Analyzer) getSenderType(
	senderSide src.ConnectionSenderSide,
	inports map[string]src.Port,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	if senderSide.PortAddr == nil && senderSide.ConstRef == nil {
		return ts.Expr{}, &compiler.Error{
			Err:      ErrSenderIsEmpty,
			Location: &scope.Location,
			Meta:     &senderSide.Meta,
		}
	}

	if senderSide.ConstRef != nil {
		expr, err := a.getResolvedConstType(*senderSide.ConstRef, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &senderSide.ConstRef.Meta,
			}.Merge(err)
		}
		return expr, nil
	}

	if senderSide.PortAddr.Node == "out" {
		return ts.Expr{}, &compiler.Error{
			Err:      ErrReadSelfOut,
			Location: &scope.Location,
			Meta:     &senderSide.PortAddr.Meta,
		}
	}

	if senderSide.PortAddr.Node == "in" {
		inport, ok := inports[senderSide.PortAddr.Port]
		if !ok {
			return ts.Expr{}, &compiler.Error{
				Err:      ErrInportNotFound,
				Location: &scope.Location,
				Meta:     &senderSide.PortAddr.Meta,
			}
		}
		return inport.TypeExpr, nil
	}

	nodeOutportType, err := a.getNodeOutportType(*senderSide.PortAddr, nodes, nodesIfaces, scope)
	if err != nil {
		return ts.Expr{}, compiler.Error{
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
) (ts.Expr, *compiler.Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, &compiler.Error{
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
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("get resolved outport type: %w", err),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	return typ, err
}

func (a Analyzer) getResolvedConstType(ref src.EntityRef, scope src.Scope) (ts.Expr, *compiler.Error) {
	entity, location, err := scope.Entity(ref)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &ref.Meta,
		}
	}

	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", errors.New("Entity found but is not constant"), entity.Kind),
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if entity.Const.Ref != nil {
		expr, err := a.getResolvedConstType(*entity.Const.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Location: &location,
				Meta:     &entity.Const.Meta,
			}.Merge(err)
		}
		return expr, nil
	}

	resolvedExpr, err := a.resolver.ResolveExpr(entity.Const.Value.TypeExpr, scope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &entity.Const.Value.Meta,
		}
	}

	return resolvedExpr, nil
}

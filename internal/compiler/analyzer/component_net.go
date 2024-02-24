package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var (
	ErrStructFieldNotFound    = errors.New("Struct field not found")
	ErrUnusedOutports         = errors.New("All component's outports are unused")
	ErrUnusedOutport          = errors.New("Unused outport found")
	ErrUnusedInports          = errors.New("All component inports are unused")
	ErrUnusedInport           = errors.New("Unused inport found")
	ErrLiteralSenderTypeEmpty = errors.New("Literal network sender must contain message value")
	ErrLiteralSenderKind      = errors.New("Literal network sender must have type of kind instantiation")
	ErrLiteralSenderType      = errors.New("Literal network sender must have primitive type")
)

// analyzeComponentNetwork must be called after analyzeNodes so we sure nodes are resolved.
func (a Analyzer) analyzeComponentNetwork(
	net []src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) ([]src.Connection, *compiler.Error) {
	nodesUsage := make(map[string]NodeNetUsage, len(nodes)) // we create it here because there's recursion down there

	if err := a.analyzeConnections(net, compInterface, nodes, nodesIfaces, nodesUsage, scope); err != nil {
		return nil, compiler.Error{Location: &scope.Location}.Wrap(err)
	}

	if err := a.checkNetPortsUsage(compInterface, nodesIfaces, scope, nodesUsage); err != nil {
		return nil, compiler.Error{Location: &scope.Location}.Wrap(err)
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
		if err := a.analyzeConnection(conn, compInterface, nodes, nodesIfaces, scope, nodesUsage); err != nil {
			return err
		}
	}
	return nil
}

func (a Analyzer) analyzeConnection( //nolint:funlen
	conn src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
	nodesUsage map[string]NodeNetUsage,
) *compiler.Error {
	outportTypeExpr, err := a.getSenderType(conn.Normal.SenderSide, compInterface, nodes, nodesIfaces, scope)
	if err != nil {
		return compiler.Error{
			Location: &scope.Location,
			Meta:     &conn.Normal.SenderSide.Meta,
		}.Wrap(err)
	}

	// mark node's outport as used if sender isn't const ref
	if conn.Normal.SenderSide.PortAddr != nil {
		senderNodeName := conn.Normal.SenderSide.PortAddr.Node
		outportName := conn.Normal.SenderSide.PortAddr.Port
		if _, ok := nodesUsage[senderNodeName]; !ok {
			nodesUsage[senderNodeName] = NodeNetUsage{
				In:  map[string]struct{}{},
				Out: map[string]struct{}{},
			}
		}
		nodesUsage[senderNodeName].Out[outportName] = struct{}{}
	}

	if len(conn.Normal.SenderSide.Selectors) > 0 {
		lastFieldType, err := ts.GetStructFieldTypeByPath(outportTypeExpr, conn.Normal.SenderSide.Selectors)
		if err != nil {
			return &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}
		outportTypeExpr = lastFieldType
	}

	if len(conn.Normal.ReceiverSide.ThenConnections) == 0 && len(conn.Normal.ReceiverSide.Receivers) == 0 {
		if err != nil {
			return &compiler.Error{
				Err: errors.New(
					"Connection's receiver side cannot be empty, it must either have then connection or receivers",
				),
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}
	} else if len(conn.Normal.ReceiverSide.ThenConnections) != 0 && len(conn.Normal.ReceiverSide.Receivers) != 0 {
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

	if conn.Normal.ReceiverSide.ThenConnections != nil {
		// note that we call analyzeConnections instead of analyzeComponentNetwork
		// because we only need to analyze connections and update nodesUsage
		// analyzeComponentNetwork OTOH will also validate nodesUsage by itself
		return a.analyzeConnections( // indirect recursion
			conn.Normal.ReceiverSide.ThenConnections,
			compInterface,
			nodes,
			nodesIfaces,
			nodesUsage,
			scope,
		)
	}

	for _, receiver := range conn.Normal.ReceiverSide.Receivers {
		inportTypeExpr, err := a.getReceiverType(
			receiver,
			compInterface,
			nodes,
			nodesIfaces,
			scope,
		)
		if err != nil {
			return compiler.Error{
				Location: &scope.Location,
				Meta:     &receiver.Meta,
			}.Wrap(err)
		}

		if err := a.resolver.IsSubtypeOf(outportTypeExpr, inportTypeExpr, scope); err != nil {
			return &compiler.Error{
				Err: fmt.Errorf(
					"Subtype checking failed: %v -> %v: %w",
					conn.Normal.SenderSide, receiver, err,
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

// checkNetPortsUsage ensures that:
// Every component's inport and outport is used;
// Every sub-node's inport is used;
// For every  sub-node's there's at least one used outport.
func (Analyzer) checkNetPortsUsage( //nolint:funlen
	compInterface src.Interface,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
	nodesUsage map[string]NodeNetUsage,
) *compiler.Error {
	inportsUsage, ok := nodesUsage["in"]
	if !ok {
		return &compiler.Error{
			Err:      ErrUnusedInports,
			Location: &scope.Location,
			Meta:     &compInterface.Meta,
		}
	}
	for inportName := range compInterface.IO.In {
		if _, ok := inportsUsage.Out[inportName]; !ok { // note that self inports are outports for the network
			return &compiler.Error{
				Err:      fmt.Errorf("%w '%v'", ErrUnusedInport, inportName),
				Location: &scope.Location,
			}
		}
	}

	outportsUsage, ok := nodesUsage["out"]
	if !ok {
		return &compiler.Error{
			Err:      ErrUnusedOutports,
			Location: &scope.Location,
		}
	}
	for outportName := range compInterface.IO.Out {
		if _, ok := outportsUsage.In[outportName]; !ok { // note that self outports are inports for the network
			return &compiler.Error{
				Err:      fmt.Errorf("%w '%v'", ErrUnusedOutport, outportName),
				Location: &scope.Location,
			}
		}
	}

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

		if len(nodeIface.IO.Out) == 0 { // std/builtin.Void
			continue
		}

		atLeastOneOutportIsUsed := false
		for outportName := range nodeIface.IO.Out {
			if _, ok := nodeUsage.Out[outportName]; ok {
				atLeastOneOutportIsUsed = true
				break
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

func (a Analyzer) getReceiverType(
	receiverSide src.ConnectionReceiver,
	iface src.Interface,
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
		outports := iface.IO.Out

		outport, ok := outports[receiverSide.PortAddr.Port]
		if !ok {
			return ts.Expr{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrOutportNotFound, receiverSide.PortAddr.Port),
				Location: &scope.Location,
				Meta:     &receiverSide.PortAddr.Meta,
			}
		}

		resolvedOutportType, err := a.resolver.ResolveExprWithFrame(
			outport.TypeExpr,
			iface.TypeParams.ToFrame(),
			scope,
		)
		if err != nil {
			return ts.Expr{}, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &receiverSide.PortAddr.Meta,
			}
		}

		return resolvedOutportType, nil
	}

	nodeInportType, err := a.getNodeInportType(receiverSide.PortAddr, nodes, nodesIfaces, scope)
	if err != nil {
		return ts.Expr{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &receiverSide.PortAddr.Meta,
		}.Wrap(err)
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

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w '%v'", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	// TODO optimize:
	// we can resolve every node's interface just once
	// before processing the network
	resolvedInportType, aerr := a.getResolvedPortType(
		nodeIface.IO.In,
		nodeIface.TypeParams.Params,
		portAddr,
		node,
		scope,
	)
	if aerr != nil {
		return ts.Expr{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}.Wrap(aerr)
	}

	return resolvedInportType, nil
}

func (a Analyzer) getResolvedPortType(
	ports map[string]src.Port,
	nodeIfaceParams []ts.Param,
	portAddr src.PortAddr,
	node src.Node,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	port, ok := ports[portAddr.Port]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err: fmt.Errorf(
				"Port not found `%v`",
				portAddr,
			),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	// we don't resolve node's args assuming they resolved already

	// create frame `nodeParam:resolvedArg` to get resolved port type
	frame := make(map[string]ts.Def, len(nodeIfaceParams))
	for i, param := range nodeIfaceParams {
		arg := node.TypeArgs[i]
		frame[param.Name] = ts.Def{
			BodyExpr: &arg,
			Meta:     arg.Meta,
		}
	}

	resolvedPortType, err := a.resolver.ResolveExprWithFrame(
		port.TypeExpr,
		frame,
		scope,
	)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &port.Meta,
		}
	}

	return resolvedPortType, nil
}

func (a Analyzer) getSenderType( //nolint:funlen
	senderSide src.ConnectionSenderSide,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	if senderSide.PortAddr == nil && senderSide.Const == nil {
		return ts.Expr{}, &compiler.Error{
			Err:      ErrSenderIsEmpty,
			Location: &scope.Location,
			Meta:     &senderSide.Meta,
		}
	}

	if senderSide.Const != nil {
		return a.getResolvedSenderConstType(senderSide, scope)
	}

	if senderSide.PortAddr.Node == "out" {
		return ts.Expr{}, &compiler.Error{
			Err:      ErrReadSelfOut,
			Location: &scope.Location,
			Meta:     &senderSide.PortAddr.Meta,
		}
	}

	if senderSide.PortAddr.Node == "in" {
		inports := iface.IO.In

		inport, ok := inports[senderSide.PortAddr.Port]
		if !ok {
			return ts.Expr{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrInportNotFound, senderSide.PortAddr.Port),
				Location: &scope.Location,
				Meta:     &senderSide.PortAddr.Meta,
			}
		}

		resolvedInportType, err := a.resolver.ResolveExprWithFrame(
			inport.TypeExpr,
			iface.TypeParams.ToFrame(),
			scope,
		)
		if err != nil {
			return ts.Expr{}, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &senderSide.PortAddr.Meta,
			}
		}

		return resolvedInportType, nil
	}

	nodeOutportType, err := a.getNodeOutportType(*senderSide.PortAddr, nodes, nodesIfaces, scope)
	if err != nil {
		return ts.Expr{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &senderSide.PortAddr.Meta,
		}.Wrap(err)
	}

	return nodeOutportType, nil
}

func (a Analyzer) getResolvedSenderConstType(
	senderSide src.ConnectionSenderSide,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	if senderSide.Const.Ref != nil {
		expr, err := a.getResolvedConstTypeByRef(*senderSide.Const.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &senderSide.Const.Ref.Meta,
			}.Wrap(err)
		}
		return expr, nil
	}
	if senderSide.Const.Value != nil {
		if err := a.validateLiteralSender(senderSide.Const); err != nil {
			return ts.Expr{}, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &senderSide.Const.Value.Meta,
			}
		}
		resolvedExpr, err := a.resolver.ResolveExpr(senderSide.Const.Value.TypeExpr, scope)
		if err != nil {
			return ts.Expr{}, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &senderSide.Const.Value.Meta,
			}
		}
		return resolvedExpr, nil
	}
	return ts.Expr{}, &compiler.Error{
		Err:      ErrLiteralSenderTypeEmpty,
		Location: &scope.Location,
		Meta:     &senderSide.Meta,
	}
}

func (a Analyzer) validateLiteralSender(cnst *src.Const) error {
	if cnst.Value.TypeExpr.Inst == nil {
		return ErrLiteralSenderKind
	}
	switch cnst.Value.TypeExpr.Inst.Ref.String() {
	case "bool", "int", "float", "string":
		return nil
	}
	return ErrLiteralSenderType
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

	resolvedPortType, err := a.getResolvedPortType(
		nodeIface.IO.Out,
		nodeIface.TypeParams.Params,
		portAddr,
		node,
		scope,
	)
	if err != nil {
		return ts.Expr{}, compiler.Error{
			Err:      fmt.Errorf("get resolved outport type: %v", portAddr),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}.Wrap(err)
	}

	return resolvedPortType, err
}

func (a Analyzer) getResolvedConstTypeByRef(ref src.EntityRef, scope src.Scope) (ts.Expr, *compiler.Error) {
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
		expr, err := a.getResolvedConstTypeByRef(*entity.Const.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Location: &location,
				Meta:     &entity.Const.Meta,
			}.Wrap(err)
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

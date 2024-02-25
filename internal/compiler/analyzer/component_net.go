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
	nodesUsage := make(map[string]nodeNetUsage, len(nodes)) // we create it here because there's recursion down there

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
	nodesUsage map[string]nodeNetUsage,
	scope src.Scope,
) *compiler.Error {
	for _, conn := range net {
		if err := a.analyzeConnection(conn, compInterface, nodes, nodesIfaces, scope, nodesUsage); err != nil {
			return err
		}
	}
	return nil
}

type nodesNetUsage map[string]nodeNetUsage

func (n nodesNetUsage) AddOutport(node, port string) {
	if _, ok := n[node]; !ok {
		defaultValue := nodeNetUsage{
			In:  map[string]struct{}{},
			Out: map[string]struct{}{},
		}
		n[node] = defaultValue
	}
	n[node].Out[port] = struct{}{}
}

func (n nodesNetUsage) AddInport(node, port string) {
	if _, ok := n[node]; !ok {
		defaultValue := nodeNetUsage{
			In:  map[string]struct{}{},
			Out: map[string]struct{}{},
		}
		n[node] = defaultValue
	}
	n[node].In[port] = struct{}{}
}

func (a Analyzer) analyzeConnection( //nolint:funlen
	conn src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
	nodesUsage map[string]nodeNetUsage,
) *compiler.Error {
	if arrBypassConn := conn.ArrayBypass; arrBypassConn != nil {
		senderType, isArray, err := a.getSenderPortAddrType(
			arrBypassConn.SenderOutport,
			scope,
			compInterface,
			nodes,
			nodesIfaces,
		)
		if err != nil {
			return compiler.Error{
				Location: &scope.Location,
				Meta:     &conn.Normal.SenderSide.Meta,
			}.Wrap(err)
		}
		if !isArray {
			return &compiler.Error{
				Err:      errors.New("Non-array outport in array-bypass connection"),
				Location: &scope.Location,
				Meta:     &arrBypassConn.SenderOutport.Meta,
			}
		}

		receiverType, isArray, err := a.getReceiverType(
			arrBypassConn.ReceiverInport,
			compInterface,
			nodes,
			nodesIfaces,
			scope,
		)
		if err != nil {
			return compiler.Error{
				Location: &scope.Location,
				Meta:     &conn.Normal.SenderSide.Meta,
			}.Wrap(err)
		}
		if !isArray {
			return &compiler.Error{
				Err:      errors.New("Non-array outport in array-bypass connection"),
				Location: &scope.Location,
				Meta:     &arrBypassConn.SenderOutport.Meta,
			}
		}

		if err := a.resolver.IsSubtypeOf(
			senderType,
			receiverType,
			scope,
		); err != nil {
			return &compiler.Error{
				Err: fmt.Errorf(
					"Incompatible types: %v -> %v: %w",
					arrBypassConn.SenderOutport, arrBypassConn.ReceiverInport, err,
				),
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}

		nodesNetUsage(nodesUsage).AddInport(
			arrBypassConn.SenderOutport.Node,
			arrBypassConn.SenderOutport.Port,
		)
		nodesNetUsage(nodesUsage).AddOutport(
			arrBypassConn.ReceiverInport.Node,
			arrBypassConn.ReceiverInport.Port,
		)

		return nil
	}

	outportTypeExpr, err := a.getSenderSideType(
		conn.Normal.SenderSide,
		compInterface,
		nodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return compiler.Error{
			Location: &scope.Location,
			Meta:     &conn.Normal.SenderSide.Meta,
		}.Wrap(err)
	}

	// mark node's outport as used if sender isn't const ref
	if conn.Normal.SenderSide.PortAddr != nil {
		nodesNetUsage(nodesUsage).AddOutport(
			conn.Normal.SenderSide.PortAddr.Node,
			conn.Normal.SenderSide.PortAddr.Port,
		)
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
		inportTypeExpr, _, err := a.getReceiverType(
			receiver.PortAddr,
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
					"Incompatible types: %v -> %v: %w",
					conn.Normal.SenderSide, receiver, err,
				),
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}

		nodesNetUsage(nodesUsage).AddInport(
			receiver.PortAddr.Node,
			receiver.PortAddr.Port,
		)
	}

	return nil
}

// nodeNetUsage shows which ports was used by the network
type nodeNetUsage struct {
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
	nodesUsage map[string]nodeNetUsage,
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
	receiverSide src.PortAddr,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, bool, *compiler.Error) {
	if receiverSide.Node == "in" {
		return ts.Expr{}, false, &compiler.Error{
			Err:      ErrWriteSelfIn,
			Location: &scope.Location,
			Meta:     &receiverSide.Meta,
		}
	}

	if receiverSide.Node == "out" {
		outports := iface.IO.Out

		outport, ok := outports[receiverSide.Port]
		if !ok {
			return ts.Expr{}, false, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrOutportNotFound, receiverSide.Port),
				Location: &scope.Location,
				Meta:     &receiverSide.Meta,
			}
		}

		resolvedOutportType, err := a.resolver.ResolveExprWithFrame(
			outport.TypeExpr,
			iface.TypeParams.ToFrame(),
			scope,
		)
		if err != nil {
			return ts.Expr{}, false, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &receiverSide.Meta,
			}
		}

		return resolvedOutportType, outport.IsArray, nil
	}

	nodeInportType, isArray, err := a.getNodeInportType(receiverSide, nodes, nodesIfaces, scope)
	if err != nil {
		return ts.Expr{}, false, compiler.Error{
			Location: &scope.Location,
			Meta:     &receiverSide.Meta,
		}.Wrap(err)
	}

	return nodeInportType, isArray, nil
}

func (a Analyzer) getNodeInportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, bool, *compiler.Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, false, &compiler.Error{
			Err:      fmt.Errorf("Node not found '%v'", portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, false, &compiler.Error{
			Err:      fmt.Errorf("%w '%v'", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	// TODO optimize:
	// we can resolve every node's interface just once
	// before processing the network
	resolvedInportType, isArray, aerr := a.getResolvedPortType(
		nodeIface.IO.In,
		nodeIface.TypeParams.Params,
		portAddr,
		node,
		scope,
	)
	if aerr != nil {
		return ts.Expr{}, false, compiler.Error{
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}.Wrap(aerr)
	}

	return resolvedInportType, isArray, nil
}

// getResolvedPortType returns port's type and isArray bool
func (a Analyzer) getResolvedPortType(
	ports map[string]src.Port,
	nodeIfaceParams []ts.Param,
	portAddr src.PortAddr,
	node src.Node,
	scope src.Scope,
) (ts.Expr, bool, *compiler.Error) {
	port, ok := ports[portAddr.Port]
	if !ok {
		return ts.Expr{}, false, &compiler.Error{
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
		return ts.Expr{}, false, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &port.Meta,
		}
	}

	return resolvedPortType, port.IsArray, nil
}

func (a Analyzer) getSenderSideType( //nolint:funlen
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

	typ, _, err := a.getSenderPortAddrType(*senderSide.PortAddr, scope, iface, nodes, nodesIfaces)
	return typ, err
}

// getSenderPortAddrType returns port's type and isArray bool
func (a Analyzer) getSenderPortAddrType(
	senderSide src.PortAddr,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
) (ts.Expr, bool, *compiler.Error) {
	if senderSide.Node == "out" {
		return ts.Expr{}, false, &compiler.Error{
			Err:      ErrReadSelfOut,
			Location: &scope.Location,
			Meta:     &senderSide.Meta,
		}
	}

	if senderSide.Node == "in" {
		inports := iface.IO.In

		inport, ok := inports[senderSide.Port]
		if !ok {
			return ts.Expr{}, false, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrInportNotFound, senderSide.Port),
				Location: &scope.Location,
				Meta:     &senderSide.Meta,
			}
		}

		resolvedInportType, err := a.resolver.ResolveExprWithFrame(
			inport.TypeExpr,
			iface.TypeParams.ToFrame(),
			scope,
		)
		if err != nil {
			return ts.Expr{}, false, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &senderSide.Meta,
			}
		}

		return resolvedInportType, inport.IsArray, nil
	}

	return a.getNodeOutportType(
		senderSide, nodes, nodesIfaces, scope,
	)
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

// getNodeOutportType returns port's type and isArray bool
func (a Analyzer) getNodeOutportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) (ts.Expr, bool, *compiler.Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, false, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, false, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	return a.getResolvedPortType(
		nodeIface.IO.Out,
		nodeIface.TypeParams.Params,
		portAddr,
		node,
		scope,
	)
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

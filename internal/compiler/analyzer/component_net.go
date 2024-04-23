package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

var (
	ErrStructFieldNotFound       = errors.New("Struct field not found")
	ErrUnusedOutports            = errors.New("All component's outports are unused")
	ErrUnusedOutport             = errors.New("Unused outport found")
	ErrUnusedInports             = errors.New("All component inports are unused")
	ErrUnusedInport              = errors.New("Unused inport found")
	ErrLiteralSenderTypeEmpty    = errors.New("Literal network sender must contain message value")
	ErrComplexLiteralSender      = errors.New("Literal network sender must have primitive type")
	ErrIllegalPortlessConnection = errors.New("Connection to a node, with more than one port, must always has a port name")
)

// analyzeComponentNetwork must be called after analyzeNodes so we sure nodes are resolved.
func (a Analyzer) analyzeComponentNetwork(
	net []src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
	scope src.Scope,
) ([]src.Connection, *compiler.Error) {
	// we create it here because there's recursion down there
	nodesUsage := make(map[string]nodeNetUsage, len(nodes))

	resolvedNet, err := a.analyzeConnections(net, compInterface, nodes, nodesIfaces, nodesUsage, scope)
	if err != nil {
		return nil, compiler.Error{Location: &scope.Location}.Wrap(err)
	}

	if err := a.checkNetPortsUsage(compInterface, nodesIfaces, scope, nodesUsage); err != nil {
		return nil, compiler.Error{Location: &scope.Location}.Wrap(err)
	}

	return resolvedNet, nil
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
) ([]src.Connection, *compiler.Error) {
	resolvedNet := make([]src.Connection, 0, len(net))

	for _, conn := range net {
		resolvedConn, err := a.analyzeConnection(
			conn,
			compInterface,
			nodes,
			nodesIfaces,
			scope,
			nodesUsage,
		)
		if err != nil {
			return nil, err
		}
		resolvedNet = append(resolvedNet, resolvedConn)
	}

	return resolvedNet, nil
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
) (src.Connection, *compiler.Error) {
	// first handle array bypass connection, they are simple
	if conn.ArrayBypass != nil {
		arrBypassConn := conn.ArrayBypass

		senderType, isArray, err := a.getSenderPortAddrType(
			arrBypassConn.SenderOutport,
			scope,
			compInterface,
			nodes,
			nodesIfaces,
		)
		if err != nil {
			return src.Connection{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &conn.Normal.SenderSide.Meta,
			}.Wrap(err)
		}
		if !isArray {
			return src.Connection{}, &compiler.Error{
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
			return src.Connection{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &conn.Normal.SenderSide.Meta,
			}.Wrap(err)
		}
		if !isArray {
			return src.Connection{}, &compiler.Error{
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
			return src.Connection{}, &compiler.Error{
				Err: fmt.Errorf(
					"Incompatible types: %v -> %v: %w",
					arrBypassConn.SenderOutport, arrBypassConn.ReceiverInport, err,
				),
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}

		nodesNetUsage(nodesUsage).AddOutport(
			arrBypassConn.SenderOutport.Node,
			arrBypassConn.SenderOutport.Port,
		)
		nodesNetUsage(nodesUsage).AddInport(
			arrBypassConn.ReceiverInport.Node,
			arrBypassConn.ReceiverInport.Port,
		)

		return conn, nil
	}

	// now handle normal connections
	normConn := conn.Normal

	// TODO mark portless connections as used

	resolvedSender, resolvedSenderType, isSenderArr, err := a.getSenderSideType(
		normConn.SenderSide,
		compInterface,
		nodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return src.Connection{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &normConn.SenderSide.Meta,
		}.Wrap(err)
	}

	if normConn.SenderSide.PortAddr != nil {
		// make sure only array outports has indexes
		if !isSenderArr && normConn.SenderSide.PortAddr.Idx != nil {
			return src.Connection{}, &compiler.Error{
				Err:      errors.New("Index for non-array port"),
				Meta:     &normConn.SenderSide.PortAddr.Meta,
				Location: &scope.Location,
			}
		}

		// make sure array outports always has indexes (it's not arr-bypass)
		if isSenderArr && normConn.SenderSide.PortAddr.Idx == nil {
			return src.Connection{}, &compiler.Error{
				Err:      errors.New("Index needed for array outport"),
				Meta:     &normConn.SenderSide.PortAddr.Meta,
				Location: &scope.Location,
			}
		}

		// mark node's outport as used since sender isn't const ref
		nodesNetUsage(nodesUsage).AddOutport(
			normConn.SenderSide.PortAddr.Node,
			normConn.SenderSide.PortAddr.Port,
		)
	}

	if len(normConn.SenderSide.Selectors) > 0 {
		lastFieldType, err := a.getStructFieldTypeByPath(
			resolvedSenderType,
			normConn.SenderSide.Selectors,
			scope,
		)
		if err != nil {
			return src.Connection{}, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}
		resolvedSenderType = lastFieldType
	}

	if len(normConn.ReceiverSide.DeferredConnections) == 0 && len(normConn.ReceiverSide.Receivers) == 0 {
		return src.Connection{}, &compiler.Error{
			Err: errors.New(
				"Connection's receiver side cannot be empty, it must either have deferred connection or receivers",
			),
			Location: &scope.Location,
			Meta:     &conn.Meta,
		}
	}

	resolvedDefConns := make([]src.Connection, 0, len(normConn.ReceiverSide.DeferredConnections))
	if len(normConn.ReceiverSide.DeferredConnections) != 0 {
		// note that we call analyzeConnections instead of analyzeComponentNetwork
		// because we only need to analyze connections and update nodesUsage
		// analyzeComponentNetwork OTOH will also validate nodesUsage by itself
		var err *compiler.Error
		resolvedDefConns, err = a.analyzeConnections( // indirect recursion
			normConn.ReceiverSide.DeferredConnections,
			compInterface,
			nodes,
			nodesIfaces,
			nodesUsage,
			scope,
		)
		if err != nil {
			return src.Connection{}, err
		}
		// receiver side can contain both deferred connections and receivers so we don't return yet
	}

	// TODO handle chain connection

	for _, receiver := range normConn.ReceiverSide.Receivers {
		inportTypeExpr, isReceiverArr, err := a.getReceiverType(
			receiver.PortAddr,
			compInterface,
			nodes,
			nodesIfaces,
			scope,
		)
		if err != nil {
			return src.Connection{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &receiver.Meta,
			}.Wrap(err)
		}

		// make sure only array outports has indexes
		if !isReceiverArr && receiver.PortAddr.Idx != nil {
			return src.Connection{}, &compiler.Error{
				Err:      errors.New("Index for non-array port"),
				Meta:     &receiver.PortAddr.Meta,
				Location: &scope.Location,
			}
		}

		// make sure array inports always has indexes (it's not arr-bypass)
		if isReceiverArr && receiver.PortAddr.Idx == nil {
			return src.Connection{}, &compiler.Error{
				Err:      errors.New("Index needed for array inport"),
				Meta:     &receiver.PortAddr.Meta,
				Location: &scope.Location,
			}
		}

		if err := a.resolver.IsSubtypeOf(resolvedSenderType, inportTypeExpr, scope); err != nil {
			return src.Connection{}, &compiler.Error{
				Err: fmt.Errorf(
					"Incompatible types: %v -> %v: %w",
					normConn.SenderSide, receiver, err,
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

	return src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: resolvedSender,
			ReceiverSide: src.ConnectionReceiverSide{
				DeferredConnections: resolvedDefConns,
				Receivers:           normConn.ReceiverSide.Receivers,
			},
		},
		Meta: conn.Meta,
	}, nil
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
			Meta:     &compInterface.Meta,
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
				// maybe it's portless connection
				if _, ok := nodeUsage.In[""]; ok && len(nodeIface.IO.In) == 1 {
					continue
				}

				meta := nodeIface.IO.In[inportName].Meta
				return &compiler.Error{
					Err:      fmt.Errorf("%w: %v:%v", ErrUnusedNodeInport, nodeName, inportName),
					Location: &scope.Location,
					Meta:     &meta,
				}
			}
		}

		if len(nodeIface.IO.Out) == 0 { // e.g. std/builtin.Del
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
			// maybe it's portless connection
			if _, ok := nodeUsage.Out[""]; ok && len(nodeIface.IO.Out) == 1 {
				continue
			}

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
	if portAddr.Port == "" {
		if len(ports) > 1 {
			return ts.Expr{}, false, &compiler.Error{
				Err:      ErrIllegalPortlessConnection,
				Location: &scope.Location,
				Meta:     &portAddr.Meta,
			}
		}

		for name := range ports {
			portAddr.Port = name
			break
		}
	}

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
) (src.ConnectionSenderSide, ts.Expr, bool, *compiler.Error) {
	if senderSide.PortAddr == nil && senderSide.Const == nil {
		return src.ConnectionSenderSide{}, ts.Expr{}, false, &compiler.Error{
			Err:      ErrSenderIsEmpty,
			Location: &scope.Location,
			Meta:     &senderSide.Meta,
		}
	}

	if senderSide.Const != nil {
		resolvedConst, resolvedExpr, err := a.getResolvedSenderConstType(*senderSide.Const, scope)
		if err != nil {
			return src.ConnectionSenderSide{}, ts.Expr{}, false, err
		}

		return src.ConnectionSenderSide{
			Const:     &resolvedConst,
			Selectors: senderSide.Selectors,
			Meta:      senderSide.Meta,
		}, resolvedExpr, false, nil
	}

	resolvedExpr, isArr, err := a.getSenderPortAddrType(
		*senderSide.PortAddr,
		scope,
		iface,
		nodes,
		nodesIfaces,
	)
	if err != nil {
		return src.ConnectionSenderSide{}, ts.Expr{}, false, err
	}

	return src.ConnectionSenderSide{
		PortAddr:  senderSide.PortAddr,
		Selectors: senderSide.Selectors,
		Meta:      senderSide.Meta,
	}, resolvedExpr, isArr, nil
}

// getSenderPortAddrType returns port's type and isArray bool
func (a Analyzer) getSenderPortAddrType(
	senderSidePortAddr src.PortAddr,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]src.Interface,
) (ts.Expr, bool, *compiler.Error) {
	if senderSidePortAddr.Node == "out" {
		return ts.Expr{}, false, &compiler.Error{
			Err:      ErrReadSelfOut,
			Location: &scope.Location,
			Meta:     &senderSidePortAddr.Meta,
		}
	}

	if senderSidePortAddr.Node == "in" {
		inports := iface.IO.In

		inport, ok := inports[senderSidePortAddr.Port]
		if !ok {
			return ts.Expr{}, false, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrInportNotFound, senderSidePortAddr.Port),
				Location: &scope.Location,
				Meta:     &senderSidePortAddr.Meta,
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
				Meta:     &senderSidePortAddr.Meta,
			}
		}

		return resolvedInportType, inport.IsArray, nil
	}

	return a.getNodeOutportType(
		senderSidePortAddr, nodes, nodesIfaces, scope,
	)
}

func (a Analyzer) getResolvedSenderConstType(
	constSender src.Const,
	scope src.Scope,
) (src.Const, ts.Expr, *compiler.Error) {
	if constSender.Ref != nil {
		expr, err := a.getResolvedConstTypeByRef(*constSender.Ref, scope)
		if err != nil {
			return src.Const{}, ts.Expr{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &constSender.Ref.Meta,
			}.Wrap(err)
		}
		return constSender, expr, nil
	}

	if constSender.Message == nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Err:      ErrLiteralSenderTypeEmpty,
			Location: &scope.Location,
			Meta:     &constSender.Meta,
		}
	}

	resolvedExpr, err := a.resolver.ResolveExpr(
		constSender.Message.TypeExpr,
		scope,
	)
	if err != nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &constSender.Message.Meta,
		}
	}

	if err := a.validateLiteralSender(resolvedExpr); err != nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &constSender.Message.Meta,
		}
	}

	return src.Const{
		Message: &src.Message{
			TypeExpr:    resolvedExpr,
			Bool:        constSender.Message.Bool,
			Int:         constSender.Message.Int,
			Float:       constSender.Message.Float,
			Str:         constSender.Message.Str,
			List:        constSender.Message.List,
			MapOrStruct: constSender.Message.MapOrStruct,
			Enum:        constSender.Message.Enum,
			Meta:        constSender.Message.Meta,
		},
		Meta: constSender.Meta,
	}, resolvedExpr, nil
}

func (a Analyzer) validateLiteralSender(resolvedExpr ts.Expr) error {
	if resolvedExpr.Inst != nil {
		switch resolvedExpr.Inst.Ref.String() {
		case "bool", "int", "float", "string":
			return nil
		}
		return ErrComplexLiteralSender
	}

	if resolvedExpr.Lit == nil ||
		resolvedExpr.Lit.Enum == nil {
		return ErrComplexLiteralSender
	}

	return nil
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

func (a Analyzer) getResolvedConstTypeByRef(ref core.EntityRef, scope src.Scope) (ts.Expr, *compiler.Error) {
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

	resolvedExpr, err := a.resolver.ResolveExpr(entity.Const.Message.TypeExpr, scope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &entity.Const.Message.Meta,
		}
	}

	return resolvedExpr, nil
}

func (a Analyzer) getStructFieldTypeByPath(
	senderType ts.Expr,
	path []string,
	scope src.Scope,
) (ts.Expr, error) {
	if len(path) == 0 {
		return senderType, nil
	}

	if senderType.Lit == nil || senderType.Lit.Struct == nil {
		return ts.Expr{}, fmt.Errorf("Type not struct: %v", senderType.String())
	}

	curField := path[0]
	fieldType, ok := senderType.Lit.Struct[curField]
	if !ok {
		return ts.Expr{}, fmt.Errorf("struct field '%v' not found", curField)
	}

	return a.getStructFieldTypeByPath(fieldType, path[1:], scope)
}

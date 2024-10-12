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
	ErrUnusedOutports                = errors.New("All flow's outports are unused")
	ErrUnusedOutport                 = errors.New("Unused outport found")
	ErrUnusedInports                 = errors.New("All flow inports are unused")
	ErrUnusedInport                  = errors.New("Unused inport found")
	ErrLiteralSenderTypeEmpty        = errors.New("Literal network sender must contain message value")
	ErrComplexLiteralSender          = errors.New("Literal network sender must have primitive type")
	ErrIllegalPortlessConnection     = errors.New("Connection to a node, with more than one port, must always has a port name")
	ErrGuardMixedWithExplicitErrConn = errors.New("If node has error guard '?' it's ':err' outport must not be explicitly used in the network")
)

// analyzeNetwork must be called after analyzeNodes so we sure nodes are resolved.
func (a Analyzer) analyzeNetwork(
	net []src.Connection,
	compInterface src.Interface,
	hasGuard bool,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
) ([]src.Connection, *compiler.Error) {
	// we create it here because there's recursion down there
	nodesUsage := make(map[string]nodeNetUsage, len(nodes))

	analyzedConnections, err := a.analyzeConnections(net, compInterface, nodes, nodesIfaces, nodesUsage, scope)
	if err != nil {
		return nil, compiler.Error{Location: &scope.Location}.Wrap(err)
	}

	if err := a.checkNetPortsUsage(
		compInterface,
		nodesIfaces,
		hasGuard,
		scope,
		nodesUsage,
		nodes,
	); err != nil {
		return nil, compiler.Error{Location: &scope.Location}.Wrap(err)
	}

	return analyzedConnections, nil
}

// analyzeConnections does two things:
// 1. Analyzes every connection and terminates with non-nil error if any of them is invalid.
// 2. Updates nodesUsage (we mutate it in-place instead of returning to avoid merging across recursive calls).
func (a Analyzer) analyzeConnections(
	net []src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]nodeNetUsage,
	scope src.Scope,
) ([]src.Connection, *compiler.Error) {
	analyzedConnections := make([]src.Connection, 0, len(net))

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

		analyzedConnections = append(analyzedConnections, resolvedConn)
	}

	return analyzedConnections, nil
}

func (a Analyzer) analyzeConnection(
	conn src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
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
				Meta:     &conn.Meta,
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
	senderSide := conn.Normal.SenderSide

	if senderSide.PortAddr == nil && senderSide.Const == nil && senderSide.Range == nil {
		return src.Connection{}, &compiler.Error{
			Err:      ErrEmptySender,
			Location: &scope.Location,
			Meta:     &senderSide.Meta,
		}
	}

	resolvedSenderSide, resolvedSenderType, isSenderArr, err := a.getSenderSideType(
		senderSide,
		compInterface,
		nodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return src.Connection{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &senderSide.Meta,
		}.Wrap(err)
	}

	if senderSide.PortAddr != nil {
		// make sure only array outports has indexes
		if !isSenderArr && senderSide.PortAddr.Idx != nil {
			return src.Connection{}, &compiler.Error{
				Err:      errors.New("Index for non-array port"),
				Meta:     &senderSide.PortAddr.Meta,
				Location: &scope.Location,
			}
		}

		// make sure array outports always has indexes (it's not arr-bypass)
		if isSenderArr && senderSide.PortAddr.Idx == nil {
			return src.Connection{}, &compiler.Error{
				Err:      errors.New("Index needed for array outport"),
				Meta:     &senderSide.PortAddr.Meta,
				Location: &scope.Location,
			}
		}

		// mark node's outport as used since sender isn't const ref
		nodesNetUsage(nodesUsage).AddOutport(
			senderSide.PortAddr.Node,
			senderSide.PortAddr.Port,
		)
	}

	if len(senderSide.Selectors) > 0 {
		lastFieldType, err := a.getStructFieldTypeByPath(
			resolvedSenderType,
			senderSide.Selectors,
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

	receiverSide := conn.Normal.ReceiverSide

	if len(receiverSide.Receivers) == 0 &&
		len(receiverSide.DeferredConnections) == 0 &&
		receiverSide.ChainedConnection == nil {
		return src.Connection{}, &compiler.Error{
			Err:      errors.New("Connection must have receiver-side"),
			Location: &scope.Location,
			Meta:     &senderSide.Meta,
		}
	}

	analyzedDeferredConnections := make([]src.Connection, 0, len(receiverSide.DeferredConnections))
	if len(receiverSide.DeferredConnections) != 0 {
		// note that we call analyzeConnections instead of analyzeNetwork
		// because we only need to analyze connections and update nodesUsage
		// analyzeNetwork OTOH will also validate nodesUsage by itself
		var err *compiler.Error
		analyzedDeferredConnections, err = a.analyzeConnections( // indirect recursion
			receiverSide.DeferredConnections,
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

	if receiverSide.ChainedConnection != nil {
		// analyze the sender side compatibility with the chained connection's receiver
		chainedConn := receiverSide.ChainedConnection
		if chainedConn.Normal == nil {
			return src.Connection{}, &compiler.Error{
				Err:      errors.New("Chained connection must be a normal connection"),
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}

		chainHead := chainedConn.Normal.SenderSide
		chainHeadType, err := a.getChainHeadType(
			chainHead,
			nodes,
			nodesIfaces,
			scope,
		)
		if err != nil {
			return src.Connection{}, err
		}

		if err := a.resolver.IsSubtypeOf(resolvedSenderType, chainHeadType, scope); err != nil {
			return src.Connection{}, &compiler.Error{
				Err: fmt.Errorf(
					"Incompatible types: %v -> %v: %w",
					senderSide, chainHead, err,
				),
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}
		}

		// Recursively analyze the chained connection
		analyzedChainedConn, err := a.analyzeConnection(
			*chainedConn,
			compInterface,
			nodes,
			nodesIfaces,
			scope,
			nodesUsage,
		)
		if err != nil {
			return src.Connection{}, err
		}

		// Update the chained connection with the analyzed version
		receiverSide.ChainedConnection = &analyzedChainedConn

		if chainHead.PortAddr != nil {
			// don't forget to mark the receiver side as used
			// (in recursive call it's marked as sender, but it's also receiver)
			nodesNetUsage(nodesUsage).AddInport(
				chainHead.PortAddr.Node,
				chainHead.PortAddr.Port,
			)
		}
	}

	for _, receiver := range receiverSide.Receivers {
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
					senderSide, receiver, err,
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
			SenderSide: resolvedSenderSide,
			ReceiverSide: src.ConnectionReceiverSide{
				DeferredConnections: analyzedDeferredConnections,
				Receivers:           receiverSide.Receivers,
				ChainedConnection:   receiverSide.ChainedConnection,
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
// Every self inport and outport is used;
// Every sub-node's inport is used;
// For every sub-node's there's at least one used outport.
func (Analyzer) checkNetPortsUsage(
	compInterface src.Interface,
	nodesIfaces map[string]foundInterface,
	hasGuard bool,
	scope src.Scope,
	nodesUsage map[string]nodeNetUsage,
	nodes map[string]src.Node,
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
		// note that self outports are inports for the network
		if _, ok := outportsUsage.In[outportName]; ok {
			continue
		}

		if outportName == "err" && hasGuard {
			continue
		}

		return &compiler.Error{
			Err:      fmt.Errorf("%w '%v'", ErrUnusedOutport, outportName),
			Location: &scope.Location,
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

		for inportName := range nodeIface.iface.IO.In {
			if _, ok := nodeUsage.In[inportName]; !ok {
				// maybe it's portless connection
				if _, ok := nodeUsage.In[""]; ok && len(nodeIface.iface.IO.In) == 1 {
					continue
				}

				meta := nodeIface.iface.IO.In[inportName].Meta

				return &compiler.Error{
					Err: fmt.Errorf(
						"Unused node inport: %v:%v",
						nodeName,
						inportName,
					),
					Location: &scope.Location,
					Meta:     &meta,
				}
			}
		}

		if len(nodeIface.iface.IO.Out) == 0 { // e.g. std/builtin.Del
			continue
		}

		atLeastOneOutportIsUsed := false
		for outportName := range nodeIface.iface.IO.Out {
			if _, ok := nodeUsage.Out[outportName]; !ok {
				continue
			}

			atLeastOneOutportIsUsed = true

			if outportName != "err" {
				continue
			}

			meta := nodes[nodeName].Meta

			if nodes[nodeName].ErrGuard {
				return &compiler.Error{
					Err: fmt.Errorf(
						"%w: %v",
						ErrGuardMixedWithExplicitErrConn,
						nodeName,
					),
					Location: &scope.Location,
					Meta:     &meta,
				}
			}

		}

		if !atLeastOneOutportIsUsed {
			// maybe it's portless connection
			if _, ok := nodeUsage.Out[""]; ok && len(nodeIface.iface.IO.Out) == 1 {
				continue
			}

			return &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrUnusedNodeOutports, nodeName),
				Location: &scope.Location,
				Meta:     &nodeIface.iface.Meta,
			}
		}
	}

	// TODO make sure we don't use same sender twice or more outside of the fan-out expression
	
	return nil
}

func (a Analyzer) getReceiverType(
	receiverSide src.PortAddr,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
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
	nodesIfaces map[string]foundInterface,
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
		nodeIface.iface.IO.In,
		nodeIface.iface.TypeParams.Params,
		portAddr,
		node,
		scope.WithLocation(nodeIface.location),
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
				Err:      fmt.Errorf("%w: node '%v'", ErrIllegalPortlessConnection, portAddr.Node),
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

func (a Analyzer) getSenderSideType(
	senderSide src.ConnectionSenderSide,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
) (src.ConnectionSenderSide, ts.Expr, bool, *compiler.Error) {
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

	if senderSide.Range != nil {
		// range sends stream<int> from its :data outport
		rangeType := ts.Expr{
			Inst: &ts.InstExpr{
				Ref: core.EntityRef{Name: "stream"},
				Args: []ts.Expr{{
					Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
				}},
			},
		}
		return senderSide, rangeType, false, nil
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
	nodesIfaces map[string]foundInterface,
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
	if constSender.Value.Ref != nil {
		expr, err := a.getResolvedConstTypeByRef(*constSender.Value.Ref, scope)
		if err != nil {
			return src.Const{}, ts.Expr{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &constSender.Value.Ref.Meta,
			}.Wrap(err)
		}
		return constSender, expr, nil
	}

	if constSender.Value.Message == nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Err:      ErrLiteralSenderTypeEmpty,
			Location: &scope.Location,
			Meta:     &constSender.Meta,
		}
	}

	resolvedExpr, err := a.resolver.ResolveExpr(
		constSender.TypeExpr,
		scope,
	)
	if err != nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &constSender.Value.Message.Meta,
		}
	}

	if err := a.validateLiteralSender(resolvedExpr); err != nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &constSender.Value.Message.Meta,
		}
	}

	return src.Const{
		TypeExpr: resolvedExpr,
		Value: src.ConstValue{
			Message: &src.MsgLiteral{
				Bool:         constSender.Value.Message.Bool,
				Int:          constSender.Value.Message.Int,
				Float:        constSender.Value.Message.Float,
				Str:          constSender.Value.Message.Str,
				List:         constSender.Value.Message.List,
				DictOrStruct: constSender.Value.Message.DictOrStruct,
				Enum:         constSender.Value.Message.Enum,
				Meta:         constSender.Value.Message.Meta,
			},
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
	nodesIfaces map[string]foundInterface,
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
		nodeIface.iface.IO.Out,
		nodeIface.iface.TypeParams.Params,
		portAddr,
		node,
		scope.WithLocation(nodeIface.location),
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

	if entity.Const.Value.Ref != nil {
		expr, err := a.getResolvedConstTypeByRef(*entity.Const.Value.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Location: &location,
				Meta:     &entity.Const.Meta,
			}.Wrap(err)
		}
		return expr, nil
	}

	scope = scope.WithLocation(location)

	resolvedExpr, err := a.resolver.ResolveExpr(entity.Const.TypeExpr, scope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &entity.Const.Value.Message.Meta,
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

func (a Analyzer) getChainHeadType(
	chainHead src.ConnectionSenderSide,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	if chainHead.PortAddr == nil && chainHead.Range == nil {
		return ts.Expr{}, &compiler.Error{
			Err:      errors.New("Chained connection must start with a port or range"),
			Location: &scope.Location,
			Meta:     &chainHead.Meta,
		}
	}

	if chainHead.PortAddr != nil {
		resolvedType, _, err := a.getNodeInportType(*chainHead.PortAddr, nodes, nodesIfaces, scope)
		if err != nil {
			return ts.Expr{}, err
		}
		return resolvedType, nil
	}

	// range receives any to its :sig inport
	return ts.Expr{
		Inst: &ts.InstExpr{
			Ref: core.EntityRef{Name: "any"},
		},
	}, nil
}

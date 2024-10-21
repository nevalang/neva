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
	ErrUnusedOutports                = errors.New("All component outports are unused")
	ErrUnusedOutport                 = errors.New("unused outport found")
	ErrUnusedInports                 = errors.New("all flow inports are unused")
	ErrUnusedInport                  = errors.New("unused inport found")
	ErrLiteralSenderTypeEmpty        = errors.New("literal network sender must contain message value")
	ErrComplexLiteralSender          = errors.New("literal network sender must have primitive type")
	ErrIllegalPortlessConnection     = errors.New("connection to a node, with more than one port, must always has a port name")
	ErrGuardMixedWithExplicitErrConn = errors.New("if node has error guard '?' it's ':err' outport must not be explicitly used in the network")
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
	nodesUsage := make(map[string]netNodeUsage, len(nodes))

	analyzedConnections, err := a.analyzeConnections(
		net,
		compInterface,
		nodes,
		nodesIfaces,
		nodesUsage,
		scope,
	)
	if err != nil {
		return nil, compiler.Error{Location: &scope.Location}.Wrap(err)
	}

	if err := a.analyzeNetPortsUsage(
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
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
	scope src.Scope,
) ([]src.Connection, *compiler.Error) {
	analyzedConnections := make([]src.Connection, 0, len(net))

	for _, conn := range net {
		resolvedConn, err := a.analyzeConnection(
			conn,
			iface,
			nodes,
			nodesIfaces,
			scope,
			nodesUsage,
			nil,
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
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
	nodesUsage map[string]netNodeUsage,
	prevChainLink []src.ConnectionSender,
) (src.Connection, *compiler.Error) {
	if conn.ArrayBypass != nil {
		if err := a.analyzeArrayBypassConnection(
			conn,
			scope,
			iface,
			nodes,
			nodesIfaces,
			nodesUsage,
		); err != nil {
			return src.Connection{}, err
		}
		return conn, nil
	}

	analyzedNormalConn, err := a.analyzeNormalConnection(
		conn.Normal,
		iface,
		nodes,
		nodesIfaces,
		scope,
		nodesUsage,
		prevChainLink,
	)
	if err != nil {
		return src.Connection{}, err
	}

	return src.Connection{
		Normal: analyzedNormalConn,
		Meta:   conn.Meta,
	}, nil
}

func (a Analyzer) analyzeNormalConnection(
	normConn *src.NormalConnection,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
	nodesUsage map[string]netNodeUsage,
	prevChainLink []src.ConnectionSender,
) (*src.NormalConnection, *compiler.Error) {
	analyzedSenders, resolvedSenderTypes, err := a.analyzeSenderSide(
		normConn.SenderSide,
		scope,
		iface,
		nodes,
		nodesIfaces,
		nodesUsage,
		prevChainLink,
	)
	if err != nil {
		return nil, err
	}

	analyzedReceiverSide, err := a.analyzeReceiverSide(
		normConn.ReceiverSide,
		scope,
		iface,
		nodes,
		nodesIfaces,
		nodesUsage,
		resolvedSenderTypes,
		analyzedSenders,
	)
	if err != nil {
		return nil, err
	}

	return &src.NormalConnection{
		SenderSide:   analyzedSenders,
		ReceiverSide: analyzedReceiverSide,
	}, nil
}

func (a Analyzer) analyzeReceiverSide(
	receiverSide []src.ConnectionReceiver,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
	resolvedSenderTypes []*ts.Expr,
	analyzedSenders []src.ConnectionSender,
) ([]src.ConnectionReceiver, *compiler.Error) {
	analyzedReceivers := make([]src.ConnectionReceiver, 0, len(receiverSide))

	for _, receiver := range receiverSide {
		analyzedReceiver, err := a.analyzeReceiver(
			receiver,
			scope,
			iface,
			nodes,
			nodesIfaces,
			nodesUsage,
			resolvedSenderTypes,
			analyzedSenders,
		)
		if err != nil {
			return nil, err
		}

		analyzedReceivers = append(analyzedReceivers, *analyzedReceiver)
	}

	return analyzedReceivers, nil
}

func (a Analyzer) analyzeReceiver(
	receiver src.ConnectionReceiver,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
	resolvedSenderTypes []*ts.Expr,
	analyzedSenders []src.ConnectionSender,
) (*src.ConnectionReceiver, *compiler.Error) {
	if receiver.PortAddr == nil &&
		receiver.ChainedConnection == nil &&
		receiver.DeferredConnection == nil {
		return nil, &compiler.Error{
			Err:      errors.New("Connection must have receiver-side"),
			Location: &scope.Location,
			Range:    &receiver.Meta,
		}
	}

	switch {
	case receiver.PortAddr != nil:
		analyzedPortAddr, err := a.analyzePortAddrReceiver(
			*receiver.PortAddr,
			scope,
			iface,
			nodes,
			nodesIfaces,
			nodesUsage,
			resolvedSenderTypes,
			analyzedSenders,
		)
		if err != nil {
			return nil, err
		}
		return &src.ConnectionReceiver{
			PortAddr: &analyzedPortAddr,
		}, nil
	case receiver.ChainedConnection != nil:
		analyzedChainedConn, err := a.analyzeChainedConnectionReceiver(
			*receiver.ChainedConnection,
			scope,
			iface,
			nodes,
			nodesIfaces,
			nodesUsage,
			resolvedSenderTypes,
			analyzedSenders,
		)
		if err != nil {
			return nil, err
		}
		return &src.ConnectionReceiver{
			ChainedConnection: &analyzedChainedConn,
		}, nil
	case receiver.DeferredConnection != nil:
		analyzedDeferredConn, err := a.analyzeConnection(
			*receiver.DeferredConnection,
			iface,
			nodes,
			nodesIfaces,
			scope,
			nodesUsage,
			nil,
		)
		if err != nil {
			return nil, err
		}
		return &src.ConnectionReceiver{
			DeferredConnection: &analyzedDeferredConn,
		}, nil
	}

	return nil, &compiler.Error{
		Err:      errors.New("Connection must have receiver-side"),
		Location: &scope.Location,
		Range:    &receiver.Meta,
	}
}

func (a Analyzer) analyzePortAddrReceiver(
	portAddr src.PortAddr,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
	resolvedSenderTypes []*ts.Expr,
	analyzedSenders []src.ConnectionSender,
) (src.PortAddr, *compiler.Error) {
	typeExpr, isArrPort, err := a.getReceiverPortType(
		portAddr,
		iface,
		nodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return src.PortAddr{}, compiler.Error{
			Location: &scope.Location,
			Range:    &portAddr.Meta,
		}.Wrap(err)
	}

	if !isArrPort && portAddr.Idx != nil {
		return src.PortAddr{}, &compiler.Error{
			Err:      errors.New("Index for non-array port"),
			Range:    &portAddr.Meta,
			Location: &scope.Location,
		}
	}

	if isArrPort && portAddr.Idx == nil {
		return src.PortAddr{}, &compiler.Error{
			Err:      errors.New("Index needed for array inport"),
			Range:    &portAddr.Meta,
			Location: &scope.Location,
		}
	}

	for i, resolvedSenderType := range resolvedSenderTypes {
		if err := a.resolver.IsSubtypeOf(*resolvedSenderType, typeExpr, scope); err != nil {
			return src.PortAddr{}, &compiler.Error{
				Err: fmt.Errorf(
					"Incompatible types: %v -> %v: %w",
					analyzedSenders[i], portAddr, err,
				),
				Location: &scope.Location,
				Range:    &portAddr.Meta,
			}
		}
	}

	if err := netNodesUsage(nodesUsage).trackInportUsage(portAddr); err != nil {
		return src.PortAddr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Range:    &portAddr.Meta,
		}
	}

	return portAddr, nil
}

func (a Analyzer) analyzeChainedConnectionReceiver(
	chainedConn src.Connection,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
	resolvedSenderTypes []*ts.Expr,
	analyzedSenders []src.ConnectionSender,
) (src.Connection, *compiler.Error) {
	if chainedConn.Normal == nil {
		return src.Connection{}, &compiler.Error{
			Err:      errors.New("chained connection must be a normal connection"),
			Location: &scope.Location,
			Range:    &chainedConn.Meta,
		}
	}

	if len(chainedConn.Normal.SenderSide) != 1 {
		return src.Connection{}, &compiler.Error{
			Err:      errors.New("multiple senders are only allowed at the start of a connection"),
			Location: &scope.Location,
			Range:    &chainedConn.Normal.Meta,
		}
	}

	chainHead := chainedConn.Normal.SenderSide[0]

	chainHeadType, err := a.getChainHeadType(
		chainHead,
		nodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return src.Connection{}, err
	}

	for i, resolvedSenderType := range resolvedSenderTypes {
		if err := a.resolver.IsSubtypeOf(*resolvedSenderType, chainHeadType, scope); err != nil {
			return src.Connection{}, &compiler.Error{
				Err: fmt.Errorf(
					"Incompatible types: %v -> %v: %w",
					analyzedSenders[i], chainHead, err,
				),
				Location: &scope.Location,
				Range:    &chainedConn.Meta,
			}
		}
	}

	analyzedChainedConn, err := a.analyzeConnection(
		chainedConn,
		iface,
		nodes,
		nodesIfaces,
		scope,
		nodesUsage,
		analyzedSenders,
	)
	if err != nil {
		return src.Connection{}, err
	}

	if chainHead.PortAddr != nil {
		if err := netNodesUsage(nodesUsage).trackInportUsage(*chainHead.PortAddr); err != nil {
			return src.Connection{}, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Range:    &chainedConn.Meta,
			}
		}
	}

	return analyzedChainedConn, nil
}

func (a Analyzer) analyzeSenderSide(
	senders []src.ConnectionSender,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
	prevChainLink []src.ConnectionSender,
) ([]src.ConnectionSender, []*ts.Expr, *compiler.Error) {
	analyzedSenders := make([]src.ConnectionSender, 0, len(senders))
	resolvedSenderTypes := make([]*ts.Expr, 0, len(senders))

	for _, sender := range senders {
		analyzedSender, expr, err := a.analyzeSender(
			sender,
			scope,
			iface,
			nodes,
			nodesIfaces,
			nodesUsage,
			prevChainLink,
		)
		if err != nil {
			return nil, nil, compiler.Error{
				Location: &scope.Location,
				Range:    &sender.Meta,
			}.Wrap(err)
		}

		analyzedSenders = append(analyzedSenders, *analyzedSender)
		resolvedSenderTypes = append(resolvedSenderTypes, expr)

		if sender.PortAddr == nil {
			continue
		}
	}

	return analyzedSenders, resolvedSenderTypes, nil
}

func (a Analyzer) analyzeSender(
	sender src.ConnectionSender,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
	prevChainLink []src.ConnectionSender,
) (*src.ConnectionSender, *ts.Expr, *compiler.Error) {
	if sender.PortAddr == nil &&
		sender.Const == nil &&
		sender.Range == nil &&
		len(sender.StructSelector) == 0 {
		return nil, nil, &compiler.Error{
			Err:      ErrEmptySender,
			Location: &scope.Location,
			Range:    &sender.Meta,
		}
	}

	if sender.Range != nil && prevChainLink == nil {
		return nil, nil, &compiler.Error{
			Err:      errors.New("range expression cannot be used in non-chained connection"),
			Location: &scope.Location,
			Range:    &sender.Meta,
		}
	}

	if sender.Const != nil && len(prevChainLink) != 0 {
		return nil, nil, &compiler.Error{
			Err:      errors.New("constant cannot be used in chained connection"),
			Location: &scope.Location,
			Range:    &sender.Meta,
		}
	}

	if len(sender.StructSelector) > 0 && prevChainLink == nil {
		return nil, nil, &compiler.Error{
			Err:      errors.New("struct selectors cannot be used in non-chained connection"),
			Location: &scope.Location,
			Range:    &sender.Meta,
		}
	}

	resolvedSender, resolvedSenderType, isSenderArr, err := a.getSenderSideType(
		sender,
		iface,
		nodes,
		nodesIfaces,
		scope,
		prevChainLink,
	)
	if err != nil {
		return nil, nil, compiler.Error{
			Location: &scope.Location,
			Range:    &sender.Meta,
		}.Wrap(err)
	}

	if sender.PortAddr != nil {
		if !isSenderArr && sender.PortAddr.Idx != nil {
			return nil, nil, &compiler.Error{
				Err:      errors.New("Index for non-array port"),
				Range:    &sender.PortAddr.Meta,
				Location: &scope.Location,
			}
		}

		if isSenderArr && sender.PortAddr.Idx == nil {
			return nil, nil, &compiler.Error{
				Err:      errors.New("Index needed for array outport"),
				Range:    &sender.PortAddr.Meta,
				Location: &scope.Location,
			}
		}

		if err := netNodesUsage(nodesUsage).trackOutportUsage(*sender.PortAddr); err != nil {
			return nil, nil, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Range:    &sender.PortAddr.Meta,
			}
		}
	}

	return &resolvedSender, &resolvedSenderType, nil
}

func (a Analyzer) analyzeArrayBypassConnection(
	conn src.Connection,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
) *compiler.Error {
	arrBypassConn := conn.ArrayBypass

	senderType, isArray, err := a.getPortSenderType(
		arrBypassConn.SenderOutport,
		scope,
		iface,
		nodes,
		nodesIfaces,
	)
	if err != nil {
		return compiler.Error{
			Location: &scope.Location,
			Range:    &conn.Meta,
		}.Wrap(err)
	}
	if !isArray {
		return &compiler.Error{
			Err:      errors.New("Non-array outport in array-bypass connection"),
			Location: &scope.Location,
			Range:    &arrBypassConn.SenderOutport.Meta,
		}
	}

	receiverType, isArray, err := a.getReceiverPortType(
		arrBypassConn.ReceiverInport,
		iface,
		nodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return compiler.Error{
			Location: &scope.Location,
			Range:    &conn.Meta,
		}.Wrap(err)
	}
	if !isArray {
		return &compiler.Error{
			Err:      errors.New("Non-array outport in array-bypass connection"),
			Location: &scope.Location,
			Range:    &arrBypassConn.SenderOutport.Meta,
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
			Range:    &conn.Meta,
		}
	}

	if err := netNodesUsage(nodesUsage).trackOutportUsage(arrBypassConn.SenderOutport); err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Range:    &conn.Meta,
		}
	}

	if err := netNodesUsage(nodesUsage).trackInportUsage(arrBypassConn.ReceiverInport); err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Range:    &conn.Meta,
		}
	}

	return nil
}

func (a Analyzer) analyzeNetPortsUsage(
	compInterface src.Interface,
	nodesIfaces map[string]foundInterface,
	hasGuard bool,
	scope src.Scope,
	nodesUsage map[string]netNodeUsage,
	nodes map[string]src.Node,
) *compiler.Error {
	// 1. every self inport must be used
	inportsUsage, ok := nodesUsage["in"]
	if !ok {
		return &compiler.Error{
			Err:      ErrUnusedInports,
			Location: &scope.Location,
			Range:    &compInterface.Meta,
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

	// 2. every self-outport must be used
	outportsUsage, ok := nodesUsage["out"]
	if !ok {
		return &compiler.Error{
			Err:      ErrUnusedOutports,
			Location: &scope.Location,
			Range:    &compInterface.Meta,
		}
	}

	for outportName := range compInterface.IO.Out {
		if _, ok := outportsUsage.In[outportName]; ok { // self outports are inports in network
			continue
		}

		// err outport is allowed to be unused if parent uses guard
		if outportName == "err" && hasGuard {
			continue
		}

		return &compiler.Error{
			Err:      fmt.Errorf("%w '%v'", ErrUnusedOutport, outportName),
			Location: &scope.Location,
		}
	}

	// 3. check sub-nodes usage in network
	for nodeName, nodeIface := range nodesIfaces {
		// every sub-node must be used
		nodeUsage, ok := nodesUsage[nodeName]
		if !ok {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrUnusedNode, nodeName),
				Location: &scope.Location,
			}
		}

		// every sub-node's inport must be used
		for inportName := range nodeIface.iface.IO.In {
			if _, ok := nodeUsage.In[inportName]; ok {
				continue
			}

			_, portless := nodeUsage.In[""]
			if portless && len(nodeIface.iface.IO.In) == 1 {
				continue
			}

			return &compiler.Error{
				Err: fmt.Errorf(
					"Unused node inport: %v:%v",
					nodeName,
					inportName,
				),
				Location: &scope.Location,
				Range:    compiler.Pointer(nodeIface.iface.IO.In[inportName].Meta),
			}
		}

		if len(nodeIface.iface.IO.Out) == 0 { // e.g. Del
			continue
		}

		// :err outport must always be used + at least one outport must be used in general
		atLeastOneOutportIsUsed := false
		for outportName, port := range nodeIface.iface.IO.Out {
			_, portUsed := nodeUsage.Out[outportName]

			if portUsed {
				if outportName == "err" && nodes[nodeName].ErrGuard {
					return &compiler.Error{
						Err: fmt.Errorf(
							"if node has error guard '?' it's ':err' outport must not be explicitly used in the network: %v",
							nodeName,
						),
						Location: &scope.Location,
						Range:    &port.Meta,
					}
				}

				atLeastOneOutportIsUsed = true
				continue
			}

			if outportName == "err" && !nodes[nodeName].ErrGuard {
				return &compiler.Error{
					Err:      fmt.Errorf("unhandled error: %v:err", nodeName),
					Location: &scope.Location,
					Range:    &port.Meta,
				}
			}
		}

		if !atLeastOneOutportIsUsed {
			if _, ok := nodeUsage.Out[""]; ok && len(nodeIface.iface.IO.Out) == 1 {
				continue
			}

			return &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrUnusedNodeOutports, nodeName),
				Location: &scope.Location,
				Range:    &nodeIface.iface.Meta,
			}
		}
	}

	// 4. check that array ports are used correctly (from 0 and without holes)
	for nodeName, nodeUsage := range nodesUsage {
		for portName, usedSlots := range nodeUsage.In {
			if usedSlots == nil {
				continue // skip non-array ports
			}

			maxSlot := uint8(0)
			for slot := range usedSlots {
				if slot > maxSlot {
					maxSlot = slot
				}
			}

			for i := uint8(0); i <= maxSlot; i++ {
				if _, ok := usedSlots[i]; !ok {
					return &compiler.Error{
						Err: fmt.Errorf(
							"array inport '%s:%s' is used incorrectly: slot %d is missing",
							nodeName,
							portName,
							i,
						),
						Location: &scope.Location,
					}
				}
			}
		}

		for portName, usedSlots := range nodeUsage.Out {
			if usedSlots == nil {
				continue // skip non-array ports
			}

			maxSlot := uint8(0)
			for slot := range usedSlots {
				if slot > maxSlot {
					maxSlot = slot
				}
			}

			for i := uint8(0); i <= maxSlot; i++ {
				if _, ok := usedSlots[i]; !ok {
					return &compiler.Error{
						Err: fmt.Errorf(
							"array outport '%s:%s' is used incorrectly: slot %d is missing",
							nodeName,
							portName,
							i,
						),
						Location: &scope.Location,
					}
				}
			}
		}
	}

	return nil
}

func (a Analyzer) getReceiverPortType(
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
			Range:    &receiverSide.Meta,
		}
	}

	if receiverSide.Node == "out" {
		outports := iface.IO.Out

		outport, ok := outports[receiverSide.Port]
		if !ok {
			return ts.Expr{}, false, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrOutportNotFound, receiverSide.Port),
				Location: &scope.Location,
				Range:    &receiverSide.Meta,
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
				Range:    &receiverSide.Meta,
			}
		}

		return resolvedOutportType, outport.IsArray, nil
	}

	nodeInportType, isArray, err := a.getNodeInportType(receiverSide, nodes, nodesIfaces, scope)
	if err != nil {
		return ts.Expr{}, false, compiler.Error{
			Location: &scope.Location,
			Range:    &receiverSide.Meta,
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
			Range:    &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, false, &compiler.Error{
			Err:      fmt.Errorf("%w '%v'", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Range:    &portAddr.Meta,
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
			Range:    &portAddr.Meta,
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
				Range:    &portAddr.Meta,
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
			Range:    &portAddr.Meta,
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
			Range:    &port.Meta,
		}
	}

	return resolvedPortType, port.IsArray, nil
}

func (a Analyzer) getSenderSideType(
	senderSide src.ConnectionSender,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
	prevChainLink []src.ConnectionSender,
) (src.ConnectionSender, ts.Expr, bool, *compiler.Error) {
	if senderSide.Const != nil {
		resolvedConst, resolvedExpr, err := a.getConstSenderType(*senderSide.Const, scope)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, err
		}

		return src.ConnectionSender{
			Const: &resolvedConst,
			Meta:  senderSide.Meta,
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

	if len(senderSide.StructSelector) > 0 {
		if len(prevChainLink) != 1 {
			return src.ConnectionSender{}, ts.Expr{}, false, &compiler.Error{
				Err:      errors.New("fan-in with struct selectors is not supported"),
				Location: &scope.Location,
				Range:    &senderSide.Meta,
			}
		}

		_, chainLinkType, _, err := a.getSenderSideType(
			prevChainLink[0],
			iface,
			nodes,
			nodesIfaces,
			scope,
			prevChainLink,
		)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, err
		}

		lastFieldType, err := a.getSelectorsSenderType(
			chainLinkType,
			senderSide.StructSelector,
			scope,
		)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Range:    &senderSide.Meta,
			}
		}

		return senderSide, lastFieldType, false, nil
	}

	resolvedExpr, isArr, err := a.getPortSenderType(
		*senderSide.PortAddr,
		scope,
		iface,
		nodes,
		nodesIfaces,
	)
	if err != nil {
		return src.ConnectionSender{}, ts.Expr{}, false, err
	}

	return senderSide, resolvedExpr, isArr, nil
}

// getPortSenderType returns port's type and isArray bool
func (a Analyzer) getPortSenderType(
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
			Range:    &senderSidePortAddr.Meta,
		}
	}

	if senderSidePortAddr.Node == "in" {
		inports := iface.IO.In

		inport, ok := inports[senderSidePortAddr.Port]
		if !ok {
			return ts.Expr{}, false, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrInportNotFound, senderSidePortAddr.Port),
				Location: &scope.Location,
				Range:    &senderSidePortAddr.Meta,
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
				Range:    &senderSidePortAddr.Meta,
			}
		}

		return resolvedInportType, inport.IsArray, nil
	}

	return a.getNodeOutportType(
		senderSidePortAddr, nodes, nodesIfaces, scope,
	)
}

func (a Analyzer) getConstSenderType(
	constSender src.Const,
	scope src.Scope,
) (src.Const, ts.Expr, *compiler.Error) {
	if constSender.Value.Ref != nil {
		expr, err := a.getResolvedConstTypeByRef(*constSender.Value.Ref, scope)
		if err != nil {
			return src.Const{}, ts.Expr{}, compiler.Error{
				Location: &scope.Location,
				Range:    &constSender.Value.Ref.Meta,
			}.Wrap(err)
		}
		return constSender, expr, nil
	}

	if constSender.Value.Message == nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Err:      ErrLiteralSenderTypeEmpty,
			Location: &scope.Location,
			Range:    &constSender.Meta,
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
			Range:    &constSender.Value.Message.Meta,
		}
	}

	if err := a.validateLiteralSender(resolvedExpr); err != nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Range:    &constSender.Value.Message.Meta,
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
			Range:    &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, false, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrNodeNotFound, portAddr.Node),
			Location: &scope.Location,
			Range:    &portAddr.Meta,
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
			Range:    &ref.Meta,
		}
	}

	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", errors.New("Entity found but is not constant"), entity.Kind),
			Location: &location,
			Range:    entity.Meta(),
		}
	}

	if entity.Const.Value.Ref != nil {
		expr, err := a.getResolvedConstTypeByRef(*entity.Const.Value.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Location: &location,
				Range:    &entity.Const.Meta,
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
			Range:    &entity.Const.Value.Message.Meta,
		}
	}

	return resolvedExpr, nil
}

func (a Analyzer) getSelectorsSenderType(
	senderType ts.Expr,
	selectors []string,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	if len(selectors) == 0 {
		return senderType, nil
	}

	if senderType.Lit == nil || senderType.Lit.Struct == nil {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("Type not struct: %v", senderType.String()),
			Location: &scope.Location,
		}
	}

	curField := selectors[0]
	fieldType, ok := senderType.Lit.Struct[curField]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("struct field '%v' not found", curField),
			Location: &scope.Location,
		}
	}

	return a.getSelectorsSenderType(fieldType, selectors[1:], scope)
}

func (a Analyzer) getChainHeadType(
	chainHead src.ConnectionSender,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	if chainHead.PortAddr != nil {
		resolvedType, _, err := a.getNodeInportType(*chainHead.PortAddr, nodes, nodesIfaces, scope)
		if err != nil {
			return ts.Expr{}, err
		}
		return resolvedType, nil
	}

	if chainHead.Range != nil {
		return ts.Expr{
			Inst: &ts.InstExpr{
				Ref: core.EntityRef{Name: "any"}, // :sig
			},
		}, nil
	}

	if len(chainHead.StructSelector) > 0 {
		return a.getStructSelectorInportType(chainHead), nil
	}

	return ts.Expr{}, &compiler.Error{
		Err:      errors.New("Chained connection must start with port address or range expression"),
		Location: &scope.Location,
		Range:    &chainHead.Meta,
	}
}

func (Analyzer) getStructSelectorInportType(chainHead src.ConnectionSender) ts.Expr {
	// build nested struct type for selectors
	typeExpr := ts.Expr{
		Lit: &ts.LitExpr{
			Struct: make(map[string]ts.Expr),
		},
	}

	currentStruct := typeExpr.Lit.Struct
	for i, selector := range chainHead.StructSelector {
		if i == len(chainHead.StructSelector)-1 {
			// last selector, use any type
			currentStruct[selector] = ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "any"},
				},
			}
		} else {
			// create nested struct
			nestedStruct := make(map[string]ts.Expr)
			currentStruct[selector] = ts.Expr{
				Lit: &ts.LitExpr{
					Struct: nestedStruct,
				},
			}
			currentStruct = nestedStruct
		}
	}

	return typeExpr
}

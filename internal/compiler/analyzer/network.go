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
	ErrComplexLiteralSender = errors.New("literal network sender must have primitive type")
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
			Message:  "Connection must have receiver-side",
			Location: &scope.Location,
			Meta:     &receiver.Meta,
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
		Message:  "Connection must have receiver-side",
		Location: &scope.Location,
		Meta:     &receiver.Meta,
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
			Meta:     &portAddr.Meta,
		}.Wrap(err)
	}

	if !isArrPort && portAddr.Idx != nil {
		return src.PortAddr{}, &compiler.Error{
			Message:  "Index for non-array port",
			Meta:     &portAddr.Meta,
			Location: &scope.Location,
		}
	}

	if isArrPort && portAddr.Idx == nil {
		return src.PortAddr{}, &compiler.Error{
			Message:  "Index needed for array inport",
			Meta:     &portAddr.Meta,
			Location: &scope.Location,
		}
	}

	for i, resolvedSenderType := range resolvedSenderTypes {
		if err := a.resolver.IsSubtypeOf(*resolvedSenderType, typeExpr, scope); err != nil {
			return src.PortAddr{}, &compiler.Error{
				Message: fmt.Sprintf(
					"Incompatible types: %v -> %v: %v",
					analyzedSenders[i], portAddr, err.Error(),
				),
				Location: &scope.Location,
				Meta:     &portAddr.Meta,
			}
		}
	}

	if err := netNodesUsage(nodesUsage).trackInportUsage(portAddr); err != nil {
		return src.PortAddr{}, &compiler.Error{
			Message:  err.Error(),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
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
			Message:  "chained connection must be a normal connection",
			Location: &scope.Location,
			Meta:     &chainedConn.Meta,
		}
	}

	if len(chainedConn.Normal.SenderSide) != 1 {
		return src.Connection{}, &compiler.Error{
			Message:  "multiple senders are only allowed at the start of a connection",
			Location: &scope.Location,
			Meta:     &chainedConn.Normal.Meta,
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
				Message: fmt.Sprintf(
					"Incompatible types: %v -> %v: %v",
					analyzedSenders[i], chainHead, err.Error(),
				),
				Location: &scope.Location,
				Meta:     &chainedConn.Meta,
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
				Message:  err.Error(),
				Location: &scope.Location,
				Meta:     &chainedConn.Meta,
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
				Meta:     &sender.Meta,
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
		sender.Binary == nil &&
		sender.Ternary == nil &&
		len(sender.StructSelector) == 0 {
		return nil, nil, &compiler.Error{
			Message:  "Sender in network must contain port address, constant reference or message literal",
			Location: &scope.Location,
			Meta:     &sender.Meta,
		}
	}

	if sender.Range != nil && prevChainLink == nil {
		return nil, nil, &compiler.Error{
			Message:  "range expression cannot be used in non-chained connection",
			Location: &scope.Location,
			Meta:     &sender.Meta,
		}
	}

	if sender.Const != nil && len(prevChainLink) != 0 {
		return nil, nil, &compiler.Error{
			Message:  "constant cannot be used in chained connection",
			Location: &scope.Location,
			Meta:     &sender.Meta,
		}
	}

	if len(sender.StructSelector) > 0 && prevChainLink == nil {
		return nil, nil, &compiler.Error{
			Message:  "struct selectors cannot be used in non-chained connection",
			Location: &scope.Location,
			Meta:     &sender.Meta,
		}
	}

	if sender.Ternary != nil {
		// analyze the condition part
		_, condType, err := a.analyzeSender(
			sender.Ternary.Condition,
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
				Meta:     &sender.Ternary.Meta,
			}.Wrap(err)
		}

		// ensure the condition is of boolean type
		boolType := ts.Expr{
			Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "bool"}},
		}
		if err := a.resolver.IsSubtypeOf(*condType, boolType, scope); err != nil {
			return nil, nil, &compiler.Error{
				Message:  "Condition of ternary expression must be of boolean type",
				Location: &scope.Location,
				Meta:     &sender.Ternary.Meta,
			}
		}

		// analyze the trueVal part
		_, trueValType, err := a.analyzeSender(
			sender.Ternary.Left,
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
				Meta:     &sender.Ternary.Meta,
			}.Wrap(err)
		}

		// analyze the falseVal part
		_, _, err = a.analyzeSender(
			sender.Ternary.Right,
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
				Meta:     &sender.Ternary.Meta,
			}.Wrap(err)
		}

		return &sender, trueValType, nil
	}

	if sender.Binary != nil {
		_, leftType, err := a.analyzeSender(
			sender.Binary.Left,
			scope,
			iface,
			nodes,
			nodesIfaces,
			nodesUsage,
			prevChainLink,
		)
		if err != nil {
			return nil, nil, err
		}

		_, rightType, err := a.analyzeSender(
			sender.Binary.Right,
			scope,
			iface,
			nodes,
			nodesIfaces,
			nodesUsage,
			prevChainLink,
		)
		if err != nil {
			return nil, nil, err
		}

		var constr ts.Expr
		switch sender.Binary.Operator {
		case src.AddOp:
			constr = ts.Expr{
				Lit: &ts.LitExpr{
					Union: []ts.Expr{
						{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
						{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "float"}}},
						{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
					},
				},
			}
		case src.SubOp, src.MulOp, src.DivOp:
			constr = ts.Expr{
				Lit: &ts.LitExpr{
					Union: []ts.Expr{
						{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
						{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "float"}}},
					},
				},
			}
		case src.ModOp:
			constr = ts.Expr{
				Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
			}
		case src.EqOp:
			constr = ts.Expr{
				Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "any"}},
			}
		default:
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf(
					"Unsupported binary operator: %v",
					sender.Binary.Operator,
				),
				Location: &scope.Location,
				Meta:     &sender.Binary.Meta,
			}
		}

		if err := a.resolver.IsSubtypeOf(*leftType, constr, scope); err != nil {
			return nil, nil, &compiler.Error{
				Message:  fmt.Sprintf("Invalid left operand type for %s: %v", sender.Binary.Operator, err),
				Location: &scope.Location,
				Meta:     &sender.Binary.Meta,
			}
		}

		if err := a.resolver.IsSubtypeOf(*rightType, constr, scope); err != nil {
			return nil, nil, &compiler.Error{
				Message:  fmt.Sprintf("Invalid right operand type for %s: %v", sender.Binary.Operator, err),
				Location: &scope.Location,
				Meta:     &sender.Binary.Meta,
			}
		}

		// desugarer needs this information to use overloaded components
		// it could figure this out itself but it's extra work
		sender.Binary.AnalyzedType = *leftType

		var resultType ts.Expr
		if sender.Binary.Operator == src.EqOp {
			resultType = ts.Expr{
				Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "bool"}},
			}
		} else {
			resultType = *leftType
		}

		return &sender, &resultType, nil
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
			Meta:     &sender.Meta,
		}.Wrap(err)
	}

	if sender.PortAddr != nil {
		if !isSenderArr && sender.PortAddr.Idx != nil {
			return nil, nil, &compiler.Error{
				Message:  "Index for non-array port",
				Meta:     &sender.PortAddr.Meta,
				Location: &scope.Location,
			}
		}

		if isSenderArr && sender.PortAddr.Idx == nil {
			return nil, nil, &compiler.Error{
				Message:  "Index needed for array outport",
				Meta:     &sender.PortAddr.Meta,
				Location: &scope.Location,
			}
		}

		if err := netNodesUsage(nodesUsage).trackOutportUsage(*sender.PortAddr); err != nil {
			return nil, nil, &compiler.Error{
				Message:  err.Error(),
				Location: &scope.Location,
				Meta:     &sender.PortAddr.Meta,
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
			Meta:     &conn.Meta,
		}.Wrap(err)
	}
	if !isArray {
		return &compiler.Error{
			Message:  "Non-array outport in array-bypass connection",
			Location: &scope.Location,
			Meta:     &arrBypassConn.SenderOutport.Meta,
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
			Meta:     &conn.Meta,
		}.Wrap(err)
	}
	if !isArray {
		return &compiler.Error{
			Message:  "Non-array outport in array-bypass connection",
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
			Message: fmt.Sprintf(
				"Incompatible types: %v -> %v: %v",
				arrBypassConn.SenderOutport, arrBypassConn.ReceiverInport, err.Error(),
			),
			Location: &scope.Location,
			Meta:     &conn.Meta,
		}
	}

	if err := netNodesUsage(nodesUsage).trackOutportUsage(arrBypassConn.SenderOutport); err != nil {
		return &compiler.Error{
			Message:  err.Error(),
			Location: &scope.Location,
			Meta:     &conn.Meta,
		}
	}

	if err := netNodesUsage(nodesUsage).trackInportUsage(arrBypassConn.ReceiverInport); err != nil {
		return &compiler.Error{
			Message:  err.Error(),
			Location: &scope.Location,
			Meta:     &conn.Meta,
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
			Message:  "Unused inports",
			Location: &scope.Location,
			Meta:     &compInterface.Meta,
		}
	}

	for inportName := range compInterface.IO.In {
		if _, ok := inportsUsage.Out[inportName]; !ok { // note that self inports are outports for the network
			return &compiler.Error{
				Message:  fmt.Sprintf("Unused inport: %v", inportName),
				Location: &scope.Location,
			}
		}
	}

	// 2. every self-outport must be used
	outportsUsage, ok := nodesUsage["out"]
	if !ok {
		return &compiler.Error{
			Message:  "Unused outports",
			Location: &scope.Location,
			Meta:     &compInterface.Meta,
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
			Message:  fmt.Sprintf("Unused outport: %v", outportName),
			Location: &scope.Location,
		}
	}

	// 3. check sub-nodes usage in network
	for nodeName, nodeIface := range nodesIfaces {
		// every sub-node must be used
		nodeUsage, ok := nodesUsage[nodeName]
		if !ok {
			return &compiler.Error{
				Message:  fmt.Sprintf("Unused node found: %v", nodeName),
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
				Message: fmt.Sprintf(
					"Unused node inport: %v:%v",
					nodeName,
					inportName,
				),
				Location: &scope.Location,
				Meta:     compiler.Pointer(nodeIface.iface.IO.In[inportName].Meta),
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
						Message: fmt.Sprintf(
							"if node has error guard '?' it's ':err' outport must not be explicitly used in the network: %v",
							nodeName,
						),
						Location: &scope.Location,
						Meta:     &port.Meta,
					}
				}

				atLeastOneOutportIsUsed = true
				continue
			}

			if outportName == "err" && !nodes[nodeName].ErrGuard {
				return &compiler.Error{
					Message:  fmt.Sprintf("unhandled error: %v:err", nodeName),
					Location: &scope.Location,
					Meta:     &port.Meta,
				}
			}
		}

		if !atLeastOneOutportIsUsed {
			if _, ok := nodeUsage.Out[""]; ok && len(nodeIface.iface.IO.Out) == 1 {
				continue
			}

			return &compiler.Error{
				Message:  fmt.Sprintf("All node's outports are unused: %v", nodeName),
				Location: &scope.Location,
				Meta:     &nodeIface.iface.Meta,
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
						Message: fmt.Sprintf(
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
						Message: fmt.Sprintf(
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
			Message:  "Component cannot read from self inport",
			Location: &scope.Location,
			Meta:     &receiverSide.Meta,
		}
	}

	if receiverSide.Node == "out" {
		outports := iface.IO.Out

		outport, ok := outports[receiverSide.Port]
		if !ok {
			return ts.Expr{}, false, &compiler.Error{
				Message:  fmt.Sprintf("Referenced inport not found in component's interface: %v", receiverSide.Port),
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
				Message:  err.Error(),
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
			Message:  fmt.Sprintf("Node not found '%v'", portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, false, &compiler.Error{
			Message:  fmt.Sprintf("Referenced node not found: %v", portAddr.Node),
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
		scope.Relocate(nodeIface.location),
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
				Message:  fmt.Sprintf("node '%v' has multiple ports but no port name", portAddr.Node),
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
			Message: fmt.Sprintf(
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
			Message:  err.Error(),
			Location: &scope.Location,
			Meta:     &port.Meta,
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
				Message:  "fan-in with struct selectors is not supported",
				Location: &scope.Location,
				Meta:     &senderSide.Meta,
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
			return src.ConnectionSender{}, ts.Expr{}, false, compiler.Error{
				Location: &scope.Location,
				Meta:     &senderSide.Meta,
			}.Wrap(err)
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
			Message:  "Component cannot read from self outport",
			Location: &scope.Location,
			Meta:     &senderSidePortAddr.Meta,
		}
	}

	if senderSidePortAddr.Node == "in" {
		inports := iface.IO.In

		inport, ok := inports[senderSidePortAddr.Port]
		if !ok {
			return ts.Expr{}, false, &compiler.Error{
				Message:  fmt.Sprintf("Referenced inport not found in component's interface: %v", senderSidePortAddr.Port),
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
				Message:  err.Error(),
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

func (a Analyzer) getConstSenderType(
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
			Message:  "Literal sender type is empty",
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
			Message:  err.Error(),
			Location: &scope.Location,
			Meta:     &constSender.Value.Message.Meta,
		}
	}

	if err := a.validateLiteralSender(resolvedExpr); err != nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Message:  err.Error(),
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
			Message:  fmt.Sprintf("Referenced node not found: %v", portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return ts.Expr{}, false, &compiler.Error{
			Message:  fmt.Sprintf("Referenced node not found: %v", portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	return a.getResolvedPortType(
		nodeIface.iface.IO.Out,
		nodeIface.iface.TypeParams.Params,
		portAddr,
		node,
		scope.Relocate(nodeIface.location),
	)
}

func (a Analyzer) getResolvedConstTypeByRef(ref core.EntityRef, scope src.Scope) (ts.Expr, *compiler.Error) {
	entity, location, err := scope.Entity(ref)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message:  err.Error(),
			Location: &scope.Location,
			Meta:     &ref.Meta,
		}
	}

	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, &compiler.Error{
			Message:  fmt.Sprintf("%v: %v", errors.New("Entity found but is not constant"), entity.Kind),
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

	scope = scope.Relocate(location)

	resolvedExpr, err := a.resolver.ResolveExpr(entity.Const.TypeExpr, scope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message:  err.Error(),
			Location: &scope.Location,
			Meta:     &entity.Const.Value.Message.Meta,
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
			Message:  fmt.Sprintf("Type not struct: %v", senderType.String()),
			Location: &scope.Location,
		}
	}

	curField := selectors[0]
	fieldType, ok := senderType.Lit.Struct[curField]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Message:  fmt.Sprintf("struct field '%v' not found", curField),
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
		Message:  "Chained connection must start with port address or range expression",
		Location: &scope.Location,
		Meta:     &chainHead.Meta,
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

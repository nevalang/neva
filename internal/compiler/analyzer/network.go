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
		return nil, err
	}

	if err := a.analyzeNetPortsUsage(
		compInterface,
		nodesIfaces,
		hasGuard,
		nodesUsage,
		nodes,
	); err != nil {
		return nil, err
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
	analyzedSenders, resolvedSenderTypes, err := a.analyzeSenders(
		normConn.Senders,
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
		normConn.Receivers,
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
		Senders:   analyzedSenders,
		Receivers: analyzedReceiverSide,
		Meta:      normConn.Meta,
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
	switch {
	case receiver.PortAddr != nil:
		err := a.analyzePortAddrReceiver(
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
			PortAddr: receiver.PortAddr, // no need to change anything
			Meta:     receiver.Meta,
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
	case receiver.Switch != nil:
		analyzedSwitchConns, analyzedDefault, err := a.analyzeSwitchReceiver(
			receiver,
			iface,
			nodes,
			nodesIfaces,
			scope,
			nodesUsage,
			analyzedSenders,
			resolvedSenderTypes,
		)
		if err != nil {
			return nil, err
		}
		return &src.ConnectionReceiver{
			Switch: &src.Switch{
				Cases:   analyzedSwitchConns,
				Default: analyzedDefault,
			},
			Meta: receiver.Meta,
		}, nil
	}

	return nil, &compiler.Error{
		Message: "Connection must have receiver-side",
		Meta:    &receiver.Meta,
	}
}

func (a Analyzer) analyzeSwitchReceiver(
	receiver src.ConnectionReceiver,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
	nodesUsage map[string]netNodeUsage,
	analyzedSenders []src.ConnectionSender,
	resolvedSenderTypes []*ts.Expr,
) ([]src.NormalConnection, []src.ConnectionReceiver, *compiler.Error) {
	analyzedSwitchConns := make([]src.NormalConnection, 0, len(receiver.Switch.Cases))

	for _, switchConn := range receiver.Switch.Cases {
		// all option-senders must be subtypes of their branch-receivers
		analyzedSwitchConn, err := a.analyzeNormalConnection(
			&switchConn,
			iface,
			nodes,
			nodesIfaces,
			scope,
			nodesUsage,
			nil,
		)
		if err != nil {
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf("Invalid switch case: %v", err),
				Meta:    &switchConn.Meta,
			}
		}

		analyzedSwitchConns = append(analyzedSwitchConns, *analyzedSwitchConn)

		// all incoming senders must be subtypes of each option-sender
		// (both incoming senders and switch's option-senders might be slice)
		for _, switchSender := range switchConn.Senders {
			_, switchSenderType, _, err := a.getResolvedSenderType(
				switchSender,
				iface,
				nodes,
				nodesIfaces,
				scope,
				nil,
			)
			if err != nil {
				return nil, nil, &compiler.Error{
					Message: fmt.Sprintf("Invalid switch case sender: %v", err),
					Meta:    &switchSender.Meta,
				}
			}

			for i, resolvedSenderType := range resolvedSenderTypes {
				if err := a.resolver.IsSubtypeOf(*resolvedSenderType, switchSenderType, scope); err != nil {
					return nil, nil, &compiler.Error{
						Message: fmt.Sprintf(
							"Incompatible types in switch: %v -> %v: %v",
							analyzedSenders[i], switchSender, err.Error(),
						),
						Meta: &switchSender.Meta,
					}
				}
			}
		}
	}

	if receiver.Switch.Default == nil {
		return nil, nil, &compiler.Error{
			Message: "Switch must have a default case",
			Meta:    &receiver.Meta,
		}
	}

	analyzedDefault, err := a.analyzeReceiverSide(
		receiver.Switch.Default,
		scope,
		iface,
		nodes,
		nodesIfaces,
		nodesUsage,
		resolvedSenderTypes,
		analyzedSenders,
	)
	if err != nil {
		return nil, nil, err
	}

	return analyzedSwitchConns, analyzedDefault, nil
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
) *compiler.Error {
	resolvedPortAddr, typeExpr, isArrPort, err := a.getReceiverPortType(
		portAddr,
		iface,
		nodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return compiler.Error{
			Meta: &portAddr.Meta,
		}.Wrap(err)
	}

	if !isArrPort && portAddr.Idx != nil {
		return &compiler.Error{
			Message: "Index for non-array port",
			Meta:    &portAddr.Meta,
		}
	}

	if isArrPort && portAddr.Idx == nil {
		return &compiler.Error{
			Message: "Index needed for array inport",
			Meta:    &portAddr.Meta,
		}
	}

	for i, resolvedSenderType := range resolvedSenderTypes {
		if err := a.resolver.IsSubtypeOf(*resolvedSenderType, typeExpr, scope); err != nil {
			return &compiler.Error{
				Message: fmt.Sprintf(
					"Incompatible types: %v -> %v: %v",
					analyzedSenders[i], portAddr, err.Error(),
				),
				Meta: &portAddr.Meta,
			}
		}
	}

	// sometimes port name is omitted and we need to resolve it first
	// but it's important not to return it, so syntax sugar remains untouched
	// otherwise desugarer won't be able to properly desugar such port-addresses
	if err := netNodesUsage(nodesUsage).trackInportUsage(resolvedPortAddr); err != nil {
		return &compiler.Error{
			Message: err.Error(),
			Meta:    &portAddr.Meta,
		}
	}

	return nil
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
			Message: "chained connection must be a normal connection",
			Meta:    &chainedConn.Meta,
		}
	}

	if len(chainedConn.Normal.Senders) != 1 {
		return src.Connection{}, &compiler.Error{
			Message: "multiple senders are only allowed at the start of a connection",
			Meta:    &chainedConn.Normal.Meta,
		}
	}

	chainHead := chainedConn.Normal.Senders[0]

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
				Meta: &chainedConn.Meta,
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
				Message: err.Error(),
				Meta:    &chainedConn.Meta,
			}
		}
	}

	return analyzedChainedConn, nil
}

func (a Analyzer) analyzeSenders(
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
				Meta: &sender.Meta,
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

// analyzeSender validates sender, marks it as used if it's port-address and returns its type.
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
		sender.Unary == nil &&
		sender.Ternary == nil &&
		len(sender.StructSelector) == 0 {
		return nil, nil, &compiler.Error{
			Message: "Sender in network must contain port address, constant reference or message literal",
			Meta:    &sender.Meta,
		}
	}

	// TODO support unary

	if sender.Range != nil && len(prevChainLink) == 0 {
		return nil, nil, &compiler.Error{
			Message: "range expression cannot be used in non-chained connection",
			Meta:    &sender.Meta,
		}
	}

	if len(sender.StructSelector) > 0 && len(prevChainLink) == 0 {
		return nil, nil, &compiler.Error{
			Message: "struct selectors cannot be used in non-chained connection",
			Meta:    &sender.Meta,
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
				Meta: &sender.Ternary.Meta,
			}.Wrap(err)
		}

		// ensure the condition is of boolean type
		boolType := ts.Expr{
			Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "bool"}},
		}
		if err := a.resolver.IsSubtypeOf(*condType, boolType, scope); err != nil {
			return nil, nil, &compiler.Error{
				Message: "Condition of ternary expression must be of boolean type",
				Meta:    &sender.Ternary.Meta,
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
				Meta: &sender.Ternary.Meta,
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
				Meta: &sender.Ternary.Meta,
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

		constr, err := a.getOperatorConstraint(*sender.Binary)
		if err != nil {
			return nil, nil, err
		}

		if err := a.resolver.IsSubtypeOf(*leftType, constr, scope); err != nil {
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf("Invalid left operand type for %s: %v", sender.Binary.Operator, err),
				Meta:    &sender.Binary.Meta,
			}
		}

		if err := a.resolver.IsSubtypeOf(*rightType, constr, scope); err != nil {
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf("Invalid right operand type for %s: %v", sender.Binary.Operator, err),
				Meta:    &sender.Binary.Meta,
			}
		}

		// desugarer needs this information to use overloaded components
		// it could figure this out itself but it's extra work
		sender.Binary.AnalyzedType = *leftType

		resultType := a.getBinaryExprType(sender.Binary.Operator, *leftType)

		return &sender, &resultType, nil
	}

	resolvedSenderAddr, resolvedSenderType, isSenderArr, err := a.getResolvedSenderType(
		sender,
		iface,
		nodes,
		nodesIfaces,
		scope,
		prevChainLink,
	)
	if err != nil {
		return nil, nil, compiler.Error{
			Meta: &sender.Meta,
		}.Wrap(err)
	}

	if sender.PortAddr != nil {
		if !isSenderArr && sender.PortAddr.Idx != nil {
			return nil, nil, &compiler.Error{
				Message: "Index for non-array port",
				Meta:    &sender.PortAddr.Meta,
			}
		}

		if isSenderArr && sender.PortAddr.Idx == nil {
			return nil, nil, &compiler.Error{
				Message: "Index needed for array outport",
				Meta:    &sender.PortAddr.Meta,
			}
		}

		if sender.PortAddr.Port == "err" && nodes[sender.PortAddr.Node].ErrGuard {
			return nil, nil, &compiler.Error{
				Message: "if node has error guard '?' it's ':err' outport must not be explicitly used in the network",
				Meta:    &sender.PortAddr.Meta,
			}
		}

		// it's important to track resolved port address here
		// because sometimes port name is omitted and we need to resolve it first
		// but it's important not to return it, so syntax sugar remains untouched
		// otherwise desugarer won't be able to properly desugar such port-addresses
		if err := netNodesUsage(nodesUsage).trackOutportUsage(*resolvedSenderAddr.PortAddr); err != nil {
			return nil, nil, &compiler.Error{
				Message: err.Error(),
				Meta:    &sender.PortAddr.Meta,
			}
		}

		return &src.ConnectionSender{
			PortAddr: sender.PortAddr,
			Meta:     sender.Meta,
		}, &resolvedSenderType, nil
	}

	return &resolvedSenderAddr, &resolvedSenderType, nil
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

	_, senderType, isArray, err := a.getPortSenderType(
		arrBypassConn.SenderOutport,
		scope,
		iface,
		nodes,
		nodesIfaces,
	)
	if err != nil {
		return compiler.Error{
			Meta: &conn.Meta,
		}.Wrap(err)
	}
	if !isArray {
		return &compiler.Error{
			Message: "Non-array outport in array-bypass connection",
			Meta:    &arrBypassConn.SenderOutport.Meta,
		}
	}

	_, receiverType, isArray, err := a.getReceiverPortType(
		arrBypassConn.ReceiverInport,
		iface,
		nodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return compiler.Error{
			Meta: &conn.Meta,
		}.Wrap(err)
	}
	if !isArray {
		return &compiler.Error{
			Message: "Non-array outport in array-bypass connection",
			Meta:    &arrBypassConn.SenderOutport.Meta,
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
			Meta: &conn.Meta,
		}
	}

	if err := netNodesUsage(nodesUsage).trackOutportUsage(arrBypassConn.SenderOutport); err != nil {
		return &compiler.Error{
			Message: err.Error(),
			Meta:    &conn.Meta,
		}
	}

	if err := netNodesUsage(nodesUsage).trackInportUsage(arrBypassConn.ReceiverInport); err != nil {
		return &compiler.Error{
			Message: err.Error(),
			Meta:    &conn.Meta,
		}
	}

	return nil
}

func (a Analyzer) analyzeNetPortsUsage(
	compInterface src.Interface,
	nodesIfaces map[string]foundInterface,
	hasGuard bool,
	nodesUsage map[string]netNodeUsage,
	nodes map[string]src.Node,
) *compiler.Error {
	// 1. every self inport must be used
	inportsUsage, ok := nodesUsage["in"]
	if !ok {
		return &compiler.Error{
			Message: "Unused inports",
			Meta:    &compInterface.Meta,
		}
	}

	for inportName := range compInterface.IO.In {
		if _, ok := inportsUsage.Out[inportName]; !ok { // note that self inports are outports for the network
			return &compiler.Error{
				Message: fmt.Sprintf("Unused inport: %v", inportName),
			}
		}
	}

	// 2. every self-outport must be used
	outportsUsage, ok := nodesUsage["out"]
	if !ok {
		return &compiler.Error{
			Message: "Component must use its outports",
			Meta:    &compInterface.Meta,
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
			Message: fmt.Sprintf("Unused outport: %v", outportName),
		}
	}

	// 3. check sub-nodes usage in network
	for nodeName, nodeIface := range nodesIfaces {
		nodeMeta := nodes[nodeName].Meta

		// every sub-node must be used
		nodeUsage, ok := nodesUsage[nodeName]
		if !ok {
			return &compiler.Error{
				Message: fmt.Sprintf("Unused node found: %v", nodeName),
				Meta:    &nodeMeta,
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
				Meta: &nodeMeta,
			}
		}

		if len(nodeIface.iface.IO.Out) == 0 { // e.g. Del
			continue
		}

		// :err outport must always be used + at least one outport must be used in general
		atLeastOneOutportIsUsed := false
		for outportName := range nodeIface.iface.IO.Out {
			if _, ok := nodeUsage.Out[outportName]; ok {
				atLeastOneOutportIsUsed = true
				continue
			}

			if outportName == "err" && !nodes[nodeName].ErrGuard {
				return &compiler.Error{
					Message: fmt.Sprintf("unhandled error: %v:err", nodeName),
					Meta:    &nodeMeta,
				}
			}
		}

		if !atLeastOneOutportIsUsed {
			if _, ok := nodeUsage.Out[""]; ok && len(nodeIface.iface.IO.Out) == 1 {
				continue
			}
			return &compiler.Error{
				Message: fmt.Sprintf("All node's outports are unused: %v", nodeName),
				Meta:    &nodeMeta,
			}
		}
	}

	// 4. check that array ports are used correctly (from 0 and without holes)
	for nodeName, nodeUsage := range nodesUsage {
		nodeMeta := nodes[nodeName].Meta

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
						Meta: &nodeMeta,
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
						Meta: &nodeMeta,
					}
				}
			}
		}
	}

	return nil
}

// getReceiverPortType returns resolved port-addr, type expr and isArray bool.
// Resolved port is equal to the given one unless it was an "" empty string.
func (a Analyzer) getReceiverPortType(
	receiverSide src.PortAddr,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
) (src.PortAddr, ts.Expr, bool, *compiler.Error) {
	if receiverSide.Node == "in" {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: "Component cannot read from self inport",
			Meta:    &receiverSide.Meta,
		}
	}

	if receiverSide.Node == "out" {
		outports := iface.IO.Out

		outport, ok := outports[receiverSide.Port]
		if !ok {
			return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
				Message: fmt.Sprintf("Referenced inport not found in component's interface: %v", receiverSide.Port),
				Meta:    &receiverSide.Meta,
			}
		}

		resolvedOutportType, err := a.resolver.ResolveExprWithFrame(
			outport.TypeExpr,
			iface.TypeParams.ToFrame(),
			scope,
		)
		if err != nil {
			return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
				Message: err.Error(),
				Meta:    &receiverSide.Meta,
			}
		}

		return receiverSide, resolvedOutportType, outport.IsArray, nil
	}

	resolvedReceiver, nodeInportType, isArray, err := a.getNodeInportType(
		receiverSide, nodes, nodesIfaces, scope,
	)
	if err != nil {
		return src.PortAddr{}, ts.Expr{}, false, compiler.Error{
			Meta: &receiverSide.Meta,
		}.Wrap(err)
	}

	return resolvedReceiver, nodeInportType, isArray, nil
}

// getNodeInportType returns resolved port-addr, type expr and isArray bool.
// Resolved port is equal to the given one unless it was an "" empty string.
func (a Analyzer) getNodeInportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
) (src.PortAddr, ts.Expr, bool, *compiler.Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: fmt.Sprintf("Node not found '%v'", portAddr.Node),
			Meta:    &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: fmt.Sprintf("Referenced node not found: %v", portAddr.Node),
			Meta:    &portAddr.Meta,
		}
	}

	resolvedPortAddr, resolvedInportType, isArray, err := a.getResolvedPortType(
		nodeIface.iface.IO.In,
		nodeIface.iface.TypeParams.Params,
		portAddr,
		node,
		scope.Relocate(nodeIface.location),
		true,
	)
	if err != nil {
		return src.PortAddr{}, ts.Expr{}, false, compiler.Error{
			Meta: &portAddr.Meta,
		}.Wrap(err)
	}

	return resolvedPortAddr, resolvedInportType, isArray, nil
}

// getResolvedPortType returns resolved port-addr, type expr and isArray bool.
// Resolved port is equal to the given one unless it was an "" empty string.
func (a Analyzer) getResolvedPortType(
	ports map[string]src.Port,
	nodeIfaceParams []ts.Param,
	portAddr src.PortAddr,
	node src.Node,
	scope src.Scope,
	isInput bool,
) (src.PortAddr, ts.Expr, bool, *compiler.Error) {
	if portAddr.Port == "" {
		if len(ports) == 1 || (!isInput && len(ports) == 2 && node.ErrGuard) {
			for name := range ports {
				portAddr.Port = name
				break
			}
		} else {
			return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
				Message: fmt.Sprintf("node '%v' has multiple ports - port name must be specified", portAddr.Node),
				Meta:    &portAddr.Meta,
			}
		}
	}

	port, ok := ports[portAddr.Port]
	if !ok {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: fmt.Sprintf(
				"Port not found '%v'",
				portAddr,
			),
			Meta: &portAddr.Meta,
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
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: err.Error(),
			Meta:    &port.Meta,
		}
	}

	return portAddr, resolvedPortType, port.IsArray, nil
}

func (a Analyzer) getResolvedSenderType(
	sender src.ConnectionSender,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
	prevChainLink []src.ConnectionSender,
) (src.ConnectionSender, ts.Expr, bool, *compiler.Error) {
	if sender.Const != nil {
		resolvedConst, resolvedExpr, err := a.getConstSenderType(*sender.Const, scope)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, err
		}

		return src.ConnectionSender{
			Const: &resolvedConst,
			Meta:  sender.Meta,
		}, resolvedExpr, false, nil
	}

	if sender.Range != nil {
		// range sends stream<int> from its :data outport
		rangeType := ts.Expr{
			Inst: &ts.InstExpr{
				Ref: core.EntityRef{Name: "stream"},
				Args: []ts.Expr{{
					Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
				}},
			},
		}
		return sender, rangeType, false, nil
	}

	if len(sender.StructSelector) > 0 {
		_, chainLinkType, _, err := a.getResolvedSenderType(
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
			sender.StructSelector,
			scope,
		)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, compiler.Error{
				Meta: &sender.Meta,
			}.Wrap(err)
		}

		return sender, lastFieldType, false, nil
	}

	// logic of getting type for binary expr partially duplicates logic of validating it
	if sender.Binary != nil {
		_, leftType, _, err := a.getResolvedSenderType(
			sender.Binary.Left,
			iface,
			nodes,
			nodesIfaces,
			scope,
			prevChainLink,
		)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, err
		}
		return sender, a.getBinaryExprType(sender.Binary.Operator, leftType), false, nil
	}

	// logic of getting type for ternary expr partially duplicates logic of validating it
	// so we have to duplicate some code from "analyzeSender", but it should be possible to refactor
	if sender.Ternary != nil {
		_, trueValType, _, err := a.getResolvedSenderType(
			sender.Ternary.Left,
			iface,
			nodes,
			nodesIfaces,
			scope,
			prevChainLink,
		)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, compiler.Error{
				Meta: &sender.Ternary.Meta,
			}.Wrap(err)
		}
		// FIXME: perfectly we need to return union, but it's not yet supported
		// see https://github.com/nevalang/neva/issues/737
		return sender, trueValType, false, nil
	}

	// handle port-address sender
	resolvedPort, resolvedExpr, isArr, err := a.getPortSenderType(
		*sender.PortAddr,
		scope,
		iface,
		nodes,
		nodesIfaces,
	)
	if err != nil {
		return src.ConnectionSender{}, ts.Expr{}, false, err
	}

	return src.ConnectionSender{
		PortAddr: &resolvedPort,
		Meta:     sender.Meta,
	}, resolvedExpr, isArr, nil
}

// getBinaryExprType returns type of the result of binary expression.
func (Analyzer) getBinaryExprType(operator src.BinaryOperator, leftType ts.Expr) ts.Expr {
	if operator == src.EqOp ||
		operator == src.NeOp ||
		operator == src.GtOp ||
		operator == src.LtOp ||
		operator == src.GeOp ||
		operator == src.LeOp ||
		operator == src.AndOp ||
		operator == src.OrOp {
		return ts.Expr{
			Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "bool"}},
		}
	}
	return leftType
}

func (Analyzer) getOperatorConstraint(binary src.Binary) (ts.Expr, *compiler.Error) {
	switch binary.Operator {
	case src.AddOp:
		return ts.Expr{
			Lit: &ts.LitExpr{
				Union: []ts.Expr{
					{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
					{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "float"}}},
					{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
				},
			},
		}, nil
	case src.SubOp, src.MulOp, src.DivOp:
		return ts.Expr{
			Lit: &ts.LitExpr{
				Union: []ts.Expr{
					{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
					{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "float"}}},
				},
			},
		}, nil
	case src.ModOp, src.PowOp:
		return ts.Expr{
			Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
		}, nil
	case src.EqOp, src.NeOp:
		return ts.Expr{
			Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "any"}},
		}, nil
	case src.GtOp, src.LtOp, src.GeOp, src.LeOp:
		return ts.Expr{
			Lit: &ts.LitExpr{
				Union: []ts.Expr{
					{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
					{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "float"}}},
					{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
				},
			},
		}, nil
	case src.AndOp, src.OrOp:
		return ts.Expr{
			Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "bool"}},
		}, nil
	case src.BitAndOp, src.BitOrOp, src.BitXorOp, src.BitLshOp, src.BitRshOp:
		return ts.Expr{
			Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
		}, nil
	default:
		return ts.Expr{}, &compiler.Error{
			Message: fmt.Sprintf(
				"Unsupported binary operator: %v",
				binary.Operator,
			),
			Meta: &binary.Meta,
		}
	}
}

// getPortSenderType returns resolved port-addr, type expr and isArray bool.
// Resolved port is equal to the given one unless it was an "" empty string.
func (a Analyzer) getPortSenderType(
	senderSidePortAddr src.PortAddr,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
) (src.PortAddr, ts.Expr, bool, *compiler.Error) {
	if senderSidePortAddr.Node == "out" {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: "Component cannot read from self outport",
			Meta:    &senderSidePortAddr.Meta,
		}
	}

	if senderSidePortAddr.Node == "in" {
		inports := iface.IO.In

		inport, ok := inports[senderSidePortAddr.Port]
		if !ok {
			return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
				Message: fmt.Sprintf("Referenced inport not found in component's interface: %v", senderSidePortAddr.Port),
				Meta:    &senderSidePortAddr.Meta,
			}
		}

		resolvedInportType, err := a.resolver.ResolveExprWithFrame(
			inport.TypeExpr,
			iface.TypeParams.ToFrame(),
			scope,
		)
		if err != nil {
			return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
				Message: err.Error(),
				Meta:    &senderSidePortAddr.Meta,
			}
		}

		return senderSidePortAddr, resolvedInportType, inport.IsArray, nil
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
				Meta: &constSender.Value.Ref.Meta,
			}.Wrap(err)
		}
		return constSender, expr, nil
	}

	if constSender.Value.Message == nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Message: "Literal sender type is empty",
			Meta:    &constSender.Meta,
		}
	}

	resolvedExpr, err := a.resolver.ResolveExpr(
		constSender.TypeExpr,
		scope,
	)
	if err != nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &constSender.Value.Message.Meta,
		}
	}

	if err := a.validateLiteralSender(resolvedExpr); err != nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &constSender.Value.Message.Meta,
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

// getNodeOutportType returns resolved port-addr, type expr and isArray bool.
// Resolved port is equal to the given one unless it was an "" empty string.
func (a Analyzer) getNodeOutportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
) (src.PortAddr, ts.Expr, bool, *compiler.Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: fmt.Sprintf("Referenced node not found: %v", portAddr.Node),
			Meta:    &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: fmt.Sprintf("Referenced node not found: %v", portAddr.Node),
			Meta:    &portAddr.Meta,
		}
	}

	return a.getResolvedPortType(
		nodeIface.iface.IO.Out,
		nodeIface.iface.TypeParams.Params,
		portAddr,
		node,
		scope.Relocate(nodeIface.location),
		false,
	)
}

func (a Analyzer) getResolvedConstTypeByRef(ref core.EntityRef, scope src.Scope) (ts.Expr, *compiler.Error) {
	entity, location, err := scope.Entity(ref)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &ref.Meta,
		}
	}

	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, &compiler.Error{
			Message: fmt.Sprintf("%v: %v", errors.New("Entity found but is not constant"), entity.Kind),
			Meta:    entity.Meta(),
		}
	}

	if entity.Const.Value.Ref != nil {
		expr, err := a.getResolvedConstTypeByRef(*entity.Const.Value.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Meta: &entity.Const.Meta,
			}.Wrap(err)
		}
		return expr, nil
	}

	scope = scope.Relocate(location)

	resolvedExpr, err := a.resolver.ResolveExpr(entity.Const.TypeExpr, scope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &entity.Const.Value.Message.Meta,
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
			Message: fmt.Sprintf("Type not struct: %v", senderType.String()),
		}
	}

	curField := selectors[0]
	fieldType, ok := senderType.Lit.Struct[curField]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Message: fmt.Sprintf("struct field '%v' not found", curField),
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
		_, resolvedType, _, err := a.getNodeInportType(
			*chainHead.PortAddr,
			nodes,
			nodesIfaces,
			scope,
		)
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

	if chainHead.Const != nil {
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
		Message: "Chained connection must start with port address or range expression",
		Meta:    &chainHead.Meta,
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

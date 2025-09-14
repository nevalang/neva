// This file contains the logic for analyzing network connections.
// Some methods here might look like they are related to senders or receivers specifically,
// but they are actually related to both, so they are placed here.
package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

var ErrComplexLiteralSender = errors.New("literal network sender must have primitive type")

// analyzeNetwork must be called after analyzeNodes so we sure nodes are resolved.
func (a Analyzer) analyzeNetwork(
	net []src.Connection, // network to analyze
	iface src.Interface, // resolved interface of the component that contains the network
	hasGuard bool, // whether `?` is used by at least one node in the network
	nodes map[string]src.Node, // nodes of the component that contains the network
	nodesIfaces map[string]foundInterface, // resolved interfaces of the nodes
	scope src.Scope,
) ([]src.Connection, *compiler.Error) {
	nodesUsage := make(map[string]netNodeUsage, len(nodes))

	analyzedConnections, err := a.analyzeConnections(
		net,
		iface,
		nodes,
		nodesIfaces,
		nodesUsage,
		scope,
	)
	if err != nil {
		return nil, err
	}

	if err := a.analyzeNetPortsUsage(
		iface,
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

	analyzedReceivers, err := a.analyzeReceivers(
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
		Receivers: analyzedReceivers,
		Meta:      normConn.Meta,
	}, nil
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
	iface src.Interface, // resolved interface of the component that contains the network
	nodesIfaces map[string]foundInterface, // resolved interfaces of the nodes in the network
	hasGuard bool, // whether `?` is used by at least one node in the network
	nodesUsage map[string]netNodeUsage,
	nodes map[string]src.Node,
) *compiler.Error {
	// 1. every self inport must be used
	inportsUsage, ok := nodesUsage["in"]
	if !ok {
		return &compiler.Error{
			Message: "Unused inports",
			Meta:    &iface.Meta,
		}
	}

	for inportName := range iface.IO.In {
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
			Meta:    &iface.Meta,
		}
	}

	for outportName := range iface.IO.Out {
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
	nodesUsage map[string]netNodeUsage,
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
			nodesUsage,
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
			nodesUsage,
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
			nodesUsage,
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

	// logic of getting type for union sender partially duplicates logic of validating it
	// so we have to duplicate some code from "analyzeSender", but it should be possible to refactor
	if sender.Union != nil {
		entity, _, err := scope.GetType(sender.Union.EntityRef)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, &compiler.Error{
				Message: fmt.Sprintf("failed to resolve union type: %v", err),
				Meta:    &sender.Meta,
			}
		}

		resolvedUnionType, typeExprErr := a.analyzeTypeExpr(*entity.BodyExpr, scope)
		if typeExprErr != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, &compiler.Error{
				Message: fmt.Sprintf("failed to resolve union type: %v", typeExprErr),
				Meta:    &sender.Meta,
			}
		}

		tagTypeExpr := resolvedUnionType.Lit.Union[sender.Union.Tag] // we assume it's analyzed already

		// if there's no type-expr (and thus no wrapped sender)
		// it's a "enum-like" union, so there's nothing to analyze
		if tagTypeExpr == nil {
			return sender, resolvedUnionType, false, nil
		}

		// analyze wrapped-sender and get its resolved type
		resolvedWrappedSender, _, analyzeWrappedErr := a.analyzeSender(
			*sender.Union.Data,
			scope,
			iface,
			nodes,
			nodesIfaces,
			nodesUsage,
			prevChainLink,
		)
		if analyzeWrappedErr != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, analyzeWrappedErr
		}

		// return fully resolved union sender, including its wrapped-sender
		return src.ConnectionSender{
			Union: &src.UnionSender{
				EntityRef: sender.Union.EntityRef,
				Tag:       sender.Union.Tag,
				Data:      resolvedWrappedSender,
			},
			Meta: sender.Meta,
		}, resolvedUnionType, false, nil
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
				Union: map[string]*ts.Expr{
					"int": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
					},
					"float": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "float"}},
					},
					"string": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}},
					},
				},
			},
		}, nil
	case src.SubOp, src.MulOp, src.DivOp:
		return ts.Expr{
			Lit: &ts.LitExpr{
				Union: map[string]*ts.Expr{
					"int": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
					},
					"float": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "float"}},
					},
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
				Union: map[string]*ts.Expr{
					"int": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
					},
					"float": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "float"}},
					},
					"string": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}},
					},
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

// checkOperatorOperandTypesWithTypeSystem validates that the operand types are compatible with the operator
// using the type system with proper union construction
func (a Analyzer) checkOperatorOperandTypesWithTypeSystem(binary src.Binary, leftType, rightType ts.Expr, scope src.Scope) *compiler.Error {
	// get the operator constraint (union of allowed types)
	constraint, err := a.getOperatorConstraint(binary)
	if err != nil {
		return err
	}

	// create single-element unions for the operand types to work with type system
	leftUnion := a.createSingleElementUnion(leftType)
	rightUnion := a.createSingleElementUnion(rightType)

	// check that both operands are subtypes of the constraint using the type system
	// this leverages the full power of the type system for complex type resolution
	if err := a.resolver.IsSubtypeOf(leftUnion, constraint, scope); err != nil {
		return &compiler.Error{
			Message: fmt.Sprintf("Invalid left operand type for %s: %v (leftType: %v, leftUnion: %v, constraint: %v)", binary.Operator, err, leftType.String(), leftUnion.String(), constraint.String()),
			Meta:    &binary.Meta,
		}
	}

	if err := a.resolver.IsSubtypeOf(rightUnion, constraint, scope); err != nil {
		return &compiler.Error{
			Message: fmt.Sprintf("Invalid right operand type for %s: %v", binary.Operator, err),
			Meta:    &binary.Meta,
		}
	}

	return nil
}

// createSingleElementUnion creates a union type with a single element matching the given type
func (a Analyzer) createSingleElementUnion(expr ts.Expr) ts.Expr {
	// if the expression is already a union, return it as-is
	if expr.Lit != nil && expr.Lit.Union != nil {
		return expr
	}

	// create a single-element union
	// for primitive types like int, create union { int }
	// for complex types, create union with the type name as the tag
	if expr.Inst != nil {
		typeName := expr.Inst.Ref.String()
		// create a new instance expression with the same type
		tagExpr := ts.Expr{
			Inst: &ts.InstExpr{
				Ref:  expr.Inst.Ref,
				Args: expr.Inst.Args,
			},
		}
		return ts.Expr{
			Lit: &ts.LitExpr{
				Union: map[string]*ts.Expr{
					typeName: &tagExpr,
				},
			},
		}
	}

	// if the expression is a literal, we need to handle it differently
	if expr.Lit != nil {
		// for literal expressions, we can't easily create a union
		// this shouldn't happen for operator operands, but let's handle it
		return expr
	}

	// fallback: return the expression as-is if we can't create a union
	return expr
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
				Union:        constSender.Value.Message.Union,
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
		resolvedExpr.Lit.Union == nil {
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

func (a Analyzer) getResolvedConstTypeByRef(
	ref core.EntityRef,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	constant, loc, err := scope.GetConst(ref)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &ref.Meta,
		}
	}

	if constant.Value.Ref != nil {
		expr, err := a.getResolvedConstTypeByRef(*constant.Value.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Meta: &constant.Meta,
			}.Wrap(err)
		}
		return expr, nil
	}

	scope = scope.Relocate(loc)

	resolvedExpr, err := a.resolver.ResolveExpr(constant.TypeExpr, scope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &constant.Value.Message.Meta,
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

package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

// getSwitchCaseOutportType resolves the specialized type for a Switch case output.
// It inspects the Switch generic type T. If T is a Union, it looks up the input connection
// to determine which Union tag matches this case, and returns that tag's type.
// If T is not a Union, it returns nil (meaning standard T applies).
func (a Analyzer) getSwitchCaseOutportType(
	switchPort src.PortAddr,
	nodes map[string]src.Node,
	scope src.Scope,
	net []src.Connection,
) (*ts.Expr, *compiler.Error) {
	// 1. Get Switch Node to check if T is Union
	node, ok := nodes[switchPort.Node]
	if !ok {
		panic("switch node not found")
	}

	// Switch<T> has 1 type parameter.
	typeArg := node.TypeArgs[0]

	// Resolve T to check if it's a Union
	resolvedT, err := a.resolver.ResolveExpr(typeArg, scope)
	if err != nil {
		return nil, &compiler.Error{
			Message: err.Error(),
			Meta:    &switchPort.Meta,
		}
	}

	if resolvedT.Lit == nil || resolvedT.Lit.Union == nil {
		// T is not a Union (e.g. Switch<int>), so output type is just T.
		return &resolvedT, nil
	}

	// 2. Scan net to find the input connection to this Switch case
	// We are looking for a connection where Receiver is switchPort (case[idx])
	inputSender, ferr := a.findSenderForSwitchCaseInput(
		net,
		switchPort.Node,
		switchPort.Idx,
		switchPort,
	)
	if ferr != nil {
		return nil, ferr
	}

	// 3. Extract the Tag from the Input Sender
	tag, tagErr := a.getUnionTagFromSender(*inputSender, scope)
	if tagErr != nil {
		return nil, &compiler.Error{
			Message: "Switch case input must be a Union pattern when T passed to Switch<T> is a union",
			Meta:    &inputSender.Meta,
		}
	}

	// 4. Lookup Tag in T (Union)
	memberType, ok := resolvedT.Lit.Union[tag]
	if !ok {
		return nil, &compiler.Error{
			Message: fmt.Sprintf("Tag %q not found in union %v", tag, resolvedT),
			Meta:    &inputSender.Meta,
		}
	}

	if memberType == nil {
		// Tag-only member.
		// We return the Union Type itself
		return &resolvedT, nil
	}

	// Return the payload type.
	return memberType, nil
}

func (a Analyzer) getUnionTagFromSender(
	sender src.ConnectionSender,
	scope src.Scope,
) (string, *compiler.Error) {
	if sender.Const == nil {
		return "", &compiler.Error{
			Message: "sender is not a union constant",
			Meta:    &sender.Meta,
		}
	}

	if sender.Const.Value.Message != nil && sender.Const.Value.Message.Union != nil {
		return sender.Const.Value.Message.Union.Tag, nil
	}

	if sender.Const.Value.Ref != nil {
		unionLiteral, err := a.resolveUnionConstRef(*sender.Const.Value.Ref, scope)
		if err != nil {
			return "", err
		}
		return unionLiteral.Tag, nil
	}

	return "", &compiler.Error{
		Message: "sender is not a union constant",
		Meta:    &sender.Meta,
	}
}

func (a Analyzer) resolveUnionConstRef(
	ref core.EntityRef,
	scope src.Scope,
) (*src.UnionLiteral, *compiler.Error) {
	constant, loc, err := scope.GetConst(ref)
	if err != nil {
		return nil, &compiler.Error{
			Message: err.Error(),
			Meta:    &ref.Meta,
		}
	}

	if constant.Value.Ref != nil {
		return a.resolveUnionConstRef(*constant.Value.Ref, scope.Relocate(loc))
	}

	if constant.Value.Message == nil || constant.Value.Message.Union == nil {
		return nil, &compiler.Error{
			Message: "sender is not a union constant",
			Meta:    &ref.Meta,
		}
	}

	return constant.Value.Message.Union, nil
}

func (a Analyzer) findSenderForSwitchCaseInput(
	net []src.Connection,
	nodeName string,
	idx *uint8,
	switchPort src.PortAddr,
) (*src.ConnectionSender, *compiler.Error) {
	for _, conn := range net {
		sender, err := a.findSenderForSwitchCaseInputInConn(&conn, nodeName, idx)
		if err != nil {
			return nil, err
		}
		if sender != nil {
			return sender, nil
		}
	}

	return nil, &compiler.Error{
		Message: "sender for switch case inport not found",
		Meta:    &switchPort.Meta,
	}
}

func (a Analyzer) findSenderForSwitchCaseInputInConn(
	conn *src.Connection,
	nodeName string,
	idx *uint8,
) (*src.ConnectionSender, *compiler.Error) {
	if conn.Normal == nil {
		return nil, nil
	}

	return a.findSenderForSwitchCaseInportInConn(*conn.Normal, nodeName, idx)
}

func (a Analyzer) findSenderForSwitchCaseInportInConn(
	normConn src.NormalConnection,
	nodeName string,
	idx *uint8,
) (*src.ConnectionSender, *compiler.Error) {
	for _, receiver := range normConn.Receivers {
		// Receiver is `... -> switch:case[i]`
		if a.isSwitchCaseReceiver(receiver, nodeName, idx) {
			if len(normConn.Senders) != 1 {
				return nil, &compiler.Error{
					Message: "switch case connection must have exactly one sender when its union tag",
					Meta:    &normConn.Meta,
				}
			}
			return &normConn.Senders[0], nil
		}

		// Check chained connection
		// ... -> switch:case[i] -> ...
		if receiver.ChainedConnection != nil && receiver.ChainedConnection.Normal != nil {
			if len(receiver.ChainedConnection.Normal.Senders) > 0 {
				chainHead := receiver.ChainedConnection.Normal.Senders[0]
				if a.isSwitchCaseSender(chainHead, nodeName, idx) {
					if len(normConn.Senders) != 1 {
						return nil, &compiler.Error{
							Message: "switch case connection must have exactly one sender when its union tag",
							Meta:    &normConn.Meta,
						}
					}
					// Found the connection! (Switch case is head of chain)
					return &normConn.Senders[0], nil
				}
			}

			sender, err := a.findSenderForSwitchCaseInportInConn(
				*receiver.ChainedConnection.Normal,
				nodeName,
				idx,
			)
			if err != nil {
				return nil, err
			}
			if sender != nil {
				return sender, nil
			}
		}
	}

	return nil, nil
}

func (Analyzer) isSwitchCaseSender(sender src.ConnectionSender, nodeName string, idx *uint8) bool {
	return sender.PortAddr != nil &&
		sender.PortAddr.Node == nodeName &&
		sender.PortAddr.Port == "case" &&
		sender.PortAddr.Idx != nil &&
		*sender.PortAddr.Idx == *idx
}

func (Analyzer) isSwitchCaseReceiver(receiver src.ConnectionReceiver, nodeName string, idx *uint8) bool {
	return receiver.PortAddr != nil &&
		receiver.PortAddr.Node == nodeName &&
		receiver.PortAddr.Port == "case" &&
		receiver.PortAddr.Idx != nil &&
		*receiver.PortAddr.Idx == *idx
}

// isSwitchCasePort checks if a port address refers to a Switch component's case port.
func isSwitchCasePort(portAddr src.PortAddr, nodes map[string]src.Node) bool {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return false
	}
	return node.EntityRef.Name == "Switch" && portAddr.Port == "case"
}

package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
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
	nodeName := switchPort.Node
	idx := switchPort.Idx
	if idx == nil {
		panic("switch case port must have index")
	}

	// 1. Get Switch Node to check if T is Union
	node, ok := nodes[nodeName]
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
	inputSender, ferr := a.findSenderForSwitchCaseInput(net, nodeName, idx, switchPort)
	if ferr != nil {
		return nil, ferr
	}

	// 3. Extract the Tag from the Input Sender
	if inputSender.Union == nil {
		return nil, &compiler.Error{
			Message: "Switch case input must be a Union pattern when T passed to Switch<T> is a union",
			Meta:    &inputSender.Meta,
		}
	}

	tag := inputSender.Union.Tag

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

func (a Analyzer) findSenderForSwitchCaseInput(
	net []src.Connection,
	nodeName string,
	idx *uint8,
	switchPort src.PortAddr,
) (*src.ConnectionSender, *compiler.Error) {
	for _, conn := range net {
		if conn.Normal == nil {
			continue
		}

		for _, receiver := range conn.Normal.Receivers {
			// Receiver is `... -> switch:case[i]`
			if a.isSwitchCaseReceiver(receiver, nodeName, idx) {
				// If T is union then we expect case
				if len(conn.Normal.Senders) != 1 {
					return nil, &compiler.Error{
						Message: "switch case connection must have exactly one sender when its union tag",
						Meta:    &conn.Meta,
					}
				}
				return &conn.Normal.Senders[0], nil
			}

			// Check chained connection
			// ... -> switch:case[i] -> ...
			if receiver.ChainedConnection != nil && receiver.ChainedConnection.Normal != nil {
				chainHead := receiver.ChainedConnection.Normal.Senders[0]
				if a.isSwitchCaseSender(chainHead, nodeName, idx) {
					// Found the connection! (Switch case is head of chain)
					return &conn.Normal.Senders[0], nil
				}
			}
		}
	}

	return nil, &compiler.Error{
		Message: "sender for switch case inport not found",
		Meta:    &switchPort.Meta,
	}
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

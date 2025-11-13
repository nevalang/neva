package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

func (a Analyzer) analyzeReceivers(
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
	receiver src.ConnectionReceiver, // switch receiver
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
	nodesUsage map[string]netNodeUsage,
	analyzedSenders []src.ConnectionSender, // input senders for switch
	resolvedSwitchInputTypes []*ts.Expr, // types of input senders for switch
) ([]src.NormalConnection, []src.ConnectionReceiver, *compiler.Error) {
	analyzedSwitchConns := make([]src.NormalConnection, 0, len(receiver.Switch.Cases))

	for _, switchCaseBranch := range receiver.Switch.Cases {
		// step 1: analyze each branch as a normal connection in pattern-matching mode
		// - inside a branch, union patterns with typed members must behave as payload types
		//   (e.g. `foo::bar -> yyy` where `bar` carries `int` makes the branch sender `int`),
		//   because runtime unwraps the matched variant before forwarding to the branch receiver.
		// - this is why we pass `isPatternMatchingBranch = true` here.
		//
		// note: this step validates branch-internal sender->receiver compatibility only.
		// the global compatibility of switch inputs vs branch patterns is handled below.
		//
		// all option-senders must be subtypes of their branch-receivers
		analyzedSwitchBranch, err := a.analyzeNormalConnection(
			&switchCaseBranch,
			iface,
			nodes,
			nodesIfaces,
			scope,
			nodesUsage,
			nil,
			true, // isPatternMatchingBranch = true for switch
		)
		if err != nil {
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf("Invalid switch case: %v", err),
				Meta:    &switchCaseBranch.Meta,
			}
		}

		analyzedSwitchConns = append(analyzedSwitchConns, *analyzedSwitchBranch)

		// step 2: ensure each switch input sender is compatible with each branch pattern sender
		//
		// important: union patterns must be compared as the union type (not payload).
		// - calling `getResolvedSenderType` with isPattern=false would reject typed tag-only patterns
		//   (it demands wrapped data for typed members), which is incorrect for pattern syntax.
		// - calling it with isPattern=true would yield the payload type (e.g. `int`) which is
		//   also incorrect for matchability (we need to compare incoming `foo` against pattern `foo`,
		//   not `int`).
		// therefore we resolve union patterns to the union type manually here; for non-union senders
		// the flag is irrelevant, so we safely call the generic resolver with `false`.
		//
		// all switch branch senders must be compatible with switch input senders
		for _, branchSender := range switchCaseBranch.Senders {
			var branchSenderType ts.Expr
			if branchSender.Union != nil {
				// for pattern matching compatibility, compare against the union type itself
				typeDef, _, err := scope.GetType(branchSender.Union.EntityRef)
				if err != nil {
					return nil, nil, &compiler.Error{
						Message: fmt.Sprintf("Invalid switch case sender: failed to resolve union type: %v", err),
						Meta:    &branchSender.Meta,
					}
				}
				resolvedUnion, analyzeErr := a.analyzeTypeExpr(*typeDef.BodyExpr, scope)
				if analyzeErr != nil {
					return nil, nil, &compiler.Error{
						Message: fmt.Sprintf("Invalid switch case sender: failed to resolve union type: %v", analyzeErr),
						Meta:    &branchSender.Meta,
					}
				}
				branchSenderType = resolvedUnion
			} else {
				_, resolvedType, _, err := a.getResolvedSenderType(
					branchSender,
					iface,
					nodes,
					nodesIfaces,
					scope,
					nil,
					nodesUsage,
					false,
				)
				if err != nil {
					return nil, nil, &compiler.Error{
						Message: fmt.Sprintf("Invalid switch case sender: %v", err),
						Meta:    &branchSender.Meta,
					}
				}

				branchSenderType = resolvedType
			}

			for i, resolverSwitchInputType := range resolvedSwitchInputTypes {
				if err := a.resolver.IsSubtypeOf(*resolverSwitchInputType, branchSenderType, scope); err != nil {
					return nil, nil, &compiler.Error{
						Message: fmt.Sprintf(
							"Incompatible types in switch: %v -> %v: %v",
							analyzedSenders[i], branchSender, err.Error(),
						),
						Meta: &branchSender.Meta,
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

	analyzedDefault, err := a.analyzeReceivers(
		receiver.Switch.Default,
		scope,
		iface,
		nodes,
		nodesIfaces,
		nodesUsage,
		resolvedSwitchInputTypes,
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

	chainHeadType, err := a.getChainHeadInputType(
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

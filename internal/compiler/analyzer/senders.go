package analyzer

import (
	"github.com/nevalang/neva/internal/compiler"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
)

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
	// Some top level validation

	if sender.PortAddr == nil &&
		sender.Const == nil &&
		len(sender.StructSelector) == 0 {
		return nil, nil, &compiler.Error{
			Message: "sender is empty",
			Meta:    &sender.Meta,
		}
	}

	if sender.Const != nil && len(prevChainLink) == 0 {
		return nil, nil, &compiler.Error{
			Message: "Constants must be triggered by a signal (e.g. :start -> 42 -> ...)",
			Meta:    &sender.Meta,
		}
	}

	if len(sender.StructSelector) > 0 && len(prevChainLink) == 0 {
		return nil, nil, &compiler.Error{
			Message: "struct selectors cannot be used in non-chained connection",
			Meta:    &sender.Meta,
		}
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

func (a Analyzer) getChainHeadInputType(
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

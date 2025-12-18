package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

func (a Analyzer) analyzeSenders(
	senders []src.ConnectionSender,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
	prevChainLink []src.ConnectionSender,
	isPatternSenders bool,
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
			isPatternSenders,
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
	isPatternSender bool,
) (*src.ConnectionSender, *ts.Expr, *compiler.Error) {
	if sender.PortAddr == nil &&
		sender.Const == nil &&
		sender.Range == nil &&
		sender.Binary == nil &&

		sender.Ternary == nil &&
		sender.Union == nil &&
		len(sender.StructSelector) == 0 {
		return nil, nil, &compiler.Error{
			Message: "invalid sender",
			Meta:    &sender.Meta,
		}
	}

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
			isPatternSender,
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
			isPatternSender,
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
			isPatternSender,
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
			isPatternSender,
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
			isPatternSender,
		)
		if err != nil {
			return nil, nil, err
		}

		// check operator operand types using type system with proper union construction
		if err := a.checkOperatorOperandTypesWithTypeSystem(*sender.Binary, *leftType, *rightType, scope); err != nil {
			return nil, nil, err
		}

		// desugarer needs this information to use overloaded components
		// it could figure this out itself but it's extra work
		sender.Binary.AnalyzedType = *leftType

		resultType := a.getBinaryExprType(sender.Binary.Operator, *leftType)

		return &sender, &resultType, nil
	}

	if sender.Union != nil {
		// check that entity we are referring to is existing type definition
		typeDef, _, err := scope.GetType(sender.Union.EntityRef)
		if err != nil {
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf("failed to resolve union type: %v", err),
				Meta:    &sender.Meta,
			}
		}

		// resolve type we are referring to
		unionTypeExpr, analyzeExprErr := a.analyzeTypeExpr(*typeDef.BodyExpr, scope)
		if analyzeExprErr != nil {
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf("failed to resolve union type: %v", analyzeExprErr),
				Meta:    &sender.Meta,
			}
		}

		// check that type we refer to resolves to union
		if unionTypeExpr.Lit == nil || unionTypeExpr.Lit.Union == nil {
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf("type %v is not a union", sender.Union.EntityRef),
				Meta:    &sender.Meta,
			}
		}

		// check that tag we refer to exists in union
		memberTypeExpr, ok := unionTypeExpr.Lit.Union[sender.Union.Tag]
		if !ok {
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf("tag %q not found in union %v", sender.Union.Tag, sender.Union.EntityRef),
				Meta:    &sender.Meta,
			}
		}

		// if there's no type-expr (and thus no wrapped sender)
		// it's a tag-only union, so there's nothing to analyze
		if memberTypeExpr == nil {
			return &sender, &unionTypeExpr, nil
		}

		// Sometimes union member has type expr but union sender doesn't wrap another sender
		// This is allowed in pattern matching contexts (like switch cases), but not in other contexts
		if sender.Union.Data == nil {
			if isPatternSender {
				// in pattern matching, the switch runtime unwraps the union
				// so the type that flows to the receiver is the tag's data type
				return &sender, memberTypeExpr, nil
			}
			// if tag has type-expr and it's not pattern matching, then this union-sender must wrap another sender
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf(
					"tag %q requires a wrapped value of type %v",
					sender.Union.Tag,
					memberTypeExpr,
				),
				Meta: &sender.Meta,
			}
		}

		// analyze wrapped-sender and get its resolved type
		resolvedWrappedSender, resolvedWrappedType, analyzeWrappedErr := a.analyzeSender(
			*sender.Union.Data,
			scope,
			iface,
			nodes,
			nodesIfaces,
			nodesUsage,
			prevChainLink,
			isPatternSender,
		)
		if analyzeWrappedErr != nil {
			return nil, nil, analyzeWrappedErr
		}

		// check that wrapped-sender is sub-type of tag's constraint
		if err := a.resolver.IsSubtypeOf(*resolvedWrappedType, *memberTypeExpr, scope); err != nil {
			return nil, nil, &compiler.Error{
				Message: fmt.Sprintf(
					"wrapped sender type %v is not compatible with union tag type %v: %v",
					resolvedWrappedType,
					memberTypeExpr,
					err,
				),
				Meta: &sender.Meta,
			}
		}

		return &src.ConnectionSender{
			Union: &src.UnionSender{
				EntityRef: sender.Union.EntityRef,
				Tag:       sender.Union.Tag,
				Data:      resolvedWrappedSender,
				Meta:      sender.Union.Meta,
			},
			Meta: sender.Meta,
		}, &unionTypeExpr, nil // return type of the union, not specific tag
	}

	resolvedSenderAddr, resolvedSenderType, isSenderArr, err := a.getResolvedSenderType(
		sender,
		iface,
		nodes,
		nodesIfaces,
		scope,
		prevChainLink,
		nodesUsage,
		isPatternSender,
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

	if chainHead.Range != nil || chainHead.Const != nil {
		return ts.Expr{
			Inst: &ts.InstExpr{
				Ref: core.EntityRef{Name: "any"}, // :sig
			},
		}, nil
	}

	if len(chainHead.StructSelector) > 0 {
		return a.getStructSelectorInportType(chainHead), nil
	}

	if chainHead.Union != nil {
		typeDef, _, err := scope.GetType(chainHead.Union.EntityRef)
		if err != nil {
			return ts.Expr{}, &compiler.Error{
				Message: fmt.Sprintf("failed to resolve union type: %v", err),
				Meta:    &chainHead.Meta,
			}
		}

		unionTypeExpr, analyzeExprErr := a.analyzeTypeExpr(*typeDef.BodyExpr, scope)
		if analyzeExprErr != nil {
			return ts.Expr{}, &compiler.Error{
				Message: fmt.Sprintf("failed to resolve union type: %v", analyzeExprErr),
				Meta:    &chainHead.Meta,
			}
		}

		return unionTypeExpr, nil
	}

	return ts.Expr{}, &compiler.Error{
		Message: "Chained connection must start with port address or range expression",
		Meta:    &chainHead.Meta,
	}
}

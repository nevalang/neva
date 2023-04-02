package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/nevalang/internal/compiler"
	ts "github.com/nevalang/nevalang/pkg/types"
)

var (
	ErrReferencedMsg    = errors.New("msg not found by ref")
	ErrScopeRebase      = errors.New("scope rebase")
	ErrVecArgsLen       = errors.New("unexpected count of args for vec")
	ErrUnknownMsgType   = errors.New("unknown msg type")
	ErrUnwantedMsgField = errors.New("unwanted msg field")
	ErrMissingMsgField  = errors.New("missing msg field")
	ErrVecEl            = errors.New("vector element")
	ErrNestedMsg        = errors.New("sub message")
	ErrInassignMsg      = errors.New("msg is not assignable")
)

func (a Analyzer) analyzeMsg(
	msg compiler.Msg,
	scope Scope,
	resolvedConstr *ts.Expr,
) (compiler.Msg, map[compiler.EntityRef]struct{}, error) {
	if msg.Ref != nil {
		subMsg, err := scope.getMsg(*msg.Ref)
		if err != nil {
			return compiler.Msg{}, nil, fmt.Errorf("%w: %v", ErrReferencedMsg, err)
		}
		if msg.Ref.Pkg != "" {
			scope, err = scope.rebase(msg.Ref.Pkg)
			if err != nil {
				return compiler.Msg{}, nil, fmt.Errorf("%w: %v", ErrScopeRebase, err)
			}
		}
		resolvedSubMsg, used, err := a.analyzeMsg(subMsg, scope, resolvedConstr)
		if err != nil {
			return compiler.Msg{}, nil, fmt.Errorf("%w: %v: %v", ErrNestedMsg, err, msg.Ref)
		}
		used[*msg.Ref] = struct{}{}
		return resolvedSubMsg, used, nil // TODO do we really want unpacking here?
	}

	resolvedType, err := a.resolver.Resolve(msg.Value.TypeExpr, scope)
	if err != nil {
		return compiler.Msg{}, nil, fmt.Errorf("%w: ", err)
	}

	if resolvedConstr != nil {
		if err := a.checker.Check(resolvedType, *resolvedConstr, ts.TerminatorParams{}); err != nil {
			return compiler.Msg{}, nil, fmt.Errorf("%w: %v", ErrInassignMsg, err)
		}
	}

	msgToReturn := compiler.Msg{
		Value: compiler.MsgValue{
			TypeExpr: resolvedType,
		},
	}

	switch resolvedType.Inst.Ref {
	case "int":
		if msg.Value.Vec != nil {
			return compiler.Msg{}, nil, fmt.Errorf("%w: %v", ErrUnwantedMsgField, msg.Value.Vec)
		}
		msgToReturn.Value.Int = msg.Value.Int

	case "vec":
		if msg.Value.Int != 0 {
			return compiler.Msg{}, nil, fmt.Errorf("%w: %v", ErrUnwantedMsgField, msg.Value.Vec)
		}
		if msg.Value.Vec == nil {
			return compiler.Msg{}, nil, fmt.Errorf("%w: %v", ErrMissingMsgField, msg.Value.Vec)
		}
		if l := len(resolvedType.Inst.Args); l != 1 {
			return compiler.Msg{}, nil, fmt.Errorf("%w: want 1, got %v", ErrVecArgsLen, l)
		}
		vecType := resolvedType.Inst.Args[0]
		resolvedVec := make([]compiler.Msg, 0, len(msg.Value.Vec))
		for i, el := range msg.Value.Vec {
			analyzedEl, _, err := a.analyzeMsg(el, scope, &vecType)
			if err != nil {
				return compiler.Msg{}, nil, fmt.Errorf("%w: #%d: %v", ErrVecEl, i, err)
			}
			resolvedVec = append(resolvedVec, analyzedEl)
		}
		msgToReturn.Value.Vec = resolvedVec

	default:
		return compiler.Msg{}, nil, fmt.Errorf("%w: %v", ErrUnknownMsgType, resolvedType.Inst.Ref)
	}

	return msgToReturn, scope.visited, nil
}

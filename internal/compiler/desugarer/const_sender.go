package desugarer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var constComponentRef = src.EntityRef{
	Pkg:  "builtin",
	Name: "Const",
}

type handleConstSenderResult struct {
	desugaredConstConn src.Connection
	constNodeName      string
	constNode          src.Node
}

func (d Desugarer) handleConstSender(conn src.Connection, scope src.Scope) (handleConstSenderResult, *compiler.Error) {
	constTypeExpr, err := d.getConstType(*conn.SenderSide.ConstRef, scope)
	if err != nil {
		return handleConstSenderResult{}, compiler.Error{
			Err:      fmt.Errorf("Unable to get constant type by reference '%v'", *conn.SenderSide.ConstRef),
			Location: &scope.Location,
			Meta:     &conn.SenderSide.ConstRef.Meta,
		}.Merge(err)
	}

	constRefStr := conn.SenderSide.ConstRef.String()
	constNodeName := fmt.Sprintf("__%v__", constRefStr)
	constNode := src.Node{
		Directives: map[src.Directive][]string{
			compiler.RuntimeFuncMsgDirective: {constRefStr},
		},
		EntityRef: constComponentRef,
		TypeArgs:  []ts.Expr{constTypeExpr},
	}
	constNodeOutportAddr := src.PortAddr{
		Node: constNodeName,
		Port: "v",
	}

	return handleConstSenderResult{
		desugaredConstConn: src.Connection{
			SenderSide: src.ConnectionSenderSide{
				PortAddr:  &constNodeOutportAddr,
				Selectors: conn.SenderSide.Selectors,
				Meta:      conn.SenderSide.Meta,
			},
			ReceiverSide: conn.ReceiverSide,
			Meta:         conn.Meta,
		},
		constNodeName: constNodeName,
		constNode:     constNode,
	}, nil
}

// getConstType is needed to figure out type parameters for Const node
func (d Desugarer) getConstType(ref src.EntityRef, scope src.Scope) (ts.Expr, *compiler.Error) {
	entity, _, err := scope.Entity(ref)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &ref.Meta,
		}
	}

	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrConstSenderEntityKind, entity.Kind),
			Location: &scope.Location,
			Meta:     entity.Meta(),
		}
	}

	if entity.Const.Ref != nil {
		expr, err := d.getConstType(*entity.Const.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Location: &scope.Location,
				Meta:     entity.Meta(),
			}.Merge(err)
		}
		return expr, nil
	}

	return entity.Const.Value.TypeExpr, nil
}

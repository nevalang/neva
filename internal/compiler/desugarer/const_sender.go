package desugarer

import (
	"fmt"
	"sync/atomic"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var emitterComponentRef = src.EntityRef{
	Pkg:  "builtin",
	Name: "Emitter",
}

type handleLiteralSenderResult struct {
	handleConstRefSenderResult // conceptually incorrrect but convenient to reuse
	// constant                   src.Const
	constName string
}

type handleConstRefSenderResult struct {
	desugaredConn   src.Connection
	emitterNodeName string
	emitterNode     src.Node
}

// In the future compiler can operate in concurrently
var litSendersCount atomic.Uint32

func (d Desugarer) handleLiteralSender(
	conn src.Connection,
	scope src.Scope,
) (
	handleLiteralSenderResult,
	*compiler.Error,
) {
	counter := litSendersCount.Load()
	litSendersCount.Store(counter + 1)
	constName := fmt.Sprintf("literal-%d", counter)

	// we can't call d.handleConstRefSender()
	// because our virtual const isn't in the scope

	emitterNodeName := "$" + constName
	emitterNode := src.Node{
		Directives: map[src.Directive][]string{
			compiler.BindDirective: {constName},
		},
		EntityRef: emitterComponentRef,
		TypeArgs: []ts.Expr{
			conn.
				SenderSide.
				Const.
				Value.
				TypeExpr,
		},
	}
	emitterNodeOutportAddr := src.PortAddr{
		Node: emitterNodeName,
		Port: "msg",
	}

	return handleLiteralSenderResult{
		constName: constName,
		handleConstRefSenderResult: handleConstRefSenderResult{
			desugaredConn: src.Connection{
				SenderSide: src.ConnectionSenderSide{
					PortAddr:  &emitterNodeOutportAddr,
					Selectors: conn.SenderSide.Selectors,
					Meta:      conn.SenderSide.Meta,
				},
				ReceiverSide: conn.ReceiverSide,
				Meta:         conn.Meta,
			},
			emitterNodeName: emitterNodeName,
			emitterNode:     emitterNode,
		},
	}, nil
}

func (d Desugarer) handleConstRefSender(
	conn src.Connection,
	scope src.Scope,
) (
	handleConstRefSenderResult,
	*compiler.Error,
) {
	constTypeExpr, err := d.getConstTypeByRef(*conn.SenderSide.Const.Ref, scope)
	if err != nil {
		return handleConstRefSenderResult{}, compiler.Error{
			Err: fmt.Errorf(
				"Unable to get constant type by reference '%v'",
				*conn.SenderSide.Const.Ref,
			),
			Location: &scope.Location,
			Meta:     &conn.SenderSide.Const.Ref.Meta,
		}.Merge(err)
	}

	constRefStr := conn.SenderSide.Const.Ref.String()

	emitterNodeName := "$" + constRefStr
	emitterNode := src.Node{
		Directives: map[src.Directive][]string{
			compiler.BindDirective: {constRefStr},
		},
		EntityRef: emitterComponentRef,
		TypeArgs:  []ts.Expr{constTypeExpr},
	}
	emitterNodeOutportAddr := src.PortAddr{
		Node: emitterNodeName,
		Port: "msg",
	}

	return handleConstRefSenderResult{
		desugaredConn: src.Connection{
			SenderSide: src.ConnectionSenderSide{
				PortAddr:  &emitterNodeOutportAddr,
				Selectors: conn.SenderSide.Selectors,
				Meta:      conn.SenderSide.Meta,
			},
			ReceiverSide: conn.ReceiverSide,
			Meta:         conn.Meta,
		},
		emitterNodeName: emitterNodeName,
		emitterNode:     emitterNode,
	}, nil
}

// getConstTypeByRef is needed to figure out type parameters for Const node
func (d Desugarer) getConstTypeByRef(ref src.EntityRef, scope src.Scope) (ts.Expr, *compiler.Error) {
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
		expr, err := d.getConstTypeByRef(*entity.Const.Ref, scope)
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

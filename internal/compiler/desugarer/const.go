package desugarer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

// handleConst inserts nodes and connections and returns resulted network (because new array could be allocated).
func (d Desugarer) handleConst(
	conn src.Connection,
	scope src.Scope,
	desugaredNodes map[string]src.Node,
	desugaredNet []src.Connection,
) ([]src.Connection, *compiler.Error) {
	constTypeExpr, err := d.getConstType(*conn.SenderSide.ConstRef, scope)
	if err != nil {
		return nil, compiler.Error{
			Err:      fmt.Errorf("Unable to get constant type by reference '%v'", *conn.SenderSide.ConstRef),
			Location: &scope.Location,
			Meta:     &conn.SenderSide.ConstRef.Meta,
		}.Merge(err)
	}

	constRefStr := conn.SenderSide.ConstRef.String()
	constNodeName := fmt.Sprintf("__%v__", constRefStr)

	desugaredNodes[constNodeName] = src.Node{
		Directives: map[src.Directive][]string{
			compiler.RuntimeFuncMsgDirective: {constRefStr},
		},
		EntityRef: src.EntityRef{
			Pkg:  "builtin",
			Name: "Const",
		},
		TypeArgs: []ts.Expr{constTypeExpr},
	}

	constNodeOutportAddr := src.PortAddr{
		Node: constNodeName,
		Port: "v",
	}

	desugaredNet = append(desugaredNet, src.Connection{
		SenderSide: src.SenderConnectionSide{
			PortAddr:  &constNodeOutportAddr,
			Selectors: conn.SenderSide.Selectors,
			Meta:      conn.SenderSide.Meta,
		},
		ReceiverSides: conn.ReceiverSides,
		Meta:          conn.Meta,
	})

	return desugaredNet, nil
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

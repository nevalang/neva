package desugarer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

var emitterFlowRef = core.EntityRef{
	Pkg:  "builtin",
	Name: "New",
}

// In the future compiler can operate in concurrently
var (
	virtualEmittersCount uint64
	virtualConstCount    uint64
)

func (d Desugarer) handleLiteralSender(
	constant src.Const,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (src.PortAddr, *compiler.Error) {
	virtualConstCount++
	constName := fmt.Sprintf("__const__%d", virtualConstCount)

	// we can't call d.handleConstRefSender()
	// because our virtual const isn't in the scope

	virtualEmittersCount++
	emitterNodeName := fmt.Sprintf("__new__%d", virtualEmittersCount)

	emitterNode := src.Node{
		Directives: map[src.Directive][]string{
			compiler.BindDirective: {constName},
		},
		EntityRef: emitterFlowRef,
		TypeArgs:  []ts.Expr{constant.TypeExpr},
	}

	nodesToInsert[emitterNodeName] = emitterNode
	constsToInsert[constName] = constant

	emitterNodeOutportAddr := src.PortAddr{
		Node: emitterNodeName,
		Port: "msg",
	}

	return emitterNodeOutportAddr, nil
}

func (d Desugarer) handleConstRefSender(
	ref core.EntityRef,
	nodesToInsert map[string]src.Node,
	scope src.Scope,
) (src.PortAddr, *compiler.Error) {
	constTypeExpr, err := d.getConstTypeByRef(ref, scope)
	if err != nil {
		return src.PortAddr{}, compiler.Error{
			Message:  fmt.Sprintf("Unable to get constant type by reference '%v'", ref),
			Location: &scope.Location,
			Meta:     &ref.Meta,
		}.Wrap(err)
	}

	virtualEmittersCount++
	virtualEmitterName := fmt.Sprintf("__new__%d", virtualEmittersCount)

	emitterNode := src.Node{
		// don't forget to bind
		Directives: map[src.Directive][]string{
			compiler.BindDirective: {ref.String()},
		},
		EntityRef: emitterFlowRef,
		TypeArgs:  []ts.Expr{constTypeExpr},
	}

	emitterNodeOutportAddr := src.PortAddr{
		Node: virtualEmitterName,
		Port: "msg",
	}

	nodesToInsert[virtualEmitterName] = emitterNode

	return emitterNodeOutportAddr, nil
}

// getConstTypeByRef is needed to figure out type parameters for Const node
func (d Desugarer) getConstTypeByRef(ref core.EntityRef, scope src.Scope) (ts.Expr, *compiler.Error) {
	entity, _, err := scope.Entity(ref)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message:  err.Error(),
			Location: &scope.Location,
			Meta:     &ref.Meta,
		}
	}

	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, &compiler.Error{
			Message: fmt.Sprintf(
				"Entity that is used as a const reference in flow's network must be of kind constant: %v",
				entity.Kind,
			),
			Location: &scope.Location,
			Meta:     entity.Meta(),
		}
	}

	if entity.Const.Value.Ref != nil {
		expr, err := d.getConstTypeByRef(*entity.Const.Value.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Location: &scope.Location,
				Meta:     entity.Meta(),
			}.Wrap(err)
		}
		return expr, nil
	}

	return entity.Const.TypeExpr, nil
}

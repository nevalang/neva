package irgen

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler/ir"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	src "github.com/nevalang/neva/pkg/ast"
)

//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (Generator) getFuncRef(versions []src.Component, node src.Node) (string, src.Component) {
	var version src.Component
	if len(versions) == 1 {
		version = versions[0]
	} else {
		version = versions[*node.OverloadIndex]
	}

	extern, hasExtern := version.Directives.Find(src.DirectiveKindExtern)
	if !hasExtern {
		return "", version
	}

	return extern.Extern.Ref, version
}

// getConfigMsg lowers the one static configuration producer of an extern call.
// A node-level #bind supplies an ordinary constant. A component-level
// #bind_type supplies the fully resolved TypeExpr selected for this node as a
// std/reflect.Type-shaped message. The analyzer rejects the combination, so a
// runtime function never has two competing configuration values.
//
//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (g Generator) getConfigMsg(
	node src.Node,
	component src.Component,
	callerScope src.Scope,
	componentScope src.Scope,
) (*ir.Message, error) {
	bind, hasBind := node.Directives.Find(src.DirectiveKindBind)
	bindType, hasBindType := component.Directives.Find(src.DirectiveKindBindType)

	switch {
	case hasBind:
		entity, location, err := callerScope.Entity(bind.Bind.ConstRef)
		if err != nil {
			//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			return nil, err
		}

		return getIRMsgBySrcRef(
			entity.Const.Value,
			callerScope.Relocate(location),
			entity.Const.TypeExpr,
		)
	case hasBindType:
		resolvedType, err := g.resolveBoundType(
			bindType.BindType.TypeExpr,
			component,
			node,
			componentScope,
		)
		if err != nil {
			return nil, err
		}
		return reflectTypeMessage(resolvedType, g.resolver, componentScope)
	default:
		//nolint:nilnil // nil config is expected when no directive produces it
		return nil, nil
	}
}

// resolveBoundType substitutes the extern component's resolved node arguments
// into its #bind_type expression. The analyzer has already proved that the
// expression is valid in the declaration scope; reaching an unresolved type
// here means a compiler cross-stage invariant was violated.
//
//nolint:gocritic // Compiler pipeline values are passed by value consistently.
func (g Generator) resolveBoundType(
	expr ts.Expr,
	component src.Component,
	node src.Node,
	scope src.Scope,
) (ts.Expr, error) {
	params := component.TypeParams.Params
	frame := make(map[string]ts.Def, len(params))
	for index := range node.TypeArgs {
		if index >= len(params) {
			return ts.Expr{}, fmt.Errorf(
				"#bind_type component has %d declared type parameters, got argument %d",
				len(params),
				index+1,
			)
		}
		argument := node.TypeArgs[index]
		frame[params[index].Name] = ts.Def{BodyExpr: &argument}
	}

	resolved, err := g.resolver.ResolveExprWithFrame(expr, frame, scope)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("resolve #bind_type expression: %w", err)
	}
	return resolved, nil
}

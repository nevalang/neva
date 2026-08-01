package irgen

import (
	"context"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/ir"
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

	extern, hasExtern := version.Directives.Find(src.ExternDirective)
	if !hasExtern {
		return "", version
	}

	return extern.Identifier.Value, version
}

//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func getConfigMsg(node src.Node, scope src.Scope) (*ir.Message, error) {
	bind, hasBind := node.Directives.Find(src.BindDirective)
	if !hasBind {
		//nolint:nilnil // nil config is expected when no bind directive is present
		return nil, nil
	}

	entityRef, err := compiler.ParseEntityRef(context.Background(), bind.Identifier.Value)
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	entity, location, err := scope.Entity(entityRef)
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return getIRMsgBySrcRef(
		entity.Const.Value,
		scope.Relocate(location),
		entity.Const.TypeExpr,
	)
}

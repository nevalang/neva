package irgen

import (
	"context"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/ir"
	src "github.com/nevalang/neva/pkg/ast"
)

func (Generator) getFuncRef(versions []src.Component, node src.Node) (string, src.Component) {
	var version src.Component
	if len(versions) == 1 {
		version = versions[0]
	} else {
		version = versions[*node.OverloadIndex]
	}

	externArg, hasExtern := version.Directives[compiler.ExternDirective]
	if !hasExtern {
		return "", version
	}

	return externArg, version
}

func getConfigMsg(node src.Node, scope src.Scope) (*ir.Message, error) {
	bindArg, hasBind := node.Directives[compiler.BindDirective]
	if !hasBind {
		//nolint:nilnil // nil config is expected when no bind directive is present
		return nil, nil
	}

	entityRef, err := compiler.ParseEntityRef(context.Background(), bindArg)
	if err != nil {
		return nil, err
	}

	entity, location, err := scope.Entity(entityRef)
	if err != nil {
		return nil, err
	}

	return getIRMsgBySrcRef(
		entity.Const.Value,
		scope.Relocate(location),
		entity.Const.TypeExpr,
	)
}

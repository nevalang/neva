package irgen

import (
	"context"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ir"
)

func (Generator) getFuncRef(versions []src.Component, node src.Node) (string, src.Component, error) {
	var version src.Component
	if len(versions) == 1 {
		version = versions[0]
	} else {
		version = versions[*node.OverloadIndex]
	}

	externArg, hasExtern := version.Directives[compiler.ExternDirective]
	if !hasExtern {
		return "", version, nil
	}

	return externArg, version, nil
}

func getConfigMsg(node src.Node, scope src.Scope) (*ir.Message, error) {
	bindArg, hasBind := node.Directives[compiler.BindDirective]
	if !hasBind {
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

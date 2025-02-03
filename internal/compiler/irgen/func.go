package irgen

import (
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/ir"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
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
		return "", src.Component{}, nil
	}

	return externArg, version, nil
}

func getConfigMsg(node src.Node, scope src.Scope) (*ir.Message, error) {
	bindArg, hasBind := node.Directives[compiler.BindDirective]
	if !hasBind {
		return nil, nil
	}

	entity, location, err := scope.Entity(compiler.ParseEntityRef(bindArg))
	if err != nil {
		return nil, err
	}

	return getIRMsgBySrcRef(
		entity.Const.Value,
		scope.Relocate(location),
		entity.Const.TypeExpr,
	)
}

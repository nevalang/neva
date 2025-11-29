package irgen

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/internal/compiler/utils/generated"
	"github.com/nevalang/neva/internal/runtime"
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

	// Call the generated Neva function
	out, err := generated.ParseEntityRef(context.Background(), generated.ParseEntityRefInput{Ref: bindArg})
	if err != nil {
		return nil, err
	}

	// Unmarshal the result
	msg, ok := out.Res.(runtime.StructMsg)
	if !ok {
		return nil, fmt.Errorf("expected struct msg, got %T", out.Res)
	}
	entityRef := core.EntityRef{
		Pkg:  msg.Get("pkg").Str(),
		Name: msg.Get("name").Str(),
		Meta: core.Meta{Text: msg.Get("metaText").Str()},
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

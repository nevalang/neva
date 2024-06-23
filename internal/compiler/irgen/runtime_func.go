package irgen

import (
	"errors"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	"github.com/nevalang/neva/internal/runtime/ir"
)

// scope must contain location where node found, not function
func (g Generator) getFuncCall(
	nodeCtx nodeContext,
	scope src.Scope,
	funcRef string,
) (ir.FuncCall, *compiler.Error) {
	cfgMsg, err := getCfgMsg(nodeCtx.node, scope)
	if err != nil {
		return ir.FuncCall{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
		}
	}
	return ir.FuncCall{
		Ref: funcRef,
		IO: ir.FuncIO{
			In:  g.getFuncInports(nodeCtx),
			Out: g.getFuncOutports(nodeCtx),
		},
		Msg: cfgMsg,
	}, nil
}

func getFuncRef(flow src.Component, nodeTypeArgs []ts.Expr) (string, error) {
	args, ok := flow.Directives[compiler.ExternDirective]
	if !ok {
		return "", nil
	}

	if len(args) == 1 {
		return args[0], nil
	}

	if len(nodeTypeArgs) == 0 || nodeTypeArgs[0].Inst == nil {
		// FIXME sometimes we have union here
		// we must use node argument instead of flow type param
		return "", nil
	}

	firstTypeArg := nodeTypeArgs[0].Inst.Ref.String()
	for _, arg := range args {
		parts := strings.Split(arg, " ")
		if firstTypeArg == parts[0] {
			return parts[1], nil
		}
	}

	return "", errors.New("type argument mismatches runtime func directive")
}

func getCfgMsg(node src.Node, scope src.Scope) (*ir.Message, *compiler.Error) {
	args, ok := node.Directives[compiler.BindDirective]
	if !ok {
		return nil, nil
	}

	entity, location, err := scope.Entity(compiler.ParseEntityRef(args[0]))
	if err != nil {
		return nil, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
		}
	}

	return getIRMsgBySrcRef(entity.Const, scope.WithLocation(location))
}

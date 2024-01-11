package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

//nolint:lll
var (
	ErrStructInportsArgNonStruct               = errors.New("Type argument for component with struct inports directive must be struct")
	ErrStructInportsNodeTypeArgsCount          = errors.New("Note that uses component with struct inports directive must pass exactly one type argument")
	ErrStructInportsTypeParamConstr            = errors.New("Component that uses struct inports directive must have type parameter with struct constraint")
	ErrStructInportsTypeParamsCount            = errors.New("Component that uses struct inports directive must have type parameter with have exactly one type parameter")
	ErrNormalInportsWithStructInportsDirective = errors.New("Component that uses struct inports directive must have no defined inports")
)

func (a Analyzer) analyzeComponentNodes(
	nodes map[string]src.Node,
	scope src.Scope,
) (map[string]src.Node, map[string]src.Interface, *compiler.Error) {
	analyzedNodes := make(map[string]src.Node, len(nodes))
	nodesInterfaces := make(map[string]src.Interface, len(nodes))

	for nodeName, node := range nodes {
		analyzedNode, nodeInterface, err := a.analyzeComponentNode(node, scope)
		if err != nil {
			return nil, nil, compiler.Error{
				Err:      fmt.Errorf("Invalid node '%v %v", nodeName, node),
				Location: &scope.Location,
				Meta:     &node.Meta,
			}.Merge(err)
		}

		nodesInterfaces[nodeName] = nodeInterface
		analyzedNodes[nodeName] = analyzedNode
	}

	return analyzedNodes, nodesInterfaces, nil
}

//nolint:funlen
func (a Analyzer) analyzeComponentNode(node src.Node, scope src.Scope) (src.Node, src.Interface, *compiler.Error) {
	entity, location, err := scope.Entity(node.EntityRef)
	if err != nil {
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &node.Meta,
		}
	}

	if entity.Kind != src.ComponentEntity && entity.Kind != src.InterfaceEntity {
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrNodeWrongEntity, entity.Kind),
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	runtimeMsgArgs, hasRuntimeMsg := node.Directives[compiler.RuntimeFuncMsgDirective]
	if hasRuntimeMsg && len(runtimeMsgArgs) != 1 {
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err:      ErrRuntimeMsgArgs,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	iface, aerr := a.getResolvedNodeInterface(entity, hasRuntimeMsg, location, node, scope)
	if aerr != nil {
		return src.Node{}, src.Interface{}, aerr
	}

	if len(node.TypeArgs) != len(iface.TypeParams.Params) {
		var err error
		if len(node.TypeArgs) < len(iface.TypeParams.Params) {
			err = ErrNodeTypeArgsMissing
		} else {
			err = ErrNodeTypeArgsTooMuch
		}
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err: fmt.Errorf(
				"%w: want %v, got %v",
				err, iface.TypeParams, node.TypeArgs,
			),
			Location: &location,
			Meta:     &node.Meta,
		}
	}

	resolvedArgs, _, err := a.resolver.ResolveFrame(node.TypeArgs, iface.TypeParams.Params, scope)
	if err != nil {
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err:      err,
			Location: &location,
			Meta:     &node.Meta,
		}
	}

	if node.Deps == nil {
		return src.Node{
			Directives: node.Directives,
			EntityRef:  node.EntityRef,
			TypeArgs:   resolvedArgs,
			Meta:       node.Meta,
		}, iface, nil
	}

	resolvedComponentDI := make(map[string]src.Node, len(node.Deps))
	for depName, depNode := range node.Deps {
		resolvedDep, _, err := a.analyzeComponentNode(depNode, scope)
		if err != nil {
			return src.Node{}, src.Interface{}, compiler.Error{
				Err:      fmt.Errorf("Invalid node dependency: node '%v'", depNode),
				Location: &location,
				Meta:     &depNode.Meta,
			}.Merge(err)
		}
		resolvedComponentDI[depName] = resolvedDep
	}

	return src.Node{
		Directives: node.Directives,
		EntityRef:  node.EntityRef,
		TypeArgs:   resolvedArgs,
		Deps:       resolvedComponentDI,
		Meta:       node.Meta,
	}, iface, nil
}

func (a Analyzer) getResolvedNodeInterface( //nolint:funlen
	entity src.Entity,
	hasRuntimeMsg bool,
	location src.Location,
	node src.Node,
	scope src.Scope,
) (src.Interface, *compiler.Error) {
	if entity.Kind == src.InterfaceEntity {
		if hasRuntimeMsg {
			return src.Interface{}, &compiler.Error{
				Err:      ErrInterfaceNodeWithRuntimeMsg,
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		if node.Deps != nil {
			return src.Interface{}, &compiler.Error{
				Err:      ErrNonComponentNodeWithDI,
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		return entity.Interface, nil
	}

	runtimeFuncArgs, isRuntimeFunc := entity.Component.Directives[compiler.RuntimeFuncDirective]

	if hasRuntimeMsg && !isRuntimeFunc {
		return src.Interface{}, &compiler.Error{
			Err:      ErrNormNodeRuntimeMsg,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if len(runtimeFuncArgs) > 1 && len(node.TypeArgs) != 1 {
		return src.Interface{}, &compiler.Error{
			Err:      ErrRuntimeFuncOverloadingNodeArgs,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	iface := entity.Component.Interface

	_, hasStructInportsDirective := entity.Component.Directives[compiler.StructInports]

	if !hasStructInportsDirective {
		return iface, nil
	}

	// handle struct inports directive case:

	if len(iface.IO.In) != 0 {
		return src.Interface{}, &compiler.Error{
			Err:      ErrNormalInportsWithStructInportsDirective,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if len(iface.TypeParams.Params) != 1 {
		return src.Interface{}, &compiler.Error{
			Err:      ErrStructInportsTypeParamsCount,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	resolvedTypeParamConstr, err := a.resolver.ResolveExpr(*iface.TypeParams.Params[0].Constr, scope)
	if err != nil {
		return src.Interface{}, &compiler.Error{
			Err:      err,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if resolvedTypeParamConstr.Lit == nil || resolvedTypeParamConstr.Lit.Struct == nil {
		return src.Interface{}, &compiler.Error{
			Err:      ErrStructInportsTypeParamConstr,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if len(node.TypeArgs) != 1 {
		return src.Interface{}, &compiler.Error{
			Err:      ErrStructInportsNodeTypeArgsCount,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	resolvedNodeArg, err := a.resolver.ResolveExpr(node.TypeArgs[0], scope)
	if err != nil {
		return src.Interface{}, &compiler.Error{
			Err:      err,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if resolvedNodeArg.Lit == nil || resolvedNodeArg.Lit.Struct == nil {
		return src.Interface{}, &compiler.Error{
			Err:      ErrStructInportsArgNonStruct,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	structFields := resolvedNodeArg.Lit.Struct
	inports := make(map[string]src.Port, len(structFields))
	for fieldName, fieldTypeExpr := range structFields {
		inports[fieldName] = src.Port{
			TypeExpr: fieldTypeExpr,
		}
	}

	iface = src.Interface{
		TypeParams: iface.TypeParams,
		IO: src.IO{
			In:  inports,
			Out: iface.IO.Out,
		},
		Meta: iface.Meta,
	}

	return iface, nil
}

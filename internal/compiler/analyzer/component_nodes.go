package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

//nolint:lll
var (
	ErrAutoPortsArgNonStruct               = errors.New("Type argument for component with struct inports directive must be struct")
	ErrAutoPortsNodeTypeArgsCount          = errors.New("Note that uses component with struct inports directive must pass exactly one type argument")
	ErrAutoPortsTypeParamConstr            = errors.New("Component that uses struct inports directive must have type parameter with struct constraint")
	ErrAutoPortsTypeParamsCount            = errors.New("Component that uses struct inports directive must have type parameter with have exactly one type parameter")
	ErrNormalInportsWithAutoPortsDirective = errors.New("Component that uses struct inports directive must have no defined inports")
)

func (a Analyzer) analyzeComponentNodes(
	parentTypeParams src.TypeParams,
	nodes map[string]src.Node,
	scope src.Scope,
) (
	map[string]src.Node, // resolved nodes
	map[string]src.Interface, // resolved nodes interfaces
	*compiler.Error, // err
) {
	analyzedNodes := make(map[string]src.Node, len(nodes))
	nodesInterfaces := make(map[string]src.Interface, len(nodes))

	for nodeName, node := range nodes {
		analyzedNode, nodeInterface, err := a.analyzeComponentNode(node, parentTypeParams, scope)
		if err != nil {
			return nil, nil, compiler.Error{
				Location: &scope.Location,
				Meta:     &node.Meta,
			}.Wrap(err)
		}

		nodesInterfaces[nodeName] = nodeInterface
		analyzedNodes[nodeName] = analyzedNode
	}

	return analyzedNodes, nodesInterfaces, nil
}

//nolint:funlen
func (a Analyzer) analyzeComponentNode(
	node src.Node,
	parentTypeParams src.TypeParams,
	scope src.Scope,
) (src.Node, src.Interface, *compiler.Error) {
	nodeEntity, location, err := scope.Entity(node.EntityRef)
	if err != nil {
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &node.Meta,
		}
	}

	if nodeEntity.Kind != src.ComponentEntity &&
		nodeEntity.Kind != src.InterfaceEntity {
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrNodeWrongEntity, nodeEntity.Kind),
			Location: &location,
			Meta:     nodeEntity.Meta(),
		}
	}

	bindDirectiveArgs, usesBindDirective := node.Directives[compiler.BindDirective]
	if usesBindDirective && len(bindDirectiveArgs) != 1 {
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err:      ErrBindDirectiveArgs,
			Location: &location,
			Meta:     nodeEntity.Meta(),
		}
	}

	nodeIface, aerr := a.getNodeInterface(
		nodeEntity,
		usesBindDirective,
		location,
		node,
		scope,
	)
	if aerr != nil {
		return src.Node{}, src.Interface{}, aerr
	}

	// We need to get resolved frame from parent type parameters
	// in order to be able to resolve node's args
	// since they can refer to type parameter of the parent (interface)
	_, resolvedParentParamsFrame, err := a.resolver.ResolveParams(
		parentTypeParams.Params,
		scope,
	)
	if err != nil {
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err:      err,
			Location: &location,
			Meta:     &node.Meta,
		}
	}

	// Now when we have frame made of parent type parameters constraints
	// we can resolve cases like `subnode SubComponent<T>`
	// where `T` refers to type parameter of the component/interface we're in.
	resolvedNodeArgs, err := a.resolver.ResolveExprsWithFrame(
		node.TypeArgs,
		resolvedParentParamsFrame,
		scope,
	)
	if err != nil {
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err:      err,
			Location: &location,
			Meta:     &node.Meta,
		}
	}

	// Finally check that every argument is compatible
	// with corresponding parameter of the node's interface.
	if err = a.resolver.CheckArgsCompatibility(
		resolvedNodeArgs,
		nodeIface.TypeParams.Params,
		scope,
	); err != nil {
		return src.Node{}, src.Interface{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &node.Meta,
		}
	}

	if node.Deps == nil {
		return src.Node{
			Directives: node.Directives,
			EntityRef:  node.EntityRef,
			TypeArgs:   resolvedNodeArgs,
			Meta:       node.Meta,
		}, nodeIface, nil
	}

	resolvedComponentDI := make(map[string]src.Node, len(node.Deps))
	for depName, depNode := range node.Deps {
		resolvedDep, _, err := a.analyzeComponentNode(depNode, parentTypeParams, scope)
		if err != nil {
			return src.Node{}, src.Interface{}, compiler.Error{
				Location: &location,
				Meta:     &depNode.Meta,
			}.Wrap(err)
		}
		resolvedComponentDI[depName] = resolvedDep
	}

	return src.Node{
		Directives: node.Directives,
		EntityRef:  node.EntityRef,
		TypeArgs:   resolvedNodeArgs,
		Deps:       resolvedComponentDI,
		Meta:       node.Meta,
	}, nodeIface, nil
}

func (a Analyzer) getNodeInterface( //nolint:funlen
	entity src.Entity,
	hasConfigMsg bool,
	location src.Location,
	node src.Node,
	scope src.Scope,
) (src.Interface, *compiler.Error) {
	if entity.Kind == src.InterfaceEntity {
		if hasConfigMsg {
			return src.Interface{}, &compiler.Error{
				Err:      ErrInterfaceNodeBindDirective,
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

	externArgs, hasExternDirective := entity.Component.Directives[compiler.ExternDirective]

	if hasConfigMsg && !hasExternDirective {
		return src.Interface{}, &compiler.Error{
			Err:      ErrNormNodeBind,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if len(externArgs) > 1 && len(node.TypeArgs) != 1 {
		return src.Interface{}, &compiler.Error{
			Err:      ErrExternOverloadingNodeArgs,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	iface := entity.Component.Interface

	_, hasAutoPortsDirective := entity.Component.Directives[compiler.AutoportsDirective]
	if !hasAutoPortsDirective {
		return iface, nil
	}

	// if we here then we have #autoports

	if len(iface.IO.In) != 0 {
		return src.Interface{}, &compiler.Error{
			Err:      ErrNormalInportsWithAutoPortsDirective,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if len(iface.TypeParams.Params) != 1 {
		return src.Interface{}, &compiler.Error{
			Err:      ErrAutoPortsTypeParamsCount,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	resolvedTypeParamConstr, err := a.resolver.ResolveExpr(iface.TypeParams.Params[0].Constr, scope)
	if err != nil {
		return src.Interface{}, &compiler.Error{
			Err:      err,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if resolvedTypeParamConstr.Lit == nil || resolvedTypeParamConstr.Lit.Struct == nil {
		return src.Interface{}, &compiler.Error{
			Err:      ErrAutoPortsTypeParamConstr,
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if len(node.TypeArgs) != 1 {
		return src.Interface{}, &compiler.Error{
			Err:      ErrAutoPortsNodeTypeArgsCount,
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
			Err:      ErrAutoPortsArgNonStruct,
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

	return src.Interface{
		TypeParams: iface.TypeParams,
		IO: src.IO{
			In: inports,
			Out: map[string]src.Port{
				"msg": {
					TypeExpr: resolvedNodeArg,
					IsArray:  false,
					Meta:     iface.IO.Out["v"].Meta,
				},
			},
		},
		Meta: iface.Meta,
	}, nil
}

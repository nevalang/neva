package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	"github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type foundInterface struct {
	iface    src.Interface
	location src.Location
}

func (a Analyzer) analyzeNodes(
	flowIface src.Interface,
	nodes map[string]src.Node,
	scope src.Scope,
) (
	map[string]src.Node, // resolved nodes
	map[string]foundInterface, // resolved nodes interfaces with locations
	bool, // one of the nodes has error guard
	*compiler.Error, // err
) {
	analyzedNodes := make(map[string]src.Node, len(nodes))
	nodesInterfaces := make(map[string]foundInterface, len(nodes))
	hasErrGuard := false

	for nodeName, node := range nodes {
		if node.ErrGuard {
			hasErrGuard = true
		}

		analyzedNode, nodeInterface, err := a.analyzeNode(
			flowIface,
			node,
			scope,
		)
		if err != nil {
			return nil, nil, false, compiler.Error{
				Location: &scope.Location,
				Meta:     &node.Meta,
			}.Wrap(err)
		}

		nodesInterfaces[nodeName] = nodeInterface
		analyzedNodes[nodeName] = analyzedNode
	}

	return analyzedNodes, nodesInterfaces, hasErrGuard, nil
}

func (a Analyzer) analyzeNode(
	flowIface src.Interface,
	node src.Node,
	scope src.Scope,
) (src.Node, foundInterface, *compiler.Error) {
	parentTypeParams := flowIface.TypeParams

	nodeEntity, location, err := scope.Entity(node.EntityRef)
	if err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message:  err.Error(),
			Location: &scope.Location,
			Meta:     &node.Meta,
		}
	}

	if nodeEntity.Kind != src.ComponentEntity &&
		nodeEntity.Kind != src.InterfaceEntity {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message:  fmt.Sprintf("Node can only refer to flows or interfaces: %v", nodeEntity.Kind),
			Location: &location,
			Meta:     nodeEntity.Meta(),
		}
	}

	bindDirectiveArgs, usesBindDirective := node.Directives[compiler.BindDirective]
	if usesBindDirective && len(bindDirectiveArgs) != 1 {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message:  "Node with #bind directive must provide exactly one argument",
			Location: &location,
			Meta:     nodeEntity.Meta(),
		}
	}

	// We need to get resolved frame from parent type parameters
	// in order to be able to resolve node's args
	// since they can refer to type parameter of the parent (interface)
	_, resolvedParentParamsFrame, err := a.resolver.ResolveParams(
		parentTypeParams.Params,
		scope,
	)
	if err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message:  err.Error(),
			Location: &location,
			Meta:     &node.Meta,
		}
	}

	// Now when we have frame made of parent type parameters constraints
	// we can resolve cases like `subnode SubFlow<T>`
	// where `T` refers to type parameter of the component/interface we're in.
	resolvedNodeArgs, err := a.resolver.ResolveExprsWithFrame(
		node.TypeArgs,
		resolvedParentParamsFrame,
		scope,
	)
	if err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message:  err.Error(),
			Location: &location,
			Meta:     &node.Meta,
		}
	}

	nodeIface, aerr := a.getNodeInterface(
		nodeEntity,
		usesBindDirective,
		location,
		node,
		scope,
		resolvedNodeArgs,
	)
	if aerr != nil {
		return src.Node{}, foundInterface{}, aerr
	}

	if node.ErrGuard {
		if _, ok := flowIface.IO.Out["err"]; !ok {
			return src.Node{}, foundInterface{}, &compiler.Error{
				Message:  "Guard is not allowed for nodes without 'err' output",
				Location: &scope.Location,
				Meta:     &node.Meta,
			}
		}
		if _, ok := nodeIface.IO.Out["err"]; !ok {
			return src.Node{}, foundInterface{}, &compiler.Error{
				Message:  "Guard is not allowed for nodes without ':err' output",
				Location: &scope.Location,
				Meta:     &node.Meta,
			}
		}
	}

	// default any
	if len(resolvedNodeArgs) == 0 && len(nodeIface.TypeParams.Params) == 1 {
		resolvedNodeArgs = []typesystem.Expr{
			{
				Inst: &typesystem.InstExpr{
					Ref: core.EntityRef{Name: "any"},
				},
			},
		}
	}

	// Finally check that every argument is compatible
	// with corresponding parameter of the node's interface.
	if err = a.resolver.CheckArgsCompatibility(
		resolvedNodeArgs,
		nodeIface.TypeParams.Params,
		scope,
	); err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message:  err.Error(),
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
				ErrGuard:   node.ErrGuard,
			}, foundInterface{
				iface:    nodeIface,
				location: location,
			}, nil
	}

	// TODO probably here
	// implement interface->component subtyping
	// in a way where FP possible

	resolvedFlowDI := make(map[string]src.Node, len(node.Deps))
	for depName, depNode := range node.Deps {
		resolvedDep, _, err := a.analyzeNode(
			flowIface,
			depNode,
			scope,
		)
		if err != nil {
			return src.Node{}, foundInterface{}, compiler.Error{
				Location: &location,
				Meta:     &depNode.Meta,
			}.Wrap(err)
		}
		resolvedFlowDI[depName] = resolvedDep
	}

	return src.Node{
			Directives: node.Directives,
			EntityRef:  node.EntityRef,
			TypeArgs:   resolvedNodeArgs,
			Deps:       resolvedFlowDI,
			Meta:       node.Meta,
			ErrGuard:   node.ErrGuard,
		}, foundInterface{
			iface:    nodeIface,
			location: location,
		}, nil
}

// also does validation
func (a Analyzer) getNodeInterface(
	entity src.Entity,
	usesBindDirective bool,
	location src.Location,
	node src.Node,
	scope src.Scope,
	resolvedNodeArgs []typesystem.Expr,
) (src.Interface, *compiler.Error) {
	if entity.Kind == src.InterfaceEntity {
		if usesBindDirective {
			return src.Interface{}, &compiler.Error{
				Message:  "Interface node cannot use #bind directive",
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		if node.Deps != nil {
			return src.Interface{}, &compiler.Error{
				Message:  "Only component node can have dependency injection",
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		return entity.Interface, nil
	}

	externArgs, hasExternDirective := entity.Component.Directives[compiler.ExternDirective]

	if usesBindDirective && !hasExternDirective {
		return src.Interface{}, &compiler.Error{
			Message:  "Node can't use #bind if it isn't instantiated with the component that use #extern",
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if len(externArgs) > 1 && len(resolvedNodeArgs) != 1 {
		return src.Interface{}, &compiler.Error{
			Message:  "Component that use #extern directive with > 1 argument, must have exactly one type-argument for overloading",
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	iface := entity.Component.Interface

	_, hasAutoPortsDirective := entity.Component.Directives[compiler.AutoportsDirective]
	if !hasAutoPortsDirective {
		return iface, nil
	}

	// if we here then we have #autoports (only for structs RN)

	if len(iface.IO.In) != 0 {
		return src.Interface{}, &compiler.Error{
			Message:  "Component that uses struct inports directive must have no defined inports",
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if len(iface.TypeParams.Params) != 1 {
		return src.Interface{}, &compiler.Error{
			Message:  "Exactly one type parameter expected",
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	resolvedTypeParamConstr, err := a.resolver.ResolveExpr(iface.TypeParams.Params[0].Constr, scope)
	if err != nil {
		return src.Interface{}, &compiler.Error{
			Message:  err.Error(),
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if resolvedTypeParamConstr.Lit == nil || resolvedTypeParamConstr.Lit.Struct == nil {
		return src.Interface{}, &compiler.Error{
			Message:  "Struct type expected",
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if len(resolvedNodeArgs) != 1 {
		return src.Interface{}, &compiler.Error{
			Message:  "Exactly one type argument expected",
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	resolvedNodeArg, err := a.resolver.ResolveExpr(resolvedNodeArgs[0], scope)
	if err != nil {
		return src.Interface{}, &compiler.Error{
			Message:  err.Error(),
			Location: &location,
			Meta:     entity.Meta(),
		}
	}

	if resolvedNodeArg.Lit == nil || resolvedNodeArg.Lit.Struct == nil {
		return src.Interface{}, &compiler.Error{
			Message:  "Struct argument expected",
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

	// TODO refactor (maybe work for desugarer?)
	return src.Interface{
		TypeParams: iface.TypeParams,
		IO: src.IO{
			In: inports,
			Out: map[string]src.Port{
				"msg": {
					TypeExpr: resolvedNodeArg,
					IsArray:  false,
					Meta:     iface.IO.Out["msg"].Meta,
				},
			},
		},
		Meta: iface.Meta,
	}, nil
}

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
	location core.Location
}

func (a Analyzer) analyzeNodes(component src.Component, scope src.Scope) (
	map[string]src.Node, // resolved nodes
	map[string]foundInterface, // resolved nodes interfaces with locations
	bool, // one of the nodes has error guard
	*compiler.Error, // err
) {
	analyzedNodes := make(map[string]src.Node, len(component.Nodes))
	nodesInterfaces := make(map[string]foundInterface, len(component.Nodes))
	hasErrGuard := false

	for nodeName, node := range component.Nodes {
		if node.ErrGuard {
			hasErrGuard = true
		}

		analyzedNode, nodeInterface, err := a.analyzeNode(
			nodeName,
			node,
			scope,
			component,
		)
		if err != nil {
			return nil, nil, false, compiler.Error{
				Meta: &node.Meta,
			}.Wrap(err)
		}

		nodesInterfaces[nodeName] = nodeInterface
		analyzedNodes[nodeName] = analyzedNode
	}

	return analyzedNodes, nodesInterfaces, hasErrGuard, nil
}

func (a Analyzer) analyzeNode(
	name string,
	node src.Node,
	scope src.Scope,
	parent src.Component,
) (src.Node, foundInterface, *compiler.Error) {
	parentTypeParams := parent.Interface.TypeParams

	nodeEntity, location, err := scope.Entity(node.EntityRef)
	if err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &node.Meta,
		}
	}

	if nodeEntity.Kind != src.ComponentEntity &&
		nodeEntity.Kind != src.InterfaceEntity {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: fmt.Sprintf("Node can only refer to flows or interfaces: %v", nodeEntity.Kind),
			Meta:    nodeEntity.Meta(),
		}
	}

	bindArg, hasBind := node.Directives[compiler.BindDirective]
	if hasBind && len(bindArg) != 1 {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: "Node with #bind directive must provide exactly one argument",
			Meta:    nodeEntity.Meta(),
		}
	}

	if hasBind && nodeEntity.Kind == src.InterfaceEntity {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: "Interface node cannot use #bind directive",
			Meta:    nodeEntity.Meta(),
		}
	}

	if nodeEntity.Kind == src.InterfaceEntity && node.DIArgs != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: "Only component node can have dependency injection",
			Meta:    nodeEntity.Meta(),
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
			Message: err.Error(),
			Meta:    &node.Meta,
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
			Message: err.Error(),
			Meta:    &node.Meta,
		}
	}

	var (
		nodeIface     src.Interface
		overloadIndex *int // only for component nodes
	)
	if nodeEntity.Kind == src.InterfaceEntity {
		nodeIface = nodeEntity.Interface
	} else {
		var err *compiler.Error
		nodeIface, overloadIndex, err = a.getComponentNodeInterface(
			name,
			nodeEntity,
			hasBind,
			node,
			scope,
			resolvedNodeArgs,
			parent,
		)
		if err != nil {
			return src.Node{}, foundInterface{}, err
		}
	}

	if node.ErrGuard {
		if _, ok := parent.Interface.IO.Out["err"]; !ok {
			return src.Node{}, foundInterface{}, &compiler.Error{
				Message: "Error-guard operator '?' can only be used in components with ':err' outport to propagate errors",
				Meta:    &node.Meta,
			}
		}
		if _, ok := nodeIface.IO.Out["err"]; !ok {
			return src.Node{}, foundInterface{}, &compiler.Error{
				Message: "Error-guard operator '?' requires node to have ':err' outport to propagate errors",
				Meta:    &node.Meta,
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
			Message: err.Error(),
			Meta:    &node.Meta,
		}
	}

	if node.DIArgs == nil {
		return src.Node{
				Directives:    node.Directives,
				EntityRef:     node.EntityRef,
				TypeArgs:      resolvedNodeArgs,
				Meta:          node.Meta,
				OverloadIndex: overloadIndex,
				ErrGuard:      node.ErrGuard,
			}, foundInterface{
				iface:    nodeIface,
				location: location,
			}, nil
	}

	resolvedFlowDI := make(map[string]src.Node, len(node.DIArgs))
	for depName, depNode := range node.DIArgs {
		resolvedDep, _, err := a.analyzeNode(
			name, // TODO make sure DI works with overloading (example: Reduce{Add})
			depNode,
			scope,
			parent, // TODO make sure DI works with overloading (example: Reduce{Add})
		)
		if err != nil {
			return src.Node{}, foundInterface{}, compiler.Error{
				Meta: &depNode.Meta,
			}.Wrap(err)
		}
		resolvedFlowDI[depName] = resolvedDep
	}

	return src.Node{
			Directives:    node.Directives,
			EntityRef:     node.EntityRef,
			TypeArgs:      resolvedNodeArgs,
			DIArgs:        resolvedFlowDI,
			ErrGuard:      node.ErrGuard,
			OverloadIndex: overloadIndex,
			Meta:          node.Meta,
		}, foundInterface{
			iface:    nodeIface,
			location: location,
		}, nil
}

// getComponentNodeInterface returns interface of the component node.
// It also performs some validation.
// Overloading at the level of sourcecode is implemented here.
func (a Analyzer) getComponentNodeInterface(
	name string,
	entity src.Entity,
	hasBind bool,
	node src.Node,
	scope src.Scope,
	resolvedNodeArgs []typesystem.Expr,
	parent src.Component,
) (src.Interface, *int, *compiler.Error) {
	var (
		overloadIndex *int
		version       src.Component
	)
	if len(entity.Component) == 1 {
		version = entity.Component[0]
	} else {
		v, err := a.getNodeOverloadIndex(name, node, parent, entity.Component, scope)
		if err != nil {
			return src.Interface{}, nil, &compiler.Error{
				Message: "Node can't use #bind if it isn't instantiated with the component that use #extern",
				Meta:    &node.Meta,
			}
		}
		version = entity.Component[v]
		overloadIndex = &v
	}

	_, hasExtern := version.Directives[compiler.ExternDirective]
	if hasBind && !hasExtern {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Node can't use #bind if it isn't instantiated with the component that use #extern",
			Meta:    entity.Meta(),
		}
	}

	iface := version.Interface

	_, hasAutoPortsDirective := version.Directives[compiler.AutoportsDirective]
	if !hasAutoPortsDirective {
		return iface, overloadIndex, nil
	}

	// if we here then we have #autoports (only for structs)

	if len(iface.IO.In) != 0 {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Component that uses struct inports directive must have no defined inports",
			Meta:    entity.Meta(),
		}
	}

	if len(iface.TypeParams.Params) != 1 {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Exactly one type parameter expected",
			Meta:    entity.Meta(),
		}
	}

	resolvedTypeParamConstr, err := a.resolver.ResolveExpr(iface.TypeParams.Params[0].Constr, scope)
	if err != nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: err.Error(),
			Meta:    entity.Meta(),
		}
	}

	if resolvedTypeParamConstr.Lit == nil || resolvedTypeParamConstr.Lit.Struct == nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Struct type expected",
			Meta:    entity.Meta(),
		}
	}

	if len(resolvedNodeArgs) != 1 {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Exactly one type argument expected",
			Meta:    entity.Meta(),
		}
	}

	resolvedNodeArg, err := a.resolver.ResolveExpr(resolvedNodeArgs[0], scope)
	if err != nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: err.Error(),
			Meta:    entity.Meta(),
		}
	}

	if resolvedNodeArg.Lit == nil || resolvedNodeArg.Lit.Struct == nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Struct argument expected",
			Meta:    entity.Meta(),
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
				// struct builder has exactly one outport - created structure
				"res": {
					TypeExpr: resolvedNodeArg,
					IsArray:  false,
					Meta:     iface.IO.Out["res"].Meta,
				},
			},
		},
		Meta: iface.Meta,
	}, overloadIndex, nil
}

// getNodeOverloadIndex returns index of the overload of the node.
// It must only be called for components that are overloaded.
// It determines which overload to use based on node's usage in parent component.
func (a Analyzer) getNodeOverloadIndex(
	name string,
	node src.Node,
	parent src.Component,
	versions []src.Component,
	scope src.Scope,
) (int, error) {
	panic("not implemented")
}

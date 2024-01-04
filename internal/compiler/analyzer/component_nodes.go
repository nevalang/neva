package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
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
				Err:      fmt.Errorf("Invalid node '%v'", nodeName),
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

	var iface src.Interface
	if entity.Kind == src.ComponentEntity { //nolint:nestif
		runtimeFuncArgs, isRuntimeFunc := entity.Component.Directives[compiler.RuntimeFuncDirective]

		if hasRuntimeMsg && !isRuntimeFunc {
			return src.Node{}, src.Interface{}, &compiler.Error{
				Err:      ErrNormNodeRuntimeMsg,
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		if len(runtimeFuncArgs) > 1 && len(node.TypeArgs) != 1 {
			return src.Node{}, src.Interface{}, &compiler.Error{
				Err:      ErrRuntimeFuncOverloadingNodeArgs,
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		iface = entity.Component.Interface
	} else {
		if hasRuntimeMsg {
			return src.Node{}, src.Interface{}, &compiler.Error{
				Err:      ErrInterfaceNodeWithRuntimeMsg,
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		if node.Deps != nil {
			return src.Node{}, src.Interface{}, &compiler.Error{
				Err:      ErrNonComponentNodeWithDI,
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		iface = entity.Interface
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

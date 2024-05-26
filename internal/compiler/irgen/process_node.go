package irgen

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/runtime/ir"
)

func (g Generator) processComponentNode(
	nodeCtx nodeContext,
	scope src.Scope,
	result *ir.Program,
) *compiler.Error {
	ref := nodeCtx.node.EntityRef

	entity, found, err := scope.Entity(ref)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &scope.Location,
		}
	}

	comp := entity.Component
	in := g.insertAndReturnInports(nodeCtx, result)   // for inports we only use parent context because all inports are used
	out := g.insertAndReturnOutports(nodeCtx, result) //  for outports we use both parent context and component's interface

	funcRef, err := getFuncRef(
		comp,
		nodeCtx.node.TypeArgs,
	)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &found,
			Meta:     &comp.Meta,
		}
	}

	if funcRef != "" {
		funcCall, err := getFuncCall(nodeCtx, scope, funcRef, in, out)
		if err != nil {
			return err
		}
		result.Funcs = append(result.Funcs, funcCall)
		return nil
	}

	newScope := scope.WithLocation(found)

	subnodesPortsUsage, err := g.processNetwork(comp.Net, nodeCtx, result)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &newScope.Location,
		}
	}

	for nodeName, node := range comp.Nodes {
		nodePortsUsage := subnodesPortsUsage[nodeName]

		subNodeCtx := nodeContext{
			path:       append(nodeCtx.path, nodeName),
			portsUsage: nodePortsUsage,
			node:       node,
		}

		scopeToUse := getSubnodeScope(
			nodeCtx.node.Deps,
			nodeName,
			subNodeCtx,
			scope,
			newScope,
		)

		if err := g.processComponentNode(
			subNodeCtx,
			scopeToUse,
			result,
		); err != nil {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: node '%v'", err, nodeName),
				Location: &found,
				Meta:     &comp.Meta,
			}
		}
	}

	return nil
}

func getSubnodeScope(
	deps map[string]src.Node,
	nodeName string,
	subNodeCtx nodeContext,
	scope src.Scope,
	newScope src.Scope,
) src.Scope {
	var scopeToUseThisTime src.Scope
	if dep, ok := deps[nodeName]; ok { // is interface node
		subNodeCtx.node = dep
		scopeToUseThisTime = scope
	} else {
		scopeToUseThisTime = newScope
	}
	return scopeToUseThisTime
}

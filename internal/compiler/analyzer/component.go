package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
)

func (a Analyzer) analyzeComponent(comp src.Component, prog src.Program) (src.Component, error) {
	resolvedInterface, err := a.analyzeInterface(comp.Interface)
	if err != nil {
		return src.Component{}, fmt.Errorf("analyze interface: %w", err)
	}

	if err := a.analyzeComponentNodes(comp.Nodes, prog); err != nil {
		return src.Component{}, fmt.Errorf("analyze component nodes: %w", err)
	}

	normalizedNetwork, err := a.analyzeComponentNet(comp.Net, resolvedInterface, comp.Nodes)
	if err != nil {
		return src.Component{}, fmt.Errorf("analyze component network: %w", err)
	}

	return src.Component{
		Interface: resolvedInterface,
		Nodes:     comp.Nodes,
		Net:       normalizedNetwork,
	}, nil
}

func (a Analyzer) analyzeComponentNodes(nodes map[string]src.Node, prog src.Program) (map[string]src.Node, error) {
	resolvedNodes := make(map[string]src.Node, len(nodes))
	for name, node := range nodes {
		resolvedNode, err := a.analyzeComponentNode(node, prog)
		if err != nil {
			return nil, fmt.Errorf("analyze node: %w", err)
		}
		resolvedNodes[name] = resolvedNode
	}
	return resolvedNodes, nil
}

var (
	ErrNodeWrongEntity           = fmt.Errorf("node entity is not a component or interface")
	ErrNodeTypeArgsCountMismatch = errors.New("node type args count mismatch")
	ErrNodeInterfaceDI           = errors.New("interface node cannot have dependency injection")
)

func (a Analyzer) analyzeComponentNode(node src.Node, prog src.Program) (src.Node, error) {
	entity, err := prog.Entity(node.EntityRef)
	if err != nil {
		return src.Node{}, fmt.Errorf("entity: %w", err)
	}

	if entity.Kind != src.ComponentEntity && entity.Kind != src.InterfaceEntity {
		return src.Node{}, fmt.Errorf("%w: %v", ErrNodeWrongEntity, entity.Kind)
	}

	var compInterface src.Interface
	if entity.Kind == src.ComponentEntity {
		compInterface = entity.Component.Interface
	} else {
		if node.ComponentDI != nil {
			return src.Node{}, ErrNodeInterfaceDI
		}
		compInterface = entity.Interface
	}

	if len(node.TypeArgs) != len(compInterface.TypeParams) {
		return src.Node{}, fmt.Errorf(
			"%w: want %v, got %v",
			ErrNodeTypeArgsCountMismatch, compInterface.TypeParams, node.TypeArgs,
		)
	}

	resolvedArgs, _, err := a.resolver.ResolveArgs(node.TypeArgs, compInterface.TypeParams, nil)
	if err != nil {
		return src.Node{}, fmt.Errorf("resolve args: %w", err)
	}

	if node.ComponentDI == nil {
		return src.Node{
			EntityRef: node.EntityRef,
			TypeArgs:  resolvedArgs,
		}, nil
	}

	resolvedComponentDI := make(map[string]src.Node, len(node.ComponentDI)) // TODO track recursion
	for depName, depNode := range node.ComponentDI {
		resolvedDep, err := a.analyzeComponentNode(depNode, prog)
		if err != nil {
			return src.Node{}, fmt.Errorf("analyze dependency node: %w", err)
		}
		resolvedComponentDI[depName] = resolvedDep
	}

	return src.Node{
		EntityRef:   node.EntityRef,
		TypeArgs:    resolvedArgs,
		ComponentDI: resolvedComponentDI,
	}, nil
}

func (a Analyzer) analyzeComponentNet(
	net []src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
) ([]src.Connection, error) {
	return net, nil // TODO
}

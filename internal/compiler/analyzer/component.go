package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
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

func (a Analyzer) analyzeComponentNodes(nodes map[string]src.Node, prog src.Program) error {
	for _, node := range nodes {
		if err := a.analyzeComponentNode(node, prog); err != nil {
			return fmt.Errorf("analyze node: %w", err)
		}
	}
	return nil
}

var (
	ErrNodeEntity                = fmt.Errorf("node entity is not a component or interface")
	ErrNodeTypeArgsCountMismatch = errors.New("node type args count mismatch")
)

func (a Analyzer) analyzeComponentNode(node src.Node, prog src.Program) error {
	entity, err := prog.Entity(node.EntityRef)
	if err != nil {
		return fmt.Errorf("entity: %w", err)
	}

	if entity.Kind != src.ComponentEntity && entity.Kind != src.InterfaceEntity {
		return fmt.Errorf("%w: %v", ErrNodeEntity, entity.Kind)
	}

	var compInterface src.Interface
	if entity.Kind == src.ComponentEntity {
		compInterface = entity.Component.Interface
	} else {
		compInterface = entity.Interface
	}

	if len(node.TypeArgs) != len(compInterface.TypeParams) {
		return fmt.Errorf(
			"%w: want %v, got %v",
			ErrNodeTypeArgsCountMismatch, compInterface.TypeParams, node.TypeArgs,
		)
	}

	resolvedArgs := make([]ts.Expr, 0, len(node.TypeArgs))
	for _, arg := range node.TypeArgs {
		resolvedArg, err := a.analyzeTypeExpr(arg)
		if err != nil {
			return fmt.Errorf("analyze type expr: %w", err)
		}
		resolvedArgs = append(resolvedArgs, resolvedArg)
	}

	// check that args are compatible with type params
	// this can be done by creating

	return nil
}

func (a Analyzer) analyzeComponentNet(
	net []src.Connection,
	compInterface src.Interface,
	nodes map[string]src.Node,
) ([]src.Connection, error) {
	return net, nil // TODO
}

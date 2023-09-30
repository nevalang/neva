package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
)

var (
	ErrMainComponentWithTypeParams     = errors.New("main component can't have type parameters")
	ErrMainComponentNodes              = errors.New("something wrong with main component's nodes")
	ErrEntityNotFoundByNodeRef         = errors.New("entity not found by node ref")
	ErrMainComponentInportsCount       = errors.New("main component must have one inport")
	ErrMainComponentOutportsCount      = errors.New("main component must have exactly one outport")
	ErrMainComponentWithoutEnterInport = errors.New("main component must have 'enter' inport")
	ErrMainComponentWithoutExitOutport = errors.New("main component must have 'exit' outport")
	ErrMainPortIsArray                 = errors.New("main component's ports cannot not be arrays")
	ErrMainComponentPortTypeNotAny     = errors.New("main component's ports must be of type any")
	ErrMainComponentNodeNotComponent   = errors.New("main component's nodes must be components only")
)

func (a Analyzer) analyzeMainComponent(cmp src.Component, pkg src.Package, pkgs map[string]src.Package) error { //nolint:unparam,lll
	if len(cmp.Interface.TypeParams) != 0 {
		return fmt.Errorf("%w: %v", ErrMainComponentWithTypeParams, cmp.Interface.TypeParams)
	}

	if err := a.analyzeMainComponentIO(cmp.Interface.IO); err != nil {
		return fmt.Errorf("main component io: %w", err)
	}

	if err := a.analyzeMainComponentNodes(cmp.Nodes, pkg, pkgs); err != nil {
		return fmt.Errorf("%w: %v", ErrMainComponentNodes, err)
	}

	return nil
}

func (a Analyzer) analyzeMainComponentIO(io src.IO) error {
	if len(io.Out) != 1 {
		return fmt.Errorf("%w: %v", ErrMainComponentOutportsCount, io.Out)
	}

	if len(io.In) != 1 {
		return fmt.Errorf("%w: %v", ErrMainComponentInportsCount, io.In)
	}

	enterInport, ok := io.In["enter"]
	if !ok {
		return ErrMainComponentWithoutEnterInport
	}

	if enterInport.IsArray {
		return ErrMainPortIsArray
	}

	if enterInport.TypeExpr != nil {
		return ErrMainComponentPortTypeNotAny
	}

	exitInport, ok := io.In["exit"]
	if !ok {
		return ErrMainComponentWithoutExitOutport
	}

	if exitInport.IsArray {
		return ErrMainPortIsArray
	}

	if exitInport.TypeExpr != nil {
		return ErrMainComponentPortTypeNotAny
	}

	return nil
}

func (Analyzer) analyzeMainComponentNodes(nodes map[string]src.Node, pkg src.Package, prog src.Program) error {
	for name, node := range nodes {
		nodeEntity, err := prog.Entity(node.EntityRef)
		if err != nil {
			return fmt.Errorf("%w: %v: %v", ErrEntityNotFoundByNodeRef, name, node.EntityRef)
		}
		if nodeEntity.Kind != src.ComponentEntity {
			return fmt.Errorf("%w: %v: %v", ErrMainComponentNodeNotComponent, name, node.EntityRef)
		}
	}
	return nil
}

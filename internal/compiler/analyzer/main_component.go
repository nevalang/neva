package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

func (a Analyzer) analyzeMainComponent(cmp src.Component, scope src.Scope) *compiler.Error {
	if len(cmp.Interface.TypeParams.Params) != 0 {
		return &compiler.Error{
			Message: "Main component cannot have type parameters",
			Meta:    &cmp.Interface.Meta,
		}
	}

	if err := a.analyzeMainFlowIO(cmp.Interface.IO); err != nil {
		return compiler.Error{Meta: &cmp.Interface.Meta}.Wrap(err)
	}

	if err := a.analyzeMainFlowNodes(cmp.Nodes, scope); err != nil {
		return compiler.Error{Meta: &cmp.Meta}.Wrap(err)
	}

	return nil
}

func (a Analyzer) analyzeMainFlowIO(io src.IO) *compiler.Error {
	if len(io.In) != 1 {
		return &compiler.Error{
			Message: fmt.Sprintf("Main component must have exactly 1 inport: got %v", len(io.In)),
		}
	}
	if len(io.Out) != 1 {
		return &compiler.Error{
			Message: fmt.Sprintf("Main component must have exactly 1 outport: got %v", len(io.Out)),
		}
	}

	enterInport, ok := io.In["start"]
	if !ok {
		return &compiler.Error{Message: "Main component must have 'start' inport"}
	}

	if err := a.analyzeMainComponentPort(enterInport); err != nil {
		return err
	}

	exitOutport, ok := io.Out["stop"]
	if !ok {
		return &compiler.Error{Message: "Main component must have 'stop' outport"}
	}

	if err := a.analyzeMainComponentPort(exitOutport); err != nil {
		return err
	}

	return nil
}

func (a Analyzer) analyzeMainComponentPort(port src.Port) *compiler.Error {
	if port.IsArray {
		return &compiler.Error{
			Message: "Main component's ports cannot be arrays",
			Meta:    &port.Meta,
		}
	}
	if !(src.Scope{}).IsTopType(port.TypeExpr) {
		return &compiler.Error{
			Message: "Main component's ports must be of type any",
			Meta:    &port.Meta,
		}
	}
	return nil
}

func (Analyzer) analyzeMainFlowNodes(
	nodes map[string]src.Node,
	scope src.Scope,
) *compiler.Error {
	for nodeName, node := range nodes {
		nodeEntity, loc, err := scope.Entity(node.EntityRef)
		if err != nil {
			return &compiler.Error{
				Message: fmt.Sprintf(
					"Referenced entity not found: node '%v', ref '%v', details '%v'",
					nodeName,
					node.EntityRef,
					err,
				),
				Location: &loc,
				Meta:     &node.EntityRef.Meta,
			}
		}

		if nodeEntity.Kind != src.ComponentEntity {
			return &compiler.Error{
				Message:  fmt.Sprintf("Main component's nodes must only refer to component entities: %v: %v", nodeName, node.EntityRef),
				Location: &loc,
				Meta:     nodeEntity.Meta(),
			}
		}
	}

	return nil
}

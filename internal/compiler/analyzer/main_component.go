package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

func (a Analyzer) analyzeMainComponent(cmp src.Component, scope src.Scope) *compiler.Error {
	if len(cmp.Interface.TypeParams.Params) != 0 {
		return &compiler.Error{
			Err:  errors.New("Main flow cannot have type parameters"),
			Meta: &cmp.Interface.Meta,
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
			Err: fmt.Errorf("Main flow must have exactly 1 inport: got %v", len(io.In)),
		}
	}
	if len(io.Out) != 1 {
		return &compiler.Error{
			Err: fmt.Errorf("Main flow must have exactly 1 outport: got %v", len(io.Out)),
		}
	}

	enterInport, ok := io.In["start"]
	if !ok {
		return &compiler.Error{Err: errors.New("Main flow must have 'start' inport")}
	}
	if err := a.analyzeMainFlowPort(enterInport); err != nil {
		return &compiler.Error{
			Err:  err,
			Meta: &enterInport.Meta,
		}
	}

	exitOutport, ok := io.Out["stop"]
	if !ok {
		return &compiler.Error{Err: errors.New("Main flow must have 'stop' outport")}
	}
	if err := a.analyzeMainFlowPort(exitOutport); err != nil {
		return &compiler.Error{
			Err:  err,
			Meta: &exitOutport.Meta,
		}
	}

	return nil
}

func (a Analyzer) analyzeMainFlowPort(port src.Port) error {
	if port.IsArray {
		return errors.New("Main flow's ports cannot be arrays")
	}
	if !(src.Scope{}).IsTopType(port.TypeExpr) {
		return errors.New("Main flow's ports must be of type any")
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
				Err: fmt.Errorf(
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
				Err:      fmt.Errorf("Main flow's nodes must only refer to flow entities: %v: %v", nodeName, node.EntityRef),
				Location: &loc,
				Meta:     nodeEntity.Meta(),
			}
		}
	}

	return nil
}

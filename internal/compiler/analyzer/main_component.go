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
			Meta:    &cmp.Interface.TypeParams.Meta,
		}
	}

	if err := a.analyzeMainFlowIO(cmp.Interface.IO); err != nil {
		return compiler.Error{Meta: &cmp.Interface.Meta}.Wrap(err)
	}

	if err := a.analyzeMainComponentNodes(cmp.Nodes, scope); err != nil {
		return compiler.Error{Meta: &cmp.Meta}.Wrap(err)
	}

	return nil
}

func (a Analyzer) analyzeMainFlowIO(io src.IO) *compiler.Error {
	if len(io.In) != 1 {
		return &compiler.Error{
			Message: fmt.Sprintf("Main component must have exactly 1 inport: got %v", len(io.In)),
			Meta:    &io.Meta,
		}
	}
	if len(io.Out) != 1 {
		return &compiler.Error{
			Message: fmt.Sprintf("Main component must have exactly 1 outport: got %v", len(io.Out)),
			Meta:    &io.Meta,
		}
	}

	enterInport, ok := io.In["start"]
	if !ok {
		return &compiler.Error{Message: "Main component must have 'start' inport", Meta: &io.Meta}
	}

	if err := a.analyzeMainComponentPort(enterInport); err != nil {
		return err
	}

	exitOutport, ok := io.Out["stop"]
	if !ok {
		return &compiler.Error{Message: "Main component must have 'stop' outport", Meta: &io.Meta}
	}

	if err := a.analyzeMainComponentPort(exitOutport); err != nil {
		return compiler.Error{Meta: &exitOutport.Meta}.Wrap(err)
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

func (Analyzer) analyzeMainComponentNodes(
	nodes map[string]src.Node,
	scope src.Scope,
) *compiler.Error {
	for _, node := range nodes {
		if _, err := scope.GetComponent(node.EntityRef); err != nil {
			return &compiler.Error{
				Message: err.Error(),
				Meta:    &node.EntityRef.Meta,
			}
		}
	}
	return nil
}

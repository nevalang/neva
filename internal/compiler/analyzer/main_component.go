package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

var (
	ErrMainFlowWithTypeParams     = errors.New("Main flow cannot have type parameters")
	ErrEntityNotFoundByNodeRef    = errors.New("Node references to entity that cannot be found")
	ErrMainFlowInportsCount       = errors.New("Main flow must have exactly 1 inport")
	ErrMainFlowOutportsCount      = errors.New("Main flow must have exactly 1 outport")
	ErrMainFlowWithoutEnterInport = errors.New("Main flow must have 'enter' inport")
	ErrMainFlowWithoutExitOutport = errors.New("Main flow must have 'exit' outport")
	ErrMainPortIsArray            = errors.New("Main flow cannot have array ports")
	ErrMainFlowPortTypeNotAny     = errors.New("Main flow's ports must be of type any")
	ErrMainNodeEntityNotFlow      = errors.New("Main flow's nodes must only refer to flow entities")
)

func (a Analyzer) analyzeMainComponent(cmp src.Component, scope src.Scope) *compiler.Error {
	if len(cmp.Interface.TypeParams.Params) != 0 {
		return &compiler.Error{
			Err:   ErrMainFlowWithTypeParams,
			Range: &cmp.Interface.Meta,
		}
	}

	if err := a.analyzeMainFlowIO(cmp.Interface.IO); err != nil {
		return compiler.Error{Range: &cmp.Interface.Meta}.Wrap(err)
	}

	if err := a.analyzeMainFlowNodes(cmp.Nodes, scope); err != nil {
		return compiler.Error{Range: &cmp.Meta}.Wrap(err)
	}

	return nil
}

func (a Analyzer) analyzeMainFlowIO(io src.IO) *compiler.Error {
	if len(io.In) != 1 {
		return &compiler.Error{
			Err: fmt.Errorf("%w: got %v", ErrMainFlowInportsCount, len(io.In)),
		}
	}
	if len(io.Out) != 1 {
		return &compiler.Error{
			Err: fmt.Errorf("%w: got %v", ErrMainFlowOutportsCount, len(io.Out)),
		}
	}

	enterInport, ok := io.In["start"]
	if !ok {
		return &compiler.Error{Err: ErrMainFlowWithoutEnterInport}
	}
	if err := a.analyzeMainFlowPort(enterInport); err != nil {
		return &compiler.Error{
			Err:   err,
			Range: &enterInport.Meta,
		}
	}

	exitOutport, ok := io.Out["stop"]
	if !ok {
		return &compiler.Error{Err: ErrMainFlowWithoutExitOutport}
	}
	if err := a.analyzeMainFlowPort(exitOutport); err != nil {
		return &compiler.Error{
			Err:   err,
			Range: &exitOutport.Meta,
		}
	}

	return nil
}

func (a Analyzer) analyzeMainFlowPort(port src.Port) error {
	if port.IsArray {
		return ErrMainPortIsArray
	}
	if !(src.Scope{}).IsTopType(port.TypeExpr) {
		return ErrMainFlowPortTypeNotAny
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
					"%w: node '%v', ref '%v', details '%v'",
					ErrEntityNotFoundByNodeRef,
					nodeName,
					node.EntityRef,
					err,
				),
				Location: &loc,
				Range:    &node.EntityRef.Meta,
			}
		}

		if nodeEntity.Kind != src.ComponentEntity {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: %v: %v", ErrMainNodeEntityNotFlow, nodeName, node.EntityRef),
				Location: &loc,
				Range:    nodeEntity.Meta(),
			}
		}
	}

	return nil
}

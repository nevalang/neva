package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

var (
	ErrMainComponentWithTypeParams     = errors.New("Main component cannot have type parameters")
	ErrEntityNotFoundByNodeRef         = errors.New("Node references to entity that cannot be found")
	ErrMainComponentInportsCount       = errors.New("Main component must have exactly 1 inport")
	ErrMainComponentOutportsCount      = errors.New("Main component must have exactly 1 outport")
	ErrMainComponentWithoutEnterInport = errors.New("Main component must have 'enter' inport")
	ErrMainComponentWithoutExitOutport = errors.New("Main component must have 'exit' outport")
	ErrMainPortIsArray                 = errors.New("Main component cannot have array ports")
	ErrMainComponentPortTypeNotAny     = errors.New("Main component's ports must be of type any")
	ErrMainComponentNodeNotComponent   = errors.New("Main component's nodes must only refer to components")
)

func (a Analyzer) analyzeMainComponent(cmp src.Component, scope src.Scope) *compiler.Error {
	if len(cmp.Interface.TypeParams.Params) != 0 {
		return &compiler.Error{
			Err:  ErrMainComponentWithTypeParams,
			Meta: &cmp.Interface.Meta,
		}
	}

	if err := a.analyzeMainComponentIO(cmp.Interface.IO); err != nil {
		return compiler.Error{Meta: &cmp.Interface.Meta}.Wrap(err)
	}

	if err := a.analyzeMainComponentNodes(cmp.Nodes, scope); err != nil {
		return compiler.Error{Meta: &cmp.Meta}.Wrap(err)
	}

	return nil
}

func (a Analyzer) analyzeMainComponentIO(io src.IO) *compiler.Error {
	if len(io.In) != 1 {
		return &compiler.Error{
			Err: fmt.Errorf("%w: got %v", ErrMainComponentInportsCount, len(io.In)),
		}
	}
	if len(io.Out) != 1 {
		return &compiler.Error{
			Err: fmt.Errorf("%w: got %v", ErrMainComponentOutportsCount, len(io.Out)),
		}
	}

	enterInport, ok := io.In["start"]
	if !ok {
		return &compiler.Error{Err: ErrMainComponentWithoutEnterInport}
	}
	if err := a.analyzeMainComponentPort(enterInport); err != nil {
		return &compiler.Error{
			Err:  err,
			Meta: &enterInport.Meta,
		}
	}

	exitOutport, ok := io.Out["stop"]
	if !ok {
		return &compiler.Error{Err: ErrMainComponentWithoutExitOutport}
	}
	if err := a.analyzeMainComponentPort(exitOutport); err != nil {
		return &compiler.Error{
			Err:  err,
			Meta: &exitOutport.Meta,
		}
	}

	return nil
}

func (a Analyzer) analyzeMainComponentPort(port src.Port) error {
	if port.IsArray {
		return ErrMainPortIsArray
	}
	if !(src.Scope{}).IsTopType(port.TypeExpr) {
		return ErrMainComponentPortTypeNotAny
	}
	return nil
}

func (Analyzer) analyzeMainComponentNodes(nodes map[string]src.Node, scope src.Scope) *compiler.Error {
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
				Meta:     &node.EntityRef.Meta,
			}
		}

		if nodeEntity.Kind != src.ComponentEntity {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: %v: %v", ErrMainComponentNodeNotComponent, nodeName, node.EntityRef),
				Location: &loc,
				Meta:     nodeEntity.Meta(),
			}
		}
	}

	return nil
}

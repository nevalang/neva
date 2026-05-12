package analyzer

import (
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/ast"
)

func validateEntityComments(entity *src.Entity) *compiler.Error {
	if entity.Comments == nil {
		return nil
	}

	hasInportsComments := len(entity.Comments.Inports) > 0
	hasOutportsComments := len(entity.Comments.Outports) > 0
	if hasInportsComments != hasOutportsComments {
		return &compiler.Error{
			Message: "comment tags must include both @inport and @outport sections",
			Meta:    &entity.Comments.Meta,
		}
	}

	if !hasInportsComments {
		return nil
	}

	entityIO, err := entityIOForComments(entity)
	if err != nil {
		return err
	}

	if err := validateCommentPortsExist(entity.Comments, entityIO); err != nil {
		return err
	}

	return validateCommentPortsAreComplete(entity.Comments, entityIO)
}

func entityIOForComments(entity *src.Entity) (*src.IO, *compiler.Error) {
	switch entity.Kind {
	case src.InterfaceEntity:
		return &entity.Interface.IO, nil
	case src.ComponentEntity:
		if len(entity.Component) == 0 {
			panic("component entity has no component versions")
		}
		return &entity.Component[0].IO, nil
	case src.ConstEntity, src.TypeEntity:
		return nil, &compiler.Error{
			Message: "comment tags @inport/@outport are allowed only on interface or component entities",
			Meta:    &entity.Comments.Meta,
		}
	default:
		panic("unknown entity kind")
	}
}

func validateCommentPortsExist(comments *src.Comments, entityIO *src.IO) *compiler.Error {
	for portName := range comments.Inports {
		if _, ok := entityIO.In[portName]; !ok {
			return &compiler.Error{
				Message: "comment references unknown inport: " + portName,
				Meta:    &comments.Meta,
			}
		}
	}

	for portName := range comments.Outports {
		if _, ok := entityIO.Out[portName]; !ok {
			return &compiler.Error{
				Message: "comment references unknown outport: " + portName,
				Meta:    &comments.Meta,
			}
		}
	}

	return nil
}

func validateCommentPortsAreComplete(comments *src.Comments, entityIO *src.IO) *compiler.Error {
	for portName := range entityIO.In {
		if _, ok := comments.Inports[portName]; ok {
			continue
		}

		return &compiler.Error{
			Message: "comment must document all inports; missing @inport " + portName,
			Meta:    &comments.Meta,
		}
	}

	for portName := range entityIO.Out {
		if _, ok := comments.Outports[portName]; ok {
			continue
		}

		return &compiler.Error{
			Message: "comment must document all outports; missing @outport " + portName,
			Meta:    &comments.Meta,
		}
	}

	return nil
}

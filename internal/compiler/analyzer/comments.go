package analyzer

import (
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/ast"
)

func validateEntityComments(entity *src.Entity) *compiler.Error {
	if entity.Comments == nil {
		return nil
	}

	hasPortTagComments := len(entity.Comments.Inports) > 0 || len(entity.Comments.Outports) > 0
	if !hasPortTagComments {
		return nil
	}

	if !supportsPortCommentTags(entity.Kind) {
		return &compiler.Error{
			Message: "comment tags @inport/@outport are allowed only on interface or component entities",
			Meta:    &entity.Comments.Meta,
		}
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

	entityIO := entityIOForComments(entity)

	if err := validateCommentPortsExist(entity.Comments, entityIO); err != nil {
		return err
	}

	return validateCommentPortsAreComplete(entity.Comments, entityIO)
}

func supportsPortCommentTags(kind src.EntityKind) bool {
	switch kind {
	case src.InterfaceEntity, src.ComponentEntity:
		return true
	case src.ConstEntity, src.TypeEntity:
		return false
	default:
		panic("unknown entity kind")
	}
}

func entityIOForComments(entity *src.Entity) *src.IO {
	switch entity.Kind {
	case src.InterfaceEntity:
		return &entity.Interface.IO
	case src.ComponentEntity:
		if len(entity.Component) == 0 {
			panic("component entity has no component versions")
		}
		return &entity.Component[0].IO
	case src.ConstEntity, src.TypeEntity:
		panic("unexpected entity kind for port comment tags")
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

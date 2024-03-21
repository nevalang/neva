package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

var (
	ErrEmptyConst         = errors.New("Constant must either have value or reference to another constant")
	ErrEntityNotConst     = errors.New("Constant refers to an entity that is not constant")
	ErrResolveConstType   = errors.New("Cannot resolve constant type")
	ErrUnionConst         = errors.New("Constant cannot have type union")
	ErrConstSeveralValues = errors.New("Constant cannot have several values at once")
)

//nolint:funlen
func (a Analyzer) analyzeConst(
	constant src.Const,
	scope src.Scope,
) (src.Const, *compiler.Error) {
	if constant.Message == nil && constant.Ref == nil {
		return src.Const{}, &compiler.Error{
			Err:      ErrEmptyConst,
			Location: &scope.Location,
			Meta:     &constant.Meta,
		}
	}

	if constant.Message == nil { // is ref
		entity, location, err := scope.Entity(*constant.Ref)
		if err != nil {
			return src.Const{}, &compiler.Error{
				Err:      err,
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		if entity.Kind != src.ConstEntity {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: entity kind %v", ErrEntityNotConst, entity.Kind),
				Location: &location,
				Meta:     entity.Meta(),
			}
		}

		return a.analyzeConst(entity.Const, scope)
	}

	resolvedType, err := a.analyzeTypeExpr(constant.Message.TypeExpr, scope)
	if err != nil {
		return src.Const{}, compiler.Error{
			Err:      ErrResolveConstType,
			Location: &scope.Location,
			Meta:     &constant.Meta,
		}.Wrap(err)
	}

	if resolvedType.Lit != nil && resolvedType.Lit.Union != nil {
		return src.Const{}, &compiler.Error{
			Err:      ErrUnionConst,
			Location: &scope.Location,
			Meta:     &constant.Meta,
		}
	}

	var typeExprStrRepr string
	if inst := resolvedType.Inst; inst != nil {
		typeExprStrRepr = inst.Ref.String()
	} else if lit := resolvedType.Lit; lit != nil {
		if lit.Enum != nil {
			typeExprStrRepr = "enum"
		} else if lit.Struct != nil {
			typeExprStrRepr = "struct"
		}
	}

	switch typeExprStrRepr {
	case "bool":
		if constant.Message.Bool == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("Boolean value is missing in boolean contant: %v", constant),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
		if constant.Message.Int != nil ||
			constant.Message.Float != nil ||
			constant.Message.Str != nil ||
			constant.Message.List != nil ||
			constant.Message.MapOrStruct != nil ||
			constant.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "int":
		if constant.Message.Int == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("Integer value is missing in integer contant: %v", constant),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
		if constant.Message.Bool != nil ||
			constant.Message.Float != nil ||
			constant.Message.Str != nil ||
			constant.Message.List != nil ||
			constant.Message.MapOrStruct != nil ||
			constant.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Message),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "float":
		// Float is special case. Constant can have float type expression but integer literal.
		// We must pass this case. Desugarer will turn integer literal into float.
		if constant.Message.Float == nil && constant.Message.Int == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("Float or integer value is missing in float contant: %v", constant),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
		if constant.Message.Float != nil && constant.Message.Int != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Message),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
		if constant.Message.Bool != nil ||
			constant.Message.Str != nil ||
			constant.Message.List != nil ||
			constant.Message.MapOrStruct != nil ||
			constant.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Message),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "string":
		if constant.Message.Str == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("String value is missing in string contant: %v", constant),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
		if constant.Message.Bool != nil ||
			constant.Message.Int != nil ||
			constant.Message.Float != nil ||
			constant.Message.List != nil ||
			constant.Message.MapOrStruct != nil ||
			constant.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Message),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "list":
		if constant.Message.List == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("List value is missing in list contant: %v", constant),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
		if constant.Message.Bool != nil ||
			constant.Message.Int != nil ||
			constant.Message.Float != nil ||
			constant.Message.MapOrStruct != nil ||
			constant.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Message),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "map", "struct":
		if constant.Message.MapOrStruct == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("Map or struct value is missing in map or struct contant: %v", constant),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
		if constant.Message.Bool != nil ||
			constant.Message.Int != nil ||
			constant.Message.Float != nil ||
			constant.Message.List != nil ||
			constant.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Message),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "enum":
		if constant.Message.Enum == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("Enum value is missing in enum contant: %v", constant),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
		if constant.Message.Bool != nil ||
			constant.Message.Int != nil ||
			constant.Message.Float != nil ||
			constant.Message.List != nil ||
			constant.Message.MapOrStruct != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Message),
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	}

	valueCopy := *constant.Message
	valueCopy.TypeExpr = resolvedType

	return src.Const{
		Message: &valueCopy,
		Meta:    constant.Meta,
	}, nil
}

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

func (a Analyzer) analyzeConst(
	constant src.Const,
	scope src.Scope,
) (src.Const, *compiler.Error) {
	if constant.Value.Message == nil && constant.Value.Ref == nil {
		return src.Const{}, &compiler.Error{
			Err:      ErrEmptyConst,
			Location: &scope.Location,
			Range:    &constant.Meta,
		}
	}

	if constant.Value.Message == nil { // is ref
		entity, location, err := scope.Entity(*constant.Value.Ref)
		if err != nil {
			return src.Const{}, &compiler.Error{
				Err:      err,
				Location: &location,
				Range:    entity.Meta(),
			}
		}

		if entity.Kind != src.ConstEntity {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: entity kind %v", ErrEntityNotConst, entity.Kind),
				Location: &location,
				Range:    entity.Meta(),
			}
		}

		return a.analyzeConst(entity.Const, scope)
	}

	resolvedType, err := a.analyzeTypeExpr(constant.TypeExpr, scope)
	if err != nil {
		return src.Const{}, compiler.Error{
			Err:      ErrResolveConstType,
			Location: &scope.Location,
			Range:    &constant.Meta,
		}.Wrap(err)
	}

	if resolvedType.Lit != nil && resolvedType.Lit.Union != nil {
		return src.Const{}, &compiler.Error{
			Err:      ErrUnionConst,
			Location: &scope.Location,
			Range:    &constant.Meta,
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
		if constant.Value.Message.Bool == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("Boolean value is missing in boolean contant: %v", constant),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
		if constant.Value.Message.Int != nil ||
			constant.Value.Message.Float != nil ||
			constant.Value.Message.Str != nil ||
			constant.Value.Message.List != nil ||
			constant.Value.Message.DictOrStruct != nil ||
			constant.Value.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
	case "int":
		if constant.Value.Message.Int == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("Integer value is missing in integer contant: %v", constant),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
		if constant.Value.Message.Bool != nil ||
			constant.Value.Message.Float != nil ||
			constant.Value.Message.Str != nil ||
			constant.Value.Message.List != nil ||
			constant.Value.Message.DictOrStruct != nil ||
			constant.Value.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Value.Message),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
	case "float":
		// Float is special case. Constant can have float type expression but integer literal.
		// We must pass this case. Desugarer will turn integer literal into float.
		if constant.Value.Message.Float == nil && constant.Value.Message.Int == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("Float or integer value is missing in float contant: %v", constant),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
		if constant.Value.Message.Float != nil && constant.Value.Message.Int != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Value.Message),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
		if constant.Value.Message.Bool != nil ||
			constant.Value.Message.Str != nil ||
			constant.Value.Message.List != nil ||
			constant.Value.Message.DictOrStruct != nil ||
			constant.Value.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Value.Message),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
	case "string":
		if constant.Value.Message.Str == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("String value is missing in string contant: %v", constant),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
		if constant.Value.Message.Bool != nil ||
			constant.Value.Message.Int != nil ||
			constant.Value.Message.Float != nil ||
			constant.Value.Message.List != nil ||
			constant.Value.Message.DictOrStruct != nil ||
			constant.Value.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Value.Message),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
	case "list":
		if constant.Value.Message.List == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("List value is missing in list contant: %v", constant),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
		if constant.Value.Message.Bool != nil ||
			constant.Value.Message.Int != nil ||
			constant.Value.Message.Float != nil ||
			constant.Value.Message.DictOrStruct != nil ||
			constant.Value.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Value.Message),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
	case "map", "struct":
		if constant.Value.Message.DictOrStruct == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("Map or struct value is missing in map or struct contant: %v", constant),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
		if constant.Value.Message.Bool != nil ||
			constant.Value.Message.Int != nil ||
			constant.Value.Message.Float != nil ||
			constant.Value.Message.List != nil ||
			constant.Value.Message.Enum != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Value.Message),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
	case "enum":
		if constant.Value.Message.Enum == nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("Enum value is missing in enum contant: %v", constant),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
		if constant.Value.Message.Bool != nil ||
			constant.Value.Message.Int != nil ||
			constant.Value.Message.Float != nil ||
			constant.Value.Message.List != nil ||
			constant.Value.Message.DictOrStruct != nil {
			return src.Const{}, &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrConstSeveralValues, constant.Value.Message),
				Location: &scope.Location,
				Range:    &constant.Meta,
			}
		}
	}

	return src.Const{
		TypeExpr: resolvedType,
		Value:    constant.Value,
		Meta:     constant.Meta,
	}, nil
}

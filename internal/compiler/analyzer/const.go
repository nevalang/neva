package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

var (
	ErrEmptyConst         = errors.New("Constant must either have value or reference to another constant")
	ErrEntityNotConst     = errors.New("Constant refers to an entity that is not constant")
	ErrResolveConstType   = errors.New("Cannot resolve constant type")
	ErrUnionConst         = errors.New("Constant cannot have type union")
	ErrConstSeveralValues = errors.New("Constant cannot have several values at once")
)

//nolint:funlen
func (a Analyzer) analyzeConst(constant src.Const, scope src.Scope) (src.Const, *compiler.Error) { //nolint:gocyclo,gocognit,lll
	if constant.Value == nil && constant.Ref == nil {
		return src.Const{}, &compiler.Error{
			Err:      ErrEmptyConst,
			Location: &scope.Location,
			Meta:     &constant.Meta,
		}
	}

	if constant.Value == nil { // is ref
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
	}

	resolvedType, err := a.analyzeTypeExpr(constant.Value.TypeExpr, scope)
	if err != nil {
		return src.Const{}, compiler.Error{
			Err:      ErrResolveConstType,
			Location: &scope.Location,
			Meta:     &constant.Meta,
		}.Merge(err)
	}

	if resolvedType.Lit != nil && resolvedType.Lit.Union != nil {
		return src.Const{}, &compiler.Error{
			Err:      ErrUnionConst,
			Location: &scope.Location,
			Meta:     &constant.Meta,
		}
	}

	var typ string
	if inst := resolvedType.Inst; inst != nil {
		typ = inst.Ref.String()
	} else if lit := resolvedType.Lit; lit != nil {
		if lit.Enum != nil {
			typ = "enum"
		} else if lit.Struct != nil {
			typ = "struct"
		}
	}

	switch typ {
	case "bool":
		if constant.Value.Int != nil ||
			constant.Value.Float != nil ||
			constant.Value.Str != nil ||
			constant.Value.List != nil ||
			constant.Value.Map != nil {
			return src.Const{}, &compiler.Error{
				Err:      ErrConstSeveralValues,
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "int", "enum":
		if constant.Value.Bool != nil ||
			constant.Value.Float != nil ||
			constant.Value.Str != nil ||
			constant.Value.List != nil ||
			constant.Value.Map != nil {
			return src.Const{}, &compiler.Error{
				Err:      ErrConstSeveralValues,
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "float":
		if constant.Value.Bool != nil ||
			constant.Value.Int != nil ||
			constant.Value.Str != nil ||
			constant.Value.List != nil ||
			constant.Value.Map != nil {
			return src.Const{}, &compiler.Error{
				Err:      ErrConstSeveralValues,
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "str":
		if constant.Value.Bool != nil ||
			constant.Value.Int != nil ||
			constant.Value.Float != nil ||
			constant.Value.List != nil ||
			constant.Value.Map != nil {
			return src.Const{}, &compiler.Error{
				Err:      ErrConstSeveralValues,
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "list":
		if constant.Value.Bool != nil ||
			constant.Value.Int != nil ||
			constant.Value.Float != nil ||
			constant.Value.Map != nil {
			return src.Const{}, &compiler.Error{
				Err:      ErrConstSeveralValues,
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	case "map", "struct":
		if constant.Value.Bool != nil ||
			constant.Value.Int != nil ||
			constant.Value.Float != nil ||
			constant.Value.List != nil {
			return src.Const{}, &compiler.Error{
				Err:      ErrConstSeveralValues,
				Location: &scope.Location,
				Meta:     &constant.Meta,
			}
		}
	}

	valueCopy := *constant.Value
	valueCopy.TypeExpr = resolvedType

	return src.Const{
		Value: &valueCopy,
	}, nil
}

package analyzer

import (
	"errors"
	"fmt"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

var (
	ErrEmptyConst         = errors.New("Constant must either have value or reference to another constant")
	ErrEntityNotConst     = errors.New("Constant refers to an entity that is not constant")
	ErrResolveConstType   = errors.New("Cannot resolve constant type")
	ErrConstSeveralValues = errors.New("Constant cannot have several values at once")
)

func (a Analyzer) analyzeConst(cnst src.Const, scope src.Scope) (src.Const, error) {
	if cnst.Value == nil && cnst.Ref == nil {
		return src.Const{}, ErrEmptyConst
	}

	if cnst.Value == nil { // is ref
		entity, _, err := scope.Entity(*cnst.Ref)
		if err != nil {
			return src.Const{}, fmt.Errorf("entity: %w", err)
		}
		if entity.Kind != src.ConstEntity {
			return src.Const{}, fmt.Errorf("%w: %v", ErrEntityNotConst, entity.Kind)
		}
	}

	resolvedType, err := a.analyzeTypeExpr(cnst.Value.TypeExpr, scope)
	if err != nil {
		return src.Const{}, fmt.Errorf("%w: %v", ErrResolveConstType, err)
	}

	switch resolvedType.Inst.Ref.String() {
	case "bool":
		if cnst.Value.Int != 0 || cnst.Value.Float != 0 || cnst.Value.Str != "" {
			return src.Const{}, fmt.Errorf("%w: %v", ErrConstSeveralValues, cnst.Value)
		}
	case "int":
		if cnst.Value.Bool != false || cnst.Value.Float != 0 || cnst.Value.Str != "" {
			return src.Const{}, fmt.Errorf("%w: %v", ErrConstSeveralValues, cnst.Value)
		}
	case "float":
		if cnst.Value.Bool != false || cnst.Value.Int != 0 || cnst.Value.Str != "" {
			return src.Const{}, fmt.Errorf("%w: %v", ErrConstSeveralValues, cnst.Value)
		}
	case "str":
		if cnst.Value.Bool != false || cnst.Value.Int != 0 || cnst.Value.Float != 0 {
			return src.Const{}, fmt.Errorf("%w: %v", ErrConstSeveralValues, cnst.Value)
		}
	}

	valueCopy := *cnst.Value
	valueCopy.TypeExpr = resolvedType

	return src.Const{
		Value: &valueCopy,
	}, nil
}

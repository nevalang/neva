package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
)

var (
	ErrEmptyConst                  = errors.New("const must have value or reference to another const")
	ErrResolveConstType            = errors.New("cannot resolve constant type")
	ErrConstValuesOfDifferentTypes = errors.New("constant cannot have values of different types at once")
)

func (a Analyzer) analyzeConst(constant src.Const) (src.Const, error) {
	if constant.Value == nil && constant.Ref == nil {
		return src.Const{}, ErrEmptyConst
	}

	if constant.Value == nil {
		panic("// TODO: references for constants not implemented yet")
	}

	resolvedType, err := a.resolveTypeExpr(constant.Value.TypeExpr)
	if err != nil {
		return src.Const{}, fmt.Errorf("%w: %v", ErrResolveConstType, err)
	}

	switch resolvedType.Inst.Ref {
	case "bool":
		if constant.Value.Int != 0 || constant.Value.Float != 0 || constant.Value.Str != "" {
			return src.Const{}, fmt.Errorf("%w: %v", ErrConstValuesOfDifferentTypes, constant.Value)
		}
	case "int":
		if constant.Value.Bool != false || constant.Value.Float != 0 || constant.Value.Str != "" {
			return src.Const{}, fmt.Errorf("%w: %v", ErrConstValuesOfDifferentTypes, constant.Value)
		}
	case "float":
		if constant.Value.Bool != false || constant.Value.Int != 0 || constant.Value.Str != "" {
			return src.Const{}, fmt.Errorf("%w: %v", ErrConstValuesOfDifferentTypes, constant.Value)
		}
	case "str":
		if constant.Value.Bool != false || constant.Value.Int != 0 || constant.Value.Float != 0 {
			return src.Const{}, fmt.Errorf("%w: %v", ErrConstValuesOfDifferentTypes, constant.Value)
		}
	}

	valueCopy := *constant.Value
	valueCopy.TypeExpr = resolvedType

	return src.Const{
		Value: &valueCopy,
	}, nil
}

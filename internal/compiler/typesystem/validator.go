package typesystem

import (
	"errors"
	"fmt"
)

type Validator struct{}

var (
	ErrExprMustBeInstOrLit = errors.New("expr must be ether literal or instantiation, not both and not neither")
	ErrParamDuplicate      = errors.New("params must have unique names")
	ErrParams              = errors.New("bad params")
	ErrUnionTag            = errors.New("union tag must be non-empty")
	ErrUnionTagType        = errors.New("union tag type must be valid")
)

// ValidateDef makes sure that type supports recursion only if it's base type and that parameters are valid
func (v Validator) ValidateDef(def Def) error {
	if err := v.CheckParamUnique(def.Params); err != nil {
		return errors.Join(ErrParams, err)
	}
	return nil
}

// CheckParamUnique doesn't validate constraints, only ensures uniqueness
func (v Validator) CheckParamUnique(params []Param) error {
	m := make(map[string]struct{}, len(params))
	for _, param := range params {
		if _, ok := m[param.Name]; ok {
			return fmt.Errorf("%w: param", ErrParamDuplicate)
		}
	}
	return nil
}

// Validate makes shallow validation of expr.
// It checks that it's inst or literal, not both and not neither; All insts are valid by default;
func (v Validator) Validate(expr Expr) error {
	if expr.Lit.Empty() == (expr.Inst == nil) {
		return ErrExprMustBeInstOrLit
	}

	if expr.Inst != nil || expr.Lit.Type() == StructLitType {
		return nil
	}

	if expr.Lit.Type() == UnionLitType {
		for tag, tagExpr := range expr.Lit.Union {
			if tagExpr == nil {
				continue
			}
			if err := v.Validate(*tagExpr); err != nil {
				return fmt.Errorf("%s: invalid type for tag %s: %w", ErrUnionTagType.Error(), tag, err)
			}
		}
	}

	return nil
}

package typesystem

import (
	"errors"
	"fmt"
)

type Validator struct{}

var (
	ErrExprMustBeInstOrLit = errors.New("expr must be ether literal or instantiation, not both and not neither")
	ErrUnknownLit          = errors.New("expr literal must be known")
	ErrArrSize             = errors.New("arr size must be >= 2")
	ErrArrLitKind          = errors.New("array literal must have no enum, union or struct")
	ErrUnionLitKind        = errors.New("union literal must have no enum, array or struct")
	ErrEnumLitKind         = errors.New("enum literal must have no union, array or struct")
	ErrEnumLen             = errors.New("enum len must be >= 2")
	ErrUnionLen            = errors.New("union len must be >= 1")
	ErrEnumDupl            = errors.New("enum contains duplicate elements")
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
// Arr, union and enum must have size >= 2; Enum must have no duplicate elements.
func (v Validator) Validate(expr Expr) error {
	if expr.Lit.Empty() == (expr.Inst == nil) {
		return ErrExprMustBeInstOrLit
	}

	if expr.Inst != nil || expr.Lit.Type() == StructLitType {
		return nil
	}

	switch expr.Lit.Type() {
	case UnionLitType:
		var (
			hasTagOnlyMembers  bool
			hasMembersWithData bool
		)
		for tag, tagExpr := range expr.Lit.Union {
			if tagExpr == nil {
				hasTagOnlyMembers = true
				continue
			}
			if err := v.Validate(*tagExpr); err != nil {
				return fmt.Errorf("%w: invalid type for tag %s: %v", ErrUnionTagType, tag, err)
			}
			hasMembersWithData = true
		}
		if hasTagOnlyMembers && hasMembersWithData {
			return fmt.Errorf("%w: all union memers must either have data or not", ErrUnionTag)
		}
	}

	return nil
}

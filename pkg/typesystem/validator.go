package typesystem

import (
	"errors"
	"fmt"
)

type Validator struct{}

var (
	ErrExprMustBeInstOrLit          = errors.New("expr must be ether literal or instantiation, not both and not neither")
	ErrUnknownLit                   = errors.New("expr literal must be known")
	ErrArrSize                      = errors.New("arr size must be >= 2")
	ErrArrLitKind                   = errors.New("array literal must have no enum, union or record")
	ErrUnionLitKind                 = errors.New("union literal must have no enum, array or record")
	ErrEnumLitKind                  = errors.New("enum literal must have no union, array or record")
	ErrEnumLen                      = errors.New("enum len must be >= 2")
	ErrUnionLen                     = errors.New("union len must be >= 2")
	ErrEnumDupl                     = errors.New("enum contains duplicate elements")
	ErrNotBaseTypeSupportsRecursion = errors.New("only base type definitions can have support for recursion")
	ErrParamDuplicate               = errors.New("params must have unique names")
	ErrParams                       = errors.New("bad params")
)

// ValidateDef makes sure that type supports recursion only if it's base type and that parameters are valid
func (v Validator) ValidateDef(def Def) error {
	if def.CanBeUsedForRecursiveDefinitions && def.BodyExpr != nil {
		return fmt.Errorf("%w: %v", ErrNotBaseTypeSupportsRecursion, def)
	}
	if err := v.ValidateParams(def.Params); err != nil {
		return errors.Join(ErrParams, err)
	}
	return nil
}

// ValidateParams doesn't validate constraints, only ensures uniqueness
func (v Validator) ValidateParams(params []Param) error {
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

	if expr.Inst != nil || expr.Lit.Type() == RecLitType {
		return nil
	}

	switch expr.Lit.Type() { // by now we know it's not empty literal
	case ArrLitType:
		if expr.Lit.Arr.Size < 2 {
			return fmt.Errorf("%w: got %d", ErrArrSize, expr.Lit.Arr.Size)
		}
		switch {
		case expr.Lit.Enum != nil, expr.Lit.Rec != nil, expr.Lit.Union != nil:
			return ErrArrLitKind
		}
	case UnionLitType:
		if l := len(expr.Lit.Union); l < 2 {
			return fmt.Errorf("%w: got %d", ErrUnionLen, l)
		}
		switch {
		case expr.Lit.Enum != nil, expr.Lit.Rec != nil, expr.Lit.Arr != nil:
			return ErrUnionLitKind
		}
	case EnumLitType:
		if l := len(expr.Lit.Enum); l < 2 {
			return fmt.Errorf("%w: got %d", ErrEnumLen, l)
		}
		set := make(map[string]struct{}, len(expr.Lit.Enum))
		for _, el := range expr.Lit.Enum { // look for duplicate
			if _, ok := set[el]; ok {
				return fmt.Errorf("%w: %s", ErrEnumDupl, el)
			}
			set[el] = struct{}{}
		}
		switch {
		case expr.Lit.Union != nil, expr.Lit.Rec != nil, expr.Lit.Arr != nil:
			return ErrEnumLitKind
		}
	}

	return nil
}

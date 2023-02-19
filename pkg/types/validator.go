package types

import (
	"errors"
	"fmt"
)

type Validator struct{}

var (
	ErrInvalidExprType              = errors.New("expr must be ether literal or instantiation, not both and not neither")
	ErrUnknownLit                   = errors.New("expr literal must be known")
	ErrArrSize                      = errors.New("arr size must be >= 2")
	ErrEnumLen                      = errors.New("enum len must be >= 2")
	ErrUnionLen                     = errors.New("union len must be >= 2")
	ErrEnumDupl                     = errors.New("enum contains duplicate elements")
	ErrNotBaseTypeSupportsRecursion = errors.New("only base type definitions can have support for recursion")
)

// ValidateDef checks that type doesn't support recursion if it's not a base type (it's body isn't empty)
// and that there're no duplicate names among type parameters. It doesn't validate body or constrs
// because there's ValidateExpr that called by Resolver already.
func (v Validator) ValidateDef(def Def) error {
	if def.IsRecursionAllowed && !def.BodyExpr.Empty() {
		return fmt.Errorf("%w: %v", ErrNotBaseTypeSupportsRecursion, def)
	}

	m := make(map[string]struct{}, len(def.Params))
	for _, param := range def.Params {
		if _, ok := m[param.Name]; ok {
			return errors.New("params must have unique names")
		}
	}

	return nil
}

// Validate makes shallow validation of expr.
// It checks that expr is inst or literal, not both, not neither.
// All insts are valid by default.
// Arr, union and enum must have size >= 2. Enum must have no duplicate elements.
func (v Validator) Validate(expr Expr) error {
	// FIXME empty expr considered invalid but native types has empty exprs as body which means their body invalid
	if expr.Lit.Empty() == expr.Inst.Empty() { // we don't use expr.Empty() because constraint can be empty
		return ErrInvalidExprType
	}

	if expr.Lit.Empty() || expr.Lit.Type() == RecLitType {
		return nil
	}

	switch expr.Lit.Type() { // by now we know it's not empty literal
	case ArrLitType:
		if expr.Lit.Arr.Size < 2 {
			return fmt.Errorf("%w: got %d", ErrArrSize, expr.Lit.Arr.Size)
		}
	case UnionLitType:
		if l := len(expr.Lit.Union); l < 2 {
			return fmt.Errorf("%w: got %d", ErrUnionLen, l)
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
	}

	return nil
}

package types

import "fmt"

type Validator struct{}

// Checks that expr is either an instantiation or a literal.
// If literal is empty then it's an instantiation. All insts are valid because it's ok to have "" type with 0 args.
// For arr, union and enum lits it checks that their size is >= 2. For enum it ensures there no duplicate elements.
func (v Validator) Validate(expr Expr) error {
	if expr.Lit.Empty() { // if it's inst
		return nil // then nothing to validate, resolving needed
	} // by now we know it's not empty lit

	if expr.Inst.Ref != "" || len(expr.Inst.Args) != 0 { // must not be both lit and inst
		return ErrInvalidExprType
	}

	// we don't check recs and empty lits
	switch expr.Lit.Type() { // because we know lit isn't empty and because empty recs are fine
	case ArrLitType:
		if expr.Lit.ArrLit.Size < 2 {
			return fmt.Errorf("%w: got %d", ErrArrSize, expr.Lit.ArrLit.Size)
		}
	case UnionLitType:
		if l := len(expr.Lit.UnionLit); l < 2 {
			return fmt.Errorf("%w: got %d", ErrUnionLen, l)
		}
	case EnumLitType:
		if l := len(expr.Lit.EnumLit); l < 2 {
			return fmt.Errorf("%w: got %d", ErrEnumLen, l)
		}
		set := make(map[string]struct{}, len(expr.Lit.EnumLit))
		for _, el := range expr.Lit.EnumLit { // look for duplicate
			if _, ok := set[el]; ok {
				return fmt.Errorf("%w: %s", ErrEnumDupl, el)
			}
			set[el] = struct{}{}
		}
	}

	return nil // valid lit
}

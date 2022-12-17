package types

import (
	"errors"
	"fmt"
)

var (
	ErrDiffTypes    = errors.New("expr and constraint must both be literals or instantiations except if constr is union")
	ErrDiffRefs     = errors.New("expr instantiation must have same ref as constraint")
	ErrArgsLen      = errors.New("expr instantiation must have >= args than constraint")
	ErrArg          = errors.New("expr arg must be subtype of it's corresponding constraint arg")
	ErrLitNotArr    = errors.New("expr is literal but not arr")
	ErrLitArrSize   = errors.New("expr arr size must be >= constraint")
	ErrArrDiffType  = errors.New("expr arr must have same type as constraint")
	ErrLitNotEnum   = errors.New("expr is literal but not enum")
	ErrBigEnum      = errors.New("expr enum must be <= constraint enum")
	ErrEnumEl       = errors.New("expr enum el doesn't match constraint")
	ErrLitNotRec    = errors.New("expr is lit but not rec")
	ErrRecLen       = errors.New("expr record must contain >= fields than constraint")
	ErrRecField     = errors.New("expr rec field must be subtype of corresponding constraint field")
	ErrRecNoField   = errors.New("expr rec is missing field of constraint")
	ErrUnion        = errors.New("expr must be subtype of constraint union")
	ErrUnionsLen    = errors.New("expr union must be <= constraint union")
	ErrUnions       = errors.New("expr union el must be subtype of constraint union")
	ErrInvariant    = errors.New("expr's invariant is broken")
	ErrDiffLitTypes = errors.New("expr and constraint literals must be of the same type")
)

// Both expression and constraint must be resolved
func (expr Expr) IsSubTypeOf(constr Expr) error { //nolint:funlen,gocognit,gocyclo
	if expr.Lit.Empty() != constr.Lit.Empty() && constr.Lit.Type() != UnionLitType { // expr can be inst if constr is union
		return fmt.Errorf("%w: expr %v, constaint %v", ErrDiffTypes, expr.Lit, constr.Lit)
	}

	if constr.Lit.Empty() { // expr and constr insts
		if expr.Inst.Ref != constr.Inst.Ref {
			return fmt.Errorf(
				"%w: got %v, want %v", ErrDiffRefs, expr.Inst.Ref, constr.Inst.Ref,
			)
		}
		if len(expr.Inst.Args) < len(constr.Inst.Args) {
			return fmt.Errorf(
				"%w: got %v, want %v", ErrArgsLen, len(expr.Inst.Args), len(constr.Inst.Args),
			)
		}
		for i, constraintArg := range constr.Inst.Args {
			if err := constraintArg.IsSubTypeOf(expr.Inst.Args[i]); err != nil {
				return fmt.Errorf("%w: #%d, got %v, want %v", ErrArg, i, constraintArg, expr.Inst.Args[i])
			}
		}
		return nil
	} // we know constr is lit by now

	exprLitType := expr.Lit.Type()
	constrLitType := constr.Lit.Type()
	if constrLitType != UnionLitType && exprLitType != constrLitType { // if it's not union, expr must be same lit
		return fmt.Errorf("%w: got %v, want %v", ErrDiffLitTypes, exprLitType, constrLitType)
	}

	switch constrLitType { // it's union constr and whatever expr
	case ArrLitType: // [5]int <: [4]int|float ???
		if expr.Lit.ArrLit.Size < constr.Lit.ArrLit.Size {
			return fmt.Errorf(
				"%w: got %d, want %d", ErrLitArrSize, expr.Lit.ArrLit.Size, constr.Lit.ArrLit.Size,
			)
		}
		if err := expr.Lit.ArrLit.Expr.IsSubTypeOf(constr.Lit.ArrLit.Expr); err != nil {
			return fmt.Errorf("%w: %v", ErrArrDiffType, err)
		}
	case EnumLitType: // {a b c} <: {a b c d}
		if len(expr.Lit.EnumLit) > len(constr.Lit.EnumLit) {
			return fmt.Errorf(
				"%w: got %d, want %d", ErrBigEnum, len(expr.Lit.EnumLit), len(constr.Lit.EnumLit),
			)
		}
		for i, exprEl := range expr.Lit.EnumLit {
			if exprEl != constr.Lit.EnumLit[i] {
				return fmt.Errorf("%w: #%d got %s, want %s", ErrEnumEl, i, exprEl, constr.Lit.EnumLit[i])
			}
		}
	case RecLitType: // {x int, y float} <: {x int|str}
		if len(expr.Lit.RecLit) < len(constr.Lit.RecLit) {
			return fmt.Errorf(
				"%w: got %v, want %v", ErrRecLen, len(expr.Lit.RecLit), len(constr.Lit.RecLit),
			)
		}
		for constraintFieldName, constraintField := range constr.Lit.RecLit {
			exprField, ok := expr.Lit.RecLit[constraintFieldName]
			if !ok {
				return fmt.Errorf("%w: %v", ErrRecNoField, constraintFieldName)
			}
			if err := exprField.IsSubTypeOf(constraintField); err != nil {
				return fmt.Errorf("%w: field '%s': %v", ErrRecField, constraintFieldName, err)
			}
		}
	case UnionLitType: // 1) int <: str | int 2) int | str <: str | bool | int
		if expr.Lit.UnionLit == nil { // constraint is union, expr is not
			for _, constraintUnionEl := range constr.Lit.UnionLit {
				if expr.IsSubTypeOf(constraintUnionEl) == nil {
					return nil
				}
			}
			return fmt.Errorf("%w: got %v", ErrUnion, expr.Lit)
		}
		if len(expr.Lit.UnionLit) > len(constr.Lit.UnionLit) {
			return fmt.Errorf(
				"%w: got %d, want %d", ErrUnionsLen, len(expr.Lit.UnionLit), len(constr.Lit.UnionLit),
			)
		}
		for _, exprEl := range expr.Lit.UnionLit {
			var b bool
			for _, constraintEl := range constr.Lit.UnionLit {
				if exprEl.IsSubTypeOf(constraintEl) == nil {
					b = true
					break
				}
			}
			if !b {
				return fmt.Errorf(
					"%w: got %v, want %v",
					ErrUnions, exprEl, constr.Lit.UnionLit,
				)
			}
		}
	}

	return nil
}

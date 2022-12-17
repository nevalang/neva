package types

import (
	"errors"
	"fmt"
)

var (
	ErrDiffTypes   = errors.New("expr and constraint must both be literals or instantiations")
	ErrDiffRefs    = errors.New("expr instantiation must have same ref as constraint")
	ErrArgsLen     = errors.New("expr instantiation must have >= args than constraint")
	ErrArg         = errors.New("expr arg must be subtype of it's corresponding constraint arg")
	ErrLitNotArr   = errors.New("expr is literal but not arr")
	ErrLitArrSize  = errors.New("expr arr size must be >= constraint")
	ErrArrDiffType = errors.New("expr arr must have same type as constraint")
	ExprLitNotEnum = errors.New("expr literal must be enum")
	ExprLitBigEnum = errors.New("expr enum must be <= constraint enum")
	ErrEnumEl      = errors.New("expr enum el doesn't match constraint")
	ErrLitNotRec   = errors.New("expr is lit but not rec")
	ErrRecLen      = errors.New("expr record must contain >= fields than constraint")
	ErrRecField    = errors.New("expr rec field must be subtype of corresponding constraint field")
	ErrRecNoField  = errors.New("expr rec is missing field of constraint")
	ErrUnion       = errors.New("expr must be subtype of constraint union")
	ErrUnionsLen   = errors.New("expr union must be <= constraint union")
	ErrUnions      = errors.New("expr union el must be subtype of constraint union")
	ErrInvariant   = errors.New("expr's invariant is broken")
)

// Both expression and constraint must be resolved
func (expr Expr) IsSubTypeOf(constraint Expr) error { //nolint:funlen,gocognit
	isExprLit := !expr.Lit.Empty()
	isConstraintLit := !constraint.Lit.Empty()

	if isExprLit != isConstraintLit {
		return fmt.Errorf("%w: expr %v, constaint %v", ErrDiffTypes, expr.Lit, constraint.Lit)
	}

	if !isConstraintLit {
		if expr.Inst.Ref != constraint.Inst.Ref {
			return fmt.Errorf(
				"%w: got %v, want %v", ErrDiffRefs, expr.Inst.Ref, constraint.Inst.Ref,
			)
		}

		if len(expr.Inst.Args) < len(constraint.Inst.Args) {
			return fmt.Errorf(
				"%w: got %v, want %v", ErrArgsLen, len(expr.Inst.Args), len(constraint.Inst.Args),
			)
		}

		for i, constraintArg := range constraint.Inst.Args {
			if err := constraintArg.IsSubTypeOf(expr.Inst.Args[i]); err != nil {
				return fmt.Errorf("%w: #%d, got %v, want %v", ErrArg, i, constraintArg, expr.Inst.Args[i])
			}
		}
	}

	switch {
	case constraint.Lit.ArrLit != nil: // [5]int <: [4]int|float ???
		if expr.Lit.ArrLit == nil {
			return fmt.Errorf("%w: got %v", ErrLitNotArr, expr.Lit)
		}
		if expr.Lit.ArrLit.Size < constraint.Lit.ArrLit.Size {
			return fmt.Errorf(
				"%w: got %d, want %d", ErrLitArrSize, expr.Lit.ArrLit.Size, constraint.Lit.ArrLit.Size,
			)
		}
		if err := expr.Lit.ArrLit.Expr.IsSubTypeOf(constraint.Lit.ArrLit.Expr); err != nil {
			return fmt.Errorf("%w: %v", ErrArrDiffType, err)
		}
	case constraint.Lit.EnumLit != nil: // {a b c} <: {a b c d}
		if expr.Lit.EnumLit == nil {
			return fmt.Errorf("%w: got %v", ExprLitNotEnum, expr.Lit)
		}
		if len(expr.Lit.EnumLit) > len(constraint.Lit.EnumLit) {
			return fmt.Errorf(
				"%w: got %d, want %d", ExprLitBigEnum, len(expr.Lit.EnumLit), len(constraint.Lit.EnumLit),
			)
		}
		for i, exprEl := range expr.Lit.EnumLit {
			if exprEl != constraint.Lit.EnumLit[i] {
				return fmt.Errorf("%w: #%d got %s, want %s", ErrEnumEl, i, exprEl, constraint.Lit.EnumLit[i])
			}
		}
	case constraint.Lit.RecLit != nil: // {x int, y float} <: {x int|str}
		if expr.Lit.RecLit == nil {
			return fmt.Errorf("%w: %v", ErrLitNotRec, expr.Lit)
		}
		if len(expr.Lit.RecLit) < len(constraint.Lit.RecLit) {
			return fmt.Errorf(
				"%w: got %v, want %v", ErrRecLen, len(expr.Lit.RecLit), len(constraint.Lit.RecLit),
			)
		}
		for constraintFieldName, constraintField := range constraint.Lit.RecLit {
			exprField, ok := expr.Lit.RecLit[constraintFieldName]
			if !ok {
				return fmt.Errorf("%w: %v", ErrRecNoField, constraintFieldName)
			}
			if err := exprField.IsSubTypeOf(constraintField); err != nil {
				return fmt.Errorf("%w: field '%s': %v", ErrRecField, constraintFieldName, err)
			}
		}
	case constraint.Lit.UnionLit != nil: // 1) int <: str | int 2) int | str <: str | bool | int
		if expr.Lit.UnionLit == nil { // constraint is union, expr is not
			for _, constraintUnionEl := range constraint.Lit.UnionLit {
				if expr.IsSubTypeOf(constraintUnionEl) == nil {
					return nil
				}
			}
			return fmt.Errorf("%w: got %v", ErrUnion, expr.Lit)
		}
		if len(expr.Lit.UnionLit) > len(constraint.Lit.UnionLit) {
			return fmt.Errorf(
				"%w: got %d, want %d", ErrUnionsLen, len(expr.Lit.UnionLit), len(constraint.Lit.UnionLit),
			)
		}
		for _, exprEl := range expr.Lit.UnionLit {
			var b bool
			for _, constraintEl := range constraint.Lit.UnionLit {
				if exprEl.IsSubTypeOf(constraintEl) == nil {
					b = true
					break
				}
			}
			if !b {
				return fmt.Errorf(
					"%w: got %v, want %v",
					ErrUnions, exprEl, constraint.Lit.UnionLit,
				)
			}
		}
	}

	return ErrInvariant
}

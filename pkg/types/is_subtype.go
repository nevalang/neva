package types

import (
	"fmt"
)

// // Compare checks whether two expressions have compatible types.
// func Compare(subTyp, typ Expr) error {
// return nil
// }

// Both expressions must be resolved.
func (expr Expr) IsSubType(constraint Expr) error { //nolint:funlen,gocognit
	isExprLit := !expr.Literal.Empty()
	isConstraintLit := !constraint.Literal.Empty()

	if isExprLit != isConstraintLit {
		return fmt.Errorf(
			"expr and constraint must both be literals or instantiations: expr %v, constaint %v",
			expr, constraint,
		)
	}

	if !isConstraintLit {
		if expr.Instantiation.Ref != constraint.Instantiation.Ref {
			return fmt.Errorf(
				"expr must have same ref type as constraint: got %v, want %v",
				expr.Instantiation.Ref, constraint.Instantiation.Ref,
			)
		}

		if len(expr.Instantiation.Args) != len(constraint.Instantiation.Args) {
			// точно == а не <>?
			return fmt.Errorf("...")
		}

		for i, exprArg := range expr.Instantiation.Args {
			if err := exprArg.IsSubType(constraint.Instantiation.Args[i]); err != nil {
				return err
			}
		}
	}

	switch {
	case constraint.Literal.ArrLir != nil: // [5]int <: [4]int|float ???
		if expr.Literal.ArrLir == nil {
			return fmt.Errorf("expr must be array, got %v", expr)
		}
		if expr.Literal.ArrLir.Size < constraint.Literal.ArrLir.Size {
			return fmt.Errorf(
				"expr array size must be >= constraint array size: got %d, want %d",
				expr.Literal.ArrLir.Size, constraint.Literal.ArrLir.Size,
			)
		}
		if err := expr.Literal.ArrLir.Expr.IsSubType(constraint.Literal.ArrLir.Expr); err != nil {
			return fmt.Errorf("expr array type must be subtype of constraint array type: %w", err)
		}
	case constraint.Literal.EnumLit != nil: // {a b c} <: {a b c d}
		if expr.Literal.EnumLit == nil {
			return fmt.Errorf("expr must be enum, got %v", expr)
		}
		if len(expr.Literal.EnumLit) > len(constraint.Literal.EnumLit) {
			return fmt.Errorf("expr enum must be <= constraint enum")
		}
		for i, exprEl := range expr.Literal.EnumLit {
			if exprEl != constraint.Literal.EnumLit[i] {
				return fmt.Errorf(
					"expr enum el #%d doesn't match constraint: got %s, want %s",
					i, exprEl, constraint.Literal.EnumLit[i],
				)
			}
		}
	case constraint.Literal.RecLit != nil: // {x int, y float} <: {x int|str}
		if expr.Literal.RecLit == nil {
			return fmt.Errorf("expr must be record, got %v", expr)
		}
		if len(expr.Literal.RecLit) < len(constraint.Literal.RecLit) {
			return fmt.Errorf(
				"expr record must be >= constraint record: got %v, want %v",
				len(expr.Literal.RecLit), len(constraint.Literal.RecLit),
			)
		}
		for name, field := range constraint.Literal.RecLit {
			if err := expr.Literal.RecLit[name].IsSubType(field); err != nil {
				return fmt.Errorf("expr record field %s is not a sub type of constraint", err)
			}
		}
	case constraint.Literal.UnionLit != nil: // 1) int <: str | int 2) int | str <: str | bool | int
		if expr.Literal.UnionLit == nil {
			for _, el := range constraint.Literal.UnionLit {
				if expr.IsSubType(el) == nil {
					return nil
				}
			}
		}
		if len(expr.Literal.UnionLit) > len(constraint.Literal.UnionLit) {
			return fmt.Errorf(
				"expr union must be <= constraint union: got %d, want %d",
				len(expr.Literal.UnionLit), len(constraint.Literal.UnionLit),
			)
		}
		for _, exprEl := range expr.Literal.UnionLit {
			var b bool
			for _, constraintEl := range constraint.Literal.UnionLit {
				if exprEl.IsSubType(constraintEl) == nil {
					b = true
					break
				}
			}
			if !b {
				return fmt.Errorf(
					"expr union el must be subtype of constraint union: got %v, want %v",
					exprEl, constraint.Literal.UnionLit,
				)
			}
		}
	}

	return nil
}

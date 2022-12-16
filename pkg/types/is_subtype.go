package types

import (
	"fmt"
)

// // Compare checks whether two expressions have compatible types.
// func Compare(subTyp, typ Expr) error {
// return nil
// }

// Both expressions must be resolved.
func (expr Expr) IsSubType(constraint Expr) error {
	switch {
	case constraint.Literal.ArrLir != nil: // [5]int <: [4]int|float
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
	case constraint.Literal.EnumLit != nil:
		if expr.Literal.EnumLit == nil {
			return fmt.Errorf("expr must be enum, got %v", expr)
		}
		if len(expr.Literal.EnumLit) < len(constraint.Literal.EnumLit) {
			return fmt.Errorf("expr enum must be >= constraint enum")
		}
		for i, el := range constraint.Literal.EnumLit {
			if el != expr.Literal.EnumLit[i] {
				return fmt.Errorf(
					"expr enum el #%d doesn't match constraint: got %s, want %s",
					i, expr.Literal.EnumLit[i], el,
				)
			}
		}
	case constraint.Literal.RecLit != nil:
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
	case constraint.Literal.UnionLit != nil:
		if expr.Literal.UnionLit == nil {
			if len(expr.Literal.UnionLit) > len(constraint.Literal.UnionLit) {
				return fmt.Errorf("expr union must be <= constraint union: got %d, want %d",
					len(expr.Literal.UnionLit), len(constraint.Literal.UnionLit),
				)
			}
		}

	}

	return nil
}

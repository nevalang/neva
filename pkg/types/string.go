package types

import "fmt"

// String formats expression in a TS manner
func (expr Expr) String() string { // todo move?
	var s string

	switch expr.Lit.Type() {
	case ArrLitType:
		return fmt.Sprintf(
			"[%d]%s",
			expr.Lit.ArrLit.Size, expr.Lit.ArrLit.Expr.String(),
		)
	case EnumLitType:
		s += "{"
		for i, el := range expr.Lit.EnumLit {
			s += " " + el
			if i == len(expr.Lit.EnumLit)-1 {
				s += " "
			}
		}
		return s + "}"
	case RecLitType:
		s += "{"
		c := 0
		for fieldName, fieldExpr := range expr.Lit.RecLit {
			s += fmt.Sprintf(" %s %s", fieldName, fieldExpr)
			if c < len(expr.Lit.RecLit)-1 {
				s += ","
			} else {
				s += " "
			}
			c++
		}
		return s + "}"
	case UnionLitType:
		for i, el := range expr.Lit.UnionLit {
			s += el.String()
			if i < len(expr.Lit.UnionLit)-1 {
				s += " | "
			}
		}
		return s
	}

	if len(expr.Inst.Args) == 0 {
		if expr.Inst.Args != nil {
			return expr.Inst.Ref + "<>"
		}
		return expr.Inst.Ref
	}

	s = expr.Inst.Ref + "<"
	for i, arg := range expr.Inst.Args {
		s += arg.String()
		if i < len(expr.Inst.Args)-1 {
			s += ", "
		}
	}
	s += ">"

	return s
}

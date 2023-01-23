package types

import "fmt"

// String formats expression in a TS manner
func (expr Expr) String() string { // todo move?
	var s string

	switch expr.Lit.Type() {
	case ArrLitType:
		return fmt.Sprintf(
			"[%d]%s",
			expr.Lit.Arr.Size, expr.Lit.Arr.Expr.String(),
		)
	case EnumLitType:
		s += "{"
		for i, el := range expr.Lit.Enum {
			s += " " + el
			if i == len(expr.Lit.Enum)-1 {
				s += " "
			}
		}
		return s + "}"
	case RecLitType:
		s += "{"
		c := 0
		for fieldName, fieldExpr := range expr.Lit.Rec {
			s += fmt.Sprintf(" %s %s", fieldName, fieldExpr)
			if c < len(expr.Lit.Rec)-1 {
				s += ","
			} else {
				s += " "
			}
			c++
		}
		return s + "}"
	case UnionLitType:
		for i, el := range expr.Lit.Union {
			s += el.String()
			if i < len(expr.Lit.Union)-1 {
				s += " | "
			}
		}
		return s
	}

	if len(expr.Inst.Args) == 0 {
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

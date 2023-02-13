package types

import "fmt"

// String formats expression in a TS manner
func (expr Expr) String() string { // todo move?
	var str string

	switch expr.Lit.Type() {
	case EmptyLitType:
		return "unknown"
	case ArrLitType:
		return fmt.Sprintf(
			"[%d]%s",
			expr.Lit.Arr.Size, expr.Lit.Arr.Expr.String(),
		)
	case EnumLitType:
		str += "{"
		for i, el := range expr.Lit.Enum {
			str += " " + el
			if i == len(expr.Lit.Enum)-1 {
				str += " "
			} else {
				str += ","
			}
		}
		return str + "}"
	case RecLitType:
		str += "{"
		count := 0
		for fieldName, fieldExpr := range expr.Lit.Rec {
			str += fmt.Sprintf(" %s %s", fieldName, fieldExpr)
			if count < len(expr.Lit.Rec)-1 {
				str += ","
			} else {
				str += " "
			}
			count++
		}
		return str + "}"
	case UnionLitType:
		for i, el := range expr.Lit.Union {
			str += el.String()
			if i < len(expr.Lit.Union)-1 {
				str += " | "
			}
		}
		return str
	}

	if len(expr.Inst.Args) == 0 {
		return expr.Inst.Ref
	}

	str = expr.Inst.Ref + "<"
	for i, arg := range expr.Inst.Args {
		str += arg.String()
		if i < len(expr.Inst.Args)-1 {
			str += ", "
		}
	}
	str += ">"

	return str
}

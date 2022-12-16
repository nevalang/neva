package types

// Выводы:
// Разрешив union только как constraint
// встретив юнион, мы будем вынуждены рекурсивно проверять
// точно ли текущий expr это constraint или expr вложенный в constraint

func (expr Expr) String() string {
	var s string

	if expr.Literal.RecLit != nil {
		s += "{"
		for fieldName, fieldExpr := range expr.Literal.RecLit {
			s += " " + fieldName + ": " + fieldExpr.String() + " "
		}
		s += "}"
		return s
	}

	if len(expr.Instantiation.Args) == 0 {
		return expr.Instantiation.Ref
	}

	s = expr.Instantiation.Ref + "<"
	for i, arg := range expr.Instantiation.Args {
		s += arg.String()
		if i < len(expr.Instantiation.Args)-1 {
			s += ", "
		}
	}
	s += ">"

	return s
}

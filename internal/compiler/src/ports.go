package src

type IO struct {
	In, Out Ports
}

type Ports map[string]PortDef

type PortDef struct {
	TypeExpr    TypeExpr
	IsArray     bool
	Description string
}

func NewPort(expr TypeExpr, isArray bool, desc string) PortDef {
	return PortDef{
		TypeExpr:    expr,
		IsArray:     isArray,
		Description: desc,
	}
}

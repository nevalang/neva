package src

type MsgDef struct {
	TypeExpr TypeExpr
	Int      int64
	Float    float64
	List     []MsgDef
	Dict     map[string]MsgDef
}

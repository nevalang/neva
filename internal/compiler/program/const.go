package program

type Const interface {
	Type() DataType
	Int() int
	Str() string
	Bool() bool
}

type emptyConst struct{}

func (emptyConst) Int() int    { return 0 }
func (emptyConst) Str() string { return "" }
func (emptyConst) Bool() bool  { return false }

type IntConst struct {
	emptyConst
	v int
}

func (i IntConst) Type() DataType { return TypeInt }
func (i IntConst) Int() int       { return i.v }

func NewIntConst(v int) IntConst {
	return IntConst{
		v:          v,
		emptyConst: emptyConst{},
	}
}

type StrConst struct {
	emptyConst
	v string
}

func (s StrConst) Type() DataType { return TypeStr }
func (s StrConst) Str() string    { return s.v }

func NewStrConst(s string) StrConst {
	return StrConst{
		v:          s,
		emptyConst: emptyConst{},
	}
}

type BoolConst struct {
	emptyConst
	v bool
}

func (b BoolConst) Type() DataType { return TypeBool }
func (b BoolConst) Bool() bool     { return b.v }

func NewBoolConst(b bool) BoolConst {
	return BoolConst{
		v:          b,
		emptyConst: emptyConst{},
	}
}

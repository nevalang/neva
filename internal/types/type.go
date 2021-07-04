package types

type Type uint8

func (t Type) String() string {
	return tn[t]
}

const (
	Unknown Type = iota
	Int
	Str
	Bool
)

type typeNames map[Type]string

var (
	tn = typeNames{
		Unknown: "unknown",
		Int:     "int",
		Str:     "str",
		Bool:    "bool",
	}
)

func ByName(s string) Type {
	for k, v := range tn {
		if v == s {
			return k
		}
	}
	return Unknown
}

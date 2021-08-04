package types

type Type uint8

func (t Type) String() string {
	switch t {
	case Int:
		return "int"
	case Str:
		return "str"
	case Bool:
		return "bool"
	}

	return "unknown"
}

const (
	Unknown Type = iota
	Int
	Str
	Bool
)

func ByName(s string) Type {
	switch s {
	case "int":
		return Int
	case "str":
		return Str
	case "bool":
		return Bool
	}

	return Unknown
}

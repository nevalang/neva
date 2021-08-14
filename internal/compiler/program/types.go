package program

type Type uint8

func (t Type) String() string {
	switch t {
	case Int:
		return "int"
	case Str:
		return "str"
	case Bool:
		return "bool"
	case Struct:
		return "struct"
	}

	return "unknown"
}

const (
	Unknown Type = iota
	Int
	Str
	Bool
	Struct
)

func ByName(s string) Type {
	switch s {
	case "int":
		return Int
	case "str":
		return Str
	case "bool":
		return Bool
	case "struct":
		return Struct
	}

	return Unknown
}

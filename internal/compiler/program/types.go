package program

type Type uint8

func (t Type) String() string {
	switch t {
	case IntType:
		return "int"
	case StrType:
		return "str"
	case BoolType:
		return "bool"
	}

	return "unknown"
}

const (
	UnknownType Type = iota
	IntType
	StrType
	BoolType
)

func TypeByName(name string) Type {
	switch name {
	case "int":
		return IntType
	case "str":
		return StrType
	case "bool":
		return BoolType
	}

	return UnknownType
}

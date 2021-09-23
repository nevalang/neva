package program

type TypeDescriptor struct {
	base   Type
	params []TypeDescriptor
}

type Type uint8

func (t Type) String() string {
	switch t {
	case IntType:
		return "int"
	case StrType:
		return "str"
	case BoolType:
		return "bool"
	case StructType:
		return "struct"
	}

	return "unknown"
}

const (
	UnknownType Type = iota
	IntType
	StrType
	BoolType
	StructType
)

func TypeByName(name string) Type {
	switch name {
	case "int":
		return IntType
	case "str":
		return StrType
	case "bool":
		return BoolType
	case "struct":
		return StructType
	}

	return UnknownType
}

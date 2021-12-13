package program

import "fmt"

type Type uint8

const (
	TypeInt Type = iota + 1
	TypeStr
	TypeBool
	TypeSig
)

func (t Type) String() string {
	var s typeName = "unknown"

	switch t {
	case TypeInt:
		s = intType
	case TypeStr:
		s = strType
	case TypeBool:
		s = boolType
	}

	return string(s)
}

type typeName string

const (
	intType  typeName = "int"
	strType  typeName = "str"
	boolType typeName = "bool"
	sigType  typeName = "sig"
)

func TypeByName(name string) (Type, error) {
	switch typeName(name) {
	case intType:
		return TypeInt, nil
	case strType:
		return TypeStr, nil
	case boolType:
		return TypeBool, nil
	}

	return 0, fmt.Errorf("unknown type %s", name)
}

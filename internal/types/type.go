package types

import "fmt"

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
	Int Type = iota + 1
	Str
	Bool
)

// TODO do nut return error
func ByName(s string) (Type, error) {
	switch s {
	case "int":
		return Int, nil
	case "str":
		return Str, nil
	case "bool":
		return Bool, nil
	}

	return 0, fmt.Errorf("no type has name '%s'", s)
}

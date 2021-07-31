package types

import "fmt"

type Type uint8

func (t Type) String() string {
	v, ok := tn[t]
	if !ok {
		return "unknown"
	}
	return v
}

const (
	Int Type = iota + 1
	Str
	Bool
)

type typeNames map[Type]string

var (
	tn = typeNames{
		Int:  "int",
		Str:  "str",
		Bool: "bool",
	}
)

func ByName(s string) (Type, error) {
	for k, v := range tn {
		if v == s {
			return k, nil
		}
	}
	return 0, fmt.Errorf("no type has name '%s'", s)
}

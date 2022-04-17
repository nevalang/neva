package main

func check(a, b TypeInstance) bool {
	if a.Type != b.Type || len(a.Args) != len(b.Args) {
		return false
	}

	for i, arg := range a.Args {
		if !check(arg, b.Args[i]) {
			return false
		}
	}

	return true
}

type TypeInstance struct {
	Type Type
	Args []TypeInstance
}

type Type uint8

const (
	Bool Type = iota + 1
	Int
	Struct
	List
	Map
)

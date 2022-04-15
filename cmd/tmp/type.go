package main

type TypeInstance struct {
	Type Type
	Args []TypeInstance
}

// type TypeInstance interface {
// 	Type() Type
// 	Args() []TypeInstance
// }

type Type uint8

const (
	Bool Type = iota + 1
	Int
	Struct
	List
	Map
)

type IntTypeInstance struct{}

func (IntTypeInstance) Type() Type           { return Int }
func (IntTypeInstance) Args() []TypeInstance { return nil }

// Complex type without args
type StructTypeInstance struct{ fields map[string]TypeInstance }

func (StructTypeInstance) Type() Type           { return Struct }
func (StructTypeInstance) Args() []TypeInstance { return nil }

// Complex type with args
type ListTypeInstance struct{ elType TypeInstance }

func (ListTypeInstance) Type() Type           { return List }
func (ListTypeInstance) Args() []TypeInstance { return nil }

func NewListTypeInstance(elType TypeInstance) (ListTypeInstance, error) {
}

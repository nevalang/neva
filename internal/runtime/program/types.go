package program

type Type uint8

const (
	UnknownType Type = iota
	IntType
	StrType
	BoolType
	StructType
)

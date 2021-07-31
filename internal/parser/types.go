package parser

import "github.com/emil14/stream/internal/types"

type typeNames map[types.Type]string

var (
	tn = typeNames{
		types.Int:  "int",
		types.Str:  "str",
		types.Bool: "bool",
	}
)

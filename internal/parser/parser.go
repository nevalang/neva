package parser

import "github.com/emil14/refactored-garbanzo/internal/runtime"

type Parser interface {
	Parse([]byte) (runtime.Module, error)
}

func NewParser(v Validator) Parser {
	return jsonParser{v}
}

package app

import (
	"github.com/emil14/refactored-garbanzo/internal/core"
	"github.com/emil14/refactored-garbanzo/internal/parser"
)

type System struct {
	p parser.Parser
	r core.Runtime
}

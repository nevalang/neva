package app

import (
	"github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/parser"
)

type System struct {
	p parser.Parser
	r core.Runtime
}

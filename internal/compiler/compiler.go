package compiler

import (
	"github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/runtime"
)

type compiler struct {
	p Parser
}

func (c compiler) Compile(prog map[string]core.Component) (runtime.Program, error) {
	return runtime.Program{}, nil
}

package decoder

import (
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/pkg/runtimesdk"
)

type caster struct{}

func (c caster) Cast(runtimesdk.Program) (runtime.Program, error) {
	return runtime.Program{}, nil // TODO
}

func NewCaster() caster {
	return caster{}
}

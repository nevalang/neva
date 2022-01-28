package parser

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/compiler"
)

type (
	Caster interface {
		Cast(Module) compiler.Module
	}
	Unmarshaler interface {
		Unmarshal([]byte) (Module, error)
	}
)

var ErrUnmarshaler = errors.New("unmarshaler")

type parser struct {
	unmarshaler Unmarshaler
	caster      Caster
}

func (p parser) Parse(mods map[string][]byte) (map[string]compiler.Module, error) {
	compilerMods := make(map[string]compiler.Module, len(mods))

	for name, bb := range mods {
		mod, err := p.unmarshaler.Unmarshal(bb)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrUnmarshaler, err)
		}
		compilerMods[name] = p.caster.Cast(mod)
	}

	return compilerMods, nil
}

func MustNewYaml() parser {
	return parser{
		yamlUnmarshaler{},
		caster{},
	}
}

package parser

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler"
)

type (
	Caster interface {
		Cast(Module) (compiler.Module, error)
	}
	Unmarshaler interface {
		Unmarshal([]byte) (Module, error)
	}
)

var (
	ErrUnmarshaler = errors.New("unmarshaler")
	ErrCaster      = errors.New("caster")
)

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

		cmod, err := p.caster.Cast(mod)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCaster, err)
		}

		compilerMods[name] = cmod
	}

	return compilerMods, nil
}

func MustNewYaml() parser {
	return parser{
		yamlUnmarshaler{},
		caster{},
	}
}

package parser

import (
	"fmt"

	"github.com/emil14/neva/internal/old/compiler/program"
)

type parser struct {
	marshal   marshal
	unmarshal unmarshal
	caster    caster
}

func (p parser) Module(bb []byte) (program.Module, error) {
	var mod module
	if err := p.unmarshal(bb, &mod); err != nil {
		return program.Module{}, err
	}

	return p.caster.to(mod)
}

func (p parser) Pkg(pkg []byte) (program.Module, error) {
	var mod pkg
	if err := p.unmarshal(bb, &mod); err != nil {
		return program.Module{}, err
	}

	return p.caster.to(mod)
}

type (
	marshal   func(interface{}) ([]byte, error)
	unmarshal func([]byte, interface{}) (err error)
	caster    interface {
		to(module) (program.Module, error)
		from(program.Module) module
	}
)

func New(u unmarshal, m marshal, c caster) (parser, error) {
	if u == nil || m == nil || c == nil {
		return parser{}, fmt.Errorf("parser constructor err")
	}

	return parser{
		unmarshal: u,
		marshal:   m,
		caster:    c,
	}, nil
}

func MustNew(u unmarshal, m marshal, c caster) parser {
	p, err := New(u, m, c)
	if err != nil {
		panic(err)
	}
	return p
}

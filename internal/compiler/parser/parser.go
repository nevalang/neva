package parser

import (
	"fmt"

	"github.com/emil14/respect/internal/compiler"
	"github.com/emil14/respect/internal/compiler/program"
)

type Caster interface {
	From(program.Module) module
	To(module) program.Module
}

type parser struct {
	marshal   Marshal
	unmarshal Unmarshal
	caster    Caster
}

func (p parser) Module(bb []byte) (program.Module, error) {
	var mod module
	if err := p.unmarshal(bb, &mod); err != nil {
		return program.Module{}, err
	}
	to := p.caster.To(mod)

	return to, nil
}

func (p parser) Program(bb program.Program) ([]byte, error) {
	return nil, nil // TODO
}

func (p parser) Unparse(mod program.Module) ([]byte, error) {
	bb, err := p.marshal(mod)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

type Unmarshal func([]byte, interface{}) (err error)

type Marshal func(interface{}) ([]byte, error)

func New(u Unmarshal, m Marshal, c Caster) (compiler.SRCParser, error) {
	if u == nil || m == nil || c == nil {
		return parser{}, fmt.Errorf("parser constructor err")
	}

	return parser{
		unmarshal: u,
		marshal:   m,
		caster:    c,
	}, nil
}

func MustNew(u Unmarshal, m Marshal, c Caster) compiler.SRCParser {
	p, err := New(u, m, c)
	if err != nil {
		panic(err)
	}
	return p
}

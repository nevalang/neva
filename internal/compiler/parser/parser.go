package parser

import (
	"fmt"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/program"
)

type parser struct {
	marshal   Marshal
	unmarshal Unmarshal
	caster    Caster
}

type Caster interface {
	From(program.Module) module
	To(module) program.Module
}

func (p parser) ParseModule(bb []byte) (program.Module, error) {
	var mod module
	if err := p.unmarshal(bb, &mod); err != nil {
		return program.Module{}, err
	}
	return p.caster.To(mod), nil
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

func New(u Unmarshal, m Marshal, c Caster) (compiler.Parser, error) {
	if u == nil || m == nil || c == nil {
		return parser{}, fmt.Errorf("parser constructor err")
	}

	return parser{
		unmarshal: u,
		marshal:   m,
		caster:    c,
	}, nil
}

func MustNew(u Unmarshal, m Marshal, c Caster) compiler.Parser {
	p, err := New(u, m, c)
	if err != nil {
		panic(err)
	}
	return p
}

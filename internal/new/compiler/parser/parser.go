package parser

import (
	"fmt"

	"github.com/emil14/neva/internal/new/compiler"
	"github.com/emil14/neva/internal/old/compiler/program"
)

type Caster interface {
	From(program.Module) Module
	To(Module) program.Module
}

type parser struct {
	marshal   Marshal
	unmarshal Unmarshal
	caster    Caster
}

func (p parser) Parse(bb []byte) (program.Module, error) {
	var mod Module
	if err := p.unmarshal(bb, &mod); err != nil {
		return program.Module{}, err
	}
	return p.caster.To(mod), nil
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

func New(u Unmarshal, m Marshal, c Caster) (compiler.ModuleParser, error) {
	if u == nil || m == nil || c == nil {
		return parser{}, fmt.Errorf("parser constructor err")
	}

	return parser{
		unmarshal: u,
		marshal:   m,
		caster:    c,
	}, nil
}

func MustNew(u Unmarshal, m Marshal, c Caster) compiler.ModuleParser {
	p, err := New(u, m, c)
	if err != nil {
		panic(err)
	}
	return p
}

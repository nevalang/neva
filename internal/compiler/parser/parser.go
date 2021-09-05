package parser

import (
	"fmt"

	"github.com/emil14/neva/internal/compiler/program"
)

type parser struct {
	marshal   Marshal
	unmarshal Unmarshal
	cast      Cast
}

func (p parser) Parse(bb []byte) (program.Module, error) {
	var mod module
	if err := p.unmarshal(bb, &mod); err != nil {
		return program.Module{}, err
	}
	return p.cast(mod), nil
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

type Cast func(module) program.Module

func New(u Unmarshal, m Marshal, c Cast) (parser, error) {
	if u == nil || m == nil || c == nil {
		return parser{}, fmt.Errorf("parser constructor err")
	}

	return parser{
		unmarshal: u,
		marshal:   m,
		cast:      c,
	}, nil
}

func MustNew(u Unmarshal, m Marshal, c Cast) parser {
	p, err := New(u, m, c)
	if err != nil {
		panic(err)
	}
	return p
}

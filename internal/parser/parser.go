package parser

import (
	"fmt"

	"github.com/emil14/refactored-garbanzo/internal/core"
	"gopkg.in/yaml.v2"
)

func MustNewYAML() Parser {
	return MustNew(yaml.Unmarshal, yaml.Marshal, cast)
}

type Parser struct {
	marshal   Marshal
	unmarshal Unmarshal
	cast      Cast
}

func (p Parser) Parse(bb []byte) (core.Module, error) {
	var mod module
	if err := p.unmarshal(bb, &mod); err != nil {
		return nil, err
	}

	return p.cast(mod)
}

// func (p Parser) Encode(mod core.Module) ([]byte, error) {
//  bb := []byte
//  dto := p.caster.From(mod)
// 	if err := p.marshal(dto, bb); err != nil {
// 		return nil, err
// 	}
// 	return bb, nil
// }

type Unmarshal func([]byte, interface{}) (err error)

type Marshal func(interface{}) ([]byte, error)

type Cast func(module) (core.Module, error)

func New(u Unmarshal, c Cast) (Parser, error) {
	if u == nil || c == nil {
		return Parser{}, fmt.Errorf("unmarshal or cast is nil")
	}

	return Parser{
		unmarshal: u,
		cast:      c,
	}, nil
}

func MustNew(u Unmarshal, m Marshal, c Cast) Parser {
	p, err := New(u, c)
	if err != nil {
		panic(err)
	}

	return p
}
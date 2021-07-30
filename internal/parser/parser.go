package parser

import (
	"fmt"

	"github.com/emil14/stream/internal/core"
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

type Unmarshal func([]byte, interface{}) (err error)

type Marshal func(interface{}) ([]byte, error)

type Cast func(module) (core.Module, error)

func New(u Unmarshal, m Marshal, c Cast) (Parser, error) {
	if u == nil || m == nil || c == nil {
		return Parser{}, fmt.Errorf("parser constructor err")
	}

	return Parser{
		unmarshal: u,
		marshal:   m,
		cast:      c,
	}, nil
}

func MustNew(u Unmarshal, m Marshal, c Cast) Parser {
	p, err := New(u, m, c)
	if err != nil {
		panic(err)
	}

	return p
}

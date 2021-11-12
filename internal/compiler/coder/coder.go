package coder

import (
	"encoding/json"

	"github.com/emil14/respect/internal/runtime/program"
)

type jsonCoder struct {
	caster interface {
		Cast(program.Program) Program
	}
	marshal func(interface{}) ([]byte, error)
}

func (c jsonCoder) Code(prog program.Program) ([]byte, error) {
	bb, err := c.marshal(c.caster.Cast(prog))
	if err != nil {
		return nil, err
	}

	return bb, nil
}

func New() jsonCoder {
	return jsonCoder{
		marshal: json.Marshal,
		// caster:  caster{},
	}
}

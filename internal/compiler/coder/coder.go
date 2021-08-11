package coder

import (
	"encoding/json"

	"github.com/emil14/stream/internal/runtime/program"
)

type jsonCoder struct{}

func (c jsonCoder) Code(prog program.Program) ([]byte, error) {
	bb, err := json.Marshal(prog)
	if err != nil {
		return nil, err
	}

	return bb, nil
}

func MustNewJSON() jsonCoder {
	return jsonCoder{}
}

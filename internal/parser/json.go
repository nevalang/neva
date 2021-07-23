package parser

import (
	"encoding/json"
)

type jsonParser struct {
}

func (p jsonParser) Parse(bb []byte) (Module, error) {
	var mod Module
	if err := json.Unmarshal(bb, &mod); err != nil {
		return Module{}, err
	}
	return mod, nil
}

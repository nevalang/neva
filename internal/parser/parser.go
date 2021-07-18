package parser

import "encoding/json"

type Parser interface {
	Parse([]byte) (Module, error)
}

type jsonParser struct {
	validator Validator
}

func (jp jsonParser) Parse(bb []byte) (Module, error) {
	var mod Module
	if err := json.Unmarshal(bb, &mod); err != nil {
		return Module{}, err
	}
	if err := jp.validator.Validate(mod); err != nil {
		return Module{}, err
	}
	return mod, nil
}

func NewParser(v Validator) Parser {
	return jsonParser{v}
}

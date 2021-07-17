package parser

import "encoding/json"

type Parser interface {
	Parse([]byte) (CustomModule, error)
}

type jsonParser struct {
	validator Validator
}

func (jp jsonParser) Parse(bb []byte) (CustomModule, error) {
	var mod CustomModule
	if err := json.Unmarshal(bb, &mod); err != nil {
		return CustomModule{}, err
	}
	if err := jp.validator.Validate(mod); err != nil {
		return CustomModule{}, err
	}
	return mod, nil
}

func NewParser(v Validator) Parser {
	return jsonParser{v}
}

package parser

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

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

func NewJSONParser(v Validator) Parser {
	return jsonParser{v}
}

type yamlParser struct {
	validator Validator
}

func (yp yamlParser) Parse(bb []byte) (Module, error) {
	var mod Module
	if err := yaml.Unmarshal(bb, &mod); err != nil {
		return Module{}, err
	}
	if err := yp.validator.Validate(mod); err != nil {
		return Module{}, err
	}
	return mod, nil
}

func NewYAMLParser(v Validator) Parser {
	return yamlParser{}
}

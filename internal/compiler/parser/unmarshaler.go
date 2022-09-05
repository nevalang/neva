package parser

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

var ErrYaml = errors.New("yaml")

type yamlUnmarshaler struct{}

func (u yamlUnmarshaler) Unmarshal(bb []byte) (Module, error) {
	var mod Module

	if err := yaml.Unmarshal(bb, &mod); err != nil {
		return Module{}, fmt.Errorf("%w: %v", ErrYaml, err)
	}

	return mod, nil
}

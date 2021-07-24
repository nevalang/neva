package parser

import yaml "gopkg.in/yaml.v2"

type yamlParser struct{}

func (p yamlParser) Parse(bb []byte) (Module, error) {
	var mod Module
	if err := yaml.Unmarshal(bb, &mod); err != nil {
		return Module{}, err
	}
	return mod, nil
}

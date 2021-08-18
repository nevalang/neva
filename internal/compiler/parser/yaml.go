package parser

import "gopkg.in/yaml.v2"

func MustNewYAML() parser {
	return MustNew(yaml.Unmarshal, yaml.Marshal, castModule)
}

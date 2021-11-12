package parser

import (
	"encoding/json"

	"github.com/emil14/respect/internal/compiler"
	"gopkg.in/yaml.v2"
)

func MustNewYAML() compiler.Parser {
	return MustNew(yaml.Unmarshal, yaml.Marshal, caster{})
}

func MustNewJSON() compiler.Parser {
	return MustNew(json.Unmarshal, json.Marshal, caster{})
}

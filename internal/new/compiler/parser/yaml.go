package parser

import (
	"encoding/json"

	"github.com/emil14/neva/internal/new/compiler"
	"gopkg.in/yaml.v2"
)

func MustNewYAML() compiler.ModuleParser {
	return MustNew(yaml.Unmarshal, yaml.Marshal, caster{})
}

func MustNewJSON() compiler.ModuleParser {
	return MustNew(json.Unmarshal, json.Marshal, caster{})
}

package parser

import (
	"encoding/json"

	"github.com/emil14/neva/internal/compiler"
	"gopkg.in/yaml.v2"
)

func MustNewYAML() compiler.SRCParser {
	return MustNew(yaml.Unmarshal, yaml.Marshal, caster{})
}

func MustNewJSON() compiler.SRCParser {
	return MustNew(json.Unmarshal, json.Marshal, caster{})
}

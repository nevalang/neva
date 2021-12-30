package parser

type module struct {
	IO      io                  `yaml:"io"`
	Deps    deps                `yaml:"deps"`
	Const   map[string]constant `yaml:"const"`
	Workers map[string]string   `yaml:"workers"`
	Start   bool                `yaml:"start"`
	Net     map[string][]string `yaml:"net,required"`
}

type deps map[string]io

type io struct {
	Params []string          `yaml:"params"`
	In     map[string]string `yaml:"in,required"`
	Out    map[string]string `yaml:"out,required"`
}

type constant struct {
	Type      string `yaml:"type"`
	IntValue  int    `yaml:"intValue"`
	StrValue  string `yaml:"strValue"`
	BoolValue bool   `yaml:"boolValue"`
}

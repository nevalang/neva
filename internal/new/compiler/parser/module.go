package parser

type Module struct {
	In      map[string]string   `yaml:"in"`
	Out     map[string]string   `yaml:"out"`
	Deps    map[string]IO       `yaml:"deps"`
	Const   map[string]Const    `yaml:"const"`
	Workers map[string]string   `yaml:"workers"`
	Net     map[string][]string `yaml:"net"`
}

type IO struct {
	In  map[string]string `yaml:"in"`
	Out map[string]string `yaml:"out"`
}

type Const struct {
	Type     string `yaml:"type"`
	IntValue int    `yaml:"intValue"`
}

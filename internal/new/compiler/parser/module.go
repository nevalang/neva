package parser

type Module struct {
	Meta  Meta
	IO    IO                  `yaml:"io"`
	Deps  map[string]IO       `yaml:"deps"`
	Nodes Nodes               `yaml:"nodes"`
	Net   map[string][]string `yaml:"net"`
}

type IO struct {
	In  map[string]string `yaml:"in"`
	Out map[string]string `yaml:"out"`
}

type Nodes struct {
	Const   map[string]ConstOutPort `yaml:"const"`
	Workers map[string]string       `yaml:"workers"`
}

type ConstOutPort struct {
	Type string `yaml:"type"`
	Int  int    `yaml:"int"`
	Str  string `yaml:"int"`
	Bool bool   `yaml:"int"`
}

type Meta struct {
	Compiler string `yaml:"compiler"`
}

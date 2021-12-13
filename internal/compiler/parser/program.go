package parser

type module struct {
	In      inports          `yaml:"in"`
	Out     outports         `yaml:"out"`
	Deps    moduleDeps       `yaml:"deps"`
	Const   map[string]Const `yaml:"const"`
	Workers workers          `yaml:"workers"`
	Net     net              `yaml:"net,required"`
}

type inports ports

type outports ports

type ports map[string]string

type moduleDeps map[string]io

type io struct {
	In  inports  `yaml:"in"`
	Out outports `yaml:"out"`
}

type Const struct {
	Type      string `yaml:"type"`
	IntValue  int    `yaml:"intValue"`
	StrValue  string `yaml:"strValue"`
	BoolValue bool   `yaml:"boolValue"`
}

type workers map[string]string

// senders -> outports -> receivers -> inports
type net map[string]map[string]map[string][]string

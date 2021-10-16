package parser

type Program struct {
	Deps   deps              `yaml:"deps"`
	Import map[string]string `yaml:"import"`
	Root   string            `yaml:"root"`
}

type deps map[string]struct {
	Repo    string `yaml:"repo"`
	Version string `yaml:"v"`
}

type module struct {
	In      inports          `yaml:"in"`
	Out     outports         `yaml:"out"`
	Deps    moduleDeps       `yaml:"deps"`
	Const   map[string]Const `yaml:"const"`
	Workers workers          `yaml:"workers"`
	Net     net              `yaml:"net"`
}

type inports ports

type outports ports

type ports map[string]string // name -> type

// module -> interface
type moduleDeps map[string]io

type io struct {
	In  inports  `yaml:"in"`
	Out outports `yaml:"out"`
}

type Const struct {
	Type  string      `yaml:"type"`
	Value interface{} `yaml:"value"`
}

type workers map[string]string // worker -> dep

// senders -> outports -> receivers -> inports
type net map[string]map[string]map[string][]string

package parser

type Program struct {
	Deps   progDeps          `yaml:"deps" json:"deps"`
	Import map[string]string `yaml:"import" json:"import"`
	Root   string            `yaml:"root" json:"root"`
}

type progDeps map[string]struct {
	Repo    string `yaml:"repo" json:"repo"`
	Version string `yaml:"v" json:"v"`
}

type Module struct {
	Out     outports   `yaml:"out" json:"out"`         // input ports
	In      inports    `yaml:"in" json:"in"`           // output ports
	Deps    moduleDeps `yaml:"deps" json:"deps"`       // deps interfaces
	Workers workers    `yaml:"workers" json:"workers"` // maps workers to components
	Net     net        `yaml:"net" json:"net"`         // data flow
}

type inports ports

type outports ports

type ports map[string]string // name -> type

// module -> interface
type moduleDeps map[string]IO

type IO struct {
	In  inports  `yaml:"in" json:"in"`
	Out outports `yaml:"out" json:"out"`
}

type workers map[string]string // worker -> dep

// senders -> outports -> receivers -> inports
type net map[string]map[string]map[string][]string

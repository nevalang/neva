package parser

type program struct {
	Deps   progDeps          `yaml:"deps"`
	Import map[string]string `yaml:"import"`
	Root   string            `yaml:"root"`
}

type progDeps map[string]struct {
	Repo    string `yaml:"repo"`
	Version string `yaml:"v"`
}

type module struct {
	Out     outports   `yaml:"out"`     // input ports
	In      inports    `yaml:"in"`      // output ports
	Deps    moduleDeps `yaml:"deps"`    // deps interfaces
	Workers workers    `yaml:"workers"` // maps workers to components
	Net     net        `yaml:"net"`     // data flow
}

type inports ports

type outports ports

type ports map[string]string // name -> type

// module -> interface
type moduleDeps map[string]struct {
	In  inports  `yaml:"in"`
	Out outports `yaml:"out"`
}

type workers map[string]string // worker -> dep

// senders -> outports -> receivers -> inports
type net map[string]map[string]map[string][]string

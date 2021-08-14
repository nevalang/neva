package parser

type module struct {
	Out     outports `yaml:"out"`     // input ports
	In      inports  `yaml:"in"`      // output ports
	Deps    deps     `yaml:"deps"`    // deps interfaces
	Workers workers  `yaml:"workers"` // maps workers to components
	Net     net      `yaml:"net"`     // data flow
}

type inports ports

type outports ports

type ports map[string]string // name -> type

// module -> interface
type deps map[string]struct {
	In  inports  `yaml:"in"`
	Out outports `yaml:"out"`
}

type workers map[string]string // worker -> dep

// senders -> outports -> receivers -> inports
type net map[string]map[string]map[string][]string

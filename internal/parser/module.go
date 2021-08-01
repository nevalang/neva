package parser

// custom module.
type module struct {
	In      inports  `yaml:"in"`
	Out     outports `yaml:"out"`
	Deps    deps     `yaml:"deps"`
	Workers workers  `yaml:"workers"`
	Net     net      `yaml:"net"`
}

// input ports.
type inports ports

// output ports.
type outports ports

// port name -> type.
type ports map[string]string

// module name -> interface.
type deps map[string]struct {
	In  inports  `yaml:"in"`
	Out outports `yaml:"out"`
}

// worker -> dep.
type workers map[string]string

// senders -> outports -> receivers -> inports.
type net map[string]map[string]map[string][]string

package parser

// user defined component.
type module struct {
	In      inports  `yaml:"in"`      // output ports
	Out     outports `yaml:"out"`     // input ports
	Deps    deps     `yaml:"deps"`    // describes dependencies interfaces
	Workers workers  `yaml:"workers"` // maps workers to components
	Net     net      `yaml:"net"`     // describes data flow
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

// worker -> dependency.
type workers map[string]string

// senders -> outports -> receivers -> inports.
type net map[string]map[string]map[string][]string

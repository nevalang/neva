package parser

type module struct {
	In      inports  `yaml:"in"`
	Out     outports `yaml:"out"`
	Deps    deps     `yaml:"deps"`
	Workers workers  `yaml:"workers"`
	Net     net      `yaml:"net"`
}

type inports Ports

type outports Ports

type Ports map[string]string

type deps map[string]struct {
	In  inports  `yaml:"in"`
	Out outports `yaml:"out"`
}

type workers map[string]string

type net map[string]conns

type conns map[string]conn

type conn map[string][]string

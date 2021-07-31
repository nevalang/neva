package parser

type module struct {
	In      inports  `json:"in"`
	Out     outports `json:"out"`
	Deps    deps     `json:"deps"`
	Workers workers  `json:"workers"`
	Net     net      `json:"net"`
}

type inports Ports

type outports Ports

type Ports map[string]string

type deps map[string]struct {
	In  inports  `json:"in"`
	Out outports `json:"out"`
}

type workers map[string]string

type net map[string]conns

type conns map[string]conn

type conn map[string][]string

package parser

type module struct {
	in      inports
	out     outports
	deps    deps
	workers workers
	net     net
}

type inports Ports

type outports Ports

type Ports map[string]string

type deps map[string]struct {
	In  inports
	Out outports
}

type workers map[string]string

type net map[string]conns

type conns map[string]conn

type conn map[string][]string


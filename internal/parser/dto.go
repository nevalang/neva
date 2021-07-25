package parser

type Module struct {
	In      Inports
	Out     Outports
	Deps    Deps
	Workers Workers
	Net     Net
}

type Inports Ports

type Outports Ports

type Ports map[string]string

type Deps map[string]struct {
	In  Inports
	Out Outports
}

type Workers map[string]string

type Net map[string]Conns

type Conns map[string]Conn

type Conn map[string][]string

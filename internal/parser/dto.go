package parser

type Module struct {
	In      InportsInterface
	Out     OutportsInterface
	Deps    Deps
	Workers Workers
	Net     Net
}

type InportsInterface PortsInterface

type OutportsInterface PortsInterface

type PortsInterface map[string]string

type Deps map[string]struct {
	In  InportsInterface
	Out OutportsInterface
}

type Workers map[string]string

type Net map[string]Connections

type Connections map[string]Connection

type Connection map[string][]string

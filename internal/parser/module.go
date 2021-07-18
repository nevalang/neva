package parser

type Module struct {
	Deps    Deps              `json:"deps"`
	In      InportsInterface  `json:"in"`
	Out     OutportsInterface `json:"out"`
	Workers Workers           `json:"workers"`
	Net     Net               `json:"net"`
}

type Deps map[string]struct {
	In  InportsInterface
	Out OutportsInterface
}

type InportsInterface PortsInterface

type OutportsInterface PortsInterface

type PortsInterface map[string]string

type Workers map[string]string

type Net []Subscription

type Subscription struct {
	Sender    PortPoint   `json:"sender"`
	Recievers []PortPoint `json:"recievers"`
}

type PortPoint struct {
	Node string `json:"node"`
	Port string `json:"port"`
}

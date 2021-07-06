package parsing

import "encoding/json"

// Module represents module after parsing.
type Module struct {
	Deps      Deps              `json:"deps"`
	In        InPorts           `json:"in"`      // inports map
	Out       OutPorts          `json:"out"`     // outports map
	WorkerMap map[string]string `json:"workers"` // maps workername to depname
	Net       Net               `json:"net"`
}

type Deps map[string]struct {
	In  InPorts
	Out OutPorts
}

type InPorts Ports

type OutPorts Ports

type Ports map[string]string

type Net []struct {
	Sender    Conn   `json:"sender"`
	Recievers []Conn `json:"recievers"`
}

type Conn struct {
	Node string `json:"node"`
	Port string `json:"port"`
}

func FromJSON(bb []byte) (Module, error) {
	var m Module
	if err := json.Unmarshal(bb, &m); err != nil {
		return Module{}, err
	}
	return m, nil
}

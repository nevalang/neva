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

type InPorts PortMap

type OutPorts PortMap

type PortMap map[string]string

type Net []struct {
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

package parsing

// Parsed represents module after parsing.
type Parsed struct {
	Deps      Deps              `json:"deps"`
	In        PortMap           `json:"in"`      // inports map
	Out       PortMap           `json:"out"`     // outports map
	WorkerMap map[string]string `json:"workers"` // maps workername to depname
	Net       Net               `json:"net"`
}

type Deps map[string]struct {
	In, Out PortMap
}

type PortMap map[string]string

type Net []struct {
	Node string `json:"node"`
	Port string `json:"port"`
}

package coder

type Program struct {
	Components map[string]Component `json:"components"`
	Root       NodeMeta             `json:"root"`
}

type Component struct {
	Operator string              `json:"operator,omitempty"`
	Workers  map[string]NodeMeta `json:"workers,omitempty"`
	Net      []Connection        `json:"net,omitempty"`
}

type NodeMeta struct {
	In        map[string]uint8 `json:"in"`
	Out       map[string]uint8 `json:"out"`
	Component string           `json:"component"`
}

type Connection struct {
	From PortAddr   `json:"from"`
	To   []PortAddr `json:"to"`
}

type PortAddr struct {
	Node string `json:"node"`
	Port string `json:"port"`
	Idx  uint8  `json:"idx,omitempty"`
}

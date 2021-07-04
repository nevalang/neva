package runtime

import "fbp/internal/types"

type Module struct {
	Deps      Deps              `json:"deps"`
	In        InPorts           `json:"in"`
	Out       OutPorts          `json:"out"`
	WorkerMap map[string]string `json:"workers"`
	Net       []Connection      `json:"net"`
}

type Deps map[string]Module

type InPorts PortMap

type OutPorts PortMap

type PortMap map[string]types.Type

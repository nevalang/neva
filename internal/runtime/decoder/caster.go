package decoder

import (
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/pkg/runtimesdk"
)

type caster struct{}

func (c caster) Cast(runtimesdk.Program) (runtime.Program, error) {
	return runtime.Program{}, nil // TODO
}

func (c caster) components(from map[string]Component) map[string]program.Component {
	to := make(map[string]program.Component, len(from))
	for name, component := range from {
		to[name] = program.Component{
			OperatorName:        component.Operator,
			WorkerNodesMeta: c.nodesMeta(component.Workers),
			Net:             c.net(component.Net),
		}
	}
	return to
}

func (c caster) nodesMeta(workers map[string]NodeMeta) map[string]program.WorkerNodeMeta {
	result := make(map[string]program.WorkerNodeMeta, len(workers))

	for w, nodeMeta := range workers {
		result[w] = program.WorkerNodeMeta{
			// Node:      w,
			In:            nodeMeta.In,
			Out:           nodeMeta.Out,
			ComponentName: nodeMeta.Component,
		}
	}

	return result
}

func (c caster) net(model []Connection) []program.Connection {
	result := make([]program.Connection, len(model))
	for i := range model {
		result[i] = program.Connection{
			From: program.PortAddr(model[i].From),
			To:   c.portAddrs(model[i].To),
		}
	}
	return result
}

func (c caster) portAddrs(from []PortAddr) []program.PortAddr {
	to := make([]program.PortAddr, len(from))
	for i := range from {
		to[i] = program.PortAddr(from[i])
	}
	return to
}

func NewCaster() caster {
	return caster{}
}

package decoder

import "github.com/emil14/neva/internal/runtime/program"

type caster struct{}

func (c caster) Cast(prog Program) program.Program {
	return program.Program{
		RootNode: program.NodeMeta{
			Node:      "root",
			In:        prog.RootNode.In,
			Out:       prog.RootNode.Out,
			Component: prog.RootNode.Component,
		},
		Scope: c.components(prog.Scope),
	}
}

func (c caster) components(from map[string]Component) map[string]program.Component {
	to := make(map[string]program.Component, len(from))
	for name, component := range from {
		to[name] = program.Component{
			Operator:        component.Operator,
			WorkerNodesMeta: c.workerNodesMeta(component.Workers),
			Connections:     c.net(component.Net),
		}
	}
	return to
}

func (c caster) workerNodesMeta(workers map[string]NodeMeta) map[string]program.NodeMeta {
	result := make(map[string]program.NodeMeta, len(workers))

	for w, nodeMeta := range workers {
		result[w] = program.NodeMeta{
			Node:      w,
			In:        nodeMeta.In,
			Out:       nodeMeta.Out,
			Component: nodeMeta.Component,
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

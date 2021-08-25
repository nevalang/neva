package decoder

import "github.com/emil14/neva/internal/runtime/program"

type caster struct{}

func (c caster) Cast(prog Program) program.Program {
	return program.Program{
		Root:       program.NodeMeta(prog.Root),
		Components: c.components(prog.Components),
	}
}

func (c caster) components(from map[string]Component) map[string]program.Component {
	to := make(map[string]program.Component, len(from))
	for name, component := range from {
		to[name] = program.Component{
			Operator: component.Operator,
			Workers:  c.workers(component.Workers),
			Net:      c.net(component.Net),
		}
	}
	return to
}

func (c caster) workers(from map[string]NodeMeta) map[string]program.NodeMeta {
	to := make(map[string]program.NodeMeta, len(from))
	for k, v := range from {
		to[k] = program.NodeMeta(v)
	}
	return to
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

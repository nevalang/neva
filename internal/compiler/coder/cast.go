package coder

import "github.com/emil14/respect/internal/runtime/program"

type caster struct{}

func (c caster) Cast(prog program.Program) Program {
	return Program{
		Root:       c.castNode(prog.RootNodeMeta),
		Components: c.castComponents(prog.Scope),
	}
}

func (c caster) castNode(node program.WorkerNodeMeta) NodeMeta {
	return NodeMeta{
		In:        node.In,
		Out:       node.Out,
		Component: node.ComponentName,
	}
}

func (c caster) castComponents(cc map[string]program.Component) map[string]Component {
	res := make(map[string]Component, len(cc))

	for name := range cc {
		res[name] = Component{
			Operator: cc[name].OperatorName,
			Workers:  c.castNodes(cc[name].WorkerNodesMeta),
			Net:      c.castNet(cc[name].Net),
		}
	}

	return res
}

func (c caster) castNodes(nodes map[string]program.WorkerNodeMeta) map[string]NodeMeta {
	res := make(map[string]NodeMeta, len(nodes))

	for name := range nodes {
		res[name] = c.castNode(nodes[name])
	}

	return res
}

func (c caster) castNet(cc []program.Connection) []Connection {
	res := make([]Connection, len(cc))

	for i := range cc {
		res[i] = Connection{
			From: PortAddr(cc[i].From),
			To:   c.castPortAddrs(cc[i].To),
		}
	}

	return res
}

func (c caster) castPortAddrs(aa []program.PortAddr) []PortAddr {
	res := make([]PortAddr, len(aa))

	for i := range aa {
		res[i] = PortAddr(aa[i])
	}

	return res
}

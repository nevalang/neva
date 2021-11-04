package runtime

import (
	"fmt"

	"github.com/emil14/respect/internal/core"
)

type (
	Runtime struct {
		connector Connector
		operators map[string]OperatorFunc
	}

	Connector interface {
		Connect([]Connection)
	}

	OperatorFunc func(core.IO) error
)

func (r Runtime) Run(p Program) (core.IO, error) {
	return r.spawnWorkerNode("root", p.RootNodeMeta, p.Scope)
}

func (r Runtime) spawnWorkerNode(
	nodeName string,
	nodeMeta WorkerNodeMeta,
	scope map[string]Component,
) (core.IO, error) {
	component, ok := scope[nodeMeta.ComponentName]
	if !ok {
		return core.IO{}, fmt.Errorf("component not found: %s", nodeMeta.ComponentName)
	}

	var io = r.nodeIO(nodeMeta)

	if component.Type == OperatorComponent {
		op, ok := r.operators[component.Operator.Name]
		if !ok {
			return core.IO{}, fmt.Errorf("operator not found: %s", component.Operator.Name)
		}

		if err := op(io); err != nil {
			return core.IO{}, fmt.Errorf("connect operator: %w", err)
		}

		return r.patchIO(nodeMeta, io, nodeName), nil
	}

	if component.Type != ModuleComponent {
		return core.IO{}, fmt.Errorf("%s has unknown type %d", nodeMeta.ComponentName, component.Type)
	}

	netIO := map[string]core.IO{
		"in":  {Out: io.In},
		"out": {In: io.Out},
	}

	if l := len(component.Module.Const); l > 0 {
		constOutPorts := make(core.Ports, l)

		for name, cnst := range component.Module.Const {
			constOutPorts[name] = r.constOutPort(cnst)
		}

		netIO["const"] = core.IO{Out: constOutPorts}
	}

	for workerNodeName, workerNodeMeta := range component.Module.Workers {
		nodeIO, err := r.spawnWorkerNode(workerNodeName, workerNodeMeta, scope)
		if err != nil {
			return core.IO{}, err
		}
		netIO[workerNodeName] = nodeIO
	}

	cc, err := r.connections(netIO, component.Module.Net)
	if err != nil {
		return core.IO{}, err
	}

	r.connector.Connect(cc)

	return r.patchIO(nodeMeta, io, nodeName), nil
}

func (r Runtime) constOutPort(c ConstValue) chan core.Msg {
	var msg core.Msg

	switch c.Type {
	case IntValue:
		msg = core.NewIntMsg(c.IntValue)
	case BoolValue:
		msg = core.NewBoolMsg(c.BoolValue)
	case StrValue:
		msg = core.NewStrMsg(c.StrValue)
	}

	out := make(chan core.Msg)
	go func() {
		for {
			out <- msg
		}
	}()

	return out
}

// patchIO replaces "in" and "out" node names with worker name from parent network
func (r Runtime) patchIO(meta WorkerNodeMeta, io core.IO, nodeName string) core.IO {
	patch := core.IO{
		In:  map[PortAddr]chan core.Msg{},
		Out: map[PortAddr]chan core.Msg{},
	}

	for addr, ch := range io.In {
		addr.Node = nodeName
		patch.In[addr] = ch
	}
	for addr, ch := range io.Out {
		addr.Node = nodeName
		patch.Out[addr] = ch
	}

	return patch
}

// connections initializes channels for network.
func (r Runtime) connections(nodesIO map[string]core.IO, net []ConnectionDef) ([]Connection, error) {
	cc := make([]Connection, len(net))

	for i, c := range net {
		fromNodeIO, ok := nodesIO[c.From.Node]
		if !ok {
			return nil, fmt.Errorf("not found IO for node %s", c.From.Node)
		}

		sender, ok := fromNodeIO.Out[PortAddr(c.From)]
		if !ok {
			return nil, fmt.Errorf("outport %s not found", c.From)
		}

		receivers := make([]Port, len(c.To))
		for j, toAddr := range c.To {
			toNodeIO, ok := nodesIO[toAddr.Node]
			if !ok {
				return nil, fmt.Errorf("io for receiver node not found: %s", toAddr.Node)
			}

			receiver, ok := toNodeIO.In[toAddr.Port]
			if !ok {
				return nil, fmt.Errorf("inport not found %s", toAddr)
			}

			receivers[j] = Port{Ch: receiver, Addr: toAddr}
		}

		cc[i] = Connection{
			From: Port{Ch: sender, Addr: c.From},
			To:   receivers,
		}
	}

	return cc, nil
}

func (r Runtime) nodeIO(nodeMeta WorkerNodeMeta) core.IO {
	inPorts := make(map[PortAddr]chan core.Msg)

	for port, slots := range nodeMeta.In {
		addr := PortAddr{Port: port, Node: "in"}

		if slots == 0 {
			inPorts[addr] = make(chan core.Msg)
			continue
		}

		for idx := uint8(0); idx < slots; idx++ {
			addr.Slot = idx
			inPorts[addr] = make(chan core.Msg)
		}
	}

	outPorts := make(map[PortAddr]chan core.Msg)

	for port, slots := range nodeMeta.Out {
		addr := PortAddr{Port: port, Node: "out"}

		if slots == 0 {
			outPorts[addr] = make(chan core.Msg)
			continue
		}

		for idx := uint8(0); idx < slots; idx++ {
			addr.Slot = idx
			outPorts[addr] = make(chan core.Msg)
		}
	}

	return core.IO{
		In:  inPorts,
		Out: outPorts,
	}
}

type Connection struct {
	From Port
	To   []Port
}

type ConnectionDef struct {
	From PortAddr
	To   []PortAddr
}

type PortAddr struct {
	Node, Port string
	Slot       uint8
}

type Port struct {
	Ch   chan core.Msg
	Addr PortAddr
}

func New(connector Connector) Runtime {
	return Runtime{
		connector: connector,
	}
}

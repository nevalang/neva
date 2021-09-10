package runtime

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime/program"
)

var (
	ErrPortNotFound = errors.New("port not found")
)

type Runtime struct {
	connector Connector
}

func (r Runtime) Run(p program.Program) (IO, error) {
	return r.run(p.Scope, p.RootNode)
}

func (r Runtime) run(scope map[string]program.Component, node program.NodeMeta) (IO, error) {
	component, ok := scope[node.Component]
	if !ok {
		return IO{}, fmt.Errorf("component not found: %s", node.Component)
	}

	io := r.nodeIO(node)

	if component.Operator != "" {
		if err := r.connector.ConnectOperator(component.Operator, io); err != nil {
			return IO{}, fmt.Errorf("could not connect operator")
		}
		return io, nil
	}

	// for this subnet 'in' is sender and 'out' is receiver nodes
	nodesIO := map[string]IO{
		"in":  {Out: io.In},
		"out": {In: io.Out},
	}

	for workerNode, meta := range component.WorkerNodesMeta {
		io, err := r.run(scope, meta)
		if err != nil {
			return IO{}, err
		}
		nodesIO[workerNode] = io
	}

	r.connector.ConnectSubnet(
		r.connections(nodesIO, component.Connections),
	)

	return io, nil
}

func (r Runtime) connections(nodesIO map[string]IO, net []program.Connection) []Connection {
	ss := make([]Connection, len(net))

	for i := range net {
		fromNodeIO, ok := nodesIO[net[i].From.Node]
		if !ok {
			panic("not ok")
		}

		fromOutportAddr := PortAddr{port: net[i].From.Port, idx: net[i].From.Idx, node: net[i].From.Node}
		from, ok := fromNodeIO.Out[fromOutportAddr]
		if !ok {
			panic("not ok")
		}

		to := make([]Port, len(net[i].To))
		for j := range net[i].To {
			toInportAddr := PortAddr{
				node: net[i].To[j].Node,
				port: net[i].To[j].Port,
				idx:  net[i].To[j].Idx,
			}

			toNodeIO := nodesIO[net[i].To[j].Node]
			receiver, ok := toNodeIO.In[toInportAddr]
			if !ok {
				panic("not ok")
			}

			to[j] = Port{Ch: receiver, Addr: toInportAddr}
		}

		ss[i] = Connection{
			From: Port{Ch: from, Addr: fromOutportAddr},
			To:   to,
		}
	}

	return ss
}

func (r Runtime) nodeIO(nodeMeta program.NodeMeta) IO {
	inports := make(map[PortAddr]chan Msg)

	for port, slots := range nodeMeta.In {
		addr := PortAddr{port: port, node: nodeMeta.Node}
		if addr.node == "root" {
			addr.node = "in"
		}

		if slots == 0 {
			inports[addr] = make(chan Msg)
			continue
		}

		for i := uint8(0); i < slots; i++ {
			addr.idx = i
			inports[addr] = make(chan Msg, slots)
		}
	}

	outports := make(map[PortAddr]chan Msg)

	for port, slots := range nodeMeta.Out {
		addr := PortAddr{port: port, node: nodeMeta.Node}
		if addr.node == "root" {
			addr.node = "out"
		}

		if slots == 0 {
			outports[addr] = make(chan Msg)
			continue
		}

		for idx := uint8(0); idx < slots; idx++ {
			outports[addr] = make(chan Msg)
		}
	}

	return IO{inports, outports}
}

// Operator spawns a goroutine where the real computation happens.
// That goroutine will use given IO (usually provided by runtime) to receive and send data.
// If IO doesn't satisfy the interface - error is returned.
type Operator func(IO) error

// Connection represents sender-receiver pair.
type Connection struct {
	From Port
	To   []Port
}

type Port struct {
	Ch   chan Msg
	Addr PortAddr
}

// IO represents node's input and output ports.
type IO struct {
	In, Out Ports
}

// Ports maps ports to their network addresses.
type Ports map[PortAddr]chan Msg

// Slots returns all port-chanells associated with the given array port name.
func (p Ports) Slots(arrPort string) ([]chan Msg, error) {
	cc := []chan Msg{}
	for addr, ch := range p {
		if addr.port == arrPort {
			cc = append(cc, ch)
		}
	}

	if len(cc) == 0 {
		return nil, fmt.Errorf("ErrArrPortNotFound: %s", arrPort)
	}

	return cc, nil
}

// Chan returns chanell associated with the given normal port name.
func (p Ports) Chan(port string) (chan Msg, error) {
	for addr, ch := range p {
		if addr.port == port {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("%w: looking for %s, have: %v", ErrPortNotFound, port, p)
}

// PortAddr describes port address in the
type PortAddr struct {
	node string
	port string
	idx  uint8 // always 0 for normal ports
}

func New(connector Connector) Runtime {
	return Runtime{
		connector: connector,
	}
}

// AbsPortAddr represents absolute port address in the program's network.
type AbsPortAddr struct {
	port Port
	path []string
}

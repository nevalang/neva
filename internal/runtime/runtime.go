package runtime

import (
	"fmt"

	"github.com/emil14/neva/internal/runtime/program"
)

// AbsPortAddr represents absolute port address in the program's network.
type AbsPortAddr struct {
	port Port
	path []string
}

type Runtime struct {
	Operators map[string]Operator
}

func (r Runtime) Run(p program.Program) (IO, error) {
	return r.run(p.Scope, p.Root)
}

func (r Runtime) run(components map[string]program.Component, node program.NodeMeta) (IO, error) {
	component, ok := components[node.Component]
	if !ok {
		return IO{}, fmt.Errorf("component not found: %s", node.Component)
	}

	if component.Operator != "" {
		io := r.nodeIO(node)
		if err := r.connectOperator(component.Operator, io); err != nil {
			return IO{}, fmt.Errorf("could not connect operator")
		}
		return io, nil
	}

	io := r.nodeIO(node)

	nodesIO := map[string]IO{
		"in":  {Out: io.In}, // for this net 'in' is sender
		"out": {In: io.Out}, // and 'out' is receiver
	}
	for workerNode, meta := range component.WorkerNodes {
		io, err := r.run(components, meta)
		if err != nil {
			return IO{}, err
		}
		nodesIO[workerNode] = io
	}

	r.connectMany(
		r.connections(nodesIO, component.Net),
	)

	return io, nil
}

func (r Runtime) connections(nodesIO map[string]IO, net []program.Connection) []pair {
	ss := make([]pair, len(net))

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

			to[j] = Port{ch: receiver, addr: toInportAddr}
		}

		ss[i] = pair{
			from: Port{ch: from, addr: fromOutportAddr},
			to:   to,
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

func (r Runtime) connectOperator(name string, io IO) error {
	op, ok := r.Operators[name]
	if !ok {
		return fmt.Errorf("ErrUnknownOperator: %s", name)
	}

	if err := op(io); err != nil {
		return err
	}

	return nil
}

func (r Runtime) connectMany(cc []pair) {
	for i := range cc {
		go r.connect(cc[i])
	}
}

func (r Runtime) connect(s pair) {
	for msg := range s.from.ch {
		for _, recv := range s.to {
			select {
			case recv.ch <- msg:
				continue
			default:
				go func(to chan Msg, m Msg) {
					to <- m
				}(recv.ch, msg)
			}
		}
	}
}

// Operator spawns a goroutine where the real computation happens.
// It uses an IO (usually provided by runtime) to receive and send data.
// If given io won't fit the interface then error should be returned.
type Operator func(IO) error

// pair represents pair betwen sender and receiver.
type pair struct {
	from Port
	to   []Port
}

type Port struct {
	ch   chan Msg
	addr PortAddr
}

// IO represents node's input and output ports.
type IO struct {
	In, Out Ports
}

// Ports maps port-channels with their network addresses.
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

// Chan returns all port-chanells associated with the given normal port name.
func (p Ports) Chan(port string) (chan Msg, error) {
	for addr, ch := range p {
		if addr.port == port {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("ErrPortNotFound: %s", port)
}

// PortAddr describes port address in the network.
type PortAddr struct {
	node string
	port string
	idx  uint8 // always 0 for normal ports
}

func New(ops map[string]Operator) Runtime {
	return Runtime{
		Operators: ops,
	}
}

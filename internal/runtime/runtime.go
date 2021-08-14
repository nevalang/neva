package runtime

import (
	"fmt"

	"github.com/emil14/neva/internal/runtime/program"
)

type Runtime struct {
	Operators map[string]Operator
}

func (r Runtime) Run(p program.Program) (IO, error) {
	return r.connectNode(p.Components, p.Root)
}

func (r Runtime) connectNode(scope map[string]program.Component, node program.NodeMeta) (IO, error) {
	component, ok := scope[node.Component]
	if !ok {
		return IO{}, fmt.Errorf(node.Component)
	}

	if component.Operator != "" {
		io, err := r.connectOperator(component.Operator, node)
		if err != nil {
			return IO{}, fmt.Errorf("...")
		}
		return io, nil
	}

	in, out := r.nodeIO(node)
	net := map[string]IO{
		"in":  {Out: in}, // for this module's net 'in' is sender
		"out": {In: out}, // for this module's net 'out' is receiver
	}

	for workerNode, meta := range component.Workers {
		io, err := r.connectNode(scope, meta)
		if err != nil {
			return IO{}, err
		}

		net[workerNode] = io
	}

	r.startStreams(
		r.connections(net, component.Net),
	)

	return IO{in, out}, nil
}

func (r Runtime) connections(nodesIO map[string]IO, net []program.Connection) []connection {
	ss := make([]connection, len(net))

	for i := range net {
		nodeIO, ok := nodesIO[net[i].From.Node]
		if !ok {
			panic("not ok")
		}

		from, ok := nodeIO.Out[PortAddr{
			port: net[i].From.Port,
			idx:  net[i].From.Idx,
		}]
		if !ok {
			panic("not ok")
		}

		to := make([]chan Msg, len(net[i].To))
		for j := range net[i].To {
			receiver, ok := nodeIO.In[PortAddr{
				port: net[i].To[j].Port,
				idx:  net[i].To[j].Idx,
			}]
			if !ok {
				panic("not ok")
			}

			to[j] = receiver
		}

		ss[i] = connection{
			from: from,
			to:   to,
		}
	}

	return nil
}

func (r Runtime) nodeIO(node program.NodeMeta) (in Ports, out Ports) {
	inports := make(map[PortAddr]chan Msg)
	for port, size := range node.In {
		if size > 0 {
			for i := uint8(0); i < size; i++ {
				inports[PortAddr{
					port: port,
					idx:  i,
				}] = make(chan Msg)
			}
			continue
		}

		inports[PortAddr{
			port: port,
		}] = make(chan Msg)
	}

	outports := make(map[PortAddr]chan Msg)
	for port, size := range node.In {
		if size > 0 {
			for i := uint8(0); i < size; i++ {
				outports[PortAddr{
					port: port,
					idx:  i,
				}] = make(chan Msg)
			}
			continue
		}

		outports[PortAddr{
			port: port,
		}] = make(chan Msg)
	}

	return inports, outports
}

func (r Runtime) connectOperator(name string, node program.NodeMeta) (IO, error) {
	connector, ok := r.Operators[name]
	if !ok {
		return IO{}, fmt.Errorf("ErrUnknownOperator: %s", name)
	}

	in, out := r.nodeIO(node)
	io := IO{in, out}
	if err := connector(io); err != nil {
		return IO{}, err
	}

	return io, nil
}

func (r Runtime) startStreams(ss []connection) {
	for i := range ss {
		go r.startStream(ss[i])
	}
}

func (r Runtime) startStream(s connection) {
	for msg := range s.from {
		for _, receiver := range s.to {
			select {
			case receiver <- msg:
				continue
			default:
				go func(to chan Msg, m Msg) {
					to <- m
				}(receiver, msg)
			}
		}
	}
}

type Operator func(IO) error

type connection struct {
	from chan Msg
	to   []chan Msg
}

type IO struct {
	In, Out Ports
}

type Ports map[PortAddr]chan Msg

// Slots returns all channels associated with the given array port name.
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

// Port returns all channels associated with the given normal port name.
func (p Ports) Port(port string) (chan Msg, error) {
	for addr, ch := range p {
		if addr.port == port {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("ErrPortNotFound: %s", port)
}

type PortAddr struct {
	port string
	idx  uint8 // always 0 for normal ports
}

func New(ops map[string]Operator) Runtime {
	return Runtime{
		Operators: ops,
	}
}

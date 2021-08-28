package runtime

import (
	"fmt"

	"github.com/emil14/neva/internal/runtime/program"
)

type Runtime struct {
	Operators map[string]Operator
}

func (r Runtime) Run(p program.Program) (IO, error) {
	return r.run(p.Components, p.Root)
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
		"out": {In: io.Out}, // for this net 'out' is receiver
	}
	for workerNode, meta := range component.Workers {
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

func (r Runtime) connections(nodesIO map[string]IO, net []program.Connection) []connection {
	ss := make([]connection, len(net))

	for i := range net {
		fromNodeIO, ok := nodesIO[net[i].From.Node]
		if !ok {
			panic("not ok")
		}

		fromOutportAddr := PortAddr{port: net[i].From.Port, idx: net[i].From.Idx}
		from, ok := fromNodeIO.Out[fromOutportAddr]
		if !ok {
			panic("not ok")
		}

		to := make([]chan Msg, len(net[i].To))
		for j := range net[i].To {
			toInportAddr := PortAddr{
				port: net[i].To[j].Port,
				idx:  net[i].To[j].Idx,
			}

			toNodeIO := nodesIO[net[i].To[j].Node]
			receiver, ok := toNodeIO.In[toInportAddr]
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

	return ss
}

func (r Runtime) nodeIO(node program.NodeMeta) IO {
	in := make(map[PortAddr]chan Msg)
	for port, size := range node.In {
		if size > 0 {
			for i := uint8(0); i < size; i++ {
				in[PortAddr{
					port: port,
					idx:  i,
				}] = make(chan Msg)
			}
			continue
		}

		in[PortAddr{
			port: port,
		}] = make(chan Msg)
	}

	out := make(map[PortAddr]chan Msg)
	for port, size := range node.Out {
		if size > 0 {
			for i := uint8(0); i < size; i++ {
				out[PortAddr{
					port: port,
					idx:  i,
				}] = make(chan Msg)
			}
			continue
		}

		out[PortAddr{
			port: port,
		}] = make(chan Msg)
	}

	return IO{in, out}
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

func (r Runtime) connectMany(cc []connection) {
	for i := range cc {
		go r.connect(cc[i])
	}
}

func (r Runtime) connect(s connection) {
	for msg := range s.from {
		for _, recv := range s.to {
			select {
			case recv <- msg:
				continue
			default:
				go func(to chan Msg, m Msg) {
					to <- m
				}(recv, msg)
			}
		}
	}
}

// Operator is a function that uses io provided by runtime.
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

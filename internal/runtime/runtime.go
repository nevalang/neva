package runtime

import (
	"fmt"

	"github.com/emil14/stream/internal/runtime/program"
)

type Runtime struct {
	Operators map[string]func(IO) error
}

func (r Runtime) Run(p program.Program) (IO, error) {
	return r.connectNode(p.Components, p.Root)
}

func (r Runtime) connectNode(scope map[string]program.Component, node program.Node) (IO, error) {
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
	nodesIO := map[string]IO{
		"in":  IO{Out: in},
		"out": IO{In: out},
	}

	for workerNode, meta := range component.Workers {
		io, err := r.connectNode(scope, meta)
		if err != nil {
			return IO{}, err
		}
		nodesIO[workerNode] = io
	}

	r.startStreams(
		r.streams(nodesIO, component.Net),
	)

	return IO{in, out}, nil
}

func (r Runtime) streams(nodesIO map[string]IO, net []program.Stream) []stream {
	ss := make([]stream, len(net))
	for i := range net {
		ss[i] = stream{
			// from: ,
			// to: ,
		}
		// s.From
		// s.To
	}

	return nil
}

func (r Runtime) nodeIO(node program.Node) (in Ports, out Ports) {
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

func (r Runtime) connectOperator(name string, node program.Node) (IO, error) {
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

func (r Runtime) startStreams(ss []stream) {
	for i := range ss {
		go r.startStream(ss[i])
	}
}

func (r Runtime) startStream(s stream) {
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

type Msg struct {
	Str  string
	Int  int
	Bool bool
}

type stream struct {
	from chan Msg
	to   []chan Msg
}

type IO struct {
	In, Out Ports
}

type Ports map[PortAddr]chan Msg

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

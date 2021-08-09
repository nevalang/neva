package runtime

import (
	"fmt"
)

type Runtime struct {
	Operators map[string]func(RuntimeIO) error
}

func (r Runtime) Run(scope map[string]Component, meta Node) (RuntimeIO, error) {
	component, ok := scope[meta.Component]
	if !ok {
		return RuntimeIO{}, fmt.Errorf(meta.Component)
	}

	if component.Operator != "" {
		io, err := r.connectOperator(component.Operator, meta)
		if err != nil {
			return RuntimeIO{}, fmt.Errorf("...")
		}
		return io, nil
	}

	nodesIO := map[string]RuntimeIO{}
	for node, nodeMeta := range component.Workers {
		io, err := r.Run(scope, nodeMeta)
		if err != nil {
			return RuntimeIO{}, err
		}

		nodesIO[node] = io
	}

	inportsNodeIO := RuntimeIO{}

	for ioPort, nodePortAddr := range meta.In {
		node, ok := component.Workers[nodePortAddr.node]
		if !ok {
			return RuntimeIO{}, nil
		}

		size, ok := node.Out[ioPort]
		if !ok {
			return RuntimeIO{}, nil
		}

		if size > 0 {
			for i := uint8(0); i < size; i++ {
				inportsNodeIO.In[IOPortAddr{
					port: ioPort,
					idx:  i,
				}] = make(chan Msg)
			}
			continue
		}

		inportsNodeIO.In[IOPortAddr{
			port: ioPort,
		}] = make(chan Msg)
	}

	outportsNodeIO := RuntimeIO{}

	for ioPort, nodePortAddr := range component.IO.out {
		node, ok := component.Workers[nodePortAddr.node]
		if !ok {
			return RuntimeIO{}, nil
		}

		size, ok := node.In[ioPort]
		if !ok {
			return RuntimeIO{}, nil
		}

		if size > 0 {
			for i := uint8(0); i < size; i++ {
				outportsNodeIO.Out[IOPortAddr{
					port: ioPort,
					idx:  i,
				}] = make(chan Msg)
			}
			continue
		}

		outportsNodeIO.Out[IOPortAddr{
			port: ioPort,
		}] = make(chan Msg)
	}

	r.startStreams(
		r.streams(nodesIO, component.Net),
	)

	return RuntimeIO{}, nil
}

func (r Runtime) connectOperator(name string, meta Node) (RuntimeIO, error) {
	connector, ok := r.Operators[name]
	if !ok {
		return RuntimeIO{}, fmt.Errorf("ErrUnknownOperator: %s", name)
	}

	io := RuntimeIO{
		In:  map[IOPortAddr]chan Msg{},
		Out: map[IOPortAddr]chan Msg{},
	}

	for port, size := range meta.In {
		if size > 0 {
			for i := uint8(0); i < size; i++ {
				io.In[IOPortAddr{
					port: port,
					idx:  i,
				}] = make(chan Msg)
			}
			continue
		}

		io.In[IOPortAddr{
			port: port,
		}] = make(chan Msg)
	}

	for port, size := range meta.Out {
		if size > 0 {
			for i := uint8(0); i < size; i++ {
				io.Out[IOPortAddr{
					port: port,
					idx:  i,
				}] = make(chan Msg)
			}
			continue
		}

		io.In[IOPortAddr{
			port: port,
		}] = make(chan Msg)
	}

	if err := connector(io); err != nil {
		return RuntimeIO{}, err
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

type RuntimeIO struct {
	In, Out Ports
}

type Ports map[IOPortAddr]chan Msg

func (p Ports) Slots(arrPort string) []chan Msg {
	cc := []chan Msg{}
	for addr, ch := range p {
		if addr.port == arrPort {
			cc = append(cc, ch)
		}
	}
	return cc
}

func (p Ports) Port(port string) chan Msg {
	for addr, ch := range p {
		if addr.port == port {
			return ch
		}
	}
	return nil
}

type IOPortAddr struct {
	port string
	idx  uint8
}

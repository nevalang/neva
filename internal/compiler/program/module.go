package program

import (
	"fmt"
	"sync"
)

type Module struct {
	IO      IO
	DepsIO  map[string]IO
	Const   map[string]Const
	Workers map[string]string
	Net     Connections
	Start   bool
}

func (mod Module) PairPortTypes(pair PortAddrPair) (PortType, PortType, error) {
	fromType, err := mod.NodeOutportType(pair.From.Node, pair.From.Port)
	if err != nil {
		return PortType{}, PortType{}, fmt.Errorf("outport: %w", err)
	}

	toType, err := mod.NodeInportType(pair.To.Node, pair.To.Port)
	if err != nil {
		return PortType{}, PortType{}, fmt.Errorf("inport: %w", err)
	}

	return fromType, toType, nil
}

func (m Module) NodeInportType(node, port string) (PortType, error) {
	inports, err := m.NodeInports(node)
	if err != nil {
		return PortType{}, fmt.Errorf("could not get inports for node %s: %w", node, err)
	}

	portType, ok := inports[port]
	if !ok {
		return portType, fmt.Errorf("unknown port %s on node %s", port, node)
	}

	return portType, nil
}

func (m Module) NodeOutportType(node, port string) (PortType, error) {
	ports, err := m.NodeOutports(node)
	if err != nil {
		return PortType{}, fmt.Errorf("get outports for node %s: %w", node, err)
	}

	portType, ok := ports[port]
	if !ok {
		return portType, fmt.Errorf("unknown port %s on node %s", port, node)
	}

	return portType, nil
}

func (m Module) NodeOutports(node string) (Ports, error) {
	io, err := m.NodeIO(node)
	if err != nil {
		return nil, err
	}
	return io.Out, nil
}

func (m Module) NodeInports(node string) (Ports, error) {
	io, err := m.NodeIO(node)
	if err != nil {
		return nil, err
	}
	return io.In, nil
}

func (m Module) NodeIO(node string) (IO, error) {
	if node == "in" {
		return IO{
			Out: m.IO.In,
		}, nil
	}

	if node == "out" {
		return IO{
			In: m.IO.Out,
		}, nil
	}

	if node == "const" {
		return m.ConstIO(), nil
	}

	dep, ok := m.Workers[node]
	if !ok {
		return IO{}, fmt.Errorf("unknown worker node %s", node)
	}

	io, ok := m.DepsIO[dep]
	if !ok {
		return IO{}, fmt.Errorf("unknown worker dep %s", dep)
	}

	return io, nil
}

func (m Module) ConstIO() IO {
	out := Ports{}
	for k, cnst := range m.Const {
		out[k] = PortType{Type: cnst.Type()}
	}
	return IO{Out: out}
}

type Connections map[PortAddr]map[PortAddr]struct{}

func (net Connections) CountOutgoing(node, outport string) uint8 {
	var c uint8
	for from := range net {
		if from.Node == node && from.Port == outport {
			c++
		}
	}
	return c
}

func (net Connections) IncomingConnections() IncomingConnections {
	incoming := IncomingConnections{}

	for pair := range net.Walk() {
		if incoming[pair.To] == nil {
			incoming[pair.To] = map[PortAddr]struct{}{}
		}
		incoming[pair.To][pair.From] = struct{}{}
	}

	return incoming
}

func (net Connections) Walk() <-chan PortAddrPair {
	ch := make(chan PortAddrPair, len(net))
	wg := sync.WaitGroup{}

	go func() {
		for from, receivers := range net {
			for to := range receivers {
				wg.Add(1)
				go func(from, to PortAddr) {
					ch <- PortAddrPair{from, to}
					wg.Done()
				}(from, to)
			}
		}

		wg.Wait()
		close(ch)
	}()

	return ch
}

type IncomingConnections Connections

func (rnet IncomingConnections) add(pair PortAddrPair) {
	if rnet[pair.To] == nil {
		rnet[pair.To] = map[PortAddr]struct{}{}
	}
	rnet[pair.To][pair.From] = struct{}{}
}

type PortAddrPair struct{ From, To PortAddr }

type PortAddr struct {
	Node string
	Port string
	Idx  uint8
}

func NewModule(
	io IO,
	deps map[string]IO,
	workers map[string]string,
	net Connections,
) Module {
	return Module{
		DepsIO:  deps,
		IO:      io,
		Workers: workers,
		Net:     net,
	}
}

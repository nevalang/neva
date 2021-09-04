package program

import (
	"fmt"
	"sync"
)

// Modules is a component that depends on other components.
type Modules struct {
	IO      IO
	Deps    map[string]IO
	Workers map[string]string
	Net     Net
}

func (cm Modules) Interface() IO {
	return cm.IO
}

func (m Modules) NodePortType(dir, node, port string) (PortType, error) {
	var f func(string) (Ports, error)

	switch dir {
	case "in":
		f = m.NodeInports
	case "out":
		f = m.NodeOutports
	default:
		return PortType{}, fmt.Errorf("dir should be in or out, got: %s", dir)
	}

	ports, err := f(node)
	if err != nil {
		return PortType{}, fmt.Errorf("could not get %s ports for node %s: %w", dir, node, err)
	}

	portType, ok := ports[port]
	if !ok {
		return portType, fmt.Errorf("unknown port %s on node %s", port, node)
	}

	return portType, nil
}

func (m Modules) NodeOutports(node string) (Ports, error) {
	io, err := m.NodeIO(node)
	if err != nil {
		return nil, err
	}
	return io.Out, nil
}

func (m Modules) NodeInports(node string) (Ports, error) {
	io, err := m.NodeIO(node)
	if err != nil {
		return nil, err
	}
	return io.In, nil
}

func (m Modules) NodeIO(node string) (IO, error) {
	if node == "in" || node == "out" {
		return m.IO, nil
	}

	dep, ok := m.Workers[node]
	if !ok {
		return IO{}, fmt.Errorf("unknown worker node %s", node)
	}

	io, ok := m.Deps[dep]
	if !ok {
		return IO{}, fmt.Errorf("unknown worker dep %s", dep)
	}

	return io, nil
}

// Net maps outport to set of inports.
type Net map[PortAddr]map[PortAddr]struct{}

// Walk implements iterator pattern to allow traverse network.
func (net Net) Walk() <-chan Rendezvous {
	ch := make(chan Rendezvous, len(net))
	wg := sync.WaitGroup{}

	go func() {
		for from, receivers := range net {
			for to := range receivers {
				wg.Add(1)
				go func(to PortAddr) {
					ch <- Rendezvous{from, to}
					wg.Done()
				}(to)
			}
		}

		wg.Wait()
		close(ch)
	}()

	return ch
}

// Incoming returns count of incoming connections for the given port.
// It works for array ports as well.
// It always returns 0 when non-existing port given.
func (net Net) Incoming(node string, inport string) uint8 {
	var c uint8
	for _, to := range net {
		for portAddr := range to {
			if portAddr.Node == node && portAddr.Port == inport {
				c++
			}
		}
	}
	return c
}

// Rendezvous represents from-to port addresses pair.
type Rendezvous struct{ From, To PortAddr }

// PortAddr is a point on a network graph.
type PortAddr struct {
	Node string
	Port string
	Idx  uint8
}

func (p PortAddr) String() string {
	return fmt.Sprintf("%s.%s[%d]", p.Node, p.Port, p.Idx)
}

func NewModule(
	io IO,
	deps map[string]IO,
	workers map[string]string,
	net Net,
) Modules {
	return Modules{
		Deps:    deps,
		IO:      io,
		Workers: workers,
		Net:     net,
	}
}

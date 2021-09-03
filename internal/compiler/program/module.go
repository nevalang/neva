package program

import "fmt"

// Modules is a component that depends on other components.
type Modules struct {
	IO      IO
	Deps    ComponentsIO
	Workers map[string]string
	Net     Net
}

// IO returns Module input-output interface.
func (cm Modules) Interface() IO {
	return cm.IO
}

// ComponentsIO maps component name with it's io interface.
type ComponentsIO map[string]IO

// Net maps outport to set of inports.
type Net map[PortAddr]map[PortAddr]struct{}

// Incoming returns count of incoming connections for the given port.
// It also works for array ports.
// If non-existing port given it always returns 0.
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

// PortAddr is a point on a network graph.
type PortAddr struct {
	Node string
	Port string
	Idx  uint8
}

func (p PortAddr) String() string {
	return fmt.Sprintf("%s.%s[%d]", p.Node, p.Port, p.Idx)
}

// todo need?
func NewModule(
	io IO,
	deps ComponentsIO,
	workers map[string]string,
	net Net,
) (Modules, error) {
	mod := Modules{
		Deps:    deps,
		IO:      io,
		Workers: workers,
		Net:     net,
	}

	return mod, nil // todo err?
}

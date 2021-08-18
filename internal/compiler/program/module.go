package program

// Module is a component that depends on other components.
type Module struct {
	io      IO // fixme (field name)
	Deps    ComponentsIO
	Workers map[string]string
	Net     Net
}

// IO returns Module input-output interface.
func (cm Module) IO() IO {
	return cm.io
}

// ComponentsIO maps component name with it's io interface.
type ComponentsIO map[string]IO

// Net maps outport to set of inports.
type Net map[PortAddr]map[PortAddr]struct{}

// PortAddr is a point on a network graph.
type PortAddr struct {
	Node string
	Port string
	Idx  uint8
}

// todo need?
func NewModule(
	io IO,
	deps ComponentsIO,
	workers map[string]string,
	net Net,
) (Module, error) {
	mod := Module{
		Deps:    deps,
		io:      io,
		Workers: workers,
		Net:     net,
	}

	return mod, nil // todo err?
}

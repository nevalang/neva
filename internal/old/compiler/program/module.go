package program

import (
	"fmt"
)

type Module struct {
	IO    ComponentIO
	Deps  map[string]ComponentIO
	Nodes ModuleNodes
	Net   ModuleNet
}

type ModuleNodes struct {
	Const   map[string]Const
	Workers map[string]string
}

func (mod Module) ConnectionTypes(c Connection) (PortType, PortType, error) {
	fromType, err := mod.OutPortType(c.From.Node, c.From.Port)
	if err != nil {
		return PortType{}, PortType{}, fmt.Errorf("outport: %w", err)
	}

	toType, err := mod.InPortType(c.To.Node, c.To.Port)
	if err != nil {
		return PortType{}, PortType{}, fmt.Errorf("inport: %w", err)
	}

	return fromType, toType, nil
}

func (m Module) InPortType(node, inport string) (PortType, error) {
	inports, err := m.NodeInports(node)
	if err != nil {
		return PortType{}, fmt.Errorf("could not get inports for node %s: %w", node, err)
	}

	inPortType, ok := inports[inport]
	if !ok {
		return inPortType, fmt.Errorf("unknown port %s on node %s", inport, node)
	}

	return inPortType, nil
}

func (m Module) OutPortType(node, outport string) (PortType, error) {
	outports, err := m.NodeOutports(node)
	if err != nil {
		return PortType{}, fmt.Errorf("get outports for node %s: %w", node, err)
	}

	outPortType, ok := outports[outport]
	if !ok {
		return outPortType, fmt.Errorf("unknown port %s on node %s", outport, node)
	}

	return outPortType, nil
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

func (m Module) NodeIO(node string) (ComponentIO, error) {
	if node == "in" {
		return ComponentIO{
			Out: m.IO.In,
		}, nil
	}

	if node == "out" {
		return ComponentIO{
			In: m.IO.Out,
		}, nil
	}

	if node == "const" {
		return m.ConstIO(), nil
	}

	dep, ok := m.Nodes.Workers[node]
	if !ok {
		return ComponentIO{}, fmt.Errorf("unknown worker node %s", node)
	}

	io, ok := m.Deps[dep]
	if !ok {
		return ComponentIO{}, fmt.Errorf("unknown worker dep %s", dep)
	}

	return io, nil
}

func (m Module) ConstIO() ComponentIO {
	out := Ports{}
	for k, cnst := range m.Nodes.Const {
		out[k] = PortType{DataType: cnst.Type()}
	}
	return ComponentIO{Out: out}
}

type ModuleNet map[ConnectionPoint][]ConnectionPoint

type ConnectionPoint struct {
	Type ConnectionPointType
	Node string
	Port string
	Idx  uint8
}

type ConnectionPointType uint8

const (
	SinglePort  ConnectionPointType = iota + 1 // use idx
	ArrayBypass                                // don't use idx
)

type Connection struct{ From, To ConnectionPoint }

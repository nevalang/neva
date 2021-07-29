package core

import (
	"errors"
)

// type NodesIO struct {
// 	m map[string]NodeIO
// }

// func newNodeIO() NodeIO {
// 	return NetworkIO{
// 		m: map[string]NetworkIO{},
// 	}
// }

type NodesIO map[string]NodeIO

// func (io NodesIO) Set(k string, v NodeIO) error {
// 	switch k {
// 	case "in":
// 		if v.in != nil {
// 			return fmt.Errorf("in not nil")
// 		}
// 		io[k] = v
// 		return nil
// 	case "out":
// 		if v.out != nil {
// 			return fmt.Errorf("in not nil")
// 		}
// 		io[k] = v
// 		return nil
// 	}

// 	if v.in == nil {
// 		return fmt.Errorf("in nil")
// 	}
// 	if v.out == nil {
// 		return fmt.Errorf("out nil")
// 	}

// 	io[k] = v
// 	return nil
// }

type NodeIO struct {
	in  nodeInports
	out nodeOutports
}

// func NewNodeIO(in nodeInports, out nodeOutports) (NodeIO, error) {
// 	if in == nil && out == nil {
// 		return NodeIO{}, fmt.Errorf("node io: in and out are both nil")
// 	}

// 	return NodeIO{in, out}, nil
// }

func (io NodeIO) Inport(name string) (chan Msg, error) {
	np, err := io.normPort(nodePorts(io.in), name)
	if err != nil {
		return nil, errors.New("")
	}

	return np, nil
}

func (io NodeIO) Outport(name string) (chan Msg, error) {
	np, err := io.normPort(nodePorts(io.out), name)
	if err != nil {
		return nil, errors.New("")
	}

	return np, nil
}

func (io NodeIO) ArrInport(name string) ([]chan Msg, error) {
	np, err := io.arrPort(nodePorts(io.in), name)
	if err != nil {
		return nil, errors.New("")
	}

	return np, nil
}

func (io NodeIO) ArrOutport(name string) ([]chan Msg, error) {
	np, err := io.arrPort(nodePorts(io.out), name)
	if err != nil {
		return nil, errors.New("")
	}

	return np, nil
}

func (io NodeIO) normPort(ports nodePorts, name string) (chan Msg, error) {
	port, ok := ports[name]
	if !ok {
		return nil, errors.New("")
	}

	norm, ok := port.(chan Msg)
	if !ok {
		return nil, errors.New("")
	}

	return norm, nil
}

func (io NodeIO) arrPort(ports nodePorts, name string) ([]chan Msg, error) {
	port, ok := ports[name]
	if !ok {
		return nil, errors.New("")
	}

	arr, ok := port.([]chan Msg)
	if !ok {
		return nil, errors.New("")
	}

	return arr, nil
}

type nodeInports nodePorts

type nodeOutports nodePorts

type nodePorts map[string]interface{}

type Msg struct {
	Str  string
	Int  int
	Bool bool
}

type relations struct {
	Sender    chan Msg
	Receivers []chan Msg
}

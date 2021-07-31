package core

import (
	"fmt"
)

type NodeIO struct {
	in  nodeInports
	out nodeOutports
}

func (io NodeIO) NormInport(name string) (chan Msg, error) {
	return nodePorts(io.in).normPort(name)
}

func (io NodeIO) NormOut(name string) (chan Msg, error) {
	return nodePorts(io.out).normPort(name)
}

func (io NodeIO) ArrIn(name string) ([]chan Msg, error) {
	return nodePorts(io.in).arrPort(name)
}

func (io NodeIO) ArrOutport(name string) ([]chan Msg, error) {
	return nodePorts(io.out).arrPort(name)
}

type nodeInports nodePorts

type nodeOutports nodePorts

type nodePorts map[string]interface{}

func (ports nodePorts) normPort(name string) (chan Msg, error) {
	port, ok := ports[name]
	if !ok {
		return nil, fmt.Errorf("port '%s' not found", name)
	}

	norm, ok := port.(chan Msg)
	if !ok {
		return nil, fmt.Errorf("normal port expected, got %T", port)
	}

	return norm, nil
}

func (ports nodePorts) arrPort(name string) ([]chan Msg, error) {
	port, ok := ports[name]
	if !ok {
		return nil, fmt.Errorf("port '%s' not found", name)
	}

	arr, ok := port.([]chan Msg)
	if !ok {
		return nil, fmt.Errorf("array port expected, got %T", port)
	}

	return arr, nil
}

type Msg struct {
	Str  string
	Int  int
	Bool bool
}

type relations struct {
	Sender    chan Msg
	Receivers []chan Msg
}

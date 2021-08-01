package core

import (
	"fmt"
)

type NodeIO struct {
	In  nodeInports
	Out nodeOutports
}

func (io NodeIO) NormIn(name string) (chan Msg, error) {
	return nodePorts(io.In).normPort(name)
}

func (io NodeIO) NormOut(name string) (chan Msg, error) {
	return nodePorts(io.Out).normPort(name)
}

func (io NodeIO) ArrIn(name string) ([]chan Msg, error) {
	return nodePorts(io.In).arrPort(name)
}

func (io NodeIO) ArrOut(name string) ([]chan Msg, error) {
	return nodePorts(io.Out).arrPort(name)
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

type stream struct {
	Sender    chan Msg
	Receivers []chan Msg
}

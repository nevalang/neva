package core

import (
	"fmt"
)

type NodeIO struct {
	in  nodeInports
	out nodeOutports
}

func (io NodeIO) NormIn(name string) (chan Msg, error) {
	return nodePorts(io.in).norm(name)
}

func (io NodeIO) NormOut(name string) (chan Msg, error) {
	return nodePorts(io.out).norm(name)
}

func (io NodeIO) ArrIn(name string) ([]chan Msg, error) {
	return nodePorts(io.in).arr(name)
}

func (io NodeIO) ArrOut(name string) ([]chan Msg, error) {
	return nodePorts(io.out).arr(name)
}

type nodeInports nodePorts

type nodeOutports nodePorts

type nodePorts map[string]interface{}

func (ports nodePorts) norm(name string) (chan Msg, error) {
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

func (ports nodePorts) arr(name string) ([]chan Msg, error) {
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

// func NewStrMsg(s string) Msg {
// 	return Msg{Str: s}
// }

// func NewIntMsg(n int) Msg {
// 	return Msg{Int: n}
// }

// func NewBoolMsg(b bool) Msg {
// 	return Msg{Bool: b}
// }

type stream struct {
	Sender    chan Msg
	Receivers []chan Msg
}

package core

import (
	"errors"
	"fmt"
)

var ErrPortNotFound = errors.New("port not found")

type IO struct {
	In, Out Ports
}

// FIXME: map[string][]chan Msg (ordering)
type Ports map[RelativePortAddr]chan Msg

// FIXME: must only use port name (not full path) because of how this is used by operators
func (p Ports) Port(name string) (chan Msg, error) {
	port, ok := p[RelativePortAddr{Port: name}]
	if !ok {
		return nil, fmt.Errorf("%w: in: %v", ErrPortNotFound, name)
	}
	return port, nil
}

// FIXME add sorting?
func (p Ports) ArrPortSlots(name string) ([]chan Msg, error) {
	pp := make([]chan Msg, 0, len(p))

	for addr, port := range p {
		if addr.Port == name {
			pp = append(pp, port)
		}
	}

	if len(pp) == 0 {
		return nil, fmt.Errorf("%w: in: %v", ErrPortNotFound, name)
	}

	return pp, nil
}

type RelativePortAddr struct {
	Port string
	Idx  uint8
}

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
type Ports map[PortAddr]chan Msg

func (p Ports) Port(name string) (chan Msg, error) {
	port, ok := p[PortAddr{Port: name}]
	if !ok {
		return nil, fmt.Errorf("%w: in: %v", ErrPortNotFound, name)
	}
	return port, nil
}

// FIXME sorting?
func (p Ports) ArrPort(name string) ([]chan Msg, error) {
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

type PortAddr struct {
	Port string
	Idx  uint8
}

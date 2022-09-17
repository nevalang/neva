package core

import (
	"errors"
	"fmt"
	"sort"
)

var ErrPortNotFound = errors.New("port not found")

type IO struct {
	In, Out Ports
}

type Ports map[RelativePortAddr]chan Msg

func (p Ports) Port(name string) (chan Msg, error) {
	port, ok := p[RelativePortAddr{Port: name}]
	if !ok {
		return nil, fmt.Errorf("%w: in: %v", ErrPortNotFound, name)
	}
	return port, nil
}

func (p Ports) ArrPortSlots(name string) ([]chan Msg, error) {
	type port struct {
		addr RelativePortAddr
		ch   chan Msg
	}

	pp := make([]port, 0, len(p))

	for addr, ch := range p {
		if addr.Port == name {
			pp = append(pp, port{addr, ch})
		}
	}

	if len(pp) == 0 {
		return nil, fmt.Errorf("%w: in: %v", ErrPortNotFound, name)
	}

	sort.Slice(pp, func(i, j int) bool {
		return pp[i].addr.Idx < pp[j].addr.Idx
	})

	cc := make([]chan Msg, len(pp))
	for i := range pp {
		cc[i] = pp[i].ch
	}

	return cc, nil
}

type RelativePortAddr struct {
	Port string
	Idx  uint8
}

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

// Port returns port with given name and idx == 0 or non-nil err
func (p Ports) Port(name string) (chan Msg, error) {
	port, ok := p[RelativePortAddr{Port: name}]
	if !ok {
		return nil, fmt.Errorf("%w: in: %v", ErrPortNotFound, name)
	}
	return port, nil
}

// ArrPortSlots returns all ports with given name sorted by idx or non-nil err
func (p Ports) ArrPortSlots(name string) ([]chan Msg, error) {
	type port struct {
		idx uint8
		ch  chan Msg
	}

	pp := make([]port, 0, len(p))

	for addr, ch := range p {
		if addr.Port == name {
			pp = append(pp, port{addr.Idx, ch})
		}
	}

	if len(pp) == 0 {
		return nil, fmt.Errorf("%w: in: %v", ErrPortNotFound, name)
	}

	sort.Slice(pp, func(i, j int) bool {
		return pp[i].idx < pp[j].idx
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

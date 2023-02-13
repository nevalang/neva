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

type Ports map[PortAddr]chan Msg

// SinglePort returns port with given name and idx 0 if found, otherwise error.
// It doesn't care if there are more ports with given name.
func (p Ports) SinglePort(name string) (chan Msg, error) {
	v, ok := p[PortAddr{Port: name}]
	if !ok {
		return nil, fmt.Errorf("port %s not found", name)
	}
	return v, nil
}

// ArrPortSlots returns slice of ports with given name sorted by idx.
// If zero ports found it returns error.
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
		return nil, fmt.Errorf("array port %s not found", name)
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

type PortAddr struct {
	Port string
	Idx  uint8
}

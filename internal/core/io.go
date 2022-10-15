package core

import (
	"errors"
	"sort"
)

var ErrPortNotFound = errors.New("port not found")

type IO struct {
	In, Out Ports
}

type Ports map[PortAddr]chan Msg

func (p Ports) Port(name string) chan Msg {
	return p[PortAddr{Port: name}]
}

func (p Ports) ArrPortSlots(name string) []chan Msg {
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

	sort.Slice(pp, func(i, j int) bool {
		return pp[i].idx < pp[j].idx
	})

	cc := make([]chan Msg, len(pp))
	for i := range pp {
		cc[i] = pp[i].ch
	}

	return cc
}

type PortAddr struct {
	Port string
	Idx  uint8
}

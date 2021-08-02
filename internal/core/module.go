package core

import (
	"errors"
	"fmt"
)

type module struct {
	deps    Interfaces
	in      InportsInterface
	out     OutportsInterface
	workers map[string]string
	net     Net
}

func (mod module) Incoming(p Port) uint8 {
	return 0
}

func (cm module) Interface() Interface {
	return Interface{
		In:  cm.in,
		Out: cm.out,
	}
}

func (mod module) Validate() error {
	if err := mod.validatePorts(mod.in, mod.out); err != nil {
		return err
	}

	return nil
}

func (mod module) validatePorts(in InportsInterface, out OutportsInterface) error {
	if len(in) == 0 || len(out) == 0 {
		return fmt.Errorf("ports len 0")
	}

	// TODO check arr points - no holes should be

	return nil
}

type Interfaces map[string]Interface

type Net map[PortPoint]map[PortPoint]struct{}

func (net Net) Incoming(p PortPoint) (c uint8) {
	for _, rr := range net {
		_, ok := rr[p]
		if ok {
			c++
		}
	}

	return c
}

type PortPoint interface {
	Node() string
	// Port() string
}

type NormPortPoint struct {
	node string
	port string
}

func NewNormPortPoint(node, port string) (NormPortPoint, error) {
	if node == "" || port == "" {
		return NormPortPoint{}, fmt.Errorf("invalid normal port point")
	}

	return NormPortPoint{
		port: port,
		node: node,
	}, nil
}

func (p NormPortPoint) Node() string {
	return p.node
}

func (p NormPortPoint) Port() string {
	return p.port
}

type ArrPortPoint struct {
	node string
	port string
	idx  uint8
}

func NewArrPortPoint(node, port string, idx uint64) (ArrPortPoint, error) {
	if node == "" || port == "" || idx > 255 {
		return ArrPortPoint{}, errors.New("invalid array port point")
	}

	return ArrPortPoint{
		node: node,
		port: port,
		idx:  uint8(idx),
	}, nil
}

func (p ArrPortPoint) Node() string {
	return p.node
}

func (p ArrPortPoint) Port() string {
	return p.port
}

func (p ArrPortPoint) Idx() uint8 {
	return p.idx
}

func NewCustomModule(
	deps Interfaces,
	in InportsInterface,
	out OutportsInterface,
	workers map[string]string,
	net Net,
) (Component, error) {
	mod := module{
		deps:    deps,
		in:      in,
		out:     out,
		workers: workers,
		net:     net,
	}

	if err := mod.Validate(); err != nil {
		return nil, err
	}

	return mod, nil
}

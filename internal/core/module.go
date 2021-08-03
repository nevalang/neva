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

type Net map[PortAddr]map[PortAddr]struct{}

// TODO: check if that is not arrport point.
func (net Net) ArrInSize(node, port string) uint8 {
	var size uint8

	for _, rr := range net {
		for receiver := range rr {
			if receiver.Node() == node && receiver.Port() == port {
				size++
			}
		}
	}

	return size
}

func (net Net) ArrOutSize(node, port string) uint8 {
	var size uint8

	for sender := range net {
		if sender.Node() == node && sender.Port() == port {
			size++
		}
	}

	return size
}

type PortAddr interface {
	Node() string
	Port() string
	Compare(PortAddr) bool
}

type NormPortAddr struct {
	node string
	port string
}

func NewNormPortPoint(node, port string) (NormPortAddr, error) {
	if node == "" || port == "" {
		return NormPortAddr{}, fmt.Errorf("invalid normal port point")
	}

	return NormPortAddr{
		port: port,
		node: node,
	}, nil
}

func (p NormPortAddr) Node() string {
	return p.node
}

func (p NormPortAddr) Port() string {
	return p.port
}

func (p NormPortAddr) Compare(got PortAddr) bool {
	norm, ok := got.(NormPortAddr)
	if !ok {
		return false
	}

	return norm.node == got.Node() && norm.port == got.Port()
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

func (p ArrPortPoint) Compare(got PortAddr) bool {
	arr, ok := got.(ArrPortPoint)
	if !ok {
		return false
	}

	return arr.node == got.Node() && arr.port == got.Port() && arr.idx == arr.Idx()
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

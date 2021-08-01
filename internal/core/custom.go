package core

import (
	"errors"
	"fmt"
)

type customModule struct {
	deps    Interfaces
	in      InportsInterface
	out     OutportsInterface
	workers Workers
	net     []StreamDef
}

func (cm customModule) Interface() Interface {
	return Interface{
		In:  cm.in,
		Out: cm.out,
	}
}

func (mod customModule) Validate() error {
	if err := mod.validatePorts(mod.in, mod.out); err != nil {
		return err
	}

	return nil
}

func (mod customModule) validatePorts(in InportsInterface, out OutportsInterface) error {
	if len(in) == 0 || len(out) == 0 {
		return fmt.Errorf("ports len 0")
	}
	return nil
}

type Interfaces map[string]Interface

func (d Interfaces) Compare(name string, io Interface) error {
	for port, t := range io.In {
		if err := d[name].In[port].Compare(t); err != nil {
			return err
		}
	}
	for port, t := range io.Out {
		if err := d[name].Out[port].Compare(t); err != nil {
			return err
		}
	}
	return nil
}

type Workers map[string]string

func (w Workers) Interface(name string, deps Interfaces) (Interface, error) {
	i, ok := deps[name]
	if !ok {
		return Interface{}, errors.New("..")
	}
	return i, nil
}

type StreamDef struct {
	Sender    PortPoint
	Recievers []PortPoint
}

type PortPoint interface {
	Node() string
	Port() string
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
	workers Workers,
	net []StreamDef,
) (Module, error) {
	mod := customModule{
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

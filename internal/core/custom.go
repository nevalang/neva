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
	net     Net
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

type Net []RelationsDef

type RelationsDef struct {
	Sender    PortPoint
	Recievers []PortPoint
}

type PortPoint interface {
	NodeName() string
}

type NormPortPoint struct {
	Node string
	Port string
}

func (p NormPortPoint) NodeName() string {
	return p.Node
}

type ArrPortPoint struct {
	Node  string
	Port  string
	Index uint8
}

func (p ArrPortPoint) NodeName() string {
	return p.Node
}

func NewCustomModule(
	deps Interfaces,
	in InportsInterface,
	out OutportsInterface,
	workers Workers,
	net Net,
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

package core

import (
	"errors"
	"fmt"

	"github.com/emil14/refactored-garbanzo/internal/types"
)

var (
	ErrModNotFound = errors.New("module not found in scope")
)

type Module interface {
	Interface() Interface
}

type Interface struct {
	In  InportsInterface
	Out OutportsInterface
}

type Deps map[string]Interface

func (d Deps) compat(name string, io Interface) error {
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

type InportsInterface PortsInterface

type OutportsInterface PortsInterface

type PortsInterface map[string]PortType

func (p1 PortsInterface) Compare(p2 PortsInterface) error {
	if len(p1) != len(p2) {
		return fmt.Errorf("differenet len")
	}
	for k, v := range p1 {
		if err := v.Compare(p2[k]); err != nil {
			return err
		}
	}
	return nil
}

type PortType struct {
	Type types.Type
	Arr  bool
}

func (p1 PortType) Compare(p2 PortType) error {
	if p1.Arr != p2.Arr || p1.Type != p2.Type {
		return fmt.Errorf("want: %v, got: %v", p1, p2)
	}
	return nil
}

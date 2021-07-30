package core

import (
	"fmt"

	"github.com/emil14/refactored-garbanzo/internal/types"
)

type Module interface {
	Interface() Interface
}

type Interface struct {
	In  InportsInterface
	Out OutportsInterface
}

func (want Interface) Compare(got Interface) error {
	if err := want.In.Compare(got.In); err != nil {
		return err
	}
	return want.Out.Compare(got.Out)
}

type InportsInterface PortsInterface

func (want InportsInterface) Compare(got InportsInterface) error {
	err := PortsInterface(want).Compare(PortsInterface(got))
	if err != nil {
		return fmt.Errorf("incompatible inports: %w", err)
	}

	return nil
}

type OutportsInterface PortsInterface

func (want OutportsInterface) Compare(got OutportsInterface) error {
	err := PortsInterface(want).Compare(PortsInterface(got))
	if err != nil {
		return fmt.Errorf("incompatible outports: %w", err)
	}

	return nil
}

type PortsInterface map[string]PortInterface

func (want PortsInterface) Compare(got PortsInterface) error {
	len1 := len(want)
	len2 := len(got)
	if len1 != len2 {
		return errPortsLen(len1, len2)
	}

	for name, typ := range want {
		if err := typ.Compare(got[name]); err != nil {
			return errPortInvalid(name, err)
		}
	}

	return nil
}

type PortType interface {
	Compare(PortType) error
}

type PortInterface struct { // TODO rename to ArrPortInterface
	Type types.Type
	Arr  bool
	Size uint8
}

func (pt PortInterface) String() (s string) {
	if pt.Arr {
		s += " array"
	}
	s += "port of type " + pt.Type.String()
	return s
}

func (p1 PortInterface) Compare(p2 PortInterface) error {
	if p1.Arr != p2.Arr || p1.Type != p2.Type {
		return errPortTypes(p1, p2)
	}
	return nil
}

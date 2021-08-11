package core

import (
	"fmt"

	"github.com/emil14/stream/internal/core/types"
)

type IO struct {
	In  Inports
	Out Outports
}

func (want IO) Compare(got IO) error {
	if err := Ports(want.In).Compare(
		Ports(got.In),
	); err != nil {
		return fmt.Errorf("inport: %w", err)
	}

	if err := Ports(want.Out).Compare(Ports(got.Out)); err != nil {
		return fmt.Errorf("outport: %w", err)
	}

	return nil
}

type Inports Ports

type Outports Ports

type Ports map[string]PortType

func (want Ports) Compare(got Ports) error {
	if len(want) != len(got) {
		return ErrPortsLen
	}

	for name, typ := range want {
		_, ok := got[name]
		if !ok {
			return ErrPortNotFound
		}

		if err := typ.Compare(got[name]); err != nil {
			return ErrPortInvalid
		}
	}

	return nil
}

func (ports Ports) ArrPorts() map[string]PortType {
	m := map[string]PortType{}

	for name, typ := range ports {
		if typ.Arr {
			m[name] = typ
		}
	}

	return m
}

type PortType struct {
	Type types.Type
	Arr  bool
}

func (want PortType) Compare(got PortType) error {
	if want.Arr != got.Arr || want.Type != got.Type {
		return ErrPortTypes
	}

	return nil
}

func (pt PortType) String() (s string) {
	if pt.Arr {
		s += "array"
	}

	return s + "port of type " + pt.Type.String()
}
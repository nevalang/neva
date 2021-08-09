package core

import (
	"fmt"

	"github.com/emil14/stream/internal/core/types"
)

type Component interface {
	Interface() ComponentInterface
}

type ComponentInterface struct {
	In  InportsInterface
	Out OutportsInterface
}

func (want ComponentInterface) Compare(got ComponentInterface) error {
	if err := PortsInterface(want.In).Compare(
		PortsInterface(got.In),
	); err != nil {
		return fmt.Errorf("inport: %w", err)
	}

	if err := PortsInterface(want.Out).Compare(PortsInterface(got.Out)); err != nil {
		return fmt.Errorf("outport: %w", err)
	}

	return nil
}

type InportsInterface PortsInterface

type OutportsInterface PortsInterface

type PortsInterface map[string]PortType

func (want PortsInterface) Compare(got PortsInterface) error {
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

func (ports PortsInterface) ArrPorts() map[string]PortType {
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

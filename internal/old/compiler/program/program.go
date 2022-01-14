package program

import "fmt"

type Program struct {
	Components    map[string]Component
	RootComponent string
}

type Component struct {
	Type       ComponentType
	OperatorIO ComponentIO
	Module     Module
}

func (c Component) IO() ComponentIO {
	if c.Type == ModuleComponent {
		return c.Module.IO
	}
	return c.OperatorIO
}

type ComponentType uint8

const (
	ModuleComponent ComponentType = iota + 1
	OperatorComponent
)

type ComponentIO struct {
	Params map[string]struct{}
	In     Ports
	Out    Ports
}

func (want ComponentIO) Compare(got ComponentIO) error {
	wantIn := Ports(want.In)
	gotIn := Ports(got.In)

	if err := wantIn.Compare(gotIn); err != nil {
		return fmt.Errorf("inport: %w", err)
	}

	wantOut := Ports(want.Out)
	gotOut := Ports(got.Out)

	if err := wantOut.Compare(gotOut); err != nil {
		return fmt.Errorf("outport: %w", err)
	}

	return nil
}

type Ports map[string]PortType

func (want Ports) Compare(got Ports) error {
	if len(want) != len(got) {
		return ErrPortsLen
	}

	for name, typ := range want {
		_, ok := got[name]
		if !ok {
			return fmt.Errorf("%w: %s", ErrPortNotFound, name)
		}

		if err := typ.Compare(got[name]); err != nil {
			return fmt.Errorf("%w: %s", ErrPortInvalid, name)
		}
	}

	return nil
}

type PortType struct {
	DataType DataType
	IsArr    bool
}

func (want PortType) Compare(got PortType) error {
	if want.IsArr != got.IsArr || want.DataType != got.DataType {
		return fmt.Errorf("%w: got %v, want %v", ErrPortTypes, got, want)
	}

	return nil
}

type DataType uint8

const (
	TypeInt DataType = iota + 1
	TypeStr
	TypeBool
	TypeSig
)

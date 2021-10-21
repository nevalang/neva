package program

import "fmt"

type Program struct {
	Root  string
	Scope map[string]Component
}

type Component interface {
	Interface() IO
}

type IO struct {
	In  Ports
	Out Ports
}

func (want IO) Compare(got IO) error {
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
	Type Type
	Arr  bool
}

func (want PortType) Compare(got PortType) error {
	if want.Arr != got.Arr || want.Type != got.Type {
		return fmt.Errorf("%w: got %v, want %v", ErrPortTypes, got, want)
	}

	return nil
}

func (p PortType) String() string {
	s := p.Type.String()
	if p.Arr {
		s += "[]"
	}
	return s
}

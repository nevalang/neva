package generator

import (
	"fmt"

	"github.com/emil14/refactored-garbanzo/internal/parser"
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

type Validator interface {
	ValidateDeps(pdeps parser.Deps, env runtime.Env) error
}

type validator struct{}

func (v validator) ValidateDeps(pdeps parser.Deps, env runtime.Env) error {
	for name := range pdeps {
		rmod, ok := env[name]
		if !ok {
			return fmt.Errorf("unresolved dep: '%s'", name)
		}

		rin, rout := rmod.Interface()
		if err := compareAllPorts(
			pdeps[name].In, pdeps[name].Out, rin, rout,
		); err != nil {
			return fmt.Errorf("incompatible ports on module '%s': %w", name, err)
		}
	}

	return nil
}

func compareAllPorts(
	pin parser.InportsInterface,
	pout parser.OutportsInterface,
	rin runtime.InportsInterface,
	rout runtime.OutportsInterface,
) error {
	if err := comparePorts(
		parser.PortsInterface(pin),
		runtime.PortsInterface(rin),
	); err != nil {
		return fmt.Errorf("incompatible inPorts: %w", err)
	}

	if err := comparePorts(
		parser.PortsInterface(pout),
		runtime.PortsInterface(rout),
	); err != nil {
		return fmt.Errorf("incompatible outPorts: %w", err)
	}

	return nil
}

func comparePorts(pports parser.PortsInterface, rports runtime.PortsInterface) error {
	if len(pports) != len(rports) {
		return fmt.Errorf(
			"different number of ports: want %d, got %d",
			len(rports),
			len(pports),
		)
	}

	for name := range pports {
		t := types.ByName(pports[name])
		if t == types.Unknown {
			return fmt.Errorf("unknown type '%s' on port '%s'", pports[name], name)
		}

		if t != rports[name] {
			return fmt.Errorf(
				"incompatible types on port '%s': want '%s', got '%s'",
				name,
				pports[name],
				t,
			)
		}
	}

	return nil
}

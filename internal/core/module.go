package core

import (
	"fmt"

	"github.com/emil14/stream/internal/core/types"
)

type module struct {
	deps    Interfaces
	in      InportsInterface
	out     OutportsInterface
	workers map[string]string
	net     Net
}

func (cm module) Interface() ComponentInterface {
	return ComponentInterface{
		In:  cm.in,
		Out: cm.out,
	}
}

func (mod module) Validate() error {
	if err := mod.validatePorts(mod.in, mod.out); err != nil {
		return err
	}

	if err := mod.validateDeps(mod.deps); err != nil {
		return err
	}

	if err := mod.validateWorkers(mod.deps, mod.workers); err != nil {
		return err
	}

	// TODO check arr points - should be no holes

	return nil
}

// validatePorts checks that ports are not empty and there is no unknown types.
func (mod module) validatePorts(in InportsInterface, out OutportsInterface) error {
	if len(in) == 0 || len(out) == 0 {
		return fmt.Errorf("ports len 0")
	}

	for port, t := range in {
		if t.Type == types.Unknown {
			return fmt.Errorf("unknown type " + port)
		}
	}

	for port, t := range out {
		if t.Type == types.Unknown {
			return fmt.Errorf("unknown type " + port)
		}
	}

	return nil
}

// validateWorkers checks that every worker points to existing dependency.
func (v module) validateWorkers(deps Interfaces, workers map[string]string) error {
	for workerName, depName := range workers {
		if _, ok := deps[depName]; !ok {
			return fmt.Errorf("invalid workers: worker '%s' points to unknown dependency '%s'", workerName, depName)
		}
	}

	return nil
}

// validateDeps validates ports of every dependency.
func (v module) validateDeps(deps Interfaces) error {
	for name, dep := range deps {
		if err := v.validatePorts(dep.In, dep.Out); err != nil {
			return fmt.Errorf("invalid dep '%s': %w", name, err)
		}
	}

	return nil
}

type Interfaces map[string]ComponentInterface

type Net map[PortAddr]map[PortAddr]struct{}

type PortAddr struct {
	Node string
	Port string
	Idx  uint8
}

func NewModule(
	io ComponentInterface,
	deps Interfaces,
	workers map[string]string,
	net Net,
) (Component, error) {
	mod := module{
		deps:    deps,
		in:      io.In,
		out:     io.Out,
		workers: workers,
		net:     net,
	}

	if err := mod.Validate(); err != nil {
		return nil, err
	}

	return mod, nil
}

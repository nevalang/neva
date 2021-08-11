package program

import (
	"fmt"
)

type Module struct {
	deps    Interfaces
	io      IO
	workers map[string]string
	net     Net
}

func (cm Module) IO() IO {
	return cm.io
}

// TODO check arr points - should be no holes
func (mod Module) Validate() error {
	if err := mod.validatePorts(mod.io); err != nil {
		return err
	}

	if err := mod.validateDeps(mod.deps); err != nil {
		return err
	}

	if err := mod.validateWorkers(mod.deps, mod.workers); err != nil {
		return err
	}

	return nil
}

// validatePorts checks that ports are not empty and there is no unknown
func (mod Module) validatePorts(io IO) error {
	if len(io.In) == 0 || len(io.Out) == 0 {
		return fmt.Errorf("ports len 0")
	}

	for port, t := range io.In {
		if t.Type == Unknown {
			return fmt.Errorf("unknown type " + port)
		}
	}

	for port, t := range io.Out {
		if t.Type == Unknown {
			return fmt.Errorf("unknown type " + port)
		}
	}

	return nil
}

// validateWorkers checks that every worker points to existing dependency.
func (v Module) validateWorkers(deps Interfaces, workers map[string]string) error {
	for workerName, depName := range workers {
		if _, ok := deps[depName]; !ok {
			return fmt.Errorf("invalid workers: worker '%s' points to unknown dependency '%s'", workerName, depName)
		}
	}

	return nil
}

// validateDeps validates ports of every dependency.
func (v Module) validateDeps(deps Interfaces) error {
	for name, dep := range deps {
		if err := v.validatePorts(dep); err != nil {
			return fmt.Errorf("invalid dep '%s': %w", name, err)
		}
	}

	return nil
}

type Interfaces map[string]IO

type Net map[PortAddr]map[PortAddr]struct{}

type PortAddr struct {
	Node string
	Port string
	Idx  uint8
}

func NewModule(
	io IO,
	deps Interfaces,
	workers map[string]string,
	net Net,
) (Module, error) {
	mod := Module{
		deps:    deps,
		io:      io,
		workers: workers,
		net:     net,
	}

	return mod, nil
}

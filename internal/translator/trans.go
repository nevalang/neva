package generator

import (
	"fmt"

	"fbp/internal/parsing"
	"fbp/internal/runtime"
	"fbp/internal/types"
)

type Translator interface {
	Translate(parsing.Module) (runtime.Module, error)
}

func New(env runtime.WorkerMap) Translator {
	return translator{env}
}

type translator struct {
	env runtime.WorkerMap
}

func (t translator) Translate(m parsing.Module) (runtime.Module, error) {
	rin, rout, err := t.translateAllPorts(m.In, m.Out)
	if err != nil {
		return runtime.Module{}, fmt.Errorf("invalid ports: %w", err)
	}

	rdeps, err := t.translateDeps(m.Deps)
	if err != nil {
		return runtime.Module{}, fmt.Errorf("unresolved deps: %w", err)
	}

	rwm, err := t.translateWorkers(rdeps, m.WorkerMap)
	if err != nil {
		return runtime.Module{}, fmt.Errorf("could not translate workers: %w", err)
	}

	rnet, err := t.translateNet(m.Net, rwm)
	if err != nil {
		return runtime.Module{}, fmt.Errorf("could not translate net: %w", err)
	}

	return runtime.NewModule(rin, rout, rwm, rnet), nil
}

func (t translator) translateNet(net parsing.Net, wm runtime.WorkerMap) ([]runtime.Conn, error) {
	return nil, nil // TODO
}

func (t translator) translateWorkers(deps map[string]runtime.AbsModule, wm map[string]string) (runtime.WorkerMap, error) {
	rwm := runtime.WorkerMap{}

	for workerName, depName := range wm {
		depMod, ok := deps[depName]
		if !ok {
			return nil, fmt.Errorf("dep '%s' not found for worker '%s'", depName, workerName)
		}
		rwm[workerName] = depMod
	}

	return rwm, nil
}

func (t translator) translateDeps(pdeps parsing.Deps) (map[string]runtime.AbsModule, error) {
	rdeps := map[string]runtime.AbsModule{}

	for name := range pdeps {
		rmod, ok := t.env[name]
		if !ok {
			return nil, fmt.Errorf("unresolved dep: '%s'", name)
		}

		rin, rout := rmod.Ports()
		if err := compatAllPorts(
			pdeps[name].In, pdeps[name].Out, rin, rout,
		); err != nil {
			return nil, fmt.Errorf("incompatible ports on module '%s': %w", name, err)
		}

		rdeps[name] = rmod
	}

	return rdeps, nil
}

// translateAllPorts returns error if at least one of ports has unknown type.
func (t translator) translateAllPorts(in parsing.InPorts, out parsing.OutPorts) (runtime.InPorts, runtime.OutPorts, error) {
	inPorts, err := t.translatePorts(parsing.PortMap(in))
	if err != nil {
		return runtime.InPorts{}, runtime.OutPorts{}, fmt.Errorf("could not translate inPorts: %w", err)
	}

	outPorts, err := t.translatePorts(parsing.PortMap(out))
	if err != nil {
		return runtime.InPorts{}, runtime.OutPorts{}, fmt.Errorf("could not translate outPorts: %w", err)
	}

	return runtime.InPorts(inPorts), runtime.OutPorts(outPorts), nil
}

func (t translator) translatePorts(ppm parsing.PortMap) (runtime.PortMap, error) {
	rpm := runtime.PortMap{}

	for portName, typeName := range ppm {
		t := types.ByName(typeName)
		if t == types.Unknown {
			return runtime.PortMap{}, fmt.Errorf("unknown type %s", typeName)
		}
		rpm[portName] = t
	}

	return rpm, nil
}

func compatAllPorts(
	pin parsing.InPorts,
	pout parsing.OutPorts,
	rin runtime.InPorts,
	rout runtime.OutPorts,
) error {
	if err := compatPorts(
		parsing.PortMap(pin),
		runtime.PortMap(rin),
	); err != nil {
		return fmt.Errorf("incompatible inPorts: %w", err)
	}

	if err := compatPorts(
		parsing.PortMap(pout),
		runtime.PortMap(rout),
	); err != nil {
		return fmt.Errorf("incompatible outPorts: %w", err)
	}

	return nil
}

func compatPorts(parsedPorts parsing.PortMap, runtimePorts runtime.PortMap) error {
	if len(parsedPorts) != len(runtimePorts) {
		return fmt.Errorf(
			"different number of ports: want %d, got %d",
			len(runtimePorts),
			len(parsedPorts),
		)
	}

	for name := range parsedPorts {
		t := types.ByName(parsedPorts[name])
		if t == types.Unknown {
			return fmt.Errorf("unknown type '%s' on port '%s'", parsedPorts[name], name)
		}

		if t != runtimePorts[name] {
			return fmt.Errorf(
				"incompatible types on port '%s': want '%s', got '%s'",
				name,
				parsedPorts[name],
				t,
			)
		}
	}

	return nil
}

package translator

import (
	"fmt"

	"fbp/internal/parsing"
	"fbp/internal/runtime"
	"fbp/internal/types"
)

type (
	Translator interface {
		Translate(parsing.Module) (runtime.ComplexModule, error)
	}

	translator struct {
		env runtime.Env
	}
)

func New(env runtime.Env) Translator {
	return translator{env}
}

func (t translator) Translate(pmod parsing.Module) (runtime.ComplexModule, error) {
	rdeps, err := t.translateDeps(pmod.Deps)
	if err != nil {
		return runtime.ComplexModule{}, fmt.Errorf("unresolved deps: %w", err)
	}

	rin, rout := t.translateAllPorts(pmod.In, pmod.Out)
	rworkers := t.translateWorkers(rdeps, pmod.Workers)
	rnet := t.translateNet(pmod.Net, rin, rout, rworkers)

	return runtime.NewModule(
		rin,
		rout,
		rworkers,
		rnet,
	), nil
}

func (t translator) translateDeps(pdeps parsing.Deps) (map[string]runtime.AbstractModule, error) {
	rdeps := map[string]runtime.AbstractModule{}

	for name := range pdeps {
		rmod, ok := t.env[name]
		if !ok {
			return nil, fmt.Errorf("unresolved dep: '%s'", name)
		}

		rin, rout := rmod.Ports()
		if err := t.compareAllPorts(
			pdeps[name].In, pdeps[name].Out, rin, rout,
		); err != nil {
			return nil, fmt.Errorf("incompatible ports on module '%s': %w", name, err)
		}

		rdeps[name] = rmod
	}

	return rdeps, nil
}

func (t translator) translateAllPorts(in parsing.InPorts, out parsing.OutPorts) (runtime.InPorts, runtime.OutPorts) {
	inPorts := t.translatePorts(parsing.Ports(in))
	outPorts := t.translatePorts(parsing.Ports(out))
	return runtime.InPorts(inPorts), runtime.OutPorts(outPorts)
}

func (t translator) translateWorkers(deps map[string]runtime.AbstractModule, wm map[string]string) runtime.Workers {
	rwm := runtime.Workers{}
	for workerName, depName := range wm {
		depMod, _ := deps[depName]
		rwm[workerName] = depMod
	}
	return rwm
}

func (t translator) translateNet(
	pnet parsing.Net,
	rin runtime.InPorts,
	rout runtime.OutPorts,
	rworkers runtime.Workers,
) []runtime.Conn {
	cc := make([]runtime.Conn, len(pnet))
	for i := range pnet {
		cc[i] = runtime.Conn{
			Sender:    make(<-chan runtime.Msg),
			Receivers: []chan<- runtime.Msg{},
		}
	}
	return cc
}

func (t translator) translatePorts(pports parsing.Ports) runtime.Ports {
	rports := runtime.Ports{}
	for name, typ := range pports {
		rports[name] = types.ByName(typ)
	}
	return rports
}

func (t translator) compareAllPorts(
	pin parsing.InPorts,
	pout parsing.OutPorts,
	rin runtime.InPorts,
	rout runtime.OutPorts,
) error {
	if err := t.comparePorts(
		parsing.Ports(pin),
		runtime.Ports(rin),
	); err != nil {
		return fmt.Errorf("incompatible inPorts: %w", err)
	}

	if err := t.comparePorts(
		parsing.Ports(pout),
		runtime.Ports(rout),
	); err != nil {
		return fmt.Errorf("incompatible outPorts: %w", err)
	}

	return nil
}

func (t translator) comparePorts(pports parsing.Ports, rports runtime.Ports) error {
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

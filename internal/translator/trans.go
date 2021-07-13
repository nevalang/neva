package translator

import (
	"fmt"

	"github.com/emil14/refactored-garbanzo/internal/parsing"
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/types"
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
	rdeps, err := t.resolveEnv(pmod.Deps)
	if err != nil {
		return runtime.ComplexModule{}, fmt.Errorf("unresolved deps: %w", err)
	}

	rin, rout := t.translateAllPorts(pmod.In, pmod.Out)
	rworkers := t.translateWorkers(rdeps, pmod.Workers)
	rnet := t.translateNet(pmod.Net, rin, rout, rworkers)

	return runtime.NewComplexModule(
		rin,
		rout,
		rnet,
	), nil
}

// resolveEnv returns error if at least one dependency could not be resolved.
func (t translator) resolveEnv(pdeps parsing.Deps) (runtime.Env, error) {
	rdeps := runtime.Env{}

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

func (t translator) translateWorkers(deps runtime.Env, pworkers map[string]string) runtime.Env {
	rworkers := make(runtime.Env, len(pworkers))
	for worker, dep := range pworkers {
		rworkers[worker] = deps[dep]
	}
	return rworkers
}

// TODO: rename
func (t translator) translateNet(
	pnet parsing.Net,
	rin runtime.InPorts,
	rout runtime.OutPorts,
	rworkers runtime.Env,
) []runtime.ChanRel {
	nio := make(map[string]nodeIO, len(rworkers)+2)

	nio["in"] = nodeIO{
		out: make(map[string]chan runtime.Msg),
	}
	nio["out"] = nodeIO{
		in: make(map[string]chan runtime.Msg),
	}

	for name, mod := range rworkers {
		nio[name] = nodeIO{
			in:  make(map[string]chan runtime.Msg),
			out: make(map[string]chan runtime.Msg),
		}

		in, out := mod.Ports()
		for portName := range in {
			nio[name].in[portName] = make(chan runtime.Msg)
		}
		for portName := range out {
			nio[name].out[portName] = make(chan runtime.Msg)
		}
	}

	cc := make([]runtime.ChanRel, len(pnet))

	for i, sub := range pnet {
		sender := nio[sub.Sender.Node].out[sub.Sender.Port]

		recievers := make([]chan runtime.Msg, len(sub.Recievers))
		for i, receiver := range sub.Recievers {
			recievers[i] = nio[receiver.Node].in[receiver.Port]
		}

		cc[i] = runtime.ChanRel{
			Sender:    sender,
			Receivers: recievers,
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

type nodeIO struct {
	in  map[string]chan runtime.Msg
	out map[string]chan runtime.Msg
}

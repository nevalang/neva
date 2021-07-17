package generator

import (
	"fmt"

	"github.com/emil14/refactored-garbanzo/internal/parser"
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

type (
	Generator interface {
		Generate(mod parser.CustomModule, env runtime.Env) (runtime.Module, error)
	}

	generator struct {
		validator Validator
	}
)

func New(v Validator) Generator {
	return generator{v}
}

func (t generator) Generate(pmod parser.CustomModule, env runtime.Env) (runtime.Module, error) {
	if err := t.validator.ValidateDeps(pmod.Deps, env); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	envDeps := t.depsToEnv(pmod.Deps, env)
	// rin, rout := t.translateAllPorts(pmod.In, pmod.Out)
	// // rworkers := t.translateWorkers(rdeps, pmod.Workers)
	// // rnet := t.translateNet(pmod.Net, rin, rout, rworkers)

	// return runtime.NewCustomModule(
	// 	rin,
	// 	rout,
	// )

	return nil, nil
}

func (t generator) depsToEnv(pdeps parser.Deps, env runtime.Env) runtime.Env {
	rdeps := make(runtime.Env, len(pdeps))
	for name := range pdeps {
		rdeps[name] = env[name]
	}
	return rdeps
}

func (t generator) translateAllPorts(
	in parser.InportsInterface,
	out parser.OutportsInterface,
) (runtime.InportsInterface, runtime.OutportsInterface) {
	inPorts := t.translatePorts(parser.PortsInterface(in))
	outPorts := t.translatePorts(parser.PortsInterface(out))
	return runtime.InportsInterface(inPorts), runtime.OutportsInterface(outPorts)
}

func (t generator) translateWorkers(deps runtime.Env, pworkers map[string]string) runtime.Env {
	rworkers := make(runtime.Env, len(pworkers))
	for worker, dep := range pworkers {
		rworkers[worker] = deps[dep]
	}
	return rworkers
}

func (t generator) translateNet(
	pnet parser.Net,
	rin runtime.InportsInterface,
	rout runtime.OutportsInterface,
	rworkers runtime.Env,
) []runtime.ChanRelation {
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

		in, out := mod.Interface()
		for portName := range in {
			nio[name].in[portName] = make(chan runtime.Msg)
		}
		for portName := range out {
			nio[name].out[portName] = make(chan runtime.Msg)
		}
	}

	cc := make([]runtime.ChanRelation, len(pnet))

	for i, sub := range pnet {
		sender := nio[sub.Sender.Node].out[sub.Sender.Port]

		recievers := make([]chan runtime.Msg, len(sub.Recievers))
		for i, receiver := range sub.Recievers {
			recievers[i] = nio[receiver.Node].in[receiver.Port]
		}

		cc[i] = runtime.ChanRelation{
			Sender:    sender,
			Receivers: recievers,
		}
	}

	return cc
}

func (t generator) translatePorts(pports parser.PortsInterface) runtime.PortsInterface {
	rports := runtime.PortsInterface{}
	for name, typ := range pports {
		rports[name] = types.ByName(typ)
	}
	return rports
}

type nodeIO struct {
	in  map[string]chan runtime.Msg
	out map[string]chan runtime.Msg
}

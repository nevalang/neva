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
		env runtime.Workers
	}
)

func New(env runtime.Workers) Translator {
	return translator{env}
}

func (t translator) Translate(pmod parsing.Module) (runtime.ComplexModule, error) {
	rin, rout, err := t.translateAllPorts(pmod.In, pmod.Out)
	if err != nil {
		return runtime.ComplexModule{}, fmt.Errorf("invalid ports: %w", err)
	}

	rdeps, err := t.translateDeps(pmod.Deps)
	if err != nil {
		return runtime.ComplexModule{}, fmt.Errorf("unresolved deps: %w", err)
	}

	rworkers, err := t.translateWorkers(rdeps, pmod.WorkerMap)
	if err != nil {
		return runtime.ComplexModule{}, fmt.Errorf("could not translate workers: %w", err)
	}

	rnet, err := t.translateNet(pmod.Net, rin, rout, rworkers)
	if err != nil {
		return runtime.ComplexModule{}, fmt.Errorf("could not translate net: %w", err)
	}

	return runtime.NewModule(rin, rout, rworkers, rnet), nil
}

func (t translator) translateWorkers(deps map[string]runtime.AbstractModule, wm map[string]string) (runtime.Workers, error) {
	rwm := runtime.Workers{}

	for workerName, depName := range wm {
		depMod, ok := deps[depName]
		if !ok {
			return nil, fmt.Errorf("dep '%s' not found for worker '%s'", depName, workerName)
		}
		rwm[workerName] = depMod
	}

	return rwm, nil
}

func (t translator) translateNet(
	pnet parsing.Net,
	rin runtime.InPorts,
	rout runtime.OutPorts,
	rworkers runtime.Workers,
) ([]runtime.Conn, error) {
	cc := make([]runtime.Conn, len(pnet))

	for i, pconn := range pnet {
		sender := make(<-chan runtime.Msg)
		receivers := []chan<- runtime.Msg{}

		for i := range pconn.Recievers {
			if err := t.checkConn(
				pconn.Sender, pconn.Recievers[i], rin, rout, rworkers,
			); err != nil {
				return nil, fmt.Errorf("mismatched port types on connection: %w", err)
			}
		}

		cc[i] = runtime.Conn{
			Sender:    sender,
			Receivers: receivers,
		}
	}

	return cc, nil
}

func (t translator) checkConn(
	psender parsing.Conn,
	preceiver parsing.Conn,
	rin runtime.InPorts, rout runtime.OutPorts, rworkers runtime.Workers,
) error {
	senderPortType := getPortType(psender.Node, psender.Port, rin, rout, rworkers, true)
	if senderPortType == types.Unknown {
		return fmt.Errorf("unknown sender port type")
	}

	receiverPortType := getPortType(preceiver.Node, preceiver.Port, rin, rout, rworkers, false)
	if receiverPortType == types.Unknown {
		return fmt.Errorf("unknown receiver port type")
	}

	if senderPortType != receiverPortType {
		return fmt.Errorf("mismatched port types: %s and %s", senderPortType, receiverPortType)
	}

	return nil
}

func getPortType(
	node string, port string,
	rin runtime.InPorts, rout runtime.OutPorts, rworkers runtime.Workers,
	isSender bool,
) types.Type {
	switch node {
	case "in":
		typ, ok := rin[node]
		if !ok {
			return types.Unknown
		}
		return typ
	case "out":
		typ, ok := rout[node]
		if !ok {
			return types.Unknown
		}
		return typ
	}

	in, out := rworkers[node].Ports()
	if isSender {
		typ, ok := out[node]
		if !ok {
			return types.Unknown
		}
		return typ
	}

	typ, ok := in[node]
	if !ok {
		return types.Unknown
	}

	return typ
}

func (t translator) translateDeps(pdeps parsing.Deps) (map[string]runtime.AbstractModule, error) {
	rdeps := map[string]runtime.AbstractModule{}

	for name := range pdeps {
		rmod, ok := t.env[name]
		if !ok {
			return nil, fmt.Errorf("unresolved dep: '%s'", name)
		}

		rin, rout := rmod.Ports()
		if err := compareAllPorts(
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
	inPorts, err := t.translatePorts(parsing.Ports(in))
	if err != nil {
		return runtime.InPorts{}, runtime.OutPorts{}, fmt.Errorf("could not translate inPorts: %w", err)
	}

	outPorts, err := t.translatePorts(parsing.Ports(out))
	if err != nil {
		return runtime.InPorts{}, runtime.OutPorts{}, fmt.Errorf("could not translate outPorts: %w", err)
	}

	return runtime.InPorts(inPorts), runtime.OutPorts(outPorts), nil
}

func (t translator) translatePorts(pports parsing.Ports) (runtime.Ports, error) {
	rports := runtime.Ports{}

	for name, typ := range pports {
		t := types.ByName(typ)
		if t == types.Unknown {
			return runtime.Ports{}, fmt.Errorf("unknown type %s", typ)
		}
		rports[name] = t
	}

	return rports, nil
}

func compareAllPorts(
	pin parsing.InPorts,
	pout parsing.OutPorts,
	rin runtime.InPorts,
	rout runtime.OutPorts,
) error {
	if err := comparePorts(
		parsing.Ports(pin),
		runtime.Ports(rin),
	); err != nil {
		return fmt.Errorf("incompatible inPorts: %w", err)
	}

	if err := comparePorts(
		parsing.Ports(pout),
		runtime.Ports(rout),
	); err != nil {
		return fmt.Errorf("incompatible outPorts: %w", err)
	}

	return nil
}

func comparePorts(pports parsing.Ports, rports runtime.Ports) error {
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

func compareTwoTypes(got string, want types.Type) error {
	if t := types.ByName(got); t != want {
		return fmt.Errorf("mismatched types: want %s got %s", want, t)
	}
	return nil
}

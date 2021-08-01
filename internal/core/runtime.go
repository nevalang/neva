package core

import (
	"fmt"
)

type Runtime struct {
	env   map[string]Module // TODO move out
	cache map[string]bool
}

const (
	tmpBuf     = 0
	tmpArrSize = 10
)

func (r Runtime) Run(name string) (NodeIO, error) {
	mod, ok := r.env[name]
	if !ok {
		return NodeIO{}, errModNotFound(name)
	}

	modInterface := mod.Interface()

	if nmod, ok := mod.(nativeModule); ok {
		io := r.nodeIO(modInterface.In, modInterface.Out)
		go nmod.impl(io)

		return io, nil
	}

	cmod, ok := mod.(customModule)
	if !ok {
		return NodeIO{}, errUnknownModType(name, mod)
	}

	if !r.cache[name] {
		if err := r.resolveDeps(cmod.deps); err != nil {
			return NodeIO{}, err
		}

		r.cache[name] = true
	}

	nodesIO := make(map[string]NodeIO, 2+len(cmod.workers))

	nodesIO["in"] = r.nodeIO(
		nil,
		OutportsInterface(modInterface.In),
	)
	nodesIO["out"] = r.nodeIO(
		InportsInterface(modInterface.Out),
		nil,
	)

	for w, dep := range cmod.workers {
		io, err := r.Run(dep)
		if err != nil {
			return NodeIO{}, err
		}

		nodesIO[w] = io
	}

	cc, err := r.connections(nodesIO, cmod.net)
	if err != nil {
		return NodeIO{}, err
	}

	r.connectAll(cc)

	return NodeIO{
		In:  nodeInports(nodesIO["in"].Out),
		Out: nodeOutports(nodesIO["out"].In),
	}, nil
}

func (rt Runtime) connections(io map[string]NodeIO, net []StreamDef) ([]connection, error) {
	rels := make([]connection, len(net))

	for i, rel := range net {
		sender, err := rt.chanByPoint(rel.Sender, nodePorts(io[rel.Sender.Node()].Out))
		if err != nil {
			return nil, fmt.Errorf("invalid sender, %w", err)
		}

		receivers := make([]chan Msg, len(rel.Receivers))
		for i, r := range rel.Receivers {
			receiver, err := rt.chanByPoint(r, nodePorts(io[r.Node()].In))
			if err != nil {
				return nil, fmt.Errorf("invalid receiver, %w", err)
			}

			receivers[i] = receiver
		}

		rels[i] = connection{
			Sender:    sender,
			Receivers: receivers,
		}
	}

	return rels, nil
}

func (r Runtime) chanByPoint(point PortPoint, ports nodePorts) (chan Msg, error) {
	arrpoint, ok := point.(ArrPortPoint)
	if ok {
		arrport, err := ports.arrPort(arrpoint.port)
		if err != nil {
			return nil, err
		}

		if uint8(len(arrport)) < arrpoint.idx {
			return nil, fmt.Errorf("arrport to small")
		}

		return arrport[arrpoint.idx], nil
	}

	normPoint, ok := point.(NormPortPoint)
	if !ok {
		return nil, fmt.Errorf("port point of unknown type %T", point)
	}

	normPort, err := ports.normPort(normPoint.port)
	if err != nil {
		return nil, err
	}

	return normPort, nil
}

func (r Runtime) resolveDeps(deps Interfaces) error {
	for dep := range deps {
		mod, ok := r.env[dep]
		if !ok {
			return errModNotFound(dep)
		}

		err := mod.Interface().Compare(deps[dep])
		if err != nil {
			return fmt.Errorf("unresolved dependency '%s': %w", dep, err)
		}
	}

	return nil
}

func (r Runtime) nodeIO(in InportsInterface, out OutportsInterface) NodeIO {
	inports := r.Ports(PortsInterface(in))
	outports := r.Ports(PortsInterface(out))

	return NodeIO{
		nodeInports(inports),
		nodeOutports(outports),
	}
}

func (r Runtime) Ports(ports PortsInterface) nodePorts {
	result := make(nodePorts, len(ports))

	for port, typ := range ports {
		if typ.Arr {
			cc := make([]chan Msg, tmpArrSize)
			for i := range cc {
				cc[i] = make(chan Msg)
			}
			result[port] = cc
			continue
		}

		result[port] = make(chan Msg)
	}

	return result
}

func (r Runtime) connectAll(cc []connection) {
	for i := range cc {
		go r.connect(cc[i])
	}
}

func (m Runtime) connect(c connection) {
	for msg := range c.Sender {
		for _, r := range c.Receivers {
			select {
			case r <- msg:
				continue
			default:
				go func() { r <- msg }()
			}
		}
	}
}

type Port interface{}

func New(env map[string]Module) Runtime {
	return Runtime{
		env:   env,
		cache: map[string]bool{},
	}
}

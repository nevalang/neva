package core

import (
	"fmt"
)

type Runtime struct {
	env   map[string]Component // TODO move out
	cache map[string]bool
}

const (
	tmpBuf     = 0
	tmpArrSize = 3 // FIXME deadlock possibilities
)

func (r Runtime) Run(name string) (NodeIO, error) {
	mod, ok := r.env[name]
	if !ok {
		return NodeIO{}, errModNotFound(name)
	}

	modInterface := mod.Interface()

	if nmod, ok := mod.(operator); ok {
		io := r.nodeIO(modInterface.In, modInterface.Out)
		if err := nmod.impl(io); err != nil {
			return NodeIO{}, err
		}

		return io, nil
	}

	cmod, ok := mod.(module)
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

	ss, err := r.streams(nodesIO, cmod.net)
	if err != nil {
		return NodeIO{}, err
	}

	r.startStreams(ss)

	return NodeIO{
		In:  nodeInports(nodesIO["in"].Out),
		Out: nodeOutports(nodesIO["out"].In),
	}, nil
}

func (rt Runtime) streams(io map[string]NodeIO, net Net) ([]stream, error) {
	ss := make([]stream, 0, len(net))

	for senderPoint, receiversPoints := range net {
		senderPort, err := rt.chanByPoint(senderPoint, nodePorts(io[senderPoint.Node()].Out))
		if err != nil {
			return nil, fmt.Errorf("invalid sender, %w", err)
		}

		receivers := make([]chan Msg, 0, len(receiversPoints))

		for receiverPoint := range receiversPoints {
			receiver, err := rt.chanByPoint(receiverPoint, nodePorts(io[receiverPoint.Node()].In))
			if err != nil {
				return nil, fmt.Errorf("invalid receiver, %w", err)
			}

			receivers = append(receivers, receiver)
		}

		ss = append(ss, stream{
			Sender:    senderPort,
			Receivers: receivers,
		})
	}

	return ss, nil
}

func (r Runtime) chanByPoint(point PortPoint, ports nodePorts) (chan Msg, error) {
	arrpoint, ok := point.(ArrPortPoint)
	if ok {
		arrport, err := ports.arr(arrpoint.port)
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

	normPort, err := ports.norm(normPoint.port)
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
			cc := make([]chan Msg, tmpArrSize) // FIXME
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

func (r Runtime) startStreams(ss []stream) {
	for i := range ss {
		go r.startStream(ss[i])
	}
}

func (m Runtime) startStream(s stream) {
	for msg := range s.Sender {
		for _, r := range s.Receivers {
			select {
			case r <- msg:
				continue
				// default:
				// 	go func() { r <- msg }()
			}
		}
	}
}

type Port interface{}

func New(env map[string]Component) Runtime {
	return Runtime{
		env:   env,
		cache: map[string]bool{},
	}
}

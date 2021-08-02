package core

import (
	"fmt"
)

type Runtime struct {
	env   map[string]Component
	cache map[string]bool
}

type Meta struct {
	arrPortsSize arrPortsSize
}

type arrPortsSize struct {
	In, Out map[string]uint8
}

func (r Runtime) Start(name string, meta Meta) (NodeIO, error) {
	c, ok := r.env[name]
	if !ok {
		return NodeIO{}, errModNotFound(name)
	}

	componentIO := c.Interface()

	if op, ok := c.(operator); ok {
		nodeIO := r.nodeIO(componentIO.In, componentIO.Out, meta.arrPortsSize)
		if err := op.impl(nodeIO); err != nil {
			return NodeIO{}, err
		}

		return nodeIO, nil
	}

	mod, ok := c.(module)
	if !ok {
		return NodeIO{}, errUnknownModType(name, c)
	}

	if !r.cache[name] {
		if err := r.resolveDeps(mod.deps); err != nil {
			return NodeIO{}, err
		}

		r.cache[name] = true
	}

	nodesIO := make(map[string]NodeIO, 2+len(mod.workers))

	nodesIO["in"] = r.nodeIO(
		nil,
		OutportsInterface(componentIO.In),
		meta.arrPortsSize,
	)
	nodesIO["out"] = r.nodeIO(
		InportsInterface(componentIO.Out),
		nil,
		meta.arrPortsSize,
	)

	for worker, dep := range mod.workers {
		meta := r.Meta(r.env[dep].Interface(), mod.net, worker) // TODO rewrite

		nodeIO, err := r.Start(dep, meta)
		if err != nil {
			return NodeIO{}, err
		}

		nodesIO[worker] = nodeIO
	}

	ss, err := r.streams(nodesIO, mod.net)
	if err != nil {
		return NodeIO{}, err
	}

	r.startStreams(ss)

	return NodeIO{
		in:  nodeInports(nodesIO["in"].out),
		out: nodeOutports(nodesIO["out"].in),
	}, nil
}

func (rt Runtime) Meta(io Interface, net Net, node string) Meta {
	m := arrPortsSize{
		In:  map[string]uint8{},
		Out: map[string]uint8{},
	}

	for port := range PortsInterface(io.In).Arr() {
		m.In[port] = net.ArrInSize(node, port)
	}

	for port := range PortsInterface(io.Out).Arr() {
		m.Out[port] = net.ArrOutSize(node, port)
	}

	return Meta{arrPortsSize: m}
}

func (rt Runtime) streams(io map[string]NodeIO, net Net) ([]stream, error) {
	ss := make([]stream, 0, len(net))

	for senderPoint, receiversPoints := range net {
		senderPort, err := rt.chanByPoint(senderPoint, nodePorts(io[senderPoint.Node()].out))
		if err != nil {
			return nil, fmt.Errorf("invalid sender, %w", err)
		}

		receivers := make([]chan Msg, 0, len(receiversPoints))

		for receiverPoint := range receiversPoints {
			receiver, err := rt.chanByPoint(receiverPoint, nodePorts(io[receiverPoint.Node()].in))
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

func (r Runtime) nodeIO(in InportsInterface, out OutportsInterface, size arrPortsSize) NodeIO {
	inports := r.ports(PortsInterface(in), size.In)
	outports := r.ports(PortsInterface(out), size.Out)

	return NodeIO{
		nodeInports(inports),
		nodeOutports(outports),
	}
}

func (r Runtime) ports(ports PortsInterface, size map[string]uint8) nodePorts {
	result := make(nodePorts, len(ports))

	for port, typ := range ports {
		if !typ.Arr {
			result[port] = make(chan Msg)

			continue
		}

		s, ok := size[port]
		if !ok {
			panic("no size for port " + port)
		}

		cc := make([]chan Msg, s)
		for i := range cc {
			cc[i] = make(chan Msg)
		}

		result[port] = cc
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

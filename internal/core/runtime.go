package core

import (
	"fmt"
)

type Runtime struct {
	env   map[string]Module
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

	if nmod, ok := mod.(NativeModule); ok {
		io := r.nodeIO(modInterface.In, modInterface.Out)
		go nmod.connect(io)
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

	net, err := r.net(nodesIO, cmod.net)
	if err != nil {
		return NodeIO{}, err
	}

	r.connectAll(net)

	return NodeIO{
		in:  nodeInports(nodesIO["in"].out),
		out: nodeOutports(nodesIO["out"].in),
	}, nil
}

func (r Runtime) net(io map[string]NodeIO, net []RelationsDef) ([]relations, error) {
	rels := make([]relations, len(net))

	for i, rel := range net {
		sender := r.chanByPoint(rel.Sender, io[rel.Sender.NodeName()])

		receivers := make([]chan Msg, len(rel.Recievers))
		for i, receiver := range rel.Recievers {
			receivers[i] = r.chanByPoint(receiver, io[receiver.NodeName()])
		}

		rels[i] = relations{
			Sender:    sender,
			Receivers: receivers,
		}
	}

	return rels, nil
}

func (r Runtime) chanByPoint(p PortPoint, io NodeIO) chan Msg {
	var result chan Msg

	arrprot, ok := p.(ArrPortPoint)
	if ok {
		arrport, err := io.ArrOutport(arrprot.Port)
		if err != nil {
			panic(err)
		}

		if uint8(len(arrport)) < arrprot.Index {
			panic("arrport to small")
		}

		result = arrport[arrprot.Index]
	} else {
		normport, err := io.NormOutport(arrprot.Port)
		if err != nil {
			panic(err)
		}

		result = normport
	}

	return result
}

func (r Runtime) resolveDeps(deps Interfaces) error {
	for dep := range deps {
		mod, ok := r.env[dep]
		if !ok {
			return errModNotFound(dep)
		}

		err := mod.Interface().Compare(deps[dep])
		if err != nil {
			return fmt.Errorf("ports incompatible on module '%s': %w", dep, err)
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

func (r Runtime) connectAll(rels []relations) {
	for i := range rels {
		go r.connect(rels[i])
	}
}

func (m Runtime) connect(c relations) {
	for msg := range c.Sender {
		for i := range c.Receivers {
			r := c.Receivers[i]
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

// func checkAllPorts(got, want Interface) error {
// 	if err := checkPorts(
// 		PortsInterface(got.In),
// 		PortsInterface(want.In),
// 	); err != nil {
// 		return fmt.Errorf("incompatible inPorts: %w", err)
// 	}

// 	if err := checkPorts(
// 		PortsInterface(got.Out),
// 		PortsInterface(want.Out),
// 	); err != nil {
// 		return fmt.Errorf("incompatible inPorts: %w", err)
// 	}

// 	return nil
// }

// func checkPorts(got, want PortsInterface) error {
// 	if len(got) < len(want) {
// 		return fmt.Errorf(
// 			"not enough ports: got %d, want %d",
// 			len(got),
// 			len(want),
// 		)
// 	}

// 	for name := range want {
// 		if want[name] != got[name] {
// 			return fmt.Errorf(
// 				"incompatible types on port '%s': got '%v', want '%v'",
// 				name,
// 				want[name],
// 				got[name],
// 			)
// 		}
// 	}

// 	return nil
// }

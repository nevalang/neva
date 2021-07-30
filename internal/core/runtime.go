package core

import (
	"errors"
	"fmt"
)

type Runtime struct {
	env Env
	buf uint8
}

type Env map[string]Module

const tmpBuf = 0

// TODO test only native module
func (r Runtime) Run(root string) (NodeIO, error) {
	mod, ok := r.env[root]
	if !ok {
		return NodeIO{}, errModNotFound(root)
	}

	modInterface := mod.Interface()

	if nmod, ok := mod.(NativeModule); ok {
		io := r.nodeIO(modInterface.In, modInterface.Out)
		go nmod.connect(io)
		return io, nil
	}

	cmod, ok := mod.(customModule)
	if !ok {
		return NodeIO{}, errUnknownModType(root, mod)
	}

	if err := r.checkDeps(cmod.deps); err != nil {
		return NodeIO{}, err
	}

	nodesIO := make(NodesIO, 2+len(cmod.workers))

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

	net := r.net(nodesIO, cmod.net)
	r.connectAll(net)

	return NodeIO{
		in:  nodeInports(nodesIO["in"].out),
		out: nodeOutports(nodesIO["out"].in),
	}, nil
}

func (r Runtime) net(io NodesIO, net Net) []relations {
	rels := []relations{}

	// for _, s := range net {
	// 	receivers := []chan Msg{}

	// 	for _, receiver := range s.Recievers {

	// 		aport, err := io[receiver.Node].ArrInport(receiver.Port)
	// 		if err == nil {
	// 			for _, p := range aport {
	// 				receivers = append(receivers, p)
	// 			}
	// 			continue
	// 		}
	// 		nport, _ := io[receiver.Node].Inport(receiver.Port)
	// 		receivers = append(receivers, nport)
	// 	}

	// 	sender := r.Sender(io, s.Sender.Node, s.Sender.Port)

	// 	rels = append(rels, relations{
	// 		Sender:    io[s.Sender.Node].out[s.Sender.Port],
	// 		Receivers: receivers,
	// 	})
	// }

	return rels
}

func (r Runtime) Sender(io NodesIO, node string, port string) chan Msg {
	// port := io[port].out[node]
	// if isArrPort(port) {
	// }
	return nil // TODO
}

// func isArrPort(port string) bool {

// }

// checkDeps checks that scope contains all the required modules.
func (r Runtime) checkDeps(deps Interfaces) error {
	for dep := range deps {
		mod, ok := r.env[dep]
		if !ok {
			return errModNotFound(dep)
		}

		if err := mod.Interface().Compare(deps[dep]); err != nil {
			return fmt.Errorf("ports incompatible on module '%s': %w", dep, err)
		}
	}

	return nil
}

func (r Runtime) nodeIO(in InportsInterface, out OutportsInterface) NodeIO {
	inports := make(nodeInports, len(in))
	outports := make(nodeOutports, len(in))

	for port, typ := range in {
		if typ.Arr {
			cc := make([]chan Msg, typ.Size)
			for i := range cc {
				cc[i] = make(chan Msg)
			}
			inports[port] = cc
			continue
		}

		inports[port] = make(chan Msg)
	}

	for port, typ := range out {
		if typ.Arr {
			cc := make([]chan Msg, typ.Size)
			for i := range cc {
				cc[i] = make(chan Msg)
			}
			outports[port] = cc
			continue
		}

		outports[port] = make(chan Msg)
	}

	return NodeIO{inports, outports}
}

func (r Runtime) Ports(ports PortsInterface) nodePorts {
	result := nodePorts{}

	for port, typ := range ports {
		if typ.Arr {
			cc := make([]chan Msg, typ.Size)
			for i := range cc {
				cc[i] = make(chan Msg)
			}
			result[port] = cc
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

func checkAllPorts(got, want Interface) error {
	if err := checkPorts(
		PortsInterface(got.In),
		PortsInterface(want.In),
	); err != nil {
		return fmt.Errorf("incompatible inPorts: %w", err)
	}

	if err := checkPorts(
		PortsInterface(got.Out),
		PortsInterface(want.Out),
	); err != nil {
		return fmt.Errorf("incompatible inPorts: %w", err)
	}

	return nil
}

func checkPorts(got, want PortsInterface) error {
	if len(got) < len(want) {
		return fmt.Errorf(
			"not enough ports: got %d, want %d",
			len(got),
			len(want),
		)
	}

	for name := range want {
		if want[name] != got[name] {
			return fmt.Errorf(
				"incompatible types on port '%s': got '%v', want '%v'",
				name,
				want[name],
				got[name],
			)
		}
	}

	return nil
}

func (io NodeIO) Port(name string) (NormalPort, error) {
	p, ok := io.in[name]
	if !ok {
		return nil, errors.New("...")
	}

	c, ok := p.(NormalPort)
	if !ok {
		return nil, errors.New("...")
	}

	return c, nil
}

func (io NodeIO) ArrPort(name string) (ArrPort, error) {
	p, ok := io.in[name]
	if !ok {
		return nil, errors.New("...")
	}

	cc, ok := p.(ArrPort)
	if !ok {
		return nil, errors.New("...")
	}

	return cc, nil
}

type NodeInports map[string]Port

type NodeOutports map[string]Port

type Port interface{}

type NormalPort chan Msg

type ArrPort []chan Msg

type Relations struct {
	Sender    chan Msg
	Receivers []chan Msg
}

func New(env Env) Runtime {
	return Runtime{env, tmpBuf}
}

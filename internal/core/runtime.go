package core

import "fmt"

type Runtime struct {
	scope Scope
	buf   uint8
}

type Scope map[string]Module

const tmpBuf = 0

func (r Runtime) Run(root string) (NodeIO, error) {
	mod, ok := r.scope[root]
	if !ok {
		return NodeIO{}, fmt.Errorf("mod not found")
	}

	if native, ok := mod.(NativeModule); ok {
		io := r.nodeIO(native.in, native.out)
		go native.impl(io)
		return io, nil
	}

	custom, ok := r.scope[root].(customModule)
	if !ok {
		return NodeIO{}, fmt.Errorf("mod unknown type")
	}

	if err := r.checkDeps(custom.deps); err != nil {
		return NodeIO{}, err
	}

	nodesIO := make(NodesIO, 2+len(custom.workers))

	nodesIO["in"] = r.nodeIO(
		nil,
		OutportsInterface(custom.in),
	)
	nodesIO["out"] = r.nodeIO(
		InportsInterface(custom.out),
		nil,
	)

	for w, dep := range custom.workers {
		io, err := r.Run(dep)
		if err != nil {
			return NodeIO{}, err
		}

		nodesIO[w] = io
	}

	net := r.net(custom.net, nodesIO)

	r.connectAll(net)

	return NodeIO{
		In:  NodeInports(nodesIO["in"].Out),
		Out: NodeOutports(nodesIO["out"].In),
	}, nil
}

func (r Runtime) net(net Net, nodesIO NodesIO) []Relations {
	rels := make([]Relations, len(net))
	for i, s := range net {
		receivers := make([]chan Msg, len(s.Recievers))
		for i, receiver := range s.Recievers {
			receivers[i] = nodesIO[receiver.Node].In[receiver.Port]
		}

		rels[i] = Relations{
			Sender:    nodesIO[s.Sender.Node].Out[s.Sender.Port],
			Receivers: receivers,
		}
	}
	return rels
}

// checkDeps checks that scope contains all the required modules.
func (r Runtime) checkDeps(deps Deps) error {
	for dep := range deps {
		mod, ok := r.scope[dep]
		if !ok {
			return fmt.Errorf("%w: '%s'", ErrModNotFound, dep)
		}

		if err := checkAllPorts(mod.Interface(), deps[dep]); err != nil {
			return fmt.Errorf("ports incompatibility on module '%s': %w", dep, err)
		}
	}

	return nil
}

func (r Runtime) nodeIO(in InportsInterface, out OutportsInterface) NodeIO {
	inports := make(map[string]chan Msg, len(in))
	outports := make(map[string]chan Msg, len(in))

	for port := range in {
		inports[port] = make(chan Msg)
	}
	for port := range out {
		outports[port] = make(chan Msg)
	}

	return NodeIO{inports, outports}
}

func (r Runtime) connectAll(rels []Relations) {
	for i := range rels {
		go r.connect(rels[i])
	}
}

func (m Runtime) connect(c Relations) {
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

type NodeIO struct {
	In  NodeInports
	Out NodeOutports
}

type NodeInports map[string]chan Msg

type NodeOutports map[string]chan Msg

type Msg struct {
	Str  string
	Int  int
	Bool bool
}

type Relations struct {
	Sender    chan Msg
	Receivers []chan Msg
}

type NodesIO map[string]NodeIO

func NewRuntime(scope Scope) Runtime {
	return Runtime{scope, tmpBuf}
}

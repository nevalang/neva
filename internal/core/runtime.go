package core

import "fmt"

type Runtime struct {
	scope Modules
	buf   uint8
}

type Modules map[string]Module

const tmpBuf = 0

func NewRuntime(scope Modules) Runtime {
	return Runtime{scope, tmpBuf}
}

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

	// workers + io
	nodesIO := make(map[string]NodeIO, len(custom.workers)+2)

	// io nodes
	nodesIO["in"] = r.nodeIO(
		nil,
		Outports(custom.in),
	)
	nodesIO["out"] = r.nodeIO(
		Inport(custom.out),
		nil,
	)

	// worker nodes
	for w, dep := range custom.workers {
		io, err := r.Run(dep)
		if err != nil {
			return NodeIO{}, err
		}
		nodesIO[w] = io
	}

	// net
	net := make([]Relations, len(custom.net))
	for i, s := range custom.net {
		receivers := make([]chan Msg, len(s.Recievers))
		for i, receiver := range s.Recievers {
			receivers[i] = nodesIO[receiver.Node].In[receiver.Port]
		}

		net[i] = Relations{
			Sender:    nodesIO[s.Sender.Node].Out[s.Sender.Port],
			Receivers: receivers,
		}
	}

	r.connectAll(net)

	return NodeIO{
		In:  NodeInports(nodesIO["in"].Out),
		Out: NodeOutports(nodesIO["out"].In),
	}, nil
}

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

func (r Runtime) nodeIO(in Inport, out Outports) NodeIO {
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
				"incompatible types on port '%s': got '%s', want '%s'",
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

func NewStrMsg(msg string) Msg {
	return Msg{Str: msg}
}

type Relations struct {
	Sender    chan Msg
	Receivers []chan Msg
}

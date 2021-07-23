package core

import "fmt"

type Runtime struct {
	env map[string]Module
}

func NewRuntime(env map[string]Module) Runtime {
	return Runtime{env}
}

func (r Runtime) Run(root string) (NodeIO, error) {
	mod := r.env[root].(CustomModule)
	if err := r.resolveDeps(mod.deps); err != nil {
		return NodeIO{}, err
	}

	nodesIO := make(map[string]NodeIO, len(mod.workers)+2) // workers + io

	// worker nodes
	for w, dep := range mod.workers {
		io, err := r.env[dep].SpawnWorker(r.env)
		if err != nil {
			return NodeIO{}, err
		}
		nodesIO[w] = io
	}

	// io nodes
	nodesIO["in"] = NodeIO{
		Out: make(map[string]chan Msg, len(mod.in)),
	}
	for port := range mod.in {
		nodesIO["in"].Out[port] = make(chan Msg, tmpBuf)
	}

	nodesIO["out"] = NodeIO{
		In: make(map[string]chan Msg),
	}
	for port := range mod.out {
		nodesIO["out"].In[port] = make(chan Msg, tmpBuf)
	}

	net := make([]Connection, len(mod.net))
	for i, s := range mod.net {
		receivers := make([]chan Msg, len(s.Recievers))
		for i, receiver := range s.Recievers {
			receivers[i] = nodesIO[receiver.Node].In[receiver.Port]
		}

		net[i] = Connection{
			Sender:    nodesIO[s.Sender.Node].Out[s.Sender.Port],
			Receivers: receivers,
		}
	}

	mod.connectAll(net)

	return NodeIO{
		In:  NodeInports(nodesIO["in"].Out),
		Out: NodeOutports(nodesIO["out"].In),
	}, nil
}

func (r Runtime) resolveDeps(deps Deps) error {
	for dep := range deps {
		mod, ok := r.env[dep]
		if !ok {
			return fmt.Errorf("%w: '%s'", ErrModNotFound, dep)
		}
		if err := checkAllPorts(mod.Interface(), deps[dep]); err != nil {
			return fmt.Errorf("ports incompatibility on module '%s': %w", dep, err)
		}
	}
	return nil
}

func (r Runtime) createIO(in InportsInterface, out OutportsInterface) NodeIO {
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

func (r Runtime) connectAll(rels []Connection) {
	for i := range rels {
		go r.connect(rels[i])
	}
}

func (m Runtime) connect(c Connection) {
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

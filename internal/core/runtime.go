package core

import (
	"fmt"
)

type Runtime struct {
	env Env
	buf uint8
}

type Env map[string]Module

const tmpBuf = 0

func (r Runtime) Run(root string) (NodeIO, error) {
	mod, ok := r.env[root]
	if !ok {
		return NodeIO{}, fmt.Errorf("mod not found")
	}

	if native, ok := mod.(NativeModule); ok {
		io := r.nodeIO(native.in, native.out)
		go native.impl(io)
		return io, nil
	}

	custom, ok := r.env[root].(customModule)
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
		In:  nodeInports(nodesIO["in"].Out),
		Out: nodeOutports(nodesIO["out"].In),
	}, nil
}

func (r Runtime) net(net Net, nodesIO NodesIO) []relations {
	rels := []relations{}

	for _, s := range net {
		receivers := []chan Msg{}

		for _, receiver := range s.Recievers {
			arrport, err := nodesIO[receiver.Node].ArrInport(receiver.Port)
			if err == nil {
				for _, p := range arrport {
					receivers = append(receivers, p)
				}
				continue
			}
			normport, _ := nodesIO[receiver.Node].Inport(receiver.Port)
			receivers = append(receivers, normport)
		}

		rels = append(rels, relations{
			Sender:    nodesIO[s.Sender.Node].out[s.Sender.Port],
			Receivers: receivers,
		})
	}

	return rels
}

// checkDeps checks that scope contains all the required modules.
func (r Runtime) checkDeps(deps Deps) error {
	for dep := range deps {
		mod, ok := r.env[dep]
		if !ok {
			return fmt.Errorf("%w: '%s'", ErrModNotFound, dep)
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

	for port := range in {
		inports[port] = make(chan Msg)
	}
	for port := range out {
		outports[port] = make(chan Msg)
	}

	return NodeIO{inports, outports}
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

func New(scope Env) Runtime {
	return Runtime{scope, tmpBuf}
}

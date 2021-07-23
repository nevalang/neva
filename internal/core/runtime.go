package core

import "fmt"

type Runtime struct {
	env map[string]Module
}

func NewRuntime(env map[string]Module) Runtime {
	return Runtime{env}
}

func (r Runtime) Run(root string) (NodeIO, error) {
	mod := r.env[root]

	if err := r.resolveDeps(mod.Deps()); err != nil {
		return NodeIO{}, err
	}

	io := r.createIO(mod.Interface().In, mod.Interface().Out)
	w := r.worker(mod)

	switch v := mod.(type) {
	case NativeModule:
		go v.impl(io)
		return io, nil
	case CustomModule:
		return NodeIO{}, nil
	}
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

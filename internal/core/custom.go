package core

import (
	"fmt"
)

type CustomModule struct {
	deps    Deps
	in      InportsInterface
	out     OutportsInterface
	workers Workers
	net     Net
}

func (cm CustomModule) Interface() Interface {
	return Interface{
		In:  cm.in,
		Out: cm.out,
	}
}

type Workers map[string]string

type Net []Subscription

type Subscription struct {
	Sender    PortPoint
	Recievers []PortPoint
}

type PortPoint struct {
	Node string
	Port string
}

const tmpBuf = 0

func (m CustomModule) SpawnWorker(scope map[string]Module) (NodeIO, error) {
	if err := m.resolveDeps(scope); err != nil {
		return NodeIO{}, err
	}

	nodesIO := make(map[string]NodeIO, len(m.workers)+2) // workers + io

	// worker nodes
	for w, dep := range m.workers {
		io, err := scope[dep].SpawnWorker(scope)
		if err != nil {
			return NodeIO{}, err
		}
		nodesIO[w] = io
	}

	// io nodes
	nodesIO["in"] = NodeIO{
		Out: make(map[string]chan Msg, len(m.in)),
	}
	for port := range m.in {
		nodesIO["in"].Out[port] = make(chan Msg, tmpBuf)
	}

	nodesIO["out"] = NodeIO{
		In: make(map[string]chan Msg),
	}
	for port := range m.out {
		nodesIO["out"].In[port] = make(chan Msg, tmpBuf)
	}

	net := make([]Connection, len(m.net))
	for i, s := range m.net {
		receivers := make([]chan Msg, len(s.Recievers))
		for i, receiver := range s.Recievers {
			receivers[i] = nodesIO[receiver.Node].In[receiver.Port]
		}

		net[i] = Connection{
			Sender:    nodesIO[s.Sender.Node].Out[s.Sender.Port],
			Receivers: receivers,
		}
	}

	m.connectAll(net)

	return NodeIO{
		In:  NodeInports(nodesIO["in"].Out),
		Out: NodeOutports(nodesIO["out"].In),
	}, nil
}

func (m CustomModule) connectAll(rels []Connection) {
	for i := range rels {
		go m.connect(rels[i])
	}
}

func (m CustomModule) connect(c Connection) {
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

func (m CustomModule) resolveDeps(env map[string]Module) error {
	for dep := range m.deps {
		mod, ok := env[dep]
		if !ok {
			return fmt.Errorf("%w: '%s'", ErrModNotFound, dep)
		}
		if err := checkAllPorts(mod.Interface(), m.deps[dep]); err != nil {
			return fmt.Errorf("ports incompatibility on module '%s': %w", dep, err)
		}
	}
	return nil
}

type Connection struct {
	Sender    chan Msg
	Receivers []chan Msg
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

func NewCustomModule(Deps, InportsInterface, OutportsInterface, Workers, Net) Module {
	return nil
}

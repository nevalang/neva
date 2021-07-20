package runtime

import (
	"fmt"
)

type CustomModule struct {
	Deps    Deps
	In      InportsInterface
	Out     OutportsInterface
	Workers Workers
	Net     Net
}

type Deps map[string]ModuleInterface

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

func (cm CustomModule) Interface() ModuleInterface {
	return ModuleInterface{
		In:  cm.In,
		Out: cm.Out,
	}
}

func (m CustomModule) SpawnWorker(env map[string]Module) (NodeIO, error) {
	if err := m.checkDeps(env); err != nil {
		return NodeIO{}, err
	}

	nodesIO := make(map[string]NodeIO, len(m.Workers)+2) // workers + io

	// create nodes for workers
	for w, dep := range m.Workers {
		io, err := env[dep].SpawnWorker(env)
		if err != nil {
			return NodeIO{}, err
		}
		nodesIO[w] = io
	}

	// create io nodes
	nodesIO["in"] = NodeIO{
		Out: make(map[string]chan Msg, len(m.In)),
	}
	for port := range m.In {
		nodesIO["in"].Out[port] = make(chan Msg)
	}
	nodesIO["out"] = NodeIO{
		In: make(map[string]chan Msg),
	}
	for port := range m.Out {
		nodesIO["out"].In[port] = make(chan Msg)
	}

	net := []ChanRel{}
	for _, s := range m.Net {
		receivers := make([]chan Msg, len(s.Recievers))
		for i, receiver := range s.Recievers {
			receivers[i] = nodesIO[receiver.Node].In[receiver.Port]
		}

		net = append(net, ChanRel{
			Sender:    nodesIO[s.Sender.Node].Out[s.Sender.Port],
			Receivers: receivers,
		})
	}

	m.connectAll(net)

	return NodeIO{
		In:  nodesIO["in"].Out,
		Out: nodesIO["out"].In,
	}, nil
}

func (m CustomModule) connectAll(rels []ChanRel) {
	for i := range rels {
		go m.connect(rels[i])
	}
}

func (m CustomModule) connect(c ChanRel) {
	for msg := range c.Sender {
		for i := range c.Receivers {
			r := c.Receivers[i]
			go func() { r <- msg }()
		}
	}
}

type NodeIO struct {
	In, Out map[string]chan Msg
}

func (m CustomModule) checkDeps(env map[string]Module) error {
	for dep := range m.Deps {
		mod, ok := env[dep]
		if !ok {
			return fmt.Errorf("%w: '%s'", ErrModNotFound, dep)
		}
		if err := compareAllPorts(mod.Interface(), m.Deps[dep]); err != nil {
			return fmt.Errorf("ports incompatibility on module '%s': %w", dep, err)
		}
	}
	return nil
}

type ChanRel struct {
	Sender    chan Msg
	Receivers []chan Msg
}

func compareAllPorts(got, want ModuleInterface) error {
	if err := comparePorts(
		PortsInterface(got.In),
		PortsInterface(want.In),
	); err != nil {
		return fmt.Errorf("incompatible inPorts: %w", err)
	}

	if err := comparePorts(
		PortsInterface(got.Out),
		PortsInterface(want.Out),
	); err != nil {
		return fmt.Errorf("incompatible inPorts: %w", err)
	}

	return nil
}

func comparePorts(got PortsInterface, want PortsInterface) error {
	if len(want) != len(got) {
		return fmt.Errorf(
			"different number of ports: got %d, want %d",
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

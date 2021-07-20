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

const tmpBuf = 0

func (m CustomModule) SpawnWorker(env map[string]Module) (NodeIO, error) {
	if err := m.checkDeps(env); err != nil {
		return NodeIO{}, err
	}

	nodesIO := make(map[string]NodeIO, len(m.Workers)+2) // workers + io

	// worker nodes
	for w, dep := range m.Workers {
		io, err := env[dep].SpawnWorker(env)
		if err != nil {
			return NodeIO{}, err
		}
		nodesIO[w] = io
	}

	// io nodes
	nodesIO["in"] = NodeIO{
		Out: make(map[string]chan Msg, len(m.In)),
	}
	for port := range m.In {
		nodesIO["in"].Out[port] = make(chan Msg, tmpBuf)
	}

	nodesIO["out"] = NodeIO{
		In: make(map[string]chan Msg),
	}
	for port := range m.Out {
		nodesIO["out"].In[port] = make(chan Msg, tmpBuf)
	}

	net := make([]Relation, len(m.Net))
	for i, s := range m.Net {
		receivers := make([]chan Msg, len(s.Recievers))
		for i, receiver := range s.Recievers {
			receivers[i] = nodesIO[receiver.Node].In[receiver.Port]
		}

		fmt.Println(
			"===\n",
			s.Sender,
			"-->",
			s.Recievers,
			"===\n",
		)

		net[i] = Relation{
			Sender:    nodesIO[s.Sender.Node].Out[s.Sender.Port],
			Receivers: receivers,
		}

		fmt.Println(net)
	}

	for i, s := range net {
		if s.Sender == nil {
			fmt.Println("betrayer", m.Net[i].Sender)
		}
	}

	m.connectAll(net)

	return NodeIO{
		In:  NodeInports(nodesIO["in"].Out),
		Out: NodeOutports(nodesIO["out"].In),
	}, nil
}

func (m CustomModule) connectAll(rels []Relation) {
	for i := range rels {
		go m.connect(rels[i])
	}
}

func (m CustomModule) connect(c Relation) {
	for msg := range c.Sender {
		for i := range c.Receivers {
			r := c.Receivers[i]
			go func() {
				r <- msg
			}()
		}
	}
}

func (m CustomModule) checkDeps(env map[string]Module) error {
	for dep := range m.Deps {
		mod, ok := env[dep]
		if !ok {
			return fmt.Errorf("%w: '%s'", ErrModNotFound, dep)
		}
		if err := checkAllPorts(mod.Interface(), m.Deps[dep]); err != nil {
			return fmt.Errorf("ports incompatibility on module '%s': %w", dep, err)
		}
	}
	return nil
}

type Relation struct {
	Sender    chan Msg
	Receivers []chan Msg
}

func checkAllPorts(got, want ModuleInterface) error {
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

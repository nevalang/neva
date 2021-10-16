package runtime

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime/program"
)

var ErrPortNotFound = errors.New("port not found")

type (
	// TODO: hide implementation behind interface in order to free operators from this code
	// or move operators related code to separate package
	Runtime struct {
		cnctr Connector
	}

	Connector interface {
		ConnectSubnet([]Connection)
		ConnectOperator(string, IO) error
	}
)

// Run creates root node of the program and returns its io.
func (r Runtime) Run(p program.Program) (IO, error) {
	return r.spawn(p.Scope, p.RootNodeMeta)
}

// spawn uses nodemeta to find component in scope and create a node and returns io of that node.
func (r Runtime) spawn(scope map[string]program.Component, nodeMeta program.NodeMeta) (IO, error) {
	component, ok := scope[nodeMeta.ComponentName]
	if !ok {
		return IO{}, fmt.Errorf("component not found: %s", nodeMeta.ComponentName)
	}

	io := r.nodeIO(nodeMeta)

	if component.Operator != "" {
		if err := r.cnctr.ConnectOperator(component.Operator, io); err != nil {
			return IO{}, fmt.Errorf("connect operator: %w", err)
		}
		return r.asSubNode(nodeMeta, io), nil
	}

	// it's a module so it has subnetwork and in-out nodes are part it
	subnetNodesIO := map[string]IO{
		"in":  {Out: io.In}, // for subnet 'in' node is sender
		"out": {In: io.Out}, // and 'out' is receiver
	}

	// repeat this algorithm for every worker to collect their io
	for workerNodeName, workerNodeMeta := range component.WorkerNodesMeta {
		workerNodeIO, err := r.spawn(scope, workerNodeMeta) // <- recursion
		if err != nil {
			return IO{}, err
		}
		subnetNodesIO[workerNodeName] = workerNodeIO
	}

	cc, err := r.connections(subnetNodesIO, component.Net)
	if err != nil {
		return IO{}, err
	}

	r.cnctr.ConnectSubnet(cc)

	return r.asSubNode(nodeMeta, io), nil
}

// parent network will use this io by worker name
func (r Runtime) asSubNode(meta program.NodeMeta, io IO) IO {
	io2 := IO{
		In:  map[program.PortAddr]chan Msg{},
		Out: map[program.PortAddr]chan Msg{},
	}
	for addr, ch := range io.In {
		addr.Node = meta.Name
		io2.In[addr] = ch
	}
	for addr, ch := range io.Out {
		addr.Node = meta.Name
		io2.Out[addr] = ch
	}
	return io2
}

// connections initializes channels for network.
func (r Runtime) connections(nodesIO map[string]IO, net []program.Connection) ([]Connection, error) {
	cc := make([]Connection, len(net))

	for i, c := range net {
		fromNodeIO, ok := nodesIO[c.From.Node]
		if !ok {
			return nil, fmt.Errorf("fromNodeIO, ok := nodesIO[c.From.Node]")
		}

		sender, ok := fromNodeIO.Out[c.From] // has in/out names
		if !ok {
			return nil, fmt.Errorf("from, ok := fromNodeIO.Out[fromAddr]")
		}

		receivers := make([]Port, len(c.To))
		for j, toAddr := range c.To {
			toNodeIO, ok := nodesIO[toAddr.Node]
			if !ok {
				return nil, fmt.Errorf("toNodeIO, ok := nodesIO[to.Node]")
			}

			receiver, ok := toNodeIO.In[toAddr]
			if !ok {
				return nil, fmt.Errorf("receiver, ok := toNodeIO.In[toAddr]")
			}

			receivers[j] = Port{Ch: receiver, Addr: toAddr}
		}

		cc[i] = Connection{
			From: Port{Ch: sender, Addr: c.From},
			To:   receivers,
		}
	}

	return cc, nil
}

// nodeIO creates channels for node.
func (r Runtime) nodeIO(nodeMeta program.NodeMeta) IO {
	in := make(map[program.PortAddr]chan Msg)

	for port, slots := range nodeMeta.In {
		addr := program.PortAddr{Port: port, Node: "in"}

		if slots == 0 {
			in[addr] = make(chan Msg)
			continue
		}

		for idx := uint8(0); idx < slots; idx++ {
			addr.Idx = idx
			in[addr] = make(chan Msg)
		}
	}

	outports := make(map[program.PortAddr]chan Msg)

	for port, slots := range nodeMeta.Out {
		addr := program.PortAddr{Port: port, Node: "out"}

		if slots == 0 {
			outports[addr] = make(chan Msg)
			continue
		}

		for idx := uint8(0); idx < slots; idx++ {
			addr.Idx = idx
			outports[addr] = make(chan Msg)
		}
	}

	return IO{in, outports}
}

// Operator returns error if IO doesn't fit.
type Operator func(IO) error

// Connection represents sender-receiver pair.
type Connection struct {
	From Port
	To   []Port
}

// Port maps network address with the real channel.
type Port struct {
	Ch   chan Msg
	Addr program.PortAddr
}

// IO represents node's input and output ports.
type IO struct {
	In, Out Ports
}

// Ports maps network addresses to real channels.
type Ports map[program.PortAddr]chan Msg

// PortGroup returns all port-chanells associated with the given array port name.
func (ports Ports) PortGroup(arrPort string) ([]chan Msg, error) {
	cc := []chan Msg{}

	for addr, ch := range ports {
		if addr.Port == arrPort {
			cc = append(cc, ch)
		}
	}

	if len(cc) == 0 {
		return nil, fmt.Errorf("ErrArrPortNotFound: %s", arrPort)
	}

	return cc, nil
}

func (ports Ports) Port(port string) (chan Msg, error) {
	for addr, ch := range ports {
		if addr.Port != port {
			continue
		}
		if addr.Idx > 0 {
			return nil, fmt.Errorf("unexpected port group %v", addr)
		}
		return ch, nil
	}
	return nil, fmt.Errorf("%w: looking for %s, have: %v", ErrPortNotFound, port, ports)
}

func New(connector Connector) Runtime {
	return Runtime{
		cnctr: connector,
	}
}

type AbsPortAddr struct {
	port Port
	path []string
}

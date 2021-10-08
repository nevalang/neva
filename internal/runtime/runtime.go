package runtime

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime/program"
)

var ErrPortNotFound = errors.New("port not found")

type (
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
	return r.run(p.Scope, p.RootNodeMeta)
}

// run uses nodemeta to find component in scope and create a node and returns io of that node.
func (r Runtime) run(scope map[string]program.Component, nodeMeta program.NodeMeta) (IO, error) {
	component, ok := scope[nodeMeta.Component]
	if !ok {
		return IO{}, fmt.Errorf("component not found: %s", nodeMeta.Component)
	}

	nodeIO := r.nodeIO(nodeMeta)

	if component.Operator != "" {
		if err := r.cnctr.ConnectOperator(component.Operator, nodeIO); err != nil {
			return IO{}, fmt.Errorf("connect operator: %w", err)
		}
		return nodeIO, nil
	}

	// it's a module so it has subnetwork and in-out nodes are part it
	subnetIO := map[string]IO{
		"in":  {Out: nodeIO.In}, // for subnet 'in' node is sender
		"out": {In: nodeIO.Out}, // and 'out' is receiver
	}

	// repeat this algorithm for every worker to collect their io
	for workerNodeName, workerNodeMeta := range component.WorkerNodesMeta {
		workerNodeIO, err := r.run(scope, workerNodeMeta) // <- recursion
		if err != nil {
			return IO{}, err
		}
		subnetIO[workerNodeName] = workerNodeIO
	}

	r.cnctr.ConnectSubnet( // connect all channels
		r.connections(subnetIO, component.Connections), // map connections map to real channels
	)

	return nodeIO, nil
}

// connections maps network schema with real channels.
func (r Runtime) connections(nodesIO map[string]IO, net []program.Connection) []Connection {
	cc := make([]Connection, len(net))

	for i := range net {
		fromNodeIO, ok := nodesIO[net[i].From.Node]
		if !ok {
			panic("not ok")
		}

		fromAddr := PortAddr{port: net[i].From.Port, idx: net[i].From.Idx, node: net[i].From.Node}
		from, ok := fromNodeIO.Out[fromAddr]
		if !ok {
			panic("not ok")
		}

		to := make([]Port, len(net[i].To))
		for j := range net[i].To {
			toInportAddr := PortAddr{
				node: net[i].To[j].Node,
				port: net[i].To[j].Port,
				idx:  net[i].To[j].Idx,
			}

			toNodeIO := nodesIO[net[i].To[j].Node]
			receiver, ok := toNodeIO.In[toInportAddr]
			if !ok {
				panic("not ok")
			}

			to[j] = Port{Ch: receiver, Addr: toInportAddr}
		}

		cc[i] = Connection{
			From: Port{Ch: from, Addr: fromAddr},
			To:   to,
		}
	}

	return cc
}

// nodeIO creates channels following node meta.
func (r Runtime) nodeIO(nodeMeta program.NodeMeta) IO {
	inports := make(map[PortAddr]chan Msg)

	for port, slots := range nodeMeta.In {
		addr := PortAddr{port: port, node: nodeMeta.Node}
		if addr.node == "root" {
			addr.node = "in"
		}

		if slots == 0 {
			inports[addr] = make(chan Msg)
			continue
		}

		for i := uint8(0); i < slots; i++ {
			addr.idx = i
			inports[addr] = make(chan Msg, slots)
		}
	}

	outports := make(map[PortAddr]chan Msg)

	for port, slots := range nodeMeta.Out {
		addr := PortAddr{port: port, node: nodeMeta.Node}
		if addr.node == "root" {
			addr.node = "out"
		}

		if slots == 0 {
			outports[addr] = make(chan Msg)
			continue
		}

		for idx := uint8(0); idx < slots; idx++ {
			outports[addr] = make(chan Msg)
		}
	}

	return IO{inports, outports}
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
	Addr PortAddr
}

// IO represents node's input and output ports.
type IO struct {
	In, Out Ports
}

// Ports maps network addresses to real channels.
type Ports map[PortAddr]chan Msg

// Slots returns all port-chanells associated with the given array port name.
func (ports Ports) Slots(arrPort string) ([]chan Msg, error) {
	cc := []chan Msg{}
	for addr, ch := range ports {
		if addr.port == arrPort {
			cc = append(cc, ch)
		}
	}

	if len(cc) == 0 {
		return nil, fmt.Errorf("ErrArrPortNotFound: %s", arrPort)
	}

	return cc, nil
}

// Port returns channel associated with the given port address.
func (ports Ports) Port(port string, idx uint8) (chan Msg, error) {
	for addr, ch := range ports {
		if addr.port == port && addr.idx == idx {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("%w: looking for %s, have: %v", ErrPortNotFound, port, ports)
}

// PortAddr describes port address in the
type PortAddr struct {
	node string
	port string
	idx  uint8 // always 0 for normal ports
}

func (addr PortAddr) String() string {
	return fmt.Sprintf("%s.%s", addr.node, addr.port)
}

func New(connector Connector) Runtime {
	return Runtime{
		cnctr: connector,
	}
}

// AbsPortAddr represents absolute port address in the program's network.
type AbsPortAddr struct {
	port Port
	path []string
}

package runtime

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime/program"
)

type (
	Runtime struct {
		connector Connector
		opsRepo   OpsRepo
	}

	OpsRepo interface {
		Opfunc(pkg, name string) Opfunc
	}

	Connector interface {
		Connect([]Connection)
	}

	Connection struct {
		From ConnectionPoint
		To   []ConnectionPoint
	}

	ConnectionPoint struct {
		Ch   chan Msg
		Addr program.FullPortAddr
	}
)

var ErrPortNotFound = errors.New("port not found")

func (r Runtime) Run(prog program.Program) (IO, error) {
	nodes := map[string]IO{}
	for name, node := range prog.Nodes {
		io := r.newNodeIO(node.In, node.Out)

		if node.Type == program.OperatorNode {
			op := r.opsRepo.Opfunc(node.OpRef.Name, node.OpRef.Name)
			if err := op(io); err != nil {
				return IO{}, err
			}
		}

		if node.Type == program.ConstNode {
			for name, cnst := range node.Const {
				out, err := io.Out.Port(PortAddr{Port: name})
				if err != nil {
					return IO{}, err
				}

				var msg Msg
				switch cnst.Type {
				case program.IntType:
					msg = NewIntMsg(cnst.IntValue)
				default:
					return IO{}, errors.New("")
				}

				go func() {
					for {
						out <- msg
					}
				}()
			}
		}

		nodes[name] = io
	}

	conns := make([]Connection, 0, len(prog.Net))
	for _, conn := range prog.Net {
		senderIO, ok := nodes[conn.From.Node]
		if !ok {
			return IO{}, fmt.Errorf("")
		}

		fromPortAddr := PortAddr{Port: conn.From.Port, Slot: conn.From.Slot}
		fromPortCh, err := senderIO.Out.Port(fromPortAddr)
		if err != nil {
			return IO{}, err
		}

		receivers := make([]ConnectionPoint, 0, len(conn.To))
		for _, connToAddr := range conn.To {
			receiverIO, ok := nodes[connToAddr.Node]
			if !ok {
				return IO{}, fmt.Errorf("")
			}

			toPort, err := receiverIO.In.Port(PortAddr{
				Port: connToAddr.Port,
				Slot: connToAddr.Slot,
			})
			if err != nil {
				return IO{}, err
			}

			receivers = append(receivers, ConnectionPoint{Addr: program.FullPortAddr{
				Node: connToAddr.Node,
				Port: connToAddr.Port,
				Slot: connToAddr.Slot,
			}, Ch: toPort})
		}

		conns = append(conns, Connection{
			To: receivers,
			From: ConnectionPoint{
				Ch: fromPortCh,
				Addr: program.FullPortAddr{
					Node: conn.From.Node,
					Port: conn.From.Port,
					Slot: conn.From.Slot,
				},
			},
		})
	}

	io, err := r.rootIO(prog.IORef, nodes)
	if err != nil {
		return IO{}, fmt.Errorf("")
	}

	return io, nil
}

func (Runtime) rootIO(ioRef program.IORef, nodes map[string]IO) (IO, error) {
	io := IO{}

	for _, absAddr := range ioRef.In {
		node, ok := nodes[absAddr.Node]
		if !ok {
			return IO{}, fmt.Errorf("")
		}

		addr := PortAddr{
			Port: absAddr.Port,
			Slot: absAddr.Slot,
		}

		port, err := node.Out.Port(addr)
		if err != nil {
			return IO{}, err
		}

		io.In[addr] = port
	}

	for _, absAddr := range ioRef.Out {
		node, ok := nodes[absAddr.Node]
		if !ok {
			return IO{}, fmt.Errorf("")
		}

		addr := PortAddr{
			Port: absAddr.Port,
			Slot: absAddr.Slot,
		}

		port, err := node.In.Port(addr)
		if err != nil {
			return IO{}, err
		}

		io.In[addr] = port
	}

	return io, nil
}

func (r Runtime) newNodeIO(in, out map[string]program.PortMeta) IO {
	resultIn := make(map[PortAddr]chan Msg)
	for port, meta := range in {
		addr := PortAddr{Port: port}
		if meta.Slots == 0 {
			resultIn[addr] = make(chan Msg, meta.Buf)
			continue
		}
		for slot := uint8(0); slot < meta.Slots; slot++ {
			addr.Slot = slot
			resultIn[addr] = make(chan Msg, meta.Buf)
		}
	}

	resultOut := make(map[PortAddr]chan Msg)
	for port, meta := range out {
		addr := PortAddr{Port: port}
		if meta.Slots == 0 {
			resultOut[addr] = make(chan Msg, meta.Buf)
			continue
		}
		for slot := uint8(0); slot < meta.Slots; slot++ {
			addr.Slot = slot
			resultOut[addr] = make(chan Msg, meta.Buf)
		}
	}

	return IO{resultIn, resultOut}
}

type Opfunc func(IO) error

type IO struct {
	In, Out Ports
}

type PortAddr struct {
	Port string
	Slot uint8
}

type Ports map[PortAddr]chan Msg

func (ports Ports) Port(addr PortAddr) (chan Msg, error) {
	for addr, ch := range ports {
		if addr.Port != addr.Port || addr.Slot != addr.Slot {
			continue
		}
		return ch, nil
	}

	return nil, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
}

func (ports Ports) PortGroup(name string) ([]chan Msg, error) {
	g := []chan Msg{}

	for addr, ch := range ports {
		if addr.Port == name {
			g = append(g, ch)
		}
	}

	if len(g) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrPortNotFound, name)
	}

	return g, nil
}

func New(connector Connector) Runtime {
	return Runtime{
		connector: connector,
	}
}

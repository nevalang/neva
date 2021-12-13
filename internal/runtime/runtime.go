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
		Addr program.PortAddr
	}
)

var ErrPortNotFound = errors.New("port not found")

func (r Runtime) Run(prog program.Program) error {
	nodes := map[string]IO{}

	for name, node := range prog.Nodes {
		io := r.newNodeIO(node.In, node.Out)

		if node.Type == program.OperatorNode {
			op := r.opsRepo.Opfunc(node.OpRef.Name, node.OpRef.Name)
			if err := op(io); err != nil {
				return err
			}
		}

		if node.Type == program.ConstNode {
			for name, cnst := range node.Const {
				out, err := io.Out.Port(PortAddr{Port: name})
				if err != nil {
					return err
				}

				var msg Msg
				switch cnst.Type {
				case program.IntType:
					msg = NewIntMsg(cnst.IntValue)
				default:
					return err
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

	conns := make([]Connection, 0, len(prog.Connections))
	for _, conn := range prog.Connections {
		senderIO, ok := nodes[conn.From.Node]
		if !ok {
			return fmt.Errorf("not ok")
		}

		fromPortAddr := PortAddr{Port: conn.From.Port, Slot: conn.From.Slot}
		fromPortCh, err := senderIO.Out.Port(fromPortAddr)
		if err != nil {
			return err
		}

		receivers := make([]ConnectionPoint, 0, len(conn.To))
		for _, connToAddr := range conn.To {
			receiverIO, ok := nodes[connToAddr.Node]
			if !ok {
				return err
			}

			toPort, err := receiverIO.In.Port(PortAddr{
				Port: connToAddr.Port,
				Slot: connToAddr.Slot,
			})
			if err != nil {
				return err
			}

			receivers = append(receivers, ConnectionPoint{Addr: program.PortAddr{
				Node: connToAddr.Node,
				Port: connToAddr.Port,
				Slot: connToAddr.Slot,
			}, Ch: toPort})
		}

		conns = append(conns, Connection{
			To: receivers,
			From: ConnectionPoint{
				Ch: fromPortCh,
				Addr: program.PortAddr{
					Node: conn.From.Node,
					Port: conn.From.Port,
					Slot: conn.From.Slot,
				},
			},
		})
	}

	return nil
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

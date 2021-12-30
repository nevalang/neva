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
	nodes := make(map[string]IO, len(prog.Nodes))

	for name, node := range prog.Nodes {
		io := r.newNodeIO(node.IO.In, node.IO.Out)

		if node.Type == program.OperatorNode {
			opfunc := r.opsRepo.Opfunc(node.Operator.Name, node.Operator.Name)
			if err := opfunc(io); err != nil {
				return err
			}
		}

		if node.Type == program.ConstNode {
			out, err := io.Out.Port(PortAddr{Port: "out"})
			if err != nil {
				return err
			}

			var msg Msg
			switch node.Const.Type {
			case program.IntType:
				msg = NewIntMsg(node.Const.IntValue)
			default:
				return err
			}

			go func() {
				for {
					out <- msg
				}
			}()
		}

		nodes[name] = io
	}

	conns := make([]Connection, 0, len(prog.Connections))

	for _, conn := range prog.Connections {
		senderIO, ok := nodes[conn.From.Node]
		if !ok {
			return fmt.Errorf("from node not found")
		}

		fromPortCh, err := senderIO.Out.Port(PortAddr{conn.From.Port, conn.From.Idx})
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
				Idx:  connToAddr.Idx,
			})
			if err != nil {
				return err
			}

			receivers = append(receivers, ConnectionPoint{Addr: program.PortAddr{
				Node: connToAddr.Node,
				Port: connToAddr.Port,
				Idx:  connToAddr.Idx,
			}, Ch: toPort})
		}

		conns = append(conns, Connection{
			To: receivers,
			From: ConnectionPoint{
				Ch: fromPortCh,
				Addr: program.PortAddr{
					Node: conn.From.Node,
					Port: conn.From.Port,
					Idx:  conn.From.Idx,
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
		for idx := uint8(0); idx < meta.Slots; idx++ {
			addr.Idx = idx
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
			addr.Idx = slot
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
	Idx  uint8
}

type Ports map[PortAddr]chan Msg

func (ports Ports) Port(addr PortAddr) (chan Msg, error) {
	for addr, ch := range ports {
		if addr.Port == addr.Port && addr.Idx == addr.Idx {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
}

func (ports Ports) PortArray(name string) ([]chan Msg, error) {
	arr := make([]chan Msg, 0, len(ports))

	for addr, ch := range ports {
		if addr.Port == name {
			arr = append(arr, ch)
		}
	}

	if len(arr) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrPortNotFound, name)
	}

	return arr, nil
}

func New(connector Connector) Runtime {
	return Runtime{
		connector: connector,
	}
}

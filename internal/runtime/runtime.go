package runtime

import (
	"errors"
	"fmt"
	"strings"

	"github.com/emil14/respect/internal/runtime/program"
)

type (
	Runtime struct {
		connector Connector
		opsRepo   OperatorsRepo
	}

	OperatorsRepo interface {
		OperatorFunc(pkg, name string) OperatorFunc
	}

	Connector interface {
		Connect([]Connection)
	}
)

var ErrPortNotFound = errors.New("port not found")

func (r Runtime) Run(prog program.Program) (IO, error) {
	nodes := map[string]IO{}
	for _, node := range prog.Nodes {
		io := r.nodeIO(node.In, node.Out)

		if node.Type == program.OperatorNode {
			op := r.opsRepo.OperatorFunc(node.Operator.Name, node.Operator.Name)
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

		nodes[strings.Join(node.Path, ".")] = io
	}

	net := make([]Connection, 0, len(prog.Net))
	for _, conn := range prog.Net {
		senderIO, ok := nodes[strings.Join(conn.From.NodePath, ".")]
		if !ok {
			return IO{}, fmt.Errorf("")
		}

		outportAddr := PortAddr{
			Port: conn.From.Port,
			Slot: conn.From.Slot,
		}

		outport, err := senderIO.Out.Port(outportAddr)
		if err != nil {
			return IO{}, err
		}

		receivers := make([]Port, 0, len(conn.To))
		for _, inportAddr := range conn.To {
			receiverIO, ok := nodes[strings.Join(inportAddr.NodePath, ".")]
			if !ok {
				return IO{}, fmt.Errorf("")
			}

			to, err := receiverIO.In.Port(PortAddr{
				Port: inportAddr.Port,
				Slot: inportAddr.Slot,
			})
			if err != nil {
				return IO{}, err
			}

			receivers = append(receivers, Port{Addr: PortAddr{
				Port: inportAddr.Port,
				Slot: inportAddr.Slot,
			}, Ch: to})
		}

		net = append(net, Connection{
			To: receivers,
			From: Port{
				Ch:   outport,
				Addr: outportAddr,
			},
		})
	}

	io := IO{}
	for _, absAddr := range prog.IO.In {
		node, ok := nodes[strings.Join(absAddr.NodePath, ".")]
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
	for _, absAddr := range prog.IO.Out {
		node, ok := nodes[strings.Join(absAddr.NodePath, ".")]
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

func (r Runtime) nodeIO(in, out map[string]program.PortMeta) IO {
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

type OperatorFunc func(IO) error

type Connection struct {
	From Port
	To   []Port
}

type Port struct {
	Ch   chan Msg
	Addr PortAddr
}

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

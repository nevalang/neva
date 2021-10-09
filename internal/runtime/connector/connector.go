package connector

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/program"
)

type Connector struct {
	operators   map[string]runtime.Operator
	interceptor Interceptor
}

func (cnctr Connector) ConnectSubnet(cc []runtime.Connection) {
	for _, c := range cc {
		go cnctr.connectConnection(c)
	}
}

func (cnctr Connector) connectConnection(conn runtime.Connection) {
	for msg := range conn.From.Ch {
		cnctr.interceptor.onSend(msg, conn.From.Addr, nil)

		for i := range conn.To {
			to := conn.To[i]

			go func(m runtime.Msg) {
				to.Ch <- m
				cnctr.interceptor.onReceive(m, conn.From.Addr, to.Addr)
			}(msg)
		}
	}
}

func (c Connector) ConnectOperator(name string, io runtime.IO) error {
	op, ok := c.operators[name]
	if !ok {
		return fmt.Errorf("ErrUnknownOperator: %s", name)
	}

	if err := op(io); err != nil {
		return err
	}

	return nil
}

type Interceptor interface {
	onSend(msg runtime.Msg, from program.PortAddr, to []program.PortAddr) runtime.Msg
	onReceive(msg runtime.Msg, from, to program.PortAddr)
}

type interceptor struct {
	send    func(msg runtime.Msg, from program.PortAddr, to []program.PortAddr) runtime.Msg
	receive func(msg runtime.Msg, from, to program.PortAddr)
}

func (i interceptor) onSend(msg runtime.Msg, from program.PortAddr, to []program.PortAddr) runtime.Msg {
	return i.send(msg, from, to)
}

func (i interceptor) onReceive(msg runtime.Msg, from, to program.PortAddr) {
	i.receive(msg, from, to)
}

func New(
	operators map[string]runtime.Operator,
	// TODO: use Interceptor types
	onSend func(msg runtime.Msg, from program.PortAddr, to []program.PortAddr) runtime.Msg,
	onReceive func(msg runtime.Msg, from, to program.PortAddr),
) (Connector, error) {
	if operators == nil || onSend == nil || onReceive == nil {
		return Connector{}, errors.New("init connector")
	}

	return Connector{
		operators,
		interceptor{onSend, onReceive},
	}, nil
}

func MustNew(
	ops map[string]runtime.Operator,
	onSend func(msg runtime.Msg, from program.PortAddr, to []program.PortAddr) runtime.Msg,
	onReceive func(msg runtime.Msg, from, to program.PortAddr),
) Connector {
	c, err := New(ops, onSend, onReceive)
	if err != nil {
		panic(err)
	}
	return c
}

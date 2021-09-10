package connector

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
)

type Connector struct {
	ops       map[string]runtime.Operator
	onSend    func(msg runtime.Msg, from runtime.PortAddr)
	onReceive func(msg runtime.Msg, from, to runtime.PortAddr)
}

func (cnctr Connector) ConnectSubnet(pp []runtime.Connection) {
	for i := range pp {
		go cnctr.connectConnection(pp[i])
	}
}

func (cnctr Connector) connectConnection(conn runtime.Connection) {
	for msg := range conn.From.Ch {
		cnctr.onSend(msg, conn.From.Addr)

		for i := range conn.To {
			to := conn.To[i]

			go func(m runtime.Msg) {
				to.Ch <- m
				cnctr.onReceive(m, conn.From.Addr, to.Addr)
			}(msg)
		}
	}
}

func (c Connector) ConnectOperator(name string, io runtime.IO) error {
	op, ok := c.ops[name]
	if !ok {
		return fmt.Errorf("ErrUnknownOperator: %s", name)
	}

	if err := op(io); err != nil {
		return err
	}

	return nil
}

func New(
	ops map[string]runtime.Operator,
	onSend func(msg runtime.Msg, from runtime.PortAddr),
	onReceive func(msg runtime.Msg, from, to runtime.PortAddr),
) (Connector, error) {
	if ops == nil || onSend == nil || onReceive == nil {
		return Connector{}, errors.New("init connector")
	}

	return Connector{ops, onSend, onReceive}, nil
}

func MustNew(
	ops map[string]runtime.Operator,
	onSend func(msg runtime.Msg, from runtime.PortAddr),
	onReceive func(msg runtime.Msg, from, to runtime.PortAddr),
) Connector {
	c, err := New(ops, onSend, onReceive)
	if err != nil {
		panic(err)
	}
	return c
}

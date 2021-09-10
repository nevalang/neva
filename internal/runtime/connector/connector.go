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

func (c Connector) ConnectSubnet(pp []runtime.Connection) {
	for i := range pp {
		go c.connectPair(pp[i])
	}
}

func (c Connector) connectPair(con runtime.Connection) {
	for msg := range con.From.Ch {
		c.onSend(msg, con.From.Addr)

		for _, recv := range con.To {
			select {
			case recv.Ch <- msg:
				c.onReceive(msg, con.From.Addr, recv.Addr)
			default:
				go func(to runtime.Port, m runtime.Msg) {
					to.Ch <- m
					c.onReceive(m, con.From.Addr, to.Addr)
				}(recv, msg)
			}
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

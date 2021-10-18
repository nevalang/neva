package connector

import (
	"errors"
	"fmt"
	"sync"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/program"
)

type Connector struct {
	operators   map[string]runtime.Operator
	interceptor Interceptor
}

func (cnctr Connector) ConnectSubnet(cc []runtime.Connection) {
	for _, c := range cc {
		go cnctr.loop(c)
	}
}

func (cnctr Connector) loop(conn runtime.Connection) {
	for msg := range conn.From.Ch {
		cnctr.interceptor.OnSend(msg, conn.From.Addr)

		wg := sync.WaitGroup{}
		wg.Add(len(conn.To))

		for i := range conn.To {
			to := conn.To[i]
			m := msg

			go func() {
				fmt.Printf("start %s from %s to %s\n", m, conn.From.Addr, to.Addr)
				to.Ch <- m
				cnctr.interceptor.OnReceive(m, conn.From.Addr, to.Addr)
				wg.Done()
			}()
		}

		wg.Wait()
	}
}

func (c Connector) ConnectOperator(name string, io runtime.IO) error {
	op, ok := c.operators[name]
	if !ok {
		return fmt.Errorf("unknown operator: %s", name)
	}

	if err := op(io); err != nil {
		return fmt.Errorf("operator '%s': %w", name, err)
	}

	return nil
}

type Interceptor interface {
	OnSend(msg runtime.Msg, from program.PortAddr) runtime.Msg
	OnReceive(msg runtime.Msg, from, to program.PortAddr)
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

func New(ops map[string]runtime.Operator, interceptor Interceptor) (Connector, error) {
	if ops == nil || interceptor == nil {
		return Connector{}, errors.New("init connector")
	}

	return Connector{
		ops,
		interceptor,
	}, nil
}

func MustNew(
	ops map[string]runtime.Operator,
	interceptor Interceptor,
) Connector {
	c, err := New(ops, interceptor)
	if err != nil {
		panic(err)
	}
	return c
}

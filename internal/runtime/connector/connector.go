package connector

import (
	"errors"
	"sync"

	"github.com/emil14/respect/internal/core"
	"github.com/emil14/respect/internal/runtime"
)

type Connector struct {
	interceptor Interceptor
}

func (c Connector) Connect(cc []runtime.Connection) {
	for _, connection := range cc {
		go c.loop(connection)
	}
}

func (cnctr Connector) loop(connection runtime.Connection) {
	for msg := range connection.From.Ch {
		cnctr.interceptor.OnSend(msg, connection.From.Addr)

		wg := sync.WaitGroup{}
		wg.Add(len(connection.To))

		for i := range connection.To {
			to := connection.To[i]
			m := msg

			go func() {
				to.Ch <- m
				cnctr.interceptor.OnReceive(m, connection.From.Addr, to.Addr)
				wg.Done()
			}()
		}

		wg.Wait()
	}
}

type Interceptor interface {
	OnSend(msg core.Msg, from core.PortAddr) core.Msg
	OnReceive(msg core.Msg, from, to core.PortAddr)
}

type interceptor struct {
	send    func(msg core.Msg, from core.PortAddr, to []core.PortAddr) core.Msg
	receive func(msg core.Msg, from, to core.PortAddr)
}

func (i interceptor) onSend(msg core.Msg, from core.PortAddr, to []core.PortAddr) core.Msg {
	return i.send(msg, from, to)
}

func (i interceptor) onReceive(msg core.Msg, from, to core.PortAddr) {
	i.receive(msg, from, to)
}

func New(ops map[string]runtime.OperatorFunc, interceptor Interceptor) (Connector, error) {
	if interceptor == nil {
		return Connector{}, errors.New("init connector")
	}

	return Connector{
		interceptor,
	}, nil
}

func MustNew(
	interceptor Interceptor,
) Connector {
	c, err := New(interceptor)
	if err != nil {
		panic(err)
	}
	return c
}

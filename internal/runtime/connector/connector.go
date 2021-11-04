package connector

import (
	"errors"
	"sync"

	"github.com/emil14/respect/internal/core"
	"github.com/emil14/respect/internal/runtime"
)

type (
	Connector struct {
		interceptor Interceptor
	}

	Interceptor interface {
		OnSend(msg core.Msg, from core.PortAddr) core.Msg
		OnReceive(msg core.Msg, from, to core.PortAddr)
	}
)

func (c Connector) Connect(cc []runtime.Connection) {
	for _, connection := range cc {
		go c.loop(connection)
	}
}

func (c Connector) loop(connection runtime.Connection) {
	for msg := range connection.From.Ch {
		c.interceptor.OnSend(msg, connection.From.Addr)

		wg := sync.WaitGroup{}
		wg.Add(len(connection.To))

		for i := range connection.To {
			to := connection.To[i]
			m := msg

			go func() {
				to.Ch <- m
				c.interceptor.OnReceive(m, connection.From.Addr, to.Addr)
				wg.Done()
			}()
		}

		wg.Wait()
	}
}

func New(interceptor Interceptor) (Connector, error) {
	if interceptor == nil {
		return Connector{}, errors.New("nil interceptor")
	}

	return Connector{
		interceptor,
	}, nil
}

func MustNew(interceptor Interceptor) Connector {
	c, err := New(interceptor)
	if err != nil {
		panic(err)
	}
	return c
}

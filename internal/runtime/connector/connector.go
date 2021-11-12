package connector

import (
	"errors"
	"sync"

	"github.com/emil14/respect/internal/runtime"
)

type Connector struct {
	interceptor Interceptor
}

func (c Connector) Connect(conns []runtime.Connection) {
	for _, conn := range conns {
		go c.loop(conn)
	}
}

func (c Connector) loop(conn runtime.Connection) {
	for msg := range conn.From.Ch {
		c.interceptor.OnSend(msg, conn.From.Addr)

		wg := sync.WaitGroup{}
		wg.Add(len(conn.To))

		for i := range conn.To {
			to := conn.To[i]
			m := msg

			go func() {
				to.Ch <- m
				c.interceptor.OnReceive(m, conn.From.Addr, to.Addr)
				wg.Done()
			}()
		}

		wg.Wait()
	}
}

type Interceptor interface {
	OnSend(msg runtime.Msg, from runtime.PortAddr) runtime.Msg
	OnReceive(msg runtime.Msg, from, to runtime.PortAddr)
}

func New(interceptor Interceptor) (Connector, error) {
	if interceptor == nil {
		return Connector{}, errors.New("nil interceptor")
	}
	return Connector{interceptor}, nil
}

func MustNew(
	ops map[string]runtime.OperatorFunc,
	interceptor Interceptor,
) Connector {
	c, err := New(interceptor)
	if err != nil {
		panic(err)
	}
	return c
}

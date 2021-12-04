package connector

import (
	"errors"
	"sync"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/program"
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
	OnSend(msg runtime.Msg, from program.FullPortAddr) runtime.Msg
	OnReceive(msg runtime.Msg, from, to program.FullPortAddr)
}

func New(i Interceptor) (Connector, error) {
	if i == nil {
		return Connector{}, errors.New("nil interceptor")
	}
	return Connector{i}, nil
}

func MustNew(i Interceptor) Connector {
	c, err := New(i)
	if err != nil {
		panic(err)
	}
	return c
}

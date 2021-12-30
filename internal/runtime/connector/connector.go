package connector

import (
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/program"
)

type Interceptor interface {
	Send(msg runtime.Msg, from program.PortAddr) runtime.Msg
	Receive(msg runtime.Msg, from, to program.PortAddr) runtime.Msg
}

type Connector struct {
	i Interceptor
}

func (c Connector) Connect(cc []runtime.Connection) {
	for _, conn := range cc {
		go c.rwloop(conn)
	}
}

func (c Connector) rwloop(conn runtime.Connection) {
	for msg := range conn.From.Ch {
		msg = c.i.Send(msg, conn.From.Addr)

		for i := range conn.To {
			m := msg
			to := conn.To[i]

			go func() {
				to.Ch <- c.i.Receive(m, conn.From.Addr, to.Addr)
			}()
		}
	}
}

func New(i Interceptor) (Connector, error) {
	if err := utils.NilArgs(i); err != nil {
		return Connector{}, err
	}
	return Connector{i}, nil
}

func MustNew(i Interceptor) Connector {
	c, err := New(i)
	utils.MaybePanic(err)
	return c
}

package connector

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
	"github.com/emil14/neva/internal/pkg/utils"
)

type (
	Interceptor interface {
		AfterSend(runtime.Connection, core.Msg) core.Msg
		BeforeReceive(from, to runtime.PortAddr, msg core.Msg) core.Msg
	}

	Mapper interface {
		Net(map[runtime.PortAddr]chan core.Msg, []runtime.Connection) ([]Connection, error)
	}

	Connection struct {
		original runtime.Connection
		from     chan core.Msg
		to       []chan core.Msg
	}
)

var ErrMapper = errors.New("mapper")

type Connector struct {
	interceptor Interceptor
	mapper      Mapper
}

func (c Connector) Connect(ports map[runtime.PortAddr]chan core.Msg, net []runtime.Connection) error {
	connections, err := c.mapper.Net(ports, net)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMapper, err)
	}

	for i := range connections {
		go c.stream(connections[i])
	}

	return nil
}

func (c Connector) stream(connection Connection) {
	for msg := range connection.from {
		msg = c.interceptor.AfterSend(connection.original, msg)

		for i := range connection.to {
			toAddr := connection.original.To[i]
			toPort := connection.to[i]

			go func(m core.Msg) {
				toPort <- c.interceptor.BeforeReceive(connection.original.From, toAddr, m)
			}(msg)
		}
	}
}

func MustNew(i Interceptor) Connector {
	utils.NilArgsFatal(i)

	return Connector{
		interceptor: i,
		mapper:      mapper{},
	}
}

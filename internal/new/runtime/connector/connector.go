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
		BeforeReceive(from, to runtime.AbsPortAddr, msg core.Msg) core.Msg
	}

	Mapper interface {
		Net(map[string]core.IO, []runtime.Connection) ([]Connection, error)
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

func (c Connector) Connect(nodesIO map[string]core.IO, net []runtime.Connection) error {
	connections, err := c.mapper.Net(nodesIO, net)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMapper, err)
	}

	for i := range connections {
		go c.stream(connections[i], nodesIO)
	}

	return nil
}

func (c Connector) stream(connection Connection, nodesIO map[string]core.IO) {
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

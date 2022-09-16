package connector

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime"
)

type (
	Interceptor interface {
		AfterSending(runtime.Connection, core.Msg) core.Msg
		BeforeReceiving(from, to runtime.AbsolutePortAddr, msg core.Msg) core.Msg
		AfterReceiving(from, to runtime.AbsolutePortAddr, msg core.Msg)
	}

	ChanMapper interface {
		ConnectionsWithChans(
			portChans map[runtime.AbsolutePortAddr]chan core.Msg,
			connections []runtime.Connection,
		) ([]ConnectionWithChans, error)
	}

	ConnectionWithChans struct {
		sender    chan core.Msg
		receivers []chan core.Msg
		meta      runtime.Connection
	}
)

var ErrMapper = errors.New("mapper")

type Connector struct {
	interceptor Interceptor
	mapper      ChanMapper
}

func (c Connector) Connect(ports map[runtime.AbsolutePortAddr]chan core.Msg, rels []runtime.Connection) error {
	connections, err := c.mapper.ConnectionsWithChans(ports, rels)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMapper, err)
	}

	for i := range connections {
		go c.linkConnection(connections[i])
	}

	return nil
}

func (c Connector) linkConnection(connection ConnectionWithChans) {
	guard := make(chan struct{}, len(connection.receivers))

	for msg := range connection.sender {
		msg = c.interceptor.AfterSending(connection.meta, msg)

		for i := range connection.receivers {
			receiverConnPoint := connection.meta.ReceiversConnectionPoints[i]
			receiverPortChan := connection.receivers[i]

			if receiverConnPoint.Type == runtime.StructFieldReading {
				parts := receiverConnPoint.StructFieldPath
				for _, part := range parts[:len(parts)-1] {
					msg = msg.Struct()[part]
				}
			}

			guard <- struct{}{}

			go func(m core.Msg) {
				receiverPortChan <- c.interceptor.BeforeReceiving(connection.meta.SenderPortAddr, receiverConnPoint.PortAddr, m)
				c.interceptor.AfterReceiving(connection.meta.SenderPortAddr, receiverConnPoint.PortAddr, m)
				<-guard
			}(msg)
		}
	}
}

func MustNew(i Interceptor) Connector {
	utils.PanicOnNil(i)

	return Connector{
		interceptor: i,
		mapper:      mapper{},
	}
}

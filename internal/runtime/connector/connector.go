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
		info      runtime.Connection
	}
)

var ErrMapper = errors.New("mapper")

type Connector struct {
	interceptor Interceptor
	mapper      ChanMapper
}

func (c Connector) Connect(ports map[runtime.AbsolutePortAddr]chan core.Msg, connections []runtime.Connection) error {
	connectionsWithChans, err := c.mapper.ConnectionsWithChans(ports, connections)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMapper, err)
	}

	for i := range connectionsWithChans {
		go c.connect(connectionsWithChans[i])
	}

	return nil
}

func (c Connector) connect(connection ConnectionWithChans) {
	guard := make(chan struct{}, len(connection.receivers))

	for msg := range connection.sender {
		msg = c.interceptor.AfterSending(connection.info, msg)

		for i := range connection.receivers {
			receiverPortChan := connection.receivers[i]
			receiverConnPoint := connection.info.ReceiversConnectionPoints[i]

			if receiverConnPoint.Type == runtime.StructFieldReading {
				parts := receiverConnPoint.StructFieldPath
				for _, part := range parts[:len(parts)-1] {
					msg = msg.Struct()[part] // FIXME possible panic
				}
			}

			guard <- struct{}{}

			go func(m core.Msg) {
				receiverPortChan <- c.interceptor.BeforeReceiving(connection.info.SenderPortAddr, receiverConnPoint.PortAddr, m)
				c.interceptor.AfterReceiving(connection.info.SenderPortAddr, receiverConnPoint.PortAddr, m)
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

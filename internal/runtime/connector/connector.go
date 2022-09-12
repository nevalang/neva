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
		BeforeReceiving(from, to runtime.PortAddr, msg core.Msg) core.Msg
	}

	Mapper interface {
		Net(map[runtime.PortAddr]chan core.Msg, []runtime.Connection) ([]Connection, error)
	}
)

var ErrMapper = errors.New("mapper")

type Connector struct {
	interceptor Interceptor
	mapper      Mapper
}

type Connection struct {
	sender    chan core.Msg
	receivers []chan core.Msg
	meta      runtime.Connection
}

func (c Connector) Connect(ports map[runtime.PortAddr]chan core.Msg, rels []runtime.Connection) error {
	connections, err := c.mapper.Net(ports, rels)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMapper, err)
	}

	for i := range connections {
		go c.linkConnection(connections[i])
	}

	return nil
}

func (c Connector) linkConnection(relation Connection) {
	guard := make(chan struct{}, len(relation.receivers))

	for msg := range relation.sender {
		msg = c.interceptor.AfterSending(relation.meta, msg)

		for i := range relation.receivers {
			receiverConnPoint := relation.meta.Receivers[i]
			receiverPortChan := relation.receivers[i]

			if receiverConnPoint.Type == runtime.FieldReading {
				parts := receiverConnPoint.StructFieldPath
				for _, part := range parts[:len(parts)-1] {
					msg = msg.Struct()[part]
				}
			}

			guard <- struct{}{}

			go func(m core.Msg) {
				receiverPortChan <- c.interceptor.BeforeReceiving(relation.meta.Sender, receiverConnPoint.PortAddr, m)
				<-guard
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

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
		AfterSend(runtime.Connection, core.Msg) core.Msg                // After Sending
		BeforeReceive(from, to runtime.PortAddr, msg core.Msg) core.Msg // Before receiving
	}

	RelationMapper interface { // TODO do we need business-logic behind interface?
		Net(map[runtime.PortAddr]chan core.Msg, []runtime.Connection) ([]Relation, error) // TODO rename?
	}
)

var ErrMapper = errors.New("mapper")

type Connector struct {
	interceptor Interceptor
	mapper      RelationMapper
}

type Relation struct {
	sender    chan core.Msg
	receivers []chan core.Msg
	meta      runtime.Connection
}

func (c Connector) Connect(ports map[runtime.PortAddr]chan core.Msg, rels []runtime.Connection) error {
	connections, err := c.mapper.Net(ports, rels) // TODO refactor?
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMapper, err)
	}

	for i := range connections {
		go c.linkConnection(connections[i])
	}

	return nil
}

func (c Connector) linkConnection(relation Relation) {
	guard := make(chan struct{}, len(relation.receivers))

	for msg := range relation.sender { // receiving
		msg = c.interceptor.AfterSend(relation.meta, msg)

		for i := range relation.receivers { // delivery
			receiverConnPoint := relation.meta.Receivers[i]
			receiverPortChan := relation.receivers[i]

			if receiverConnPoint.Type == runtime.FieldReading { // TODO move to core?
				parts := receiverConnPoint.StructFieldPath
				for _, part := range parts[:len(parts)-1] {
					msg = msg.Struct()[part]
				}
			}

			guard <- struct{}{}

			go func(m core.Msg) {
				receiverPortChan <- c.interceptor.BeforeReceive(relation.meta.Sender, receiverConnPoint.PortAddr, m)
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

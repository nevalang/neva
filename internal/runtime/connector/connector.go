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
		AfterSend(runtime.Relation, core.Msg) core.Msg                  // After Sending
		BeforeReceive(from, to runtime.PortAddr, msg core.Msg) core.Msg // Before receiving
	}

	RelationMapper interface { // TODO do we need business-logic behind interface?
		Net(map[runtime.PortAddr]chan core.Msg, []runtime.Relation) ([]Relation, error) // TODO rename?
	}
}

var ErrMapper = errors.New("mapper")

type Connector struct {
	interceptor Interceptor
	mapper      RelationMapper
}

func (c Connector) Connect(ports map[runtime.PortAddr]chan core.Msg, rels []runtime.Relation) error {
	connections, err := c.mapper.Net(ports, rels)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMapper, err)
	}

	for i := range connections {
		go c.linkConnection(connections[i])
	}

		wg.Wait()
	}
}

func (c Connector) linkConnection(relation Relation) {
	guard := make(chan struct{}, len(relation.receivers))

	for msg := range relation.sender { // receiving
		msg = c.interceptor.AfterSend(relation.meta, msg)

		for i := range relation.receivers { // delivery
			guard <- struct{}{}

			toAddr := relation.meta.Receivers[i]
			toPort := relation.receivers[i]

			go func(m core.Msg) {
				toPort <- c.interceptor.BeforeReceive(relation.meta.Sender, toAddr, m) // FIXME possible memory leak
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

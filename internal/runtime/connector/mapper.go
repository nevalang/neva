package connector

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

var ErrPortNotFound = errors.New("port not found by addr")

type mapper struct{}

func (m mapper) Net(ports map[runtime.PortAddr]chan core.Msg, net []runtime.Relation) ([]Relation, error) {
	connections := make([]Relation, len(net))

	for i := range net {
		from, ok := ports[net[i].Sender]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrPortNotFound, net[i].Sender)
		}

		connections[i] = Relation{
			meta:      net[i],
			sender:    from,
			receivers: make([]chan core.Msg, len(net[i].Receivers)),
		}

		for j := range net[i].Receivers {
			to, ok := ports[net[i].Receivers[j]]
			if !ok {
				return nil, fmt.Errorf("%w: %v", ErrPortNotFound, net[i].Receivers[j])
			}

			connections[i].receivers[j] = to
		}
	}

	return connections, nil
}

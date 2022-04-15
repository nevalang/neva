package connector

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

var ErrPortNotFound = errors.New("port not found by addr")

type mapper struct{}

func (m mapper) Net(ports map[runtime.PortAddr]chan core.Msg, net []runtime.Connection) ([]Connection, error) {
	connections := make([]Connection, len(net))

	for i := range net {
		from, ok := ports[net[i].From]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrPortNotFound, net[i].From)
		}

		connections[i] = Connection{
			original: net[i],
			from:     from,
			to:       make([]chan core.Msg, len(net[i].To)),
		}

		for j := range net[i].To {
			to, ok := ports[net[i].To[j]]
			if !ok {
				return nil, fmt.Errorf("%w: %v", ErrPortNotFound, net[i].To[j])
			}

			connections[i].to[j] = to
		}
	}

	return connections, nil
}

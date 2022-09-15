package connector

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

var ErrPortNotFound = errors.New("port not found by addr")

type mapper struct{}

func (m mapper) ConnectionsWithChans(
	ports map[runtime.AbsolutePortAddr]chan core.Msg,
	net []runtime.Connection,
) ([]ConnectionWithChans, error) {
	connections := make([]ConnectionWithChans, len(net))

	for i := range net {
		from, ok := ports[net[i].SenderPortAddr]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrPortNotFound, net[i].SenderPortAddr)
		}

		connections[i] = ConnectionWithChans{
			meta:      net[i],
			sender:    from,
			receivers: make([]chan core.Msg, len(net[i].ReceiversConnectionPoints)),
		}

		for j := range net[i].ReceiversConnectionPoints {
			to, ok := ports[net[i].ReceiversConnectionPoints[j].PortAddr]
			if !ok {
				return nil, fmt.Errorf("%w: %v", ErrPortNotFound, net[i].ReceiversConnectionPoints[j])
			}

			connections[i].receivers[j] = to
		}
	}

	return connections, nil
}

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
	connections []runtime.Connection,
) ([]ConnectionWithChans, error) {
	result := make([]ConnectionWithChans, len(connections))

	for i := range connections {
		from, ok := ports[connections[i].SenderPortAddr]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrPortNotFound, connections[i].SenderPortAddr)
		}

		result[i] = ConnectionWithChans{
			meta:      connections[i],
			sender:    from,
			receivers: make([]chan core.Msg, len(connections[i].ReceiversConnectionPoints)),
		}

		for j := range connections[i].ReceiversConnectionPoints {
			to, ok := ports[connections[i].ReceiversConnectionPoints[j].PortAddr]
			if !ok {
				return nil, fmt.Errorf("%w: %v", ErrPortNotFound, connections[i].ReceiversConnectionPoints[j])
			}

			result[i].receivers[j] = to
		}
	}

	return result, nil
}

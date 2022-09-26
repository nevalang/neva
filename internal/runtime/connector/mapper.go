package connector

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/src"
)

var ErrPortNotFound = errors.New("port not found by addr")

type mapper struct{}

func (m mapper) MapPortsToConnections(
	ports map[src.AbsolutePortAddr]chan core.Msg,
	connections []runtime.Connection,
) ([]Connection, error) {
	result := make([]Connection, len(connections))

	for i := range connections {
		from, ok := ports[connections[i].SenderPortAddr]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrPortNotFound, connections[i].SenderPortAddr)
		}

		result[i] = Connection{
			info:      connections[i],
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

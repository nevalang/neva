package connector

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

var (
	ErrSelectPort   = errors.New("select port")
	ErrNodeNotFound = errors.New("node not found by addr")
	ErrPortNotFound = errors.New("port not found by addr")
)

type mapper struct{}

func (m mapper) Net(nodesIO map[string]core.IO, net []runtime.Connection) ([]Connection, error) {
	result := make([]Connection, len(net))

	for i := range net {
		from, err := m.inport(net[i].From, nodesIO)
		if err != nil {
			return nil, fmt.Errorf("inport: %v: %w", net[i].From, err)
		}

		result[i] = Connection{
			original: net[i],
			from:     from,
			to:       make([]chan core.Msg, len(net[i].To)),
		}

		for j := range net[i].To {
			to, err := m.outport(net[j].To[j], nodesIO)
			if err != nil {
				return nil, fmt.Errorf("outport: %v: %w", net[j].To[j], err)
			}

			result[i].to[j] = to
		}
	}

	return result, nil
}

func (m mapper) inport(addr runtime.FullPortAddr, nodesIO map[string]core.IO) (chan core.Msg, error) {
	io, ok := nodesIO[addr.Path]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrNodeNotFound, addr)
	}

	port, ok := io.In[core.PortAddr{Port: addr.Port, Idx: addr.Idx}]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
	}

	return port, nil
}

func (m mapper) outport(addr runtime.FullPortAddr, nodesIO map[string]core.IO) (chan core.Msg, error) {
	io, ok := nodesIO[addr.Path]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrNodeNotFound, addr)
	}

	port, ok := io.Out[core.PortAddr{Port: addr.Port, Idx: addr.Idx}]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
	}

	return port, nil
}

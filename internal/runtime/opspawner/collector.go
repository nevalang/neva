package opspawner

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

var ErrPortNotFound = errors.New("port not found")

type collector struct{}

func (c collector) Collect(wantIO runtime.OpPortAddrs, ports map[runtime.PortAddr]chan core.Msg) (core.IO, error) {
	io := core.IO{}

	for _, addr := range wantIO.In {
		port, ok := ports[addr]
		if !ok {
			return core.IO{}, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		io.In[core.PortAddr{
			Port: addr.Name,
			Idx:  addr.Idx,
		}] = port
	}

	for _, addr := range wantIO.Out {
		port, ok := ports[addr]
		if !ok {
			return core.IO{}, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		io.Out[core.PortAddr{
			Port: addr.Name,
			Idx:  addr.Idx,
		}] = port
	}

	return io, nil
}

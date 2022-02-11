package opspawner

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

var ErrPortNotFound = errors.New("port not found")

type collector struct{}

func (c collector) Collect(wantIO runtime.OperatorIO, ports map[runtime.FullPortAddr]chan core.Msg) (core.IO, error) {
	io := core.IO{}

	for _, addr := range wantIO.In {
		port, ok := ports[addr]
		if !ok {
			return core.IO{}, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		io.In[core.PortAddr{
			Port: addr.Port,
			Idx:  addr.Idx,
		}] = port
	}

	for _, addr := range wantIO.Out {
		port, ok := ports[addr]
		if !ok {
			return core.IO{}, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		io.Out[core.PortAddr{
			Port: addr.Port,
			Idx:  addr.Idx,
		}] = port
	}

	return io, nil
}

package opspawner

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

var ErrPortNotFound = errors.New("port not found")

type Searcher struct{}

func (s Searcher) SearchPorts( // we lookup ports inside operator funcs and we also lookup them here?
	wantIO runtime.OperatorPortAddrs,
	ports map[runtime.AbsolutePortAddr]chan core.Msg,
) (core.IO, error) {
	io := core.IO{}

	for _, addr := range wantIO.In {
		port, ok := ports[addr]
		if !ok {
			return core.IO{}, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		io.In[core.RelativePortAddr{
			Port: addr.Port,
			Idx:  addr.Idx,
		}] = port
	}

	for _, addr := range wantIO.Out {
		port, ok := ports[addr]
		if !ok {
			return core.IO{}, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		io.Out[core.RelativePortAddr{
			Port: addr.Port,
			Idx:  addr.Idx,
		}] = port
	}

	return io, nil
}

package portgen

import (
	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type PortGen struct{}

func (p PortGen) Ports(io runtime.NodeIO) core.IO {
	return core.IO{
		In:  p.ports(io.In),
		Out: p.ports(io.Out),
	}
}

func (PortGen) ports(ports map[string]runtime.PortMeta) map[core.PortAddr]chan core.Msg {
	res := make(map[core.PortAddr]chan core.Msg)

	for portName, meta := range ports {
		addr := core.PortAddr{Port: portName}

		if meta.Slots == 0 {
			res[addr] = make(chan core.Msg, meta.Buf)
			continue
		}

		for idx := uint8(0); idx < meta.Slots; idx++ {
			addr.Idx = idx
			res[addr] = make(chan core.Msg, meta.Buf)
		}
	}

	return res
}

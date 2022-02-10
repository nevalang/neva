package portgen

import (
	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type PortGen struct{}

func (p PortGen) Ports(io runtime.OperatorIO) core.IO {
	return core.IO{
		In:  p.ports(io.In),
		Out: p.ports(io.Out),
	}
}

func (PortGen) ports(ports map[runtime.PortAddr]runtime.Port) map[core.PortAddr]chan core.Msg {
	corePorts := make(map[core.PortAddr]chan core.Msg)

	for addr, port := range ports {
		corePorts[core.PortAddr(addr)] = make(chan core.Msg, port.Buf)
	}

	return corePorts
}

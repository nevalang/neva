package portgen

import (
	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type PortGen struct{}

func (p PortGen) Ports(addrs []runtime.PortAddr) map[runtime.PortAddr]chan core.Msg {
	ports := make(map[runtime.PortAddr]chan core.Msg, len(addrs))
	for i := range addrs {
		ports[addrs[i]] = make(chan core.Msg)
	}
	return ports
}

func New() PortGen {
	return PortGen{}
}

package portgen

import (
	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type PortGen struct{}

func (p PortGen) Ports(src map[runtime.FullPortAddr]runtime.Port) map[runtime.FullPortAddr]chan core.Msg {
	ports := make(map[runtime.FullPortAddr]chan core.Msg, len(src))

	for addr := range src {
		ports[addr] = make(chan core.Msg, src[addr].Buf)
	}

	return ports
}

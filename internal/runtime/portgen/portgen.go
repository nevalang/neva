package portgen

import (
	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime/src"
)

type PortGen struct{}

func (p PortGen) Ports(in []src.Port) map[src.AbsolutePortAddr]chan core.Msg {
	out := make(
		map[src.AbsolutePortAddr]chan core.Msg,
		len(in),
	)
	for i := range in {
		out[in[i].Addr] = make(chan core.Msg, in[i].Buf)
	}
	return out
}

func New() PortGen {
	return PortGen{}
}

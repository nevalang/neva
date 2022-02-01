package constspawner

import (
	"errors"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type Spawner struct{}

var ErrPortNotFound = errors.New("port not found")

func (c Spawner) Spawn(messages map[runtime.RelPortAddr]runtime.ConstMsg, ports map[core.PortAddr]chan core.Msg) error {
	for addr := range messages {
		port, ok := ports[core.PortAddr(addr)]
		if !ok {
			return ErrPortNotFound
		}

		var msg core.Msg
		switch messages[addr].Type {
		case runtime.IntMsg:
			msg = core.NewIntMsg(messages[addr].Int)
		case runtime.BoolMsg:
			msg = core.NewBoolMsg(messages[addr].Bool)
		case runtime.StrMsg:
			msg = core.NewStrMsg(messages[addr].Str)
		case runtime.SigMsg:
			msg = core.NewSigMsg()
		}

		go func() {
			for {
				port <- msg
			}
		}()
	}

	return nil
}

package constspawner

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

type Spawner struct{}

var (
	ErrPortNotFound   = errors.New("port not found")
	ErrUnknownMsgType = errors.New("unknown message type")
)

func (c Spawner) Spawn(
	constData map[runtime.PortAddr]runtime.ConstMsg,
	ports map[runtime.PortAddr]chan core.Msg,
) error {
	for addr := range constData {
		port, ok := ports[addr]
		if !ok {
			return fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		var msg core.Msg
		switch constData[addr].Type {
		case runtime.IntMsg:
			msg = core.NewIntMsg(constData[addr].Int)
		case runtime.BoolMsg:
			msg = core.NewBoolMsg(constData[addr].Bool)
		case runtime.StrMsg:
			msg = core.NewStrMsg(constData[addr].Str)
		case runtime.SigMsg:
			msg = core.NewSigMsg()
		default:
			return fmt.Errorf("%w: %v", ErrUnknownMsgType, constData[addr].Type)
		}

		go func() {
			for {
				port <- msg
			}
		}()
	}

	return nil
}

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

func (s Spawner) Spawn(
	outportsAndValues map[runtime.PortAddr]runtime.Msg,
	outportsAndChans map[runtime.PortAddr]chan core.Msg,
) error {
	for addr := range outportsAndValues {
		port, ok := outportsAndChans[addr]
		if !ok {
			return fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		msg, err := s.coreMsg(outportsAndValues, addr)
		if err != nil {
			return fmt.Errorf("core msg: %w", err)
		}

		go func() {
			for {
				port <- msg
			}
		}()
	}

	return nil
}

func (s Spawner) coreMsg(outportsAndValues map[runtime.PortAddr]runtime.Msg, addr runtime.PortAddr) (core.Msg, error) {
	var msg core.Msg
	switch outportsAndValues[addr].Type {
	case runtime.IntMsg:
		msg = core.NewIntMsg(outportsAndValues[addr].Int)
	case runtime.BoolMsg:
		msg = core.NewBoolMsg(outportsAndValues[addr].Bool)
	case runtime.StrMsg:
		msg = core.NewStrMsg(outportsAndValues[addr].Str)
	case runtime.StructMsg:
		msg = core.NewStructMsg(
			s.structMsg(outportsAndValues[addr].Struct),
		)
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnknownMsgType, outportsAndValues[addr].Type)
	}
	return msg, nil
}

func (s Spawner) structMsg(map[string]runtime.Msg) map[string]core.Msg {
	return map[string]core.Msg{} // TODO
}

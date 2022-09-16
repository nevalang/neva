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
	outportsAndValues map[runtime.AbsolutePortAddr]runtime.Msg,
	outportsAndChans map[runtime.AbsolutePortAddr]chan core.Msg,
) error {
	for addr := range outportsAndValues {
		out, ok := outportsAndChans[addr]
		if !ok {
			return fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		msg, err := s.coreMsg(outportsAndValues[addr])
		if err != nil {
			return fmt.Errorf("core msg: %w", err)
		}

		go func() {
			for {
				out <- msg
			}
		}()
	}

	return nil
}

func (s Spawner) coreMsg(in runtime.Msg) (core.Msg, error) {
	var out core.Msg

	switch in.Type {
	case runtime.IntMsg:
		out = core.NewIntMsg(in.Int)
	case runtime.BoolMsg:
		out = core.NewBoolMsg(in.Bool)
	case runtime.StrMsg:
		out = core.NewStrMsg(in.Str)
	case runtime.StructMsg:
		structMsg := make(map[string]core.Msg, len(in.Struct))
		for field, value := range in.Struct {
			v, err := s.coreMsg(value)
			if err != nil {
				return nil, fmt.Errorf("core msg: %w", err)
			}
			structMsg[field] = v
		}
		out = core.NewStructMsg(structMsg)
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnknownMsgType, in.Type)
	}

	return out, nil
}

func (s Spawner) structMsg(in map[string]runtime.Msg) core.StructMsg {
	out := make(map[string]core.Msg, len(in))

	for field, value := range in {
		switch value.Type {
		case runtime.BoolMsg:
			out[field] = core.NewBoolMsg(value.Bool)
		case runtime.IntMsg:
			out[field] = s.structMsg(value.Struct)
		case runtime.StrMsg:
			out[field] = s.structMsg(value.Struct)
		case runtime.StructMsg:
			out[field] = s.structMsg(value.Struct)
		}
	}

	return core.NewStructMsg(out)
}

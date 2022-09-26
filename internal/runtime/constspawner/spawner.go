package constspawner

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime/src"
)

type Spawner struct{}

var (
	ErrPortNotFound   = errors.New("port not found")
	ErrUnknownMsgType = errors.New("unknown message type")
)

func (s Spawner) Spawn(
	ctx context.Context,
	messages map[src.AbsolutePortAddr]src.Msg,
	chans map[src.AbsolutePortAddr]chan core.Msg,
) error {
	for addr := range messages {
		out, ok := chans[addr]
		if !ok {
			return fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		msg, err := s.coreMsg(messages[addr])
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

func (s Spawner) coreMsg(in src.Msg) (core.Msg, error) {
	var out core.Msg

	switch in.Type {
	case src.IntMsg:
		out = core.NewIntMsg(in.Int)
	case src.BoolMsg:
		out = core.NewBoolMsg(in.Bool)
	case src.StrMsg:
		out = core.NewStrMsg(in.Str)
	case src.StructMsg:
		structMsg := make(map[string]core.Msg, len(in.Struct))
		for field, value := range in.Struct {
			v, err := s.coreMsg(value)
			if err != nil {
				return nil, fmt.Errorf("core msg: %w", err)
			}
			structMsg[field] = v
		}
		out = core.NewDictMsg(structMsg)
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnknownMsgType, in.Type)
	}

	return out, nil
}

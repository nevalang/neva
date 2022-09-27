package constspawner

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/src"
)

type Spawner struct{}

var (
	ErrPortNotFound   = errors.New("port not found")
	ErrUnknownMsgType = errors.New("unknown message type")
)

func (s Spawner) Spawn(ctx context.Context, consts []runtime.Const) error {
	for i := range consts {
		cnst := consts[i]

		msg, err := s.coreMsg(cnst.Msg)
		if err != nil {
			return fmt.Errorf("%w: %d", ErrUnknownMsgType, i)
		}

		go func() {
			for {
				cnst.Port <- msg
			}
		}()
	}

	return nil // TODO wait
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

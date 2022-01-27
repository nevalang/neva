package constspawner

import (
	"errors"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type Spawner struct{}

var ErrPortNotFound = errors.New("port not found")

func (c Spawner) Spawn(values map[string]runtime.ConstValue, ports map[core.PortAddr]chan core.Msg) error {
	for name := range values {
		port, ok := ports[core.PortAddr{Port: name}]
		if !ok {
			return ErrPortNotFound
		}

		var msg core.Msg
		switch values[name].Type {
		case runtime.IntMsg:
			msg = core.NewIntMsg(values[name].Int)
		case runtime.BoolMsg:
			msg = core.NewBoolMsg(values[name].Bool)
		case runtime.StrMsg:
			msg = core.NewStrMsg(values[name].Str)
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

package constspawner

import (
	"errors"

	"github.com/emil14/neva/internal/new/core"
)

type ConstSpawner struct{}

var ErrPortNotFound = errors.New("port not found")

func (c ConstSpawner) Spawn(data map[string]core.Msg, ports map[core.PortAddr]chan core.Msg) error {
	for name := range data {
		port, ok := ports[core.PortAddr{Port: name}]
		if !ok {
			return ErrPortNotFound
		}

		msg := data[name]
		go func() {
			for {
				port <- msg
			}
		}()
	}

	return nil
}

package loader

import (
	"errors"
	"log"
	"plugin"

	"github.com/emil14/neva/internal/runtime"
)

type Params struct {
	PluginPath     string
	ExportedEntity string
}

func Load(paths map[string]Params) (map[string]runtime.Operator, error) {
	ops := map[string]runtime.Operator{}

	for name, params := range paths {
		plgn, err := plugin.Open(params.PluginPath)
		if err != nil {
			return nil, err
		}

		sym, err := plgn.Lookup(params.ExportedEntity)
		if err != nil {
			return nil, err
		}

		op, ok := sym.(func(runtime.IO) error)
		if !ok {
			log.Printf("NOT OK: %T", sym)
			return nil, errors.New("not ok from loader")
		}

		ops[name] = op
	}

	return ops, nil
}

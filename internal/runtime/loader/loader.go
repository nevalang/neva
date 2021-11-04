package loader

import (
	"fmt"
	"plugin"

	"github.com/emil14/respect/internal/core"
	"github.com/emil14/respect/internal/runtime"
)

type LoadParams struct {
	Path   string
	Export []string
}

func Load(params map[string]LoadParams) (map[string]runtime.OperatorFunc, error) {
	ops := map[string]runtime.OperatorFunc{}

	for name, p := range params {
		plug, err := plugin.Open(p.Path)
		if err != nil {
			return nil, fmt.Errorf("plugin open: %w", err)
		}

		for _, export := range p.Export {
			symb, err := plug.Lookup(export)
			if err != nil {
				return nil, fmt.Errorf("plugin lookup: %w", err)
			}

			op, ok := symb.(func(core.IO) error)
			if !ok {
				return nil, fmt.Errorf("export not operator: %T", op)
			}

			ops[name] = op
		}

	}

	return ops, nil
}

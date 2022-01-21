package opspawner

import (
	"errors"
	"fmt"
	"log"
	"plugin"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type Plugin struct {
	cache   cache
	plugins map[string]PluginData
}

func (p Plugin) Operator(ref runtime.OpRef) (func(core.IO) error, error) {
	if op := p.cache.get(ref.Pkg, ref.Name); op != nil {
		return op, nil
	}

	pluginData, ok := p.plugins[ref.Pkg]
	if !ok {
		return nil, fmt.Errorf("unknown package %s", ref.Pkg)
	}

	plug, err := plugin.Open(pluginData.path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, pluginData.path)
	}

	for _, export := range pluginData.exports {
		sym, err := plug.Lookup(export)
		if err != nil {
			return nil, err
		}

		op, ok := sym.(func(core.IO) error)
		if !ok {
			log.Printf("NOT OK: %T", sym)
			return nil, errors.New("not ok from loader")
		}

		p.cache.set(ref.Pkg, ref.Name, op)
	}

	opfunc := p.cache.get(ref.Pkg, ref.Name)
	if opfunc == nil {
		return nil, fmt.Errorf("package %s has no operator %s", ref.Pkg, ref.Name)
	}

	return opfunc, nil
}

type PluginData struct {
	path    string
	exports []string
}

func NewPlugin(plugins map[string]PluginData) Plugin {
	return Plugin{
		plugins: plugins,
		cache:   make(cache, len(plugins)),
	}
}

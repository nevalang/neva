package repo

import (
	"errors"
	"fmt"
	"plugin"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

var (
	ErrUnknownPkg   = errors.New("operator refers to unknown package")
	ErrPluginOpen   = errors.New("plugin could not be loaded")
	ErrPluginLookup = errors.New("exported entity not found")
	ErrTypeMismatch = errors.New("exported entity doesn't match operator signature")
	ErrOpNotFound   = errors.New("package has not implemented the operator")
)

type Plugin struct {
	pkgs  map[string]PluginData
	cache map[runtime.OpRef]func(core.IO) error
}

func (r Plugin) Operator(ref runtime.OpRef) (func(core.IO) error, error) {
	if op, ok := r.cache[ref]; ok {
		return op, nil
	}

	pluginData, ok := r.pkgs[ref.Pkg]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownPkg, ref.Pkg)
	}

	plug, err := plugin.Open(pluginData.Path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPluginOpen, err)
	}

	for _, export := range pluginData.Exports {
		sym, err := plug.Lookup(export)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPluginLookup, err)
		}

		op, ok := sym.(func(core.IO) error)
		if !ok {
			return nil, fmt.Errorf("%w: %T", ErrTypeMismatch, op)
		}

		r.cache[ref] = op
	}

	op, ok := r.cache[ref]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrOpNotFound, ref)
	}

	return op, nil
}

type PluginData struct {
	Path    string
	Exports []string
}

func NewPlugin(pkgs map[string]PluginData) Plugin {
	return Plugin{
		pkgs: pkgs,
		cache: make(
			map[runtime.OpRef]func(core.IO) error,
			len(pkgs),
		),
	}
}

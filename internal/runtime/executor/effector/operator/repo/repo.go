package repo

import (
	"context"
	"errors"
	"fmt"
	"plugin"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime/executor/effector/operator"
	"github.com/emil14/neva/internal/runtime/src"
)

var (
	ErrUnknownPkg   = errors.New("operator refers to unknown package")
	ErrPluginOpen   = errors.New("plugin could not be loaded")
	ErrPluginLookup = errors.New("exported entity not found")
	ErrTypeMismatch = errors.New("exported entity doesn't match operator signature")
	ErrOpNotFound   = errors.New("operator not found")
)

type Plugin struct {
	packages map[string]File
	cache    map[src.OperatorRef]operator.Func // move cache to operator effector?
}

func (p Plugin) Operator(ref src.OperatorRef) (operator.Func, error) {
	if f, ok := p.cache[ref]; ok {
		return f, nil
	}

	file, ok := p.packages[ref.Pkg]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownPkg, ref.Pkg)
	}

	plug, err := plugin.Open(file.Path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPluginOpen, err)
	}

	for _, export := range file.Exports {
		sym, err := plug.Lookup(export)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPluginLookup, err)
		}

		f, ok := sym.(func(context.Context, core.IO) error)
		if !ok {
			return nil, fmt.Errorf("%w: %T", ErrTypeMismatch, f)
		}

		p.cache[src.OperatorRef{
			Pkg:  ref.Pkg,
			Name: export,
		}] = f
	}

	op, ok := p.cache[ref]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrOpNotFound, ref)
	}

	return op, nil
}

type File struct {
	Path    string
	Exports []string
}

func NewPlugin(files map[string]File) Plugin {
	return Plugin{
		packages: files,
		cache:    map[src.OperatorRef]operator.Func{}, // calc size? (files*exports)
	}
}

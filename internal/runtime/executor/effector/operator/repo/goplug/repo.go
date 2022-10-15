package goplug

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

type Repo struct {
	packages map[string]File
	cache    map[src.FuncRef]operator.FuncFx // move cache to operator effector?
}

func (p Repo) Func(ref src.FuncRef) (operator.FuncFx, error) {
	if f, ok := p.cache[ref]; ok {
		return f, nil
	}

	file, ok := p.packages[ref.Class]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownPkg, ref.Class)
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

		p.cache[src.FuncRef{
			Class: ref.Class,
			Name:  export,
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

func MustNewRepo(files map[string]File) Repo {
	if len(files) < 1 {
		panic(errors.New("files < 1"))
	}
	return Repo{
		packages: files,
		cache:    map[src.FuncRef]operator.FuncFx{}, // calc size? (files*exports)
	}
}

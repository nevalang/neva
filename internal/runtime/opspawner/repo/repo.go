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
	ErrOpNotFound   = errors.New("operator not found")
)

type Plugin struct {
	pkgs  map[string]Package
	cache map[runtime.OperatorRef]func(core.IO) error
}

func (p Plugin) Operator(ref runtime.OperatorRef) (func(core.IO) error, error) {
	if op, ok := p.cache[ref]; ok {
		return op, nil
	}

	pkg, ok := p.pkgs[ref.Pkg]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownPkg, ref.Pkg)
	}

	plug, err := plugin.Open(pkg.Filepath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPluginOpen, err)
	}

	for _, export := range pkg.Exports {
		sym, err := plug.Lookup(export)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPluginLookup, err)
		}

		op, ok := sym.(func(core.IO) error)
		if !ok {
			return nil, fmt.Errorf("%w: %T", ErrTypeMismatch, op)
		}

		p.cache[runtime.OperatorRef{
			Pkg:  ref.Pkg,
			Name: export,
		}] = op
	}

	op, ok := p.cache[ref]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrOpNotFound, ref)
	}

	return op, nil
}

type Package struct {
	Filepath string
	Exports  []string
}

func NewPlugin(pkgs map[string]Package) Plugin {
	return Plugin{
		pkgs: pkgs,
		cache: make(
			map[runtime.OperatorRef]func(core.IO) error,
			len(pkgs),
		),
	}
}

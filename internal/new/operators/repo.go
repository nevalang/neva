package operators

import (
	"errors"
	"fmt"
	"log"
	"plugin"

	"github.com/emil14/neva/internal/runtime"
)

type Repo struct {
	cache    cache
	packages map[string]PluginData
}

func (r Repo) Operator(pkg, name string) (runtime.Opfunc, error) {
	if opfunc := r.cache.get(pkg, name); opfunc != nil {
		return opfunc, nil
	}

	pluginData, ok := r.packages[pkg]
	if !ok {
		return nil, fmt.Errorf("unknown package %s", pkg)
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

		op, ok := sym.(func(runtime.IO) error)
		if !ok {
			log.Printf("NOT OK: %T", sym)
			return nil, errors.New("not ok from loader")
		}

		r.cache.set(pkg, name, op)
	}

	opfunc := r.cache.get(pkg, name)
	if opfunc == nil {
		return nil, fmt.Errorf("package %s has no operator %s", pkg, name)
	}

	return opfunc, nil
}

type PluginData struct {
	path    string
	exports []string
}

type cache map[string]map[string]runtime.Opfunc

func (c cache) get(pkg, name string) runtime.Opfunc {
	p, ok := c[pkg]
	if !ok {
		return nil
	}
	return p[name]
}

func (c cache) set(pkg, name string, opfunc runtime.Opfunc) {
	if c[pkg] == nil {
		c[pkg] = map[string]runtime.Opfunc{}
	}
	c[pkg][name] = opfunc
}

func NewRepo(packages map[string]PluginData) Repo {
	return Repo{
		packages: packages,
		cache:    make(cache, len(packages)),
	}
}

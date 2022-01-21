package opspawner

import "github.com/emil14/neva/internal/new/core"

type cache map[string]map[string]func(core.IO) error

func (c cache) get(pkg, name string) func(core.IO) error {
	p, ok := c[pkg]
	if !ok {
		return nil
	}
	return p[name]
}

func (c cache) set(pkg, name string, opfunc func(core.IO) error) {
	if c[pkg] == nil {
		c[pkg] = map[string]func(core.IO) error{}
	}
	c[pkg][name] = opfunc
}

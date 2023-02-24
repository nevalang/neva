package types

import "errors"

var ErrDefaultScope = errors.New("default scope")

// DefaultScope implements Scope interface for internal usage
type DefaultScope map[string]Def

func (d DefaultScope) Get(ref string) (Def, error) {
	v, ok := d[ref]
	if !ok {
		return Def{}, ErrDefaultScope
	}
	return v, nil
}

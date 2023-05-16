package types_test

import (
	"errors"

	ts "github.com/nevalang/neva/pkg/types"
)

var ErrDefaultScope = errors.New("default scope")

type Scope map[string]ts.Def

func (d Scope) GetType(ref string) (ts.Def, error) {
	v, ok := d[ref]
	if !ok {
		return ts.Def{}, ErrDefaultScope
	}
	return v, nil
}

func (d Scope) Update(string) (ts.Scope, error) {
	return d, nil
}

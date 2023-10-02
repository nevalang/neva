// This file implements Scope interface specifically for tests.
package typesystem_test

import (
	"errors"

	ts "github.com/nevalang/neva/pkg/typesystem"
)

var ErrDefaultScope = errors.New("default scope")

type Scope map[string]ts.Def

func (s Scope) GetType(ref string) (ts.Def, ts.Scope, error) {
	v, ok := s[ref]
	if !ok {
		return ts.Def{}, nil, ErrDefaultScope
	}
	return v, s, nil
}

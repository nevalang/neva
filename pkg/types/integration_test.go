//go:build integration
// +build integration

package types_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/stretchr/testify/assert"
)

func TestDefaultResolver(t *testing.T) {
	r := ts.NewDefaultResolver()

	got, err := r.Resolve(
		h.Inst("vec", h.Inst("t1")),
		map[string]ts.Def{
			"vec": h.BaseDefWithRecursion(h.ParamWithoutConstr("t")),
			"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
			"t2":  h.Def(h.Inst("vec", h.Inst("t2"))),
		},
	)

	assert.Equal(t, h.Inst("vec", h.Inst("t1")), got) // FIXME we have vec<vec<t1>>
	assert.Equal(t, nil, err)
}

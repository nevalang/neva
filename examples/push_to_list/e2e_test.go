package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out := e2e.RunExample(t, "push_to_list")

	require.Equal(
		t,
		"[320,420,100,-100,0,5,69]\n",
		out,
	)
}

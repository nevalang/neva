package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 10; i++ {
		out := e2e.RunExample(t, "match")
		require.Equal(
			t,
			"one\ntwo\nthree\nfour\n",
			out,
		)
	}
}

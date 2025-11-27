package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, err := e2e.RunExpectingError(t, "run", "main")
	require.Contains(
		t,
		out+err,
		"main/main.neva:4:1: array outport 'fanOut:data' is used incorrectly: slot 1 is missing\n",
	)
}

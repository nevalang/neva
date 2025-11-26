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
		"main/main.neva:9:1: port 'println:res' is used twice\n",
	)
}

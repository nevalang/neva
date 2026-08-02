package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Test verifies the public wire representation of a finite descriptor with a
// recursive back-edge. The value must pass through a generic component without
// changing its graph indexes before Println renders it.
func Test(t *testing.T) {
	const want = `[` +
		`{"tag":"Struct","data":[{"name":"text","node":1},{"name":"child","node":2}]},` +
		`{"tag":"String"},` +
		`{"tag":"Union","data":[` +
		`{"data":{"tag":"Some","data":0},"tag":"Some"},` +
		`{"data":{"tag":"None"},"tag":"None"}` +
		`]}` +
		`]` + "\n"

	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Equal(t, want, out)
}

package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	// TODO simplify when struct fields will have order: https://github.com/nevalang/neva/issues/698

	expectedErrorPrefix := "main/main.neva:16:8 Incompatible types: in:in -> add: Subtype and supertype must both be either literals or instances, except if supertype is union: expression any, constraint {"
	expectedErrorSuffix := "}\n"
	expectedFields := []string{"data int", "idx int", "last bool"}

	actualOutput := string(out)
	require.True(
		t,
		strings.HasPrefix(actualOutput, expectedErrorPrefix),
		"Error message should start with expected prefix",
	)
	require.True(
		t,
		strings.HasSuffix(actualOutput, expectedErrorSuffix),
		"Error message should end with expected suffix")

	// Check if all expected fields are present in the error message
	for _, field := range expectedFields {
		require.Contains(t, actualOutput, field, "Error message should contain field: "+field)
	}

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

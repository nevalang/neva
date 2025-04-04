package test

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestEmitDefault(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.yml"))
	}()

	// Create new project
	cmd := exec.Command("neva", "new")
	require.NoError(t, cmd.Run())

	// Run with IR emission
	cmd = exec.Command("neva", "run", "--emit-ir", "src")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, "Hello, World!\n", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// Verify IR file exists and is valid YAML
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err, string(out))

	var ir struct {
		Connections map[string]string `yaml:"connections"`
		Funcs       []any             `yaml:"funcs"`
	}
	require.NoError(t, yaml.Unmarshal(irBytes, &ir))
	require.NotEmpty(t, ir.Funcs)
}

func TestEmitYAML(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.yml"))
	}()

	// Create new project
	cmd := exec.Command("neva", "new")
	require.NoError(t, cmd.Run())

	// Run with IR emission
	cmd = exec.Command("neva", "run", "--emit-ir", "--emit-ir-format", "yaml", "src")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, "Hello, World!\n", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// Verify IR file exists and is valid YAML
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err, string(out))

	var ir struct {
		Connections map[string]string `yaml:"connections"`
		Funcs       []any             `yaml:"funcs"`
	}
	require.NoError(t, yaml.Unmarshal(irBytes, &ir))
	require.NotEmpty(t, ir.Funcs)
}

func TestEmitJSON(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.json"))
	}()

	// Create new project
	cmd := exec.Command("neva", "new")
	require.NoError(t, cmd.Run())

	// Run with IR emission
	cmd = exec.Command("neva", "run", "--emit-ir", "--emit-ir-format", "json", "src")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command failed: %v", string(out))
	require.Equal(t, "Hello, World!\n", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// Verify IR file exists and is valid JSON
	irBytes, err := os.ReadFile("ir.json")
	require.NoError(t, err, string(out))

	var ir struct {
		Connections map[string]string `json:"connections"`
		Funcs       []any             `json:"funcs"`
	}
	require.NoError(t, json.Unmarshal(irBytes, &ir))
	require.NotEmpty(t, ir.Funcs)
}

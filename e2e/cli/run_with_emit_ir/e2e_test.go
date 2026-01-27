package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestEmitDefault(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.yml"))
	}()

	// Create new project
	e2e.Run(t, []string{"new", "."})

	// Run with IR emission
	out, _ := e2e.Run(t, []string{"run", "--emit-ir", "src"})
	require.Equal(t, "Hello, World!\n", out)

	// Verify IR file exists and is valid YAML
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err, out)

	var ir struct {
		Connections []struct {
			From string `yaml:"from"`
			To   string `yaml:"to"`
		} `yaml:"connections"`
		Funcs []struct {
			Ref string `yaml:"ref"`
			IO  struct {
				In  []string `yaml:"in"`
				Out []string `yaml:"out"`
			} `yaml:"io"`
			Msg map[string]any `yaml:"msg,omitempty"`
		} `yaml:"funcs"`
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
	e2e.Run(t, []string{"new", "."})

	// Run with IR emission
	out, _ := e2e.Run(t, []string{"run", "--emit-ir", "--emit-ir-format", "yaml", "src"})
	require.Equal(t, "Hello, World!\n", out)

	// Verify IR file exists and is valid YAML
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err, out)

	var ir struct {
		Connections []struct {
			From string `yaml:"from"`
			To   string `yaml:"to"`
		} `yaml:"connections"`
		Funcs []struct {
			Ref string `yaml:"ref"`
			IO  struct {
				In  []string `yaml:"in"`
				Out []string `yaml:"out"`
			} `yaml:"io"`
			Msg map[string]any `yaml:"msg,omitempty"`
		} `yaml:"funcs"`
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
	e2e.Run(t, []string{"new", "."})

	// Run with IR emission
	out, _ := e2e.Run(t, []string{"run", "--emit-ir", "--emit-ir-format", "json", "src"})
	require.Equal(t, "Hello, World!\n", out)

	// Verify IR file exists and is valid JSON
	irBytes, err := os.ReadFile("ir.json")
	require.NoError(t, err, out)

	var ir struct {
		Connections map[string]string `json:"connections"`
		Funcs       []any             `json:"funcs"`
	}
	require.NoError(t, json.Unmarshal(irBytes, &ir))
	require.NotEmpty(t, ir.Funcs)
}

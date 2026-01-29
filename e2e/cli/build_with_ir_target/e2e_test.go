package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestBuildIRDefault(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.yml"))
	}()

	// create new project
	e2e.Run(t, []string{"new", "."})

	// build with ir target (default yaml format)
	out, _ := e2e.Run(t, []string{"build", "-target=ir", "src"})

	// verify ir file exists and is valid yaml
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err, out)

	var ir struct {
		Connections []struct {
			From string `yaml:"from"`
			To   string `yaml:"to"`
		} `yaml:"connections"`
		Funcs []struct { //nolint:govet // fieldalignment
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

func TestBuildIRYAML(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.yml"))
	}()

	// create new project
	e2e.Run(t, []string{"new", "."})

	// build with ir target and explicit yaml format
	out, _ := e2e.Run(t, []string{"build", "-target=ir", "-target-ir-format=yaml", "src"})

	// verify ir file exists and is valid yaml
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err, out)

	var ir struct {
		Connections []struct {
			From string `yaml:"from"`
			To   string `yaml:"to"`
		} `yaml:"connections"`
		Funcs []struct { //nolint:govet // fieldalignment
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

func TestBuildIRJSON(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.json"))
	}()

	// create new project
	e2e.Run(t, []string{"new", "."})

	// build with ir target and json format
	out, _ := e2e.Run(t, []string{"build", "-target=ir", "-target-ir-format=json", "src"})

	// verify ir file exists and is valid json
	irBytes, err := os.ReadFile("ir.json")
	require.NoError(t, err, out)

	var ir struct {
		Connections map[string]string `json:"connections"`
		Funcs       []any             `json:"funcs"`
	}
	require.NoError(t, json.Unmarshal(irBytes, &ir))
	require.NotEmpty(t, ir.Funcs)
}

func TestBuildIRWithInterfaceWithImports(t *testing.T) {
	defer func() {
		os.Remove("ir.yml")
		os.Remove("ir.json")
	}()

	// build the specific case that was failing in the original issue
	// this tests the interface_with_imports case with ir target
	out, _ := e2e.Run(t, []string{"build", "-target=ir", "../../interface_with_imports/main"})

	// verify ir file exists and is valid yaml
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err, out)

	var ir struct {
		Connections []struct {
			From string `yaml:"from"`
			To   string `yaml:"to"`
		} `yaml:"connections"`
		Funcs []struct { //nolint:govet // fieldalignment
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

func TestBuildIRWithInterfaceWithImportsJSON(t *testing.T) {
	defer func() {
		require.NoError(t, os.Remove("ir.json"))
	}()

	// build the specific case that was failing in the original issue with json format
	out, _ := e2e.Run(t, []string{"build", "-target=ir", "-target-ir-format=json", "../../interface_with_imports/main"})

	// verify ir file exists and is valid json
	irBytes, err := os.ReadFile("ir.json")
	require.NoError(t, err, out)

	var ir struct {
		Connections map[string]string `json:"connections"`
		Funcs       []any             `json:"funcs"`
	}
	require.NoError(t, json.Unmarshal(irBytes, &ir))
	require.NotEmpty(t, ir.Funcs)
}

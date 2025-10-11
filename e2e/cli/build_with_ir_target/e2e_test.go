package test

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestBuildIRDefault(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.yml"))
	}()

	// create new project
	cmd := exec.Command("neva", "new")
	require.NoError(t, cmd.Run())

	// build with ir target (default yaml format)
	cmd = exec.Command("neva", "build", "-target=ir", "src")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// verify ir file exists and is valid yaml
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err, string(out))

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

func TestBuildIRYAML(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.yml"))
	}()

	// create new project
	cmd := exec.Command("neva", "new")
	require.NoError(t, cmd.Run())

	// build with ir target and explicit yaml format
	cmd = exec.Command("neva", "build", "-target=ir", "-target-ir-format=yaml", "src")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// verify ir file exists and is valid yaml
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err, string(out))

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

func TestBuildIRJSON(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.json"))
	}()

	// create new project
	cmd := exec.Command("neva", "new")
	require.NoError(t, cmd.Run())

	// build with ir target and json format
	cmd = exec.Command("neva", "build", "-target=ir", "-target-ir-format=json", "src")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command failed: %v", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// verify ir file exists and is valid json
	irBytes, err := os.ReadFile("ir.json")
	require.NoError(t, err, string(out))

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
	cmd := exec.Command("neva", "build", "-target=ir", "../../interface_with_imports/main")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command failed: %v", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// verify ir file exists and is valid yaml
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err, string(out))

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

func TestBuildIRWithInterfaceWithImportsJSON(t *testing.T) {
	defer func() {
		require.NoError(t, os.Remove("ir.json"))
	}()

	// build the specific case that was failing in the original issue with json format
	cmd := exec.Command("neva", "build", "-target=ir", "-target-ir-format=json", "../../interface_with_imports/main")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command failed: %v", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// verify ir file exists and is valid json
	irBytes, err := os.ReadFile("ir.json")
	require.NoError(t, err, string(out))

	var ir struct {
		Connections map[string]string `json:"connections"`
		Funcs       []any             `json:"funcs"`
	}
	require.NoError(t, json.Unmarshal(irBytes, &ir))
	require.NotEmpty(t, ir.Funcs)
}

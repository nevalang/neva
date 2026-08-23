package test

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func Test(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Empty(t, out)
}

func TestBuildIRBindsResolvedTypeAsConfigMessage(t *testing.T) {
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll("ir.yml"))
	})

	e2e.Run(t, []string{"build", "-target=ir", "main"})
	irBytes, err := os.ReadFile("ir.yml")
	require.NoError(t, err)

	var program struct {
		Funcs []struct {
			Ref string `yaml:"ref"`
			Msg struct {
				List []struct {
					Union struct {
						Tag string `yaml:"tag"`
					} `yaml:"union"`
				} `yaml:"list"`
			} `yaml:"msg"`
		} `yaml:"funcs"`
	}
	require.NoError(t, yaml.Unmarshal(irBytes, &program))

	for _, call := range program.Funcs {
		if call.Ref != "int_inc" {
			continue
		}
		require.Len(t, call.Msg.List, 1)
		require.Equal(t, "Int", call.Msg.List[0].Union.Tag)
		return
	}
	t.Fatal("int_inc runtime function call not found")
}

package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// this e2e covers go package mode end-to-end using a real neva module: generate, build, and call an export.
func TestGoPkgMode_EndToEnd(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("gen")
		_ = os.RemoveAll("src")
		_ = os.Remove("neva.yml")
		_ = os.Remove("neva.yaml")
	})

	// create a real neva module with manifest and src/ via cli
	cmd := exec.Command("neva", "new", ".")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	// add exported Print42(sig any) (sig any) that prints 42 and panics on error
	nevaSrc := `import { fmt, runtime }

pub def Print42(sig any) (sig any) {
    println fmt.Println<int>, panic runtime.Panic
    ---
    42 -> println:data
    println:res -> :sig
    println:err -> panic
}`
	require.NoError(t, os.WriteFile("src/print42.neva", []byte(nevaSrc), 0o644))

	// build the go package into gen/print42
	outDir := filepath.Join("gen", "print42")
	cmd = exec.Command("neva", "build", "--target=go", "--target-go-mode=pkg", "--output="+outDir, "src")
	out, err = cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	// verify expected generated files exist
	_, err = os.Stat(filepath.Join(outDir, "api.go"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(outDir, "programs.go"))
	require.NoError(t, err)

	// initialize a separate go module in gen and point neva imports to the local repo
	goModInit := exec.Command("go", "mod", "init", "example.com/tmpgen")
	goModInit.Dir = "gen"
	out, err = goModInit.CombinedOutput()
	require.NoError(t, err, string(out))

	// ensure generated imports of github.com/nevalang/neva resolve to the repo root
	goModReplace := exec.Command("go", "mod", "edit", "-replace", "github.com/nevalang/neva=../")
	goModReplace.Dir = "gen"
	out, err = goModReplace.CombinedOutput()
	require.NoError(t, err, string(out))

	// build the generated go package to ensure it compiles
	goBuild := exec.Command("go", "build", "./...")
	goBuild.Dir = "gen"
	out, err = goBuild.CombinedOutput()
	require.NoError(t, err, string(out))

	// compile and run a tiny program that imports the generated package and calls Print42
	runner := `package main
import (
    "context"
    gen "example.com/tmpgen/print42"
    "github.com/nevalang/neva/internal/runtime"
)
func main(){
    ctx := context.Background()
    _, _ = gen.Print42(ctx, runtime.NewStructMsg([]runtime.StructField{runtime.NewStructField("sig", runtime.NewIntMsg(0))}))
}`
	require.NoError(t, os.WriteFile(filepath.Join("gen", "main.go"), []byte(runner), 0o644))
	run := exec.Command("go", "run", ".")
	run.Dir = "gen"
	out, err = run.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Contains(t, string(out), "42")
}

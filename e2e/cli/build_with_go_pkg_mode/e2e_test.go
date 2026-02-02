package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// this e2e covers go package mode end-to-end using a real neva module: generate, build, and call an export.
func TestGoPkgMode_EndToEnd(t *testing.T) {
	// Clean start
	_ = os.RemoveAll("gen")
	_ = os.RemoveAll("src")
	_ = os.Remove("neva.yml")
	_ = os.Remove("neva.yaml")

	t.Cleanup(func() {
		_ = os.RemoveAll("gen")
		_ = os.RemoveAll("src")
		_ = os.Remove("neva.yml")
		_ = os.Remove("neva.yaml")
	})

	// create neva module
	e2e.Run(t, []string{"new", "."})

	// add exported PrintHello(sig int) (res string) that prints "Hello from Neva!" and panics on error
	nevaSrc := `import {
	fmt
	runtime
}

pub def PrintHello(sig int) (res string) {
    println fmt.Println<string>
    panic runtime.Panic
    ---
    :sig -> 'Hello from Neva!' -> println:data
    println:res -> :res
    println:err -> panic
}`
	// #nosec G306 -- test fixture file
	require.NoError(t, os.WriteFile("src/hello.neva", []byte(nevaSrc), 0o644))

	// go mod init in "gen" dir BEFORE build to let compiler detect module path
	require.NoError(t, os.Mkdir("gen", 0o755))
	goModInit := exec.Command("go", "mod", "init", "example.com/tmpgen")
	goModInit.Dir = "gen"
	out, err := goModInit.CombinedOutput()
	require.NoError(t, err, string(out))

	// neva build with pkg mode into "gen/hello"
	outDir := filepath.Join("gen", "hello")
	e2e.Run(t, []string{"build", "--target=go", "--target-go-mode=pkg", "--output=" + outDir, "src"})

	// No need for go mod edit -replace or go get if runtime is self-contained and imports are correct!

	// compile and run a tiny program that imports the generated package and calls PrintHello
	runner := `package main
import (
    "context"
    "fmt"
    gen "example.com/tmpgen/hello"
)
func main(){
    ctx := context.Background()
    out, err := gen.PrintHello(ctx, gen.PrintHelloInput{Sig: 0})
    if err != nil {
        panic(err)
    }
    fmt.Printf("%v", out.Res)
}`
	// #nosec G306 -- test fixture file
	require.NoError(t, os.WriteFile(filepath.Join("gen", "main.go"), []byte(runner), 0o644))

	// Just go run .
	run := exec.Command("go", "run", ".")
	run.Dir = "gen"
	out, err = run.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Contains(t, string(out), "Hello from Neva!")
}

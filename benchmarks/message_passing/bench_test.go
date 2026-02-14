package test

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
)

type runtimeBenchCase struct {
	name string
	pkg  string
}

// BenchmarkRuntimeE2E benchmarks precompiled runtime programs by data type path.
func BenchmarkRuntimeE2E(b *testing.B) {
	// Build the CLI once and reuse it for all benchmark programs.
	repoRoot := e2e.FindRepoRoot(b)
	nevaBin := e2e.BuildNevaBinary(b, repoRoot)

	cases := []runtimeBenchCase{
		{name: "int_message_passing", pkg: "message_passing"},
		{name: "bool_map", pkg: "runtime_bool"},
		{name: "string_map", pkg: "runtime_string"},
		{name: "float_parse_add", pkg: "runtime_float"},
		{name: "list_roundtrip", pkg: "runtime_list"},
		{name: "dict_lookup", pkg: "runtime_dict"},
		{name: "struct_build_select", pkg: "runtime_struct"},
		{name: "union_wrap_unwrap", pkg: "runtime_union"},
		{name: "combo_struct_union", pkg: "runtime_combo"},
	}

	for _, benchCase := range cases {
		b.Run(benchCase.name, func(b *testing.B) {
			// Build the benchmark program once outside timed iterations.
			progPath := buildProgramOnce(b, repoRoot, nevaBin, benchCase.pkg)

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				runProgramBinary(b, progPath)
			}
		})
	}
}

// buildProgramOnce prepares an isolated module and compiles one benchmark package.
func buildProgramOnce(b *testing.B, repoRoot, nevaBin, pkgName string) string {
	b.Helper()

	// Create an isolated temp module workspace for one benchmark package.
	tmpDir := b.TempDir()
	moduleDir := filepath.Join(tmpDir, "bench-module")
	progDir := filepath.Join(moduleDir, pkgName)
	if err := os.MkdirAll(progDir, 0o755); err != nil {
		b.Fatalf("create benchmark module dirs: %v", err)
	}

	// Copy benchmark fixture files into the isolated module.
	copyFile(b, filepath.Join(repoRoot, "benchmarks", "neva.yml"), filepath.Join(moduleDir, "neva.yml"))
	copyFile(
		b,
		filepath.Join(repoRoot, "benchmarks", pkgName, "main.neva"),
		filepath.Join(progDir, "main.neva"),
	)

	// Compile the benchmark program once and return its output binary.
	buildProg := exec.Command(nevaBin, "build", pkgName)
	buildProg.Dir = moduleDir
	if output, err := buildProg.CombinedOutput(); err != nil {
		b.Fatalf("build benchmark program: %v\n%s", err, output)
	}

	return filepath.Join(moduleDir, "output")
}

// runProgramBinary executes one precompiled benchmark binary.
func runProgramBinary(b *testing.B, progPath string) {
	b.Helper()

	// #nosec G204 -- benchmark executes a fixed local binary built during setup.
	cmd := exec.Command(progPath)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		b.Fatalf("run benchmark program: %v", err)
	}
}

// copyFile copies one fixture file into the temp benchmark module.
func copyFile(b *testing.B, src, dst string) {
	b.Helper()

	data, err := os.ReadFile(src)
	if err != nil {
		b.Fatalf("read %s: %v", src, err)
	}

	// #nosec G306 -- benchmark fixture file should be readable for local runs/inspection.
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		b.Fatalf("write %s: %v", dst, err)
	}
}

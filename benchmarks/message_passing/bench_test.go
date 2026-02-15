package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// BenchmarkMessagePassingE2E benchmarks end-to-end runtime throughput.
// Setup (CLI build + program compilation) is excluded from timing.
func BenchmarkMessagePassingE2E(b *testing.B) {
	progPath := buildProgramOnce(b)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		// #nosec G204 -- benchmark executes a fixed local binary
		cmd := exec.Command(progPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			b.Fatalf("run benchmark program: %v\n%s", err, output)
		}
	}
}

func buildProgramOnce(b *testing.B) string {
	b.Helper()

	repoRoot := repoRoot(b)
	tmpDir := b.TempDir()

	moduleDir := filepath.Join(tmpDir, "bench-module")
	progDir := filepath.Join(moduleDir, "message_passing")
	if err := os.MkdirAll(progDir, 0o755); err != nil {
		b.Fatalf("create benchmark module dirs: %v", err)
	}

	copyFile(b, filepath.Join(repoRoot, "benchmarks", "neva.yml"), filepath.Join(moduleDir, "neva.yml"))
	copyFile(
		b,
		filepath.Join(repoRoot, "benchmarks", "message_passing", "main.neva"),
		filepath.Join(progDir, "main.neva"),
	)

	nevaBin := filepath.Join(tmpDir, "neva")
	buildCLI := exec.Command("go", "build", "-o", nevaBin, "./cmd/neva")
	buildCLI.Dir = repoRoot
	if output, err := buildCLI.CombinedOutput(); err != nil {
		b.Fatalf("build neva cli: %v\n%s", err, output)
	}

	buildProg := exec.Command(nevaBin, "build", "message_passing")
	buildProg.Dir = moduleDir
	if output, err := buildProg.CombinedOutput(); err != nil {
		b.Fatalf("build benchmark program: %v\n%s", err, output)
	}

	return filepath.Join(moduleDir, "output")
}

func copyFile(b *testing.B, src, dst string) {
	b.Helper()

	data, err := os.ReadFile(src)
	if err != nil {
		b.Fatalf("read %s: %v", src, err)
	}
	// #nosec G306 -- benchmark fixture file should be world-readable for local inspection.
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		b.Fatalf("write %s: %v", dst, err)
	}
}

func repoRoot(b *testing.B) string {
	b.Helper()

	cmd := exec.Command("go", "env", "GOMOD")
	output, err := cmd.Output()
	if err != nil {
		b.Fatalf("go env GOMOD: %v", err)
	}

	gomodPath := string(bytes.TrimSpace(output))
	if gomodPath == "" {
		b.Fatal("empty GOMOD path")
	}

	return filepath.Dir(gomodPath)
}

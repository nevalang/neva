package test

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
)

// BenchmarkMessagePassingE2E benchmarks precompiled program execution time.
// CLI build + Neva compilation are done once outside timed iterations.
func BenchmarkMessagePassingE2E(b *testing.B) {
	// Build artifacts once to keep the timed loop focused on runtime execution.
	progPath := buildProgramOnce(b)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		// Execute only the already-compiled benchmark binary.
		// #nosec G204 -- benchmark executes a fixed local binary built during setup.
		cmd := exec.Command(progPath)
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
		if err := cmd.Run(); err != nil {
			b.Fatalf("run benchmark program: %v", err)
		}
	}
}

// buildProgramOnce prepares an isolated benchmark module and compiles it once.
func buildProgramOnce(b *testing.B) string {
	b.Helper()

	// Discover repository root and create an isolated temp module workspace.
	repoRoot := e2e.FindRepoRoot(b)
	tmpDir := b.TempDir()

	moduleDir := filepath.Join(tmpDir, "bench-module")
	progDir := filepath.Join(moduleDir, "message_passing")
	if err := os.MkdirAll(progDir, 0o755); err != nil {
		b.Fatalf("create benchmark module dirs: %v", err)
	}

	// Copy benchmark fixture files into the isolated module.
	copyFile(b, filepath.Join(repoRoot, "benchmarks", "neva.yml"), filepath.Join(moduleDir, "neva.yml"))
	copyFile(
		b,
		filepath.Join(repoRoot, "benchmarks", "message_passing", "main.neva"),
		filepath.Join(progDir, "main.neva"),
	)

	// Build the CLI once, then compile the benchmark program once.
	nevaBin := e2e.BuildNevaBinary(b, repoRoot)

	buildProg := exec.Command(nevaBin, "build", "message_passing")
	buildProg.Dir = moduleDir
	if output, err := buildProg.CombinedOutput(); err != nil {
		b.Fatalf("build benchmark program: %v\n%s", err, output)
	}

	return filepath.Join(moduleDir, "output")
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

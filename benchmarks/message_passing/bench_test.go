package test

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/nevalang/neva/pkg/e2e"
)

// BenchmarkRuntimeE2E benchmarks precompiled runtime programs by data type path.
func BenchmarkRuntimeE2E(b *testing.B) {
	// Build the CLI once and reuse it for all benchmark programs.
	repoRoot := e2e.FindRepoRoot(b)
	nevaBin := e2e.BuildNevaBinary(b, repoRoot)

	benchPkgs, err := discoverRuntimeBenchPkgs(repoRoot)
	if err != nil {
		b.Fatalf("discover runtime benchmark packages: %v", err)
	}

	for _, benchPkg := range benchPkgs {
		benchName := strings.ReplaceAll(benchPkg, string(filepath.Separator), "_")
		b.Run(benchName, func(b *testing.B) {
			// Build the benchmark program once outside timed iterations.
			progPath := buildProgramOnce(b, repoRoot, nevaBin, benchPkg)

			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				runProgramBinary(b, progPath)
			}
		})
	}
}

// discoverRuntimeBenchPkgs finds all benchmark packages under benchmarks/runtime_bench.
func discoverRuntimeBenchPkgs(repoRoot string) ([]string, error) {
	benchmarksRoot := filepath.Join(repoRoot, "benchmarks")
	runtimeRoot := filepath.Join(benchmarksRoot, "runtime_bench")
	pkgs := make([]string, 0, 64)

	walkErr := filepath.WalkDir(runtimeRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "main.neva" {
			return nil
		}

		pkgDir := filepath.Dir(path)
		relDir, relErr := filepath.Rel(benchmarksRoot, pkgDir)
		if relErr != nil {
			return fmt.Errorf("resolve relative package dir for %q: %w", path, relErr)
		}
		pkgs = append(pkgs, relDir)
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	sort.Strings(pkgs)
	return pkgs, nil
}

// buildProgramOnce prepares an isolated module and compiles one benchmark package.
func buildProgramOnce(b *testing.B, repoRoot, nevaBin, pkgName string) string {
	b.Helper()

	// Create an isolated temp module workspace for one benchmark package.
	tmpDir := b.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	moduleDir := filepath.Join(tmpDir, "bench-module")
	progDir := filepath.Join(moduleDir, pkgName)
	if err := os.MkdirAll(progDir, 0o755); err != nil {
		b.Fatalf("create benchmark module dirs: %v", err)
	}
	if err := prepareBenchmarkHome(repoRoot, homeDir); err != nil {
		b.Fatalf("prepare benchmark home: %v", err)
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
	buildProg.Env = append(os.Environ(), "HOME="+homeDir)
	if output, err := buildProg.CombinedOutput(); err != nil {
		b.Fatalf("build benchmark program: %v\n%s", err, output)
	}

	return filepath.Join(moduleDir, "output")
}

// prepareBenchmarkHome creates an isolated Neva home with local stdlib available.
func prepareBenchmarkHome(repoRoot, homeDir string) error {
	nevaHome := filepath.Join(homeDir, "neva")
	if err := os.MkdirAll(nevaHome, 0o755); err != nil {
		return err
	}

	stdSrc := filepath.Join(repoRoot, "std")
	stdDst := filepath.Join(nevaHome, "std")
	if err := os.Symlink(stdSrc, stdDst); err == nil {
		return nil
	}

	return copyDir(stdSrc, stdDst)
}

// copyDir copies a directory recursively, preserving file modes.
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		// #nosec G306 -- benchmark fixture files are read-only test assets.
		return os.WriteFile(target, data, 0o644)
	})
}

// runProgramBinary executes one precompiled benchmark binary.
func runProgramBinary(b *testing.B, progPath string) {
	b.Helper()

	// #nosec G204 -- benchmark executes a fixed local binary built during setup.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, progPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			b.Fatalf("run benchmark program timed out: %s", progPath)
		}
		b.Fatalf("run benchmark program: %v\n%s", err, output)
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

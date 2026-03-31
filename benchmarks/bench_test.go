package benchmarks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/nevalang/neva/pkg/e2e"
	nevaos "github.com/nevalang/neva/pkg/os"
)

// BenchmarkRuntimeE2E benchmarks precompiled runtime programs by data type path.
func BenchmarkRuntimeE2E(b *testing.B) {
	// Build the CLI once and reuse it for all benchmark programs.
	repoRoot := e2e.FindRepoRoot(b)
	nevaBin := e2e.BuildNevaBinary(b, repoRoot)

	benchPkgs, err := discoverBenchmarkPkgs(repoRoot)
	if err != nil {
		b.Fatalf("discover runtime benchmark packages: %v", err)
	}

	for _, benchPkg := range benchPkgs {
		benchName := strings.ReplaceAll(benchPkg, string(filepath.Separator), "_")
		b.Run(benchName, func(b *testing.B) {
			// Build the benchmark program once outside timed iterations.
			progPath := buildProgramOnce(b, repoRoot, nevaBin, benchPkg)

			for b.Loop() {
				runProgramBinary(b, progPath)
			}
		})
	}
}

// discoverBenchmarkPkgs finds all benchmark packages under benchmark tiers.
func discoverBenchmarkPkgs(repoRoot string) ([]string, error) {
	benchmarksRoot := filepath.Join(repoRoot, "benchmarks")
	pkgs := make([]string, 0, 64)

	for _, group := range []string{"atomic", "simple", "complex"} {
		groupRoot := filepath.Join(benchmarksRoot, group)
		if _, statErr := os.Stat(groupRoot); errors.Is(statErr, os.ErrNotExist) {
			// Allow partial rollout where only one group is landed.
			continue
		}

		walkErr := filepath.WalkDir(groupRoot, func(path string, d fs.DirEntry, err error) error {
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
	}

	sort.Strings(pkgs)
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no benchmark packages found under %s", benchmarksRoot)
	}
	return pkgs, nil
}

// buildProgramOnce prepares an isolated module and compiles one benchmark package.
func buildProgramOnce(b *testing.B, repoRoot, nevaBin, pkgName string) string {
	b.Helper()

	// e2e.BuildNevaBinary caches the shared CLI, but each benchmark fixture still
	// needs its own compiled output binary because the input package differs.
	// Create an isolated temp module workspace for one benchmark package.
	tmpDir := b.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	moduleDir := filepath.Join(tmpDir, "bench-module")
	progDir := filepath.Join(moduleDir, pkgName)
	if err := e2e.PrepareIsolatedNevaHome(repoRoot, homeDir); err != nil {
		b.Fatalf("prepare benchmark home: %v", err)
	}

	// Copy benchmark module config plus the whole benchmark package fixture tree.
	configPath := filepath.Join(repoRoot, "benchmarks", "neva.yml")
	configInfo, err := os.Stat(configPath)
	if err != nil {
		b.Fatalf("stat benchmark module config: %v", err)
	}
	if err := nevaos.CopyFile(configPath, filepath.Join(moduleDir, "neva.yml"), configInfo.Mode()); err != nil {
		b.Fatalf("copy benchmark module config: %v", err)
	}
	if err := nevaos.CopyDir(filepath.Join(repoRoot, "benchmarks", pkgName), progDir); err != nil {
		b.Fatalf("copy benchmark package fixture: %v", err)
	}

	// Compile the benchmark program once and return its output binary.
	buildProg := exec.Command(nevaBin, "build", pkgName)
	buildProg.Dir = moduleDir
	buildProg.Env = append(os.Environ(), "HOME="+homeDir)
	if output, err := buildProg.CombinedOutput(); err != nil {
		b.Fatalf("build benchmark program: %v\n%s", err, output)
	}

	return filepath.Join(moduleDir, "output")
}

// runProgramBinary executes one precompiled benchmark binary.
func runProgramBinary(b *testing.B, progPath string) {
	b.Helper()

	// #nosec G204 -- benchmark executes a fixed local binary built during setup.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, progPath)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			b.Fatalf("run benchmark program timed out: %s", progPath)
		}
		b.Fatalf("run benchmark program %s: %v", progPath, err)
	}
}

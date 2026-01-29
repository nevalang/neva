package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	manifestSearchLimit = 10
	watchDebounce       = 200 * time.Millisecond
)

//nolint:gocyclo // Control flow handles multiple run and watch edge cases.
func watchAndRun(ctx context.Context, moduleRoot string, run func(context.Context) error) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}
	defer watcher.Close()

	if err := addWatchersRecursively(watcher, moduleRoot); err != nil {
		return fmt.Errorf("watch module: %w", err)
	}

	fmt.Printf("Watching %s for changes. Press Ctrl+C to stop.\n", moduleRoot)

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	var (
		debounce   *time.Timer
		debounceC  <-chan time.Time
		runPending bool
	)

	scheduleRun := func() {
		runPending = true
		if debounce == nil {
			debounce = time.NewTimer(watchDebounce)
			debounceC = debounce.C
			return
		}

		if !debounce.Stop() {
			select {
			case <-debounce.C:
			default:
			}
		}

		debounce.Reset(watchDebounce)
		debounceC = debounce.C
	}

	for {
		select {
		case <-ctx.Done():
			if debounce != nil {
				debounce.Stop()
			}
			return nil
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintln(os.Stderr, "watch error:", err)
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}

			if event.Op&fsnotify.Create != 0 {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					if err := addWatchersRecursively(watcher, event.Name); err != nil {
						fmt.Fprintf(os.Stderr, "failed to watch new directory %s: %v\n", event.Name, err)
					}
				}
			}

			if !shouldTrigger(event.Name) {
				continue
			}

			scheduleRun()
		case <-debounceC:
			debounceC = nil
			if debounce != nil {
				debounce.Stop()
				debounce = nil
			}

			if !runPending {
				continue
			}

			runPending = false

			if err := run(ctx); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func addWatchersRecursively(watcher *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		return watcher.Add(path)
	})
}

func shouldTrigger(path string) bool {
	name := filepath.Base(path)
	lower := strings.ToLower(name)
	if lower == "neva.yaml" || lower == "neva.yml" {
		return true
	}

	return strings.EqualFold(filepath.Ext(name), ".neva")
}

func findModuleRoot(workdir, pkgPath string) (string, error) {
	start := pkgPath
	if start == "" {
		start = "."
	}

	if !filepath.IsAbs(start) {
		start = filepath.Join(workdir, start)
	}

	start = filepath.Clean(start)

	info, err := os.Stat(start)
	if err != nil {
		return "", fmt.Errorf("resolve module root: %w", err)
	}

	if !info.IsDir() {
		start = filepath.Dir(start)
	}

	current := start
	for depth := 0; depth <= manifestSearchLimit; depth++ {
		if hasManifest(current) {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}

		current = parent
	}

	return "", fmt.Errorf("manifest file not found near %s", start)
}

func hasManifest(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, "neva.yaml")); err == nil {
		return true
	}

	if _, err := os.Stat(filepath.Join(dir, "neva.yml")); err == nil {
		return true
	}

	return false
}

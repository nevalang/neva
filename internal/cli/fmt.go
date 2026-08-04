package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/pmezard/go-difflib/difflib"
	cli "github.com/urfave/cli/v2"

	"github.com/nevalang/neva/pkg/formatter"
)

const nevaSourceExtension = ".neva"

// fmtMode selects the formatter's output behavior.
type fmtMode string

const (
	fmtPrint fmtMode = ""
	fmtWrite fmtMode = "write"
	fmtDiff  fmtMode = "diff"
	fmtList  fmtMode = "list"
	fmtCheck fmtMode = "check"
)

// fmtOptions collects CLI flags that affect formatting output.
type fmtOptions struct {
	mode         fmtMode
	reportErrors bool
}

// fmtResult contains the result of formatting one source file.
type fmtResult struct {
	err       error
	original  []byte
	formatted []byte
}

// newFmtCmd creates the first-party source formatting command.
func newFmtCmd() *cli.Command {
	return &cli.Command{
		Name:      "fmt",
		Usage:     "Format Neva source files",
		ArgsUsage: "[file or directory ...]",
		Description: "With no paths, formats one source file from standard input. " +
			"Directories are walked recursively; only .git is skipped.",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "w", Usage: "write formatted source back to files"},
			&cli.BoolFlag{Name: "d", Usage: "display a unified diff instead of formatted source"},
			&cli.BoolFlag{Name: "l", Usage: "list files whose formatting differs"},
			&cli.BoolFlag{Name: "check", Usage: "list unformatted files and exit non-zero"},
			&cli.BoolFlag{Name: "e", Usage: "report all independent input errors"},
		},
		Action: func(cliCtx *cli.Context) error {
			options, err := parseFmtOptions(cliCtx)
			if err != nil {
				return exitWithCode(err, 2)
			}

			return runFmt(cliCtx.Args().Slice(), options, os.Stdin, os.Stdout)
		},
	}
}

// parseFmtOptions validates formatter output flags.
func parseFmtOptions(cliCtx *cli.Context) (fmtOptions, error) {
	options := fmtOptions{reportErrors: cliCtx.Bool("e")}
	modes := []struct {
		mode fmtMode
		set  bool
	}{
		{fmtWrite, cliCtx.Bool("w")},
		{fmtDiff, cliCtx.Bool("d")},
		{fmtList, cliCtx.Bool("l")},
		{fmtCheck, cliCtx.Bool("check")},
	}

	for _, candidate := range modes {
		if !candidate.set {
			continue
		}
		if options.mode != fmtPrint {
			return fmtOptions{}, errors.New("-w, -d, -l, and -check are mutually exclusive")
		}
		options.mode = candidate.mode
	}

	return options, nil
}

// runFmt formats standard input or the requested source paths.
func runFmt(paths []string, options fmtOptions, input io.Reader, output io.Writer) error {
	if len(paths) == 0 {
		return formatStandardInput(options, input, output)
	}

	files, err := collectFmtFiles(paths)
	if err != nil {
		return err
	}

	changed, err := applyFmtResults(files, formatFiles(files), options, output)
	if err != nil {
		return err
	}
	if options.mode == fmtCheck && changed {
		return exitWithCode(nil, 1)
	}
	return nil
}

// formatStandardInput renders one complete source file from standard input.
func formatStandardInput(options fmtOptions, input io.Reader, output io.Writer) error {
	if options.mode != fmtPrint {
		return errors.New("-w, -d, -l, and -check require at least one file or directory")
	}
	source, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("read standard input: %w", err)
	}
	formatted, parseErr := formatter.Format(source)
	if parseErr != nil {
		return fmt.Errorf("format standard input: %w", parseErr)
	}
	if _, err := output.Write(formatted); err != nil {
		return fmt.Errorf("write standard output: %w", err)
	}
	return nil
}

// applyFmtResults emits or writes completed formatting results in input order.
func applyFmtResults(paths []string, results []fmtResult, options fmtOptions, output io.Writer) (bool, error) {
	changed := false
	var reported []error
	for index := range results {
		result := &results[index]
		resultChanged, resultErr := applyFmtResult(paths[index], result, options.mode, output)
		changed = changed || resultChanged
		if resultErr == nil {
			continue
		}
		if !options.reportErrors {
			return changed, resultErr
		}
		reported = append(reported, resultErr)
	}
	return changed, errors.Join(reported...)
}

// applyFmtResult handles one completed result for the selected output mode.
func applyFmtResult(path string, result *fmtResult, mode fmtMode, output io.Writer) (bool, error) {
	if result.err != nil {
		return false, result.err
	}
	changed := !bytes.Equal(result.original, result.formatted)
	if !changed && mode != fmtPrint {
		return false, nil
	}
	return changed, writeFmtResult(path, result, mode, output)
}

// collectFmtFiles expands file and directory arguments in deterministic order.
func collectFmtFiles(paths []string) ([]string, error) {
	files := make([]string, 0, len(paths))
	for _, sourcePath := range paths {
		info, err := os.Stat(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("inspect %s: %w", sourcePath, err)
		}
		if !info.IsDir() {
			if err := appendFmtFile(&files, sourcePath); err != nil {
				return nil, err
			}
			continue
		}

		if err := collectFmtDirectory(sourcePath, &files); err != nil {
			return nil, fmt.Errorf("walk %s: %w", sourcePath, err)
		}
	}

	return files, nil
}

// appendFmtFile validates and appends one explicit source file.
func appendFmtFile(files *[]string, path string) error {
	if filepath.Ext(path) != nevaSourceExtension {
		return fmt.Errorf("%s is not a .neva file", path)
	}
	*files = append(*files, path)
	return nil
}

// collectFmtDirectory appends source files below root while skipping Git metadata.
func collectFmtDirectory(root string, files *[]string) error {
	err := filepath.WalkDir(root, func(entryPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entryPath != root && entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(entry.Name()) == nevaSourceExtension {
			*files = append(*files, entryPath)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory: %w", err)
	}
	return nil
}

// formatFiles formats paths concurrently and retains their input order.
func formatFiles(paths []string) []fmtResult {
	results := make([]fmtResult, len(paths))
	jobs := make(chan int)
	workers := min(runtime.GOMAXPROCS(0), len(paths))
	var group sync.WaitGroup
	for range workers {
		group.Go(func() {
			for index := range jobs {
				results[index] = formatFile(paths[index])
			}
		})
	}
	for index := range paths {
		jobs <- index
	}
	close(jobs)
	group.Wait()
	return results
}

// formatFile reads and formats one source file without writing it.
func formatFile(path string) fmtResult {
	result := fmtResult{}
	source, err := os.ReadFile(path)
	if err != nil {
		result.err = fmt.Errorf("read %s: %w", path, err)
		return result
	}
	formatted, parseErr := formatter.Format(source)
	if parseErr != nil {
		result.err = fmt.Errorf("format %s: %w", path, parseErr)
		return result
	}
	result.original = source
	result.formatted = formatted
	return result
}

// writeFmtResult emits one changed result according to mode.
func writeFmtResult(path string, result *fmtResult, mode fmtMode, output io.Writer) error {
	switch mode {
	case fmtPrint:
		_, err := output.Write(result.formatted)
		if err != nil {
			return fmt.Errorf("write standard output: %w", err)
		}
	case fmtWrite:
		if err := writeFormattedFile(path, result.formatted); err != nil {
			return err
		}
	case fmtDiff:
		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(result.original)),
			B:        difflib.SplitLines(string(result.formatted)),
			FromFile: path,
			ToFile:   path,
			Context:  3,
		})
		if err != nil {
			return fmt.Errorf("render diff for %s: %w", path, err)
		}
		if _, err := io.WriteString(output, diff); err != nil {
			return fmt.Errorf("write standard output: %w", err)
		}
	case fmtList, fmtCheck:
		if _, err := fmt.Fprintln(output, path); err != nil {
			return fmt.Errorf("write standard output: %w", err)
		}
	default:
		return fmt.Errorf("unknown formatter mode %q", mode)
	}
	return nil
}

// writeFormattedFile atomically replaces path while preserving permissions.
func writeFormattedFile(path string, formatted []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("inspect %s: %w", path, err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".fmt-*")
	if err != nil {
		return fmt.Errorf("create temporary file for %s: %w", path, err)
	}
	temporaryPath := temporary.Name()
	if err := temporary.Chmod(info.Mode().Perm()); err != nil {
		return closeAndRemoveTemporaryFile(temporary, temporaryPath, fmt.Errorf("set permissions for %s: %w", temporaryPath, err))
	}
	if _, err := temporary.Write(formatted); err != nil {
		return closeAndRemoveTemporaryFile(temporary, temporaryPath, fmt.Errorf("write temporary file for %s: %w", path, err))
	}
	if err := temporary.Close(); err != nil {
		return removeTemporaryFile(temporaryPath, fmt.Errorf("close temporary file for %s: %w", path, err))
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return removeTemporaryFile(temporaryPath, fmt.Errorf("replace %s: %w", path, err))
	}
	return nil
}

// closeAndRemoveTemporaryFile closes and removes a failed replacement file.
func closeAndRemoveTemporaryFile(file *os.File, path string, cause error) error {
	if err := file.Close(); err != nil {
		cause = fmt.Errorf("%w; close temporary file %s: %w", cause, path, err)
	}
	return removeTemporaryFile(path, cause)
}

// removeTemporaryFile removes a failed replacement file while preserving cause.
func removeTemporaryFile(path string, cause error) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("%w; remove temporary file %s: %w", cause, path, err)
	}
	return cause
}

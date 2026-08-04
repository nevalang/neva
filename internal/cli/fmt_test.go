package cli

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	urfavecli "github.com/urfave/cli/v2"
)

const unformattedFmtSource = "def Main(start any) (stop any) {\n:start->:stop\n}\n"

const formattedFmtSource = "def Main(start any) (stop any) {\n\t:start -> :stop\n}\n"

func TestRunFmtStandardInput(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	err := runFmt(nil, fmtOptions{}, bytes.NewBufferString(unformattedFmtSource), &output)
	require.NoError(t, err)
	require.Equal(t, formattedFmtSource, output.String())
}

func TestRunFmtPrintsAlreadyFormattedFile(t *testing.T) {
	t.Parallel()

	path := writeFmtTestFile(t, t.TempDir(), "main.neva", formattedFmtSource)
	var output bytes.Buffer
	err := runFmt([]string{path}, fmtOptions{}, nil, &output)
	require.NoError(t, err)
	require.Equal(t, formattedFmtSource, output.String())
}

func TestRunFmtModes(t *testing.T) {
	dir := t.TempDir()
	path := writeFmtTestFile(t, dir, "main.neva", unformattedFmtSource)

	t.Run("print", func(t *testing.T) {
		var output bytes.Buffer
		err := runFmt([]string{path}, fmtOptions{}, nil, &output)
		require.NoError(t, err)
		require.Equal(t, formattedFmtSource, output.String())
	})

	t.Run("diff", func(t *testing.T) {
		var output bytes.Buffer
		err := runFmt([]string{path}, fmtOptions{mode: fmtDiff}, nil, &output)
		require.NoError(t, err)
		require.Contains(t, output.String(), "--- "+path)
		require.Contains(t, output.String(), "+++ "+path)
		require.Contains(t, output.String(), "-:start->:stop")
		require.Contains(t, output.String(), "+\t:start -> :stop")
	})

	t.Run("list", func(t *testing.T) {
		var output bytes.Buffer
		err := runFmt([]string{path}, fmtOptions{mode: fmtList}, nil, &output)
		require.NoError(t, err)
		require.Equal(t, path+"\n", output.String())
	})

	t.Run("check", func(t *testing.T) {
		var output bytes.Buffer
		err := runFmt([]string{path}, fmtOptions{mode: fmtCheck}, nil, &output)
		require.Error(t, err)
		require.Equal(t, 1, ExitCode(err))
		require.Equal(t, path+"\n", output.String())
	})

	t.Run("write", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			require.NoError(t, os.Chmod(path, 0o751))
		}
		var output bytes.Buffer
		err := runFmt([]string{path}, fmtOptions{mode: fmtWrite}, nil, &output)
		require.NoError(t, err)
		contents, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, formattedFmtSource, string(contents))
		if runtime.GOOS != "windows" {
			info, err := os.Stat(path)
			require.NoError(t, err)
			require.Equal(t, os.FileMode(0o751), info.Mode().Perm())
		}
	})
}

func TestRunFmtWalksDirectoriesDeterministically(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	first := writeFmtTestFile(t, dir, "a.neva", unformattedFmtSource)
	second := writeFmtTestFile(t, dir, "nested/b.neva", unformattedFmtSource)
	writeFmtTestFile(t, dir, ".git/ignored.neva", unformattedFmtSource)
	third := writeFmtTestFile(t, dir, ".neva/deps/c.neva", unformattedFmtSource)
	fourth := writeFmtTestFile(t, dir, "node_modules/d.neva", unformattedFmtSource)
	fifth := writeFmtTestFile(t, dir, "vendor/e.neva", unformattedFmtSource)

	var output bytes.Buffer
	err := runFmt([]string{dir}, fmtOptions{mode: fmtList}, nil, &output)
	require.NoError(t, err)
	require.Equal(t, third+"\n"+first+"\n"+second+"\n"+fourth+"\n"+fifth+"\n", output.String())
}

func TestRunFmtReportsAllRequestedErrors(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	first := writeFmtTestFile(t, dir, "first.neva", "def First(\n")
	second := writeFmtTestFile(t, dir, "second.neva", "def Second(\n")
	valid := writeFmtTestFile(t, dir, "valid.neva", unformattedFmtSource)

	var output bytes.Buffer
	err := runFmt([]string{first, second, valid}, fmtOptions{mode: fmtList, reportErrors: true}, nil, &output)
	require.Error(t, err)
	require.ErrorContains(t, err, first)
	require.ErrorContains(t, err, second)
	require.Equal(t, valid+"\n", output.String())
}

func TestRunFmtRejectsInvalidInputsWithoutWriting(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	invalid := writeFmtTestFile(t, dir, "broken.neva", "def Broken(\n")

	var output bytes.Buffer
	err := runFmt([]string{invalid}, fmtOptions{mode: fmtWrite}, nil, &output)
	require.Error(t, err)
	contents, readErr := os.ReadFile(invalid)
	require.NoError(t, readErr)
	require.Equal(t, "def Broken(\n", string(contents))

	nonNeva := writeFmtTestFile(t, dir, "notes.txt", "not Neva")
	err = runFmt([]string{nonNeva}, fmtOptions{}, nil, io.Discard)
	require.ErrorContains(t, err, "is not a .neva file")
}

func TestRunFmtWritePreflightsAllFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	valid := writeFmtTestFile(t, dir, "valid.neva", unformattedFmtSource)
	invalid := writeFmtTestFile(t, dir, "invalid.neva", "def Invalid(\n")

	err := runFmt([]string{valid, invalid}, fmtOptions{mode: fmtWrite}, nil, io.Discard)
	require.Error(t, err)

	contents, readErr := os.ReadFile(valid)
	require.NoError(t, readErr)
	require.Equal(t, unformattedFmtSource, string(contents))
}

func TestRunFmtRejectsOutputModesForStandardInput(t *testing.T) {
	t.Parallel()

	err := runFmt(nil, fmtOptions{mode: fmtWrite}, bytes.NewBufferString(unformattedFmtSource), io.Discard)
	require.ErrorContains(t, err, "require at least one file or directory")
}

func TestParseFmtOptionsRejectsConflictingOutputModes(t *testing.T) {
	t.Parallel()

	flags := flag.NewFlagSet("fmt", flag.ContinueOnError)
	flags.Bool("w", true, "")
	flags.Bool("d", true, "")
	flags.Bool("l", false, "")
	flags.Bool("check", false, "")
	flags.Bool("e", false, "")
	_, err := parseFmtOptions(urfavecli.NewContext(nil, flags, nil))
	require.ErrorContains(t, err, "mutually exclusive")
}

func TestExitCode(t *testing.T) {
	t.Parallel()

	require.Equal(t, 1, ExitCode(errors.New("ordinary error")))
	require.Equal(t, 2, ExitCode(exitWithCode(errors.New("usage error"), 2)))
}

func writeFmtTestFile(t *testing.T, dir, name, contents string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))
	return path
}

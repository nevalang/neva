package funcs

import (
	"os"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

func benchmarkFilePath(b *testing.B, data string) string {
	b.Helper()

	path := b.TempDir() + "/file.txt"
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		b.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func openBenchmarkFileWithData(b *testing.B, data string) *os.File {
	b.Helper()

	file, err := os.Open(benchmarkFilePath(b, data))
	if err != nil {
		b.Fatalf("Open() error = %v", err)
	}
	return file
}

func assertBenchmarkOutput(b *testing.B, outChans map[string]chan runtime.OrderedMsg, outName string) {
	b.Helper()

	select {
	case <-outChans[outName]:
	case <-time.After(time.Second):
		b.Fatalf("no output on %q", outName)
	}
}

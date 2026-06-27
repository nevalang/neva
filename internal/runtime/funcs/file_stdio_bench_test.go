package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func BenchmarkFileStdioHandles(b *testing.B) {
	benchmarks := []struct {
		creator runtime.FuncCreator
		name    string
	}{
		{name: "stdin", creator: fileStdin{}},
		{name: "stdout", creator: fileStdout{}},
		{name: "stderr", creator: fileStderr{}},
	}

	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			benchmarkStdioHandle(b, benchmark.creator)
		})
	}
}

func benchmarkStdioHandle(b *testing.B, creator runtime.FuncCreator) {
	b.Helper()

	io, inChans, outChans := newIO([]string{"sig"}, []string{"res"})
	handler, err := creator.Create(io, nil)
	if err != nil {
		b.Fatalf("Create() error = %v", err)
	}
	cancel, done := runHandler(handler)
	defer waitForHandler(b, cancel, done)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		inChans["sig"] <- runtime.OrderedMsg{Msg: emptyStruct()}
		assertBenchmarkOutput(b, outChans, "res")
	}
}

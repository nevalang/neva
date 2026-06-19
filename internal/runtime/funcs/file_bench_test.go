package funcs

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

func BenchmarkFileOpenHandle(b *testing.B) {
	store := runtime.NewFileHandles()
	path := benchmarkFilePath(b, "payload")
	_, _, outChans := newIO(nil, []string{"res", "err"})
	resOut := mustSingleOutport(b, outChans, "res")
	errOut := mustSingleOutport(b, outChans, "err")
	msg := runtime.OrderedMsg{Msg: runtime.NewStringMsg(path)}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if !(fileOpen{handles: store}).handleFileMessage(context.Background(), msg, resOut, errOut) {
			b.Fatal("handleFileMessage() returned false")
		}
		handleID := receiveIntResOutput(b, outChans)

		b.StopTimer()
		mustCloseHandle(b, store, handleID)
		b.StartTimer()
	}
}

func BenchmarkFileCreateHandle(b *testing.B) {
	store := runtime.NewFileHandles()
	path := b.TempDir() + "/created.txt"
	_, _, outChans := newIO(nil, []string{"res", "err"})
	resOut := mustSingleOutport(b, outChans, "res")
	errOut := mustSingleOutport(b, outChans, "err")
	msg := runtime.OrderedMsg{Msg: runtime.NewStringMsg(path)}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if !(fileCreate{handles: store}).handleFileMessage(context.Background(), msg, resOut, errOut) {
			b.Fatal("handleFileMessage() returned false")
		}
		handleID := receiveIntResOutput(b, outChans)

		b.StopTimer()
		mustCloseHandle(b, store, handleID)
		b.StartTimer()
	}
}

func BenchmarkFileCloseHandle(b *testing.B) {
	store := runtime.NewFileHandles()
	_, _, outChans := newIO(nil, []string{"res", "err"})
	resOut := mustSingleOutport(b, outChans, "res")
	errOut := mustSingleOutport(b, outChans, "err")

	b.ReportAllocs()
	for range b.N {
		b.StopTimer()
		file, err := os.CreateTemp(b.TempDir(), "close-*.txt")
		if err != nil {
			b.Fatalf("CreateTemp() error = %v", err)
		}
		handleID := store.Add(file)
		msg := runtime.OrderedMsg{Msg: runtime.NewIntMsg(handleID)}
		b.StartTimer()

		if !(fileClose{handles: store}).handleFileMessage(context.Background(), msg, resOut, errOut) {
			b.Fatal("handleFileMessage() returned false")
		}
		assertBenchmarkOutput(b, outChans, "res")
	}
}

func BenchmarkFileReadAllHandle(b *testing.B) {
	store := runtime.NewFileHandles()
	file := openBenchmarkFileWithData(b, "payload")
	handleID := store.Add(file)
	defer mustCloseHandle(b, store, handleID)

	_, _, outChans := newIO(nil, []string{"res", "handle", "err"})
	resOut := mustSingleOutport(b, outChans, "res")
	handleOut := mustSingleOutport(b, outChans, "handle")
	errOut := mustSingleOutport(b, outChans, "err")
	msg := runtime.OrderedMsg{Msg: runtime.NewIntMsg(handleID)}

	b.ReportAllocs()
	for range b.N {
		b.StopTimer()
		if _, err := file.Seek(0, 0); err != nil {
			b.Fatalf("Seek() error = %v", err)
		}
		b.StartTimer()

		if !(fileReadAllHandle{handles: store}).handleFileMessage(context.Background(), msg, resOut, handleOut, errOut) {
			b.Fatal("handleFileMessage() returned false")
		}
		assertBenchmarkOutput(b, outChans, "res")
		assertBenchmarkOutput(b, outChans, "handle")
	}
}

func BenchmarkFileWriteAllHandle(b *testing.B) {
	store := runtime.NewFileHandles()
	file, err := os.CreateTemp(b.TempDir(), "write-*.txt")
	if err != nil {
		b.Fatalf("CreateTemp() error = %v", err)
	}
	handleID := store.Add(file)
	defer mustCloseHandle(b, store, handleID)

	_, _, outChans := newIO(nil, []string{"res", "err"})
	resOut := mustSingleOutport(b, outChans, "res")
	errOut := mustSingleOutport(b, outChans, "err")
	fileMsg := runtime.OrderedMsg{Msg: runtime.NewIntMsg(handleID)}
	dataMsg := runtime.OrderedMsg{Msg: runtime.NewBytesMsg([]byte("payload"))}

	b.ReportAllocs()
	for range b.N {
		b.StopTimer()
		if err := file.Truncate(0); err != nil {
			b.Fatalf("Truncate() error = %v", err)
		}
		if _, err := file.Seek(0, 0); err != nil {
			b.Fatalf("Seek() error = %v", err)
		}
		b.StartTimer()

		if !(fileWriteAllHandle{handles: store}).handleFileMessage(context.Background(), fileMsg, dataMsg, resOut, errOut) {
			b.Fatal("handleFileMessage() returned false")
		}
		assertBenchmarkOutput(b, outChans, "res")
	}
}

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

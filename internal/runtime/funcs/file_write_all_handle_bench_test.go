package funcs

import (
	"context"
	"os"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

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

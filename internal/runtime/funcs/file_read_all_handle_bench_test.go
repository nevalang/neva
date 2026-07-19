package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

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

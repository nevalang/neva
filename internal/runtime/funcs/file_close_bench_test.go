package funcs

import (
	"context"
	"os"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

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

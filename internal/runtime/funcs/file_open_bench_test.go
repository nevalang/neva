package funcs

import (
	"context"
	"testing"

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

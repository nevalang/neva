package funcs

import (
	"context"
	"os"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestFileHandleID(t *testing.T) {
	t.Parallel()

	id, err := fileHandleID(runtime.NewIntMsg(42))
	if err != nil {
		t.Fatalf("fileHandleID() unexpected error = %v", err)
	}
	if id != 42 {
		t.Fatalf("fileHandleID() id = %d, want 42", id)
	}

	if _, err := fileHandleID(runtime.NewStringMsg("bad")); err == nil {
		t.Fatal("fileHandleID() expected type error")
	}
}

func TestRegistryContainsFileHandleCreators(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	keys := []string{
		"file_open",
		"file_create",
		"file_close",
		"file_read_all",
		"file_write_all",
		"file_stdin",
		"file_stdout",
		"file_stderr",
	}

	for _, key := range keys {
		if _, ok := registry[key]; !ok {
			t.Fatalf("NewRegistry() missing key %q", key)
		}
	}
}

func TestFileReadAllHandlePassesHandleOnReadError(t *testing.T) {
	t.Parallel()

	store := runtime.NewFileHandles()
	handleID := addClosedTempFile(t, store)
	_, _, outChans := newIO(nil, []string{"res", "handle", "err"})

	// A read error must still expose the handle so user code can close it.
	ok := fileReadAllHandle{handles: store}.handleFileMessage(
		context.Background(),
		runtime.OrderedMsg{Msg: runtime.NewIntMsg(handleID)},
		mustSingleOutport(t, outChans, "res"),
		mustSingleOutport(t, outChans, "handle"),
		mustSingleOutport(t, outChans, "err"),
	)
	if !ok {
		t.Fatal("handleFileMessage() returned false")
	}

	assertOutputEquals(t, outChans, "handle", runtime.NewIntMsg(handleID), nil)
	assertRuntimeErrorOutput(t, outChans, "err")
}

func TestFileWriteAllHandlePassesHandleOnWriteError(t *testing.T) {
	t.Parallel()

	store := runtime.NewFileHandles()
	handleID := addClosedTempFile(t, store)
	_, _, outChans := newIO(nil, []string{"res", "err"})

	// A write error must still expose the handle so user code can close it.
	ok := fileWriteAllHandle{handles: store}.handleFileMessage(
		context.Background(),
		runtime.OrderedMsg{Msg: runtime.NewIntMsg(handleID)},
		runtime.OrderedMsg{Msg: runtime.NewBytesMsg([]byte("data"))},
		mustSingleOutport(t, outChans, "res"),
		mustSingleOutport(t, outChans, "err"),
	)
	if !ok {
		t.Fatal("handleFileMessage() returned false")
	}

	assertOutputEquals(t, outChans, "res", runtime.NewIntMsg(handleID), nil)
	assertRuntimeErrorOutput(t, outChans, "err")
}

func addClosedTempFile(t *testing.T, store *runtime.FileHandles) int64 {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "file-handle-error-*.txt")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	return store.Add(tmpFile)
}

func mustSingleOutport(
	t *testing.T,
	outChans map[string]chan runtime.OrderedMsg,
	name string,
) runtime.SingleOutport {
	t.Helper()

	tracer := runtime.NewTracer()
	port := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "test/out", Port: name},
		runtime.NoEffectInterceptor{},
		outChans[name],
	)

	return *port
}

func assertRuntimeErrorOutput(t *testing.T, outChans map[string]chan runtime.OrderedMsg, outName string) {
	t.Helper()

	select {
	case got := <-outChans[outName]:
		if _, ok := got.Msg.(runtime.StructMsg); !ok {
			t.Fatalf("error output = %T, want struct error", got.Msg)
		}
	default:
		t.Fatalf("no runtime error on %q", outName)
	}
}

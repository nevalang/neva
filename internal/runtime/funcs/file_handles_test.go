package funcs

import (
	"context"
	"os"
	"testing"
	"time"

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

func TestFileOpenHandleOpensFile(t *testing.T) {
	t.Parallel()

	store := runtime.NewFileHandles()
	path := writeTempFile(t, "hello")
	_, _, outChans := newIO(nil, []string{"res", "err"})

	ok := fileOpen{handles: store}.handleFileMessage(
		context.Background(),
		runtime.OrderedMsg{Msg: runtime.NewStringMsg(path)},
		mustSingleOutport(t, outChans, "res"),
		mustSingleOutport(t, outChans, "err"),
	)
	if !ok {
		t.Fatal("handleFileMessage() returned false")
	}

	handleID := receiveIntResOutput(t, outChans)
	file := mustGetFile(t, store, handleID)
	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("opened file data = %q, want hello", data)
	}
	mustCloseHandle(t, store, handleID)
}

func TestFileCreateHandleCreatesFile(t *testing.T) {
	t.Parallel()

	store := runtime.NewFileHandles()
	path := t.TempDir() + "/created.txt"
	_, _, outChans := newIO(nil, []string{"res", "err"})

	ok := fileCreate{handles: store}.handleFileMessage(
		context.Background(),
		runtime.OrderedMsg{Msg: runtime.NewStringMsg(path)},
		mustSingleOutport(t, outChans, "res"),
		mustSingleOutport(t, outChans, "err"),
	)
	if !ok {
		t.Fatal("handleFileMessage() returned false")
	}

	handleID := receiveIntResOutput(t, outChans)
	file := mustGetFile(t, store, handleID)
	if _, err := file.WriteString("created"); err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}
	mustCloseHandle(t, store, handleID)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "created" {
		t.Fatalf("created file data = %q, want created", data)
	}
}

func TestFileCloseHandleClosesFile(t *testing.T) {
	t.Parallel()

	store := runtime.NewFileHandles()
	file := createTempFile(t)
	handleID := store.Add(file)
	_, _, outChans := newIO(nil, []string{"res", "err"})

	ok := fileClose{handles: store}.handleFileMessage(
		context.Background(),
		runtime.OrderedMsg{Msg: runtime.NewIntMsg(handleID)},
		mustSingleOutport(t, outChans, "res"),
		mustSingleOutport(t, outChans, "err"),
	)
	if !ok {
		t.Fatal("handleFileMessage() returned false")
	}

	assertOutputEquals(t, outChans, "res", emptyStruct(), nil)
	if _, err := store.Get(handleID); err == nil {
		t.Fatal("Get() expected error after file_close")
	}
}

func TestFileReadAllHandleReturnsDataAndHandle(t *testing.T) {
	t.Parallel()

	store := runtime.NewFileHandles()
	file := openTempFileWithData(t, "payload")
	handleID := store.Add(file)
	_, _, outChans := newIO(nil, []string{"res", "handle", "err"})

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

	assertOutputEquals(t, outChans, "res", runtime.NewBytesMsg([]byte("payload")), nil)
	assertOutputEquals(t, outChans, "handle", runtime.NewIntMsg(handleID), nil)
	mustCloseHandle(t, store, handleID)
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

func TestFileWriteAllHandleWritesDataAndKeepsHandleOpen(t *testing.T) {
	t.Parallel()

	store := runtime.NewFileHandles()
	file := createTempFile(t)
	handleID := store.Add(file)
	_, _, outChans := newIO(nil, []string{"res", "err"})

	ok := fileWriteAllHandle{handles: store}.handleFileMessage(
		context.Background(),
		runtime.OrderedMsg{Msg: runtime.NewIntMsg(handleID)},
		runtime.OrderedMsg{Msg: runtime.NewBytesMsg([]byte("payload"))},
		mustSingleOutport(t, outChans, "res"),
		mustSingleOutport(t, outChans, "err"),
	)
	if !ok {
		t.Fatal("handleFileMessage() returned false")
	}

	assertOutputEquals(t, outChans, "res", runtime.NewIntMsg(handleID), nil)
	if _, err := store.Get(handleID); err != nil {
		t.Fatalf("Get() after write error = %v", err)
	}
	mustCloseHandle(t, store, handleID)

	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "payload" {
		t.Fatalf("written data = %q, want payload", data)
	}
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

func TestFileStdioComponentsReturnRuntimeHandles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		creator runtime.FuncCreator
		name    string
		want    int64
	}{
		{name: "stdin", creator: fileStdin{}, want: runtime.StdinFileHandleID},
		{name: "stdout", creator: fileStdout{}, want: runtime.StdoutFileHandleID},
		{name: "stderr", creator: fileStderr{}, want: runtime.StderrFileHandleID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			io, inChans, outChans := newIO([]string{"sig"}, []string{"res"})
			handler, err := tt.creator.Create(io, nil)
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
			cancel, done := runHandler(handler)
			defer waitForHandler(t, cancel, done)

			sendInOrder(t, inChans, []string{"sig"}, map[string]runtime.Msg{"sig": emptyStruct()})

			assertOutputEquals(t, outChans, "res", runtime.NewIntMsg(tt.want), nil)
		})
	}
}

func writeTempFile(t *testing.T, data string) string {
	t.Helper()

	path := t.TempDir() + "/file.txt"
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func openTempFileWithData(t *testing.T, data string) *os.File {
	t.Helper()

	path := writeTempFile(t, data)
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	return file
}

func createTempFile(t *testing.T) *os.File {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "file-handle-*.txt")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	return file
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

func mustGetFile(tb testing.TB, store *runtime.FileHandles, handleID int64) *os.File {
	tb.Helper()

	file, err := store.Get(handleID)
	if err != nil {
		tb.Fatalf("Get(%d) error = %v", handleID, err)
	}
	return file
}

func mustCloseHandle(tb testing.TB, store *runtime.FileHandles, handleID int64) {
	tb.Helper()

	if err := store.Close(handleID); err != nil {
		tb.Fatalf("Close(%d) error = %v", handleID, err)
	}
}

func receiveIntResOutput(tb testing.TB, outChans map[string]chan runtime.OrderedMsg) int64 {
	tb.Helper()

	select {
	case got := <-outChans["res"]:
		msg, ok := got.Msg.(runtime.IntMsg)
		if !ok {
			tb.Fatalf("output %q = %T, want runtime.IntMsg", "res", got.Msg)
		}
		return msg.Int()
	case <-time.After(time.Second):
		tb.Fatalf("no int output on %q", "res")
		return 0
	}
}

func waitForHandler(tb testing.TB, cancel context.CancelFunc, done <-chan struct{}) {
	tb.Helper()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		tb.Fatal("handler did not stop")
	}
}

func mustSingleOutport(
	tb testing.TB,
	outChans map[string]chan runtime.OrderedMsg,
	name string,
) runtime.SingleOutport {
	tb.Helper()

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

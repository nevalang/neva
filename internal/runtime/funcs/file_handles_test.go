package funcs

import (
	"io"
	"os"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestFileHandleStoreLifecycle(t *testing.T) {
	t.Parallel()

	tmpFile, err := os.CreateTemp(t.TempDir(), "file-handles-*.txt")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	if _, err := tmpFile.WriteString("hello"); err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatalf("Seek() error = %v", err)
	}
	t.Cleanup(func() { _ = tmpFile.Close() })

	store := newFileHandleStore()
	id := store.Add(tmpFile)
	if id <= stderrFileHandleID {
		t.Fatalf("Add() id = %d, expected dynamic handle > %d", id, stderrFileHandleID)
	}

	gotFile, err := store.Get(id)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	data, err := io.ReadAll(gotFile)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("io.ReadAll() = %q, want %q", string(data), "hello")
	}

	if err := store.Close(id); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if _, err := store.Get(id); err == nil {
		t.Fatal("Get() expected error after close")
	}
	if err := store.Close(id); err == nil {
		t.Fatal("Close() expected error for unknown handle")
	}
}

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

func TestFileHandleStoreHasStdioHandles(t *testing.T) {
	t.Parallel()

	store := newFileHandleStore()

	stdinFile, err := store.Get(stdinFileHandleID)
	if err != nil {
		t.Fatalf("Get(stdin) error = %v", err)
	}
	if stdinFile != os.Stdin {
		t.Fatal("stdin handle does not point to os.Stdin")
	}

	stdoutFile, err := store.Get(stdoutFileHandleID)
	if err != nil {
		t.Fatalf("Get(stdout) error = %v", err)
	}
	if stdoutFile != os.Stdout {
		t.Fatal("stdout handle does not point to os.Stdout")
	}

	stderrFile, err := store.Get(stderrFileHandleID)
	if err != nil {
		t.Fatalf("Get(stderr) error = %v", err)
	}
	if stderrFile != os.Stderr {
		t.Fatal("stderr handle does not point to os.Stderr")
	}

	if err := store.Close(stdinFileHandleID); err == nil {
		t.Fatal("Close(stdin) expected error")
	}
	if err := store.Close(stdoutFileHandleID); err == nil {
		t.Fatal("Close(stdout) expected error")
	}
	if err := store.Close(stderrFileHandleID); err == nil {
		t.Fatal("Close(stderr) expected error")
	}
}

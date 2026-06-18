package runtime

import (
	"io"
	"os"
	"testing"
)

func TestFileHandlesLifecycle(t *testing.T) {
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

	handles := NewFileHandles()
	handleID := handles.Add(tmpFile)
	if handleID <= StderrFileHandleID {
		t.Fatalf("Add() id = %d, expected dynamic handle > %d", handleID, StderrFileHandleID)
	}

	gotFile, err := handles.Get(handleID)
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

	if err := handles.Close(handleID); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if _, err := handles.Get(handleID); err == nil {
		t.Fatal("Get() expected error after close")
	}
	if err := handles.Close(handleID); err == nil {
		t.Fatal("Close() expected error for unknown handle")
	}
}

func TestFileHandlesHasStdioHandles(t *testing.T) {
	t.Parallel()

	handles := NewFileHandles()

	stdinFile, err := handles.Get(StdinFileHandleID)
	if err != nil {
		t.Fatalf("Get(stdin) error = %v", err)
	}
	if stdinFile != os.Stdin {
		t.Fatal("stdin handle does not point to os.Stdin")
	}

	stdoutFile, err := handles.Get(StdoutFileHandleID)
	if err != nil {
		t.Fatalf("Get(stdout) error = %v", err)
	}
	if stdoutFile != os.Stdout {
		t.Fatal("stdout handle does not point to os.Stdout")
	}

	stderrFile, err := handles.Get(StderrFileHandleID)
	if err != nil {
		t.Fatalf("Get(stderr) error = %v", err)
	}
	if stderrFile != os.Stderr {
		t.Fatal("stderr handle does not point to os.Stderr")
	}

	if err := handles.Close(StdinFileHandleID); err == nil {
		t.Fatal("Close(stdin) expected error")
	}
	if err := handles.Close(StdoutFileHandleID); err == nil {
		t.Fatal("Close(stdout) expected error")
	}
	if err := handles.Close(StderrFileHandleID); err == nil {
		t.Fatal("Close(stderr) expected error")
	}
}

func TestFileHandlesAllocatesAcrossShards(t *testing.T) {
	t.Parallel()

	handles := NewFileHandles()
	seenShards := map[int64]struct{}{}

	for range fileHandleShardCount * 2 {
		tmpFile, err := os.CreateTemp(t.TempDir(), "file-handles-sharded-*.txt")
		if err != nil {
			t.Fatalf("CreateTemp() error = %v", err)
		}

		handleID := handles.Add(tmpFile)
		seenShards[handleID%fileHandleShardCount] = struct{}{}
	}

	if len(seenShards) < fileHandleShardCount {
		t.Fatalf("Add() used %d shards, want %d", len(seenShards), fileHandleShardCount)
	}
}

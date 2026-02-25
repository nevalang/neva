package funcs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestFileModeFromRuntimeMsg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		perm    int64
		want    os.FileMode
		wantErr bool
	}{
		{name: "valid", perm: 0o755, want: 0o755},
		{name: "negative", perm: -1, wantErr: true},
		{name: "too_large", perm: maxUint32AsInt64 + 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := fileModeFromRuntimeMsg(runtime.NewIntMsg(tt.perm))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got mode=%v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected mode %v, got %v", tt.want, got)
			}
		})
	}
}

func TestDirEntriesMsg(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filePath := filepath.Join(root, "file.txt")
	dirPath := filepath.Join(root, "subdir")

	if err := os.WriteFile(filePath, []byte("a"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.Mkdir(dirPath, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	list := dirEntriesMsg(entries).List()
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}

	got := map[string]bool{}
	for _, item := range list {
		s := item.Struct()
		got[s.Get("name").Str()] = s.Get("isDir").Bool()
	}

	if isDir, ok := got["file.txt"]; !ok || isDir {
		t.Fatalf("expected file.txt to exist and isDir=false, got ok=%t isDir=%t", ok, isDir)
	}
	if isDir, ok := got["subdir"]; !ok || !isDir {
		t.Fatalf("expected subdir to exist and isDir=true, got ok=%t isDir=%t", ok, isDir)
	}
}

func TestFileInfoMsg(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	filePath := filepath.Join(root, "sample.txt")

	if err := os.WriteFile(filePath, []byte("hello"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}

	msg := fileInfoMsg(info)
	if got := msg.Get("name").Str(); got != "sample.txt" {
		t.Fatalf("expected name sample.txt, got %q", got)
	}
	if got := msg.Get("size").Int(); got != 5 {
		t.Fatalf("expected size 5, got %d", got)
	}
	if got := msg.Get("isDir").Bool(); got {
		t.Fatalf("expected isDir false")
	}
	if got := msg.Get("modTimeUnix").Int(); got <= 0 {
		t.Fatalf("expected modTimeUnix > 0, got %d", got)
	}
}

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

func TestDirEntries(t *testing.T) {
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

	list := dirEntries(entries)
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

//nolint:cyclop,gocognit,gocyclo // This groups OS runtime smoke cases by public std/os area.
func TestOSFilesystemRuntimeFuncs(t *testing.T) {
	t.Run("mkdir creates directory with mode", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "dir")
		runBinaryRuntimeFunc(
			t,
			osMkdir{},
			"path",
			"perm",
			runtime.NewStringMsg(path),
			runtime.NewIntMsg(0o700),
		)

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat: %v", err)
		}
		if !info.IsDir() {
			t.Fatal("mkdir result is not directory")
		}
	})

	t.Run("mkdir missing parent reports error", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing", "dir")
		got := runBinaryRuntimeFuncErr(
			t,
			osMkdir{},
			"path",
			"perm",
			runtime.NewStringMsg(path),
			runtime.NewIntMsg(0o700),
		)
		if got.Struct().Get("text").Str() == "" {
			t.Fatal("mkdir error message is empty")
		}
	})

	t.Run("mkdir_all creates nested directory", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "a", "b")
		runBinaryRuntimeFunc(
			t,
			osMkdirAll{},
			"path",
			"perm",
			runtime.NewStringMsg(path),
			runtime.NewIntMsg(0o755),
		)

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat: %v", err)
		}
		if !info.IsDir() {
			t.Fatal("mkdir_all result is not directory")
		}
	})

	t.Run("read_dir lists file and directory", func(t *testing.T) {
		root := t.TempDir()
		writeFile(t, filepath.Join(root, "file.txt"), "x")
		mkdir(t, filepath.Join(root, "subdir"))

		got := runUnaryRuntimeFunc(t, osReadDir{}, "path", runtime.NewStringMsg(root))
		if got.List().Len() != 2 {
			t.Fatalf("read_dir len = %d, want 2", got.List().Len())
		}
	})

	t.Run("remove removes file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "file.txt")
		writeFile(t, path, "x")

		runUnaryRuntimeFunc(t, osRemove{}, "path", runtime.NewStringMsg(path))
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("Stat after remove error = %v, want not exists", err)
		}
	})

	t.Run("remove_all removes tree", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "tree")
		nested := filepath.Join(root, "nested")
		mkdir(t, nested)
		writeFile(t, filepath.Join(nested, "file.txt"), "x")

		runUnaryRuntimeFunc(t, osRemoveAll{}, "path", runtime.NewStringMsg(root))
		if _, err := os.Stat(root); !os.IsNotExist(err) {
			t.Fatalf("Stat after remove_all error = %v, want not exists", err)
		}
	})

	t.Run("rename moves file", func(t *testing.T) {
		root := t.TempDir()
		oldPath := filepath.Join(root, "old.txt")
		newPath := filepath.Join(root, "new.txt")
		writeFile(t, oldPath, "x")

		runBinaryRuntimeFunc(
			t,
			osRename{},
			"oldPath",
			"newPath",
			runtime.NewStringMsg(oldPath),
			runtime.NewStringMsg(newPath),
		)
		if _, err := os.Stat(newPath); err != nil {
			t.Fatalf("Stat renamed file: %v", err)
		}
	})

	t.Run("stat returns file info", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "file.txt")
		writeFile(t, path, "hello")

		got := runUnaryRuntimeFunc(t, osStat{}, "path", runtime.NewStringMsg(path)).Struct()
		if got.Get("name").Str() != "file.txt" || got.Get("size").Int() != 5 {
			t.Fatalf("stat = %v, want file.txt size 5", got)
		}
	})

	t.Run("lstat returns symlink info", func(t *testing.T) {
		root := t.TempDir()
		target := filepath.Join(root, "target.txt")
		link := filepath.Join(root, "link.txt")
		writeFile(t, target, "hello")
		if err := os.Symlink(target, link); err != nil {
			t.Skipf("symlink unavailable: %v", err)
		}

		got := runUnaryRuntimeFunc(t, osLstat{}, "path", runtime.NewStringMsg(link)).Struct()
		if got.Get("name").Str() != "link.txt" {
			t.Fatalf("lstat name = %q, want link.txt", got.Get("name").Str())
		}
	})

	t.Run("truncate changes file size", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "file.txt")
		writeFile(t, path, "hello")

		runBinaryRuntimeFunc(
			t,
			osTruncate{},
			"path",
			"size",
			runtime.NewStringMsg(path),
			runtime.NewIntMsg(2),
		)

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat: %v", err)
		}
		if info.Size() != 2 {
			t.Fatalf("size = %d, want 2", info.Size())
		}
	})

	t.Run("temp_dir returns os temp dir", func(t *testing.T) {
		got := runSignalRuntimeFunc(t, osTempDir{})
		if got.Str() != os.TempDir() {
			t.Fatalf("temp_dir = %q, want %q", got.Str(), os.TempDir())
		}
	})

	t.Run("mkdir_temp creates directory", func(t *testing.T) {
		got := runBinaryRuntimeFunc(
			t,
			osMkdirTemp{},
			"dir",
			"pattern",
			runtime.NewStringMsg(t.TempDir()),
			runtime.NewStringMsg("neva-*"),
		)
		info, err := os.Stat(got.Str())
		if err != nil {
			t.Fatalf("Stat temp dir: %v", err)
		}
		if !info.IsDir() {
			t.Fatal("mkdir_temp result is not directory")
		}
	})

	t.Run("create_temp creates closed file", func(t *testing.T) {
		got := runBinaryRuntimeFunc(
			t,
			osCreateTemp{},
			"dir",
			"pattern",
			runtime.NewStringMsg(t.TempDir()),
			runtime.NewStringMsg("neva-*.txt"),
		)
		if err := os.WriteFile(got.Str(), []byte("x"), 0o600); err != nil {
			t.Fatalf("WriteFile to temp file: %v", err)
		}
	})
}

func writeFile(t *testing.T, path string, data string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

func mkdir(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", path, err)
	}
}

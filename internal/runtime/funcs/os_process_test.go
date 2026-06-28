package funcs

import (
	"os"
	"strings"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestOSGetwdRuntimeFuncReturnsCurrentDirectory(t *testing.T) {
	want, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	got := runSignalRuntimeFunc(t, osGetwd{})
	if got.Str() != want {
		t.Fatalf("getwd = %q, want %q", got.Str(), want)
	}
}

func TestOSChdirRuntimeFuncChangesCurrentDirectory(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	dir := t.TempDir()
	runUnaryRuntimeFunc(t, osChdir{}, "path", runtime.NewStringMsg(dir))

	got, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd after chdir: %v", err)
	}
	gotInfo, err := os.Stat(got)
	if err != nil {
		t.Fatalf("Stat cwd: %v", err)
	}
	wantInfo, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Stat temp dir: %v", err)
	}
	if !os.SameFile(gotInfo, wantInfo) {
		t.Fatalf("cwd = %q, want same file as %q", got, dir)
	}
}

func TestOSChdirRuntimeFuncReportsMissingDirectory(t *testing.T) {
	got := runUnaryRuntimeFuncErr(
		t,
		osChdir{},
		"path",
		runtime.NewStringMsg(t.TempDir()+"/missing"),
	)
	if !strings.Contains(got.Struct().Get("text").Str(), "os.Chdir") {
		t.Fatalf("chdir error = %v, want os.Chdir", got)
	}
}

func TestOSGetpidRuntimeFuncReturnsProcessID(t *testing.T) {
	got := runSignalRuntimeFunc(t, osGetpid{})
	if got.Int() != int64(os.Getpid()) {
		t.Fatalf("getpid = %d, want %d", got.Int(), os.Getpid())
	}
}

func TestOSGetppidRuntimeFuncReturnsParentProcessID(t *testing.T) {
	got := runSignalRuntimeFunc(t, osGetppid{})
	if got.Int() != int64(os.Getppid()) {
		t.Fatalf("getppid = %d, want %d", got.Int(), os.Getppid())
	}
}

func TestOSHostnameRuntimeFuncReturnsNonEmptyHost(t *testing.T) {
	got := runSignalRuntimeFunc(t, osHostname{})
	if got.Str() == "" {
		t.Fatal("hostname is empty")
	}
}

func TestOSExecutableRuntimeFuncReturnsNonEmptyPath(t *testing.T) {
	got := runSignalRuntimeFunc(t, osExecutable{})
	if got.Str() == "" {
		t.Fatal("executable path is empty")
	}
}

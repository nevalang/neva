package versionmanager

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNormalize(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in       string
		expected string
	}{
		{"0.32.0", "v0.32.0"},
		{"v0.31.1", "v0.31.1"},
		{" V0.30.0 ", "v0.30.0"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()

			got, err := Normalize(tc.in)
			if err != nil {
				t.Fatalf("Normalize(%q) returned error: %v", tc.in, err)
			}

			if got != tc.expected {
				t.Fatalf("Normalize(%q) = %q, expected %q", tc.in, got, tc.expected)
			}
		})
	}
}

func TestUseSetsActiveVersionWithoutDownload(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	baseDir := filepath.Join(tempDir, "neva")

	manager, err := newManager(managerConfig{baseDir: baseDir})
	if err != nil {
		t.Fatalf("newManager returned error: %v", err)
	}

	version, installed, err := manager.Use(context.Background(), "0.33.0", "0.33.0")
	if err != nil {
		t.Fatalf("Use returned error: %v", err)
	}

	if installed {
		t.Fatalf("Use should not install when using bundled version")
	}

	expected := "v0.33.0"
	if version != expected {
		t.Fatalf("Use returned version %q, expected %q", version, expected)
	}

	data, err := os.ReadFile(filepath.Join(baseDir, "active-version"))
	if err != nil {
		t.Fatalf("reading version file: %v", err)
	}

	if strings.TrimSpace(string(data)) != expected {
		t.Fatalf("active version file contains %q, expected %q", strings.TrimSpace(string(data)), expected)
	}
}

func TestUseAcceptsLatestAlias(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	baseDir := filepath.Join(tempDir, "neva")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/nevalang/neva/releases/latest":
			_, _ = w.Write([]byte(`{"tag_name":"v0.33.0"}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	manager, err := newManager(managerConfig{
		baseDir:         baseDir,
		httpClient:      server.Client(),
		apiBaseURL:      server.URL,
		downloadBaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("newManager returned error: %v", err)
	}

	version, installed, err := manager.Use(context.Background(), "latest", "0.33.0")
	if err != nil {
		t.Fatalf("Use returned error: %v", err)
	}

	if version != "v0.33.0" {
		t.Fatalf("Use returned version %q, expected %q", version, "v0.33.0")
	}

	if installed {
		t.Fatalf("Use should not install when latest equals bundled version")
	}
}

func TestMaybeDelegateRunsInstalledVersion(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test relies on POSIX shell")
	}

	tempDir := t.TempDir()
	homeDir := filepath.Join(tempDir, "home")
	t.Setenv("HOME", homeDir)

	manager, err := newManager(managerConfig{baseDir: filepath.Join(homeDir, "neva")})
	if err != nil {
		t.Fatalf("newManager returned error: %v", err)
	}

	if err := manager.SetActiveVersion("v0.30.0"); err != nil {
		t.Fatalf("SetActiveVersion returned error: %v", err)
	}

	binaryDir := filepath.Join(homeDir, "neva", "versions", "v0.30.0")
	if err := os.MkdirAll(binaryDir, 0o755); err != nil {
		t.Fatalf("creating binary dir: %v", err)
	}

	binaryPath := filepath.Join(binaryDir, manager.binaryName)
	script := "#!/bin/sh\necho delegated\n"
	if err := os.WriteFile(binaryPath, []byte(script), 0o755); err != nil {
		t.Fatalf("writing script: %v", err)
	}

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("creating pipe: %v", err)
	}
	defer r.Close()

	os.Stdout = w
	os.Stderr = w

	handled, err := MaybeDelegate([]string{"neva", "version"}, "0.33.0")

	w.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	if err != nil {
		t.Fatalf("MaybeDelegate returned error: %v", err)
	}

	if !handled {
		t.Fatalf("MaybeDelegate should have handled the call")
	}

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("reading pipe: %v", err)
	}

	if !strings.Contains(string(out), "delegated") {
		t.Fatalf("expected delegated binary to run, got output %q", string(out))
	}
}

func TestMaybeDelegateSkipsUseCommand(t *testing.T) {
	t.Parallel()

	handled, err := MaybeDelegate([]string{"neva", "use", "0.32.0"}, "0.32.0")
	if err != nil {
		t.Fatalf("MaybeDelegate returned error: %v", err)
	}

	if handled {
		t.Fatalf("MaybeDelegate should not handle use command")
	}
}

package os

import (
	"errors"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

func TestComputeChecksumForFS(t *testing.T) {
	tests := []struct { //nolint:govet // fieldalignment
		name        string
		setup       func(t *testing.T) fs.FS
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name: "empty_fs",
			setup: func(t *testing.T) fs.FS {
			    t.Helper()
				return fstest.MapFS{}
			},
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", // SHA-256 of empty input
		},
		{
			name: "single_file",
			setup: func(t *testing.T) fs.FS {
			    t.Helper()
				return fstest.MapFS{
					"test.txt": {
						Data: []byte("hello, world"),
					},
				}
			},
			want: "09ca7e4eaa6e8ae9c7d261167129184883644d07dfba7cbfbc4c8a2e08360d5b",
		},
		{
			name: "multiple_files",
			setup: func(t *testing.T) fs.FS {
			    t.Helper()
				return fstest.MapFS{
					"dir1/file1.txt": {
						Data: []byte("file1"),
					},
					"dir2/file2.txt": {
						Data: []byte("file2"),
					},
				}
			},
			want: "bdbfba057752aa11d361f169f579fad1c5205df3d2c7c1a1d419940a7a71180e",
		},
		{
			name: "nested_directories",
			setup: func(t *testing.T) fs.FS {
			    t.Helper()
				return fstest.MapFS{
					"a/b/c/file1.txt": {
						Data: []byte("nested"),
					},
					"a/file2.txt": {
						Data: []byte("file2"),
					},
				}
			},
			want: "33874721bebe9e3516407e03b0a1d9f36daa8db9729e9557acc5ce0b9d743ecb",
		},
		{
			name: "error_reading_file",
			setup: func(t *testing.T) fs.FS {
			    t.Helper()
				return &failingFS{
					MapFS: fstest.MapFS{
						"test.txt": {
							Data: []byte("should fail"),
						},
					},
				}
			},
			wantErr:     true,
			errContains: "read error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := tt.setup(t)
			got, err := ComputeChecksumForFS(fsys)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ComputeChecksumForFS() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Error("expected an error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error message %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if got != tt.want {
				t.Errorf("ComputeChecksumForFS() = %v, want %v", got, tt.want)
			}
		})
	}
}

// failingFS is a test filesystem that fails on Open
type failingFS struct {
	fstest.MapFS
}

func (f *failingFS) Open(name string) (fs.File, error) {
	file, err := f.MapFS.Open(name)
	if err != nil {
		return nil, err
	}
	return &failingFile{File: file}, nil
}

type failingFile struct {
	fs.File
}

func (f *failingFile) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (f *failingFile) Stat() (fs.FileInfo, error) {
	return f.File.Stat()
}

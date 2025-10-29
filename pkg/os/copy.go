package os

import (
	"io"
	"io/fs"
	stdos "os"
	"path/filepath"
)

// use the original permission bits so the copy matches the template file
func CopyFile(src, dst string, mode fs.FileMode) error {
	if err := stdos.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	srcFile, err := stdos.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := stdos.OpenFile(dst, stdos.O_CREATE|stdos.O_WRONLY|stdos.O_TRUNC, mode.Perm())
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Close()
}

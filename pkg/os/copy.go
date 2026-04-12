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

// CopyDir recursively copies a directory tree preserving file modes.
func CopyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}

		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return stdos.MkdirAll(target, 0o755)
		}

		info, statErr := d.Info()
		if statErr != nil {
			return statErr
		}

		return CopyFile(path, target, info.Mode())
	})
}

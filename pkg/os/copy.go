package os

import (
	"fmt"
	"io"
	"io/fs"
	stdos "os"
	"path/filepath"
)

// use the original permission bits so the copy matches the template file
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func CopyFile(src, dst string, mode fs.FileMode) error {
	if err := stdos.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return fmt.Errorf("mkdir for %q: %w", dst, err)
	}

	srcFile, err := stdos.Open(src)
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return fmt.Errorf("open source %q: %w", src, err)
	}
	defer srcFile.Close()

	dstFile, err := stdos.OpenFile(dst, stdos.O_CREATE|stdos.O_WRONLY|stdos.O_TRUNC, mode.Perm())
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return fmt.Errorf("open destination %q: %w", dst, err)
	}
	defer func() { _ = dstFile.Close() }()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return fmt.Errorf("copy %q to %q: %w", src, dst, err)
	}

	//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	if err := dstFile.Close(); err != nil {
		return fmt.Errorf("close destination %q: %w", dst, err)
	}

	return nil
}

// CopyDir recursively copies a directory tree preserving file modes.
func CopyDir(src, dst string) error {
	//nolint:varnamelen,wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk %q: %w", path, err)
		}

		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			return fmt.Errorf("resolve relative path for %q: %w", path, relErr)
		}

		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return stdos.MkdirAll(target, 0o755)
		}

		info, statErr := d.Info()
		if statErr != nil {
			//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			return fmt.Errorf("stat %q: %w", path, statErr)
		}

		return CopyFile(path, target, info.Mode())
	})
}

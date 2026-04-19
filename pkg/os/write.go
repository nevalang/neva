//nolint:all // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package os

import (
	"os"
	"path/filepath"
)

func SaveFilesToDir(dst string, files map[string][]byte) error {
	for path, content := range files {
		filePath := filepath.Join(dst, path)
		dirPath := filepath.Dir(filePath)

		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return err
		}

		// #nosec G306 -- build outputs are intended to be readable
		if err := os.WriteFile(filePath, content, 0o644); err != nil {
			return err
		}
	}

	return nil
}

// FileExists reports whether path exists and points to a regular file.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

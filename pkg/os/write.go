package os

import (
	"fmt"
	"os"
	"path/filepath"
)

func SaveFilesToDir(dst string, files map[string][]byte) error {
	for path, content := range files {
		filePath := filepath.Join(dst, path)
		dirPath := filepath.Dir(filePath)

		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			return fmt.Errorf("mkdir %q: %w", dirPath, err)
		}

		// #nosec G306 -- build outputs are intended to be readable
		if err := os.WriteFile(filePath, content, 0o644); err != nil {
			//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			return fmt.Errorf("write file %q: %w", filePath, err)
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

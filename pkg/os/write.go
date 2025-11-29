package os

import (
	"os"
	"path/filepath"
)

func SaveFilesToDir(dst string, files map[string][]byte) error {
	for path, content := range files {
		filePath := filepath.Join(dst, path)
		dirPath := filepath.Dir(filePath)

		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(filePath, content, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}


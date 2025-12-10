package std

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	nevaos "github.com/nevalang/neva/pkg/os"
)

//go:embed *
var FS embed.FS

// EnsureStdlib ensures the standard library is properly extracted and up-to-date
func EnsureStdlib() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	path := filepath.Join(home, "neva", "std")

	// Compute checksum of the embedded stdlib
	embeddedChecksum, err := nevaos.ComputeChecksumForFS(FS)
	if err != nil {
		return "", fmt.Errorf("compute embedded checksum: %w", err)
	}

	// Read existing checksum if it exists
	existingChecksum, err := readChecksum(path)
	if err != nil {
		return "", fmt.Errorf("read existing checksum: %w", err)
	}

	// If checksums match, return the existing path
	if existingChecksum != "" && existingChecksum == embeddedChecksum {
		return path, nil
	}

	// If we get here, we need to update the stdlib
	// Remove the existing directory if it exists
	if err := os.RemoveAll(path); err != nil {
		return "", fmt.Errorf("remove existing stdlib: %w", err)
	}

	// Create the directory structure
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("create stdlib directory: %w", err)
	}

	// Write all files from the embedded FS
	err = fs.WalkDir(FS, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk error at %s: %w", filePath, err)
		}

		targetPath := filepath.Join(path, filePath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := fs.ReadFile(FS, filePath)
		if err != nil {
			return fmt.Errorf("read embedded file %s: %w", filePath, err)
		}

		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("create directory for %s: %w", targetPath, err)
		}

		return os.WriteFile(targetPath, data, 0644)
	})

	if err != nil {
		return "", fmt.Errorf("write stdlib files: %w", err)
	}

	// Write the new checksum
	if err := writeChecksum(path, embeddedChecksum); err != nil {
		return "", fmt.Errorf("write checksum: %w", err)
	}

	return path, nil
}

// readChecksum reads the checksum from the .checksum file
func readChecksum(stdlibPath string) (string, error) {
	checksumPath := filepath.Join(stdlibPath, ".checksum")
	data, err := os.ReadFile(checksumPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil // No checksum file exists yet
		}
		return "", fmt.Errorf("read checksum: %w", err)
	}
	return string(data), nil
}

// writeChecksum writes the checksum to the .checksum file
func writeChecksum(stdlibPath, checksum string) error {
	checksumPath := filepath.Join(stdlibPath, ".checksum")
	return os.WriteFile(checksumPath, []byte(checksum), 0644)
}

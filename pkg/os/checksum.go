package os

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"sort"
)

// ComputeChecksum computes a checksum of all files in the embedded FS
func ComputeChecksumForFS(filesys fs.FS) (string, error) {
	// First pass: collect all files
	var filenames []string
	err := fs.WalkDir(filesys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			filenames = append(filenames, path)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("walk fs: %w", err)
	}

	// Sort files for consistent ordering
	sort.Strings(filenames)

	// Second pass: hash file contents
	hasher := sha256.New()
	for _, filename := range filenames {
		f, err := filesys.Open(filename)
		if err != nil {
			return "", fmt.Errorf("open %s: %w", filename, err)
		}

		if _, err := io.Copy(hasher, f); err != nil {
			if err := f.Close(); err != nil {
				return "", fmt.Errorf("close: %s: %w", filename, err)
			}
			return "", fmt.Errorf("hash %s: %w", filename, err)
		}

		if err := f.Close(); err != nil {
			return "", fmt.Errorf("close %s: %w", filename, err)
		}
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

package os

import stdos "os"

// FileExists reports whether path exists and points to a non-directory file.
func FileExists(path string) bool {
	info, err := stdos.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

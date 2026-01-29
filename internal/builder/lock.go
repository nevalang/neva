package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// acquireLockFile prevents concurrent dependency downloads by creating a lock file.
// This is used during the build process when downloading dependencies to ensure
// only one build process can download dependencies at a time.
//
// It works by attempting to create a .lock file in the neva home directory. If the file
// already exists (meaning another build process has the lock), it will retry every second
// for up to 60 seconds. Once acquired, it returns a release function that will remove
// the lock file when called.
//
// The lock is necessary because multiple concurrent builds could try to download and write
// the same dependency files simultaneously, which could corrupt the dependency cache.
func acquireLockFile() (release func(), err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filename := filepath.Join(home, "neva", ".lock")

	for i := 0; i < 60; i++ {
		f, err := os.OpenFile(
			filename,
			os.O_CREATE|os.O_EXCL,
			0755,
		)
		if err == nil {
			return func() {
				// must close first, then remove it.
				if err := f.Close(); err != nil {
					panic(err)
				}
				if err := os.Remove(filename); err != nil {
					panic(err)
				}
			}, nil
		}

		if !os.IsExist(err) {
			return nil, fmt.Errorf(
				"unexpected error acquiring neva lock file: %w",
				err,
			)
		}

		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf(
		"maximum retry attempts while acquiring the neva lock file (does %s exist?)",
		filename,
	)
}

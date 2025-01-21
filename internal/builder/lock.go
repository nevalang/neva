package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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
				//must close first, then remove it.
				if err := f.Close(); err != nil {
					panic(err)
				}
				if err := os.Remove(filename); err != nil {
					time.Sleep(0)
					panic(err)
				}
			}, nil
		}

		if !os.IsExist(err) {
			return nil, fmt.Errorf(
				"unexpected error acquiring neva lock file: %v",
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

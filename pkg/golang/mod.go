package golang

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//nolint:gocognit // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func FindModulePath(dst string) (string, error) {
	absDst, err := filepath.Abs(dst)
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return "", err
	}

	dir := absDst
	for {
		goModPath := filepath.Join(dir, "go.mod")
		//nolint:nestif // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		if _, err := os.Stat(goModPath); err == nil {
			modulePath, err := func() (string, error) {
				//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
				f, err := os.Open(goModPath)
				if err != nil {
					//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
					return "", err
				}
				defer f.Close()

				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					line := scanner.Text()
					if after, ok := strings.CutPrefix(line, "module "); ok {
						modName := strings.TrimSpace(after)
						relPath, err := filepath.Rel(dir, absDst)
						if err != nil {
							//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
							return "", err
						}
						if relPath == "." {
							return modName, nil
						}
						return modName + "/" + relPath, nil
					}
				}
				return "", fmt.Errorf("module name not found in %s", goModPath)
			}()
			if err != nil {
				return "", err
			}
			return modulePath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("go.mod not found in %s or parents", dst)
}

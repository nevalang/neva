package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nevalang/neva/pkg"
	cli "github.com/urfave/cli/v2"
)

func newNewCmd(workdir string) *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "Create new Nevalang project",
		Args:  true,
		Action: func(cCtx *cli.Context) error {
			if pathArg := cCtx.Args().First(); pathArg != "" {
				if err := os.Mkdir(pathArg, 0755); err != nil {
					return err
				}
				return createNevaMod(pathArg)
			}
			return createNevaMod(workdir)
		},
	}
}

func createNevaMod(path string) error {
	// Create neva.yml file
	nevaYmlContent := fmt.Sprintf("neva: %s", pkg.Version)
	if err := os.WriteFile(
		filepath.Join(path, "neva.yml"),
		[]byte(nevaYmlContent),
		0644,
	); err != nil {
		return err
	}

	// Create src sub-directory
	srcPath := filepath.Join(path, "src")
	if err := os.Mkdir(srcPath, 0755); err != nil {
		return err
	}

	// Create main.neva file
	mainNevaContent := `import { fmt }

def Main(start any) (stop any) {
	fmt.Println
	---
	:start -> 'Hello, World!' -> println -> :stop
}`

	if err := os.WriteFile(
		filepath.Join(srcPath, "main.neva"),
		[]byte(mainNevaContent),
		0644,
	); err != nil {
		return err
	}

	return nil
}

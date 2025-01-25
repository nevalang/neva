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
			var path string
			if pathArg := cCtx.Args().First(); pathArg != "" {
				path = pathArg
				if err := os.Mkdir(pathArg, 0755); err != nil {
					return err
				}
			} else {
				path = workdir
			}
			if err := createNevaMod(path); err != nil {
				return err
			}
			fmt.Printf("neva module created in %s\n", path)
			return nil
		},
	}
}

func createNevaMod(path string) error {
	// Create neva.yml file
	if err := os.WriteFile(
		filepath.Join(path, "neva.yml"),
		[]byte(fmt.Sprintf("neva: %s", pkg.Version)),
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

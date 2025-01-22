package cli

import (
	"fmt"
	"github.com/nevalang/neva/internal"
	"os"
	"path/filepath"

	"github.com/nevalang/neva/pkg"
	"github.com/urfave/cli/v2"
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
				return createNevaMod(pathArg, "hello-world")
			}
			return createNevaMod(workdir, "hello-world")
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "template",
				Usage:    "Specify the template for creating new projects.",
				Required: false,

				Action: func(context *cli.Context, templateName string) error {
					if templateName == "" {
						templateName = "hello-world"
					}
					pathArg := context.Args().First()
					if err := os.Mkdir(pathArg, 0755); err != nil {
						return err
					}
					return createNevaMod(pathArg, templateName)

				},
			},
		},
	}
}

func createNevaMod(path string, templateName string) error {
	// Create neva.yml file
	nevaYmlContent := fmt.Sprintf("neva: %s", pkg.Version)
	if err := os.WriteFile(filepath.Join(path, "neva.yml"), []byte(nevaYmlContent), 0644); err != nil {
		return err
	}

	// Create src sub-directory
	srcPath := filepath.Join(path, "src")
	if err := os.Mkdir(srcPath, 0755); err != nil {
		return err
	}

	mainNevaContent, err := readTemplate(templateName + ".neva")
	if err != nil {
		return err
	}
	if err := os.WriteFile(
		filepath.Join(srcPath, "main.neva"),
		[]byte(mainNevaContent),
		0644,
	); err != nil {
		return err
	}
	return nil
}
func readTemplate(templateName string) (string, error) {
	fileData, err := internal.Efs.ReadFile(fmt.Sprintf("templates/%s", templateName))
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no such template: %s, available templates see: https://github.com/nevalang/neva/tree/main/internal/templates", templateName)
		}
		return "", fmt.Errorf("failed to read template file %s: %w", templateName, err)
	}
	return string(fileData), nil
}

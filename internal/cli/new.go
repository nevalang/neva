package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gitlib "github.com/go-git/go-git/v5"

	"github.com/nevalang/neva/pkg"
	nevaGit "github.com/nevalang/neva/pkg/git"
	nevaos "github.com/nevalang/neva/pkg/os"
	cli "github.com/urfave/cli/v2"
)

func newNewCmd() *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "Create new neva project",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "Template repository, optionally suffixed with a revision",
			},
		},
		Action: func(cCtx *cli.Context) error {
			pathArg := cCtx.Args().First()
			if pathArg == "" {
				return errors.New("path argument is required")
			}

			template := cCtx.String("template")

			if template != "" {
				if err := scaffoldFromTemplate(pathArg, template); err != nil {
					return err
				}
			} else {
				if err := os.MkdirAll(pathArg, 0o755); err != nil {
					return err
				}
				if err := createNevaMod(pathArg); err != nil {
					return err
				}
			}

			fmt.Printf("neva module created in %s\n", pathArg)
			return nil
		},
	}
}

const mainNevaContent = `import {
	fmt
	runtime
}

// main prints a greeting and propagates failures to the runtime panic node.
def Main(start any) (stop any) {
	println fmt.Println<string>
	panic runtime.Panic
	---
	:start -> 'Hello, World!' -> println
	println:res -> :stop
	println:err -> panic
}`

func createNevaMod(path string) error {
	// Create neva.yml file
	// #nosec G306 -- project file is intended to be readable
	if err := os.WriteFile(
		filepath.Join(path, "neva.yml"),
		fmt.Appendf(nil, "neva: %s", pkg.Version),
		0o644,
	); err != nil {
		return err
	}

	// Create src sub-directory
	srcPath := filepath.Join(path, "src")
	if err := os.Mkdir(srcPath, 0o755); err != nil {
		return err
	}

	// Create main.neva file
	// #nosec G306 -- new files are intended to be readable
	return os.WriteFile(
		filepath.Join(srcPath, "main.neva"),
		[]byte(mainNevaContent),
		0o644,
	)
}

func scaffoldFromTemplate(path string, template string) error {
	spec, err := nevaGit.ParseRepoSpec(template)
	if err != nil {
		return err
	}

	if spec.IsLocal() {
		return fmt.Errorf("local templates are not supported, use a remote repository URL")
	}

	if err := ensureEmptyDir(path); err != nil {
		return fmt.Errorf("prepare target directory: %w", err)
	}

	cloneDir, err := os.MkdirTemp("", "neva-template-*")
	if err != nil {
		return fmt.Errorf("create temporary directory: %w", err)
	}
	defer os.RemoveAll(cloneDir)

	repo, err := gitlib.PlainClone(cloneDir, false, &gitlib.CloneOptions{
		URL: spec.CloneURL(),
	})
	if err != nil {
		return fmt.Errorf("clone template: %w", err)
	}

	if spec.Revision != "" {
		if err := nevaGit.Checkout(repo, spec.Revision); err != nil {
			return fmt.Errorf("checkout template revision: %w", err)
		}
	} else {
		if err := nevaGit.Checkout(repo, "HEAD"); err != nil {
			return fmt.Errorf("checkout template default revision: %w", err)
		}
	}

	if err := copyDir(cloneDir, path); err != nil {
		return fmt.Errorf("copy template contents: %w", err)
	}

	return nil
}

func ensureEmptyDir(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0o755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s exists and is not a directory", path)
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("directory %s is not empty", path)
	}
	return nil
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if rel == ".git" || strings.HasPrefix(rel, ".git"+string(os.PathSeparator)) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		return nevaos.CopyFile(path, target, info.Mode())
	})
}

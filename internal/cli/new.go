package cli

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gitlib "github.com/go-git/go-git/v5"

	"github.com/nevalang/neva/pkg"
	nevaGit "github.com/nevalang/neva/pkg/git"
	cli "github.com/urfave/cli/v2"
)

func newNewCmd(workdir string) *cli.Command {
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
			path := workdir
			if pathArg := cCtx.Args().First(); pathArg != "" {
				path = pathArg
			}

			template := cCtx.String("template")
			if template != "" {
				spec, err := parseTemplateSpec(template)
				if err != nil {
					return err
				}

				if err := scaffoldFromTemplate(path, spec); err != nil {
					return err
				}
			} else {
				if pathArg := cCtx.Args().First(); pathArg != "" {
					if err := os.Mkdir(pathArg, 0o755); err != nil {
						return err
					}
				}
				if err := createNevaMod(path); err != nil {
					return err
				}
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
	mainNevaContent := `import { fmt, runtime }

# Main prints a greeting and propagates failures to the runtime panic node.
def Main(start any) (stop any) {
println fmt.Println<string>
panic runtime.Panic
---
:start -> 'Hello, World!' -> println
println:res -> :stop
println:err -> panic
}`

	if err := os.WriteFile(
		filepath.Join(srcPath, "main.neva"),
		[]byte(mainNevaContent),
		0o644, // new files should be writable by the owner and readable by others
	); err != nil {
		return err
	}

	return nil
}

type templateSpec struct {
	Source   string
	Revision string
}

func parseTemplateSpec(raw string) (templateSpec, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return templateSpec{}, errors.New("template must not be empty")
	}

	spec := templateSpec{Source: trimmed}

	if hashIdx := strings.LastIndex(trimmed, "#"); hashIdx != -1 {
		spec.Source = strings.TrimSpace(trimmed[:hashIdx])
		spec.Revision = strings.TrimSpace(trimmed[hashIdx+1:])
	} else if atIdx := strings.LastIndex(trimmed, "@"); atIdx != -1 {
		slashIdx := strings.LastIndex(trimmed, "/")
		if slashIdx < atIdx {
			spec.Source = strings.TrimSpace(trimmed[:atIdx])
			spec.Revision = strings.TrimSpace(trimmed[atIdx+1:])
		}
	}

	if spec.Source == "" {
		return templateSpec{}, errors.New("template repository must not be empty")
	}

	return spec, nil
}

func scaffoldFromTemplate(path string, spec templateSpec) error {
	if path == "" {
		return errors.New("target path must not be empty")
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
		URL: spec.Source,
	})
	if err != nil {
		return fmt.Errorf("clone template: %w", err)
	}

	if spec.Revision != "" {
		if err := nevaGit.Checkout(repo, spec.Revision); err != nil {
			return fmt.Errorf("checkout template revision: %w", err)
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

		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Use the original permission bits so the copy matches the template file.
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode.Perm())
	if err != nil {
		return err
	}
	defer func() {
		_ = dstFile.Close()
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Close()
}

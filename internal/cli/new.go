package cli

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/nevalang/neva/pkg"
	cli "github.com/urfave/cli/v2"
)

func newNewCmd(workdir string) *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "Create new Nevalang project",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "Git repository to use as a template",
			},
			&cli.StringFlag{
				Name:  "template-ref",
				Usage: "Git reference (branch, tag, or commit) to checkout for the template",
			},
		},
		Action: func(cCtx *cli.Context) error {
			path := workdir
			if pathArg := cCtx.Args().First(); pathArg != "" {
				path = pathArg
			}

			template := cCtx.String("template")
			if template != "" {
				if err := scaffoldFromTemplate(path, template, cCtx.String("template-ref")); err != nil {
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
	mainNevaContent := `import { fmt }

def Main(start any) (stop any) {
fmt.Println
---
:start -> 'Hello, World!' -> println
[println:res, println:err] -> :stop
}`

	if err := os.WriteFile(
		filepath.Join(srcPath, "main.neva"),
		[]byte(mainNevaContent),
		0o644,
	); err != nil {
		return err
	}

	return nil
}

func scaffoldFromTemplate(path, template, ref string) error {
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

	repo, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
		URL: normalizeTemplateURL(template),
	})
	if err != nil {
		return fmt.Errorf("clone template: %w", err)
	}

	if ref != "" {
		if err := checkoutTemplateRef(repo, ref); err != nil {
			return err
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

func checkoutTemplateRef(repo *git.Repository, ref string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("open worktree: %w", err)
	}

	if strings.HasPrefix(ref, "refs/") {
		if err := worktree.Checkout(&git.CheckoutOptions{Branch: plumbing.ReferenceName(ref)}); err == nil {
			return nil
		}
	}

	if err := worktree.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(ref)}); err == nil {
		return nil
	}

	tagRef, err := repo.Reference(plumbing.NewTagReferenceName(ref), true)
	if err == nil {
		if err := worktree.Checkout(&git.CheckoutOptions{Hash: tagRef.Hash()}); err == nil {
			return nil
		}
	}

	hash := plumbing.NewHash(ref)
	if hash != plumbing.ZeroHash {
		if _, err := repo.CommitObject(hash); err == nil {
			if err := worktree.Checkout(&git.CheckoutOptions{Hash: hash}); err == nil {
				return nil
			}
		}
	}

	return fmt.Errorf("template reference %q not found", ref)
}

func normalizeTemplateURL(template string) string {
	if template == "" {
		return template
	}
	if strings.Contains(template, "://") || strings.HasPrefix(template, "git@") {
		return template
	}
	if filepath.IsAbs(template) {
		return template
	}
	if strings.HasPrefix(template, ".") {
		return template
	}
	return "https://" + template
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

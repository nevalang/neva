package builder

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/pkg"
	"github.com/nevalang/neva/std"
)

type Builder struct {
	manifestParser ManifestParser
	thirdPartyPath string
	stdLibPath     string
}

type ManifestParser interface {
	ParseManifest(raw []byte) (src.ModuleManifest, error)
}

func (b Builder) Build(
	ctx context.Context,
	wd string,
) (compiler.RawBuild, string, *compiler.Error) {
	// load entry module from disk
	entryMod, entryModRootPath, err := b.LoadModuleByPath(ctx, wd)
	if err != nil {
		return compiler.RawBuild{}, "", &compiler.Error{
			Err: fmt.Errorf("build entry mod: %w", err),
		}
	}

	// inject stdlib dep to entry module
	stdModRef := src.ModuleRef{
		Path:    "std",
		Version: pkg.Version,
	}
	entryMod.Manifest.Deps["std"] = stdModRef

	// inject entry mod into the build
	mods := map[src.ModuleRef]compiler.RawModule{
		{Path: "@"}: entryMod,
	}

	// load stdlib module
	stdMod, _, err := b.LoadModuleByPath(ctx, b.stdLibPath)
	if err != nil {
		return compiler.RawBuild{}, "", &compiler.Error{
			Err: fmt.Errorf("build stdlib mod: %w", err),
		}
	}

	// inject stdlib module to build
	mods[stdModRef] = stdMod

	release, err := acquireLockFile()
	if err != nil {
		return compiler.RawBuild{}, "", &compiler.Error{
			Err: fmt.Errorf("failed to acquire lock file: %w", err),
		}
	}
	defer release()

	q := newQueue(entryMod.Manifest.Deps)

	for !q.empty() {
		depModRef := q.dequeue()

		if _, ok := mods[depModRef]; ok {
			continue
		}

		depWD, _, err := b.downloadDep(depModRef)
		if err != nil {
			return compiler.RawBuild{}, "", &compiler.Error{
				Err: fmt.Errorf("download dep: %w", err),
			}
		}

		depMod, _, err := b.LoadModuleByPath(ctx, depWD)
		if err != nil {
			return compiler.RawBuild{}, "", &compiler.Error{
				Err: fmt.Errorf("build dep mod: %w", err),
			}
		}

		// inject stdlib dep into every downloaded dep mod
		depMod.Manifest.Deps["std"] = src.ModuleRef{
			Path:    "std",
			Version: pkg.Version,
		}

		mods[depModRef] = depMod

		q.enqueue(depMod.Manifest.Deps)
	}

	return compiler.RawBuild{
		EntryModRef: src.ModuleRef{Path: "@"},
		Modules:     mods,
	}, entryModRootPath, nil
}

func getThirdPartyPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, "neva", "deps")

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}

	return path, nil
}

// rewriteStdlibOntoDisk erases stdlib on the disk if it's there and writes it again.
func rewriteStdlibOntoDisk() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, "neva", "std")

	// TODO replace this dirty hack with good solution of https://github.com/nevalang/neva/issues/563
	err = os.RemoveAll(path)
	if err != nil {
		return "", err
	}

	// _, err = os.Stat(path)
	// if err == nil {
	// 	return path, nil
	// }

	// if !os.IsNotExist(err) {
	// 	return "", err
	// }

	// Inject missing stdlib files into user's home directory
	stdFS := std.FS
	err = fs.WalkDir(stdFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(stdFS, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(home, "neva", "std", path)
		dir := filepath.Dir(targetPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
		err = os.WriteFile(targetPath, data, 0644)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return path, nil
}

func New(parser ManifestParser) (Builder, error) {
	thirdParty, err := getThirdPartyPath()
	if err != nil {
		return Builder{}, err
	}

	stdlibPath, err := rewriteStdlibOntoDisk()
	if err != nil {
		return Builder{}, err
	}

	return Builder{
		manifestParser: parser,
		stdLibPath:     stdlibPath,
		thirdPartyPath: thirdParty,
	}, nil
}

func MustNew(parser ManifestParser) Builder {
	b, err := New(parser)
	if err != nil {
		panic(err)
	}
	return b
}

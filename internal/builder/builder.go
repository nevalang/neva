package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	"github.com/nevalang/neva/pkg"
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
			Message: "build entry mod: " + err.Error(),
		}
	}

	// inject stdlib dep to entry module
	if _, ok := entryMod.Manifest.Deps["std"]; ok {
		return compiler.RawBuild{}, "", &compiler.Error{
			Message: "entry module cannot depend on 'std' explicitly; it is injected automatically",
		}
	}

	stdModRef := core.ModuleRef{
		Path:    "std",
		Version: pkg.Version,
	}
	entryMod.Manifest.Deps["std"] = stdModRef

	// inject entry mod into the build
	mods := map[core.ModuleRef]compiler.RawModule{
		{Path: "@"}: entryMod,
	}

	// load stdlib module
	stdMod, _, err := b.LoadModuleByPath(ctx, b.stdLibPath)
	if err != nil {
		return compiler.RawBuild{}, "", &compiler.Error{
			Message: "build stdlib mod: " + err.Error(),
		}
	}

	// inject stdlib module to build
	mods[stdModRef] = stdMod

	release, err := acquireLockFile()
	if err != nil {
		return compiler.RawBuild{}, "", &compiler.Error{
			Message: "failed to acquire lock file: " + err.Error(),
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
				Message: "download dep: " + err.Error(),
			}
		}

		depMod, _, err := b.LoadModuleByPath(ctx, depWD)
		if err != nil {
			return compiler.RawBuild{}, "", &compiler.Error{
				Message: "build dep mod: " + err.Error(),
			}
		}

		// inject stdlib dep into every downloaded dep mod
		depMod.Manifest.Deps["std"] = core.ModuleRef{
			Path:    "std",
			Version: pkg.Version,
		}

		mods[depModRef] = depMod

		q.enqueue(depMod.Manifest.Deps)
	}

	return compiler.RawBuild{
		EntryModRef: core.ModuleRef{Path: "@"},
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

func New(parser ManifestParser) (Builder, error) {
	thirdParty, err := getThirdPartyPath()
	if err != nil {
		return Builder{}, err
	}

	// Use EnsureStdlib to handle stdlib extraction with checksum validation
	stdlibPath, err := ensureStdlib()
	if err != nil {
		return Builder{}, fmt.Errorf("ensure stdlib: %w", err)
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

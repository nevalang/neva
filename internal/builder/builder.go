package builder

import (
	"context"
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg"
	src "github.com/nevalang/neva/pkg/sourcecode"
	"github.com/nevalang/neva/std"
)

type Builder struct {
	manifestParser ManifestParser // parser is needed to parse manifest files
	thirdPartyPath string         // path to third-party modules
}

type ManifestParser interface {
	ParseManifest(raw []byte) (src.ModuleManifest, error)
}

func (p Builder) Build( //nolint:funlen
	ctx context.Context,
	workdir string,
) (compiler.RawBuild, *compiler.Error) {
	// load user's module from disk
	entryMod, err := p.LoadModuleByPath(ctx, os.DirFS(workdir))
	if err != nil {
		return compiler.RawBuild{}, &compiler.Error{
			Err: fmt.Errorf("build entry mod: %w", err),
		}
	}

	// inject stdlib dep to user's module
	entryMod.Manifest.Deps["std"] = src.ModuleRef{
		Path:    "std",
		Version: pkg.Version,
	}

	// TODO use embedded fs for stdlib
	stdMod, err := p.LoadModuleByPath(ctx, std.FS)
	if err != nil {
		return compiler.RawBuild{}, &compiler.Error{
			Err: fmt.Errorf("build stdlib mod: %w", err),
		}
	}

	mods := map[src.ModuleRef]compiler.RawModule{
		{Path: "@"}:                         entryMod,
		{Path: "std", Version: pkg.Version}: stdMod,
	}

	q := newQueue(entryMod.Manifest.Deps)

	for !q.empty() {
		depModRef := q.dequeue()

		if _, ok := mods[depModRef]; ok {
			continue
		}

		depPath, err := p.downloadDep(depModRef)
		if err != nil {
			return compiler.RawBuild{}, &compiler.Error{
				Err: fmt.Errorf("download dep: %w", err),
			}
		}

		depMod, err := p.LoadModuleByPath(ctx, os.DirFS(depPath))
		if err != nil {
			return compiler.RawBuild{}, &compiler.Error{
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
	}, nil
}

func New(thirdParty string, parser ManifestParser) Builder {
	return Builder{
		thirdPartyPath: thirdParty,
		manifestParser: parser,
	}
}

package builder

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

type Builder struct {
	stdLibLocation     string // path to standart library module
	thirdPartyLocation string // path to third-party modules
	parser             Parser // parser is needed to parse manifest files
}

type Parser interface {
	ParseManifest(raw []byte) (src.ModuleManifest, error)
}

func (p Builder) Build( //nolint:funlen
	ctx context.Context,
	workdir string,
) (compiler.RawBuild, *compiler.Error) {
	// load user's module from disk
	entryMod, err := p.LoadModuleByPath(ctx, workdir)
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
	stdMod, err := p.LoadModuleByPath(ctx, p.stdLibLocation)
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

		depMod, err := p.LoadModuleByPath(ctx, depPath)
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

func (p Builder) Install(ctx context.Context, depModRef src.ModuleRef, workdir string) error {
	manifest, err := p.retrieveManifest(workdir)
	if err != nil {
		return err
	}

	if _, err := p.downloadDep(depModRef); err != nil {
		return err
	}

	manifest.Deps[depModRef.Path] = depModRef

	return nil
}

func New(
	stdlibPath string,
	thirdpartyPath string,
	parser Parser,
) Builder {
	return Builder{
		stdLibLocation:     stdlibPath,
		thirdPartyLocation: thirdpartyPath,
		parser:             parser,
	}
}

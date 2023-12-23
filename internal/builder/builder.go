package builder

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

type Builder struct {
	stdLibLocation     string // path to standart library module
	thirdPartyLocation string // path to third-party modules
	parser             Parser // parser is needed to parse manifest files
}

type Parser interface {
	ParseManifest(raw []byte) (src.ModuleManifest, error)
}

func (b Builder) Build(ctx context.Context, workdir string) (compiler.RawBuild, error) {
	mods := map[src.ModuleRef]compiler.RawModule{}

	entryMod, err := b.buildModule(ctx, workdir)
	if err != nil {
		return compiler.RawBuild{}, fmt.Errorf("build entry mod: %w", err)
	}
	mods[src.ModuleRef{Name: "entry"}] = entryMod

	q := NewQueue(entryMod.Manifest.Deps)

	for !q.Empty() {
		depModRef := q.Dequeue()

		if _, ok := mods[depModRef]; ok {
			continue
		}

		depPath, err := b.downloadDep(depModRef)
		if err != nil {
			return compiler.RawBuild{}, fmt.Errorf("download dep: %w", err)
		}

		depMod, err := b.buildModule(ctx, depPath)
		if err != nil {
			return compiler.RawBuild{}, fmt.Errorf("build entry mod: %w", err)
		}

		mods[depModRef] = depMod

		q.Enqueue(depMod.Manifest.Deps)
	}

	stdMod, err := b.buildModule(ctx, b.stdLibLocation)
	if err != nil {
		return compiler.RawBuild{}, fmt.Errorf("build stdlib mod: %w", err)
	}
	mods[src.ModuleRef{Name: "std"}] = stdMod

	return compiler.RawBuild{
		EntryModRef: src.ModuleRef{Name: "entry"},
		Modules:     mods,
	}, nil
}

func MustNew(
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

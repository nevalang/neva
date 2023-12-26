package pkgmanager

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

type Manager struct {
	stdLibLocation     string // path to standart library module
	thirdPartyLocation string // path to third-party modules
	parser             Parser // parser is needed to parse manifest files
}

type Parser interface {
	ParseManifest(raw []byte) (src.ModuleManifest, error)
}

func (p Manager) Build(ctx context.Context, workdir string) (compiler.RawBuild, error) {
	mods := map[src.ModuleRef]compiler.RawModule{}

	entryMod, err := p.buildModule(ctx, workdir)
	if err != nil {
		return compiler.RawBuild{}, fmt.Errorf("build entry mod: %w", err)
	}
	mods[src.ModuleRef{Path: "entry"}] = entryMod

	q := newQueue(entryMod.Manifest.Deps)

	for !q.empty() {
		depModRef := q.dequeue()

		if _, ok := mods[depModRef]; ok {
			continue
		}

		depPath, err := p.downloadDep(depModRef)
		if err != nil {
			return compiler.RawBuild{}, fmt.Errorf("download dep: %w", err)
		}

		depMod, err := p.buildModule(ctx, depPath)
		if err != nil {
			return compiler.RawBuild{}, fmt.Errorf("build entry mod: %w", err)
		}

		mods[depModRef] = depMod

		q.enqueue(depMod.Manifest.Deps)
	}

	stdMod, err := p.buildModule(ctx, p.stdLibLocation)
	if err != nil {
		return compiler.RawBuild{}, fmt.Errorf("build stdlib mod: %w", err)
	}
	mods[src.ModuleRef{Path: "std"}] = stdMod

	return compiler.RawBuild{
		EntryModRef: src.ModuleRef{Path: "entry"},
		Modules:     mods,
	}, nil
}

func (p Manager) Install(ctx context.Context, depModRef src.ModuleRef, workdir string) error {
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
) Manager {
	return Manager{
		stdLibLocation:     stdlibPath,
		thirdPartyLocation: thirdpartyPath,
		parser:             parser,
	}
}

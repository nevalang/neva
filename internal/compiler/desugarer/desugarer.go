package desugarer

import (
	"maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

// FIXME since desugarer operates before analyzer
// we must inject std mod dep and builtin import before we analyze
// i.e not in desugarer
// TODO consider changing the strategy to previous one where
// scope.Entity() know about stdlib and stuff
// it should be enough to use "builtin" in resolution
// but this could broke other virtual entities we create

// Desugarer does the following:
// 1. Replaces const ref senders with normal nodes that uses Const component with compiler directive;
// 2. Inserts void nodes and connections for every unused outport in the program;
// 3. Replaces struct selectors with chain of struct selector nodes.
type Desugarer struct{}

func (d Desugarer) Desugar(build src.Build) (src.Build, *compiler.Error) {
	desugaredMods := make(map[src.ModuleRef]src.Module, len(build.Modules))

	for modRef := range build.Modules {
		desugaredMod, err := d.desugarModule(build, modRef)
		if err != nil {
			return src.Build{}, compiler.Error{
				Location: &src.Location{ModRef: modRef},
			}.Merge(err)
		}
		desugaredMods[modRef] = desugaredMod
	}

	return src.Build{
		EntryModRef: build.EntryModRef,
		Modules:     desugaredMods,
	}, nil
}

func (d Desugarer) desugarModule(build src.Build, modRef src.ModuleRef) (src.Module, *compiler.Error) {
	mod := build.Modules[modRef]

	// create manifest copy with std module dependency
	desugaredManifest := src.ModuleManifest{
		WantCompilerVersion: mod.Manifest.WantCompilerVersion,
		Deps:                make(map[string]src.ModuleRef, len(mod.Manifest.Deps)+1),
	}
	maps.Copy(desugaredManifest.Deps, mod.Manifest.Deps)
	desugaredManifest.Deps["std"] = src.ModuleRef{Path: "std", Version: "0.0.1"}

	// copy all modules but replace manifest in current one
	modsCopy := maps.Clone(build.Modules)
	modsCopy[modRef] = src.Module{
		Manifest: desugaredManifest,
		Packages: mod.Packages,
	}

	// create new build with patched modeles (current module have patched manifest with std dependency)
	build = src.Build{
		EntryModRef: modRef,
		Modules:     modsCopy,
	}

	desugaredPkgs := make(map[string]src.Package, len(mod.Packages))

	for pkgName, pkg := range mod.Packages {
		scope := src.Scope{
			Build: build, // it's important to patch build before desugar package so we can resolve references to std
			Location: src.Location{
				ModRef:  modRef,
				PkgName: pkgName,
			},
		}

		desugaredPkg, err := d.desugarPkg(pkg, scope)
		if err != nil {
			return src.Module{}, compiler.Error{
				Location: &src.Location{PkgName: pkgName},
			}.Merge(err)
		}

		desugaredPkgs[pkgName] = desugaredPkg
	}

	return src.Module{
		Manifest: desugaredManifest,
		Packages: desugaredPkgs,
	}, nil
}

func (d Desugarer) desugarPkg(pkg src.Package, scope src.Scope) (src.Package, *compiler.Error) {
	desugaredPkgs := make(src.Package, len(pkg))

	for fileName, file := range pkg {
		newScope := scope.WithLocation(src.Location{
			ModRef:   scope.Location.ModRef,
			PkgName:  scope.Location.PkgName,
			FileName: fileName,
		})

		desugaredFile, err := d.desugarFile(file, newScope)
		if err != nil {
			return nil, compiler.Error{
				Location: &src.Location{FileName: fileName},
			}.Merge(err)
		}

		desugaredPkgs[fileName] = desugaredFile
	}

	return desugaredPkgs, nil
}

// desugarFile injects import of std/builtin into every pkg's file and desugares it's every entity
func (d Desugarer) desugarFile(file src.File, scope src.Scope) (src.File, *compiler.Error) {
	desugaredEntities := make(map[string]src.Entity, len(file.Entities))

	for entityName, entity := range file.Entities {
		entityResult, err := d.desugarEntity(entity, scope)
		if err != nil {
			return src.File{}, compiler.Error{
				Meta: entity.Meta(),
			}.Merge(err)
		}

		desugaredEntities[entityName] = entityResult.entity

		for name, entityToInsert := range entityResult.entitiesToInsert {
			desugaredEntities[name] = entityToInsert
		}
	}

	desugaredImports := maps.Clone(file.Imports)
	if desugaredImports == nil {
		desugaredImports = map[string]src.Import{}
	}

	desugaredImports["builtin"] = src.Import{
		ModuleName: "std",
		PkgName:    "builtin",
	}

	return src.File{
		Imports:  desugaredImports,
		Entities: desugaredEntities,
	}, nil
}

type desugarEntityResult struct {
	entity           src.Entity
	entitiesToInsert map[string]src.Entity
}

func (d Desugarer) desugarEntity(entity src.Entity, scope src.Scope) (desugarEntityResult, *compiler.Error) {
	if entity.Kind != src.ComponentEntity {
		return desugarEntityResult{
			entity: entity,
		}, nil
	}

	componentResult, err := d.desugarComponent(entity.Component, scope)
	if err != nil {
		return desugarEntityResult{}, compiler.Error{Meta: &entity.Component.Meta}.Merge(err)
	}

	entitiesToInsert := make(map[string]src.Entity, len(componentResult.constsToInsert))
	for constName, constant := range componentResult.constsToInsert {
		entitiesToInsert[constName] = src.Entity{
			IsPublic: false,
			Kind:     src.ConstEntity,
			Const:    constant,
		}
	}

	return desugarEntityResult{
		entitiesToInsert: entitiesToInsert,
		entity: src.Entity{
			IsPublic:  entity.IsPublic,
			Kind:      entity.Kind,
			Component: componentResult.component,
		},
	}, nil
}

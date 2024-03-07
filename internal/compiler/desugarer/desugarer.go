package desugarer

import (
	"maps"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

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
			return src.Build{},
				compiler.Error{
					Location: &src.Location{
						ModRef: modRef,
					},
				}.Wrap(err)
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
		LanguageVersion: mod.Manifest.LanguageVersion,
		Deps:            make(map[string]src.ModuleRef, len(mod.Manifest.Deps)+1),
	}
	maps.Copy(desugaredManifest.Deps, mod.Manifest.Deps)
	desugaredManifest.Deps["std"] = src.ModuleRef{Path: "std", Version: pkg.Version}

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
			}.Wrap(err)
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
			}.Wrap(err)
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
			}.Wrap(err)
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

	desugaredImports["builtin"] = src.Import{ // inject std/builtin import
		Module:  "std",
		Package: "builtin",
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
	if entity.Kind != src.ComponentEntity && entity.Kind != src.ConstEntity {
		return desugarEntityResult{entity: entity}, nil
	}

	if entity.Kind == src.ConstEntity {
		desugaredConst, err := d.handleConst(entity.Const)
		if err != nil {
			return desugarEntityResult{}, compiler.Error{Meta: &entity.Component.Meta}.Wrap(err)
		}

		return desugarEntityResult{
			entity: src.Entity{
				IsPublic: entity.IsPublic,
				Kind:     entity.Kind,
				Const:    desugaredConst,
			},
		}, nil
	}

	componentResult, err := d.handleComponent(entity.Component, scope)
	if err != nil {
		return desugarEntityResult{}, compiler.Error{Meta: &entity.Component.Meta}.Wrap(err)
	}

	return desugarEntityResult{
		entitiesToInsert: componentResult.virtualEntities,
		entity: src.Entity{
			IsPublic:  entity.IsPublic,
			Kind:      entity.Kind,
			Component: componentResult.desugaredComponent,
		},
	}, nil
}

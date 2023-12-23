package desugarer

import (
	"maps"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

type Desugarer struct{}

// TODO add sugar for manifests, insert std dependency

func (d Desugarer) Desugar(build src.Build) (src.Build, error) {
	desugaredMods := make(map[src.ModuleRef]src.Module, len(build.Modules))

	for modRef := range build.Modules {
		desugaredMod, err := d.desugarModule(build, modRef)
		if err != nil {
			return src.Build{}, err
		}
		desugaredMods[modRef] = desugaredMod
	}

	return src.Build{
		EntryModRef: build.EntryModRef,
		Modules:     desugaredMods,
	}, nil
}

func (d Desugarer) desugarModule(build src.Build, modRef src.ModuleRef) (src.Module, error) {
	mod := build.Modules[modRef]
	desugaredPkgs := make(map[string]src.Package, len(mod.Packages))

	for pkgName, pkg := range mod.Packages {
		scope := src.Scope{
			Build: build,
			Location: src.Location{
				ModRef:  modRef,
				PkgName: pkgName,
			},
		}

		desugaredPkg, err := d.desugarPkg(pkg, scope)
		if err != nil {
			return src.Module{}, err
		}

		desugaredPkgs[pkgName] = desugaredPkg
	}

	desugaredManifest := src.ModuleManifest{
		WantCompilerVersion: mod.Manifest.WantCompilerVersion,
		Deps:                make(map[string]src.ModuleRef, len(mod.Manifest.Deps)+1),
	}
	maps.Copy(desugaredManifest.Deps, mod.Manifest.Deps)
	desugaredManifest.Deps["std"] = src.ModuleRef{Name: "std"} // inject stdlib dep

	return src.Module{
		Manifest: mod.Manifest,
		Packages: desugaredPkgs,
	}, nil
}

func (d Desugarer) desugarPkg(pkg src.Package, scope src.Scope) (src.Package, error) {
	desugaredPkgs := make(src.Package, len(pkg))

	for fileName, file := range pkg {
		newScope := scope.WithLocation(src.Location{
			ModRef:   scope.Location.ModRef,
			PkgName:  scope.Location.PkgName,
			FileName: fileName,
		})

		desugaredFile, err := d.desugarFile(file, newScope)
		if err != nil {
			return nil, err
		}

		desugaredPkgs[fileName] = desugaredFile
	}

	return desugaredPkgs, nil
}

func (d Desugarer) desugarFile(file src.File, scope src.Scope) (src.File, error) {
	desugaredEntities := make(map[string]src.Entity, len(file.Entities))

	for entityName, entity := range file.Entities {
		desugaredEntity, err := d.desugarEntity(entity, scope)
		if err != nil {
			return src.File{}, err
		}
		desugaredEntities[entityName] = desugaredEntity
	}

	return src.File{
		Imports:  file.Imports,
		Entities: desugaredEntities,
	}, nil
}

func (d Desugarer) desugarEntity(entity src.Entity, scope src.Scope) (src.Entity, error) {
	if entity.Kind != src.ComponentEntity {
		return entity, nil
	}

	desugarComponent, err := d.desugarComponent(entity.Component, scope)
	if err != nil {
		return src.Entity{}, err
	}

	return src.Entity{
		IsPublic:  entity.IsPublic,
		Kind:      entity.Kind,
		Component: desugarComponent,
	}, nil
}

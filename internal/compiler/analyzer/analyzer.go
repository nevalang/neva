// Package analyzer implements source code static semantic analysis.
// It's important to keep errors as human-readable as possible
// because they are what end-user is facing when something goes wrong.
package analyzer

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var (
	ErrModuleWithoutPkgs    = errors.New("Module must contain at least one package")
	ErrEntryModNotFound     = errors.New("Entry module is not found")
	ErrMainPkgNotFound      = errors.New("Main package not found")
	ErrPkgWithoutFiles      = errors.New("Package must contain at least one file")
	ErrUnknownEntityKind    = errors.New("Entity kind can only be either component, interface, type of constant")
	ErrCompilerVersion      = errors.New("Incompatible compiler version")
	ErrDepModWithoutVersion = errors.New("Every dependency module must have version")
)

type Analyzer struct {
	compilerVersion string
	resolver        ts.Resolver
}

func (a Analyzer) AnalyzeExecutableBuild(build src.Build, mainPkgName string) (src.Build, *compiler.Error) {
	location := src.Location{
		ModRef:  build.EntryModRef,
		PkgName: mainPkgName,
	}

	entryMod, ok := build.Modules[build.EntryModRef]
	if !ok {
		return src.Build{}, &compiler.Error{
			Err:      fmt.Errorf("%w: main package name '%s'", ErrEntryModNotFound, build.EntryModRef),
			Location: &location,
		}
	}

	if _, ok := entryMod.Packages[mainPkgName]; !ok {
		return src.Build{}, &compiler.Error{
			Err:      fmt.Errorf("%w: main package name '%s'", ErrMainPkgNotFound, mainPkgName),
			Location: &location,
		}
	}

	scope := src.Scope{
		Location: location,
		Build:    build,
	}

	if err := a.mainSpecificPkgValidation(mainPkgName, entryMod, scope); err != nil {
		return src.Build{}, compiler.Error{Location: &location}.Merge(err)
	}

	analyzedBuild, err := a.AnalyzeBuild(build)
	if err != nil {
		return src.Build{}, compiler.Error{Location: &location}.Merge(err)
	}

	return analyzedBuild, nil
}

func (a Analyzer) AnalyzeBuild(build src.Build) (src.Build, *compiler.Error) {
	analyzedMods := make(map[src.ModuleRef]src.Module, len(build.Modules))

	for modRef, mod := range build.Modules {
		if mod.Manifest.LanguageVersion != a.compilerVersion {
			return src.Build{}, &compiler.Error{
				Err: fmt.Errorf(
					"%w: module %v wants %v while current is %v",
					ErrCompilerVersion,
					modRef, mod.Manifest.LanguageVersion, a.compilerVersion,
				),
			}
		}

		analyzedPkgs, err := a.analyzeModule(modRef, build)
		if err != nil {
			return src.Build{}, err
		}

		analyzedMods[modRef] = src.Module{
			Manifest: mod.Manifest,
			Packages: analyzedPkgs,
		}
	}

	return src.Build{
		EntryModRef: build.EntryModRef,
		Modules:     analyzedMods,
	}, nil
}

func (a Analyzer) analyzeModule(modRef src.ModuleRef, build src.Build) (map[string]src.Package, *compiler.Error) {
	if modRef != build.EntryModRef && modRef.Version == "" {
		return nil, &compiler.Error{
			Err: ErrDepModWithoutVersion,
		}
	}

	location := src.Location{ModRef: modRef}
	mod := build.Modules[modRef]

	if len(mod.Packages) == 0 {
		return nil, &compiler.Error{
			Err:      ErrModuleWithoutPkgs,
			Location: &location,
		}
	}

	pkgsCopy := make(map[string]src.Package, len(mod.Packages))
	maps.Copy(pkgsCopy, mod.Packages)

	for pkgName, pkg := range pkgsCopy {
		scope := src.Scope{
			Location: src.Location{
				ModRef:  modRef,
				PkgName: pkgName,
			},
			Build: build,
		}

		resolvedPkg, err := a.analyzePkg(pkg, scope)
		if err != nil {
			return nil, compiler.Error{
				Location: &src.Location{
					PkgName: pkgName,
				},
			}.Merge(err)
		}

		pkgsCopy[pkgName] = resolvedPkg
	}

	return pkgsCopy, nil
}

func (a Analyzer) analyzePkg(pkg src.Package, scope src.Scope) (src.Package, *compiler.Error) {
	if len(pkg) == 0 {
		return nil, &compiler.Error{
			Err:      ErrPkgWithoutFiles,
			Location: &scope.Location,
		}
	}

	// preallocate
	analyzedFiles := make(map[string]src.File, len(pkg))
	for fileName, file := range pkg {
		analyzedFiles[fileName] = src.File{
			Imports:  file.Imports,
			Entities: make(map[string]src.Entity, len(file.Entities)),
		}
	}

	if err := pkg.Entities(func(entity src.Entity, entityName, fileName string) error {
		scopeWithFile := scope.WithLocation(src.Location{
			FileName: fileName,
			ModRef:   scope.Location.ModRef,
			PkgName:  scope.Location.PkgName,
		})

		resolvedEntity, err := a.analyzeEntity(entity, scopeWithFile)
		if err != nil {
			return compiler.Error{
				Location: &scopeWithFile.Location,
				Meta:     entity.Meta(),
			}.Merge(err)
		}

		analyzedFiles[fileName].Entities[entityName] = resolvedEntity

		return nil
	}); err != nil {
		return nil, err.(*compiler.Error) //nolint:forcetypeassert
	}

	return analyzedFiles, nil
}

func (a Analyzer) analyzeEntity(entity src.Entity, scope src.Scope) (src.Entity, *compiler.Error) {
	resolvedEntity := src.Entity{
		IsPublic: entity.IsPublic,
		Kind:     entity.Kind,
	}

	isStd := scope.Location.ModRef.Path == "std"

	switch entity.Kind {
	case src.TypeEntity:
		resolvedTypeDef, err := a.analyzeTypeDef(entity.Type, scope, analyzeTypeDefParams{allowEmptyBody: isStd})
		if err != nil {
			meta := entity.Type.Meta.(src.Meta) //nolint:forcetypeassert
			return src.Entity{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &meta,
			}.Merge(err)
		}
		resolvedEntity.Type = resolvedTypeDef
	case src.ConstEntity:
		resolvedConst, err := a.analyzeConst(entity.Const, scope)
		if err != nil {
			meta := entity.Const.Meta
			return src.Entity{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &meta,
			}.Merge(err)
		}
		resolvedEntity.Const = resolvedConst
	case src.InterfaceEntity:
		resolvedInterface, err := a.analyzeInterface(entity.Interface, scope, analyzeInterfaceParams{
			allowEmptyInports:  false,
			allowEmptyOutports: false,
		})
		if err != nil {
			meta := entity.Interface.Meta
			return src.Entity{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &meta,
			}.Merge(err)
		}
		resolvedEntity.Interface = resolvedInterface
	case src.ComponentEntity:
		analyzedComponent, err := a.analyzeComponent(entity.Component, scope)
		if err != nil {
			return src.Entity{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &entity.Component.Meta,
			}.Merge(err)
		}
		resolvedEntity.Component = analyzedComponent
	default:
		return src.Entity{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrUnknownEntityKind, entity.Kind),
			Location: &scope.Location,
			Meta:     entity.Meta(),
		}
	}

	return resolvedEntity, nil
}

func MustNew(version string, resolver ts.Resolver) Analyzer {
	return Analyzer{
		compilerVersion: version,
		resolver:        resolver,
	}
}

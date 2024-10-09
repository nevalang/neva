// Package analyzer implements source code static semantic analysis.
// It's important to keep errors as human-readable as possible
// because they are what end-user is facing when something goes wrong.
package analyzer

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

var (
	ErrModuleWithoutPkgs    = errors.New("module must contain at least one package")
	ErrEntryModNotFound     = errors.New("entry module not found")
	ErrPkgWithoutFiles      = errors.New("package must contain at least one file")
	ErrUnknownEntityKind    = errors.New("entity kind can only be either flow, interface, type of constant")
	ErrCompilerVersion      = errors.New("incompatible compiler version")
	ErrDepModWithoutVersion = errors.New("every dependency module must have version")
)

type Analyzer struct {
	resolver ts.Resolver
}

func (a Analyzer) AnalyzeExecutableBuild(build src.Build, mainPkgName string) (src.Build, *compiler.Error) {
	location := src.Location{
		ModRef:  build.EntryModRef,
		PkgName: mainPkgName,
	}

	entryMod, ok := build.Modules[build.EntryModRef]

	if !ok {
		return src.Build{}, &compiler.Error{
			Err:      fmt.Errorf("%w: entry module name '%s'", ErrEntryModNotFound, build.EntryModRef),
			Location: &location,
		}
	}

	// FIXME mainPkgName containts full path with "examples/ "
	if _, ok := entryMod.Packages[mainPkgName]; !ok {
		return src.Build{}, &compiler.Error{
			Err:      errors.New("main package not found"),
			Location: &location,
		}
	}

	scope := src.Scope{
		Location: location,
		Build:    build,
	}

	if err := a.mainSpecificPkgValidation(mainPkgName, entryMod, scope); err != nil {
		return src.Build{}, compiler.Error{Location: &location}.Wrap(err)
	}

	analyzedBuild, err := a.AnalyzeBuild(build)
	if err != nil {
		return src.Build{}, compiler.Error{Location: &location}.Wrap(err)
	}

	return analyzedBuild, nil
}

func (a Analyzer) AnalyzeBuild(build src.Build) (src.Build, *compiler.Error) {
	analyzedMods := make(map[src.ModuleRef]src.Module, len(build.Modules))

	for modRef, mod := range build.Modules {
		if err := a.semverCheck(mod, modRef); err != nil {
			return src.Build{}, err
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
			}.Wrap(err)
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
			}.Wrap(err)
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
			meta := entity.Type.Meta.(core.Meta) //nolint:forcetypeassert
			return src.Entity{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &meta,
			}.Wrap(err)
		}
		resolvedEntity.Type = resolvedTypeDef
	case src.ConstEntity:
		resolvedConst, err := a.analyzeConst(entity.Const, scope)
		if err != nil {
			meta := entity.Const.Meta
			return src.Entity{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &meta,
			}.Wrap(err)
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
			}.Wrap(err)
		}
		resolvedEntity.Interface = resolvedInterface
	case src.ComponentEntity:
		analyzedComponent, err := a.analyzeComponent(entity.Component, scope)
		if err != nil {
			return src.Entity{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &entity.Component.Meta,
			}.Wrap(err)
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

func MustNew(resolver ts.Resolver) Analyzer {
	return Analyzer{resolver: resolver}
}

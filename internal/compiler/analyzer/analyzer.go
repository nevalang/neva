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
	ErrCompilerVersion = errors.New("incompatible compiler version")
)

type Analyzer struct {
	resolver ts.Resolver
}

func (a Analyzer) AnalyzeExecutableBuild(build src.Build, mainPkgName string) (src.Build, *compiler.Error) {
	location := src.Location{
		Module:  build.EntryModRef,
		Package: mainPkgName,
	}

	entryMod, ok := build.Modules[build.EntryModRef]
	if !ok {
		return src.Build{}, &compiler.Error{
			Message:  fmt.Sprintf("entry module not found: %s", build.EntryModRef),
			Location: &location,
		}
	}

	if _, ok := entryMod.Packages[mainPkgName]; !ok {
		return src.Build{}, &compiler.Error{
			Message:  "main package not found",
			Location: &location,
		}
	}

	scope := src.NewScope(build, location)

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
			Message:  "every dependency module must have version",
			Location: &src.Location{Module: modRef},
		}
	}

	location := src.Location{Module: modRef}
	mod := build.Modules[modRef]

	if len(mod.Packages) == 0 {
		return nil, &compiler.Error{
			Message:  "module must contain at least one package",
			Location: &location,
		}
	}

	pkgsCopy := make(map[string]src.Package, len(mod.Packages))
	maps.Copy(pkgsCopy, mod.Packages)

	for pkgName, pkg := range pkgsCopy {
		scope := src.NewScope(build, src.Location{
			Module:  modRef,
			Package: pkgName,
		})

		resolvedPkg, err := a.analyzePkg(pkg, scope)
		if err != nil {
			return nil, compiler.Error{
				Location: &src.Location{
					Package: pkgName,
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
			Message:  "package must contain at least one file",
			Location: scope.Location(),
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

	for result := range pkg.Entities() {
		relocatedScope := scope.Relocate(src.Location{
			Module:   scope.Location().Module,
			Package:  scope.Location().Package,
			Filename: result.FileName,
		})

		analyzedEntity, err := a.analyzeEntity(result.Entity, relocatedScope)
		if err != nil {
			return nil, compiler.Error{
				Location: relocatedScope.Location(),
				Meta:     result.Entity.Meta(),
			}.Wrap(err)
		}

		analyzedFiles[result.FileName].Entities[result.EntityName] = analyzedEntity
	}

	return analyzedFiles, nil
}

func (a Analyzer) analyzeEntity(entity src.Entity, scope src.Scope) (src.Entity, *compiler.Error) {
	resolvedEntity := src.Entity{
		IsPublic: entity.IsPublic,
		Kind:     entity.Kind,
	}

	isStd := scope.Location().Module.Path == "std"

	switch entity.Kind {
	case src.TypeEntity:
		resolvedTypeDef, err := a.analyzeTypeDef(entity.Type, scope, analyzeTypeDefParams{allowEmptyBody: isStd})
		if err != nil {
			meta := entity.Type.Meta.(core.Meta) //nolint:forcetypeassert
			return src.Entity{}, compiler.Error{
				Location: scope.Location(),
				Meta:     &meta,
			}.Wrap(err)
		}
		resolvedEntity.Type = resolvedTypeDef
	case src.ConstEntity:
		resolvedConst, err := a.analyzeConst(entity.Const, scope)
		if err != nil {
			meta := entity.Const.Meta
			return src.Entity{}, compiler.Error{
				Location: scope.Location(),
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
				Location: scope.Location(),
				Meta:     &meta,
			}.Wrap(err)
		}
		resolvedEntity.Interface = resolvedInterface
	case src.ComponentEntity:
		analyzedComponent, err := a.analyzeComponent(entity.Component, scope)
		if err != nil {
			return src.Entity{}, compiler.Error{
				Location: scope.Location(),
				Meta:     &entity.Component.Meta,
			}.Wrap(err)
		}
		resolvedEntity.Component = analyzedComponent
	default:
		return src.Entity{}, &compiler.Error{
			Message:  fmt.Sprintf("unknown entity kind: %v", entity.Kind),
			Location: scope.Location(),
			Meta:     entity.Meta(),
		}
	}

	return resolvedEntity, nil
}

func MustNew(resolver ts.Resolver) Analyzer {
	return Analyzer{resolver: resolver}
}

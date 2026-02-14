// Package analyzer implements source code static semantic analysis.
// It's important to keep errors as human-readable as possible
// because they are what end-user is facing when something goes wrong.
package analyzer

import (
	"fmt"

	"golang.org/x/exp/maps"

	"github.com/nevalang/neva/internal/compiler"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
)

type Analyzer struct {
	resolver ts.Resolver
}

// AnalyzeBuild analyzes a build. When mainPkgName is non-empty,
// analyzer treats the build as an executable entry. When mainPkgName is empty,
// analyzer only analyzes the build as a library.
func (a Analyzer) Analyze(build src.Build, mainPkgName string) (src.Build, *compiler.Error) {
	if mainPkgName == "" {
		return a.analyzeBuild(build)
	}

	meta := core.Meta{
		Location: core.Location{
			ModRef:  build.EntryModRef,
			Package: mainPkgName,
		},
	}

	entryMod, ok := build.Modules[build.EntryModRef]
	if !ok {
		return src.Build{}, &compiler.Error{
			Message: fmt.Sprintf("entry module not found: %s", build.EntryModRef),
			Meta:    &meta,
		}
	}

	if _, ok := entryMod.Packages[mainPkgName]; !ok {
		return src.Build{}, &compiler.Error{
			Message: "main package not found",
			Meta:    &meta,
		}
	}

	scope := src.NewScope(build, meta.Location)

	if err := a.mainSpecificPkgValidation(mainPkgName, entryMod, scope); err != nil {
		return src.Build{}, compiler.Error{Meta: &meta}.Wrap(err)
	}

	analyzedBuild, err := a.analyzeBuild(build)
	if err != nil {
		return src.Build{}, compiler.Error{Meta: &meta}.Wrap(err)
	}

	return analyzedBuild, nil
}

func (a Analyzer) analyzeBuild(build src.Build) (src.Build, *compiler.Error) {
	analyzedMods := make(map[core.ModuleRef]src.Module, len(build.Modules))

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

func (a Analyzer) analyzeModule(modRef core.ModuleRef, build src.Build) (map[string]src.Package, *compiler.Error) {
	if modRef != build.EntryModRef && modRef.Version == "" {
		return nil, &compiler.Error{
			Message: "every dependency module must have version",
			Meta: &core.Meta{
				Location: core.Location{
					ModRef: modRef,
				},
			},
		}
	}

	location := core.Location{ModRef: modRef}
	mod := build.Modules[modRef]

	if len(mod.Packages) == 0 {
		return nil, &compiler.Error{
			Message: "module must contain at least one package",
			Meta: &core.Meta{
				Location: location,
			},
		}
	}

	pkgsCopy := make(map[string]src.Package, len(mod.Packages))
	maps.Copy(pkgsCopy, mod.Packages)

	for pkgName, pkg := range pkgsCopy {
		scope := src.NewScope(build, core.Location{
			ModRef:  modRef,
			Package: pkgName,
		})

		resolvedPkg, err := a.analyzePkg(pkg, scope)
		if err != nil {
			return nil, compiler.Error{
				Meta: &core.Meta{
					Location: core.Location{
						Package: pkgName,
					},
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
			Message: "package must contain at least one file",
			Meta: &core.Meta{
				Location: *scope.Location(),
			},
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
		relocatedScope := scope.Relocate(core.Location{
			ModRef:   scope.Location().ModRef,
			Package:  scope.Location().Package,
			Filename: result.FileName,
		})

		analyzedEntity, err := a.analyzeEntity(result.EntityName, result.Entity, relocatedScope)
		if err != nil {
			return nil, compiler.Error{
				Meta: result.Entity.Meta(),
			}.Wrap(err)
		}

		analyzedFiles[result.FileName].Entities[result.EntityName] = analyzedEntity
	}

	return analyzedFiles, nil
}

func (a Analyzer) analyzeEntity(entityName string, entity src.Entity, scope src.Scope) (src.Entity, *compiler.Error) {
	resolvedEntity := src.Entity{
		IsPublic: entity.IsPublic,
		Kind:     entity.Kind,
	}

	isStd := scope.Location().ModRef.Path == "std"

	switch entity.Kind {
	case src.TypeEntity:
		resolvedTypeDef, err := a.analyzeType(
			entity.Type,
			scope,
			analyzeTypeDefParams{allowEmptyBody: isStd},
		)
		if err != nil {
			meta := entity.Type.Meta
			return src.Entity{}, compiler.Error{
				Meta: &meta,
			}.Wrap(err)
		}
		resolvedEntity.Type = resolvedTypeDef
	case src.ConstEntity:
		resolvedConst, err := a.analyzeConst(entity.Const, scope)
		if err != nil {
			meta := entity.Const.Meta
			return src.Entity{}, compiler.Error{
				Meta: &meta,
			}.Wrap(err)
		}
		resolvedEntity.Const = resolvedConst
	case src.InterfaceEntity:
		resolvedInterface, err := a.analyzeInterface(
			entity.Interface,
			scope,
			analyzeInterfaceParams{
				allowEmptyInports:  false,
				allowEmptyOutports: false,
			})
		if err != nil {
			meta := entity.Interface.Meta
			return src.Entity{}, compiler.Error{
				Meta: &meta,
			}.Wrap(err)
		}
		resolvedEntity.Interface = resolvedInterface
	case src.ComponentEntity:
		analyzedVersions := make([]src.Component, 0, len(entity.Component))
		for _, component := range entity.Component {
			analyzedComponent, err := a.analyzeComponent(entityName, component, scope)
			if err != nil {
				return src.Entity{}, compiler.Error{
					Meta: &component.Meta,
				}.Wrap(err)
			}
			analyzedVersions = append(analyzedVersions, analyzedComponent)
		}
		resolvedEntity.Component = analyzedVersions
	default:
		return src.Entity{}, &compiler.Error{
			Message: fmt.Sprintf("unknown entity kind: %v", entity.Kind),
			Meta:    entity.Meta(),
		}
	}

	return resolvedEntity, nil
}

func MustNew(resolver ts.Resolver) Analyzer {
	return Analyzer{resolver: resolver}
}

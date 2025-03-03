package desugarer

import (
	"fmt"
	"maps"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	"github.com/nevalang/neva/pkg"
)

// Desugarer is NOT thread safe and must be used in single thread
type Desugarer struct {
	virtualSelectorsCount uint64
	ternaryCounter        uint64
	switchCounter         uint64
	virtualLocksCounter   uint64
	virtualEmittersCount  uint64
	virtualConstCount     uint64
	virtualTriggersCount  uint64
	fanOutCounter         uint64
	fanInCounter          uint64
	rangeCounter          uint64
	// Arithmetic
	addCounter uint64
	subCounter uint64
	mulCounter uint64
	divCounter uint64
	modCounter uint64
	powCounter uint64
	// Comparison
	eqCounter uint64
	neCounter uint64
	gtCounter uint64
	ltCounter uint64
	geCounter uint64
	leCounter uint64
	// Logical
	andCounter uint64
	orCounter  uint64
	// Bitwise
	bitAndCounter uint64
	bitOrCounter  uint64
	bitXorCounter uint64
	bitLshCounter uint64
	bitRshCounter uint64
}

func (d *Desugarer) Desugar(build src.Build) (src.Build, error) {
	desugaredMods := make(map[core.ModuleRef]src.Module, len(build.Modules))

	for modRef := range build.Modules {
		desugaredMod, err := d.desugarModule(build, modRef)
		if err != nil {
			return src.Build{}, fmt.Errorf("desugar module %s: %w", modRef, err)
		}
		desugaredMods[modRef] = desugaredMod
	}

	return src.Build{
		EntryModRef: build.EntryModRef,
		Modules:     desugaredMods,
	}, nil
}

func (d *Desugarer) desugarModule(
	build src.Build,
	modRef core.ModuleRef,
) (src.Module, error) {
	mod := build.Modules[modRef]

	// create manifest copy with std module dependency
	desugaredManifest := src.ModuleManifest{
		LanguageVersion: mod.Manifest.LanguageVersion,
		Deps:            make(map[string]core.ModuleRef, len(mod.Manifest.Deps)+1),
	}
	maps.Copy(desugaredManifest.Deps, mod.Manifest.Deps)
	desugaredManifest.Deps["std"] = core.ModuleRef{Path: "std", Version: pkg.Version}

	// copy all modules but replace manifest in current one
	modsCopy := maps.Clone(build.Modules)
	modsCopy[modRef] = src.Module{
		Manifest: desugaredManifest,
		Packages: mod.Packages,
	}

	// create new build with patched modules (current module has patched manifest with std dependency)
	build = src.Build{
		EntryModRef: modRef,
		Modules:     modsCopy,
	}

	desugaredPkgs := make(map[string]src.Package, len(mod.Packages))

	for pkgName, pkg := range mod.Packages {
		// it's important to patch build before desugar package so we can resolve references to std
		scope := src.NewScope(build, core.Location{
			ModRef:  modRef,
			Package: pkgName,
		})

		desugaredPkg, err := d.desugarPkg(pkg, scope)
		if err != nil {
			return src.Module{}, fmt.Errorf("desugar package %s: %w", pkgName, err)
		}

		desugaredPkgs[pkgName] = desugaredPkg
	}

	return src.Module{
		Manifest: desugaredManifest,
		Packages: desugaredPkgs,
	}, nil
}

// Scope interface allows to use mocks in unit tests for private methods
//
//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}
type Scope interface {
	Entity(ref core.EntityRef) (src.Entity, core.Location, error)
	Relocate(location core.Location) src.Scope
	Location() *core.Location
	GetNodeIOByPortAddr(nodes map[string]src.Node, portAddr src.PortAddr) (src.IO, error)
}

func (d *Desugarer) desugarPkg(pkg src.Package, scope Scope) (src.Package, error) {
	desugaredPkgs := make(src.Package, len(pkg))

	for fileName, file := range pkg {
		newScope := scope.Relocate(core.Location{
			ModRef:   scope.Location().ModRef,
			Package:  scope.Location().Package,
			Filename: fileName,
		})

		desugaredFile, err := d.desugarFile(file, newScope)
		if err != nil {
			return src.Package{}, fmt.Errorf("desugar file %s: %w", fileName, err)
		}

		desugaredPkgs[fileName] = desugaredFile
	}

	return desugaredPkgs, nil
}

// desugarFile injects import of std/builtin into every pkg's file and desugares it's every entity
func (d *Desugarer) desugarFile(
	file src.File,
	scope src.Scope,
) (src.File, error) {
	desugaredEntities := make(map[string]src.Entity, len(file.Entities))

	for entityName, entity := range file.Entities {
		entityResult, err := d.desugarEntity(entity, scope)
		if err != nil {
			return src.File{}, fmt.Errorf("desugar entity %s: %w", entityName, err)
		}

		// FIXMEL: https://github.com/nevalang/neva/issues/808
		// d.resetCounters()

		desugaredEntities[entityName] = entityResult.entity

		// insert virtual entities, created by desugaring of the entity
		for name, entityToInsert := range entityResult.insert {
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
		Meta:    core.Meta{Location: *scope.Location()},
	}

	return src.File{
		Imports:  desugaredImports,
		Entities: desugaredEntities,
	}, nil
}

type desugarEntityResult struct {
	entity src.Entity
	insert map[string]src.Entity
}

func (d *Desugarer) desugarEntity(
	entity src.Entity,
	scope src.Scope,
) (desugarEntityResult, error) {
	if entity.Kind != src.ComponentEntity && entity.Kind != src.ConstEntity {
		return desugarEntityResult{entity: entity}, nil
	}

	if entity.Kind == src.ConstEntity {
		desugaredConst, err := d.handleConst(entity.Const)
		if err != nil {
			return desugarEntityResult{}, fmt.Errorf("desugar const: %w", err)
		}

		return desugarEntityResult{
			entity: src.Entity{
				IsPublic: entity.IsPublic,
				Kind:     entity.Kind,
				Const:    desugaredConst,
			},
		}, nil
	}

	componentResult, err := d.desugarComponent(entity.Component, scope)
	if err != nil {
		return desugarEntityResult{}, fmt.Errorf("desugar component: %w", err)
	}

	return desugarEntityResult{
		insert: componentResult.virtualEntities,
		entity: src.Entity{
			IsPublic:  entity.IsPublic,
			Kind:      entity.Kind,
			Component: componentResult.desugaredFlow,
		},
	}, nil
}

// FIXME: https://github.com/nevalang/neva/issues/808
// Do NOT use this method until issue is fixed
// func (d *Desugarer) resetCounters() {
// 	d.virtualSelectorsCount = 0
// 	d.ternaryCounter = 0
// 	d.switchCounter = 0
// 	d.virtualLocksCounter = 0
// 	d.virtualEmittersCount = 0
// 	d.virtualConstCount = 0
// 	d.virtualTriggersCount = 0
// 	d.fanOutCounter = 0
// 	d.fanInCounter = 0
// 	d.rangeCounter = 0
// 	//
// 	d.addCounter = 0
// 	d.subCounter = 0
// 	d.mulCounter = 0
// 	d.divCounter = 0
// 	d.modCounter = 0
// 	d.powCounter = 0
// 	//
// 	d.eqCounter = 0
// 	d.neCounter = 0
// 	d.gtCounter = 0
// 	d.ltCounter = 0
// 	d.geCounter = 0
// 	d.leCounter = 0
// 	//
// 	d.andCounter = 0
// 	d.orCounter = 0
// 	d.bitAndCounter = 0
// 	d.bitOrCounter = 0
// 	d.bitXorCounter = 0
// 	d.bitLshCounter = 0
// 	d.bitRshCounter = 0
// 	d.mulCounter = 0
// 	d.divCounter = 0
// 	d.modCounter = 0
// 	d.powCounter = 0
// 	//
// 	d.eqCounter = 0
// 	d.neCounter = 0
// 	d.gtCounter = 0
// 	d.ltCounter = 0
// 	d.geCounter = 0
// 	d.leCounter = 0
// 	//
// 	d.andCounter = 0
// 	d.orCounter = 0
// 	//
// 	d.bitAndCounter = 0
// 	d.bitOrCounter = 0
// 	d.bitXorCounter = 0
// 	d.bitLshCounter = 0
// 	d.bitRshCounter = 0
// }

func New() Desugarer {
	return Desugarer{}
}

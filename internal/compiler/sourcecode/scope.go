package sourcecode

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	"github.com/nevalang/neva/pkg"
)

// NewScope returns a new scope with a given location
func NewScope(build Build, location Location) Scope {
	return Scope{
		loc:   location,
		build: build,
	}
}

// Scope is an object that provides access to program entities
type Scope struct {
	loc   Location
	build Build
}

// Location returns a location of the current scope
func (s Scope) Location() *Location {
	return &s.loc
}

// Relocate returns a new scope with a given location
func (s Scope) Relocate(location Location) Scope {
	return Scope{
		loc:   location,
		build: s.build,
	}
}

// Location is a source code location
type Location struct {
	Module   ModuleRef
	Package  string
	Filename string
}

// String returns location as a string
func (l Location) String() string {
	var s string
	if l.Module.Path == "@" {
		s = l.Package
	} else {
		s = filepath.Join(l.Module.String(), l.Package)
	}
	if l.Filename != "" {
		s = filepath.Join(s, l.Filename+".neva")
	}
	return s
}

// IsTopType returns true if expr is a top type (any)
func (s Scope) IsTopType(expr ts.Expr) bool {
	if expr.Inst == nil {
		return false
	}
	if expr.Inst.Ref.Name != "any" {
		return false
	}
	return expr.Inst.Ref.Pkg == "" || expr.Inst.Ref.Pkg == "builtin"
}

// GetType returns type definition by reference
func (s Scope) GetType(ref core.EntityRef) (ts.Def, ts.Scope, error) {
	entity, location, err := s.Entity(ref)
	if err != nil {
		return ts.Def{}, nil, err
	}

	return entity.Type, s.Relocate(location), nil
}

// Entity returns entity by reference
func (s Scope) Entity(entityRef core.EntityRef) (Entity, Location, error) {
	return s.entity(entityRef)
}

// entity is an alrogithm that resolves entity reference based on scope's location
func (s Scope) entity(entityRef core.EntityRef) (Entity, Location, error) {
	curMod, ok := s.build.Modules[s.loc.Module]
	if !ok {
		return Entity{}, Location{}, fmt.Errorf("module not found: %v", s.loc.Module)
	}

	curPkg := curMod.Packages[s.loc.Package]
	if !ok {
		return Entity{}, Location{}, fmt.Errorf("package not found: %v", s.loc.Package)
	}

	if entityRef.Pkg == "" { // local reference (current package or builtin)
		entity, fileName, ok := curPkg.Entity(entityRef.Name)
		if ok {
			return entity, Location{
				Module:   s.loc.Module,
				Package:  s.loc.Package,
				Filename: fileName,
			}, nil
		}

		stdModRef := ModuleRef{Path: "std", Version: pkg.Version}
		stdMod, ok := s.build.Modules[stdModRef]
		if !ok {
			return Entity{}, Location{}, fmt.Errorf("std module not found: %v", stdModRef)
		}

		entity, fileName, err := stdMod.Entity(core.EntityRef{
			Pkg:  "builtin",
			Name: entityRef.Name,
		})
		if err != nil {
			return Entity{}, Location{}, err
		}

		return entity, Location{
			Module:   stdModRef,
			Package:  "builtin",
			Filename: fileName,
		}, nil
	}

	curFile, ok := curPkg[s.loc.Filename]
	if !ok {
		return Entity{}, Location{}, fmt.Errorf("file not found: %v", s.loc.Filename)
	}

	pkgImport, ok := curFile.Imports[entityRef.Pkg]
	if !ok {
		return Entity{}, Location{}, fmt.Errorf("import not found: %v", entityRef.Pkg)
	}

	var (
		mod    Module
		modRef ModuleRef
	)
	if pkgImport.Module == "@" {
		modRef = s.loc.Module // FIXME s.Location.ModRef is where we are now (e.g. std)
		mod = curMod
	} else {
		modRef = curMod.Manifest.Deps[pkgImport.Module]
		depMod, ok := s.build.Modules[modRef]
		if !ok {
			return Entity{}, Location{}, fmt.Errorf("dependency module not found: %v", modRef)
		}
		mod = depMod
	}

	ref := core.EntityRef{
		Pkg:  pkgImport.Package,
		Name: entityRef.Name,
	}

	entity, fileName, err := mod.Entity(ref)
	if err != nil {
		return Entity{}, Location{}, err
	}

	if !entity.IsPublic {
		return Entity{}, Location{}, errors.New("entity is not public")
	}

	return entity, Location{
		Module:   modRef,
		Package:  pkgImport.Package,
		Filename: fileName,
	}, nil
}

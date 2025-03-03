package sourcecode

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	"github.com/nevalang/neva/pkg"
)

// NewScope returns a new scope with a given location
func NewScope(build Build, location core.Location) Scope {
	return Scope{
		loc:   location,
		build: build,
	}
}

// Scope is entity reference resolver
type Scope struct {
	loc   core.Location
	build Build
}

// Location returns a location of the current scope
func (s Scope) Location() *core.Location {
	return &s.loc
}

// Relocate returns a new scope with a given location
func (s Scope) Relocate(location core.Location) Scope {
	return Scope{
		loc:   location,
		build: s.build,
	}
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
	entity, location, err := s.entity(ref)
	if err != nil {
		return ts.Def{}, nil, err
	}
	return entity.Type, s.Relocate(location), nil
}

func (s Scope) GetInterface(ref core.EntityRef) (Interface, error) {
	entity, _, err := s.entity(ref)
	if err != nil {
		return Interface{}, err
	}
	return entity.Interface, nil
}

// Entity returns entity by reference
func (s Scope) Entity(entityRef core.EntityRef) (Entity, core.Location, error) {
	return s.entity(entityRef)
}

func (s Scope) GetComponent(entityRef core.EntityRef) (Component, error) {
	entity, _, err := s.entity(entityRef)
	if err != nil {
		return Component{}, err
	}
	return entity.Component, nil
}

// entity is an alrogithm that resolves entity reference based on scope's location
func (s Scope) entity(entityRef core.EntityRef) (Entity, core.Location, error) {
	curMod, ok := s.build.Modules[s.loc.ModRef]
	if !ok {
		return Entity{}, core.Location{}, fmt.Errorf("module not found: %v", s.loc.ModRef)
	}

	curPkg := curMod.Packages[s.loc.Package]
	if !ok {
		return Entity{}, core.Location{}, fmt.Errorf("package not found: %v", s.loc.Package)
	}

	if entityRef.Pkg == "" { // local reference (current package or builtin)
		entity, fileName, ok := curPkg.Entity(entityRef.Name)
		if ok {
			return entity, core.Location{
				ModRef:   s.loc.ModRef,
				Package:  s.loc.Package,
				Filename: fileName,
			}, nil
		}

		stdModRef := core.ModuleRef{Path: "std", Version: pkg.Version}
		stdMod, ok := s.build.Modules[stdModRef]
		if !ok {
			return Entity{}, core.Location{}, fmt.Errorf("std module not found: %v", stdModRef)
		}

		entity, fileName, err := stdMod.Entity(core.EntityRef{
			Pkg:  "builtin",
			Name: entityRef.Name,
		})
		if err != nil {
			return Entity{}, core.Location{}, err
		}

		return entity, core.Location{
			ModRef:   stdModRef,
			Package:  "builtin",
			Filename: fileName,
		}, nil
	}

	curFile, ok := curPkg[s.loc.Filename]
	if !ok {
		return Entity{}, core.Location{}, fmt.Errorf("file not found: %v", s.loc.Filename)
	}

	pkgImport, ok := curFile.Imports[entityRef.Pkg]
	if !ok {
		return Entity{}, core.Location{}, fmt.Errorf("import not found: %v", entityRef.Pkg)
	}

	var (
		mod    Module
		modRef core.ModuleRef
	)
	if pkgImport.Module == "@" {
		modRef = s.loc.ModRef // FIXME s.Location.ModRef is where we are now (e.g. std)
		mod = curMod
	} else {
		modRef = curMod.Manifest.Deps[pkgImport.Module]
		depMod, ok := s.build.Modules[modRef]
		if !ok {
			return Entity{}, core.Location{}, fmt.Errorf("dependency module not found: %v", modRef)
		}
		mod = depMod
	}

	ref := core.EntityRef{
		Pkg:  pkgImport.Package,
		Name: entityRef.Name,
	}

	entity, fileName, err := mod.Entity(ref)
	if err != nil {
		return Entity{}, core.Location{}, err
	}

	if !entity.IsPublic {
		return Entity{}, core.Location{}, errors.New("entity is not public")
	}

	return entity, core.Location{
		ModRef:   modRef,
		Package:  pkgImport.Package,
		Filename: fileName,
	}, nil
}

func (s Scope) GetNodeIOByPortAddr(
	nodes map[string]Node,
	portAddr PortAddr,
) (IO, error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return IO{}, fmt.Errorf("node '%s' not found", portAddr.Node)
	}

	entity, _, err := s.Entity(node.EntityRef)
	if err != nil {
		return IO{}, fmt.Errorf("get entity: %w", err)
	}

	var iface Interface
	if entity.Kind == InterfaceEntity {
		iface = entity.Interface
	} else {
		iface = entity.Component.Interface
	}

	return iface.IO, nil
}

func (s Scope) GetEntityKind(entityRef core.EntityRef) (EntityKind, error) {
	entity, _, err := s.entity(entityRef)
	if err != nil {
		return "", err
	}
	return entity.Kind, nil
}

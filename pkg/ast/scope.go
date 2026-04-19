package ast

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/pkg"
	"github.com/nevalang/neva/pkg/core"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

// NewScope returns a new scope with a given location
func NewScope(build Build, location core.Location) Scope {
	return Scope{
		loc:   location,
		build: build,
	}
}

// Scope is entity reference resolver
//

type Scope struct {
	build Build
	loc   core.Location
}

// Location returns a location of the current scope
//
//nolint:gocritic
func (s Scope) Location() *core.Location {
	return &s.loc
}

// Relocate returns a new scope with a given location
//
//nolint:gocritic
func (s Scope) Relocate(location core.Location) Scope {
	return Scope{
		loc:   location,
		build: s.build,
	}
}

// IsTopType returns true if expr is a top type (any)
//
//nolint:gocritic
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
//
//nolint:gocritic,ireturn
func (s Scope) GetType(ref core.EntityRef) (ts.Def, ts.Scope, error) {
	entity, location, err := s.entity(ref)
	if err != nil {
		return ts.Def{}, nil, err
	}
	return entity.Type, s.Relocate(location), nil
}

// Entity returns entity by reference
//
//nolint:gocritic
func (s Scope) Entity(entityRef core.EntityRef) (Entity, core.Location, error) {
	return s.entity(entityRef)
}

//nolint:gocritic
func (s Scope) GetConst(entityRef core.EntityRef) (Const, core.Location, error) {
	entity, loc, err := s.entity(entityRef)
	if err != nil {
		return Const{}, core.Location{}, err
	}

	if entity.Kind != ConstEntity {
		return Const{}, core.Location{}, fmt.Errorf("entity is not a constant: %v", entity.Kind)
	}

	return entity.Const, loc, nil
}

// TODO rename to GetComponents
//
//nolint:godoclint
//nolint:gocritic
func (s Scope) GetComponent(entityRef core.EntityRef) ([]Component, error) { //nolint:gocritic,lll
	entity, _, err := s.entity(entityRef)
	if err != nil {
		return nil, err
	}

	if entity.Kind != ComponentEntity {
		return nil, fmt.Errorf("entity is not a component: %v", entity.Kind)
	}

	return entity.Component, nil
}

// entity is an alrogithm that resolves entity reference based on scope's location
//
//nolint:cyclop,funlen,gocognit,gocritic,gocyclo
func (s Scope) entity(entityRef core.EntityRef) (Entity, core.Location, error) {
	//nolint:varnamelen
	curMod, ok := s.build.Modules[s.loc.ModRef]
	if !ok {
		return Entity{}, core.Location{}, fmt.Errorf("module not found: %v", s.loc.ModRef)
	}

	curPkg := curMod.Packages[s.loc.Package]
	if !ok {
		return Entity{}, core.Location{}, fmt.Errorf("package not found: %v", s.loc.Package)
	}

	if entityRef.Pkg == "" { // local reference (current package or builtin)
		//nolint:varnamelen
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
		modRef, ok = curMod.Manifest.Deps[pkgImport.Module]
		if !ok {
			panic(fmt.Errorf("dependency module not found: %v (module key: %v, available deps: %v)", modRef, pkgImport.Module, curMod.Manifest.Deps))
		}
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

//nolint:gocritic
func (s Scope) GetNodeIOByPortAddr(
	nodes map[string]Node,
	//nolint:gocritic
	portAddr PortAddr,
) (IO, error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return IO{}, fmt.Errorf("node '%s' not found", portAddr.Node)
	}

	entity, _, err := s.entity(node.EntityRef)
	if err != nil {
		return IO{}, fmt.Errorf("get entity: %w", err)
	}

	if entity.Kind == InterfaceEntity {
		return entity.Interface.IO, nil
	}

	if len(entity.Component) == 1 {
		return entity.Component[0].IO, nil
	} else if len(entity.Component) > 1 {
		return entity.Component[*node.OverloadIndex].IO, nil
	}

	return IO{}, errors.New("component not found")
}

//nolint:gocritic
func (s Scope) GetFirstInportName(
	nodes map[string]Node,
	//nolint:gocritic
	portAddr PortAddr,
) (string, error) {
	io, err := s.GetNodeIOByPortAddr(nodes, portAddr)
	if err != nil {
		return "", err
	}

	for inport := range io.In {
		return inport, nil
	}

	return "", errors.New("first inport not found")
}

//nolint:gocritic
func (s Scope) GetEntityKind(entityRef core.EntityRef) (EntityKind, error) {
	entity, _, err := s.entity(entityRef)
	if err != nil {
		return "", err
	}

	return entity.Kind, nil
}

//nolint:gocritic
func (s Scope) GetFirstOutportName(
	nodes map[string]Node,
	//nolint:gocritic
	portAddr PortAddr,
) (string, error) {
	io, err := s.GetNodeIOByPortAddr(nodes, portAddr)
	if err != nil {
		return "", err
	}

	for outport := range io.Out {
		return outport, nil
	}

	return "", errors.New("first outport not found")
}

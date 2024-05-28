package sourcecode

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	"github.com/nevalang/neva/pkg"
)

var (
	ErrModNotFound    = errors.New("module not found")
	ErrPkgNotFound    = fmt.Errorf("package not found")
	ErrImportNotFound = errors.New("import not found")
	ErrFileNotFound   = errors.New("file not found")
	ErrEntityNotPub   = errors.New("entity is not public")
)

type Scope struct {
	Location Location
	Build    Build
}

func (s Scope) WithLocation(location Location) Scope {
	return Scope{
		Location: location,
		Build:    s.Build,
	}
}

type Location struct {
	ModRef   ModuleRef
	PkgName  string
	FileName string
}

func (l Location) String() string {
	var s string
	if l.ModRef.Path == "@" {
		s = l.PkgName
	} else {
		s = filepath.Join(l.ModRef.String(), l.PkgName)
	}
	if l.FileName != "" {
		s = filepath.Join(s, l.FileName+".neva")
	}
	return s
}

func (s Scope) IsTopType(expr ts.Expr) bool {
	if expr.Inst == nil {
		return false
	}
	if expr.Inst.Ref.Name != "any" {
		return false
	}
	return expr.Inst.Ref.Pkg == "" || expr.Inst.Ref.Pkg == "builtin"
}

func (s Scope) GetType(ref core.EntityRef) (ts.Def, ts.Scope, error) {
	entity, location, err := s.Entity(ref)
	if err != nil {
		return ts.Def{}, nil, err
	}

	return entity.Type, s.WithLocation(location), nil
}

func (s Scope) Entity(entityRef core.EntityRef) (Entity, Location, error) {
	return s.entity(entityRef)
}

func (s Scope) entity(entityRef core.EntityRef) (Entity, Location, error) {
	curMod, ok := s.Build.Modules[s.Location.ModRef]
	if !ok {
		return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrModNotFound, s.Location.ModRef)
	}

	curPkg := curMod.Packages[s.Location.PkgName]
	if !ok {
		return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrPkgNotFound, s.Location.PkgName)
	}

	if entityRef.Pkg == "" { // local reference (current package or builtin)
		entity, fileName, ok := curPkg.Entity(entityRef.Name)
		if ok {
			return entity, Location{
				ModRef:   s.Location.ModRef,
				PkgName:  s.Location.PkgName,
				FileName: fileName,
			}, nil
		}

		stdModRef := ModuleRef{Path: "std", Version: pkg.Version}
		stdMod, ok := s.Build.Modules[stdModRef]
		if !ok {
			return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrModNotFound, stdModRef)
		}

		entity, fileName, err := stdMod.Entity(core.EntityRef{
			Pkg:  "builtin",
			Name: entityRef.Name,
		})
		if err != nil {
			return Entity{}, Location{}, err
		}

		return entity, Location{
			ModRef:   stdModRef,
			PkgName:  "builtin",
			FileName: fileName,
		}, nil
	}

	curFile, ok := curPkg[s.Location.FileName]
	if !ok {
		return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrFileNotFound, s.Location.FileName)
	}

	pkgImport, ok := curFile.Imports[entityRef.Pkg]
	if !ok {
		return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrImportNotFound, entityRef.Pkg)
	}

	var (
		mod    Module
		modRef ModuleRef
	)
	if pkgImport.Module == "@" {
		modRef = s.Location.ModRef // FIXME s.Location.ModRef is where we are now (e.g. std)
		mod = curMod
	} else {
		modRef = curMod.Manifest.Deps[pkgImport.Module]
		depMod, ok := s.Build.Modules[modRef]
		if !ok {
			return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrModNotFound, modRef)
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
		return Entity{}, Location{}, ErrEntityNotPub
	}

	return entity, Location{
		ModRef:   modRef,
		PkgName:  pkgImport.Package,
		FileName: fileName,
	}, nil
}

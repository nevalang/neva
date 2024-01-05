package sourcecode

import (
	"errors"
	"fmt"

	ts "github.com/nevalang/neva/pkg/typesystem"
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
	if l.ModRef.Path == "" { // cur mod
		return fmt.Sprintf("%v/%v.neva", l.PkgName, l.FileName)
	}
	return fmt.Sprintf("%v/%v/%v.neva", l.ModRef, l.PkgName, l.FileName)
}

func (s Scope) IsTopType(expr ts.Expr) bool {
	if expr.Inst == nil {
		return false
	}
	parsed, ok := expr.Inst.Ref.(EntityRef)
	if !ok {
		return false
	}
	return parsed.Name == "any" && (parsed.Pkg == "" || parsed.Pkg == "builtin")
}

func (s Scope) GetType(ref fmt.Stringer) (ts.Def, ts.Scope, error) {
	parsedRef, ok := ref.(EntityRef)
	if !ok {
		return ts.Def{}, Scope{}, fmt.Errorf("ref is not entity ref: %v", ref)
	}

	entity, location, err := s.Entity(parsedRef)
	if err != nil {
		return ts.Def{}, nil, err
	}

	return entity.Type, s.WithLocation(location), nil
}

func (s Scope) Entity(entityRef EntityRef) (Entity, Location, error) {
	return s.entity(entityRef)
}

//nolint:funlen
func (s Scope) entity(entityRef EntityRef) (Entity, Location, error) {
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

		stdModRef := ModuleRef{Path: "std", Version: "0.0.1"}
		stdMod, ok := s.Build.Modules[stdModRef]
		if !ok {
			return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrModNotFound, stdModRef)
		}

		entity, fileName, err := stdMod.Entity(EntityRef{
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
	if pkgImport.ModuleName == "@" {
		modRef = s.Location.ModRef
		mod = curMod
	} else {
		modRef = curMod.Manifest.Deps[pkgImport.ModuleName]
		depMod, ok := s.Build.Modules[modRef]
		if !ok {
			return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrModNotFound, modRef)
		}
		mod = depMod
	}

	ref := EntityRef{
		Pkg:  pkgImport.PkgName,
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
		PkgName:  pkgImport.PkgName,
		FileName: fileName,
	}, nil
}

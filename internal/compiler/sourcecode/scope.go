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

//nolint:funlen
func (s Scope) Entity(entityRef EntityRef) (Entity, Location, error) {
	curMod, ok := s.Build.Modules[s.Location.ModRef]
	if !ok {
		return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrModNotFound, s.Location.ModRef)
	}

	curPkg := curMod.Packages[s.Location.PkgName]
	if !ok {
		return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrPkgNotFound, s.Location.PkgName)
	}

	// local reference (current package or builtin)
	if entityRef.Pkg == "" {
		entity, fileName, ok := curPkg.Entity(entityRef.Name)
		if ok {
			return entity, Location{
				PkgName:  s.Location.PkgName,
				FileName: fileName,
			}, nil
		}

		stdModRef := ModuleRef{Name: "std"}
		stdMod, ok := s.Build.Modules[stdModRef]
		if !ok {
			return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrModNotFound, stdModRef)
		}

		builtinPkgName := "builtin"
		builtinPkg, ok := stdMod.Packages[builtinPkgName]
		if !ok {
			return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrPkgNotFound, "builtin")
		}

		entity, fileName, ok = builtinPkg.Entity(entityRef.Name)
		if !ok {
			return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrEntityNotFound, entityRef.Name)
		}

		return entity, Location{
			ModRef:   stdModRef,
			PkgName:  builtinPkgName,
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

	var mod Module
	if pkgImport.ModuleName == "@" {
		mod = curMod
	} else {
		depModRef := curMod.Manifest.Deps[pkgImport.ModuleName]
		depMod, ok := s.Build.Modules[depModRef]
		if !ok {
			return Entity{}, Location{}, fmt.Errorf("%w: %v", ErrModNotFound, depModRef)
		}
		mod = depMod
	}

	entity, fileName, err := mod.Entity(EntityRef{
		Pkg:  pkgImport.PkgName,
		Name: entityRef.Name,
	})
	if err != nil {
		return Entity{}, Location{}, err
	}

	if !entity.IsPublic {
		return Entity{}, Location{}, ErrEntityNotPub
	}

	return entity, Location{
		PkgName:  pkgImport.PkgName,
		FileName: fileName,
	}, nil
}

// Entity does not return package because calleer knows it, passed entityRef contains it.
// Note that this method does not know anything about imports, builtins or anything like that.
// entityRef passed must be absolute (full, "real") path to the entity.
func (mod Module) Entity(entityRef EntityRef) (entity Entity, filename string, err error) {
	pkg, ok := mod.Packages[entityRef.Pkg]
	if !ok {
		return Entity{}, "", fmt.Errorf("%w: %s", ErrPkgNotFound, entityRef.Pkg)
	}
	for filename, file := range pkg {
		entity, ok := file.Entities[entityRef.Name]
		if ok {
			return entity, filename, nil
		}
	}
	return Entity{}, "", ErrEntityNotFound
}

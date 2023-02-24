package analyze

import (
	"errors"
	"strings"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

// Scope is created ad-hoc
type Scope struct {
	imports         map[string]src.Pkg
	local, builtins map[string]ts.Def
	visited         map[src.EntityRef]struct{}
}

// Get returns type definitions by reference. Reference must be src.EntityRef
func (s Scope) Get(ref string) (ts.Def, error) {
	var entityRef src.EntityRef

	parts := strings.Split(ref, ".")
	if len(parts) == 2 {
		entityRef.Pkg = parts[0]
		entityRef.Name = parts[1]
	} else {
		entityRef.Name = ref
	}

	if entityRef.Pkg == "" {
		localDef, ok := s.local[entityRef.Name]
		if ok {
			return localDef, nil
		}

		builtinDef, ok := s.builtins[entityRef.Name]
		if !ok {
			return ts.Def{}, errors.New("ref without package not found in local and builtin")
		}

		return builtinDef, nil
	}

	pkg, ok := s.imports[entityRef.Pkg]
	if !ok {
		return ts.Def{}, errors.New("referenced package is not found among imports")
	}

	entity, ok := pkg.Entities[entityRef.Name]
	if !ok {
		return ts.Def{}, errors.New("referenced entity not found in imported package")
	}

	if entity.Kind != src.TypeEntity {
		return ts.Def{}, errors.New("referenced entity is not type")
	}

	if !entity.Exported {
		return ts.Def{}, errors.New("referenced entity not exported")
	}

	s.visited[entityRef] = struct{}{}

	return entity.Type, nil
}

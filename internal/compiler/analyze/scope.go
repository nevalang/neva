package analyze

import (
	"errors"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

type Scope struct {
	pkgs map[string]src.Pkg
}

// Get returns type definitions by reference. Reference must be src.EntityRef
func (s Scope) Get(ref ts.Ref) (ts.Def, error) {
	entityRef, ok := ref.(src.EntityRef)
	if !ok {
		return ts.Def{}, errors.New("not entity ref")
	}

	pkg, ok := s.pkgs[entityRef.Pkg]
	if !ok {
		return ts.Def{}, errors.New("no package")
	}

	entity, ok := pkg.Entities[entityRef.Name]
	if !ok {
		return ts.Def{}, errors.New("no entity")
	}

	if entity.Kind != src.TypeEntity {
		return ts.Def{}, errors.New("not type")
	}

	if !entity.Exported {
		return ts.Def{}, errors.New("not exported")
	}

	return entity.Type, nil
}


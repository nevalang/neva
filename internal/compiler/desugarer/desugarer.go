package desugarer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
)

// TODO currently desugarer does nothing. We'll keep this template-code for the future.
// If we won't need it, we will remove it.
type Desugarer struct{}

func (d Desugarer) Desugar(prog src.Program) (src.Program, error) {
	result := make(src.Program, len(prog))

	for pkgName := range prog {
		desugaredPkg, err := d.desugarPkg(pkgName, prog)
		if err != nil {
			return src.Program{}, nil
		}
		result[pkgName] = desugaredPkg
	}

	return result, nil
}

func (d Desugarer) desugarPkg(pkgName string, prog src.Program) (src.Package, error) {
	pkg := prog[pkgName]
	result := make(src.Package, len(pkg))

	for fileName, file := range pkg {
		result[fileName] = src.File{
			Imports:  file.Imports,
			Entities: make(map[string]src.Entity, len(file.Entities)),
		}

		for entityName, entity := range file.Entities {
			scope := src.Scope{
				Loc: src.ScopeLocation{
					PkgName:  pkgName,
					FileName: fileName,
				},
				Prog: prog,
			}

			desugaredEntity, err := d.desugarEntity(entity, scope)
			if err != nil {
				return src.Package{}, fmt.Errorf("desugar entity: %w", err)
			}

			result[fileName].Entities[entityName] = desugaredEntity
		}
	}

	return result, nil
}

func (d Desugarer) desugarEntity(entity src.Entity, scope src.Scope) (src.Entity, error) {
	return entity, nil
}

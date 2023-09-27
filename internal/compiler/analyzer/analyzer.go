// Package analyzer implements source code analysis.
package analyzer

import "github.com/nevalang/neva/internal/compiler/src"

type Analyzer struct{}

func (a Analyzer) Analyze(prog src.Program) (src.Program, error) {
	if prog == nil {
		panic("analyzer: nil program")
	}

	mainPkg, ok := prog["main"]
	if !ok {
		panic("analyzer: no main package")
	}

	if err := a.mainSpecificPkgValidation(mainPkg, prog); err != nil {
		panic(err)
	}

	return nil, nil
}

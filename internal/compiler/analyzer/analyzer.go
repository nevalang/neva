// Package analyzer implements source code analysis.
package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
)

type Analyzer struct{}

var (
	ErrEmptyProgram    = errors.New("empty program")
	ErrMainPkgNotFound = errors.New("main package not found")
)

func (a Analyzer) Analyze(prog src.Program) (src.Program, error) {
	if prog == nil {
		return nil, ErrEmptyProgram
	}

	mainPkg, ok := prog["main"]
	if !ok {
		return nil, ErrMainPkgNotFound
	}

	if err := a.mainSpecificPkgValidation(mainPkg, prog); err != nil {
		return nil, fmt.Errorf("main specific pkg validation: %w", err)
	}

	return nil, nil
}

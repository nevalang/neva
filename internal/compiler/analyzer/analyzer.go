// Package analyzer implements source code analysis.
package analyzer

import "github.com/nevalang/neva/internal/compiler/src"

type Analyzer struct{}

func (a Analyzer) Analyze(prog src.Program) (src.File, error) {
	return src.File{}, nil
}

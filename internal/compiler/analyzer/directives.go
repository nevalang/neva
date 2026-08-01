package analyzer

import (
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/ast"
)

func validateDirectiveCardinality(directives src.Directives) *compiler.Error {
	duplicate, ok := directives.FirstDuplicate()
	if !ok {
		return nil
	}

	return &compiler.Error{
		Message: "Duplicate #" + string(duplicate.Kind) + " directive",
		Meta:    &duplicate.Meta,
	}
}

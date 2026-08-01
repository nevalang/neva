package analyzer

import (
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/ast"
)

// validateDirectiveCardinality rejects repeated directive variants.
func validateDirectiveCardinality(directives src.Directives) *compiler.Error {
	seen := make(map[src.DirectiveKind]struct{}, len(directives))
	for i := range directives {
		directive := directives[i]
		kind := directive.Kind()
		if _, ok := seen[kind]; ok {
			return &compiler.Error{
				Message: "Duplicate #" + string(kind) + " directive",
				Meta:    &directive.Meta,
			}
		}
		seen[kind] = struct{}{}
	}

	return nil
}

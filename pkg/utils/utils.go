package utils

import (
	"strings"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

// Pointer allows to avoid creating of temporary variables just to take pointers.
func Pointer[T any](v T) *T {
	return &v
}

// ParseRef assumes string-ref has form of <pkg_name>.<entity_name≥ or just <entity_name>.
func ParseRef(ref string) src.EntityRef {
	entityRef := src.EntityRef{
		Meta: src.Meta{Text: ref},
	}

	parts := strings.Split(ref, ".")
	if len(parts) == 2 {
		entityRef.Pkg = parts[0]
		entityRef.Name = parts[1]
	} else {
		entityRef.Name = ref
	}

	return entityRef
}

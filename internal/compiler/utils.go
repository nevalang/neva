package compiler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

// Pointer allows to avoid creating of temporary variables just to take pointers.
func Pointer[T any](v T) *T {
	return &v
}

// ParseEntityRef assumes string-ref has form of <pkg_name>.<entity_name≥ or just <entity_name>.
func ParseEntityRef(ref string) core.EntityRef {
	entityRef := core.EntityRef{
		Meta: core.Meta{Text: ref},
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

// JSONDump is for debugging purposes only!
func JSONDump(v any) string {
	bb, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(bb)
}

func SaveFilesToDir(dst string, files map[string][]byte) error {
	for path, content := range files {
		filePath := filepath.Join(dst, path)
		dirPath := filepath.Dir(filePath)

		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(filePath, content, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// Package funcs implements low-level components (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

// CreatorRegistry allows to create runtime functions by accessing function creators by key.
func CreatorRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		"Reader":         reader{},
		"Printer":        printer{},
		"Blocker":        blocker{},
		"Emitter":        emitter{},
		"IntAdder":       intAdder{},
		"IntParser":      intParser{},
		"Destructor":     destructor{},
		"StructSelector": structSelector{},
		"MapSelector":    mapSelector{},
		"StructBuilder":  structBuilder{},
	}
}

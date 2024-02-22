// Package funcs implements low-level components (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

func CreatorRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		"reader":          reader{},
		"printer":         printer{},
		"blocker":         blocker{},
		"emitter":         emitter{},
		"int_adder":       intAdder{},
		"int_parser":      intParser{},
		"destructor":      destructor{},
		"struct_selector": structSelector{},
		"map_selector":    mapSelector{},
		"struct_builder":  structBuilder{},
	}
}

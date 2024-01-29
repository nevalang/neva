// Package funcs implements low-level components (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

// CreatorRegistry allows to create runtime functions by accessing function creators by key.
func CreatorRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		"Read":           read{},
		"Print":          print{},
		"Lock":           lock{},
		"Const":          constant{},
		"AddInts":        addInts{},
		// "AddFloats":      addFloats{},
		"ParseInt":       parseInt{},
		"Void":           void{},
		"StructSelector": structSelector{},
		"MapSelector":    mapSelector{},
		"StructBuilder":  structBuilder{},
	}
}

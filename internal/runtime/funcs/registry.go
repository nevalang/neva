// Package funcs implements low-level components (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

func CreatorRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		// basic stuff
		"emitter":       emitter{},
		"destructor":    destructor{},
		"blocker":       blocker{},
		"port_streamer": portStreamer{},
		// for structures
		"struct_selector": structSelector{},
		"struct_builder":  structBuilder{},
		// io
		"line_scanner": lineScanner{},
		"line_printer": linePrinter{},
		// math
		"int_adder": intAdder{},
		// parsing
		"int_parser": intParser{},
	}
}

// Package funcs implements low-level components (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

func CreatorRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		// core
		"emitter":       emitter{},
		"destructor":    destructor{},
		"blocker":       blocker{},
		"port_streamer": portStreamer{},
		// structures
		"struct_selector": structSelector{},
		"struct_builder":  structBuilder{},
		// logic
		"int_eq": intEq{},
		"match":  match{},
		"unwrap": unwrap{},
		// math
		"int_adder":       intAdder{},
		"int_subtractor":  intSubtractor{},
		"int_multiplier":  intMultiplier{},
		"int_decrementor": intDecrementor{},
		"range":           ranger{},
		"mod":             mod{},
		// io
		"line_scanner":  lineScanner{},
		"line_printer":  linePrinter{},
		"line_fprinter": lineFPrinter{},
		// strings
		"int_parser": intParser{},
		// regexp
		"regexp_submatcher": regexpSubmatcher{},
		//list
		"list_len":    listlen{},
		"list_push":   listpush{},
		"int_sort":    listsortint{},
		"float_sort":  listsortfloat{},
		"string_sort": listsortstring{},
		"indexer":     index{},
		"join":        join{},
		// time
		"time_sleep": timeSleep{},
		//string
		"split": split{},
	}
}

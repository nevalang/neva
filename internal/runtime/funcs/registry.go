// Package funcs implements low-level components (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

func CreatorRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		// core
		"new":            new{},
		"del":            del{},
		"lock":           lock{},
		"range":          ranger{},
		"port_sequencer": portSequencer{},
		// structures
		"field":          field{},
		"struct_builder": structBuilder{},
		// logic
		"int_eq": intEq{},
		"match":  match{},
		"unwrap": unwrap{},
		// math
		"int_add":  intAdd{},
		"int_sub":  intSub{},
		"int_mul":  intMul{},
		"int_decr": intDecr{},
		"int_mod":  intMod{},
		// io
		"scanln":  scanln{},
		"println": println{},
		"printf":  printf{},
		// strings
		"int_parse": intParse{},
		// regexp
		"regexp_submatch": regexpSubmatch{},
		//list
		"list_len": listlen{},
		"index":    index{},
		// time
		"time_sleep": timeSleep{},
	}
}

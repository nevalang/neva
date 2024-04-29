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
		"port_sequencer": portStreamer{},
		// structures
		"field":          field{},
		"struct_builder": structBuilder{},
		// logic
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
		"parse_int": parseInt{},
		// regexp
		"regexp_submatch": regexpSubmatch{},
		//list
		"index":       index{},
		"list_len":    listlen{},
		"list_iter":   list_iter{},
		"list_push":   listPush{},
		"int_sort":    listSortInt{},
		"float_sort":  listSortFloat{},
		"string_sort": listSortString{},
		"join":        stringJoin{},
		// time
		"time_sleep": timeSleep{},
		//string
		"split": stringSplit{},
		// io/file
		"read_all":  readAll{},
		"write_all": writeAll{},
	}
}

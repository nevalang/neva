// Package funcs implements low-level components (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

func CreatorRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		// core
		"new":    new{},
		"del":    del{},
		"lock":   lock{},
		"match":  match{},
		"unwrap": unwrap{},

		// runtime
		"panic": panicker{},

		// streamers
		"array_port_to_stream": arrayPortToStream{},
		"list_to_stream":       listToStream{},
		"stream_int_range":     streamIntRange{},

		// builders
		"struct_builder": structBuilder{},
		"stream_to_list": streamToList{},

		// structures
		"field": readStructField{},

		// math
		"int_add":  intAdd{},
		"int_sub":  intSub{},
		"int_mul":  intMul{},
		"int_decr": intDecr{},
		"int_mod":  intMod{},

		// strconv
		"parse_int": parseInt{},

		// regexp
		"regexp_submatch": regexpSubmatch{},

		// list
		"index":      index{},
		"list_len":   listlen{},
		"list_push":  listPush{},
		"int_sort":   listSortInt{},
		"float_sort": listSortFloat{},

		// time
		"time_sleep": timeSleep{},

		// strings
		"join":        stringJoin{},
		"split":       stringSplit{},
		"string_sort": listSortString{},

		// io
		"scanln":  scanln{},
		"println": println{},
		"printf":  printf{},

		// io/file
		"read_all":  readAll{},
		"write_all": writeAll{},
		// http
		"http_get": httpGet{},
		// image
		"image_encode": imageEncode{},
		"image_new":    imageNew{},
	}
}

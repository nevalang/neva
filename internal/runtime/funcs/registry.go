// Package funcs implements low-level flows (runtime functions).
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
		"unwrap": unwrap{},

		// runtime
		"panic": panicker{},

		// logic
		"match":             match{},
		"eq":                eq{},
		"if":                if_{},
		"not":               not{},
		"and":               and{},
		"or":                or{},
		"int_is_greater":    intIsGreater{},
		"int_is_lesser":     intIsLesser{},
		"string_is_greater": strIsGreater{},
		"string_is_lesser":  strIsLesser{},
		"float_is_greater":  floatIsGreater{},
		"float_is_lesser":   floatIsLesser{},
		"int_is_even":       intIsEven{},

		// streamers
		"array_port_to_stream": arrayPortToStream{},
		"list_to_stream":       listToStream{},
		"stream_int_range":     streamIntRange{},
		"stream_product":       streamProduct{},
		"stream_zip":           streamZip{},

		// builders
		"struct_builder": structBuilder{},
		"stream_to_list": streamToList{},

		// structures
		"field": readStructField{},

		// math
		"int_add":     intAdd{},
		"int_sub":     intSub{},
		"int_mul":     intMul{},
		"int_div":     intDiv{},
		"float_div":   floatDiv{},
		"int_decr":    intDecr{},
		"int_mod":     intMod{},
		"int_casemod": intCaseMod{},

		// strconv
		"parse_int": parseInt{},

		// regexp
		"regexp_submatch": regexpSubmatch{},

		// list
		"list_at":    listAt{},
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
		"args":    args{},
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

		// sync
		"wait_group": waitGroup{},
	}
}

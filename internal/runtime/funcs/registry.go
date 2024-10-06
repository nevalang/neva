// Package funcs implements low-level flows (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

func NewRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		// core
		"new":     new{},
		"del":     del{},
		"lock":    lock{},
		"unwrap":  unwrap{},
		"fan_out": fanOut{},
		"fan_in":  fanIn{},

		// runtime
		"panic": panicker{},

		// logic
		"switch":            switcher{},
		"match":             match{},
		"select":            selector{},
		"ternary":           ternary{},
		"eq":                eq{},
		"if":                if_{},
		"cond":              cond{},
		"not":               not{},
		"and":               and{},
		"or":                or{},
		"int_is_greater":    intIsGreater{},
		"int_is_lesser":     intIsLesser{},
		"string_is_greater": strIsGreater{},
		"string_is_lesser":  strIsLesser{},
		"float_is_greater":  floatIsGreater{},
		"float_is_lesser":   floatIsLesser{},

		// streamers
		"array_port_to_stream": arrayPortToStream{},
		"list_to_stream":       listToStream{},
		"stream_int_range":     rangeInt{},
		"stream_product":       streamProduct{},
		"stream_zip":           streamZip{},

		// builders
		"struct_builder": structBuilder{},
		"stream_to_list": streamToList{},

		// structures
		"field": readStructField{},

		// dictionaries
		"get_dict_value": getDictValue{},

		// math reducers
		"int_add":    intAdd{},
		"int_sub":    intSub{},
		"int_mul":    intMul{},
		"int_div":    intDiv{},
		"float_add":  floatAdd{},
		"float_sub":  floatSub{},
		"float_mul":  floatMul{},
		"float_div":  floatDiv{},
		"string_add": stringAdd{},

		// math mappers
		"int_inc":  intInc{},
		"int_decr": intDecr{},
		"int_mod":  intMod{},

		// strconv
		"parse_int": parseInt{},

		// regexp
		"regexp_submatch": regexpSubmatch{},

		// list
		"list_at":   listAt{},
		"list_len":  listlen{},
		"list_push": listPush{},

		// time
		"time_delay": timeDelay{},
		"time_after": timeAfter{},

		// strings
		"string_at": stringAt{},
		"join":      stringJoin{},
		"split":     stringSplit{},

		// io
		"scanln":  scanln{},
		"args":    args{},
		"println": println{},
		"printf":  printf{},

		// io/file
		"read_all":  fileReadAll{},
		"write_all": writeAll{},
		// http
		"http_get": httpGet{},
		// image
		"image_encode": imageEncode{},
		"image_new":    imageNew{},

		// sync
		"wait_group": waitGroup{},

		// other
		"accumulator": accumulator{},
	}
}

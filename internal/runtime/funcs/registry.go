// Package funcs implements low-level flows (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

func NewRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		"new":     new{},
		"new_v2":  newV2{},
		"del":     del{},
		"lock":    lock{},
		"unwrap":  unwrap{},
		"fan_out": fanOut{},
		"fan_in":  fanIn{},

		"panic": panicker{},

		"switch_router": switchRouter{},
		"match":         match{},
		"select":        selector{},
		"ternary":       ternary{},
		"eq":            eq{},
		"ne":            notEq{},
		"cond":          cond{},
		"not":           not{},
		"and":           and{},
		"or":            or{},

		"int_is_greater":          intIsGreater{},
		"int_is_greater_or_equal": intIsGreaterOrEqual{},

		"int_is_lesser":          intIsLesser{},
		"int_is_lesser_or_equal": intIsLesserOrEqual{},

		"string_is_greater": strIsGreater{},
		"string_is_lesser":  strIsLesser{},

		"float_is_greater": floatIsGreater{},
		"float_is_lesser":  floatIsLesser{},

		"array_port_to_stream": arrayPortToStream{},
		"list_to_stream":       listToStream{},
		"stream_int_range":     rangeInt{},
		"stream_int_range_v2":  rangeIntV2{},
		"stream_product":       streamProduct{},
		"stream_zip":           streamZip{},

		"struct_builder": structBuilder{},
		"stream_to_list": streamToList{},

		"field": readStructField{},

		"get_dict_value": getDictValue{},

		"int_add":    intAdd{},
		"int_sub":    intSub{},
		"int_mul":    intMul{},
		"int_div":    intDiv{},
		"float_add":  floatAdd{},
		"float_sub":  floatSub{},
		"float_mul":  floatMul{},
		"float_div":  floatDiv{},
		"string_add": stringAdd{},

		"int_inc": intInc{},
		"int_dec": intDec{},
		"int_mod": intMod{},

		"parse_int": parseInt{},

		"regexp_submatch": regexpSubmatch{},

		"list_at":   listAt{},
		"list_len":  listlen{},
		"list_push": listPush{},

		"time_delay": timeDelay{},
		"time_after": timeAfter{},

		"string_at":        stringAt{},
		"strings_join":     stringJoin{},
		"strings_split":    stringsSplit{},
		"strings_to_upper": stringsToUpper{},
		"strings_to_lower": stringsToLower{},

		"scanln":  scanln{},
		"args":    args{},
		"println": println{},
		"printf":  printf{},
		"print":   print{},

		"read_all":     fileReadAll{},
		"write_all":    writeAll{},
		"http_get":     httpGet{},
		"image_encode": imageEncode{},
		"image_new":    imageNew{},

		"wait_group": waitGroup{},

		"accumulator": accumulator{},

		"int_pow": intPow{},

		"int_bitwise_and": intBitwiseAnd{},
		"int_bitwise_or":  intBitwiseOr{},
		"int_bitwise_xor": intBitwiseXor{},
		"int_bitwise_lsh": intBitwiseLsh{},
		"int_bitwise_rsh": intBitwiseRsh{},

		"errors_new": errorsNew{},
	}
}

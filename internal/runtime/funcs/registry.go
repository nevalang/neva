// Package funcs implements low-level flows (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

func NewRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		"new":     newV1{},
		"new_v2":  newV2{},
		"del":     del{},
		"lock":    lock{},
		"fan_in":  fanIn{},
		"fan_out": fanOut{},

		"panic": panicker{},

		"switch_router": switchRouter{},
		"match":         matchSelector{},
		"select":        selector{},
		"ternary":       ternarySelector{},

		"eq":   eq{},
		"ne":   notEq{},
		"cond": cond{},
		"not":  not{},
		"and":  and{},
		"or":   or{},

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
		"stream_to_list":       streamToList{},

		"stream_int_range":    rangeIntV1{},
		"stream_int_range_v2": rangeIntV2{},

		"stream_zip":      streamZip{},
		"stream_zip_many": streamZipMany{},
		"stream_product":  streamProduct{},

		"field":          structField{},
		"struct_builder": structBuilder{},

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
		"int_mod":    intMod{},

		"int_pow": intPow{},

		"int_bitwise_and": intBitwiseAnd{},
		"int_bitwise_or":  intBitwiseOr{},
		"int_bitwise_xor": intBitwiseXor{},
		"int_bitwise_lsh": intBitwiseLsh{},
		"int_bitwise_rsh": intBitwiseRsh{},

		"int_inc": intInc{},
		"int_dec": intDec{},

		"atoi":        atoi{},
		"parse_int":   parseInt{},
		"parse_float": parseFloat{},

		"regexp_submatch": regexpSubmatch{},

		"list_at":   listAt{},
		"list_len":  listlen{},
		"list_push": listPush{},

		"time_delay": timeDelay{},
		"time_after": timeAfter{},

		"string_at":           stringAt{},
		"strings_join_stream": stringJoinStream{},
		"strings_join_list":   stringJoinList{},
		"strings_split":       stringsSplit{},
		"strings_to_upper":    stringsToUpper{},
		"strings_to_lower":    stringsToLower{},

		"scanln":      scanln{},
		"args":        args{},
		"dotenv_load": dotenvLoad{},
		"println":     println{},
		"printf":      printf{},
		"print":       print{},

		"read_all":     fileReadAll{},
		"write_all":    writeAll{},
		"http_get":     httpGet{},
		"image_encode": imageEncode{},
		"image_new":    imageNew{},

		"wait_group": waitGroup{},

		"accumulator": accumulator{},
		"state":       state{},
		"errors_new":  errorsNew{},

		"union_wrap_v1": unionWrapV1{},
		"union_wrap_v2": unionWrapV2{},
	}
}

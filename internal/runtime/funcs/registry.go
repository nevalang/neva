// Package funcs implements low-level flows (runtime functions).
// It has only one exported entity which is function creators registry.
package funcs

import (
	"github.com/nevalang/neva/internal/runtime"
)

func NewRegistry() map[string]runtime.FuncCreator {
	return map[string]runtime.FuncCreator{
		"new":     newV2{},
		"del":     del{},
		"lock":    lock{},
		"fan_in":  fanIn{},
		"fan_out": fanOut{},

		"panic": panicker{},

		"switch_router": switchRouter{},
		"match":         matchSelector{},
		"select":        selector{},
		"ternary":       ternarySelector{},
		"race":          race{},

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
		"array_port_to_list":   arrayPortToList{},
		"list_to_stream":       listToStream{},
		"string_to_stream":     stringToStream{},
		"stream_to_list":       streamToList{},
		"dict_to_stream":       dictToStream{},
		"stream_to_dict":       streamToDict{},

		"stream_int_range":           rangeInt{},
		"stream_just":                streamJust{},
		"stream_enumerate":           streamEnumerate{},
		"stream_for_each_controller": streamForEachController{},
		"stream_is_data":             streamIsData{},
		"stream_is_close":            streamIsClose{},
		"stream_unwrap_data":         streamUnwrapData{},

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

		"int_inc":                   intInc{},
		"int_dec":                   intDec{},
		"int_from_float":            intFromFloat{},
		"float_from_int":            floatFromInt{},
		"string_from_int_codepoint": stringFromIntCodepoint{},
		"bytes_from_string":         bytesFromString{},
		"strings_from_bytes":        stringsFromBytes{},

		"atoi":         atoi{},
		"itoa":         itoa{},
		"parse_bool":   parseBool{},
		"parse_int":    parseInt{},
		"parse_float":  parseFloat{},
		"format_bool":  formatBool{},
		"format_int":   formatInt{},
		"format_float": formatFloat{},

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

		"scanln":                    scanln{},
		"args":                      args{},
		"os_environ":                osEnviron{},
		"dotenv_load":               dotenvLoad{},
		"dotenv_load_from":          dotenvLoadFrom{},
		"dotenv_load_override":      dotenvLoad{override: true},
		"dotenv_load_from_override": dotenvLoadFrom{override: true},
		"println":                   printlnFunc{},
		"printf":                    printf{},
		"print":                     printFunc{},

		"read_all":     fileReadAll{},
		"write_all":    writeAll{},
		"http_get":     httpGet{},
		"image_encode": imageEncode{},
		"image_new":    imageNew{},

		"wait_group": waitGroup{},

		// "state":       state{},
		"accumulator": accumulator{},
		"errors_new":  errorsNew{},

		"union_wrap": unionWrapper{},
	}
}

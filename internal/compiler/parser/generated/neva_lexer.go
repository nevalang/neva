// Code generated from ./neva.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parsing

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type nevaLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var NevaLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func nevalexerLexerInit() {
	staticData := &NevaLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "'pub'", "'type'", "'struct'", "'union'", "'interface'", "'const'",
		"'def'", "'import'", "'true'", "'false'", "'-'", "'/'", "'='", "'<'",
		"'>'", "'('", "')'", "'{'", "'}'", "'['", "']'", "','", "':'", "'.'",
		"'?'", "'$'", "'@'", "'#'", "'::'", "'->'", "'=>'", "'---'",
	}
	staticData.SymbolicNames = []string{
		"", "PUB", "TYPE", "STRUCT", "UNION", "INTERFACE", "CONST", "DEF", "IMPORT",
		"TRUE", "FALSE", "MINUS", "SLASH", "EQ", "LT", "GT", "LPAREN", "RPAREN",
		"LBRACE", "RBRACE", "LBRACK", "RBRACK", "COMMA", "COLON", "DOT", "QUEST",
		"DOLLAR", "AT", "HASH", "DCOLON", "ARROW", "FAT_ARROW", "DASH3", "IDENTIFIER",
		"INT", "FLOAT", "STRING", "COMMENT", "NEWLINE", "WS",
	}
	staticData.RuleNames = []string{
		"PUB", "TYPE", "STRUCT", "UNION", "INTERFACE", "CONST", "DEF", "IMPORT",
		"TRUE", "FALSE", "MINUS", "SLASH", "EQ", "LT", "GT", "LPAREN", "RPAREN",
		"LBRACE", "RBRACE", "LBRACK", "RBRACK", "COMMA", "COLON", "DOT", "QUEST",
		"DOLLAR", "AT", "HASH", "DCOLON", "ARROW", "FAT_ARROW", "DASH3", "IDENTIFIER",
		"INT", "FLOAT", "STRING", "COMMENT", "NEWLINE", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 39, 242, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2,
		31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36,
		7, 36, 2, 37, 7, 37, 2, 38, 7, 38, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4,
		1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6,
		1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8,
		1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 11, 1, 11, 1, 12,
		1, 12, 1, 13, 1, 13, 1, 14, 1, 14, 1, 15, 1, 15, 1, 16, 1, 16, 1, 17, 1,
		17, 1, 18, 1, 18, 1, 19, 1, 19, 1, 20, 1, 20, 1, 21, 1, 21, 1, 22, 1, 22,
		1, 23, 1, 23, 1, 24, 1, 24, 1, 25, 1, 25, 1, 26, 1, 26, 1, 27, 1, 27, 1,
		28, 1, 28, 1, 28, 1, 29, 1, 29, 1, 29, 1, 30, 1, 30, 1, 30, 1, 31, 1, 31,
		1, 31, 1, 31, 1, 32, 1, 32, 5, 32, 191, 8, 32, 10, 32, 12, 32, 194, 9,
		32, 1, 33, 4, 33, 197, 8, 33, 11, 33, 12, 33, 198, 1, 34, 5, 34, 202, 8,
		34, 10, 34, 12, 34, 205, 9, 34, 1, 34, 1, 34, 4, 34, 209, 8, 34, 11, 34,
		12, 34, 210, 1, 35, 1, 35, 5, 35, 215, 8, 35, 10, 35, 12, 35, 218, 9, 35,
		1, 35, 1, 35, 1, 36, 1, 36, 1, 36, 1, 36, 5, 36, 226, 8, 36, 10, 36, 12,
		36, 229, 9, 36, 1, 37, 3, 37, 232, 8, 37, 1, 37, 1, 37, 1, 38, 4, 38, 237,
		8, 38, 11, 38, 12, 38, 238, 1, 38, 1, 38, 1, 216, 0, 39, 1, 1, 3, 2, 5,
		3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12, 25,
		13, 27, 14, 29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 39, 20, 41, 21, 43,
		22, 45, 23, 47, 24, 49, 25, 51, 26, 53, 27, 55, 28, 57, 29, 59, 30, 61,
		31, 63, 32, 65, 33, 67, 34, 69, 35, 71, 36, 73, 37, 75, 38, 77, 39, 1,
		0, 5, 3, 0, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97,
		122, 1, 0, 48, 57, 2, 0, 10, 10, 13, 13, 2, 0, 9, 9, 32, 32, 249, 0, 1,
		1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9,
		1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0,
		17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0,
		0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0,
		0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0,
		0, 0, 0, 41, 1, 0, 0, 0, 0, 43, 1, 0, 0, 0, 0, 45, 1, 0, 0, 0, 0, 47, 1,
		0, 0, 0, 0, 49, 1, 0, 0, 0, 0, 51, 1, 0, 0, 0, 0, 53, 1, 0, 0, 0, 0, 55,
		1, 0, 0, 0, 0, 57, 1, 0, 0, 0, 0, 59, 1, 0, 0, 0, 0, 61, 1, 0, 0, 0, 0,
		63, 1, 0, 0, 0, 0, 65, 1, 0, 0, 0, 0, 67, 1, 0, 0, 0, 0, 69, 1, 0, 0, 0,
		0, 71, 1, 0, 0, 0, 0, 73, 1, 0, 0, 0, 0, 75, 1, 0, 0, 0, 0, 77, 1, 0, 0,
		0, 1, 79, 1, 0, 0, 0, 3, 83, 1, 0, 0, 0, 5, 88, 1, 0, 0, 0, 7, 95, 1, 0,
		0, 0, 9, 101, 1, 0, 0, 0, 11, 111, 1, 0, 0, 0, 13, 117, 1, 0, 0, 0, 15,
		121, 1, 0, 0, 0, 17, 128, 1, 0, 0, 0, 19, 133, 1, 0, 0, 0, 21, 139, 1,
		0, 0, 0, 23, 141, 1, 0, 0, 0, 25, 143, 1, 0, 0, 0, 27, 145, 1, 0, 0, 0,
		29, 147, 1, 0, 0, 0, 31, 149, 1, 0, 0, 0, 33, 151, 1, 0, 0, 0, 35, 153,
		1, 0, 0, 0, 37, 155, 1, 0, 0, 0, 39, 157, 1, 0, 0, 0, 41, 159, 1, 0, 0,
		0, 43, 161, 1, 0, 0, 0, 45, 163, 1, 0, 0, 0, 47, 165, 1, 0, 0, 0, 49, 167,
		1, 0, 0, 0, 51, 169, 1, 0, 0, 0, 53, 171, 1, 0, 0, 0, 55, 173, 1, 0, 0,
		0, 57, 175, 1, 0, 0, 0, 59, 178, 1, 0, 0, 0, 61, 181, 1, 0, 0, 0, 63, 184,
		1, 0, 0, 0, 65, 188, 1, 0, 0, 0, 67, 196, 1, 0, 0, 0, 69, 203, 1, 0, 0,
		0, 71, 212, 1, 0, 0, 0, 73, 221, 1, 0, 0, 0, 75, 231, 1, 0, 0, 0, 77, 236,
		1, 0, 0, 0, 79, 80, 5, 112, 0, 0, 80, 81, 5, 117, 0, 0, 81, 82, 5, 98,
		0, 0, 82, 2, 1, 0, 0, 0, 83, 84, 5, 116, 0, 0, 84, 85, 5, 121, 0, 0, 85,
		86, 5, 112, 0, 0, 86, 87, 5, 101, 0, 0, 87, 4, 1, 0, 0, 0, 88, 89, 5, 115,
		0, 0, 89, 90, 5, 116, 0, 0, 90, 91, 5, 114, 0, 0, 91, 92, 5, 117, 0, 0,
		92, 93, 5, 99, 0, 0, 93, 94, 5, 116, 0, 0, 94, 6, 1, 0, 0, 0, 95, 96, 5,
		117, 0, 0, 96, 97, 5, 110, 0, 0, 97, 98, 5, 105, 0, 0, 98, 99, 5, 111,
		0, 0, 99, 100, 5, 110, 0, 0, 100, 8, 1, 0, 0, 0, 101, 102, 5, 105, 0, 0,
		102, 103, 5, 110, 0, 0, 103, 104, 5, 116, 0, 0, 104, 105, 5, 101, 0, 0,
		105, 106, 5, 114, 0, 0, 106, 107, 5, 102, 0, 0, 107, 108, 5, 97, 0, 0,
		108, 109, 5, 99, 0, 0, 109, 110, 5, 101, 0, 0, 110, 10, 1, 0, 0, 0, 111,
		112, 5, 99, 0, 0, 112, 113, 5, 111, 0, 0, 113, 114, 5, 110, 0, 0, 114,
		115, 5, 115, 0, 0, 115, 116, 5, 116, 0, 0, 116, 12, 1, 0, 0, 0, 117, 118,
		5, 100, 0, 0, 118, 119, 5, 101, 0, 0, 119, 120, 5, 102, 0, 0, 120, 14,
		1, 0, 0, 0, 121, 122, 5, 105, 0, 0, 122, 123, 5, 109, 0, 0, 123, 124, 5,
		112, 0, 0, 124, 125, 5, 111, 0, 0, 125, 126, 5, 114, 0, 0, 126, 127, 5,
		116, 0, 0, 127, 16, 1, 0, 0, 0, 128, 129, 5, 116, 0, 0, 129, 130, 5, 114,
		0, 0, 130, 131, 5, 117, 0, 0, 131, 132, 5, 101, 0, 0, 132, 18, 1, 0, 0,
		0, 133, 134, 5, 102, 0, 0, 134, 135, 5, 97, 0, 0, 135, 136, 5, 108, 0,
		0, 136, 137, 5, 115, 0, 0, 137, 138, 5, 101, 0, 0, 138, 20, 1, 0, 0, 0,
		139, 140, 5, 45, 0, 0, 140, 22, 1, 0, 0, 0, 141, 142, 5, 47, 0, 0, 142,
		24, 1, 0, 0, 0, 143, 144, 5, 61, 0, 0, 144, 26, 1, 0, 0, 0, 145, 146, 5,
		60, 0, 0, 146, 28, 1, 0, 0, 0, 147, 148, 5, 62, 0, 0, 148, 30, 1, 0, 0,
		0, 149, 150, 5, 40, 0, 0, 150, 32, 1, 0, 0, 0, 151, 152, 5, 41, 0, 0, 152,
		34, 1, 0, 0, 0, 153, 154, 5, 123, 0, 0, 154, 36, 1, 0, 0, 0, 155, 156,
		5, 125, 0, 0, 156, 38, 1, 0, 0, 0, 157, 158, 5, 91, 0, 0, 158, 40, 1, 0,
		0, 0, 159, 160, 5, 93, 0, 0, 160, 42, 1, 0, 0, 0, 161, 162, 5, 44, 0, 0,
		162, 44, 1, 0, 0, 0, 163, 164, 5, 58, 0, 0, 164, 46, 1, 0, 0, 0, 165, 166,
		5, 46, 0, 0, 166, 48, 1, 0, 0, 0, 167, 168, 5, 63, 0, 0, 168, 50, 1, 0,
		0, 0, 169, 170, 5, 36, 0, 0, 170, 52, 1, 0, 0, 0, 171, 172, 5, 64, 0, 0,
		172, 54, 1, 0, 0, 0, 173, 174, 5, 35, 0, 0, 174, 56, 1, 0, 0, 0, 175, 176,
		5, 58, 0, 0, 176, 177, 5, 58, 0, 0, 177, 58, 1, 0, 0, 0, 178, 179, 5, 45,
		0, 0, 179, 180, 5, 62, 0, 0, 180, 60, 1, 0, 0, 0, 181, 182, 5, 61, 0, 0,
		182, 183, 5, 62, 0, 0, 183, 62, 1, 0, 0, 0, 184, 185, 5, 45, 0, 0, 185,
		186, 5, 45, 0, 0, 186, 187, 5, 45, 0, 0, 187, 64, 1, 0, 0, 0, 188, 192,
		7, 0, 0, 0, 189, 191, 7, 1, 0, 0, 190, 189, 1, 0, 0, 0, 191, 194, 1, 0,
		0, 0, 192, 190, 1, 0, 0, 0, 192, 193, 1, 0, 0, 0, 193, 66, 1, 0, 0, 0,
		194, 192, 1, 0, 0, 0, 195, 197, 7, 2, 0, 0, 196, 195, 1, 0, 0, 0, 197,
		198, 1, 0, 0, 0, 198, 196, 1, 0, 0, 0, 198, 199, 1, 0, 0, 0, 199, 68, 1,
		0, 0, 0, 200, 202, 7, 2, 0, 0, 201, 200, 1, 0, 0, 0, 202, 205, 1, 0, 0,
		0, 203, 201, 1, 0, 0, 0, 203, 204, 1, 0, 0, 0, 204, 206, 1, 0, 0, 0, 205,
		203, 1, 0, 0, 0, 206, 208, 5, 46, 0, 0, 207, 209, 7, 2, 0, 0, 208, 207,
		1, 0, 0, 0, 209, 210, 1, 0, 0, 0, 210, 208, 1, 0, 0, 0, 210, 211, 1, 0,
		0, 0, 211, 70, 1, 0, 0, 0, 212, 216, 5, 39, 0, 0, 213, 215, 9, 0, 0, 0,
		214, 213, 1, 0, 0, 0, 215, 218, 1, 0, 0, 0, 216, 217, 1, 0, 0, 0, 216,
		214, 1, 0, 0, 0, 217, 219, 1, 0, 0, 0, 218, 216, 1, 0, 0, 0, 219, 220,
		5, 39, 0, 0, 220, 72, 1, 0, 0, 0, 221, 222, 5, 47, 0, 0, 222, 223, 5, 47,
		0, 0, 223, 227, 1, 0, 0, 0, 224, 226, 8, 3, 0, 0, 225, 224, 1, 0, 0, 0,
		226, 229, 1, 0, 0, 0, 227, 225, 1, 0, 0, 0, 227, 228, 1, 0, 0, 0, 228,
		74, 1, 0, 0, 0, 229, 227, 1, 0, 0, 0, 230, 232, 5, 13, 0, 0, 231, 230,
		1, 0, 0, 0, 231, 232, 1, 0, 0, 0, 232, 233, 1, 0, 0, 0, 233, 234, 5, 10,
		0, 0, 234, 76, 1, 0, 0, 0, 235, 237, 7, 4, 0, 0, 236, 235, 1, 0, 0, 0,
		237, 238, 1, 0, 0, 0, 238, 236, 1, 0, 0, 0, 238, 239, 1, 0, 0, 0, 239,
		240, 1, 0, 0, 0, 240, 241, 6, 38, 0, 0, 241, 78, 1, 0, 0, 0, 9, 0, 192,
		198, 203, 210, 216, 227, 231, 238, 1, 6, 0, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// nevaLexerInit initializes any static state used to implement nevaLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewnevaLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func NevaLexerInit() {
	staticData := &NevaLexerLexerStaticData
	staticData.once.Do(nevalexerLexerInit)
}

// NewnevaLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewnevaLexer(input antlr.CharStream) *nevaLexer {
	NevaLexerInit()
	l := new(nevaLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &NevaLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "neva.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// nevaLexer tokens.
const (
	nevaLexerPUB        = 1
	nevaLexerTYPE       = 2
	nevaLexerSTRUCT     = 3
	nevaLexerUNION      = 4
	nevaLexerINTERFACE  = 5
	nevaLexerCONST      = 6
	nevaLexerDEF        = 7
	nevaLexerIMPORT     = 8
	nevaLexerTRUE       = 9
	nevaLexerFALSE      = 10
	nevaLexerMINUS      = 11
	nevaLexerSLASH      = 12
	nevaLexerEQ         = 13
	nevaLexerLT         = 14
	nevaLexerGT         = 15
	nevaLexerLPAREN     = 16
	nevaLexerRPAREN     = 17
	nevaLexerLBRACE     = 18
	nevaLexerRBRACE     = 19
	nevaLexerLBRACK     = 20
	nevaLexerRBRACK     = 21
	nevaLexerCOMMA      = 22
	nevaLexerCOLON      = 23
	nevaLexerDOT        = 24
	nevaLexerQUEST      = 25
	nevaLexerDOLLAR     = 26
	nevaLexerAT         = 27
	nevaLexerHASH       = 28
	nevaLexerDCOLON     = 29
	nevaLexerARROW      = 30
	nevaLexerFAT_ARROW  = 31
	nevaLexerDASH3      = 32
	nevaLexerIDENTIFIER = 33
	nevaLexerINT        = 34
	nevaLexerFLOAT      = 35
	nevaLexerSTRING     = 36
	nevaLexerCOMMENT    = 37
	nevaLexerNEWLINE    = 38
	nevaLexerWS         = 39
)

// Code generated from ./neva.g4 by ANTLR 4.13.1. DO NOT EDIT.

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
		"", "'#'", "'('", "','", "')'", "'import'", "'{'", "'}'", "'/'", "'@'",
		"'.'", "'types'", "'<'", "'>'", "'enum'", "'['", "']'", "'struct'",
		"'|'", "'interfaces'", "'const'", "'true'", "'false'", "'nil'", "':'",
		"'components'", "'nodes'", "'net'", "'->'", "'$'", "", "'pub'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "COMMENT", "PUB_KW",
		"IDENTIFIER", "INT", "FLOAT", "STRING", "NEWLINE", "WS",
	}
	staticData.RuleNames = []string{
		"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8",
		"T__9", "T__10", "T__11", "T__12", "T__13", "T__14", "T__15", "T__16",
		"T__17", "T__18", "T__19", "T__20", "T__21", "T__22", "T__23", "T__24",
		"T__25", "T__26", "T__27", "T__28", "COMMENT", "PUB_KW", "IDENTIFIER",
		"LETTER", "INT", "FLOAT", "STRING", "NEWLINE", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 37, 251, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2,
		31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36,
		7, 36, 2, 37, 7, 37, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1,
		4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1,
		7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1,
		11, 1, 11, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14,
		1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 17, 1,
		17, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18,
		1, 18, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 1,
		20, 1, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22,
		1, 22, 1, 23, 1, 23, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1,
		24, 1, 24, 1, 24, 1, 24, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 26,
		1, 26, 1, 26, 1, 26, 1, 27, 1, 27, 1, 27, 1, 28, 1, 28, 1, 29, 1, 29, 1,
		29, 1, 29, 5, 29, 195, 8, 29, 10, 29, 12, 29, 198, 9, 29, 1, 30, 1, 30,
		1, 30, 1, 30, 1, 31, 1, 31, 1, 31, 5, 31, 207, 8, 31, 10, 31, 12, 31, 210,
		9, 31, 1, 32, 1, 32, 1, 33, 4, 33, 215, 8, 33, 11, 33, 12, 33, 216, 1,
		34, 5, 34, 220, 8, 34, 10, 34, 12, 34, 223, 9, 34, 1, 34, 1, 34, 4, 34,
		227, 8, 34, 11, 34, 12, 34, 228, 1, 35, 1, 35, 5, 35, 233, 8, 35, 10, 35,
		12, 35, 236, 9, 35, 1, 35, 1, 35, 1, 36, 3, 36, 241, 8, 36, 1, 36, 1, 36,
		1, 37, 4, 37, 246, 8, 37, 11, 37, 12, 37, 247, 1, 37, 1, 37, 1, 234, 0,
		38, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21,
		11, 23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 39,
		20, 41, 21, 43, 22, 45, 23, 47, 24, 49, 25, 51, 26, 53, 27, 55, 28, 57,
		29, 59, 30, 61, 31, 63, 32, 65, 0, 67, 33, 69, 34, 71, 35, 73, 36, 75,
		37, 1, 0, 4, 2, 0, 10, 10, 13, 13, 3, 0, 65, 90, 95, 95, 97, 122, 1, 0,
		48, 57, 2, 0, 9, 9, 32, 32, 258, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0,
		5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0,
		13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0,
		0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0,
		0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0,
		0, 0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0, 43, 1,
		0, 0, 0, 0, 45, 1, 0, 0, 0, 0, 47, 1, 0, 0, 0, 0, 49, 1, 0, 0, 0, 0, 51,
		1, 0, 0, 0, 0, 53, 1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 0, 57, 1, 0, 0, 0, 0,
		59, 1, 0, 0, 0, 0, 61, 1, 0, 0, 0, 0, 63, 1, 0, 0, 0, 0, 67, 1, 0, 0, 0,
		0, 69, 1, 0, 0, 0, 0, 71, 1, 0, 0, 0, 0, 73, 1, 0, 0, 0, 0, 75, 1, 0, 0,
		0, 1, 77, 1, 0, 0, 0, 3, 79, 1, 0, 0, 0, 5, 81, 1, 0, 0, 0, 7, 83, 1, 0,
		0, 0, 9, 85, 1, 0, 0, 0, 11, 92, 1, 0, 0, 0, 13, 94, 1, 0, 0, 0, 15, 96,
		1, 0, 0, 0, 17, 98, 1, 0, 0, 0, 19, 100, 1, 0, 0, 0, 21, 102, 1, 0, 0,
		0, 23, 108, 1, 0, 0, 0, 25, 110, 1, 0, 0, 0, 27, 112, 1, 0, 0, 0, 29, 117,
		1, 0, 0, 0, 31, 119, 1, 0, 0, 0, 33, 121, 1, 0, 0, 0, 35, 128, 1, 0, 0,
		0, 37, 130, 1, 0, 0, 0, 39, 141, 1, 0, 0, 0, 41, 147, 1, 0, 0, 0, 43, 152,
		1, 0, 0, 0, 45, 158, 1, 0, 0, 0, 47, 162, 1, 0, 0, 0, 49, 164, 1, 0, 0,
		0, 51, 175, 1, 0, 0, 0, 53, 181, 1, 0, 0, 0, 55, 185, 1, 0, 0, 0, 57, 188,
		1, 0, 0, 0, 59, 190, 1, 0, 0, 0, 61, 199, 1, 0, 0, 0, 63, 203, 1, 0, 0,
		0, 65, 211, 1, 0, 0, 0, 67, 214, 1, 0, 0, 0, 69, 221, 1, 0, 0, 0, 71, 230,
		1, 0, 0, 0, 73, 240, 1, 0, 0, 0, 75, 245, 1, 0, 0, 0, 77, 78, 5, 35, 0,
		0, 78, 2, 1, 0, 0, 0, 79, 80, 5, 40, 0, 0, 80, 4, 1, 0, 0, 0, 81, 82, 5,
		44, 0, 0, 82, 6, 1, 0, 0, 0, 83, 84, 5, 41, 0, 0, 84, 8, 1, 0, 0, 0, 85,
		86, 5, 105, 0, 0, 86, 87, 5, 109, 0, 0, 87, 88, 5, 112, 0, 0, 88, 89, 5,
		111, 0, 0, 89, 90, 5, 114, 0, 0, 90, 91, 5, 116, 0, 0, 91, 10, 1, 0, 0,
		0, 92, 93, 5, 123, 0, 0, 93, 12, 1, 0, 0, 0, 94, 95, 5, 125, 0, 0, 95,
		14, 1, 0, 0, 0, 96, 97, 5, 47, 0, 0, 97, 16, 1, 0, 0, 0, 98, 99, 5, 64,
		0, 0, 99, 18, 1, 0, 0, 0, 100, 101, 5, 46, 0, 0, 101, 20, 1, 0, 0, 0, 102,
		103, 5, 116, 0, 0, 103, 104, 5, 121, 0, 0, 104, 105, 5, 112, 0, 0, 105,
		106, 5, 101, 0, 0, 106, 107, 5, 115, 0, 0, 107, 22, 1, 0, 0, 0, 108, 109,
		5, 60, 0, 0, 109, 24, 1, 0, 0, 0, 110, 111, 5, 62, 0, 0, 111, 26, 1, 0,
		0, 0, 112, 113, 5, 101, 0, 0, 113, 114, 5, 110, 0, 0, 114, 115, 5, 117,
		0, 0, 115, 116, 5, 109, 0, 0, 116, 28, 1, 0, 0, 0, 117, 118, 5, 91, 0,
		0, 118, 30, 1, 0, 0, 0, 119, 120, 5, 93, 0, 0, 120, 32, 1, 0, 0, 0, 121,
		122, 5, 115, 0, 0, 122, 123, 5, 116, 0, 0, 123, 124, 5, 114, 0, 0, 124,
		125, 5, 117, 0, 0, 125, 126, 5, 99, 0, 0, 126, 127, 5, 116, 0, 0, 127,
		34, 1, 0, 0, 0, 128, 129, 5, 124, 0, 0, 129, 36, 1, 0, 0, 0, 130, 131,
		5, 105, 0, 0, 131, 132, 5, 110, 0, 0, 132, 133, 5, 116, 0, 0, 133, 134,
		5, 101, 0, 0, 134, 135, 5, 114, 0, 0, 135, 136, 5, 102, 0, 0, 136, 137,
		5, 97, 0, 0, 137, 138, 5, 99, 0, 0, 138, 139, 5, 101, 0, 0, 139, 140, 5,
		115, 0, 0, 140, 38, 1, 0, 0, 0, 141, 142, 5, 99, 0, 0, 142, 143, 5, 111,
		0, 0, 143, 144, 5, 110, 0, 0, 144, 145, 5, 115, 0, 0, 145, 146, 5, 116,
		0, 0, 146, 40, 1, 0, 0, 0, 147, 148, 5, 116, 0, 0, 148, 149, 5, 114, 0,
		0, 149, 150, 5, 117, 0, 0, 150, 151, 5, 101, 0, 0, 151, 42, 1, 0, 0, 0,
		152, 153, 5, 102, 0, 0, 153, 154, 5, 97, 0, 0, 154, 155, 5, 108, 0, 0,
		155, 156, 5, 115, 0, 0, 156, 157, 5, 101, 0, 0, 157, 44, 1, 0, 0, 0, 158,
		159, 5, 110, 0, 0, 159, 160, 5, 105, 0, 0, 160, 161, 5, 108, 0, 0, 161,
		46, 1, 0, 0, 0, 162, 163, 5, 58, 0, 0, 163, 48, 1, 0, 0, 0, 164, 165, 5,
		99, 0, 0, 165, 166, 5, 111, 0, 0, 166, 167, 5, 109, 0, 0, 167, 168, 5,
		112, 0, 0, 168, 169, 5, 111, 0, 0, 169, 170, 5, 110, 0, 0, 170, 171, 5,
		101, 0, 0, 171, 172, 5, 110, 0, 0, 172, 173, 5, 116, 0, 0, 173, 174, 5,
		115, 0, 0, 174, 50, 1, 0, 0, 0, 175, 176, 5, 110, 0, 0, 176, 177, 5, 111,
		0, 0, 177, 178, 5, 100, 0, 0, 178, 179, 5, 101, 0, 0, 179, 180, 5, 115,
		0, 0, 180, 52, 1, 0, 0, 0, 181, 182, 5, 110, 0, 0, 182, 183, 5, 101, 0,
		0, 183, 184, 5, 116, 0, 0, 184, 54, 1, 0, 0, 0, 185, 186, 5, 45, 0, 0,
		186, 187, 5, 62, 0, 0, 187, 56, 1, 0, 0, 0, 188, 189, 5, 36, 0, 0, 189,
		58, 1, 0, 0, 0, 190, 191, 5, 47, 0, 0, 191, 192, 5, 47, 0, 0, 192, 196,
		1, 0, 0, 0, 193, 195, 8, 0, 0, 0, 194, 193, 1, 0, 0, 0, 195, 198, 1, 0,
		0, 0, 196, 194, 1, 0, 0, 0, 196, 197, 1, 0, 0, 0, 197, 60, 1, 0, 0, 0,
		198, 196, 1, 0, 0, 0, 199, 200, 5, 112, 0, 0, 200, 201, 5, 117, 0, 0, 201,
		202, 5, 98, 0, 0, 202, 62, 1, 0, 0, 0, 203, 208, 3, 65, 32, 0, 204, 207,
		3, 65, 32, 0, 205, 207, 3, 67, 33, 0, 206, 204, 1, 0, 0, 0, 206, 205, 1,
		0, 0, 0, 207, 210, 1, 0, 0, 0, 208, 206, 1, 0, 0, 0, 208, 209, 1, 0, 0,
		0, 209, 64, 1, 0, 0, 0, 210, 208, 1, 0, 0, 0, 211, 212, 7, 1, 0, 0, 212,
		66, 1, 0, 0, 0, 213, 215, 7, 2, 0, 0, 214, 213, 1, 0, 0, 0, 215, 216, 1,
		0, 0, 0, 216, 214, 1, 0, 0, 0, 216, 217, 1, 0, 0, 0, 217, 68, 1, 0, 0,
		0, 218, 220, 7, 2, 0, 0, 219, 218, 1, 0, 0, 0, 220, 223, 1, 0, 0, 0, 221,
		219, 1, 0, 0, 0, 221, 222, 1, 0, 0, 0, 222, 224, 1, 0, 0, 0, 223, 221,
		1, 0, 0, 0, 224, 226, 5, 46, 0, 0, 225, 227, 7, 2, 0, 0, 226, 225, 1, 0,
		0, 0, 227, 228, 1, 0, 0, 0, 228, 226, 1, 0, 0, 0, 228, 229, 1, 0, 0, 0,
		229, 70, 1, 0, 0, 0, 230, 234, 5, 39, 0, 0, 231, 233, 9, 0, 0, 0, 232,
		231, 1, 0, 0, 0, 233, 236, 1, 0, 0, 0, 234, 235, 1, 0, 0, 0, 234, 232,
		1, 0, 0, 0, 235, 237, 1, 0, 0, 0, 236, 234, 1, 0, 0, 0, 237, 238, 5, 39,
		0, 0, 238, 72, 1, 0, 0, 0, 239, 241, 5, 13, 0, 0, 240, 239, 1, 0, 0, 0,
		240, 241, 1, 0, 0, 0, 241, 242, 1, 0, 0, 0, 242, 243, 5, 10, 0, 0, 243,
		74, 1, 0, 0, 0, 244, 246, 7, 3, 0, 0, 245, 244, 1, 0, 0, 0, 246, 247, 1,
		0, 0, 0, 247, 245, 1, 0, 0, 0, 247, 248, 1, 0, 0, 0, 248, 249, 1, 0, 0,
		0, 249, 250, 6, 37, 0, 0, 250, 76, 1, 0, 0, 0, 10, 0, 196, 206, 208, 216,
		221, 228, 234, 240, 247, 1, 0, 1, 0,
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
	nevaLexerT__0       = 1
	nevaLexerT__1       = 2
	nevaLexerT__2       = 3
	nevaLexerT__3       = 4
	nevaLexerT__4       = 5
	nevaLexerT__5       = 6
	nevaLexerT__6       = 7
	nevaLexerT__7       = 8
	nevaLexerT__8       = 9
	nevaLexerT__9       = 10
	nevaLexerT__10      = 11
	nevaLexerT__11      = 12
	nevaLexerT__12      = 13
	nevaLexerT__13      = 14
	nevaLexerT__14      = 15
	nevaLexerT__15      = 16
	nevaLexerT__16      = 17
	nevaLexerT__17      = 18
	nevaLexerT__18      = 19
	nevaLexerT__19      = 20
	nevaLexerT__20      = 21
	nevaLexerT__21      = 22
	nevaLexerT__22      = 23
	nevaLexerT__23      = 24
	nevaLexerT__24      = 25
	nevaLexerT__25      = 26
	nevaLexerT__26      = 27
	nevaLexerT__27      = 28
	nevaLexerT__28      = 29
	nevaLexerCOMMENT    = 30
	nevaLexerPUB_KW     = 31
	nevaLexerIDENTIFIER = 32
	nevaLexerINT        = 33
	nevaLexerFLOAT      = 34
	nevaLexerSTRING     = 35
	nevaLexerNEWLINE    = 36
	nevaLexerWS         = 37
)

// Code generated from neva.g4 by ANTLR 4.13.0. DO NOT EDIT.

package parsing // neva
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type nevaParser struct {
	*antlr.BaseParser
}

var NevaParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func nevaParserInit() {
	staticData := &NevaParserStaticData
	staticData.LiteralNames = []string{
		"", "'//'", "'\\n'", "'/*'", "'*/'", "'use'", "'{'", "'}'", "'/'", "'type'",
		"'pub'", "'<'", "','", "'>'", "'['", "']'", "'|'", "'io'", "'('", "')'",
		"'const'", "'='", "'true'", "'false'", "'nil'", "':'", "'comp'", "'node'",
		"'.'", "'net'", "'->'", "'in'", "'out'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "IDENTIFIER",
		"INT", "FLOAT", "STRING", "NEWLINE", "WS",
	}
	staticData.RuleNames = []string{
		"prog", "comment", "singleLineComment", "multiLineComment", "stmt",
		"useStmt", "importList", "importDef", "importPath", "typeStmt", "typeDefList",
		"typeDef", "typeParams", "typeParam", "typeExpr", "typeInstExpr", "typeArgs",
		"typeLitExpr", "arrTypeExpr", "recTypeExpr", "recTypeFields", "recTypeField",
		"unionTypeExpr", "enumTypeExpr", "nonUnionTypeExpr", "ioStmt", "interfaceDefList",
		"interfaceDef", "portsDef", "portDefList", "portDef", "constStmt", "constDefList",
		"constDef", "constValue", "arrLit", "arrItems", "recLit", "recValueFields",
		"recValueField", "compStmt", "compDefList", "compDef", "compBody", "compNodesDef",
		"compNodeDefList", "absNodeDef", "concreteNodeDef", "concreteNodeInst",
		"nodeRef", "nodeArgs", "nodeArgList", "nodeArg", "compNetDef", "connDefList",
		"connDef", "portAddr", "portDirection", "connReceiverSide", "connReceivers",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 38, 506, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 2, 57, 7, 57,
		2, 58, 7, 58, 2, 59, 7, 59, 1, 0, 1, 0, 5, 0, 123, 8, 0, 10, 0, 12, 0,
		126, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 3, 1, 132, 8, 1, 1, 2, 1, 2, 5, 2, 136,
		8, 2, 10, 2, 12, 2, 139, 9, 2, 1, 2, 1, 2, 1, 3, 1, 3, 5, 3, 145, 8, 3,
		10, 3, 12, 3, 148, 9, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 156,
		8, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 5, 6, 166, 8, 6,
		10, 6, 12, 6, 169, 9, 6, 1, 7, 3, 7, 172, 8, 7, 1, 7, 1, 7, 1, 8, 1, 8,
		1, 8, 5, 8, 179, 8, 8, 10, 8, 12, 8, 182, 9, 8, 1, 9, 1, 9, 1, 9, 1, 9,
		1, 9, 1, 10, 1, 10, 1, 10, 5, 10, 192, 8, 10, 10, 10, 12, 10, 195, 9, 10,
		1, 11, 3, 11, 198, 8, 11, 1, 11, 1, 11, 3, 11, 202, 8, 11, 1, 11, 1, 11,
		1, 12, 1, 12, 1, 12, 1, 12, 3, 12, 210, 8, 12, 1, 12, 5, 12, 213, 8, 12,
		10, 12, 12, 12, 216, 9, 12, 1, 12, 1, 12, 1, 13, 1, 13, 3, 13, 222, 8,
		13, 1, 14, 1, 14, 1, 14, 3, 14, 227, 8, 14, 1, 15, 1, 15, 3, 15, 231, 8,
		15, 1, 16, 1, 16, 1, 16, 1, 16, 5, 16, 237, 8, 16, 10, 16, 12, 16, 240,
		9, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 3, 17, 247, 8, 17, 1, 18, 1,
		18, 1, 18, 1, 18, 1, 18, 1, 19, 1, 19, 3, 19, 256, 8, 19, 1, 19, 1, 19,
		1, 20, 1, 20, 1, 20, 5, 20, 263, 8, 20, 10, 20, 12, 20, 266, 9, 20, 1,
		21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 4, 22, 274, 8, 22, 11, 22, 12, 22,
		275, 1, 23, 1, 23, 1, 24, 1, 24, 3, 24, 282, 8, 24, 1, 25, 1, 25, 1, 25,
		1, 25, 1, 25, 1, 26, 1, 26, 1, 26, 5, 26, 292, 8, 26, 10, 26, 12, 26, 295,
		9, 26, 1, 27, 3, 27, 298, 8, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1,
		28, 1, 28, 1, 28, 1, 28, 1, 29, 1, 29, 1, 29, 3, 29, 312, 8, 29, 1, 29,
		5, 29, 315, 8, 29, 10, 29, 12, 29, 318, 9, 29, 1, 30, 1, 30, 1, 30, 1,
		31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 32, 1, 32, 1, 32, 5, 32, 332,
		8, 32, 10, 32, 12, 32, 335, 9, 32, 1, 33, 3, 33, 338, 8, 33, 1, 33, 1,
		33, 1, 33, 1, 33, 1, 33, 1, 34, 1, 34, 1, 34, 1, 34, 1, 34, 1, 34, 1, 34,
		1, 34, 3, 34, 353, 8, 34, 1, 35, 1, 35, 1, 35, 1, 35, 1, 36, 1, 36, 1,
		36, 1, 36, 3, 36, 363, 8, 36, 1, 36, 5, 36, 366, 8, 36, 10, 36, 12, 36,
		369, 9, 36, 3, 36, 371, 8, 36, 1, 37, 1, 37, 1, 37, 1, 37, 1, 38, 1, 38,
		1, 38, 3, 38, 380, 8, 38, 1, 38, 5, 38, 383, 8, 38, 10, 38, 12, 38, 386,
		9, 38, 1, 39, 1, 39, 1, 39, 1, 39, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1,
		40, 1, 41, 1, 41, 1, 41, 5, 41, 401, 8, 41, 10, 41, 12, 41, 404, 9, 41,
		1, 42, 3, 42, 407, 8, 42, 1, 42, 1, 42, 1, 42, 1, 43, 1, 43, 1, 43, 1,
		43, 1, 43, 3, 43, 417, 8, 43, 1, 44, 1, 44, 1, 44, 1, 44, 1, 44, 1, 45,
		1, 45, 3, 45, 426, 8, 45, 1, 46, 1, 46, 1, 46, 1, 47, 1, 47, 1, 47, 1,
		47, 1, 48, 1, 48, 1, 48, 1, 48, 1, 49, 1, 49, 1, 49, 5, 49, 442, 8, 49,
		10, 49, 12, 49, 445, 9, 49, 1, 50, 1, 50, 1, 50, 1, 50, 1, 51, 1, 51, 1,
		51, 3, 51, 454, 8, 51, 1, 51, 1, 51, 1, 52, 1, 52, 1, 53, 1, 53, 1, 53,
		1, 53, 1, 53, 1, 54, 1, 54, 1, 54, 5, 54, 468, 8, 54, 10, 54, 12, 54, 471,
		9, 54, 1, 55, 1, 55, 1, 55, 1, 55, 1, 56, 3, 56, 478, 8, 56, 1, 56, 1,
		56, 1, 56, 1, 56, 1, 56, 3, 56, 485, 8, 56, 3, 56, 487, 8, 56, 1, 57, 1,
		57, 1, 58, 1, 58, 3, 58, 493, 8, 58, 1, 59, 1, 59, 1, 59, 1, 59, 5, 59,
		499, 8, 59, 10, 59, 12, 59, 502, 9, 59, 1, 59, 1, 59, 1, 59, 1, 146, 0,
		60, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34,
		36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70,
		72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104,
		106, 108, 110, 112, 114, 116, 118, 0, 2, 1, 0, 2, 2, 1, 0, 31, 32, 502,
		0, 124, 1, 0, 0, 0, 2, 131, 1, 0, 0, 0, 4, 133, 1, 0, 0, 0, 6, 142, 1,
		0, 0, 0, 8, 155, 1, 0, 0, 0, 10, 157, 1, 0, 0, 0, 12, 162, 1, 0, 0, 0,
		14, 171, 1, 0, 0, 0, 16, 175, 1, 0, 0, 0, 18, 183, 1, 0, 0, 0, 20, 188,
		1, 0, 0, 0, 22, 197, 1, 0, 0, 0, 24, 205, 1, 0, 0, 0, 26, 219, 1, 0, 0,
		0, 28, 226, 1, 0, 0, 0, 30, 228, 1, 0, 0, 0, 32, 232, 1, 0, 0, 0, 34, 246,
		1, 0, 0, 0, 36, 248, 1, 0, 0, 0, 38, 253, 1, 0, 0, 0, 40, 259, 1, 0, 0,
		0, 42, 267, 1, 0, 0, 0, 44, 270, 1, 0, 0, 0, 46, 277, 1, 0, 0, 0, 48, 281,
		1, 0, 0, 0, 50, 283, 1, 0, 0, 0, 52, 288, 1, 0, 0, 0, 54, 297, 1, 0, 0,
		0, 56, 304, 1, 0, 0, 0, 58, 308, 1, 0, 0, 0, 60, 319, 1, 0, 0, 0, 62, 322,
		1, 0, 0, 0, 64, 328, 1, 0, 0, 0, 66, 337, 1, 0, 0, 0, 68, 352, 1, 0, 0,
		0, 70, 354, 1, 0, 0, 0, 72, 370, 1, 0, 0, 0, 74, 372, 1, 0, 0, 0, 76, 376,
		1, 0, 0, 0, 78, 387, 1, 0, 0, 0, 80, 391, 1, 0, 0, 0, 82, 397, 1, 0, 0,
		0, 84, 406, 1, 0, 0, 0, 86, 416, 1, 0, 0, 0, 88, 418, 1, 0, 0, 0, 90, 425,
		1, 0, 0, 0, 92, 427, 1, 0, 0, 0, 94, 430, 1, 0, 0, 0, 96, 434, 1, 0, 0,
		0, 98, 438, 1, 0, 0, 0, 100, 446, 1, 0, 0, 0, 102, 450, 1, 0, 0, 0, 104,
		457, 1, 0, 0, 0, 106, 459, 1, 0, 0, 0, 108, 464, 1, 0, 0, 0, 110, 472,
		1, 0, 0, 0, 112, 486, 1, 0, 0, 0, 114, 488, 1, 0, 0, 0, 116, 492, 1, 0,
		0, 0, 118, 494, 1, 0, 0, 0, 120, 123, 3, 2, 1, 0, 121, 123, 3, 8, 4, 0,
		122, 120, 1, 0, 0, 0, 122, 121, 1, 0, 0, 0, 123, 126, 1, 0, 0, 0, 124,
		122, 1, 0, 0, 0, 124, 125, 1, 0, 0, 0, 125, 127, 1, 0, 0, 0, 126, 124,
		1, 0, 0, 0, 127, 128, 5, 0, 0, 1, 128, 1, 1, 0, 0, 0, 129, 132, 3, 4, 2,
		0, 130, 132, 3, 6, 3, 0, 131, 129, 1, 0, 0, 0, 131, 130, 1, 0, 0, 0, 132,
		3, 1, 0, 0, 0, 133, 137, 5, 1, 0, 0, 134, 136, 8, 0, 0, 0, 135, 134, 1,
		0, 0, 0, 136, 139, 1, 0, 0, 0, 137, 135, 1, 0, 0, 0, 137, 138, 1, 0, 0,
		0, 138, 140, 1, 0, 0, 0, 139, 137, 1, 0, 0, 0, 140, 141, 5, 37, 0, 0, 141,
		5, 1, 0, 0, 0, 142, 146, 5, 3, 0, 0, 143, 145, 9, 0, 0, 0, 144, 143, 1,
		0, 0, 0, 145, 148, 1, 0, 0, 0, 146, 147, 1, 0, 0, 0, 146, 144, 1, 0, 0,
		0, 147, 149, 1, 0, 0, 0, 148, 146, 1, 0, 0, 0, 149, 150, 5, 4, 0, 0, 150,
		7, 1, 0, 0, 0, 151, 156, 3, 10, 5, 0, 152, 156, 3, 18, 9, 0, 153, 156,
		3, 50, 25, 0, 154, 156, 3, 62, 31, 0, 155, 151, 1, 0, 0, 0, 155, 152, 1,
		0, 0, 0, 155, 153, 1, 0, 0, 0, 155, 154, 1, 0, 0, 0, 156, 9, 1, 0, 0, 0,
		157, 158, 5, 5, 0, 0, 158, 159, 5, 6, 0, 0, 159, 160, 3, 12, 6, 0, 160,
		161, 5, 7, 0, 0, 161, 11, 1, 0, 0, 0, 162, 167, 3, 14, 7, 0, 163, 164,
		5, 37, 0, 0, 164, 166, 3, 14, 7, 0, 165, 163, 1, 0, 0, 0, 166, 169, 1,
		0, 0, 0, 167, 165, 1, 0, 0, 0, 167, 168, 1, 0, 0, 0, 168, 13, 1, 0, 0,
		0, 169, 167, 1, 0, 0, 0, 170, 172, 5, 33, 0, 0, 171, 170, 1, 0, 0, 0, 171,
		172, 1, 0, 0, 0, 172, 173, 1, 0, 0, 0, 173, 174, 3, 16, 8, 0, 174, 15,
		1, 0, 0, 0, 175, 180, 5, 33, 0, 0, 176, 177, 5, 8, 0, 0, 177, 179, 5, 33,
		0, 0, 178, 176, 1, 0, 0, 0, 179, 182, 1, 0, 0, 0, 180, 178, 1, 0, 0, 0,
		180, 181, 1, 0, 0, 0, 181, 17, 1, 0, 0, 0, 182, 180, 1, 0, 0, 0, 183, 184,
		5, 9, 0, 0, 184, 185, 5, 6, 0, 0, 185, 186, 3, 20, 10, 0, 186, 187, 5,
		7, 0, 0, 187, 19, 1, 0, 0, 0, 188, 193, 3, 22, 11, 0, 189, 190, 5, 37,
		0, 0, 190, 192, 3, 22, 11, 0, 191, 189, 1, 0, 0, 0, 192, 195, 1, 0, 0,
		0, 193, 191, 1, 0, 0, 0, 193, 194, 1, 0, 0, 0, 194, 21, 1, 0, 0, 0, 195,
		193, 1, 0, 0, 0, 196, 198, 5, 10, 0, 0, 197, 196, 1, 0, 0, 0, 197, 198,
		1, 0, 0, 0, 198, 199, 1, 0, 0, 0, 199, 201, 5, 33, 0, 0, 200, 202, 3, 24,
		12, 0, 201, 200, 1, 0, 0, 0, 201, 202, 1, 0, 0, 0, 202, 203, 1, 0, 0, 0,
		203, 204, 3, 28, 14, 0, 204, 23, 1, 0, 0, 0, 205, 206, 5, 11, 0, 0, 206,
		214, 3, 26, 13, 0, 207, 209, 5, 12, 0, 0, 208, 210, 5, 37, 0, 0, 209, 208,
		1, 0, 0, 0, 209, 210, 1, 0, 0, 0, 210, 211, 1, 0, 0, 0, 211, 213, 3, 26,
		13, 0, 212, 207, 1, 0, 0, 0, 213, 216, 1, 0, 0, 0, 214, 212, 1, 0, 0, 0,
		214, 215, 1, 0, 0, 0, 215, 217, 1, 0, 0, 0, 216, 214, 1, 0, 0, 0, 217,
		218, 5, 13, 0, 0, 218, 25, 1, 0, 0, 0, 219, 221, 5, 33, 0, 0, 220, 222,
		3, 28, 14, 0, 221, 220, 1, 0, 0, 0, 221, 222, 1, 0, 0, 0, 222, 27, 1, 0,
		0, 0, 223, 227, 3, 30, 15, 0, 224, 227, 3, 34, 17, 0, 225, 227, 3, 44,
		22, 0, 226, 223, 1, 0, 0, 0, 226, 224, 1, 0, 0, 0, 226, 225, 1, 0, 0, 0,
		227, 29, 1, 0, 0, 0, 228, 230, 5, 33, 0, 0, 229, 231, 3, 32, 16, 0, 230,
		229, 1, 0, 0, 0, 230, 231, 1, 0, 0, 0, 231, 31, 1, 0, 0, 0, 232, 233, 5,
		11, 0, 0, 233, 238, 3, 28, 14, 0, 234, 235, 5, 12, 0, 0, 235, 237, 3, 28,
		14, 0, 236, 234, 1, 0, 0, 0, 237, 240, 1, 0, 0, 0, 238, 236, 1, 0, 0, 0,
		238, 239, 1, 0, 0, 0, 239, 241, 1, 0, 0, 0, 240, 238, 1, 0, 0, 0, 241,
		242, 5, 13, 0, 0, 242, 33, 1, 0, 0, 0, 243, 247, 3, 36, 18, 0, 244, 247,
		3, 38, 19, 0, 245, 247, 3, 46, 23, 0, 246, 243, 1, 0, 0, 0, 246, 244, 1,
		0, 0, 0, 246, 245, 1, 0, 0, 0, 247, 35, 1, 0, 0, 0, 248, 249, 5, 14, 0,
		0, 249, 250, 5, 34, 0, 0, 250, 251, 5, 15, 0, 0, 251, 252, 3, 28, 14, 0,
		252, 37, 1, 0, 0, 0, 253, 255, 5, 6, 0, 0, 254, 256, 3, 40, 20, 0, 255,
		254, 1, 0, 0, 0, 255, 256, 1, 0, 0, 0, 256, 257, 1, 0, 0, 0, 257, 258,
		5, 7, 0, 0, 258, 39, 1, 0, 0, 0, 259, 264, 3, 42, 21, 0, 260, 261, 5, 37,
		0, 0, 261, 263, 3, 42, 21, 0, 262, 260, 1, 0, 0, 0, 263, 266, 1, 0, 0,
		0, 264, 262, 1, 0, 0, 0, 264, 265, 1, 0, 0, 0, 265, 41, 1, 0, 0, 0, 266,
		264, 1, 0, 0, 0, 267, 268, 5, 33, 0, 0, 268, 269, 3, 28, 14, 0, 269, 43,
		1, 0, 0, 0, 270, 273, 3, 48, 24, 0, 271, 272, 5, 16, 0, 0, 272, 274, 3,
		48, 24, 0, 273, 271, 1, 0, 0, 0, 274, 275, 1, 0, 0, 0, 275, 273, 1, 0,
		0, 0, 275, 276, 1, 0, 0, 0, 276, 45, 1, 0, 0, 0, 277, 278, 5, 6, 0, 0,
		278, 47, 1, 0, 0, 0, 279, 282, 3, 30, 15, 0, 280, 282, 3, 34, 17, 0, 281,
		279, 1, 0, 0, 0, 281, 280, 1, 0, 0, 0, 282, 49, 1, 0, 0, 0, 283, 284, 5,
		17, 0, 0, 284, 285, 5, 6, 0, 0, 285, 286, 3, 52, 26, 0, 286, 287, 5, 7,
		0, 0, 287, 51, 1, 0, 0, 0, 288, 293, 3, 54, 27, 0, 289, 290, 5, 37, 0,
		0, 290, 292, 3, 54, 27, 0, 291, 289, 1, 0, 0, 0, 292, 295, 1, 0, 0, 0,
		293, 291, 1, 0, 0, 0, 293, 294, 1, 0, 0, 0, 294, 53, 1, 0, 0, 0, 295, 293,
		1, 0, 0, 0, 296, 298, 5, 10, 0, 0, 297, 296, 1, 0, 0, 0, 297, 298, 1, 0,
		0, 0, 298, 299, 1, 0, 0, 0, 299, 300, 5, 33, 0, 0, 300, 301, 3, 24, 12,
		0, 301, 302, 3, 56, 28, 0, 302, 303, 3, 56, 28, 0, 303, 55, 1, 0, 0, 0,
		304, 305, 5, 18, 0, 0, 305, 306, 3, 58, 29, 0, 306, 307, 5, 19, 0, 0, 307,
		57, 1, 0, 0, 0, 308, 316, 3, 60, 30, 0, 309, 311, 5, 12, 0, 0, 310, 312,
		5, 37, 0, 0, 311, 310, 1, 0, 0, 0, 311, 312, 1, 0, 0, 0, 312, 313, 1, 0,
		0, 0, 313, 315, 3, 60, 30, 0, 314, 309, 1, 0, 0, 0, 315, 318, 1, 0, 0,
		0, 316, 314, 1, 0, 0, 0, 316, 317, 1, 0, 0, 0, 317, 59, 1, 0, 0, 0, 318,
		316, 1, 0, 0, 0, 319, 320, 5, 33, 0, 0, 320, 321, 3, 28, 14, 0, 321, 61,
		1, 0, 0, 0, 322, 323, 5, 20, 0, 0, 323, 324, 5, 6, 0, 0, 324, 325, 3, 64,
		32, 0, 325, 326, 5, 7, 0, 0, 326, 327, 5, 37, 0, 0, 327, 63, 1, 0, 0, 0,
		328, 333, 3, 66, 33, 0, 329, 330, 5, 37, 0, 0, 330, 332, 3, 66, 33, 0,
		331, 329, 1, 0, 0, 0, 332, 335, 1, 0, 0, 0, 333, 331, 1, 0, 0, 0, 333,
		334, 1, 0, 0, 0, 334, 65, 1, 0, 0, 0, 335, 333, 1, 0, 0, 0, 336, 338, 5,
		10, 0, 0, 337, 336, 1, 0, 0, 0, 337, 338, 1, 0, 0, 0, 338, 339, 1, 0, 0,
		0, 339, 340, 5, 33, 0, 0, 340, 341, 3, 28, 14, 0, 341, 342, 5, 21, 0, 0,
		342, 343, 3, 68, 34, 0, 343, 67, 1, 0, 0, 0, 344, 353, 5, 22, 0, 0, 345,
		353, 5, 23, 0, 0, 346, 353, 5, 34, 0, 0, 347, 353, 5, 35, 0, 0, 348, 353,
		5, 36, 0, 0, 349, 353, 3, 70, 35, 0, 350, 353, 3, 74, 37, 0, 351, 353,
		5, 24, 0, 0, 352, 344, 1, 0, 0, 0, 352, 345, 1, 0, 0, 0, 352, 346, 1, 0,
		0, 0, 352, 347, 1, 0, 0, 0, 352, 348, 1, 0, 0, 0, 352, 349, 1, 0, 0, 0,
		352, 350, 1, 0, 0, 0, 352, 351, 1, 0, 0, 0, 353, 69, 1, 0, 0, 0, 354, 355,
		5, 14, 0, 0, 355, 356, 3, 72, 36, 0, 356, 357, 5, 15, 0, 0, 357, 71, 1,
		0, 0, 0, 358, 371, 3, 68, 34, 0, 359, 367, 3, 68, 34, 0, 360, 362, 5, 12,
		0, 0, 361, 363, 5, 37, 0, 0, 362, 361, 1, 0, 0, 0, 362, 363, 1, 0, 0, 0,
		363, 364, 1, 0, 0, 0, 364, 366, 3, 68, 34, 0, 365, 360, 1, 0, 0, 0, 366,
		369, 1, 0, 0, 0, 367, 365, 1, 0, 0, 0, 367, 368, 1, 0, 0, 0, 368, 371,
		1, 0, 0, 0, 369, 367, 1, 0, 0, 0, 370, 358, 1, 0, 0, 0, 370, 359, 1, 0,
		0, 0, 371, 73, 1, 0, 0, 0, 372, 373, 5, 6, 0, 0, 373, 374, 3, 76, 38, 0,
		374, 375, 5, 7, 0, 0, 375, 75, 1, 0, 0, 0, 376, 384, 3, 78, 39, 0, 377,
		379, 5, 12, 0, 0, 378, 380, 5, 37, 0, 0, 379, 378, 1, 0, 0, 0, 379, 380,
		1, 0, 0, 0, 380, 381, 1, 0, 0, 0, 381, 383, 3, 78, 39, 0, 382, 377, 1,
		0, 0, 0, 383, 386, 1, 0, 0, 0, 384, 382, 1, 0, 0, 0, 384, 385, 1, 0, 0,
		0, 385, 77, 1, 0, 0, 0, 386, 384, 1, 0, 0, 0, 387, 388, 5, 33, 0, 0, 388,
		389, 5, 25, 0, 0, 389, 390, 3, 68, 34, 0, 390, 79, 1, 0, 0, 0, 391, 392,
		5, 26, 0, 0, 392, 393, 5, 6, 0, 0, 393, 394, 3, 82, 41, 0, 394, 395, 5,
		7, 0, 0, 395, 396, 5, 37, 0, 0, 396, 81, 1, 0, 0, 0, 397, 402, 3, 84, 42,
		0, 398, 399, 5, 37, 0, 0, 399, 401, 3, 84, 42, 0, 400, 398, 1, 0, 0, 0,
		401, 404, 1, 0, 0, 0, 402, 400, 1, 0, 0, 0, 402, 403, 1, 0, 0, 0, 403,
		83, 1, 0, 0, 0, 404, 402, 1, 0, 0, 0, 405, 407, 5, 10, 0, 0, 406, 405,
		1, 0, 0, 0, 406, 407, 1, 0, 0, 0, 407, 408, 1, 0, 0, 0, 408, 409, 3, 54,
		27, 0, 409, 410, 3, 86, 43, 0, 410, 85, 1, 0, 0, 0, 411, 412, 5, 6, 0,
		0, 412, 417, 3, 88, 44, 0, 413, 414, 3, 106, 53, 0, 414, 415, 5, 7, 0,
		0, 415, 417, 1, 0, 0, 0, 416, 411, 1, 0, 0, 0, 416, 413, 1, 0, 0, 0, 417,
		87, 1, 0, 0, 0, 418, 419, 5, 27, 0, 0, 419, 420, 5, 6, 0, 0, 420, 421,
		3, 90, 45, 0, 421, 422, 5, 7, 0, 0, 422, 89, 1, 0, 0, 0, 423, 426, 3, 92,
		46, 0, 424, 426, 3, 94, 47, 0, 425, 423, 1, 0, 0, 0, 425, 424, 1, 0, 0,
		0, 426, 91, 1, 0, 0, 0, 427, 428, 5, 33, 0, 0, 428, 429, 3, 30, 15, 0,
		429, 93, 1, 0, 0, 0, 430, 431, 5, 33, 0, 0, 431, 432, 5, 21, 0, 0, 432,
		433, 3, 96, 48, 0, 433, 95, 1, 0, 0, 0, 434, 435, 3, 98, 49, 0, 435, 436,
		3, 100, 50, 0, 436, 437, 3, 32, 16, 0, 437, 97, 1, 0, 0, 0, 438, 443, 5,
		33, 0, 0, 439, 440, 5, 28, 0, 0, 440, 442, 5, 33, 0, 0, 441, 439, 1, 0,
		0, 0, 442, 445, 1, 0, 0, 0, 443, 441, 1, 0, 0, 0, 443, 444, 1, 0, 0, 0,
		444, 99, 1, 0, 0, 0, 445, 443, 1, 0, 0, 0, 446, 447, 5, 18, 0, 0, 447,
		448, 3, 102, 51, 0, 448, 449, 5, 19, 0, 0, 449, 101, 1, 0, 0, 0, 450, 451,
		3, 104, 52, 0, 451, 453, 5, 12, 0, 0, 452, 454, 5, 37, 0, 0, 453, 452,
		1, 0, 0, 0, 453, 454, 1, 0, 0, 0, 454, 455, 1, 0, 0, 0, 455, 456, 3, 104,
		52, 0, 456, 103, 1, 0, 0, 0, 457, 458, 3, 96, 48, 0, 458, 105, 1, 0, 0,
		0, 459, 460, 5, 29, 0, 0, 460, 461, 5, 6, 0, 0, 461, 462, 3, 108, 54, 0,
		462, 463, 5, 7, 0, 0, 463, 107, 1, 0, 0, 0, 464, 469, 3, 110, 55, 0, 465,
		466, 5, 37, 0, 0, 466, 468, 3, 110, 55, 0, 467, 465, 1, 0, 0, 0, 468, 471,
		1, 0, 0, 0, 469, 467, 1, 0, 0, 0, 469, 470, 1, 0, 0, 0, 470, 109, 1, 0,
		0, 0, 471, 469, 1, 0, 0, 0, 472, 473, 3, 112, 56, 0, 473, 474, 5, 30, 0,
		0, 474, 475, 3, 116, 58, 0, 475, 111, 1, 0, 0, 0, 476, 478, 5, 33, 0, 0,
		477, 476, 1, 0, 0, 0, 477, 478, 1, 0, 0, 0, 478, 479, 1, 0, 0, 0, 479,
		487, 3, 114, 57, 0, 480, 484, 5, 33, 0, 0, 481, 482, 5, 14, 0, 0, 482,
		483, 5, 34, 0, 0, 483, 485, 5, 15, 0, 0, 484, 481, 1, 0, 0, 0, 484, 485,
		1, 0, 0, 0, 485, 487, 1, 0, 0, 0, 486, 477, 1, 0, 0, 0, 486, 480, 1, 0,
		0, 0, 487, 113, 1, 0, 0, 0, 488, 489, 7, 1, 0, 0, 489, 115, 1, 0, 0, 0,
		490, 493, 3, 112, 56, 0, 491, 493, 3, 118, 59, 0, 492, 490, 1, 0, 0, 0,
		492, 491, 1, 0, 0, 0, 493, 117, 1, 0, 0, 0, 494, 495, 5, 6, 0, 0, 495,
		500, 3, 112, 56, 0, 496, 497, 5, 37, 0, 0, 497, 499, 3, 112, 56, 0, 498,
		496, 1, 0, 0, 0, 499, 502, 1, 0, 0, 0, 500, 498, 1, 0, 0, 0, 500, 501,
		1, 0, 0, 0, 501, 503, 1, 0, 0, 0, 502, 500, 1, 0, 0, 0, 503, 504, 5, 7,
		0, 0, 504, 119, 1, 0, 0, 0, 47, 122, 124, 131, 137, 146, 155, 167, 171,
		180, 193, 197, 201, 209, 214, 221, 226, 230, 238, 246, 255, 264, 275, 281,
		293, 297, 311, 316, 333, 337, 352, 362, 367, 370, 379, 384, 402, 406, 416,
		425, 443, 453, 469, 477, 484, 486, 492, 500,
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

// nevaParserInit initializes any static state used to implement nevaParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewnevaParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func NevaParserInit() {
	staticData := &NevaParserStaticData
	staticData.once.Do(nevaParserInit)
}

// NewnevaParser produces a new parser instance for the optional input antlr.TokenStream.
func NewnevaParser(input antlr.TokenStream) *nevaParser {
	NevaParserInit()
	this := new(nevaParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &NevaParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "neva.g4"

	return this
}

// nevaParser tokens.
const (
	nevaParserEOF        = antlr.TokenEOF
	nevaParserT__0       = 1
	nevaParserT__1       = 2
	nevaParserT__2       = 3
	nevaParserT__3       = 4
	nevaParserT__4       = 5
	nevaParserT__5       = 6
	nevaParserT__6       = 7
	nevaParserT__7       = 8
	nevaParserT__8       = 9
	nevaParserT__9       = 10
	nevaParserT__10      = 11
	nevaParserT__11      = 12
	nevaParserT__12      = 13
	nevaParserT__13      = 14
	nevaParserT__14      = 15
	nevaParserT__15      = 16
	nevaParserT__16      = 17
	nevaParserT__17      = 18
	nevaParserT__18      = 19
	nevaParserT__19      = 20
	nevaParserT__20      = 21
	nevaParserT__21      = 22
	nevaParserT__22      = 23
	nevaParserT__23      = 24
	nevaParserT__24      = 25
	nevaParserT__25      = 26
	nevaParserT__26      = 27
	nevaParserT__27      = 28
	nevaParserT__28      = 29
	nevaParserT__29      = 30
	nevaParserT__30      = 31
	nevaParserT__31      = 32
	nevaParserIDENTIFIER = 33
	nevaParserINT        = 34
	nevaParserFLOAT      = 35
	nevaParserSTRING     = 36
	nevaParserNEWLINE    = 37
	nevaParserWS         = 38
)

// nevaParser rules.
const (
	nevaParserRULE_prog              = 0
	nevaParserRULE_comment           = 1
	nevaParserRULE_singleLineComment = 2
	nevaParserRULE_multiLineComment  = 3
	nevaParserRULE_stmt              = 4
	nevaParserRULE_useStmt           = 5
	nevaParserRULE_importList        = 6
	nevaParserRULE_importDef         = 7
	nevaParserRULE_importPath        = 8
	nevaParserRULE_typeStmt          = 9
	nevaParserRULE_typeDefList       = 10
	nevaParserRULE_typeDef           = 11
	nevaParserRULE_typeParams        = 12
	nevaParserRULE_typeParam         = 13
	nevaParserRULE_typeExpr          = 14
	nevaParserRULE_typeInstExpr      = 15
	nevaParserRULE_typeArgs          = 16
	nevaParserRULE_typeLitExpr       = 17
	nevaParserRULE_arrTypeExpr       = 18
	nevaParserRULE_recTypeExpr       = 19
	nevaParserRULE_recTypeFields     = 20
	nevaParserRULE_recTypeField      = 21
	nevaParserRULE_unionTypeExpr     = 22
	nevaParserRULE_enumTypeExpr      = 23
	nevaParserRULE_nonUnionTypeExpr  = 24
	nevaParserRULE_ioStmt            = 25
	nevaParserRULE_interfaceDefList  = 26
	nevaParserRULE_interfaceDef      = 27
	nevaParserRULE_portsDef          = 28
	nevaParserRULE_portDefList       = 29
	nevaParserRULE_portDef           = 30
	nevaParserRULE_constStmt         = 31
	nevaParserRULE_constDefList      = 32
	nevaParserRULE_constDef          = 33
	nevaParserRULE_constValue        = 34
	nevaParserRULE_arrLit            = 35
	nevaParserRULE_arrItems          = 36
	nevaParserRULE_recLit            = 37
	nevaParserRULE_recValueFields    = 38
	nevaParserRULE_recValueField     = 39
	nevaParserRULE_compStmt          = 40
	nevaParserRULE_compDefList       = 41
	nevaParserRULE_compDef           = 42
	nevaParserRULE_compBody          = 43
	nevaParserRULE_compNodesDef      = 44
	nevaParserRULE_compNodeDefList   = 45
	nevaParserRULE_absNodeDef        = 46
	nevaParserRULE_concreteNodeDef   = 47
	nevaParserRULE_concreteNodeInst  = 48
	nevaParserRULE_nodeRef           = 49
	nevaParserRULE_nodeArgs          = 50
	nevaParserRULE_nodeArgList       = 51
	nevaParserRULE_nodeArg           = 52
	nevaParserRULE_compNetDef        = 53
	nevaParserRULE_connDefList       = 54
	nevaParserRULE_connDef           = 55
	nevaParserRULE_portAddr          = 56
	nevaParserRULE_portDirection     = 57
	nevaParserRULE_connReceiverSide  = 58
	nevaParserRULE_connReceivers     = 59
)

// IProgContext is an interface to support dynamic dispatch.
type IProgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	AllComment() []ICommentContext
	Comment(i int) ICommentContext
	AllStmt() []IStmtContext
	Stmt(i int) IStmtContext

	// IsProgContext differentiates from other interfaces.
	IsProgContext()
}

type ProgContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgContext() *ProgContext {
	var p = new(ProgContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_prog
	return p
}

func InitEmptyProgContext(p *ProgContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_prog
}

func (*ProgContext) IsProgContext() {}

func NewProgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgContext {
	var p = new(ProgContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_prog

	return p
}

func (s *ProgContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgContext) EOF() antlr.TerminalNode {
	return s.GetToken(nevaParserEOF, 0)
}

func (s *ProgContext) AllComment() []ICommentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICommentContext); ok {
			len++
		}
	}

	tst := make([]ICommentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICommentContext); ok {
			tst[i] = t.(ICommentContext)
			i++
		}
	}

	return tst
}

func (s *ProgContext) Comment(i int) ICommentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICommentContext)
}

func (s *ProgContext) AllStmt() []IStmtContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStmtContext); ok {
			len++
		}
	}

	tst := make([]IStmtContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStmtContext); ok {
			tst[i] = t.(IStmtContext)
			i++
		}
	}

	return tst
}

func (s *ProgContext) Stmt(i int) IStmtContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtContext)
}

func (s *ProgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterProg(s)
	}
}

func (s *ProgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitProg(s)
	}
}

func (p *nevaParser) Prog() (localctx IProgContext) {
	localctx = NewProgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, nevaParserRULE_prog)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(124)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1180202) != 0 {
		p.SetState(122)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case nevaParserT__0, nevaParserT__2:
			{
				p.SetState(120)
				p.Comment()
			}

		case nevaParserT__4, nevaParserT__8, nevaParserT__16, nevaParserT__19:
			{
				p.SetState(121)
				p.Stmt()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(126)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(127)
		p.Match(nevaParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICommentContext is an interface to support dynamic dispatch.
type ICommentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SingleLineComment() ISingleLineCommentContext
	MultiLineComment() IMultiLineCommentContext

	// IsCommentContext differentiates from other interfaces.
	IsCommentContext()
}

type CommentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommentContext() *CommentContext {
	var p = new(CommentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_comment
	return p
}

func InitEmptyCommentContext(p *CommentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_comment
}

func (*CommentContext) IsCommentContext() {}

func NewCommentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommentContext {
	var p = new(CommentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_comment

	return p
}

func (s *CommentContext) GetParser() antlr.Parser { return s.parser }

func (s *CommentContext) SingleLineComment() ISingleLineCommentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISingleLineCommentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISingleLineCommentContext)
}

func (s *CommentContext) MultiLineComment() IMultiLineCommentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMultiLineCommentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMultiLineCommentContext)
}

func (s *CommentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterComment(s)
	}
}

func (s *CommentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitComment(s)
	}
}

func (p *nevaParser) Comment() (localctx ICommentContext) {
	localctx = NewCommentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, nevaParserRULE_comment)
	p.SetState(131)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__0:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(129)
			p.SingleLineComment()
		}

	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(130)
			p.MultiLineComment()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISingleLineCommentContext is an interface to support dynamic dispatch.
type ISingleLineCommentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NEWLINE() antlr.TerminalNode

	// IsSingleLineCommentContext differentiates from other interfaces.
	IsSingleLineCommentContext()
}

type SingleLineCommentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySingleLineCommentContext() *SingleLineCommentContext {
	var p = new(SingleLineCommentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_singleLineComment
	return p
}

func InitEmptySingleLineCommentContext(p *SingleLineCommentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_singleLineComment
}

func (*SingleLineCommentContext) IsSingleLineCommentContext() {}

func NewSingleLineCommentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SingleLineCommentContext {
	var p = new(SingleLineCommentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_singleLineComment

	return p
}

func (s *SingleLineCommentContext) GetParser() antlr.Parser { return s.parser }

func (s *SingleLineCommentContext) NEWLINE() antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, 0)
}

func (s *SingleLineCommentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SingleLineCommentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SingleLineCommentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterSingleLineComment(s)
	}
}

func (s *SingleLineCommentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitSingleLineComment(s)
	}
}

func (p *nevaParser) SingleLineComment() (localctx ISingleLineCommentContext) {
	localctx = NewSingleLineCommentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, nevaParserRULE_singleLineComment)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(133)
		p.Match(nevaParserT__0)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(134)
				_la = p.GetTokenStream().LA(1)

				if _la <= 0 || _la == nevaParserT__1 {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

		}
		p.SetState(139)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	{
		p.SetState(140)
		p.Match(nevaParserNEWLINE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMultiLineCommentContext is an interface to support dynamic dispatch.
type IMultiLineCommentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsMultiLineCommentContext differentiates from other interfaces.
	IsMultiLineCommentContext()
}

type MultiLineCommentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMultiLineCommentContext() *MultiLineCommentContext {
	var p = new(MultiLineCommentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_multiLineComment
	return p
}

func InitEmptyMultiLineCommentContext(p *MultiLineCommentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_multiLineComment
}

func (*MultiLineCommentContext) IsMultiLineCommentContext() {}

func NewMultiLineCommentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MultiLineCommentContext {
	var p = new(MultiLineCommentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_multiLineComment

	return p
}

func (s *MultiLineCommentContext) GetParser() antlr.Parser { return s.parser }
func (s *MultiLineCommentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiLineCommentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MultiLineCommentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterMultiLineComment(s)
	}
}

func (s *MultiLineCommentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitMultiLineComment(s)
	}
}

func (p *nevaParser) MultiLineComment() (localctx IMultiLineCommentContext) {
	localctx = NewMultiLineCommentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, nevaParserRULE_multiLineComment)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(142)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(146)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 1 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1+1 {
			p.SetState(143)
			p.MatchWildcard()

		}
		p.SetState(148)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	{
		p.SetState(149)
		p.Match(nevaParserT__3)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStmtContext is an interface to support dynamic dispatch.
type IStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UseStmt() IUseStmtContext
	TypeStmt() ITypeStmtContext
	IoStmt() IIoStmtContext
	ConstStmt() IConstStmtContext

	// IsStmtContext differentiates from other interfaces.
	IsStmtContext()
}

type StmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStmtContext() *StmtContext {
	var p = new(StmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_stmt
	return p
}

func InitEmptyStmtContext(p *StmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_stmt
}

func (*StmtContext) IsStmtContext() {}

func NewStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtContext {
	var p = new(StmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_stmt

	return p
}

func (s *StmtContext) GetParser() antlr.Parser { return s.parser }

func (s *StmtContext) UseStmt() IUseStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUseStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUseStmtContext)
}

func (s *StmtContext) TypeStmt() ITypeStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeStmtContext)
}

func (s *StmtContext) IoStmt() IIoStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIoStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIoStmtContext)
}

func (s *StmtContext) ConstStmt() IConstStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstStmtContext)
}

func (s *StmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterStmt(s)
	}
}

func (s *StmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitStmt(s)
	}
}

func (p *nevaParser) Stmt() (localctx IStmtContext) {
	localctx = NewStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, nevaParserRULE_stmt)
	p.SetState(155)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__4:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(151)
			p.UseStmt()
		}

	case nevaParserT__8:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(152)
			p.TypeStmt()
		}

	case nevaParserT__16:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(153)
			p.IoStmt()
		}

	case nevaParserT__19:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(154)
			p.ConstStmt()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUseStmtContext is an interface to support dynamic dispatch.
type IUseStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ImportList() IImportListContext

	// IsUseStmtContext differentiates from other interfaces.
	IsUseStmtContext()
}

type UseStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUseStmtContext() *UseStmtContext {
	var p = new(UseStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_useStmt
	return p
}

func InitEmptyUseStmtContext(p *UseStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_useStmt
}

func (*UseStmtContext) IsUseStmtContext() {}

func NewUseStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UseStmtContext {
	var p = new(UseStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_useStmt

	return p
}

func (s *UseStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *UseStmtContext) ImportList() IImportListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImportListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImportListContext)
}

func (s *UseStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UseStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UseStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterUseStmt(s)
	}
}

func (s *UseStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitUseStmt(s)
	}
}

func (p *nevaParser) UseStmt() (localctx IUseStmtContext) {
	localctx = NewUseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, nevaParserRULE_useStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(157)
		p.Match(nevaParserT__4)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(158)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(159)
		p.ImportList()
	}
	{
		p.SetState(160)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IImportListContext is an interface to support dynamic dispatch.
type IImportListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllImportDef() []IImportDefContext
	ImportDef(i int) IImportDefContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsImportListContext differentiates from other interfaces.
	IsImportListContext()
}

type ImportListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportListContext() *ImportListContext {
	var p = new(ImportListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_importList
	return p
}

func InitEmptyImportListContext(p *ImportListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_importList
}

func (*ImportListContext) IsImportListContext() {}

func NewImportListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportListContext {
	var p = new(ImportListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_importList

	return p
}

func (s *ImportListContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportListContext) AllImportDef() []IImportDefContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IImportDefContext); ok {
			len++
		}
	}

	tst := make([]IImportDefContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IImportDefContext); ok {
			tst[i] = t.(IImportDefContext)
			i++
		}
	}

	return tst
}

func (s *ImportListContext) ImportDef(i int) IImportDefContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImportDefContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImportDefContext)
}

func (s *ImportListContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ImportListContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *ImportListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterImportList(s)
	}
}

func (s *ImportListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitImportList(s)
	}
}

func (p *nevaParser) ImportList() (localctx IImportListContext) {
	localctx = NewImportListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, nevaParserRULE_importList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(162)
		p.ImportDef()
	}
	p.SetState(167)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(163)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(164)
			p.ImportDef()
		}

		p.SetState(169)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IImportDefContext is an interface to support dynamic dispatch.
type IImportDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ImportPath() IImportPathContext
	IDENTIFIER() antlr.TerminalNode

	// IsImportDefContext differentiates from other interfaces.
	IsImportDefContext()
}

type ImportDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportDefContext() *ImportDefContext {
	var p = new(ImportDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_importDef
	return p
}

func InitEmptyImportDefContext(p *ImportDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_importDef
}

func (*ImportDefContext) IsImportDefContext() {}

func NewImportDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportDefContext {
	var p = new(ImportDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_importDef

	return p
}

func (s *ImportDefContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportDefContext) ImportPath() IImportPathContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImportPathContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImportPathContext)
}

func (s *ImportDefContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *ImportDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterImportDef(s)
	}
}

func (s *ImportDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitImportDef(s)
	}
}

func (p *nevaParser) ImportDef() (localctx IImportDefContext) {
	localctx = NewImportDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, nevaParserRULE_importDef)
	p.EnterOuterAlt(localctx, 1)
	p.SetState(171)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(170)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(173)
		p.ImportPath()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IImportPathContext is an interface to support dynamic dispatch.
type IImportPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode

	// IsImportPathContext differentiates from other interfaces.
	IsImportPathContext()
}

type ImportPathContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportPathContext() *ImportPathContext {
	var p = new(ImportPathContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_importPath
	return p
}

func InitEmptyImportPathContext(p *ImportPathContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_importPath
}

func (*ImportPathContext) IsImportPathContext() {}

func NewImportPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportPathContext {
	var p = new(ImportPathContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_importPath

	return p
}

func (s *ImportPathContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportPathContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(nevaParserIDENTIFIER)
}

func (s *ImportPathContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, i)
}

func (s *ImportPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterImportPath(s)
	}
}

func (s *ImportPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitImportPath(s)
	}
}

func (p *nevaParser) ImportPath() (localctx IImportPathContext) {
	localctx = NewImportPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, nevaParserRULE_importPath)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(175)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(180)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 {
		{
			p.SetState(176)
			p.Match(nevaParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(177)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(182)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeStmtContext is an interface to support dynamic dispatch.
type ITypeStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TypeDefList() ITypeDefListContext

	// IsTypeStmtContext differentiates from other interfaces.
	IsTypeStmtContext()
}

type TypeStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeStmtContext() *TypeStmtContext {
	var p = new(TypeStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeStmt
	return p
}

func InitEmptyTypeStmtContext(p *TypeStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeStmt
}

func (*TypeStmtContext) IsTypeStmtContext() {}

func NewTypeStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeStmtContext {
	var p = new(TypeStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_typeStmt

	return p
}

func (s *TypeStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeStmtContext) TypeDefList() ITypeDefListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeDefListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeDefListContext)
}

func (s *TypeStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterTypeStmt(s)
	}
}

func (s *TypeStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitTypeStmt(s)
	}
}

func (p *nevaParser) TypeStmt() (localctx ITypeStmtContext) {
	localctx = NewTypeStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, nevaParserRULE_typeStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(183)
		p.Match(nevaParserT__8)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(184)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(185)
		p.TypeDefList()
	}
	{
		p.SetState(186)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeDefListContext is an interface to support dynamic dispatch.
type ITypeDefListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTypeDef() []ITypeDefContext
	TypeDef(i int) ITypeDefContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsTypeDefListContext differentiates from other interfaces.
	IsTypeDefListContext()
}

type TypeDefListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeDefListContext() *TypeDefListContext {
	var p = new(TypeDefListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeDefList
	return p
}

func InitEmptyTypeDefListContext(p *TypeDefListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeDefList
}

func (*TypeDefListContext) IsTypeDefListContext() {}

func NewTypeDefListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeDefListContext {
	var p = new(TypeDefListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_typeDefList

	return p
}

func (s *TypeDefListContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeDefListContext) AllTypeDef() []ITypeDefContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITypeDefContext); ok {
			len++
		}
	}

	tst := make([]ITypeDefContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITypeDefContext); ok {
			tst[i] = t.(ITypeDefContext)
			i++
		}
	}

	return tst
}

func (s *TypeDefListContext) TypeDef(i int) ITypeDefContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeDefContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeDefContext)
}

func (s *TypeDefListContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *TypeDefListContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *TypeDefListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeDefListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeDefListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterTypeDefList(s)
	}
}

func (s *TypeDefListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitTypeDefList(s)
	}
}

func (p *nevaParser) TypeDefList() (localctx ITypeDefListContext) {
	localctx = NewTypeDefListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, nevaParserRULE_typeDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(188)
		p.TypeDef()
	}
	p.SetState(193)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(189)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(190)
			p.TypeDef()
		}

		p.SetState(195)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeDefContext is an interface to support dynamic dispatch.
type ITypeDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeExpr() ITypeExprContext
	TypeParams() ITypeParamsContext

	// IsTypeDefContext differentiates from other interfaces.
	IsTypeDefContext()
}

type TypeDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeDefContext() *TypeDefContext {
	var p = new(TypeDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeDef
	return p
}

func InitEmptyTypeDefContext(p *TypeDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeDef
}

func (*TypeDefContext) IsTypeDefContext() {}

func NewTypeDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeDefContext {
	var p = new(TypeDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_typeDef

	return p
}

func (s *TypeDefContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeDefContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *TypeDefContext) TypeExpr() ITypeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeExprContext)
}

func (s *TypeDefContext) TypeParams() ITypeParamsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeParamsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeParamsContext)
}

func (s *TypeDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterTypeDef(s)
	}
}

func (s *TypeDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitTypeDef(s)
	}
}

func (p *nevaParser) TypeDef() (localctx ITypeDefContext) {
	localctx = NewTypeDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, nevaParserRULE_typeDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(197)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__9 {
		{
			p.SetState(196)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(199)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(201)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__10 {
		{
			p.SetState(200)
			p.TypeParams()
		}

	}
	{
		p.SetState(203)
		p.TypeExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeParamsContext is an interface to support dynamic dispatch.
type ITypeParamsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTypeParam() []ITypeParamContext
	TypeParam(i int) ITypeParamContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsTypeParamsContext differentiates from other interfaces.
	IsTypeParamsContext()
}

type TypeParamsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeParamsContext() *TypeParamsContext {
	var p = new(TypeParamsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeParams
	return p
}

func InitEmptyTypeParamsContext(p *TypeParamsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeParams
}

func (*TypeParamsContext) IsTypeParamsContext() {}

func NewTypeParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeParamsContext {
	var p = new(TypeParamsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_typeParams

	return p
}

func (s *TypeParamsContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeParamsContext) AllTypeParam() []ITypeParamContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITypeParamContext); ok {
			len++
		}
	}

	tst := make([]ITypeParamContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITypeParamContext); ok {
			tst[i] = t.(ITypeParamContext)
			i++
		}
	}

	return tst
}

func (s *TypeParamsContext) TypeParam(i int) ITypeParamContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeParamContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeParamContext)
}

func (s *TypeParamsContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *TypeParamsContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *TypeParamsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeParamsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeParamsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterTypeParams(s)
	}
}

func (s *TypeParamsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitTypeParams(s)
	}
}

func (p *nevaParser) TypeParams() (localctx ITypeParamsContext) {
	localctx = NewTypeParamsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, nevaParserRULE_typeParams)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(205)
		p.Match(nevaParserT__10)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(206)
		p.TypeParam()
	}
	p.SetState(214)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__11 {
		{
			p.SetState(207)
			p.Match(nevaParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(209)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserNEWLINE {
			{
				p.SetState(208)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(211)
			p.TypeParam()
		}

		p.SetState(216)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(217)
		p.Match(nevaParserT__12)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeParamContext is an interface to support dynamic dispatch.
type ITypeParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeExpr() ITypeExprContext

	// IsTypeParamContext differentiates from other interfaces.
	IsTypeParamContext()
}

type TypeParamContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeParamContext() *TypeParamContext {
	var p = new(TypeParamContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeParam
	return p
}

func InitEmptyTypeParamContext(p *TypeParamContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeParam
}

func (*TypeParamContext) IsTypeParamContext() {}

func NewTypeParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeParamContext {
	var p = new(TypeParamContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_typeParam

	return p
}

func (s *TypeParamContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeParamContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *TypeParamContext) TypeExpr() ITypeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeExprContext)
}

func (s *TypeParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterTypeParam(s)
	}
}

func (s *TypeParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitTypeParam(s)
	}
}

func (p *nevaParser) TypeParam() (localctx ITypeParamContext) {
	localctx = NewTypeParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, nevaParserRULE_typeParam)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(219)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(221)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&8589951040) != 0 {
		{
			p.SetState(220)
			p.TypeExpr()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeExprContext is an interface to support dynamic dispatch.
type ITypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TypeInstExpr() ITypeInstExprContext
	TypeLitExpr() ITypeLitExprContext
	UnionTypeExpr() IUnionTypeExprContext

	// IsTypeExprContext differentiates from other interfaces.
	IsTypeExprContext()
}

type TypeExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeExprContext() *TypeExprContext {
	var p = new(TypeExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeExpr
	return p
}

func InitEmptyTypeExprContext(p *TypeExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeExpr
}

func (*TypeExprContext) IsTypeExprContext() {}

func NewTypeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeExprContext {
	var p = new(TypeExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_typeExpr

	return p
}

func (s *TypeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeExprContext) TypeInstExpr() ITypeInstExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeInstExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeInstExprContext)
}

func (s *TypeExprContext) TypeLitExpr() ITypeLitExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeLitExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeLitExprContext)
}

func (s *TypeExprContext) UnionTypeExpr() IUnionTypeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionTypeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionTypeExprContext)
}

func (s *TypeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterTypeExpr(s)
	}
}

func (s *TypeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitTypeExpr(s)
	}
}

func (p *nevaParser) TypeExpr() (localctx ITypeExprContext) {
	localctx = NewTypeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, nevaParserRULE_typeExpr)
	p.SetState(226)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(223)
			p.TypeInstExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(224)
			p.TypeLitExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(225)
			p.UnionTypeExpr()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeInstExprContext is an interface to support dynamic dispatch.
type ITypeInstExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeArgs() ITypeArgsContext

	// IsTypeInstExprContext differentiates from other interfaces.
	IsTypeInstExprContext()
}

type TypeInstExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeInstExprContext() *TypeInstExprContext {
	var p = new(TypeInstExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeInstExpr
	return p
}

func InitEmptyTypeInstExprContext(p *TypeInstExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeInstExpr
}

func (*TypeInstExprContext) IsTypeInstExprContext() {}

func NewTypeInstExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeInstExprContext {
	var p = new(TypeInstExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_typeInstExpr

	return p
}

func (s *TypeInstExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeInstExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *TypeInstExprContext) TypeArgs() ITypeArgsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeArgsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeArgsContext)
}

func (s *TypeInstExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeInstExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeInstExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterTypeInstExpr(s)
	}
}

func (s *TypeInstExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitTypeInstExpr(s)
	}
}

func (p *nevaParser) TypeInstExpr() (localctx ITypeInstExprContext) {
	localctx = NewTypeInstExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, nevaParserRULE_typeInstExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(228)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(230)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__10 {
		{
			p.SetState(229)
			p.TypeArgs()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeArgsContext is an interface to support dynamic dispatch.
type ITypeArgsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTypeExpr() []ITypeExprContext
	TypeExpr(i int) ITypeExprContext

	// IsTypeArgsContext differentiates from other interfaces.
	IsTypeArgsContext()
}

type TypeArgsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeArgsContext() *TypeArgsContext {
	var p = new(TypeArgsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeArgs
	return p
}

func InitEmptyTypeArgsContext(p *TypeArgsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeArgs
}

func (*TypeArgsContext) IsTypeArgsContext() {}

func NewTypeArgsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeArgsContext {
	var p = new(TypeArgsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_typeArgs

	return p
}

func (s *TypeArgsContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeArgsContext) AllTypeExpr() []ITypeExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITypeExprContext); ok {
			len++
		}
	}

	tst := make([]ITypeExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITypeExprContext); ok {
			tst[i] = t.(ITypeExprContext)
			i++
		}
	}

	return tst
}

func (s *TypeArgsContext) TypeExpr(i int) ITypeExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeExprContext)
}

func (s *TypeArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeArgsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterTypeArgs(s)
	}
}

func (s *TypeArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitTypeArgs(s)
	}
}

func (p *nevaParser) TypeArgs() (localctx ITypeArgsContext) {
	localctx = NewTypeArgsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, nevaParserRULE_typeArgs)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(232)
		p.Match(nevaParserT__10)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(233)
		p.TypeExpr()
	}
	p.SetState(238)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__11 {
		{
			p.SetState(234)
			p.Match(nevaParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(235)
			p.TypeExpr()
		}

		p.SetState(240)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(241)
		p.Match(nevaParserT__12)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeLitExprContext is an interface to support dynamic dispatch.
type ITypeLitExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ArrTypeExpr() IArrTypeExprContext
	RecTypeExpr() IRecTypeExprContext
	EnumTypeExpr() IEnumTypeExprContext

	// IsTypeLitExprContext differentiates from other interfaces.
	IsTypeLitExprContext()
}

type TypeLitExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeLitExprContext() *TypeLitExprContext {
	var p = new(TypeLitExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeLitExpr
	return p
}

func InitEmptyTypeLitExprContext(p *TypeLitExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeLitExpr
}

func (*TypeLitExprContext) IsTypeLitExprContext() {}

func NewTypeLitExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeLitExprContext {
	var p = new(TypeLitExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_typeLitExpr

	return p
}

func (s *TypeLitExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeLitExprContext) ArrTypeExpr() IArrTypeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrTypeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrTypeExprContext)
}

func (s *TypeLitExprContext) RecTypeExpr() IRecTypeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRecTypeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRecTypeExprContext)
}

func (s *TypeLitExprContext) EnumTypeExpr() IEnumTypeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnumTypeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnumTypeExprContext)
}

func (s *TypeLitExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeLitExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeLitExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterTypeLitExpr(s)
	}
}

func (s *TypeLitExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitTypeLitExpr(s)
	}
}

func (p *nevaParser) TypeLitExpr() (localctx ITypeLitExprContext) {
	localctx = NewTypeLitExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, nevaParserRULE_typeLitExpr)
	p.SetState(246)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 18, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(243)
			p.ArrTypeExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(244)
			p.RecTypeExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(245)
			p.EnumTypeExpr()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArrTypeExprContext is an interface to support dynamic dispatch.
type IArrTypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INT() antlr.TerminalNode
	TypeExpr() ITypeExprContext

	// IsArrTypeExprContext differentiates from other interfaces.
	IsArrTypeExprContext()
}

type ArrTypeExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrTypeExprContext() *ArrTypeExprContext {
	var p = new(ArrTypeExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_arrTypeExpr
	return p
}

func InitEmptyArrTypeExprContext(p *ArrTypeExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_arrTypeExpr
}

func (*ArrTypeExprContext) IsArrTypeExprContext() {}

func NewArrTypeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrTypeExprContext {
	var p = new(ArrTypeExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_arrTypeExpr

	return p
}

func (s *ArrTypeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrTypeExprContext) INT() antlr.TerminalNode {
	return s.GetToken(nevaParserINT, 0)
}

func (s *ArrTypeExprContext) TypeExpr() ITypeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeExprContext)
}

func (s *ArrTypeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrTypeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrTypeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterArrTypeExpr(s)
	}
}

func (s *ArrTypeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitArrTypeExpr(s)
	}
}

func (p *nevaParser) ArrTypeExpr() (localctx IArrTypeExprContext) {
	localctx = NewArrTypeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, nevaParserRULE_arrTypeExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(248)
		p.Match(nevaParserT__13)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(249)
		p.Match(nevaParserINT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(250)
		p.Match(nevaParserT__14)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(251)
		p.TypeExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRecTypeExprContext is an interface to support dynamic dispatch.
type IRecTypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RecTypeFields() IRecTypeFieldsContext

	// IsRecTypeExprContext differentiates from other interfaces.
	IsRecTypeExprContext()
}

type RecTypeExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRecTypeExprContext() *RecTypeExprContext {
	var p = new(RecTypeExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recTypeExpr
	return p
}

func InitEmptyRecTypeExprContext(p *RecTypeExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recTypeExpr
}

func (*RecTypeExprContext) IsRecTypeExprContext() {}

func NewRecTypeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RecTypeExprContext {
	var p = new(RecTypeExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_recTypeExpr

	return p
}

func (s *RecTypeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *RecTypeExprContext) RecTypeFields() IRecTypeFieldsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRecTypeFieldsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRecTypeFieldsContext)
}

func (s *RecTypeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecTypeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RecTypeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterRecTypeExpr(s)
	}
}

func (s *RecTypeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitRecTypeExpr(s)
	}
}

func (p *nevaParser) RecTypeExpr() (localctx IRecTypeExprContext) {
	localctx = NewRecTypeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, nevaParserRULE_recTypeExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(253)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(255)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserIDENTIFIER {
		{
			p.SetState(254)
			p.RecTypeFields()
		}

	}
	{
		p.SetState(257)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRecTypeFieldsContext is an interface to support dynamic dispatch.
type IRecTypeFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRecTypeField() []IRecTypeFieldContext
	RecTypeField(i int) IRecTypeFieldContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsRecTypeFieldsContext differentiates from other interfaces.
	IsRecTypeFieldsContext()
}

type RecTypeFieldsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRecTypeFieldsContext() *RecTypeFieldsContext {
	var p = new(RecTypeFieldsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recTypeFields
	return p
}

func InitEmptyRecTypeFieldsContext(p *RecTypeFieldsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recTypeFields
}

func (*RecTypeFieldsContext) IsRecTypeFieldsContext() {}

func NewRecTypeFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RecTypeFieldsContext {
	var p = new(RecTypeFieldsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_recTypeFields

	return p
}

func (s *RecTypeFieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *RecTypeFieldsContext) AllRecTypeField() []IRecTypeFieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRecTypeFieldContext); ok {
			len++
		}
	}

	tst := make([]IRecTypeFieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRecTypeFieldContext); ok {
			tst[i] = t.(IRecTypeFieldContext)
			i++
		}
	}

	return tst
}

func (s *RecTypeFieldsContext) RecTypeField(i int) IRecTypeFieldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRecTypeFieldContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRecTypeFieldContext)
}

func (s *RecTypeFieldsContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *RecTypeFieldsContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *RecTypeFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecTypeFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RecTypeFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterRecTypeFields(s)
	}
}

func (s *RecTypeFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitRecTypeFields(s)
	}
}

func (p *nevaParser) RecTypeFields() (localctx IRecTypeFieldsContext) {
	localctx = NewRecTypeFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, nevaParserRULE_recTypeFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(259)
		p.RecTypeField()
	}
	p.SetState(264)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(260)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(261)
			p.RecTypeField()
		}

		p.SetState(266)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRecTypeFieldContext is an interface to support dynamic dispatch.
type IRecTypeFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeExpr() ITypeExprContext

	// IsRecTypeFieldContext differentiates from other interfaces.
	IsRecTypeFieldContext()
}

type RecTypeFieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRecTypeFieldContext() *RecTypeFieldContext {
	var p = new(RecTypeFieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recTypeField
	return p
}

func InitEmptyRecTypeFieldContext(p *RecTypeFieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recTypeField
}

func (*RecTypeFieldContext) IsRecTypeFieldContext() {}

func NewRecTypeFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RecTypeFieldContext {
	var p = new(RecTypeFieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_recTypeField

	return p
}

func (s *RecTypeFieldContext) GetParser() antlr.Parser { return s.parser }

func (s *RecTypeFieldContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *RecTypeFieldContext) TypeExpr() ITypeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeExprContext)
}

func (s *RecTypeFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecTypeFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RecTypeFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterRecTypeField(s)
	}
}

func (s *RecTypeFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitRecTypeField(s)
	}
}

func (p *nevaParser) RecTypeField() (localctx IRecTypeFieldContext) {
	localctx = NewRecTypeFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, nevaParserRULE_recTypeField)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(267)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(268)
		p.TypeExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnionTypeExprContext is an interface to support dynamic dispatch.
type IUnionTypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNonUnionTypeExpr() []INonUnionTypeExprContext
	NonUnionTypeExpr(i int) INonUnionTypeExprContext

	// IsUnionTypeExprContext differentiates from other interfaces.
	IsUnionTypeExprContext()
}

type UnionTypeExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnionTypeExprContext() *UnionTypeExprContext {
	var p = new(UnionTypeExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_unionTypeExpr
	return p
}

func InitEmptyUnionTypeExprContext(p *UnionTypeExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_unionTypeExpr
}

func (*UnionTypeExprContext) IsUnionTypeExprContext() {}

func NewUnionTypeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionTypeExprContext {
	var p = new(UnionTypeExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_unionTypeExpr

	return p
}

func (s *UnionTypeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *UnionTypeExprContext) AllNonUnionTypeExpr() []INonUnionTypeExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INonUnionTypeExprContext); ok {
			len++
		}
	}

	tst := make([]INonUnionTypeExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INonUnionTypeExprContext); ok {
			tst[i] = t.(INonUnionTypeExprContext)
			i++
		}
	}

	return tst
}

func (s *UnionTypeExprContext) NonUnionTypeExpr(i int) INonUnionTypeExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INonUnionTypeExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(INonUnionTypeExprContext)
}

func (s *UnionTypeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionTypeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnionTypeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterUnionTypeExpr(s)
	}
}

func (s *UnionTypeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitUnionTypeExpr(s)
	}
}

func (p *nevaParser) UnionTypeExpr() (localctx IUnionTypeExprContext) {
	localctx = NewUnionTypeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, nevaParserRULE_unionTypeExpr)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(270)
		p.NonUnionTypeExpr()
	}
	p.SetState(273)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(271)
				p.Match(nevaParserT__15)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(272)
				p.NonUnionTypeExpr()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(275)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEnumTypeExprContext is an interface to support dynamic dispatch.
type IEnumTypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsEnumTypeExprContext differentiates from other interfaces.
	IsEnumTypeExprContext()
}

type EnumTypeExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEnumTypeExprContext() *EnumTypeExprContext {
	var p = new(EnumTypeExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_enumTypeExpr
	return p
}

func InitEmptyEnumTypeExprContext(p *EnumTypeExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_enumTypeExpr
}

func (*EnumTypeExprContext) IsEnumTypeExprContext() {}

func NewEnumTypeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumTypeExprContext {
	var p = new(EnumTypeExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_enumTypeExpr

	return p
}

func (s *EnumTypeExprContext) GetParser() antlr.Parser { return s.parser }
func (s *EnumTypeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EnumTypeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EnumTypeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterEnumTypeExpr(s)
	}
}

func (s *EnumTypeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitEnumTypeExpr(s)
	}
}

func (p *nevaParser) EnumTypeExpr() (localctx IEnumTypeExprContext) {
	localctx = NewEnumTypeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, nevaParserRULE_enumTypeExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(277)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INonUnionTypeExprContext is an interface to support dynamic dispatch.
type INonUnionTypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TypeInstExpr() ITypeInstExprContext
	TypeLitExpr() ITypeLitExprContext

	// IsNonUnionTypeExprContext differentiates from other interfaces.
	IsNonUnionTypeExprContext()
}

type NonUnionTypeExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNonUnionTypeExprContext() *NonUnionTypeExprContext {
	var p = new(NonUnionTypeExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_nonUnionTypeExpr
	return p
}

func InitEmptyNonUnionTypeExprContext(p *NonUnionTypeExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_nonUnionTypeExpr
}

func (*NonUnionTypeExprContext) IsNonUnionTypeExprContext() {}

func NewNonUnionTypeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NonUnionTypeExprContext {
	var p = new(NonUnionTypeExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_nonUnionTypeExpr

	return p
}

func (s *NonUnionTypeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *NonUnionTypeExprContext) TypeInstExpr() ITypeInstExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeInstExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeInstExprContext)
}

func (s *NonUnionTypeExprContext) TypeLitExpr() ITypeLitExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeLitExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeLitExprContext)
}

func (s *NonUnionTypeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NonUnionTypeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NonUnionTypeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterNonUnionTypeExpr(s)
	}
}

func (s *NonUnionTypeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitNonUnionTypeExpr(s)
	}
}

func (p *nevaParser) NonUnionTypeExpr() (localctx INonUnionTypeExprContext) {
	localctx = NewNonUnionTypeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, nevaParserRULE_nonUnionTypeExpr)
	p.SetState(281)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(279)
			p.TypeInstExpr()
		}

	case nevaParserT__5, nevaParserT__13:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(280)
			p.TypeLitExpr()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIoStmtContext is an interface to support dynamic dispatch.
type IIoStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	InterfaceDefList() IInterfaceDefListContext

	// IsIoStmtContext differentiates from other interfaces.
	IsIoStmtContext()
}

type IoStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIoStmtContext() *IoStmtContext {
	var p = new(IoStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_ioStmt
	return p
}

func InitEmptyIoStmtContext(p *IoStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_ioStmt
}

func (*IoStmtContext) IsIoStmtContext() {}

func NewIoStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IoStmtContext {
	var p = new(IoStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_ioStmt

	return p
}

func (s *IoStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *IoStmtContext) InterfaceDefList() IInterfaceDefListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInterfaceDefListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInterfaceDefListContext)
}

func (s *IoStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IoStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IoStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterIoStmt(s)
	}
}

func (s *IoStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitIoStmt(s)
	}
}

func (p *nevaParser) IoStmt() (localctx IIoStmtContext) {
	localctx = NewIoStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, nevaParserRULE_ioStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(283)
		p.Match(nevaParserT__16)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(284)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(285)
		p.InterfaceDefList()
	}
	{
		p.SetState(286)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IInterfaceDefListContext is an interface to support dynamic dispatch.
type IInterfaceDefListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllInterfaceDef() []IInterfaceDefContext
	InterfaceDef(i int) IInterfaceDefContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsInterfaceDefListContext differentiates from other interfaces.
	IsInterfaceDefListContext()
}

type InterfaceDefListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInterfaceDefListContext() *InterfaceDefListContext {
	var p = new(InterfaceDefListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_interfaceDefList
	return p
}

func InitEmptyInterfaceDefListContext(p *InterfaceDefListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_interfaceDefList
}

func (*InterfaceDefListContext) IsInterfaceDefListContext() {}

func NewInterfaceDefListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InterfaceDefListContext {
	var p = new(InterfaceDefListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_interfaceDefList

	return p
}

func (s *InterfaceDefListContext) GetParser() antlr.Parser { return s.parser }

func (s *InterfaceDefListContext) AllInterfaceDef() []IInterfaceDefContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IInterfaceDefContext); ok {
			len++
		}
	}

	tst := make([]IInterfaceDefContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IInterfaceDefContext); ok {
			tst[i] = t.(IInterfaceDefContext)
			i++
		}
	}

	return tst
}

func (s *InterfaceDefListContext) InterfaceDef(i int) IInterfaceDefContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInterfaceDefContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInterfaceDefContext)
}

func (s *InterfaceDefListContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *InterfaceDefListContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *InterfaceDefListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InterfaceDefListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *InterfaceDefListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterInterfaceDefList(s)
	}
}

func (s *InterfaceDefListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitInterfaceDefList(s)
	}
}

func (p *nevaParser) InterfaceDefList() (localctx IInterfaceDefListContext) {
	localctx = NewInterfaceDefListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, nevaParserRULE_interfaceDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(288)
		p.InterfaceDef()
	}
	p.SetState(293)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(289)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(290)
			p.InterfaceDef()
		}

		p.SetState(295)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IInterfaceDefContext is an interface to support dynamic dispatch.
type IInterfaceDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeParams() ITypeParamsContext
	AllPortsDef() []IPortsDefContext
	PortsDef(i int) IPortsDefContext

	// IsInterfaceDefContext differentiates from other interfaces.
	IsInterfaceDefContext()
}

type InterfaceDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInterfaceDefContext() *InterfaceDefContext {
	var p = new(InterfaceDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_interfaceDef
	return p
}

func InitEmptyInterfaceDefContext(p *InterfaceDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_interfaceDef
}

func (*InterfaceDefContext) IsInterfaceDefContext() {}

func NewInterfaceDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InterfaceDefContext {
	var p = new(InterfaceDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_interfaceDef

	return p
}

func (s *InterfaceDefContext) GetParser() antlr.Parser { return s.parser }

func (s *InterfaceDefContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *InterfaceDefContext) TypeParams() ITypeParamsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeParamsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeParamsContext)
}

func (s *InterfaceDefContext) AllPortsDef() []IPortsDefContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPortsDefContext); ok {
			len++
		}
	}

	tst := make([]IPortsDefContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPortsDefContext); ok {
			tst[i] = t.(IPortsDefContext)
			i++
		}
	}

	return tst
}

func (s *InterfaceDefContext) PortsDef(i int) IPortsDefContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPortsDefContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPortsDefContext)
}

func (s *InterfaceDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InterfaceDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *InterfaceDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterInterfaceDef(s)
	}
}

func (s *InterfaceDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitInterfaceDef(s)
	}
}

func (p *nevaParser) InterfaceDef() (localctx IInterfaceDefContext) {
	localctx = NewInterfaceDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, nevaParserRULE_interfaceDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(297)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__9 {
		{
			p.SetState(296)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(299)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(300)
		p.TypeParams()
	}
	{
		p.SetState(301)
		p.PortsDef()
	}
	{
		p.SetState(302)
		p.PortsDef()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPortsDefContext is an interface to support dynamic dispatch.
type IPortsDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PortDefList() IPortDefListContext

	// IsPortsDefContext differentiates from other interfaces.
	IsPortsDefContext()
}

type PortsDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPortsDefContext() *PortsDefContext {
	var p = new(PortsDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portsDef
	return p
}

func InitEmptyPortsDefContext(p *PortsDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portsDef
}

func (*PortsDefContext) IsPortsDefContext() {}

func NewPortsDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PortsDefContext {
	var p = new(PortsDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_portsDef

	return p
}

func (s *PortsDefContext) GetParser() antlr.Parser { return s.parser }

func (s *PortsDefContext) PortDefList() IPortDefListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPortDefListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPortDefListContext)
}

func (s *PortsDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PortsDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PortsDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterPortsDef(s)
	}
}

func (s *PortsDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitPortsDef(s)
	}
}

func (p *nevaParser) PortsDef() (localctx IPortsDefContext) {
	localctx = NewPortsDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, nevaParserRULE_portsDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(304)
		p.Match(nevaParserT__17)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(305)
		p.PortDefList()
	}
	{
		p.SetState(306)
		p.Match(nevaParserT__18)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPortDefListContext is an interface to support dynamic dispatch.
type IPortDefListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPortDef() []IPortDefContext
	PortDef(i int) IPortDefContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsPortDefListContext differentiates from other interfaces.
	IsPortDefListContext()
}

type PortDefListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPortDefListContext() *PortDefListContext {
	var p = new(PortDefListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portDefList
	return p
}

func InitEmptyPortDefListContext(p *PortDefListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portDefList
}

func (*PortDefListContext) IsPortDefListContext() {}

func NewPortDefListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PortDefListContext {
	var p = new(PortDefListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_portDefList

	return p
}

func (s *PortDefListContext) GetParser() antlr.Parser { return s.parser }

func (s *PortDefListContext) AllPortDef() []IPortDefContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPortDefContext); ok {
			len++
		}
	}

	tst := make([]IPortDefContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPortDefContext); ok {
			tst[i] = t.(IPortDefContext)
			i++
		}
	}

	return tst
}

func (s *PortDefListContext) PortDef(i int) IPortDefContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPortDefContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPortDefContext)
}

func (s *PortDefListContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *PortDefListContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *PortDefListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PortDefListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PortDefListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterPortDefList(s)
	}
}

func (s *PortDefListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitPortDefList(s)
	}
}

func (p *nevaParser) PortDefList() (localctx IPortDefListContext) {
	localctx = NewPortDefListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, nevaParserRULE_portDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(308)
		p.PortDef()
	}
	p.SetState(316)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__11 {
		{
			p.SetState(309)
			p.Match(nevaParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(311)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserNEWLINE {
			{
				p.SetState(310)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(313)
			p.PortDef()
		}

		p.SetState(318)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPortDefContext is an interface to support dynamic dispatch.
type IPortDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeExpr() ITypeExprContext

	// IsPortDefContext differentiates from other interfaces.
	IsPortDefContext()
}

type PortDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPortDefContext() *PortDefContext {
	var p = new(PortDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portDef
	return p
}

func InitEmptyPortDefContext(p *PortDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portDef
}

func (*PortDefContext) IsPortDefContext() {}

func NewPortDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PortDefContext {
	var p = new(PortDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_portDef

	return p
}

func (s *PortDefContext) GetParser() antlr.Parser { return s.parser }

func (s *PortDefContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *PortDefContext) TypeExpr() ITypeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeExprContext)
}

func (s *PortDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PortDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PortDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterPortDef(s)
	}
}

func (s *PortDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitPortDef(s)
	}
}

func (p *nevaParser) PortDef() (localctx IPortDefContext) {
	localctx = NewPortDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, nevaParserRULE_portDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(319)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(320)
		p.TypeExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConstStmtContext is an interface to support dynamic dispatch.
type IConstStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ConstDefList() IConstDefListContext
	NEWLINE() antlr.TerminalNode

	// IsConstStmtContext differentiates from other interfaces.
	IsConstStmtContext()
}

type ConstStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConstStmtContext() *ConstStmtContext {
	var p = new(ConstStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_constStmt
	return p
}

func InitEmptyConstStmtContext(p *ConstStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_constStmt
}

func (*ConstStmtContext) IsConstStmtContext() {}

func NewConstStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstStmtContext {
	var p = new(ConstStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_constStmt

	return p
}

func (s *ConstStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ConstStmtContext) ConstDefList() IConstDefListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstDefListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstDefListContext)
}

func (s *ConstStmtContext) NEWLINE() antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, 0)
}

func (s *ConstStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConstStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConstStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConstStmt(s)
	}
}

func (s *ConstStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConstStmt(s)
	}
}

func (p *nevaParser) ConstStmt() (localctx IConstStmtContext) {
	localctx = NewConstStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, nevaParserRULE_constStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(322)
		p.Match(nevaParserT__19)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(323)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(324)
		p.ConstDefList()
	}
	{
		p.SetState(325)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(326)
		p.Match(nevaParserNEWLINE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConstDefListContext is an interface to support dynamic dispatch.
type IConstDefListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllConstDef() []IConstDefContext
	ConstDef(i int) IConstDefContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsConstDefListContext differentiates from other interfaces.
	IsConstDefListContext()
}

type ConstDefListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConstDefListContext() *ConstDefListContext {
	var p = new(ConstDefListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_constDefList
	return p
}

func InitEmptyConstDefListContext(p *ConstDefListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_constDefList
}

func (*ConstDefListContext) IsConstDefListContext() {}

func NewConstDefListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstDefListContext {
	var p = new(ConstDefListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_constDefList

	return p
}

func (s *ConstDefListContext) GetParser() antlr.Parser { return s.parser }

func (s *ConstDefListContext) AllConstDef() []IConstDefContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IConstDefContext); ok {
			len++
		}
	}

	tst := make([]IConstDefContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IConstDefContext); ok {
			tst[i] = t.(IConstDefContext)
			i++
		}
	}

	return tst
}

func (s *ConstDefListContext) ConstDef(i int) IConstDefContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstDefContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstDefContext)
}

func (s *ConstDefListContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ConstDefListContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *ConstDefListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConstDefListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConstDefListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConstDefList(s)
	}
}

func (s *ConstDefListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConstDefList(s)
	}
}

func (p *nevaParser) ConstDefList() (localctx IConstDefListContext) {
	localctx = NewConstDefListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, nevaParserRULE_constDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(328)
		p.ConstDef()
	}
	p.SetState(333)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(329)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(330)
			p.ConstDef()
		}

		p.SetState(335)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConstDefContext is an interface to support dynamic dispatch.
type IConstDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeExpr() ITypeExprContext
	ConstValue() IConstValueContext

	// IsConstDefContext differentiates from other interfaces.
	IsConstDefContext()
}

type ConstDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConstDefContext() *ConstDefContext {
	var p = new(ConstDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_constDef
	return p
}

func InitEmptyConstDefContext(p *ConstDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_constDef
}

func (*ConstDefContext) IsConstDefContext() {}

func NewConstDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstDefContext {
	var p = new(ConstDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_constDef

	return p
}

func (s *ConstDefContext) GetParser() antlr.Parser { return s.parser }

func (s *ConstDefContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *ConstDefContext) TypeExpr() ITypeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeExprContext)
}

func (s *ConstDefContext) ConstValue() IConstValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstValueContext)
}

func (s *ConstDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConstDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConstDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConstDef(s)
	}
}

func (s *ConstDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConstDef(s)
	}
}

func (p *nevaParser) ConstDef() (localctx IConstDefContext) {
	localctx = NewConstDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, nevaParserRULE_constDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(337)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__9 {
		{
			p.SetState(336)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(339)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(340)
		p.TypeExpr()
	}
	{
		p.SetState(341)
		p.Match(nevaParserT__20)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(342)
		p.ConstValue()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConstValueContext is an interface to support dynamic dispatch.
type IConstValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INT() antlr.TerminalNode
	FLOAT() antlr.TerminalNode
	STRING() antlr.TerminalNode
	ArrLit() IArrLitContext
	RecLit() IRecLitContext

	// IsConstValueContext differentiates from other interfaces.
	IsConstValueContext()
}

type ConstValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConstValueContext() *ConstValueContext {
	var p = new(ConstValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_constValue
	return p
}

func InitEmptyConstValueContext(p *ConstValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_constValue
}

func (*ConstValueContext) IsConstValueContext() {}

func NewConstValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstValueContext {
	var p = new(ConstValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_constValue

	return p
}

func (s *ConstValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ConstValueContext) INT() antlr.TerminalNode {
	return s.GetToken(nevaParserINT, 0)
}

func (s *ConstValueContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(nevaParserFLOAT, 0)
}

func (s *ConstValueContext) STRING() antlr.TerminalNode {
	return s.GetToken(nevaParserSTRING, 0)
}

func (s *ConstValueContext) ArrLit() IArrLitContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrLitContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrLitContext)
}

func (s *ConstValueContext) RecLit() IRecLitContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRecLitContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRecLitContext)
}

func (s *ConstValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConstValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConstValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConstValue(s)
	}
}

func (s *ConstValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConstValue(s)
	}
}

func (p *nevaParser) ConstValue() (localctx IConstValueContext) {
	localctx = NewConstValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, nevaParserRULE_constValue)
	p.SetState(352)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__21:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(344)
			p.Match(nevaParserT__21)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserT__22:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(345)
			p.Match(nevaParserT__22)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserINT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(346)
			p.Match(nevaParserINT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserFLOAT:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(347)
			p.Match(nevaParserFLOAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserSTRING:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(348)
			p.Match(nevaParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserT__13:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(349)
			p.ArrLit()
		}

	case nevaParserT__5:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(350)
			p.RecLit()
		}

	case nevaParserT__23:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(351)
			p.Match(nevaParserT__23)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArrLitContext is an interface to support dynamic dispatch.
type IArrLitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ArrItems() IArrItemsContext

	// IsArrLitContext differentiates from other interfaces.
	IsArrLitContext()
}

type ArrLitContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrLitContext() *ArrLitContext {
	var p = new(ArrLitContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_arrLit
	return p
}

func InitEmptyArrLitContext(p *ArrLitContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_arrLit
}

func (*ArrLitContext) IsArrLitContext() {}

func NewArrLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrLitContext {
	var p = new(ArrLitContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_arrLit

	return p
}

func (s *ArrLitContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrLitContext) ArrItems() IArrItemsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrItemsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrItemsContext)
}

func (s *ArrLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrLitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterArrLit(s)
	}
}

func (s *ArrLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitArrLit(s)
	}
}

func (p *nevaParser) ArrLit() (localctx IArrLitContext) {
	localctx = NewArrLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, nevaParserRULE_arrLit)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(354)
		p.Match(nevaParserT__13)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(355)
		p.ArrItems()
	}
	{
		p.SetState(356)
		p.Match(nevaParserT__14)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArrItemsContext is an interface to support dynamic dispatch.
type IArrItemsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllConstValue() []IConstValueContext
	ConstValue(i int) IConstValueContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsArrItemsContext differentiates from other interfaces.
	IsArrItemsContext()
}

type ArrItemsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrItemsContext() *ArrItemsContext {
	var p = new(ArrItemsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_arrItems
	return p
}

func InitEmptyArrItemsContext(p *ArrItemsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_arrItems
}

func (*ArrItemsContext) IsArrItemsContext() {}

func NewArrItemsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrItemsContext {
	var p = new(ArrItemsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_arrItems

	return p
}

func (s *ArrItemsContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrItemsContext) AllConstValue() []IConstValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IConstValueContext); ok {
			len++
		}
	}

	tst := make([]IConstValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IConstValueContext); ok {
			tst[i] = t.(IConstValueContext)
			i++
		}
	}

	return tst
}

func (s *ArrItemsContext) ConstValue(i int) IConstValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstValueContext)
}

func (s *ArrItemsContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ArrItemsContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *ArrItemsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrItemsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrItemsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterArrItems(s)
	}
}

func (s *ArrItemsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitArrItems(s)
	}
}

func (p *nevaParser) ArrItems() (localctx IArrItemsContext) {
	localctx = NewArrItemsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, nevaParserRULE_arrItems)
	var _la int

	p.SetState(370)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 32, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(358)
			p.ConstValue()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(359)
			p.ConstValue()
		}
		p.SetState(367)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserT__11 {
			{
				p.SetState(360)
				p.Match(nevaParserT__11)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(362)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if _la == nevaParserNEWLINE {
				{
					p.SetState(361)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			}
			{
				p.SetState(364)
				p.ConstValue()
			}

			p.SetState(369)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRecLitContext is an interface to support dynamic dispatch.
type IRecLitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RecValueFields() IRecValueFieldsContext

	// IsRecLitContext differentiates from other interfaces.
	IsRecLitContext()
}

type RecLitContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRecLitContext() *RecLitContext {
	var p = new(RecLitContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recLit
	return p
}

func InitEmptyRecLitContext(p *RecLitContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recLit
}

func (*RecLitContext) IsRecLitContext() {}

func NewRecLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RecLitContext {
	var p = new(RecLitContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_recLit

	return p
}

func (s *RecLitContext) GetParser() antlr.Parser { return s.parser }

func (s *RecLitContext) RecValueFields() IRecValueFieldsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRecValueFieldsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRecValueFieldsContext)
}

func (s *RecLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecLitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RecLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterRecLit(s)
	}
}

func (s *RecLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitRecLit(s)
	}
}

func (p *nevaParser) RecLit() (localctx IRecLitContext) {
	localctx = NewRecLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, nevaParserRULE_recLit)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(372)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(373)
		p.RecValueFields()
	}
	{
		p.SetState(374)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRecValueFieldsContext is an interface to support dynamic dispatch.
type IRecValueFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRecValueField() []IRecValueFieldContext
	RecValueField(i int) IRecValueFieldContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsRecValueFieldsContext differentiates from other interfaces.
	IsRecValueFieldsContext()
}

type RecValueFieldsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRecValueFieldsContext() *RecValueFieldsContext {
	var p = new(RecValueFieldsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recValueFields
	return p
}

func InitEmptyRecValueFieldsContext(p *RecValueFieldsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recValueFields
}

func (*RecValueFieldsContext) IsRecValueFieldsContext() {}

func NewRecValueFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RecValueFieldsContext {
	var p = new(RecValueFieldsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_recValueFields

	return p
}

func (s *RecValueFieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *RecValueFieldsContext) AllRecValueField() []IRecValueFieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRecValueFieldContext); ok {
			len++
		}
	}

	tst := make([]IRecValueFieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRecValueFieldContext); ok {
			tst[i] = t.(IRecValueFieldContext)
			i++
		}
	}

	return tst
}

func (s *RecValueFieldsContext) RecValueField(i int) IRecValueFieldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRecValueFieldContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRecValueFieldContext)
}

func (s *RecValueFieldsContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *RecValueFieldsContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *RecValueFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecValueFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RecValueFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterRecValueFields(s)
	}
}

func (s *RecValueFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitRecValueFields(s)
	}
}

func (p *nevaParser) RecValueFields() (localctx IRecValueFieldsContext) {
	localctx = NewRecValueFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, nevaParserRULE_recValueFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(376)
		p.RecValueField()
	}
	p.SetState(384)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__11 {
		{
			p.SetState(377)
			p.Match(nevaParserT__11)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(379)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserNEWLINE {
			{
				p.SetState(378)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(381)
			p.RecValueField()
		}

		p.SetState(386)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRecValueFieldContext is an interface to support dynamic dispatch.
type IRecValueFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	ConstValue() IConstValueContext

	// IsRecValueFieldContext differentiates from other interfaces.
	IsRecValueFieldContext()
}

type RecValueFieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRecValueFieldContext() *RecValueFieldContext {
	var p = new(RecValueFieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recValueField
	return p
}

func InitEmptyRecValueFieldContext(p *RecValueFieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recValueField
}

func (*RecValueFieldContext) IsRecValueFieldContext() {}

func NewRecValueFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RecValueFieldContext {
	var p = new(RecValueFieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_recValueField

	return p
}

func (s *RecValueFieldContext) GetParser() antlr.Parser { return s.parser }

func (s *RecValueFieldContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *RecValueFieldContext) ConstValue() IConstValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstValueContext)
}

func (s *RecValueFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecValueFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RecValueFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterRecValueField(s)
	}
}

func (s *RecValueFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitRecValueField(s)
	}
}

func (p *nevaParser) RecValueField() (localctx IRecValueFieldContext) {
	localctx = NewRecValueFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, nevaParserRULE_recValueField)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(387)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(388)
		p.Match(nevaParserT__24)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(389)
		p.ConstValue()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompStmtContext is an interface to support dynamic dispatch.
type ICompStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CompDefList() ICompDefListContext
	NEWLINE() antlr.TerminalNode

	// IsCompStmtContext differentiates from other interfaces.
	IsCompStmtContext()
}

type CompStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompStmtContext() *CompStmtContext {
	var p = new(CompStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compStmt
	return p
}

func InitEmptyCompStmtContext(p *CompStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compStmt
}

func (*CompStmtContext) IsCompStmtContext() {}

func NewCompStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompStmtContext {
	var p = new(CompStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_compStmt

	return p
}

func (s *CompStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *CompStmtContext) CompDefList() ICompDefListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompDefListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompDefListContext)
}

func (s *CompStmtContext) NEWLINE() antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, 0)
}

func (s *CompStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterCompStmt(s)
	}
}

func (s *CompStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitCompStmt(s)
	}
}

func (p *nevaParser) CompStmt() (localctx ICompStmtContext) {
	localctx = NewCompStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, nevaParserRULE_compStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(391)
		p.Match(nevaParserT__25)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(392)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(393)
		p.CompDefList()
	}
	{
		p.SetState(394)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(395)
		p.Match(nevaParserNEWLINE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompDefListContext is an interface to support dynamic dispatch.
type ICompDefListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCompDef() []ICompDefContext
	CompDef(i int) ICompDefContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsCompDefListContext differentiates from other interfaces.
	IsCompDefListContext()
}

type CompDefListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompDefListContext() *CompDefListContext {
	var p = new(CompDefListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compDefList
	return p
}

func InitEmptyCompDefListContext(p *CompDefListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compDefList
}

func (*CompDefListContext) IsCompDefListContext() {}

func NewCompDefListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompDefListContext {
	var p = new(CompDefListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_compDefList

	return p
}

func (s *CompDefListContext) GetParser() antlr.Parser { return s.parser }

func (s *CompDefListContext) AllCompDef() []ICompDefContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICompDefContext); ok {
			len++
		}
	}

	tst := make([]ICompDefContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICompDefContext); ok {
			tst[i] = t.(ICompDefContext)
			i++
		}
	}

	return tst
}

func (s *CompDefListContext) CompDef(i int) ICompDefContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompDefContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompDefContext)
}

func (s *CompDefListContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *CompDefListContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *CompDefListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompDefListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompDefListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterCompDefList(s)
	}
}

func (s *CompDefListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitCompDefList(s)
	}
}

func (p *nevaParser) CompDefList() (localctx ICompDefListContext) {
	localctx = NewCompDefListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, nevaParserRULE_compDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(397)
		p.CompDef()
	}
	p.SetState(402)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(398)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(399)
			p.CompDef()
		}

		p.SetState(404)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompDefContext is an interface to support dynamic dispatch.
type ICompDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	InterfaceDef() IInterfaceDefContext
	CompBody() ICompBodyContext

	// IsCompDefContext differentiates from other interfaces.
	IsCompDefContext()
}

type CompDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompDefContext() *CompDefContext {
	var p = new(CompDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compDef
	return p
}

func InitEmptyCompDefContext(p *CompDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compDef
}

func (*CompDefContext) IsCompDefContext() {}

func NewCompDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompDefContext {
	var p = new(CompDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_compDef

	return p
}

func (s *CompDefContext) GetParser() antlr.Parser { return s.parser }

func (s *CompDefContext) InterfaceDef() IInterfaceDefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInterfaceDefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInterfaceDefContext)
}

func (s *CompDefContext) CompBody() ICompBodyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompBodyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompBodyContext)
}

func (s *CompDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterCompDef(s)
	}
}

func (s *CompDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitCompDef(s)
	}
}

func (p *nevaParser) CompDef() (localctx ICompDefContext) {
	localctx = NewCompDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, nevaParserRULE_compDef)
	p.EnterOuterAlt(localctx, 1)
	p.SetState(406)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 36, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(405)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(408)
		p.InterfaceDef()
	}
	{
		p.SetState(409)
		p.CompBody()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompBodyContext is an interface to support dynamic dispatch.
type ICompBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CompNodesDef() ICompNodesDefContext
	CompNetDef() ICompNetDefContext

	// IsCompBodyContext differentiates from other interfaces.
	IsCompBodyContext()
}

type CompBodyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompBodyContext() *CompBodyContext {
	var p = new(CompBodyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compBody
	return p
}

func InitEmptyCompBodyContext(p *CompBodyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compBody
}

func (*CompBodyContext) IsCompBodyContext() {}

func NewCompBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompBodyContext {
	var p = new(CompBodyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_compBody

	return p
}

func (s *CompBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *CompBodyContext) CompNodesDef() ICompNodesDefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompNodesDefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompNodesDefContext)
}

func (s *CompBodyContext) CompNetDef() ICompNetDefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompNetDefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompNetDefContext)
}

func (s *CompBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompBodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterCompBody(s)
	}
}

func (s *CompBodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitCompBody(s)
	}
}

func (p *nevaParser) CompBody() (localctx ICompBodyContext) {
	localctx = NewCompBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, nevaParserRULE_compBody)
	p.SetState(416)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__5:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(411)
			p.Match(nevaParserT__5)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(412)
			p.CompNodesDef()
		}

	case nevaParserT__28:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(413)
			p.CompNetDef()
		}
		{
			p.SetState(414)
			p.Match(nevaParserT__6)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompNodesDefContext is an interface to support dynamic dispatch.
type ICompNodesDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CompNodeDefList() ICompNodeDefListContext

	// IsCompNodesDefContext differentiates from other interfaces.
	IsCompNodesDefContext()
}

type CompNodesDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompNodesDefContext() *CompNodesDefContext {
	var p = new(CompNodesDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compNodesDef
	return p
}

func InitEmptyCompNodesDefContext(p *CompNodesDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compNodesDef
}

func (*CompNodesDefContext) IsCompNodesDefContext() {}

func NewCompNodesDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompNodesDefContext {
	var p = new(CompNodesDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_compNodesDef

	return p
}

func (s *CompNodesDefContext) GetParser() antlr.Parser { return s.parser }

func (s *CompNodesDefContext) CompNodeDefList() ICompNodeDefListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompNodeDefListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompNodeDefListContext)
}

func (s *CompNodesDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompNodesDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompNodesDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterCompNodesDef(s)
	}
}

func (s *CompNodesDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitCompNodesDef(s)
	}
}

func (p *nevaParser) CompNodesDef() (localctx ICompNodesDefContext) {
	localctx = NewCompNodesDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, nevaParserRULE_compNodesDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(418)
		p.Match(nevaParserT__26)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(419)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(420)
		p.CompNodeDefList()
	}
	{
		p.SetState(421)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompNodeDefListContext is an interface to support dynamic dispatch.
type ICompNodeDefListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AbsNodeDef() IAbsNodeDefContext
	ConcreteNodeDef() IConcreteNodeDefContext

	// IsCompNodeDefListContext differentiates from other interfaces.
	IsCompNodeDefListContext()
}

type CompNodeDefListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompNodeDefListContext() *CompNodeDefListContext {
	var p = new(CompNodeDefListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compNodeDefList
	return p
}

func InitEmptyCompNodeDefListContext(p *CompNodeDefListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compNodeDefList
}

func (*CompNodeDefListContext) IsCompNodeDefListContext() {}

func NewCompNodeDefListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompNodeDefListContext {
	var p = new(CompNodeDefListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_compNodeDefList

	return p
}

func (s *CompNodeDefListContext) GetParser() antlr.Parser { return s.parser }

func (s *CompNodeDefListContext) AbsNodeDef() IAbsNodeDefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAbsNodeDefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAbsNodeDefContext)
}

func (s *CompNodeDefListContext) ConcreteNodeDef() IConcreteNodeDefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConcreteNodeDefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConcreteNodeDefContext)
}

func (s *CompNodeDefListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompNodeDefListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompNodeDefListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterCompNodeDefList(s)
	}
}

func (s *CompNodeDefListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitCompNodeDefList(s)
	}
}

func (p *nevaParser) CompNodeDefList() (localctx ICompNodeDefListContext) {
	localctx = NewCompNodeDefListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, nevaParserRULE_compNodeDefList)
	p.SetState(425)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 38, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(423)
			p.AbsNodeDef()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(424)
			p.ConcreteNodeDef()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAbsNodeDefContext is an interface to support dynamic dispatch.
type IAbsNodeDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeInstExpr() ITypeInstExprContext

	// IsAbsNodeDefContext differentiates from other interfaces.
	IsAbsNodeDefContext()
}

type AbsNodeDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAbsNodeDefContext() *AbsNodeDefContext {
	var p = new(AbsNodeDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_absNodeDef
	return p
}

func InitEmptyAbsNodeDefContext(p *AbsNodeDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_absNodeDef
}

func (*AbsNodeDefContext) IsAbsNodeDefContext() {}

func NewAbsNodeDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AbsNodeDefContext {
	var p = new(AbsNodeDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_absNodeDef

	return p
}

func (s *AbsNodeDefContext) GetParser() antlr.Parser { return s.parser }

func (s *AbsNodeDefContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *AbsNodeDefContext) TypeInstExpr() ITypeInstExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeInstExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeInstExprContext)
}

func (s *AbsNodeDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AbsNodeDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AbsNodeDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterAbsNodeDef(s)
	}
}

func (s *AbsNodeDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitAbsNodeDef(s)
	}
}

func (p *nevaParser) AbsNodeDef() (localctx IAbsNodeDefContext) {
	localctx = NewAbsNodeDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, nevaParserRULE_absNodeDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(427)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(428)
		p.TypeInstExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConcreteNodeDefContext is an interface to support dynamic dispatch.
type IConcreteNodeDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	ConcreteNodeInst() IConcreteNodeInstContext

	// IsConcreteNodeDefContext differentiates from other interfaces.
	IsConcreteNodeDefContext()
}

type ConcreteNodeDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConcreteNodeDefContext() *ConcreteNodeDefContext {
	var p = new(ConcreteNodeDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_concreteNodeDef
	return p
}

func InitEmptyConcreteNodeDefContext(p *ConcreteNodeDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_concreteNodeDef
}

func (*ConcreteNodeDefContext) IsConcreteNodeDefContext() {}

func NewConcreteNodeDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConcreteNodeDefContext {
	var p = new(ConcreteNodeDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_concreteNodeDef

	return p
}

func (s *ConcreteNodeDefContext) GetParser() antlr.Parser { return s.parser }

func (s *ConcreteNodeDefContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *ConcreteNodeDefContext) ConcreteNodeInst() IConcreteNodeInstContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConcreteNodeInstContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConcreteNodeInstContext)
}

func (s *ConcreteNodeDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConcreteNodeDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConcreteNodeDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConcreteNodeDef(s)
	}
}

func (s *ConcreteNodeDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConcreteNodeDef(s)
	}
}

func (p *nevaParser) ConcreteNodeDef() (localctx IConcreteNodeDefContext) {
	localctx = NewConcreteNodeDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, nevaParserRULE_concreteNodeDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(430)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(431)
		p.Match(nevaParserT__20)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(432)
		p.ConcreteNodeInst()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConcreteNodeInstContext is an interface to support dynamic dispatch.
type IConcreteNodeInstContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NodeRef() INodeRefContext
	NodeArgs() INodeArgsContext
	TypeArgs() ITypeArgsContext

	// IsConcreteNodeInstContext differentiates from other interfaces.
	IsConcreteNodeInstContext()
}

type ConcreteNodeInstContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConcreteNodeInstContext() *ConcreteNodeInstContext {
	var p = new(ConcreteNodeInstContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_concreteNodeInst
	return p
}

func InitEmptyConcreteNodeInstContext(p *ConcreteNodeInstContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_concreteNodeInst
}

func (*ConcreteNodeInstContext) IsConcreteNodeInstContext() {}

func NewConcreteNodeInstContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConcreteNodeInstContext {
	var p = new(ConcreteNodeInstContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_concreteNodeInst

	return p
}

func (s *ConcreteNodeInstContext) GetParser() antlr.Parser { return s.parser }

func (s *ConcreteNodeInstContext) NodeRef() INodeRefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INodeRefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INodeRefContext)
}

func (s *ConcreteNodeInstContext) NodeArgs() INodeArgsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INodeArgsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INodeArgsContext)
}

func (s *ConcreteNodeInstContext) TypeArgs() ITypeArgsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeArgsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeArgsContext)
}

func (s *ConcreteNodeInstContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConcreteNodeInstContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConcreteNodeInstContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConcreteNodeInst(s)
	}
}

func (s *ConcreteNodeInstContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConcreteNodeInst(s)
	}
}

func (p *nevaParser) ConcreteNodeInst() (localctx IConcreteNodeInstContext) {
	localctx = NewConcreteNodeInstContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, nevaParserRULE_concreteNodeInst)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(434)
		p.NodeRef()
	}
	{
		p.SetState(435)
		p.NodeArgs()
	}
	{
		p.SetState(436)
		p.TypeArgs()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INodeRefContext is an interface to support dynamic dispatch.
type INodeRefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode

	// IsNodeRefContext differentiates from other interfaces.
	IsNodeRefContext()
}

type NodeRefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNodeRefContext() *NodeRefContext {
	var p = new(NodeRefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_nodeRef
	return p
}

func InitEmptyNodeRefContext(p *NodeRefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_nodeRef
}

func (*NodeRefContext) IsNodeRefContext() {}

func NewNodeRefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NodeRefContext {
	var p = new(NodeRefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_nodeRef

	return p
}

func (s *NodeRefContext) GetParser() antlr.Parser { return s.parser }

func (s *NodeRefContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(nevaParserIDENTIFIER)
}

func (s *NodeRefContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, i)
}

func (s *NodeRefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NodeRefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NodeRefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterNodeRef(s)
	}
}

func (s *NodeRefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitNodeRef(s)
	}
}

func (p *nevaParser) NodeRef() (localctx INodeRefContext) {
	localctx = NewNodeRefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, nevaParserRULE_nodeRef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(438)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(443)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__27 {
		{
			p.SetState(439)
			p.Match(nevaParserT__27)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(440)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(445)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INodeArgsContext is an interface to support dynamic dispatch.
type INodeArgsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NodeArgList() INodeArgListContext

	// IsNodeArgsContext differentiates from other interfaces.
	IsNodeArgsContext()
}

type NodeArgsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNodeArgsContext() *NodeArgsContext {
	var p = new(NodeArgsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_nodeArgs
	return p
}

func InitEmptyNodeArgsContext(p *NodeArgsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_nodeArgs
}

func (*NodeArgsContext) IsNodeArgsContext() {}

func NewNodeArgsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NodeArgsContext {
	var p = new(NodeArgsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_nodeArgs

	return p
}

func (s *NodeArgsContext) GetParser() antlr.Parser { return s.parser }

func (s *NodeArgsContext) NodeArgList() INodeArgListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INodeArgListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INodeArgListContext)
}

func (s *NodeArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NodeArgsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NodeArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterNodeArgs(s)
	}
}

func (s *NodeArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitNodeArgs(s)
	}
}

func (p *nevaParser) NodeArgs() (localctx INodeArgsContext) {
	localctx = NewNodeArgsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, nevaParserRULE_nodeArgs)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(446)
		p.Match(nevaParserT__17)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(447)
		p.NodeArgList()
	}
	{
		p.SetState(448)
		p.Match(nevaParserT__18)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INodeArgListContext is an interface to support dynamic dispatch.
type INodeArgListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNodeArg() []INodeArgContext
	NodeArg(i int) INodeArgContext
	NEWLINE() antlr.TerminalNode

	// IsNodeArgListContext differentiates from other interfaces.
	IsNodeArgListContext()
}

type NodeArgListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNodeArgListContext() *NodeArgListContext {
	var p = new(NodeArgListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_nodeArgList
	return p
}

func InitEmptyNodeArgListContext(p *NodeArgListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_nodeArgList
}

func (*NodeArgListContext) IsNodeArgListContext() {}

func NewNodeArgListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NodeArgListContext {
	var p = new(NodeArgListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_nodeArgList

	return p
}

func (s *NodeArgListContext) GetParser() antlr.Parser { return s.parser }

func (s *NodeArgListContext) AllNodeArg() []INodeArgContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INodeArgContext); ok {
			len++
		}
	}

	tst := make([]INodeArgContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INodeArgContext); ok {
			tst[i] = t.(INodeArgContext)
			i++
		}
	}

	return tst
}

func (s *NodeArgListContext) NodeArg(i int) INodeArgContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INodeArgContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(INodeArgContext)
}

func (s *NodeArgListContext) NEWLINE() antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, 0)
}

func (s *NodeArgListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NodeArgListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NodeArgListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterNodeArgList(s)
	}
}

func (s *NodeArgListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitNodeArgList(s)
	}
}

func (p *nevaParser) NodeArgList() (localctx INodeArgListContext) {
	localctx = NewNodeArgListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, nevaParserRULE_nodeArgList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(450)
		p.NodeArg()
	}

	{
		p.SetState(451)
		p.Match(nevaParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(453)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserNEWLINE {
		{
			p.SetState(452)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(455)
		p.NodeArg()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INodeArgContext is an interface to support dynamic dispatch.
type INodeArgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ConcreteNodeInst() IConcreteNodeInstContext

	// IsNodeArgContext differentiates from other interfaces.
	IsNodeArgContext()
}

type NodeArgContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNodeArgContext() *NodeArgContext {
	var p = new(NodeArgContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_nodeArg
	return p
}

func InitEmptyNodeArgContext(p *NodeArgContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_nodeArg
}

func (*NodeArgContext) IsNodeArgContext() {}

func NewNodeArgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NodeArgContext {
	var p = new(NodeArgContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_nodeArg

	return p
}

func (s *NodeArgContext) GetParser() antlr.Parser { return s.parser }

func (s *NodeArgContext) ConcreteNodeInst() IConcreteNodeInstContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConcreteNodeInstContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConcreteNodeInstContext)
}

func (s *NodeArgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NodeArgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NodeArgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterNodeArg(s)
	}
}

func (s *NodeArgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitNodeArg(s)
	}
}

func (p *nevaParser) NodeArg() (localctx INodeArgContext) {
	localctx = NewNodeArgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, nevaParserRULE_nodeArg)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(457)
		p.ConcreteNodeInst()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompNetDefContext is an interface to support dynamic dispatch.
type ICompNetDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ConnDefList() IConnDefListContext

	// IsCompNetDefContext differentiates from other interfaces.
	IsCompNetDefContext()
}

type CompNetDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompNetDefContext() *CompNetDefContext {
	var p = new(CompNetDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compNetDef
	return p
}

func InitEmptyCompNetDefContext(p *CompNetDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_compNetDef
}

func (*CompNetDefContext) IsCompNetDefContext() {}

func NewCompNetDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompNetDefContext {
	var p = new(CompNetDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_compNetDef

	return p
}

func (s *CompNetDefContext) GetParser() antlr.Parser { return s.parser }

func (s *CompNetDefContext) ConnDefList() IConnDefListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConnDefListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConnDefListContext)
}

func (s *CompNetDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompNetDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompNetDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterCompNetDef(s)
	}
}

func (s *CompNetDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitCompNetDef(s)
	}
}

func (p *nevaParser) CompNetDef() (localctx ICompNetDefContext) {
	localctx = NewCompNetDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, nevaParserRULE_compNetDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(459)
		p.Match(nevaParserT__28)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(460)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(461)
		p.ConnDefList()
	}
	{
		p.SetState(462)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConnDefListContext is an interface to support dynamic dispatch.
type IConnDefListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllConnDef() []IConnDefContext
	ConnDef(i int) IConnDefContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsConnDefListContext differentiates from other interfaces.
	IsConnDefListContext()
}

type ConnDefListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConnDefListContext() *ConnDefListContext {
	var p = new(ConnDefListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_connDefList
	return p
}

func InitEmptyConnDefListContext(p *ConnDefListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_connDefList
}

func (*ConnDefListContext) IsConnDefListContext() {}

func NewConnDefListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConnDefListContext {
	var p = new(ConnDefListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_connDefList

	return p
}

func (s *ConnDefListContext) GetParser() antlr.Parser { return s.parser }

func (s *ConnDefListContext) AllConnDef() []IConnDefContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IConnDefContext); ok {
			len++
		}
	}

	tst := make([]IConnDefContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IConnDefContext); ok {
			tst[i] = t.(IConnDefContext)
			i++
		}
	}

	return tst
}

func (s *ConnDefListContext) ConnDef(i int) IConnDefContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConnDefContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConnDefContext)
}

func (s *ConnDefListContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ConnDefListContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *ConnDefListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConnDefListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConnDefListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConnDefList(s)
	}
}

func (s *ConnDefListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConnDefList(s)
	}
}

func (p *nevaParser) ConnDefList() (localctx IConnDefListContext) {
	localctx = NewConnDefListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, nevaParserRULE_connDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(464)
		p.ConnDef()
	}
	p.SetState(469)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(465)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(466)
			p.ConnDef()
		}

		p.SetState(471)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConnDefContext is an interface to support dynamic dispatch.
type IConnDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PortAddr() IPortAddrContext
	ConnReceiverSide() IConnReceiverSideContext

	// IsConnDefContext differentiates from other interfaces.
	IsConnDefContext()
}

type ConnDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConnDefContext() *ConnDefContext {
	var p = new(ConnDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_connDef
	return p
}

func InitEmptyConnDefContext(p *ConnDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_connDef
}

func (*ConnDefContext) IsConnDefContext() {}

func NewConnDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConnDefContext {
	var p = new(ConnDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_connDef

	return p
}

func (s *ConnDefContext) GetParser() antlr.Parser { return s.parser }

func (s *ConnDefContext) PortAddr() IPortAddrContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPortAddrContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPortAddrContext)
}

func (s *ConnDefContext) ConnReceiverSide() IConnReceiverSideContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConnReceiverSideContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConnReceiverSideContext)
}

func (s *ConnDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConnDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConnDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConnDef(s)
	}
}

func (s *ConnDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConnDef(s)
	}
}

func (p *nevaParser) ConnDef() (localctx IConnDefContext) {
	localctx = NewConnDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, nevaParserRULE_connDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(472)
		p.PortAddr()
	}
	{
		p.SetState(473)
		p.Match(nevaParserT__29)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(474)
		p.ConnReceiverSide()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPortAddrContext is an interface to support dynamic dispatch.
type IPortAddrContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PortDirection() IPortDirectionContext
	IDENTIFIER() antlr.TerminalNode
	INT() antlr.TerminalNode

	// IsPortAddrContext differentiates from other interfaces.
	IsPortAddrContext()
}

type PortAddrContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPortAddrContext() *PortAddrContext {
	var p = new(PortAddrContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portAddr
	return p
}

func InitEmptyPortAddrContext(p *PortAddrContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portAddr
}

func (*PortAddrContext) IsPortAddrContext() {}

func NewPortAddrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PortAddrContext {
	var p = new(PortAddrContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_portAddr

	return p
}

func (s *PortAddrContext) GetParser() antlr.Parser { return s.parser }

func (s *PortAddrContext) PortDirection() IPortDirectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPortDirectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPortDirectionContext)
}

func (s *PortAddrContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *PortAddrContext) INT() antlr.TerminalNode {
	return s.GetToken(nevaParserINT, 0)
}

func (s *PortAddrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PortAddrContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PortAddrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterPortAddr(s)
	}
}

func (s *PortAddrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitPortAddr(s)
	}
}

func (p *nevaParser) PortAddr() (localctx IPortAddrContext) {
	localctx = NewPortAddrContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, nevaParserRULE_portAddr)
	var _la int

	p.SetState(486)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(477)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserIDENTIFIER {
			{
				p.SetState(476)
				p.Match(nevaParserIDENTIFIER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(479)
			p.PortDirection()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(480)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(484)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserT__13 {
			{
				p.SetState(481)
				p.Match(nevaParserT__13)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(482)
				p.Match(nevaParserINT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(483)
				p.Match(nevaParserT__14)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPortDirectionContext is an interface to support dynamic dispatch.
type IPortDirectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPortDirectionContext differentiates from other interfaces.
	IsPortDirectionContext()
}

type PortDirectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPortDirectionContext() *PortDirectionContext {
	var p = new(PortDirectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portDirection
	return p
}

func InitEmptyPortDirectionContext(p *PortDirectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portDirection
}

func (*PortDirectionContext) IsPortDirectionContext() {}

func NewPortDirectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PortDirectionContext {
	var p = new(PortDirectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_portDirection

	return p
}

func (s *PortDirectionContext) GetParser() antlr.Parser { return s.parser }
func (s *PortDirectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PortDirectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PortDirectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterPortDirection(s)
	}
}

func (s *PortDirectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitPortDirection(s)
	}
}

func (p *nevaParser) PortDirection() (localctx IPortDirectionContext) {
	localctx = NewPortDirectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, nevaParserRULE_portDirection)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(488)
		_la = p.GetTokenStream().LA(1)

		if !(_la == nevaParserT__30 || _la == nevaParserT__31) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConnReceiverSideContext is an interface to support dynamic dispatch.
type IConnReceiverSideContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PortAddr() IPortAddrContext
	ConnReceivers() IConnReceiversContext

	// IsConnReceiverSideContext differentiates from other interfaces.
	IsConnReceiverSideContext()
}

type ConnReceiverSideContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConnReceiverSideContext() *ConnReceiverSideContext {
	var p = new(ConnReceiverSideContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_connReceiverSide
	return p
}

func InitEmptyConnReceiverSideContext(p *ConnReceiverSideContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_connReceiverSide
}

func (*ConnReceiverSideContext) IsConnReceiverSideContext() {}

func NewConnReceiverSideContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConnReceiverSideContext {
	var p = new(ConnReceiverSideContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_connReceiverSide

	return p
}

func (s *ConnReceiverSideContext) GetParser() antlr.Parser { return s.parser }

func (s *ConnReceiverSideContext) PortAddr() IPortAddrContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPortAddrContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPortAddrContext)
}

func (s *ConnReceiverSideContext) ConnReceivers() IConnReceiversContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConnReceiversContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConnReceiversContext)
}

func (s *ConnReceiverSideContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConnReceiverSideContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConnReceiverSideContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConnReceiverSide(s)
	}
}

func (s *ConnReceiverSideContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConnReceiverSide(s)
	}
}

func (p *nevaParser) ConnReceiverSide() (localctx IConnReceiverSideContext) {
	localctx = NewConnReceiverSideContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, nevaParserRULE_connReceiverSide)
	p.SetState(492)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__30, nevaParserT__31, nevaParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(490)
			p.PortAddr()
		}

	case nevaParserT__5:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(491)
			p.ConnReceivers()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConnReceiversContext is an interface to support dynamic dispatch.
type IConnReceiversContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPortAddr() []IPortAddrContext
	PortAddr(i int) IPortAddrContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsConnReceiversContext differentiates from other interfaces.
	IsConnReceiversContext()
}

type ConnReceiversContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConnReceiversContext() *ConnReceiversContext {
	var p = new(ConnReceiversContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_connReceivers
	return p
}

func InitEmptyConnReceiversContext(p *ConnReceiversContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_connReceivers
}

func (*ConnReceiversContext) IsConnReceiversContext() {}

func NewConnReceiversContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConnReceiversContext {
	var p = new(ConnReceiversContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_connReceivers

	return p
}

func (s *ConnReceiversContext) GetParser() antlr.Parser { return s.parser }

func (s *ConnReceiversContext) AllPortAddr() []IPortAddrContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPortAddrContext); ok {
			len++
		}
	}

	tst := make([]IPortAddrContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPortAddrContext); ok {
			tst[i] = t.(IPortAddrContext)
			i++
		}
	}

	return tst
}

func (s *ConnReceiversContext) PortAddr(i int) IPortAddrContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPortAddrContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPortAddrContext)
}

func (s *ConnReceiversContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ConnReceiversContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *ConnReceiversContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConnReceiversContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConnReceiversContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConnReceivers(s)
	}
}

func (s *ConnReceiversContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConnReceivers(s)
	}
}

func (p *nevaParser) ConnReceivers() (localctx IConnReceiversContext) {
	localctx = NewConnReceiversContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 118, nevaParserRULE_connReceivers)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(494)
		p.Match(nevaParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(495)
		p.PortAddr()
	}
	p.SetState(500)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(496)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(497)
			p.PortAddr()
		}

		p.SetState(502)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(503)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

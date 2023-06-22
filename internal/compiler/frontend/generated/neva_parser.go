// Code generated from ./neva.g4 by ANTLR 4.13.0. DO NOT EDIT.

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
		"", "'//'", "'use'", "'{'", "'}'", "'@/'", "'/'", "'type'", "'pub'",
		"'<'", "','", "'>'", "'['", "']'", "'|'", "'io'", "'('", "')'", "'const'",
		"'='", "'true'", "'false'", "'nil'", "':'", "'comp'", "'node'", "'.'",
		"'net'", "'->'", "'in'", "'out'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "IDENTIFIER",
		"INT", "FLOAT", "STRING", "NEWLINE", "WS",
	}
	staticData.RuleNames = []string{
		"prog", "comment", "stmt", "useStmt", "importDef", "importPath", "typeStmt",
		"typeDef", "typeParams", "typeParam", "typeExpr", "typeInstExpr", "typeArgs",
		"typeLitExpr", "enumTypeExpr", "arrTypeExpr", "recTypeExpr", "recFields",
		"recField", "unionTypeExpr", "nonUnionTypeExpr", "ioStmt", "interfaceDef",
		"inPortsDef", "outPortsDef", "portsDef", "portAndType", "constStmt",
		"constDefList", "constDef", "constValue", "arrLit", "arrItems", "recLit",
		"recValueFields", "recValueField", "compStmt", "compDefList", "compDef",
		"compBody", "compNodesDef", "compNodeDefList", "absNodeDef", "concreteNodeDef",
		"concreteNodeInst", "nodeRef", "nodeArgs", "nodeArgList", "nodeArg",
		"compNetDef", "connDefList", "connDef", "portAddr", "portDirection",
		"connReceiverSide", "connReceivers",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 36, 635, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 1, 0, 1, 0, 1, 0, 5, 0, 116,
		8, 0, 10, 0, 12, 0, 119, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 5, 1, 125, 8, 1,
		10, 1, 12, 1, 128, 9, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 135, 8, 2,
		1, 3, 1, 3, 1, 3, 5, 3, 140, 8, 3, 10, 3, 12, 3, 143, 9, 3, 1, 3, 5, 3,
		146, 8, 3, 10, 3, 12, 3, 149, 9, 3, 1, 3, 1, 3, 1, 4, 3, 4, 154, 8, 4,
		1, 4, 1, 4, 5, 4, 158, 8, 4, 10, 4, 12, 4, 161, 9, 4, 1, 5, 3, 5, 164,
		8, 5, 1, 5, 1, 5, 1, 5, 5, 5, 169, 8, 5, 10, 5, 12, 5, 172, 9, 5, 1, 6,
		1, 6, 1, 6, 5, 6, 177, 8, 6, 10, 6, 12, 6, 180, 9, 6, 1, 6, 5, 6, 183,
		8, 6, 10, 6, 12, 6, 186, 9, 6, 1, 6, 1, 6, 1, 7, 3, 7, 191, 8, 7, 1, 7,
		1, 7, 3, 7, 195, 8, 7, 1, 7, 1, 7, 5, 7, 199, 8, 7, 10, 7, 12, 7, 202,
		9, 7, 1, 8, 1, 8, 5, 8, 206, 8, 8, 10, 8, 12, 8, 209, 9, 8, 1, 8, 1, 8,
		1, 8, 5, 8, 214, 8, 8, 10, 8, 12, 8, 217, 9, 8, 1, 8, 5, 8, 220, 8, 8,
		10, 8, 12, 8, 223, 9, 8, 1, 8, 5, 8, 226, 8, 8, 10, 8, 12, 8, 229, 9, 8,
		1, 8, 1, 8, 1, 9, 1, 9, 3, 9, 235, 8, 9, 1, 10, 1, 10, 1, 10, 3, 10, 240,
		8, 10, 1, 11, 1, 11, 3, 11, 244, 8, 11, 1, 12, 1, 12, 5, 12, 248, 8, 12,
		10, 12, 12, 12, 251, 9, 12, 1, 12, 1, 12, 1, 12, 5, 12, 256, 8, 12, 10,
		12, 12, 12, 259, 9, 12, 1, 12, 5, 12, 262, 8, 12, 10, 12, 12, 12, 265,
		9, 12, 1, 12, 5, 12, 268, 8, 12, 10, 12, 12, 12, 271, 9, 12, 1, 12, 1,
		12, 1, 13, 1, 13, 1, 13, 3, 13, 278, 8, 13, 1, 14, 1, 14, 5, 14, 282, 8,
		14, 10, 14, 12, 14, 285, 9, 14, 1, 14, 1, 14, 1, 14, 5, 14, 290, 8, 14,
		10, 14, 12, 14, 293, 9, 14, 1, 14, 5, 14, 296, 8, 14, 10, 14, 12, 14, 299,
		9, 14, 1, 14, 5, 14, 302, 8, 14, 10, 14, 12, 14, 305, 9, 14, 1, 14, 1,
		14, 1, 15, 1, 15, 5, 15, 311, 8, 15, 10, 15, 12, 15, 314, 9, 15, 1, 15,
		1, 15, 5, 15, 318, 8, 15, 10, 15, 12, 15, 321, 9, 15, 1, 15, 1, 15, 1,
		15, 1, 16, 1, 16, 5, 16, 328, 8, 16, 10, 16, 12, 16, 331, 9, 16, 1, 16,
		3, 16, 334, 8, 16, 1, 16, 1, 16, 1, 17, 1, 17, 4, 17, 340, 8, 17, 11, 17,
		12, 17, 341, 1, 17, 5, 17, 345, 8, 17, 10, 17, 12, 17, 348, 9, 17, 1, 18,
		1, 18, 1, 18, 5, 18, 353, 8, 18, 10, 18, 12, 18, 356, 9, 18, 1, 19, 1,
		19, 5, 19, 360, 8, 19, 10, 19, 12, 19, 363, 9, 19, 1, 19, 1, 19, 5, 19,
		367, 8, 19, 10, 19, 12, 19, 370, 9, 19, 1, 19, 4, 19, 373, 8, 19, 11, 19,
		12, 19, 374, 1, 20, 1, 20, 3, 20, 379, 8, 20, 1, 21, 1, 21, 1, 21, 5, 21,
		384, 8, 21, 10, 21, 12, 21, 387, 9, 21, 1, 21, 5, 21, 390, 8, 21, 10, 21,
		12, 21, 393, 9, 21, 1, 21, 1, 21, 1, 22, 3, 22, 398, 8, 22, 1, 22, 1, 22,
		3, 22, 402, 8, 22, 1, 22, 1, 22, 1, 22, 5, 22, 407, 8, 22, 10, 22, 12,
		22, 410, 9, 22, 1, 23, 1, 23, 1, 24, 1, 24, 1, 25, 1, 25, 5, 25, 418, 8,
		25, 10, 25, 12, 25, 421, 9, 25, 1, 25, 3, 25, 424, 8, 25, 1, 25, 1, 25,
		1, 25, 5, 25, 429, 8, 25, 10, 25, 12, 25, 432, 9, 25, 3, 25, 434, 8, 25,
		1, 25, 1, 25, 1, 26, 5, 26, 439, 8, 26, 10, 26, 12, 26, 442, 9, 26, 1,
		26, 1, 26, 1, 26, 5, 26, 447, 8, 26, 10, 26, 12, 26, 450, 9, 26, 1, 27,
		1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 28, 1, 28, 1, 28, 5, 28, 461, 8,
		28, 10, 28, 12, 28, 464, 9, 28, 1, 29, 3, 29, 467, 8, 29, 1, 29, 1, 29,
		1, 29, 1, 29, 1, 29, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1,
		30, 3, 30, 482, 8, 30, 1, 31, 1, 31, 1, 31, 1, 31, 1, 32, 1, 32, 1, 32,
		1, 32, 3, 32, 492, 8, 32, 1, 32, 5, 32, 495, 8, 32, 10, 32, 12, 32, 498,
		9, 32, 3, 32, 500, 8, 32, 1, 33, 1, 33, 1, 33, 1, 33, 1, 34, 1, 34, 1,
		34, 3, 34, 509, 8, 34, 1, 34, 5, 34, 512, 8, 34, 10, 34, 12, 34, 515, 9,
		34, 1, 35, 1, 35, 1, 35, 1, 35, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36,
		1, 37, 1, 37, 1, 37, 5, 37, 530, 8, 37, 10, 37, 12, 37, 533, 9, 37, 1,
		38, 3, 38, 536, 8, 38, 1, 38, 1, 38, 1, 38, 1, 39, 1, 39, 1, 39, 1, 39,
		1, 39, 3, 39, 546, 8, 39, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 41, 1,
		41, 3, 41, 555, 8, 41, 1, 42, 1, 42, 1, 42, 1, 43, 1, 43, 1, 43, 1, 43,
		1, 44, 1, 44, 1, 44, 1, 44, 1, 45, 1, 45, 1, 45, 5, 45, 571, 8, 45, 10,
		45, 12, 45, 574, 9, 45, 1, 46, 1, 46, 1, 46, 1, 46, 1, 47, 1, 47, 1, 47,
		3, 47, 583, 8, 47, 1, 47, 1, 47, 1, 48, 1, 48, 1, 49, 1, 49, 1, 49, 1,
		49, 1, 49, 1, 50, 1, 50, 1, 50, 5, 50, 597, 8, 50, 10, 50, 12, 50, 600,
		9, 50, 1, 51, 1, 51, 1, 51, 1, 51, 1, 52, 3, 52, 607, 8, 52, 1, 52, 1,
		52, 1, 52, 1, 52, 1, 52, 3, 52, 614, 8, 52, 3, 52, 616, 8, 52, 1, 53, 1,
		53, 1, 54, 1, 54, 3, 54, 622, 8, 54, 1, 55, 1, 55, 1, 55, 1, 55, 5, 55,
		628, 8, 55, 10, 55, 12, 55, 631, 9, 55, 1, 55, 1, 55, 1, 55, 0, 0, 56,
		0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36,
		38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72,
		74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106,
		108, 110, 0, 2, 1, 0, 35, 35, 1, 0, 29, 30, 664, 0, 117, 1, 0, 0, 0, 2,
		122, 1, 0, 0, 0, 4, 134, 1, 0, 0, 0, 6, 136, 1, 0, 0, 0, 8, 153, 1, 0,
		0, 0, 10, 163, 1, 0, 0, 0, 12, 173, 1, 0, 0, 0, 14, 190, 1, 0, 0, 0, 16,
		203, 1, 0, 0, 0, 18, 232, 1, 0, 0, 0, 20, 239, 1, 0, 0, 0, 22, 241, 1,
		0, 0, 0, 24, 245, 1, 0, 0, 0, 26, 277, 1, 0, 0, 0, 28, 279, 1, 0, 0, 0,
		30, 308, 1, 0, 0, 0, 32, 325, 1, 0, 0, 0, 34, 337, 1, 0, 0, 0, 36, 349,
		1, 0, 0, 0, 38, 357, 1, 0, 0, 0, 40, 378, 1, 0, 0, 0, 42, 380, 1, 0, 0,
		0, 44, 397, 1, 0, 0, 0, 46, 411, 1, 0, 0, 0, 48, 413, 1, 0, 0, 0, 50, 415,
		1, 0, 0, 0, 52, 440, 1, 0, 0, 0, 54, 451, 1, 0, 0, 0, 56, 457, 1, 0, 0,
		0, 58, 466, 1, 0, 0, 0, 60, 481, 1, 0, 0, 0, 62, 483, 1, 0, 0, 0, 64, 499,
		1, 0, 0, 0, 66, 501, 1, 0, 0, 0, 68, 505, 1, 0, 0, 0, 70, 516, 1, 0, 0,
		0, 72, 520, 1, 0, 0, 0, 74, 526, 1, 0, 0, 0, 76, 535, 1, 0, 0, 0, 78, 545,
		1, 0, 0, 0, 80, 547, 1, 0, 0, 0, 82, 554, 1, 0, 0, 0, 84, 556, 1, 0, 0,
		0, 86, 559, 1, 0, 0, 0, 88, 563, 1, 0, 0, 0, 90, 567, 1, 0, 0, 0, 92, 575,
		1, 0, 0, 0, 94, 579, 1, 0, 0, 0, 96, 586, 1, 0, 0, 0, 98, 588, 1, 0, 0,
		0, 100, 593, 1, 0, 0, 0, 102, 601, 1, 0, 0, 0, 104, 615, 1, 0, 0, 0, 106,
		617, 1, 0, 0, 0, 108, 621, 1, 0, 0, 0, 110, 623, 1, 0, 0, 0, 112, 116,
		5, 35, 0, 0, 113, 116, 3, 2, 1, 0, 114, 116, 3, 4, 2, 0, 115, 112, 1, 0,
		0, 0, 115, 113, 1, 0, 0, 0, 115, 114, 1, 0, 0, 0, 116, 119, 1, 0, 0, 0,
		117, 115, 1, 0, 0, 0, 117, 118, 1, 0, 0, 0, 118, 120, 1, 0, 0, 0, 119,
		117, 1, 0, 0, 0, 120, 121, 5, 0, 0, 1, 121, 1, 1, 0, 0, 0, 122, 126, 5,
		1, 0, 0, 123, 125, 8, 0, 0, 0, 124, 123, 1, 0, 0, 0, 125, 128, 1, 0, 0,
		0, 126, 124, 1, 0, 0, 0, 126, 127, 1, 0, 0, 0, 127, 3, 1, 0, 0, 0, 128,
		126, 1, 0, 0, 0, 129, 135, 3, 6, 3, 0, 130, 135, 3, 12, 6, 0, 131, 135,
		3, 42, 21, 0, 132, 135, 3, 54, 27, 0, 133, 135, 3, 72, 36, 0, 134, 129,
		1, 0, 0, 0, 134, 130, 1, 0, 0, 0, 134, 131, 1, 0, 0, 0, 134, 132, 1, 0,
		0, 0, 134, 133, 1, 0, 0, 0, 135, 5, 1, 0, 0, 0, 136, 137, 5, 2, 0, 0, 137,
		141, 5, 3, 0, 0, 138, 140, 5, 35, 0, 0, 139, 138, 1, 0, 0, 0, 140, 143,
		1, 0, 0, 0, 141, 139, 1, 0, 0, 0, 141, 142, 1, 0, 0, 0, 142, 147, 1, 0,
		0, 0, 143, 141, 1, 0, 0, 0, 144, 146, 3, 8, 4, 0, 145, 144, 1, 0, 0, 0,
		146, 149, 1, 0, 0, 0, 147, 145, 1, 0, 0, 0, 147, 148, 1, 0, 0, 0, 148,
		150, 1, 0, 0, 0, 149, 147, 1, 0, 0, 0, 150, 151, 5, 4, 0, 0, 151, 7, 1,
		0, 0, 0, 152, 154, 5, 31, 0, 0, 153, 152, 1, 0, 0, 0, 153, 154, 1, 0, 0,
		0, 154, 155, 1, 0, 0, 0, 155, 159, 3, 10, 5, 0, 156, 158, 5, 35, 0, 0,
		157, 156, 1, 0, 0, 0, 158, 161, 1, 0, 0, 0, 159, 157, 1, 0, 0, 0, 159,
		160, 1, 0, 0, 0, 160, 9, 1, 0, 0, 0, 161, 159, 1, 0, 0, 0, 162, 164, 5,
		5, 0, 0, 163, 162, 1, 0, 0, 0, 163, 164, 1, 0, 0, 0, 164, 165, 1, 0, 0,
		0, 165, 170, 5, 31, 0, 0, 166, 167, 5, 6, 0, 0, 167, 169, 5, 31, 0, 0,
		168, 166, 1, 0, 0, 0, 169, 172, 1, 0, 0, 0, 170, 168, 1, 0, 0, 0, 170,
		171, 1, 0, 0, 0, 171, 11, 1, 0, 0, 0, 172, 170, 1, 0, 0, 0, 173, 174, 5,
		7, 0, 0, 174, 178, 5, 3, 0, 0, 175, 177, 5, 35, 0, 0, 176, 175, 1, 0, 0,
		0, 177, 180, 1, 0, 0, 0, 178, 176, 1, 0, 0, 0, 178, 179, 1, 0, 0, 0, 179,
		184, 1, 0, 0, 0, 180, 178, 1, 0, 0, 0, 181, 183, 3, 14, 7, 0, 182, 181,
		1, 0, 0, 0, 183, 186, 1, 0, 0, 0, 184, 182, 1, 0, 0, 0, 184, 185, 1, 0,
		0, 0, 185, 187, 1, 0, 0, 0, 186, 184, 1, 0, 0, 0, 187, 188, 5, 4, 0, 0,
		188, 13, 1, 0, 0, 0, 189, 191, 5, 8, 0, 0, 190, 189, 1, 0, 0, 0, 190, 191,
		1, 0, 0, 0, 191, 192, 1, 0, 0, 0, 192, 194, 5, 31, 0, 0, 193, 195, 3, 16,
		8, 0, 194, 193, 1, 0, 0, 0, 194, 195, 1, 0, 0, 0, 195, 196, 1, 0, 0, 0,
		196, 200, 3, 20, 10, 0, 197, 199, 5, 35, 0, 0, 198, 197, 1, 0, 0, 0, 199,
		202, 1, 0, 0, 0, 200, 198, 1, 0, 0, 0, 200, 201, 1, 0, 0, 0, 201, 15, 1,
		0, 0, 0, 202, 200, 1, 0, 0, 0, 203, 207, 5, 9, 0, 0, 204, 206, 5, 35, 0,
		0, 205, 204, 1, 0, 0, 0, 206, 209, 1, 0, 0, 0, 207, 205, 1, 0, 0, 0, 207,
		208, 1, 0, 0, 0, 208, 210, 1, 0, 0, 0, 209, 207, 1, 0, 0, 0, 210, 221,
		3, 18, 9, 0, 211, 215, 5, 10, 0, 0, 212, 214, 5, 35, 0, 0, 213, 212, 1,
		0, 0, 0, 214, 217, 1, 0, 0, 0, 215, 213, 1, 0, 0, 0, 215, 216, 1, 0, 0,
		0, 216, 218, 1, 0, 0, 0, 217, 215, 1, 0, 0, 0, 218, 220, 3, 18, 9, 0, 219,
		211, 1, 0, 0, 0, 220, 223, 1, 0, 0, 0, 221, 219, 1, 0, 0, 0, 221, 222,
		1, 0, 0, 0, 222, 227, 1, 0, 0, 0, 223, 221, 1, 0, 0, 0, 224, 226, 5, 35,
		0, 0, 225, 224, 1, 0, 0, 0, 226, 229, 1, 0, 0, 0, 227, 225, 1, 0, 0, 0,
		227, 228, 1, 0, 0, 0, 228, 230, 1, 0, 0, 0, 229, 227, 1, 0, 0, 0, 230,
		231, 5, 11, 0, 0, 231, 17, 1, 0, 0, 0, 232, 234, 5, 31, 0, 0, 233, 235,
		3, 20, 10, 0, 234, 233, 1, 0, 0, 0, 234, 235, 1, 0, 0, 0, 235, 19, 1, 0,
		0, 0, 236, 240, 3, 22, 11, 0, 237, 240, 3, 26, 13, 0, 238, 240, 3, 38,
		19, 0, 239, 236, 1, 0, 0, 0, 239, 237, 1, 0, 0, 0, 239, 238, 1, 0, 0, 0,
		240, 21, 1, 0, 0, 0, 241, 243, 5, 31, 0, 0, 242, 244, 3, 24, 12, 0, 243,
		242, 1, 0, 0, 0, 243, 244, 1, 0, 0, 0, 244, 23, 1, 0, 0, 0, 245, 249, 5,
		9, 0, 0, 246, 248, 5, 35, 0, 0, 247, 246, 1, 0, 0, 0, 248, 251, 1, 0, 0,
		0, 249, 247, 1, 0, 0, 0, 249, 250, 1, 0, 0, 0, 250, 252, 1, 0, 0, 0, 251,
		249, 1, 0, 0, 0, 252, 263, 3, 20, 10, 0, 253, 257, 5, 10, 0, 0, 254, 256,
		5, 35, 0, 0, 255, 254, 1, 0, 0, 0, 256, 259, 1, 0, 0, 0, 257, 255, 1, 0,
		0, 0, 257, 258, 1, 0, 0, 0, 258, 260, 1, 0, 0, 0, 259, 257, 1, 0, 0, 0,
		260, 262, 3, 20, 10, 0, 261, 253, 1, 0, 0, 0, 262, 265, 1, 0, 0, 0, 263,
		261, 1, 0, 0, 0, 263, 264, 1, 0, 0, 0, 264, 269, 1, 0, 0, 0, 265, 263,
		1, 0, 0, 0, 266, 268, 5, 35, 0, 0, 267, 266, 1, 0, 0, 0, 268, 271, 1, 0,
		0, 0, 269, 267, 1, 0, 0, 0, 269, 270, 1, 0, 0, 0, 270, 272, 1, 0, 0, 0,
		271, 269, 1, 0, 0, 0, 272, 273, 5, 11, 0, 0, 273, 25, 1, 0, 0, 0, 274,
		278, 3, 28, 14, 0, 275, 278, 3, 30, 15, 0, 276, 278, 3, 32, 16, 0, 277,
		274, 1, 0, 0, 0, 277, 275, 1, 0, 0, 0, 277, 276, 1, 0, 0, 0, 278, 27, 1,
		0, 0, 0, 279, 283, 5, 3, 0, 0, 280, 282, 5, 35, 0, 0, 281, 280, 1, 0, 0,
		0, 282, 285, 1, 0, 0, 0, 283, 281, 1, 0, 0, 0, 283, 284, 1, 0, 0, 0, 284,
		286, 1, 0, 0, 0, 285, 283, 1, 0, 0, 0, 286, 297, 5, 31, 0, 0, 287, 291,
		5, 10, 0, 0, 288, 290, 5, 35, 0, 0, 289, 288, 1, 0, 0, 0, 290, 293, 1,
		0, 0, 0, 291, 289, 1, 0, 0, 0, 291, 292, 1, 0, 0, 0, 292, 294, 1, 0, 0,
		0, 293, 291, 1, 0, 0, 0, 294, 296, 5, 31, 0, 0, 295, 287, 1, 0, 0, 0, 296,
		299, 1, 0, 0, 0, 297, 295, 1, 0, 0, 0, 297, 298, 1, 0, 0, 0, 298, 303,
		1, 0, 0, 0, 299, 297, 1, 0, 0, 0, 300, 302, 5, 35, 0, 0, 301, 300, 1, 0,
		0, 0, 302, 305, 1, 0, 0, 0, 303, 301, 1, 0, 0, 0, 303, 304, 1, 0, 0, 0,
		304, 306, 1, 0, 0, 0, 305, 303, 1, 0, 0, 0, 306, 307, 5, 4, 0, 0, 307,
		29, 1, 0, 0, 0, 308, 312, 5, 12, 0, 0, 309, 311, 5, 35, 0, 0, 310, 309,
		1, 0, 0, 0, 311, 314, 1, 0, 0, 0, 312, 310, 1, 0, 0, 0, 312, 313, 1, 0,
		0, 0, 313, 315, 1, 0, 0, 0, 314, 312, 1, 0, 0, 0, 315, 319, 5, 32, 0, 0,
		316, 318, 5, 35, 0, 0, 317, 316, 1, 0, 0, 0, 318, 321, 1, 0, 0, 0, 319,
		317, 1, 0, 0, 0, 319, 320, 1, 0, 0, 0, 320, 322, 1, 0, 0, 0, 321, 319,
		1, 0, 0, 0, 322, 323, 5, 13, 0, 0, 323, 324, 3, 20, 10, 0, 324, 31, 1,
		0, 0, 0, 325, 329, 5, 3, 0, 0, 326, 328, 5, 35, 0, 0, 327, 326, 1, 0, 0,
		0, 328, 331, 1, 0, 0, 0, 329, 327, 1, 0, 0, 0, 329, 330, 1, 0, 0, 0, 330,
		333, 1, 0, 0, 0, 331, 329, 1, 0, 0, 0, 332, 334, 3, 34, 17, 0, 333, 332,
		1, 0, 0, 0, 333, 334, 1, 0, 0, 0, 334, 335, 1, 0, 0, 0, 335, 336, 5, 4,
		0, 0, 336, 33, 1, 0, 0, 0, 337, 346, 3, 36, 18, 0, 338, 340, 5, 35, 0,
		0, 339, 338, 1, 0, 0, 0, 340, 341, 1, 0, 0, 0, 341, 339, 1, 0, 0, 0, 341,
		342, 1, 0, 0, 0, 342, 343, 1, 0, 0, 0, 343, 345, 3, 36, 18, 0, 344, 339,
		1, 0, 0, 0, 345, 348, 1, 0, 0, 0, 346, 344, 1, 0, 0, 0, 346, 347, 1, 0,
		0, 0, 347, 35, 1, 0, 0, 0, 348, 346, 1, 0, 0, 0, 349, 350, 5, 31, 0, 0,
		350, 354, 3, 20, 10, 0, 351, 353, 5, 35, 0, 0, 352, 351, 1, 0, 0, 0, 353,
		356, 1, 0, 0, 0, 354, 352, 1, 0, 0, 0, 354, 355, 1, 0, 0, 0, 355, 37, 1,
		0, 0, 0, 356, 354, 1, 0, 0, 0, 357, 372, 3, 40, 20, 0, 358, 360, 5, 35,
		0, 0, 359, 358, 1, 0, 0, 0, 360, 363, 1, 0, 0, 0, 361, 359, 1, 0, 0, 0,
		361, 362, 1, 0, 0, 0, 362, 364, 1, 0, 0, 0, 363, 361, 1, 0, 0, 0, 364,
		368, 5, 14, 0, 0, 365, 367, 5, 35, 0, 0, 366, 365, 1, 0, 0, 0, 367, 370,
		1, 0, 0, 0, 368, 366, 1, 0, 0, 0, 368, 369, 1, 0, 0, 0, 369, 371, 1, 0,
		0, 0, 370, 368, 1, 0, 0, 0, 371, 373, 3, 40, 20, 0, 372, 361, 1, 0, 0,
		0, 373, 374, 1, 0, 0, 0, 374, 372, 1, 0, 0, 0, 374, 375, 1, 0, 0, 0, 375,
		39, 1, 0, 0, 0, 376, 379, 3, 22, 11, 0, 377, 379, 3, 26, 13, 0, 378, 376,
		1, 0, 0, 0, 378, 377, 1, 0, 0, 0, 379, 41, 1, 0, 0, 0, 380, 381, 5, 15,
		0, 0, 381, 385, 5, 3, 0, 0, 382, 384, 5, 35, 0, 0, 383, 382, 1, 0, 0, 0,
		384, 387, 1, 0, 0, 0, 385, 383, 1, 0, 0, 0, 385, 386, 1, 0, 0, 0, 386,
		391, 1, 0, 0, 0, 387, 385, 1, 0, 0, 0, 388, 390, 3, 44, 22, 0, 389, 388,
		1, 0, 0, 0, 390, 393, 1, 0, 0, 0, 391, 389, 1, 0, 0, 0, 391, 392, 1, 0,
		0, 0, 392, 394, 1, 0, 0, 0, 393, 391, 1, 0, 0, 0, 394, 395, 5, 4, 0, 0,
		395, 43, 1, 0, 0, 0, 396, 398, 5, 8, 0, 0, 397, 396, 1, 0, 0, 0, 397, 398,
		1, 0, 0, 0, 398, 399, 1, 0, 0, 0, 399, 401, 5, 31, 0, 0, 400, 402, 3, 16,
		8, 0, 401, 400, 1, 0, 0, 0, 401, 402, 1, 0, 0, 0, 402, 403, 1, 0, 0, 0,
		403, 404, 3, 46, 23, 0, 404, 408, 3, 48, 24, 0, 405, 407, 5, 35, 0, 0,
		406, 405, 1, 0, 0, 0, 407, 410, 1, 0, 0, 0, 408, 406, 1, 0, 0, 0, 408,
		409, 1, 0, 0, 0, 409, 45, 1, 0, 0, 0, 410, 408, 1, 0, 0, 0, 411, 412, 3,
		50, 25, 0, 412, 47, 1, 0, 0, 0, 413, 414, 3, 50, 25, 0, 414, 49, 1, 0,
		0, 0, 415, 433, 5, 16, 0, 0, 416, 418, 5, 35, 0, 0, 417, 416, 1, 0, 0,
		0, 418, 421, 1, 0, 0, 0, 419, 417, 1, 0, 0, 0, 419, 420, 1, 0, 0, 0, 420,
		434, 1, 0, 0, 0, 421, 419, 1, 0, 0, 0, 422, 424, 3, 52, 26, 0, 423, 422,
		1, 0, 0, 0, 423, 424, 1, 0, 0, 0, 424, 434, 1, 0, 0, 0, 425, 430, 3, 52,
		26, 0, 426, 427, 5, 10, 0, 0, 427, 429, 3, 52, 26, 0, 428, 426, 1, 0, 0,
		0, 429, 432, 1, 0, 0, 0, 430, 428, 1, 0, 0, 0, 430, 431, 1, 0, 0, 0, 431,
		434, 1, 0, 0, 0, 432, 430, 1, 0, 0, 0, 433, 419, 1, 0, 0, 0, 433, 423,
		1, 0, 0, 0, 433, 425, 1, 0, 0, 0, 434, 435, 1, 0, 0, 0, 435, 436, 5, 17,
		0, 0, 436, 51, 1, 0, 0, 0, 437, 439, 5, 35, 0, 0, 438, 437, 1, 0, 0, 0,
		439, 442, 1, 0, 0, 0, 440, 438, 1, 0, 0, 0, 440, 441, 1, 0, 0, 0, 441,
		443, 1, 0, 0, 0, 442, 440, 1, 0, 0, 0, 443, 444, 5, 31, 0, 0, 444, 448,
		3, 20, 10, 0, 445, 447, 5, 35, 0, 0, 446, 445, 1, 0, 0, 0, 447, 450, 1,
		0, 0, 0, 448, 446, 1, 0, 0, 0, 448, 449, 1, 0, 0, 0, 449, 53, 1, 0, 0,
		0, 450, 448, 1, 0, 0, 0, 451, 452, 5, 18, 0, 0, 452, 453, 5, 3, 0, 0, 453,
		454, 3, 56, 28, 0, 454, 455, 5, 4, 0, 0, 455, 456, 5, 35, 0, 0, 456, 55,
		1, 0, 0, 0, 457, 462, 3, 58, 29, 0, 458, 459, 5, 35, 0, 0, 459, 461, 3,
		58, 29, 0, 460, 458, 1, 0, 0, 0, 461, 464, 1, 0, 0, 0, 462, 460, 1, 0,
		0, 0, 462, 463, 1, 0, 0, 0, 463, 57, 1, 0, 0, 0, 464, 462, 1, 0, 0, 0,
		465, 467, 5, 8, 0, 0, 466, 465, 1, 0, 0, 0, 466, 467, 1, 0, 0, 0, 467,
		468, 1, 0, 0, 0, 468, 469, 5, 31, 0, 0, 469, 470, 3, 20, 10, 0, 470, 471,
		5, 19, 0, 0, 471, 472, 3, 60, 30, 0, 472, 59, 1, 0, 0, 0, 473, 482, 5,
		20, 0, 0, 474, 482, 5, 21, 0, 0, 475, 482, 5, 32, 0, 0, 476, 482, 5, 33,
		0, 0, 477, 482, 5, 34, 0, 0, 478, 482, 3, 62, 31, 0, 479, 482, 3, 66, 33,
		0, 480, 482, 5, 22, 0, 0, 481, 473, 1, 0, 0, 0, 481, 474, 1, 0, 0, 0, 481,
		475, 1, 0, 0, 0, 481, 476, 1, 0, 0, 0, 481, 477, 1, 0, 0, 0, 481, 478,
		1, 0, 0, 0, 481, 479, 1, 0, 0, 0, 481, 480, 1, 0, 0, 0, 482, 61, 1, 0,
		0, 0, 483, 484, 5, 12, 0, 0, 484, 485, 3, 64, 32, 0, 485, 486, 5, 13, 0,
		0, 486, 63, 1, 0, 0, 0, 487, 500, 3, 60, 30, 0, 488, 496, 3, 60, 30, 0,
		489, 491, 5, 10, 0, 0, 490, 492, 5, 35, 0, 0, 491, 490, 1, 0, 0, 0, 491,
		492, 1, 0, 0, 0, 492, 493, 1, 0, 0, 0, 493, 495, 3, 60, 30, 0, 494, 489,
		1, 0, 0, 0, 495, 498, 1, 0, 0, 0, 496, 494, 1, 0, 0, 0, 496, 497, 1, 0,
		0, 0, 497, 500, 1, 0, 0, 0, 498, 496, 1, 0, 0, 0, 499, 487, 1, 0, 0, 0,
		499, 488, 1, 0, 0, 0, 500, 65, 1, 0, 0, 0, 501, 502, 5, 3, 0, 0, 502, 503,
		3, 68, 34, 0, 503, 504, 5, 4, 0, 0, 504, 67, 1, 0, 0, 0, 505, 513, 3, 70,
		35, 0, 506, 508, 5, 10, 0, 0, 507, 509, 5, 35, 0, 0, 508, 507, 1, 0, 0,
		0, 508, 509, 1, 0, 0, 0, 509, 510, 1, 0, 0, 0, 510, 512, 3, 70, 35, 0,
		511, 506, 1, 0, 0, 0, 512, 515, 1, 0, 0, 0, 513, 511, 1, 0, 0, 0, 513,
		514, 1, 0, 0, 0, 514, 69, 1, 0, 0, 0, 515, 513, 1, 0, 0, 0, 516, 517, 5,
		31, 0, 0, 517, 518, 5, 23, 0, 0, 518, 519, 3, 60, 30, 0, 519, 71, 1, 0,
		0, 0, 520, 521, 5, 24, 0, 0, 521, 522, 5, 3, 0, 0, 522, 523, 3, 74, 37,
		0, 523, 524, 5, 4, 0, 0, 524, 525, 5, 35, 0, 0, 525, 73, 1, 0, 0, 0, 526,
		531, 3, 76, 38, 0, 527, 528, 5, 35, 0, 0, 528, 530, 3, 76, 38, 0, 529,
		527, 1, 0, 0, 0, 530, 533, 1, 0, 0, 0, 531, 529, 1, 0, 0, 0, 531, 532,
		1, 0, 0, 0, 532, 75, 1, 0, 0, 0, 533, 531, 1, 0, 0, 0, 534, 536, 5, 8,
		0, 0, 535, 534, 1, 0, 0, 0, 535, 536, 1, 0, 0, 0, 536, 537, 1, 0, 0, 0,
		537, 538, 3, 44, 22, 0, 538, 539, 3, 78, 39, 0, 539, 77, 1, 0, 0, 0, 540,
		541, 5, 3, 0, 0, 541, 546, 3, 80, 40, 0, 542, 543, 3, 98, 49, 0, 543, 544,
		5, 4, 0, 0, 544, 546, 1, 0, 0, 0, 545, 540, 1, 0, 0, 0, 545, 542, 1, 0,
		0, 0, 546, 79, 1, 0, 0, 0, 547, 548, 5, 25, 0, 0, 548, 549, 5, 3, 0, 0,
		549, 550, 3, 82, 41, 0, 550, 551, 5, 4, 0, 0, 551, 81, 1, 0, 0, 0, 552,
		555, 3, 84, 42, 0, 553, 555, 3, 86, 43, 0, 554, 552, 1, 0, 0, 0, 554, 553,
		1, 0, 0, 0, 555, 83, 1, 0, 0, 0, 556, 557, 5, 31, 0, 0, 557, 558, 3, 22,
		11, 0, 558, 85, 1, 0, 0, 0, 559, 560, 5, 31, 0, 0, 560, 561, 5, 19, 0,
		0, 561, 562, 3, 88, 44, 0, 562, 87, 1, 0, 0, 0, 563, 564, 3, 90, 45, 0,
		564, 565, 3, 92, 46, 0, 565, 566, 3, 24, 12, 0, 566, 89, 1, 0, 0, 0, 567,
		572, 5, 31, 0, 0, 568, 569, 5, 26, 0, 0, 569, 571, 5, 31, 0, 0, 570, 568,
		1, 0, 0, 0, 571, 574, 1, 0, 0, 0, 572, 570, 1, 0, 0, 0, 572, 573, 1, 0,
		0, 0, 573, 91, 1, 0, 0, 0, 574, 572, 1, 0, 0, 0, 575, 576, 5, 16, 0, 0,
		576, 577, 3, 94, 47, 0, 577, 578, 5, 17, 0, 0, 578, 93, 1, 0, 0, 0, 579,
		580, 3, 96, 48, 0, 580, 582, 5, 10, 0, 0, 581, 583, 5, 35, 0, 0, 582, 581,
		1, 0, 0, 0, 582, 583, 1, 0, 0, 0, 583, 584, 1, 0, 0, 0, 584, 585, 3, 96,
		48, 0, 585, 95, 1, 0, 0, 0, 586, 587, 3, 88, 44, 0, 587, 97, 1, 0, 0, 0,
		588, 589, 5, 27, 0, 0, 589, 590, 5, 3, 0, 0, 590, 591, 3, 100, 50, 0, 591,
		592, 5, 4, 0, 0, 592, 99, 1, 0, 0, 0, 593, 598, 3, 102, 51, 0, 594, 595,
		5, 35, 0, 0, 595, 597, 3, 102, 51, 0, 596, 594, 1, 0, 0, 0, 597, 600, 1,
		0, 0, 0, 598, 596, 1, 0, 0, 0, 598, 599, 1, 0, 0, 0, 599, 101, 1, 0, 0,
		0, 600, 598, 1, 0, 0, 0, 601, 602, 3, 104, 52, 0, 602, 603, 5, 28, 0, 0,
		603, 604, 3, 108, 54, 0, 604, 103, 1, 0, 0, 0, 605, 607, 5, 31, 0, 0, 606,
		605, 1, 0, 0, 0, 606, 607, 1, 0, 0, 0, 607, 608, 1, 0, 0, 0, 608, 616,
		3, 106, 53, 0, 609, 613, 5, 31, 0, 0, 610, 611, 5, 12, 0, 0, 611, 612,
		5, 32, 0, 0, 612, 614, 5, 13, 0, 0, 613, 610, 1, 0, 0, 0, 613, 614, 1,
		0, 0, 0, 614, 616, 1, 0, 0, 0, 615, 606, 1, 0, 0, 0, 615, 609, 1, 0, 0,
		0, 616, 105, 1, 0, 0, 0, 617, 618, 7, 1, 0, 0, 618, 107, 1, 0, 0, 0, 619,
		622, 3, 104, 52, 0, 620, 622, 3, 110, 55, 0, 621, 619, 1, 0, 0, 0, 621,
		620, 1, 0, 0, 0, 622, 109, 1, 0, 0, 0, 623, 624, 5, 3, 0, 0, 624, 629,
		3, 104, 52, 0, 625, 626, 5, 35, 0, 0, 626, 628, 3, 104, 52, 0, 627, 625,
		1, 0, 0, 0, 628, 631, 1, 0, 0, 0, 629, 627, 1, 0, 0, 0, 629, 630, 1, 0,
		0, 0, 630, 632, 1, 0, 0, 0, 631, 629, 1, 0, 0, 0, 632, 633, 5, 4, 0, 0,
		633, 111, 1, 0, 0, 0, 73, 115, 117, 126, 134, 141, 147, 153, 159, 163,
		170, 178, 184, 190, 194, 200, 207, 215, 221, 227, 234, 239, 243, 249, 257,
		263, 269, 277, 283, 291, 297, 303, 312, 319, 329, 333, 341, 346, 354, 361,
		368, 374, 378, 385, 391, 397, 401, 408, 419, 423, 430, 433, 440, 448, 462,
		466, 481, 491, 496, 499, 508, 513, 531, 535, 545, 554, 572, 582, 598, 606,
		613, 615, 621, 629,
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
	nevaParserIDENTIFIER = 31
	nevaParserINT        = 32
	nevaParserFLOAT      = 33
	nevaParserSTRING     = 34
	nevaParserNEWLINE    = 35
	nevaParserWS         = 36
)

// nevaParser rules.
const (
	nevaParserRULE_prog             = 0
	nevaParserRULE_comment          = 1
	nevaParserRULE_stmt             = 2
	nevaParserRULE_useStmt          = 3
	nevaParserRULE_importDef        = 4
	nevaParserRULE_importPath       = 5
	nevaParserRULE_typeStmt         = 6
	nevaParserRULE_typeDef          = 7
	nevaParserRULE_typeParams       = 8
	nevaParserRULE_typeParam        = 9
	nevaParserRULE_typeExpr         = 10
	nevaParserRULE_typeInstExpr     = 11
	nevaParserRULE_typeArgs         = 12
	nevaParserRULE_typeLitExpr      = 13
	nevaParserRULE_enumTypeExpr     = 14
	nevaParserRULE_arrTypeExpr      = 15
	nevaParserRULE_recTypeExpr      = 16
	nevaParserRULE_recFields        = 17
	nevaParserRULE_recField         = 18
	nevaParserRULE_unionTypeExpr    = 19
	nevaParserRULE_nonUnionTypeExpr = 20
	nevaParserRULE_ioStmt           = 21
	nevaParserRULE_interfaceDef     = 22
	nevaParserRULE_inPortsDef       = 23
	nevaParserRULE_outPortsDef      = 24
	nevaParserRULE_portsDef         = 25
	nevaParserRULE_portAndType      = 26
	nevaParserRULE_constStmt        = 27
	nevaParserRULE_constDefList     = 28
	nevaParserRULE_constDef         = 29
	nevaParserRULE_constValue       = 30
	nevaParserRULE_arrLit           = 31
	nevaParserRULE_arrItems         = 32
	nevaParserRULE_recLit           = 33
	nevaParserRULE_recValueFields   = 34
	nevaParserRULE_recValueField    = 35
	nevaParserRULE_compStmt         = 36
	nevaParserRULE_compDefList      = 37
	nevaParserRULE_compDef          = 38
	nevaParserRULE_compBody         = 39
	nevaParserRULE_compNodesDef     = 40
	nevaParserRULE_compNodeDefList  = 41
	nevaParserRULE_absNodeDef       = 42
	nevaParserRULE_concreteNodeDef  = 43
	nevaParserRULE_concreteNodeInst = 44
	nevaParserRULE_nodeRef          = 45
	nevaParserRULE_nodeArgs         = 46
	nevaParserRULE_nodeArgList      = 47
	nevaParserRULE_nodeArg          = 48
	nevaParserRULE_compNetDef       = 49
	nevaParserRULE_connDefList      = 50
	nevaParserRULE_connDef          = 51
	nevaParserRULE_portAddr         = 52
	nevaParserRULE_portDirection    = 53
	nevaParserRULE_connReceiverSide = 54
	nevaParserRULE_connReceivers    = 55
)

// IProgContext is an interface to support dynamic dispatch.
type IProgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
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

func (s *ProgContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ProgContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.SetState(117)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&34376810630) != 0 {
		p.SetState(115)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case nevaParserNEWLINE:
			{
				p.SetState(112)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case nevaParserT__0:
			{
				p.SetState(113)
				p.Comment()
			}

		case nevaParserT__1, nevaParserT__6, nevaParserT__14, nevaParserT__17, nevaParserT__23:
			{
				p.SetState(114)
				p.Stmt()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(119)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(120)
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *CommentContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *CommentContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(122)
		p.Match(nevaParserT__0)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(123)
				_la = p.GetTokenStream().LA(1)

				if _la <= 0 || _la == nevaParserNEWLINE {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}

		}
		p.SetState(128)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext())
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
	CompStmt() ICompStmtContext

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

func (s *StmtContext) CompStmt() ICompStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompStmtContext)
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
	p.EnterRule(localctx, 4, nevaParserRULE_stmt)
	p.SetState(134)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(129)
			p.UseStmt()
		}

	case nevaParserT__6:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(130)
			p.TypeStmt()
		}

	case nevaParserT__14:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(131)
			p.IoStmt()
		}

	case nevaParserT__17:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(132)
			p.ConstStmt()
		}

	case nevaParserT__23:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(133)
			p.CompStmt()
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
	AllImportDef() []IImportDefContext
	ImportDef(i int) IImportDefContext

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

func (s *UseStmtContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *UseStmtContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *UseStmtContext) AllImportDef() []IImportDefContext {
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

func (s *UseStmtContext) ImportDef(i int) IImportDefContext {
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
	p.EnterRule(localctx, 6, nevaParserRULE_useStmt)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(136)
		p.Match(nevaParserT__1)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(137)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(141)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(138)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(143)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(147)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__4 || _la == nevaParserIDENTIFIER {
		{
			p.SetState(144)
			p.ImportDef()
		}

		p.SetState(149)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(150)
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

// IImportDefContext is an interface to support dynamic dispatch.
type IImportDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ImportPath() IImportPathContext
	IDENTIFIER() antlr.TerminalNode
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *ImportDefContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ImportDefContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.EnterRule(localctx, 8, nevaParserRULE_importDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(153)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(152)
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
		p.SetState(155)
		p.ImportPath()
	}
	p.SetState(159)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(156)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(161)
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
	p.EnterRule(localctx, 10, nevaParserRULE_importPath)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(163)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__4 {
		{
			p.SetState(162)
			p.Match(nevaParserT__4)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(165)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__5 {
		{
			p.SetState(166)
			p.Match(nevaParserT__5)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(167)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(172)
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
	AllTypeDef() []ITypeDefContext
	TypeDef(i int) ITypeDefContext

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

func (s *TypeStmtContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *TypeStmtContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *TypeStmtContext) AllTypeDef() []ITypeDefContext {
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

func (s *TypeStmtContext) TypeDef(i int) ITypeDefContext {
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
	p.EnterRule(localctx, 12, nevaParserRULE_typeStmt)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(173)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(174)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(178)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(175)
			p.Match(nevaParserNEWLINE)
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
	}
	p.SetState(184)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 || _la == nevaParserIDENTIFIER {
		{
			p.SetState(181)
			p.TypeDef()
		}

		p.SetState(186)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(187)
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

// ITypeDefContext is an interface to support dynamic dispatch.
type ITypeDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeExpr() ITypeExprContext
	TypeParams() ITypeParamsContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *TypeDefContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *TypeDefContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.EnterRule(localctx, 14, nevaParserRULE_typeDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(190)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__7 {
		{
			p.SetState(189)
			p.Match(nevaParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(192)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(194)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__8 {
		{
			p.SetState(193)
			p.TypeParams()
		}

	}
	{
		p.SetState(196)
		p.TypeExpr()
	}
	p.SetState(200)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(197)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(202)
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
	p.EnterRule(localctx, 16, nevaParserRULE_typeParams)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(203)
		p.Match(nevaParserT__8)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(207)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(204)
			p.Match(nevaParserNEWLINE)
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
	}
	{
		p.SetState(210)
		p.TypeParam()
	}
	p.SetState(221)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__9 {
		{
			p.SetState(211)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(215)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(212)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(217)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(218)
			p.TypeParam()
		}

		p.SetState(223)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(227)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(224)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(229)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(230)
		p.Match(nevaParserT__10)
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
	p.EnterRule(localctx, 18, nevaParserRULE_typeParam)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(232)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(234)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2147487752) != 0 {
		{
			p.SetState(233)
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
	p.EnterRule(localctx, 20, nevaParserRULE_typeExpr)
	p.SetState(239)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(236)
			p.TypeInstExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(237)
			p.TypeLitExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(238)
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
	p.EnterRule(localctx, 22, nevaParserRULE_typeInstExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(241)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(243)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__8 {
		{
			p.SetState(242)
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *TypeArgsContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *TypeArgsContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.EnterRule(localctx, 24, nevaParserRULE_typeArgs)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(245)
		p.Match(nevaParserT__8)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(249)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(246)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(251)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(252)
		p.TypeExpr()
	}
	p.SetState(263)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__9 {
		{
			p.SetState(253)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(257)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(254)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(259)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(260)
			p.TypeExpr()
		}

		p.SetState(265)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(269)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(266)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(271)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(272)
		p.Match(nevaParserT__10)
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
	EnumTypeExpr() IEnumTypeExprContext
	ArrTypeExpr() IArrTypeExprContext
	RecTypeExpr() IRecTypeExprContext

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
	p.EnterRule(localctx, 26, nevaParserRULE_typeLitExpr)
	p.SetState(277)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(274)
			p.EnumTypeExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(275)
			p.ArrTypeExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(276)
			p.RecTypeExpr()
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

// IEnumTypeExprContext is an interface to support dynamic dispatch.
type IEnumTypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *EnumTypeExprContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(nevaParserIDENTIFIER)
}

func (s *EnumTypeExprContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, i)
}

func (s *EnumTypeExprContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *EnumTypeExprContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

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
	p.EnterRule(localctx, 28, nevaParserRULE_enumTypeExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(279)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(283)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(280)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(285)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(286)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(297)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__9 {
		{
			p.SetState(287)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(291)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(288)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(293)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(294)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(299)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(303)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(300)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(305)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(306)
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

// IArrTypeExprContext is an interface to support dynamic dispatch.
type IArrTypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INT() antlr.TerminalNode
	TypeExpr() ITypeExprContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *ArrTypeExprContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ArrTypeExprContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.EnterRule(localctx, 30, nevaParserRULE_arrTypeExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(308)
		p.Match(nevaParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(312)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(309)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(314)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(315)
		p.Match(nevaParserINT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(319)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(316)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(321)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(322)
		p.Match(nevaParserT__12)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(323)
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
	RecFields() IRecFieldsContext

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

func (s *RecTypeExprContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *RecTypeExprContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *RecTypeExprContext) RecFields() IRecFieldsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRecFieldsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRecFieldsContext)
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
	p.EnterRule(localctx, 32, nevaParserRULE_recTypeExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(325)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(329)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(326)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(331)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(333)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserIDENTIFIER {
		{
			p.SetState(332)
			p.RecFields()
		}

	}
	{
		p.SetState(335)
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

// IRecFieldsContext is an interface to support dynamic dispatch.
type IRecFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRecField() []IRecFieldContext
	RecField(i int) IRecFieldContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsRecFieldsContext differentiates from other interfaces.
	IsRecFieldsContext()
}

type RecFieldsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRecFieldsContext() *RecFieldsContext {
	var p = new(RecFieldsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recFields
	return p
}

func InitEmptyRecFieldsContext(p *RecFieldsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recFields
}

func (*RecFieldsContext) IsRecFieldsContext() {}

func NewRecFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RecFieldsContext {
	var p = new(RecFieldsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_recFields

	return p
}

func (s *RecFieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *RecFieldsContext) AllRecField() []IRecFieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRecFieldContext); ok {
			len++
		}
	}

	tst := make([]IRecFieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRecFieldContext); ok {
			tst[i] = t.(IRecFieldContext)
			i++
		}
	}

	return tst
}

func (s *RecFieldsContext) RecField(i int) IRecFieldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRecFieldContext); ok {
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

	return t.(IRecFieldContext)
}

func (s *RecFieldsContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *RecFieldsContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *RecFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RecFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterRecFields(s)
	}
}

func (s *RecFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitRecFields(s)
	}
}

func (p *nevaParser) RecFields() (localctx IRecFieldsContext) {
	localctx = NewRecFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, nevaParserRULE_recFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(337)
		p.RecField()
	}
	p.SetState(346)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		p.SetState(339)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == nevaParserNEWLINE {
			{
				p.SetState(338)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(341)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(343)
			p.RecField()
		}

		p.SetState(348)
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

// IRecFieldContext is an interface to support dynamic dispatch.
type IRecFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeExpr() ITypeExprContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsRecFieldContext differentiates from other interfaces.
	IsRecFieldContext()
}

type RecFieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRecFieldContext() *RecFieldContext {
	var p = new(RecFieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recField
	return p
}

func InitEmptyRecFieldContext(p *RecFieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_recField
}

func (*RecFieldContext) IsRecFieldContext() {}

func NewRecFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RecFieldContext {
	var p = new(RecFieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_recField

	return p
}

func (s *RecFieldContext) GetParser() antlr.Parser { return s.parser }

func (s *RecFieldContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *RecFieldContext) TypeExpr() ITypeExprContext {
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

func (s *RecFieldContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *RecFieldContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *RecFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RecFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterRecField(s)
	}
}

func (s *RecFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitRecField(s)
	}
}

func (p *nevaParser) RecField() (localctx IRecFieldContext) {
	localctx = NewRecFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, nevaParserRULE_recField)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(349)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(350)
		p.TypeExpr()
	}
	p.SetState(354)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(351)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		p.SetState(356)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
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

// IUnionTypeExprContext is an interface to support dynamic dispatch.
type IUnionTypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNonUnionTypeExpr() []INonUnionTypeExprContext
	NonUnionTypeExpr(i int) INonUnionTypeExprContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *UnionTypeExprContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *UnionTypeExprContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.EnterRule(localctx, 38, nevaParserRULE_unionTypeExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(357)
		p.NonUnionTypeExpr()
	}
	p.SetState(372)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			p.SetState(361)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(358)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(363)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(364)
				p.Match(nevaParserT__13)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(368)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(365)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(370)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(371)
				p.NonUnionTypeExpr()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(374)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 40, nevaParserRULE_nonUnionTypeExpr)
	p.SetState(378)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(376)
			p.TypeInstExpr()
		}

	case nevaParserT__2, nevaParserT__11:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(377)
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
	AllInterfaceDef() []IInterfaceDefContext
	InterfaceDef(i int) IInterfaceDefContext

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

func (s *IoStmtContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *IoStmtContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *IoStmtContext) AllInterfaceDef() []IInterfaceDefContext {
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

func (s *IoStmtContext) InterfaceDef(i int) IInterfaceDefContext {
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
	p.EnterRule(localctx, 42, nevaParserRULE_ioStmt)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(380)
		p.Match(nevaParserT__14)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(381)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(385)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(382)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(387)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(391)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 || _la == nevaParserIDENTIFIER {
		{
			p.SetState(388)
			p.InterfaceDef()
		}

		p.SetState(393)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(394)
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

// IInterfaceDefContext is an interface to support dynamic dispatch.
type IInterfaceDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	InPortsDef() IInPortsDefContext
	OutPortsDef() IOutPortsDefContext
	TypeParams() ITypeParamsContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *InterfaceDefContext) InPortsDef() IInPortsDefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInPortsDefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInPortsDefContext)
}

func (s *InterfaceDefContext) OutPortsDef() IOutPortsDefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOutPortsDefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOutPortsDefContext)
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

func (s *InterfaceDefContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *InterfaceDefContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.EnterRule(localctx, 44, nevaParserRULE_interfaceDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(397)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__7 {
		{
			p.SetState(396)
			p.Match(nevaParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(399)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(401)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__8 {
		{
			p.SetState(400)
			p.TypeParams()
		}

	}
	{
		p.SetState(403)
		p.InPortsDef()
	}
	{
		p.SetState(404)
		p.OutPortsDef()
	}
	p.SetState(408)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(405)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(410)
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

// IInPortsDefContext is an interface to support dynamic dispatch.
type IInPortsDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PortsDef() IPortsDefContext

	// IsInPortsDefContext differentiates from other interfaces.
	IsInPortsDefContext()
}

type InPortsDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInPortsDefContext() *InPortsDefContext {
	var p = new(InPortsDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_inPortsDef
	return p
}

func InitEmptyInPortsDefContext(p *InPortsDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_inPortsDef
}

func (*InPortsDefContext) IsInPortsDefContext() {}

func NewInPortsDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InPortsDefContext {
	var p = new(InPortsDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_inPortsDef

	return p
}

func (s *InPortsDefContext) GetParser() antlr.Parser { return s.parser }

func (s *InPortsDefContext) PortsDef() IPortsDefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPortsDefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPortsDefContext)
}

func (s *InPortsDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InPortsDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *InPortsDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterInPortsDef(s)
	}
}

func (s *InPortsDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitInPortsDef(s)
	}
}

func (p *nevaParser) InPortsDef() (localctx IInPortsDefContext) {
	localctx = NewInPortsDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, nevaParserRULE_inPortsDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(411)
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

// IOutPortsDefContext is an interface to support dynamic dispatch.
type IOutPortsDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PortsDef() IPortsDefContext

	// IsOutPortsDefContext differentiates from other interfaces.
	IsOutPortsDefContext()
}

type OutPortsDefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOutPortsDefContext() *OutPortsDefContext {
	var p = new(OutPortsDefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_outPortsDef
	return p
}

func InitEmptyOutPortsDefContext(p *OutPortsDefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_outPortsDef
}

func (*OutPortsDefContext) IsOutPortsDefContext() {}

func NewOutPortsDefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OutPortsDefContext {
	var p = new(OutPortsDefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_outPortsDef

	return p
}

func (s *OutPortsDefContext) GetParser() antlr.Parser { return s.parser }

func (s *OutPortsDefContext) PortsDef() IPortsDefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPortsDefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPortsDefContext)
}

func (s *OutPortsDefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OutPortsDefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OutPortsDefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterOutPortsDef(s)
	}
}

func (s *OutPortsDefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitOutPortsDef(s)
	}
}

func (p *nevaParser) OutPortsDef() (localctx IOutPortsDefContext) {
	localctx = NewOutPortsDefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, nevaParserRULE_outPortsDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(413)
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
	AllPortAndType() []IPortAndTypeContext
	PortAndType(i int) IPortAndTypeContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *PortsDefContext) AllPortAndType() []IPortAndTypeContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPortAndTypeContext); ok {
			len++
		}
	}

	tst := make([]IPortAndTypeContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPortAndTypeContext); ok {
			tst[i] = t.(IPortAndTypeContext)
			i++
		}
	}

	return tst
}

func (s *PortsDefContext) PortAndType(i int) IPortAndTypeContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPortAndTypeContext); ok {
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

	return t.(IPortAndTypeContext)
}

func (s *PortsDefContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *PortsDefContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.EnterRule(localctx, 50, nevaParserRULE_portsDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(415)
		p.Match(nevaParserT__15)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(433)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 50, p.GetParserRuleContext()) {
	case 1:
		p.SetState(419)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(416)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(421)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		p.SetState(423)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserIDENTIFIER || _la == nevaParserNEWLINE {
			{
				p.SetState(422)
				p.PortAndType()
			}

		}

	case 3:
		{
			p.SetState(425)
			p.PortAndType()
		}
		p.SetState(430)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserT__9 {
			{
				p.SetState(426)
				p.Match(nevaParserT__9)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(427)
				p.PortAndType()
			}

			p.SetState(432)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	{
		p.SetState(435)
		p.Match(nevaParserT__16)
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

// IPortAndTypeContext is an interface to support dynamic dispatch.
type IPortAndTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeExpr() ITypeExprContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsPortAndTypeContext differentiates from other interfaces.
	IsPortAndTypeContext()
}

type PortAndTypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPortAndTypeContext() *PortAndTypeContext {
	var p = new(PortAndTypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portAndType
	return p
}

func InitEmptyPortAndTypeContext(p *PortAndTypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_portAndType
}

func (*PortAndTypeContext) IsPortAndTypeContext() {}

func NewPortAndTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PortAndTypeContext {
	var p = new(PortAndTypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_portAndType

	return p
}

func (s *PortAndTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *PortAndTypeContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, 0)
}

func (s *PortAndTypeContext) TypeExpr() ITypeExprContext {
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

func (s *PortAndTypeContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *PortAndTypeContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *PortAndTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PortAndTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PortAndTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterPortAndType(s)
	}
}

func (s *PortAndTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitPortAndType(s)
	}
}

func (p *nevaParser) PortAndType() (localctx IPortAndTypeContext) {
	localctx = NewPortAndTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, nevaParserRULE_portAndType)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(440)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(437)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(442)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(443)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(444)
		p.TypeExpr()
	}
	p.SetState(448)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(445)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(450)
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
	p.EnterRule(localctx, 54, nevaParserRULE_constStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(451)
		p.Match(nevaParserT__17)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(452)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(453)
		p.ConstDefList()
	}
	{
		p.SetState(454)
		p.Match(nevaParserT__3)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(455)
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
	p.EnterRule(localctx, 56, nevaParserRULE_constDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(457)
		p.ConstDef()
	}
	p.SetState(462)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(458)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(459)
			p.ConstDef()
		}

		p.SetState(464)
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
	p.EnterRule(localctx, 58, nevaParserRULE_constDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(466)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__7 {
		{
			p.SetState(465)
			p.Match(nevaParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(468)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(469)
		p.TypeExpr()
	}
	{
		p.SetState(470)
		p.Match(nevaParserT__18)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(471)
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
	p.EnterRule(localctx, 60, nevaParserRULE_constValue)
	p.SetState(481)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__19:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(473)
			p.Match(nevaParserT__19)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserT__20:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(474)
			p.Match(nevaParserT__20)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserINT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(475)
			p.Match(nevaParserINT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserFLOAT:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(476)
			p.Match(nevaParserFLOAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserSTRING:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(477)
			p.Match(nevaParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserT__11:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(478)
			p.ArrLit()
		}

	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(479)
			p.RecLit()
		}

	case nevaParserT__21:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(480)
			p.Match(nevaParserT__21)
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
	p.EnterRule(localctx, 62, nevaParserRULE_arrLit)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(483)
		p.Match(nevaParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(484)
		p.ArrItems()
	}
	{
		p.SetState(485)
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
	p.EnterRule(localctx, 64, nevaParserRULE_arrItems)
	var _la int

	p.SetState(499)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 58, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(487)
			p.ConstValue()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(488)
			p.ConstValue()
		}
		p.SetState(496)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserT__9 {
			{
				p.SetState(489)
				p.Match(nevaParserT__9)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(491)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if _la == nevaParserNEWLINE {
				{
					p.SetState(490)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			}
			{
				p.SetState(493)
				p.ConstValue()
			}

			p.SetState(498)
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
	p.EnterRule(localctx, 66, nevaParserRULE_recLit)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(501)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(502)
		p.RecValueFields()
	}
	{
		p.SetState(503)
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
	p.EnterRule(localctx, 68, nevaParserRULE_recValueFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(505)
		p.RecValueField()
	}
	p.SetState(513)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__9 {
		{
			p.SetState(506)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(508)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserNEWLINE {
			{
				p.SetState(507)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(510)
			p.RecValueField()
		}

		p.SetState(515)
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
	p.EnterRule(localctx, 70, nevaParserRULE_recValueField)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(516)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(517)
		p.Match(nevaParserT__22)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(518)
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
	p.EnterRule(localctx, 72, nevaParserRULE_compStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(520)
		p.Match(nevaParserT__23)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(521)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(522)
		p.CompDefList()
	}
	{
		p.SetState(523)
		p.Match(nevaParserT__3)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(524)
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
	p.EnterRule(localctx, 74, nevaParserRULE_compDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(526)
		p.CompDef()
	}
	p.SetState(531)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(527)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(528)
			p.CompDef()
		}

		p.SetState(533)
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
	p.EnterRule(localctx, 76, nevaParserRULE_compDef)
	p.EnterOuterAlt(localctx, 1)
	p.SetState(535)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 62, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(534)
			p.Match(nevaParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(537)
		p.InterfaceDef()
	}
	{
		p.SetState(538)
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
	p.EnterRule(localctx, 78, nevaParserRULE_compBody)
	p.SetState(545)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(540)
			p.Match(nevaParserT__2)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(541)
			p.CompNodesDef()
		}

	case nevaParserT__26:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(542)
			p.CompNetDef()
		}
		{
			p.SetState(543)
			p.Match(nevaParserT__3)
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
	p.EnterRule(localctx, 80, nevaParserRULE_compNodesDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(547)
		p.Match(nevaParserT__24)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(548)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(549)
		p.CompNodeDefList()
	}
	{
		p.SetState(550)
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
	p.EnterRule(localctx, 82, nevaParserRULE_compNodeDefList)
	p.SetState(554)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 64, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(552)
			p.AbsNodeDef()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(553)
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
	p.EnterRule(localctx, 84, nevaParserRULE_absNodeDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(556)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(557)
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
	p.EnterRule(localctx, 86, nevaParserRULE_concreteNodeDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(559)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(560)
		p.Match(nevaParserT__18)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(561)
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
	p.EnterRule(localctx, 88, nevaParserRULE_concreteNodeInst)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(563)
		p.NodeRef()
	}
	{
		p.SetState(564)
		p.NodeArgs()
	}
	{
		p.SetState(565)
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
	p.EnterRule(localctx, 90, nevaParserRULE_nodeRef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(567)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(572)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__25 {
		{
			p.SetState(568)
			p.Match(nevaParserT__25)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(569)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(574)
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
	p.EnterRule(localctx, 92, nevaParserRULE_nodeArgs)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(575)
		p.Match(nevaParserT__15)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(576)
		p.NodeArgList()
	}
	{
		p.SetState(577)
		p.Match(nevaParserT__16)
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
	p.EnterRule(localctx, 94, nevaParserRULE_nodeArgList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(579)
		p.NodeArg()
	}

	{
		p.SetState(580)
		p.Match(nevaParserT__9)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(582)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserNEWLINE {
		{
			p.SetState(581)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(584)
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
	p.EnterRule(localctx, 96, nevaParserRULE_nodeArg)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(586)
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
	p.EnterRule(localctx, 98, nevaParserRULE_compNetDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(588)
		p.Match(nevaParserT__26)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(589)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(590)
		p.ConnDefList()
	}
	{
		p.SetState(591)
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
	p.EnterRule(localctx, 100, nevaParserRULE_connDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(593)
		p.ConnDef()
	}
	p.SetState(598)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(594)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(595)
			p.ConnDef()
		}

		p.SetState(600)
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
	p.EnterRule(localctx, 102, nevaParserRULE_connDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(601)
		p.PortAddr()
	}
	{
		p.SetState(602)
		p.Match(nevaParserT__27)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(603)
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
	p.EnterRule(localctx, 104, nevaParserRULE_portAddr)
	var _la int

	p.SetState(615)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 70, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(606)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserIDENTIFIER {
			{
				p.SetState(605)
				p.Match(nevaParserIDENTIFIER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(608)
			p.PortDirection()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(609)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(613)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserT__11 {
			{
				p.SetState(610)
				p.Match(nevaParserT__11)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(611)
				p.Match(nevaParserINT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(612)
				p.Match(nevaParserT__12)
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
	p.EnterRule(localctx, 106, nevaParserRULE_portDirection)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(617)
		_la = p.GetTokenStream().LA(1)

		if !(_la == nevaParserT__28 || _la == nevaParserT__29) {
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
	p.EnterRule(localctx, 108, nevaParserRULE_connReceiverSide)
	p.SetState(621)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__28, nevaParserT__29, nevaParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(619)
			p.PortAddr()
		}

	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(620)
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
	p.EnterRule(localctx, 110, nevaParserRULE_connReceivers)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(623)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(624)
		p.PortAddr()
	}
	p.SetState(629)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(625)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(626)
			p.PortAddr()
		}

		p.SetState(631)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(632)
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

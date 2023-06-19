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
		4, 1, 36, 566, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 1, 0, 1, 0, 5, 0, 115, 8,
		0, 10, 0, 12, 0, 118, 9, 0, 1, 0, 1, 0, 5, 0, 122, 8, 0, 10, 0, 12, 0,
		125, 9, 0, 5, 0, 127, 8, 0, 10, 0, 12, 0, 130, 9, 0, 1, 0, 1, 0, 1, 1,
		1, 1, 5, 1, 136, 8, 1, 10, 1, 12, 1, 139, 9, 1, 1, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 3, 2, 146, 8, 2, 1, 3, 1, 3, 1, 3, 5, 3, 151, 8, 3, 10, 3, 12, 3,
		154, 9, 3, 1, 3, 5, 3, 157, 8, 3, 10, 3, 12, 3, 160, 9, 3, 1, 3, 1, 3,
		1, 4, 3, 4, 165, 8, 4, 1, 4, 1, 4, 5, 4, 169, 8, 4, 10, 4, 12, 4, 172,
		9, 4, 1, 5, 3, 5, 175, 8, 5, 1, 5, 1, 5, 1, 5, 5, 5, 180, 8, 5, 10, 5,
		12, 5, 183, 9, 5, 1, 6, 1, 6, 1, 6, 5, 6, 188, 8, 6, 10, 6, 12, 6, 191,
		9, 6, 1, 6, 5, 6, 194, 8, 6, 10, 6, 12, 6, 197, 9, 6, 1, 6, 1, 6, 1, 7,
		3, 7, 202, 8, 7, 1, 7, 1, 7, 3, 7, 206, 8, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1,
		8, 1, 8, 5, 8, 214, 8, 8, 10, 8, 12, 8, 217, 9, 8, 1, 8, 1, 8, 1, 9, 5,
		9, 222, 8, 9, 10, 9, 12, 9, 225, 9, 9, 1, 9, 1, 9, 3, 9, 229, 8, 9, 1,
		9, 5, 9, 232, 8, 9, 10, 9, 12, 9, 235, 9, 9, 1, 10, 5, 10, 238, 8, 10,
		10, 10, 12, 10, 241, 9, 10, 1, 10, 1, 10, 1, 10, 3, 10, 246, 8, 10, 1,
		10, 5, 10, 249, 8, 10, 10, 10, 12, 10, 252, 9, 10, 1, 11, 1, 11, 3, 11,
		256, 8, 11, 1, 12, 1, 12, 1, 12, 1, 12, 5, 12, 262, 8, 12, 10, 12, 12,
		12, 265, 9, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 3, 13, 272, 8, 13, 1,
		14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 15, 1, 15, 5, 15, 281, 8, 15, 10, 15,
		12, 15, 284, 9, 15, 1, 15, 3, 15, 287, 8, 15, 1, 15, 1, 15, 1, 16, 1, 16,
		4, 16, 293, 8, 16, 11, 16, 12, 16, 294, 1, 16, 5, 16, 298, 8, 16, 10, 16,
		12, 16, 301, 9, 16, 1, 17, 1, 17, 1, 17, 1, 18, 1, 18, 1, 18, 4, 18, 309,
		8, 18, 11, 18, 12, 18, 310, 1, 19, 1, 19, 5, 19, 315, 8, 19, 10, 19, 12,
		19, 318, 9, 19, 1, 19, 4, 19, 321, 8, 19, 11, 19, 12, 19, 322, 1, 19, 5,
		19, 326, 8, 19, 10, 19, 12, 19, 329, 9, 19, 1, 19, 1, 19, 5, 19, 333, 8,
		19, 10, 19, 12, 19, 336, 9, 19, 1, 19, 1, 19, 1, 20, 1, 20, 3, 20, 342,
		8, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 5, 22, 352,
		8, 22, 10, 22, 12, 22, 355, 9, 22, 1, 23, 3, 23, 358, 8, 23, 1, 23, 1,
		23, 1, 23, 1, 23, 1, 23, 1, 24, 1, 24, 1, 24, 1, 24, 1, 25, 1, 25, 1, 25,
		3, 25, 372, 8, 25, 1, 25, 5, 25, 375, 8, 25, 10, 25, 12, 25, 378, 9, 25,
		1, 26, 1, 26, 1, 26, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 28, 1,
		28, 1, 28, 5, 28, 392, 8, 28, 10, 28, 12, 28, 395, 9, 28, 1, 29, 3, 29,
		398, 8, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 30, 1, 30, 1, 30, 1,
		30, 1, 30, 1, 30, 1, 30, 1, 30, 3, 30, 413, 8, 30, 1, 31, 1, 31, 1, 31,
		1, 31, 1, 32, 1, 32, 1, 32, 1, 32, 3, 32, 423, 8, 32, 1, 32, 5, 32, 426,
		8, 32, 10, 32, 12, 32, 429, 9, 32, 3, 32, 431, 8, 32, 1, 33, 1, 33, 1,
		33, 1, 33, 1, 34, 1, 34, 1, 34, 3, 34, 440, 8, 34, 1, 34, 5, 34, 443, 8,
		34, 10, 34, 12, 34, 446, 9, 34, 1, 35, 1, 35, 1, 35, 1, 35, 1, 36, 1, 36,
		1, 36, 1, 36, 1, 36, 1, 36, 1, 37, 1, 37, 1, 37, 5, 37, 461, 8, 37, 10,
		37, 12, 37, 464, 9, 37, 1, 38, 3, 38, 467, 8, 38, 1, 38, 1, 38, 1, 38,
		1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 3, 39, 477, 8, 39, 1, 40, 1, 40, 1,
		40, 1, 40, 1, 40, 1, 41, 1, 41, 3, 41, 486, 8, 41, 1, 42, 1, 42, 1, 42,
		1, 43, 1, 43, 1, 43, 1, 43, 1, 44, 1, 44, 1, 44, 1, 44, 1, 45, 1, 45, 1,
		45, 5, 45, 502, 8, 45, 10, 45, 12, 45, 505, 9, 45, 1, 46, 1, 46, 1, 46,
		1, 46, 1, 47, 1, 47, 1, 47, 3, 47, 514, 8, 47, 1, 47, 1, 47, 1, 48, 1,
		48, 1, 49, 1, 49, 1, 49, 1, 49, 1, 49, 1, 50, 1, 50, 1, 50, 5, 50, 528,
		8, 50, 10, 50, 12, 50, 531, 9, 50, 1, 51, 1, 51, 1, 51, 1, 51, 1, 52, 3,
		52, 538, 8, 52, 1, 52, 1, 52, 1, 52, 1, 52, 1, 52, 3, 52, 545, 8, 52, 3,
		52, 547, 8, 52, 1, 53, 1, 53, 1, 54, 1, 54, 3, 54, 553, 8, 54, 1, 55, 1,
		55, 1, 55, 1, 55, 5, 55, 559, 8, 55, 10, 55, 12, 55, 562, 9, 55, 1, 55,
		1, 55, 1, 55, 0, 0, 56, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24,
		26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60,
		62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96,
		98, 100, 102, 104, 106, 108, 110, 0, 2, 1, 0, 35, 35, 1, 0, 29, 30, 580,
		0, 128, 1, 0, 0, 0, 2, 133, 1, 0, 0, 0, 4, 145, 1, 0, 0, 0, 6, 147, 1,
		0, 0, 0, 8, 164, 1, 0, 0, 0, 10, 174, 1, 0, 0, 0, 12, 184, 1, 0, 0, 0,
		14, 201, 1, 0, 0, 0, 16, 209, 1, 0, 0, 0, 18, 223, 1, 0, 0, 0, 20, 239,
		1, 0, 0, 0, 22, 253, 1, 0, 0, 0, 24, 257, 1, 0, 0, 0, 26, 271, 1, 0, 0,
		0, 28, 273, 1, 0, 0, 0, 30, 278, 1, 0, 0, 0, 32, 290, 1, 0, 0, 0, 34, 302,
		1, 0, 0, 0, 36, 305, 1, 0, 0, 0, 38, 312, 1, 0, 0, 0, 40, 341, 1, 0, 0,
		0, 42, 343, 1, 0, 0, 0, 44, 348, 1, 0, 0, 0, 46, 357, 1, 0, 0, 0, 48, 364,
		1, 0, 0, 0, 50, 368, 1, 0, 0, 0, 52, 379, 1, 0, 0, 0, 54, 382, 1, 0, 0,
		0, 56, 388, 1, 0, 0, 0, 58, 397, 1, 0, 0, 0, 60, 412, 1, 0, 0, 0, 62, 414,
		1, 0, 0, 0, 64, 430, 1, 0, 0, 0, 66, 432, 1, 0, 0, 0, 68, 436, 1, 0, 0,
		0, 70, 447, 1, 0, 0, 0, 72, 451, 1, 0, 0, 0, 74, 457, 1, 0, 0, 0, 76, 466,
		1, 0, 0, 0, 78, 476, 1, 0, 0, 0, 80, 478, 1, 0, 0, 0, 82, 485, 1, 0, 0,
		0, 84, 487, 1, 0, 0, 0, 86, 490, 1, 0, 0, 0, 88, 494, 1, 0, 0, 0, 90, 498,
		1, 0, 0, 0, 92, 506, 1, 0, 0, 0, 94, 510, 1, 0, 0, 0, 96, 517, 1, 0, 0,
		0, 98, 519, 1, 0, 0, 0, 100, 524, 1, 0, 0, 0, 102, 532, 1, 0, 0, 0, 104,
		546, 1, 0, 0, 0, 106, 548, 1, 0, 0, 0, 108, 552, 1, 0, 0, 0, 110, 554,
		1, 0, 0, 0, 112, 116, 3, 2, 1, 0, 113, 115, 5, 35, 0, 0, 114, 113, 1, 0,
		0, 0, 115, 118, 1, 0, 0, 0, 116, 114, 1, 0, 0, 0, 116, 117, 1, 0, 0, 0,
		117, 127, 1, 0, 0, 0, 118, 116, 1, 0, 0, 0, 119, 123, 3, 4, 2, 0, 120,
		122, 5, 35, 0, 0, 121, 120, 1, 0, 0, 0, 122, 125, 1, 0, 0, 0, 123, 121,
		1, 0, 0, 0, 123, 124, 1, 0, 0, 0, 124, 127, 1, 0, 0, 0, 125, 123, 1, 0,
		0, 0, 126, 112, 1, 0, 0, 0, 126, 119, 1, 0, 0, 0, 127, 130, 1, 0, 0, 0,
		128, 126, 1, 0, 0, 0, 128, 129, 1, 0, 0, 0, 129, 131, 1, 0, 0, 0, 130,
		128, 1, 0, 0, 0, 131, 132, 5, 0, 0, 1, 132, 1, 1, 0, 0, 0, 133, 137, 5,
		1, 0, 0, 134, 136, 8, 0, 0, 0, 135, 134, 1, 0, 0, 0, 136, 139, 1, 0, 0,
		0, 137, 135, 1, 0, 0, 0, 137, 138, 1, 0, 0, 0, 138, 3, 1, 0, 0, 0, 139,
		137, 1, 0, 0, 0, 140, 146, 3, 6, 3, 0, 141, 146, 3, 12, 6, 0, 142, 146,
		3, 42, 21, 0, 143, 146, 3, 54, 27, 0, 144, 146, 3, 72, 36, 0, 145, 140,
		1, 0, 0, 0, 145, 141, 1, 0, 0, 0, 145, 142, 1, 0, 0, 0, 145, 143, 1, 0,
		0, 0, 145, 144, 1, 0, 0, 0, 146, 5, 1, 0, 0, 0, 147, 148, 5, 2, 0, 0, 148,
		152, 5, 3, 0, 0, 149, 151, 5, 35, 0, 0, 150, 149, 1, 0, 0, 0, 151, 154,
		1, 0, 0, 0, 152, 150, 1, 0, 0, 0, 152, 153, 1, 0, 0, 0, 153, 158, 1, 0,
		0, 0, 154, 152, 1, 0, 0, 0, 155, 157, 3, 8, 4, 0, 156, 155, 1, 0, 0, 0,
		157, 160, 1, 0, 0, 0, 158, 156, 1, 0, 0, 0, 158, 159, 1, 0, 0, 0, 159,
		161, 1, 0, 0, 0, 160, 158, 1, 0, 0, 0, 161, 162, 5, 4, 0, 0, 162, 7, 1,
		0, 0, 0, 163, 165, 5, 31, 0, 0, 164, 163, 1, 0, 0, 0, 164, 165, 1, 0, 0,
		0, 165, 166, 1, 0, 0, 0, 166, 170, 3, 10, 5, 0, 167, 169, 5, 35, 0, 0,
		168, 167, 1, 0, 0, 0, 169, 172, 1, 0, 0, 0, 170, 168, 1, 0, 0, 0, 170,
		171, 1, 0, 0, 0, 171, 9, 1, 0, 0, 0, 172, 170, 1, 0, 0, 0, 173, 175, 5,
		5, 0, 0, 174, 173, 1, 0, 0, 0, 174, 175, 1, 0, 0, 0, 175, 176, 1, 0, 0,
		0, 176, 181, 5, 31, 0, 0, 177, 178, 5, 6, 0, 0, 178, 180, 5, 31, 0, 0,
		179, 177, 1, 0, 0, 0, 180, 183, 1, 0, 0, 0, 181, 179, 1, 0, 0, 0, 181,
		182, 1, 0, 0, 0, 182, 11, 1, 0, 0, 0, 183, 181, 1, 0, 0, 0, 184, 185, 5,
		7, 0, 0, 185, 189, 5, 3, 0, 0, 186, 188, 5, 35, 0, 0, 187, 186, 1, 0, 0,
		0, 188, 191, 1, 0, 0, 0, 189, 187, 1, 0, 0, 0, 189, 190, 1, 0, 0, 0, 190,
		195, 1, 0, 0, 0, 191, 189, 1, 0, 0, 0, 192, 194, 3, 14, 7, 0, 193, 192,
		1, 0, 0, 0, 194, 197, 1, 0, 0, 0, 195, 193, 1, 0, 0, 0, 195, 196, 1, 0,
		0, 0, 196, 198, 1, 0, 0, 0, 197, 195, 1, 0, 0, 0, 198, 199, 5, 4, 0, 0,
		199, 13, 1, 0, 0, 0, 200, 202, 5, 8, 0, 0, 201, 200, 1, 0, 0, 0, 201, 202,
		1, 0, 0, 0, 202, 203, 1, 0, 0, 0, 203, 205, 5, 31, 0, 0, 204, 206, 3, 16,
		8, 0, 205, 204, 1, 0, 0, 0, 205, 206, 1, 0, 0, 0, 206, 207, 1, 0, 0, 0,
		207, 208, 3, 20, 10, 0, 208, 15, 1, 0, 0, 0, 209, 210, 5, 9, 0, 0, 210,
		215, 3, 18, 9, 0, 211, 212, 5, 10, 0, 0, 212, 214, 3, 18, 9, 0, 213, 211,
		1, 0, 0, 0, 214, 217, 1, 0, 0, 0, 215, 213, 1, 0, 0, 0, 215, 216, 1, 0,
		0, 0, 216, 218, 1, 0, 0, 0, 217, 215, 1, 0, 0, 0, 218, 219, 5, 11, 0, 0,
		219, 17, 1, 0, 0, 0, 220, 222, 5, 35, 0, 0, 221, 220, 1, 0, 0, 0, 222,
		225, 1, 0, 0, 0, 223, 221, 1, 0, 0, 0, 223, 224, 1, 0, 0, 0, 224, 226,
		1, 0, 0, 0, 225, 223, 1, 0, 0, 0, 226, 228, 5, 31, 0, 0, 227, 229, 3, 20,
		10, 0, 228, 227, 1, 0, 0, 0, 228, 229, 1, 0, 0, 0, 229, 233, 1, 0, 0, 0,
		230, 232, 5, 35, 0, 0, 231, 230, 1, 0, 0, 0, 232, 235, 1, 0, 0, 0, 233,
		231, 1, 0, 0, 0, 233, 234, 1, 0, 0, 0, 234, 19, 1, 0, 0, 0, 235, 233, 1,
		0, 0, 0, 236, 238, 5, 35, 0, 0, 237, 236, 1, 0, 0, 0, 238, 241, 1, 0, 0,
		0, 239, 237, 1, 0, 0, 0, 239, 240, 1, 0, 0, 0, 240, 245, 1, 0, 0, 0, 241,
		239, 1, 0, 0, 0, 242, 246, 3, 22, 11, 0, 243, 246, 3, 26, 13, 0, 244, 246,
		3, 36, 18, 0, 245, 242, 1, 0, 0, 0, 245, 243, 1, 0, 0, 0, 245, 244, 1,
		0, 0, 0, 246, 250, 1, 0, 0, 0, 247, 249, 5, 35, 0, 0, 248, 247, 1, 0, 0,
		0, 249, 252, 1, 0, 0, 0, 250, 248, 1, 0, 0, 0, 250, 251, 1, 0, 0, 0, 251,
		21, 1, 0, 0, 0, 252, 250, 1, 0, 0, 0, 253, 255, 5, 31, 0, 0, 254, 256,
		3, 24, 12, 0, 255, 254, 1, 0, 0, 0, 255, 256, 1, 0, 0, 0, 256, 23, 1, 0,
		0, 0, 257, 258, 5, 9, 0, 0, 258, 263, 3, 20, 10, 0, 259, 260, 5, 10, 0,
		0, 260, 262, 3, 20, 10, 0, 261, 259, 1, 0, 0, 0, 262, 265, 1, 0, 0, 0,
		263, 261, 1, 0, 0, 0, 263, 264, 1, 0, 0, 0, 264, 266, 1, 0, 0, 0, 265,
		263, 1, 0, 0, 0, 266, 267, 5, 11, 0, 0, 267, 25, 1, 0, 0, 0, 268, 272,
		3, 28, 14, 0, 269, 272, 3, 30, 15, 0, 270, 272, 3, 38, 19, 0, 271, 268,
		1, 0, 0, 0, 271, 269, 1, 0, 0, 0, 271, 270, 1, 0, 0, 0, 272, 27, 1, 0,
		0, 0, 273, 274, 5, 12, 0, 0, 274, 275, 5, 32, 0, 0, 275, 276, 5, 13, 0,
		0, 276, 277, 3, 20, 10, 0, 277, 29, 1, 0, 0, 0, 278, 282, 5, 3, 0, 0, 279,
		281, 5, 35, 0, 0, 280, 279, 1, 0, 0, 0, 281, 284, 1, 0, 0, 0, 282, 280,
		1, 0, 0, 0, 282, 283, 1, 0, 0, 0, 283, 286, 1, 0, 0, 0, 284, 282, 1, 0,
		0, 0, 285, 287, 3, 32, 16, 0, 286, 285, 1, 0, 0, 0, 286, 287, 1, 0, 0,
		0, 287, 288, 1, 0, 0, 0, 288, 289, 5, 4, 0, 0, 289, 31, 1, 0, 0, 0, 290,
		299, 3, 34, 17, 0, 291, 293, 5, 35, 0, 0, 292, 291, 1, 0, 0, 0, 293, 294,
		1, 0, 0, 0, 294, 292, 1, 0, 0, 0, 294, 295, 1, 0, 0, 0, 295, 296, 1, 0,
		0, 0, 296, 298, 3, 34, 17, 0, 297, 292, 1, 0, 0, 0, 298, 301, 1, 0, 0,
		0, 299, 297, 1, 0, 0, 0, 299, 300, 1, 0, 0, 0, 300, 33, 1, 0, 0, 0, 301,
		299, 1, 0, 0, 0, 302, 303, 5, 31, 0, 0, 303, 304, 3, 20, 10, 0, 304, 35,
		1, 0, 0, 0, 305, 308, 3, 40, 20, 0, 306, 307, 5, 14, 0, 0, 307, 309, 3,
		40, 20, 0, 308, 306, 1, 0, 0, 0, 309, 310, 1, 0, 0, 0, 310, 308, 1, 0,
		0, 0, 310, 311, 1, 0, 0, 0, 311, 37, 1, 0, 0, 0, 312, 316, 5, 3, 0, 0,
		313, 315, 5, 35, 0, 0, 314, 313, 1, 0, 0, 0, 315, 318, 1, 0, 0, 0, 316,
		314, 1, 0, 0, 0, 316, 317, 1, 0, 0, 0, 317, 320, 1, 0, 0, 0, 318, 316,
		1, 0, 0, 0, 319, 321, 5, 31, 0, 0, 320, 319, 1, 0, 0, 0, 321, 322, 1, 0,
		0, 0, 322, 320, 1, 0, 0, 0, 322, 323, 1, 0, 0, 0, 323, 334, 1, 0, 0, 0,
		324, 326, 5, 35, 0, 0, 325, 324, 1, 0, 0, 0, 326, 329, 1, 0, 0, 0, 327,
		325, 1, 0, 0, 0, 327, 328, 1, 0, 0, 0, 328, 330, 1, 0, 0, 0, 329, 327,
		1, 0, 0, 0, 330, 331, 5, 10, 0, 0, 331, 333, 5, 31, 0, 0, 332, 327, 1,
		0, 0, 0, 333, 336, 1, 0, 0, 0, 334, 332, 1, 0, 0, 0, 334, 335, 1, 0, 0,
		0, 335, 337, 1, 0, 0, 0, 336, 334, 1, 0, 0, 0, 337, 338, 5, 4, 0, 0, 338,
		39, 1, 0, 0, 0, 339, 342, 3, 22, 11, 0, 340, 342, 3, 26, 13, 0, 341, 339,
		1, 0, 0, 0, 341, 340, 1, 0, 0, 0, 342, 41, 1, 0, 0, 0, 343, 344, 5, 15,
		0, 0, 344, 345, 5, 3, 0, 0, 345, 346, 3, 44, 22, 0, 346, 347, 5, 4, 0,
		0, 347, 43, 1, 0, 0, 0, 348, 353, 3, 46, 23, 0, 349, 350, 5, 35, 0, 0,
		350, 352, 3, 46, 23, 0, 351, 349, 1, 0, 0, 0, 352, 355, 1, 0, 0, 0, 353,
		351, 1, 0, 0, 0, 353, 354, 1, 0, 0, 0, 354, 45, 1, 0, 0, 0, 355, 353, 1,
		0, 0, 0, 356, 358, 5, 8, 0, 0, 357, 356, 1, 0, 0, 0, 357, 358, 1, 0, 0,
		0, 358, 359, 1, 0, 0, 0, 359, 360, 5, 31, 0, 0, 360, 361, 3, 16, 8, 0,
		361, 362, 3, 48, 24, 0, 362, 363, 3, 48, 24, 0, 363, 47, 1, 0, 0, 0, 364,
		365, 5, 16, 0, 0, 365, 366, 3, 50, 25, 0, 366, 367, 5, 17, 0, 0, 367, 49,
		1, 0, 0, 0, 368, 376, 3, 52, 26, 0, 369, 371, 5, 10, 0, 0, 370, 372, 5,
		35, 0, 0, 371, 370, 1, 0, 0, 0, 371, 372, 1, 0, 0, 0, 372, 373, 1, 0, 0,
		0, 373, 375, 3, 52, 26, 0, 374, 369, 1, 0, 0, 0, 375, 378, 1, 0, 0, 0,
		376, 374, 1, 0, 0, 0, 376, 377, 1, 0, 0, 0, 377, 51, 1, 0, 0, 0, 378, 376,
		1, 0, 0, 0, 379, 380, 5, 31, 0, 0, 380, 381, 3, 20, 10, 0, 381, 53, 1,
		0, 0, 0, 382, 383, 5, 18, 0, 0, 383, 384, 5, 3, 0, 0, 384, 385, 3, 56,
		28, 0, 385, 386, 5, 4, 0, 0, 386, 387, 5, 35, 0, 0, 387, 55, 1, 0, 0, 0,
		388, 393, 3, 58, 29, 0, 389, 390, 5, 35, 0, 0, 390, 392, 3, 58, 29, 0,
		391, 389, 1, 0, 0, 0, 392, 395, 1, 0, 0, 0, 393, 391, 1, 0, 0, 0, 393,
		394, 1, 0, 0, 0, 394, 57, 1, 0, 0, 0, 395, 393, 1, 0, 0, 0, 396, 398, 5,
		8, 0, 0, 397, 396, 1, 0, 0, 0, 397, 398, 1, 0, 0, 0, 398, 399, 1, 0, 0,
		0, 399, 400, 5, 31, 0, 0, 400, 401, 3, 20, 10, 0, 401, 402, 5, 19, 0, 0,
		402, 403, 3, 60, 30, 0, 403, 59, 1, 0, 0, 0, 404, 413, 5, 20, 0, 0, 405,
		413, 5, 21, 0, 0, 406, 413, 5, 32, 0, 0, 407, 413, 5, 33, 0, 0, 408, 413,
		5, 34, 0, 0, 409, 413, 3, 62, 31, 0, 410, 413, 3, 66, 33, 0, 411, 413,
		5, 22, 0, 0, 412, 404, 1, 0, 0, 0, 412, 405, 1, 0, 0, 0, 412, 406, 1, 0,
		0, 0, 412, 407, 1, 0, 0, 0, 412, 408, 1, 0, 0, 0, 412, 409, 1, 0, 0, 0,
		412, 410, 1, 0, 0, 0, 412, 411, 1, 0, 0, 0, 413, 61, 1, 0, 0, 0, 414, 415,
		5, 12, 0, 0, 415, 416, 3, 64, 32, 0, 416, 417, 5, 13, 0, 0, 417, 63, 1,
		0, 0, 0, 418, 431, 3, 60, 30, 0, 419, 427, 3, 60, 30, 0, 420, 422, 5, 10,
		0, 0, 421, 423, 5, 35, 0, 0, 422, 421, 1, 0, 0, 0, 422, 423, 1, 0, 0, 0,
		423, 424, 1, 0, 0, 0, 424, 426, 3, 60, 30, 0, 425, 420, 1, 0, 0, 0, 426,
		429, 1, 0, 0, 0, 427, 425, 1, 0, 0, 0, 427, 428, 1, 0, 0, 0, 428, 431,
		1, 0, 0, 0, 429, 427, 1, 0, 0, 0, 430, 418, 1, 0, 0, 0, 430, 419, 1, 0,
		0, 0, 431, 65, 1, 0, 0, 0, 432, 433, 5, 3, 0, 0, 433, 434, 3, 68, 34, 0,
		434, 435, 5, 4, 0, 0, 435, 67, 1, 0, 0, 0, 436, 444, 3, 70, 35, 0, 437,
		439, 5, 10, 0, 0, 438, 440, 5, 35, 0, 0, 439, 438, 1, 0, 0, 0, 439, 440,
		1, 0, 0, 0, 440, 441, 1, 0, 0, 0, 441, 443, 3, 70, 35, 0, 442, 437, 1,
		0, 0, 0, 443, 446, 1, 0, 0, 0, 444, 442, 1, 0, 0, 0, 444, 445, 1, 0, 0,
		0, 445, 69, 1, 0, 0, 0, 446, 444, 1, 0, 0, 0, 447, 448, 5, 31, 0, 0, 448,
		449, 5, 23, 0, 0, 449, 450, 3, 60, 30, 0, 450, 71, 1, 0, 0, 0, 451, 452,
		5, 24, 0, 0, 452, 453, 5, 3, 0, 0, 453, 454, 3, 74, 37, 0, 454, 455, 5,
		4, 0, 0, 455, 456, 5, 35, 0, 0, 456, 73, 1, 0, 0, 0, 457, 462, 3, 76, 38,
		0, 458, 459, 5, 35, 0, 0, 459, 461, 3, 76, 38, 0, 460, 458, 1, 0, 0, 0,
		461, 464, 1, 0, 0, 0, 462, 460, 1, 0, 0, 0, 462, 463, 1, 0, 0, 0, 463,
		75, 1, 0, 0, 0, 464, 462, 1, 0, 0, 0, 465, 467, 5, 8, 0, 0, 466, 465, 1,
		0, 0, 0, 466, 467, 1, 0, 0, 0, 467, 468, 1, 0, 0, 0, 468, 469, 3, 46, 23,
		0, 469, 470, 3, 78, 39, 0, 470, 77, 1, 0, 0, 0, 471, 472, 5, 3, 0, 0, 472,
		477, 3, 80, 40, 0, 473, 474, 3, 98, 49, 0, 474, 475, 5, 4, 0, 0, 475, 477,
		1, 0, 0, 0, 476, 471, 1, 0, 0, 0, 476, 473, 1, 0, 0, 0, 477, 79, 1, 0,
		0, 0, 478, 479, 5, 25, 0, 0, 479, 480, 5, 3, 0, 0, 480, 481, 3, 82, 41,
		0, 481, 482, 5, 4, 0, 0, 482, 81, 1, 0, 0, 0, 483, 486, 3, 84, 42, 0, 484,
		486, 3, 86, 43, 0, 485, 483, 1, 0, 0, 0, 485, 484, 1, 0, 0, 0, 486, 83,
		1, 0, 0, 0, 487, 488, 5, 31, 0, 0, 488, 489, 3, 22, 11, 0, 489, 85, 1,
		0, 0, 0, 490, 491, 5, 31, 0, 0, 491, 492, 5, 19, 0, 0, 492, 493, 3, 88,
		44, 0, 493, 87, 1, 0, 0, 0, 494, 495, 3, 90, 45, 0, 495, 496, 3, 92, 46,
		0, 496, 497, 3, 24, 12, 0, 497, 89, 1, 0, 0, 0, 498, 503, 5, 31, 0, 0,
		499, 500, 5, 26, 0, 0, 500, 502, 5, 31, 0, 0, 501, 499, 1, 0, 0, 0, 502,
		505, 1, 0, 0, 0, 503, 501, 1, 0, 0, 0, 503, 504, 1, 0, 0, 0, 504, 91, 1,
		0, 0, 0, 505, 503, 1, 0, 0, 0, 506, 507, 5, 16, 0, 0, 507, 508, 3, 94,
		47, 0, 508, 509, 5, 17, 0, 0, 509, 93, 1, 0, 0, 0, 510, 511, 3, 96, 48,
		0, 511, 513, 5, 10, 0, 0, 512, 514, 5, 35, 0, 0, 513, 512, 1, 0, 0, 0,
		513, 514, 1, 0, 0, 0, 514, 515, 1, 0, 0, 0, 515, 516, 3, 96, 48, 0, 516,
		95, 1, 0, 0, 0, 517, 518, 3, 88, 44, 0, 518, 97, 1, 0, 0, 0, 519, 520,
		5, 27, 0, 0, 520, 521, 5, 3, 0, 0, 521, 522, 3, 100, 50, 0, 522, 523, 5,
		4, 0, 0, 523, 99, 1, 0, 0, 0, 524, 529, 3, 102, 51, 0, 525, 526, 5, 35,
		0, 0, 526, 528, 3, 102, 51, 0, 527, 525, 1, 0, 0, 0, 528, 531, 1, 0, 0,
		0, 529, 527, 1, 0, 0, 0, 529, 530, 1, 0, 0, 0, 530, 101, 1, 0, 0, 0, 531,
		529, 1, 0, 0, 0, 532, 533, 3, 104, 52, 0, 533, 534, 5, 28, 0, 0, 534, 535,
		3, 108, 54, 0, 535, 103, 1, 0, 0, 0, 536, 538, 5, 31, 0, 0, 537, 536, 1,
		0, 0, 0, 537, 538, 1, 0, 0, 0, 538, 539, 1, 0, 0, 0, 539, 547, 3, 106,
		53, 0, 540, 544, 5, 31, 0, 0, 541, 542, 5, 12, 0, 0, 542, 543, 5, 32, 0,
		0, 543, 545, 5, 13, 0, 0, 544, 541, 1, 0, 0, 0, 544, 545, 1, 0, 0, 0, 545,
		547, 1, 0, 0, 0, 546, 537, 1, 0, 0, 0, 546, 540, 1, 0, 0, 0, 547, 105,
		1, 0, 0, 0, 548, 549, 7, 1, 0, 0, 549, 107, 1, 0, 0, 0, 550, 553, 3, 104,
		52, 0, 551, 553, 3, 110, 55, 0, 552, 550, 1, 0, 0, 0, 552, 551, 1, 0, 0,
		0, 553, 109, 1, 0, 0, 0, 554, 555, 5, 3, 0, 0, 555, 560, 3, 104, 52, 0,
		556, 557, 5, 35, 0, 0, 557, 559, 3, 104, 52, 0, 558, 556, 1, 0, 0, 0, 559,
		562, 1, 0, 0, 0, 560, 558, 1, 0, 0, 0, 560, 561, 1, 0, 0, 0, 561, 563,
		1, 0, 0, 0, 562, 560, 1, 0, 0, 0, 563, 564, 5, 4, 0, 0, 564, 111, 1, 0,
		0, 0, 60, 116, 123, 126, 128, 137, 145, 152, 158, 164, 170, 174, 181, 189,
		195, 201, 205, 215, 223, 228, 233, 239, 245, 250, 255, 263, 271, 282, 286,
		294, 299, 310, 316, 322, 327, 334, 341, 353, 357, 371, 376, 393, 397, 412,
		422, 427, 430, 439, 444, 462, 466, 476, 485, 503, 513, 529, 537, 544, 546,
		552, 560,
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
	nevaParserRULE_arrTypeExpr      = 14
	nevaParserRULE_recTypeExpr      = 15
	nevaParserRULE_recTypeFields    = 16
	nevaParserRULE_recTypeField     = 17
	nevaParserRULE_unionTypeExpr    = 18
	nevaParserRULE_enumTypeExpr     = 19
	nevaParserRULE_nonUnionTypeExpr = 20
	nevaParserRULE_ioStmt           = 21
	nevaParserRULE_interfaceDefList = 22
	nevaParserRULE_interfaceDef     = 23
	nevaParserRULE_portsDef         = 24
	nevaParserRULE_portDefList      = 25
	nevaParserRULE_portDef          = 26
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
	AllComment() []ICommentContext
	Comment(i int) ICommentContext
	AllStmt() []IStmtContext
	Stmt(i int) IStmtContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *ProgContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ProgContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.SetState(128)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&17072262) != 0 {
		p.SetState(126)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case nevaParserT__0:
			{
				p.SetState(112)
				p.Comment()
			}
			p.SetState(116)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(113)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(118)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		case nevaParserT__1, nevaParserT__6, nevaParserT__14, nevaParserT__17, nevaParserT__23:
			{
				p.SetState(119)
				p.Stmt()
			}
			p.SetState(123)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(120)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(125)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(130)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(131)
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
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(134)
				_la = p.GetTokenStream().LA(1)

				if _la <= 0 || _la == nevaParserNEWLINE {
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
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext())
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
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(140)
			p.UseStmt()
		}

	case nevaParserT__6:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(141)
			p.TypeStmt()
		}

	case nevaParserT__14:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(142)
			p.IoStmt()
		}

	case nevaParserT__17:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(143)
			p.ConstStmt()
		}

	case nevaParserT__23:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(144)
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
		p.SetState(147)
		p.Match(nevaParserT__1)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(148)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(152)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(149)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(154)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(158)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__4 || _la == nevaParserIDENTIFIER {
		{
			p.SetState(155)
			p.ImportDef()
		}

		p.SetState(160)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(161)
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
	p.SetState(164)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(163)
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
		p.SetState(166)
		p.ImportPath()
	}
	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(167)
			p.Match(nevaParserNEWLINE)
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
	p.SetState(174)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__4 {
		{
			p.SetState(173)
			p.Match(nevaParserT__4)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(176)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(181)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__5 {
		{
			p.SetState(177)
			p.Match(nevaParserT__5)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(178)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(183)
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
		p.SetState(184)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(185)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(189)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(186)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(191)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(195)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 || _la == nevaParserIDENTIFIER {
		{
			p.SetState(192)
			p.TypeDef()
		}

		p.SetState(197)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(198)
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
	p.EnterRule(localctx, 14, nevaParserRULE_typeDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(201)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__7 {
		{
			p.SetState(200)
			p.Match(nevaParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(203)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(205)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__8 {
		{
			p.SetState(204)
			p.TypeParams()
		}

	}
	{
		p.SetState(207)
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
		p.SetState(209)
		p.Match(nevaParserT__8)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(210)
		p.TypeParam()
	}
	p.SetState(215)
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
		{
			p.SetState(212)
			p.TypeParam()
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
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

func (s *TypeParamContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *TypeParamContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.SetState(223)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(220)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(225)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(226)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(228)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 18, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(227)
			p.TypeExpr()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	p.SetState(233)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(230)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(235)
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

// ITypeExprContext is an interface to support dynamic dispatch.
type ITypeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TypeInstExpr() ITypeInstExprContext
	TypeLitExpr() ITypeLitExprContext
	UnionTypeExpr() IUnionTypeExprContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *TypeExprContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *TypeExprContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(239)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(236)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(241)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(245)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(242)
			p.TypeInstExpr()
		}

	case 2:
		{
			p.SetState(243)
			p.TypeLitExpr()
		}

	case 3:
		{
			p.SetState(244)
			p.UnionTypeExpr()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.SetState(250)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 22, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(247)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		p.SetState(252)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 22, p.GetParserRuleContext())
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
		p.SetState(253)
		p.Match(nevaParserIDENTIFIER)
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

	if _la == nevaParserT__8 {
		{
			p.SetState(254)
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
	p.EnterRule(localctx, 24, nevaParserRULE_typeArgs)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(257)
		p.Match(nevaParserT__8)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(258)
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
			p.SetState(259)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
	{
		p.SetState(266)
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
	p.EnterRule(localctx, 26, nevaParserRULE_typeLitExpr)
	p.SetState(271)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 25, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(268)
			p.ArrTypeExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(269)
			p.RecTypeExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(270)
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
	p.EnterRule(localctx, 28, nevaParserRULE_arrTypeExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(273)
		p.Match(nevaParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(274)
		p.Match(nevaParserINT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(275)
		p.Match(nevaParserT__12)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(276)
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

func (s *RecTypeExprContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *RecTypeExprContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

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
	p.EnterRule(localctx, 30, nevaParserRULE_recTypeExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(278)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(282)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(279)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(284)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(286)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserIDENTIFIER {
		{
			p.SetState(285)
			p.RecTypeFields()
		}

	}
	{
		p.SetState(288)
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
	p.EnterRule(localctx, 32, nevaParserRULE_recTypeFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(290)
		p.RecTypeField()
	}
	p.SetState(299)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		p.SetState(292)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == nevaParserNEWLINE {
			{
				p.SetState(291)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(294)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(296)
			p.RecTypeField()
		}

		p.SetState(301)
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
	p.EnterRule(localctx, 34, nevaParserRULE_recTypeField)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(302)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(303)
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
	p.EnterRule(localctx, 36, nevaParserRULE_unionTypeExpr)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(305)
		p.NonUnionTypeExpr()
	}
	p.SetState(308)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(306)
				p.Match(nevaParserT__13)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(307)
				p.NonUnionTypeExpr()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(310)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 30, p.GetParserRuleContext())
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

	// Getter signatures
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode

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

func (s *EnumTypeExprContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *EnumTypeExprContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *EnumTypeExprContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(nevaParserIDENTIFIER)
}

func (s *EnumTypeExprContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserIDENTIFIER, i)
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
	p.EnterRule(localctx, 38, nevaParserRULE_enumTypeExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(312)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(316)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(313)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(318)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(320)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == nevaParserIDENTIFIER {
		{
			p.SetState(319)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(322)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(334)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__9 || _la == nevaParserNEWLINE {
		p.SetState(327)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(324)
				p.Match(nevaParserNEWLINE)
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
		}
		{
			p.SetState(330)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(331)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(336)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(337)
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
	p.SetState(341)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(339)
			p.TypeInstExpr()
		}

	case nevaParserT__2, nevaParserT__11:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(340)
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
	p.EnterRule(localctx, 42, nevaParserRULE_ioStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(343)
		p.Match(nevaParserT__14)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(344)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(345)
		p.InterfaceDefList()
	}
	{
		p.SetState(346)
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
	p.EnterRule(localctx, 44, nevaParserRULE_interfaceDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(348)
		p.InterfaceDef()
	}
	p.SetState(353)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(349)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(350)
			p.InterfaceDef()
		}

		p.SetState(355)
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
	p.EnterRule(localctx, 46, nevaParserRULE_interfaceDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(357)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__7 {
		{
			p.SetState(356)
			p.Match(nevaParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(359)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(360)
		p.TypeParams()
	}
	{
		p.SetState(361)
		p.PortsDef()
	}
	{
		p.SetState(362)
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
	p.EnterRule(localctx, 48, nevaParserRULE_portsDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(364)
		p.Match(nevaParserT__15)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(365)
		p.PortDefList()
	}
	{
		p.SetState(366)
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
	p.EnterRule(localctx, 50, nevaParserRULE_portDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(368)
		p.PortDef()
	}
	p.SetState(376)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__9 {
		{
			p.SetState(369)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(371)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserNEWLINE {
			{
				p.SetState(370)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(373)
			p.PortDef()
		}

		p.SetState(378)
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
	p.EnterRule(localctx, 52, nevaParserRULE_portDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(379)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(380)
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
	p.EnterRule(localctx, 54, nevaParserRULE_constStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(382)
		p.Match(nevaParserT__17)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(383)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(384)
		p.ConstDefList()
	}
	{
		p.SetState(385)
		p.Match(nevaParserT__3)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(386)
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
		p.SetState(388)
		p.ConstDef()
	}
	p.SetState(393)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(389)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(390)
			p.ConstDef()
		}

		p.SetState(395)
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
	{
		p.SetState(400)
		p.TypeExpr()
	}
	{
		p.SetState(401)
		p.Match(nevaParserT__18)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(402)
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
	p.SetState(412)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__19:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(404)
			p.Match(nevaParserT__19)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserT__20:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(405)
			p.Match(nevaParserT__20)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserINT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(406)
			p.Match(nevaParserINT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserFLOAT:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(407)
			p.Match(nevaParserFLOAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserSTRING:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(408)
			p.Match(nevaParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserT__11:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(409)
			p.ArrLit()
		}

	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(410)
			p.RecLit()
		}

	case nevaParserT__21:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(411)
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
		p.SetState(414)
		p.Match(nevaParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(415)
		p.ArrItems()
	}
	{
		p.SetState(416)
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

	p.SetState(430)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 45, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(418)
			p.ConstValue()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(419)
			p.ConstValue()
		}
		p.SetState(427)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserT__9 {
			{
				p.SetState(420)
				p.Match(nevaParserT__9)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(422)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if _la == nevaParserNEWLINE {
				{
					p.SetState(421)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			}
			{
				p.SetState(424)
				p.ConstValue()
			}

			p.SetState(429)
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
		p.SetState(432)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(433)
		p.RecValueFields()
	}
	{
		p.SetState(434)
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
		p.SetState(436)
		p.RecValueField()
	}
	p.SetState(444)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__9 {
		{
			p.SetState(437)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(439)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserNEWLINE {
			{
				p.SetState(438)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(441)
			p.RecValueField()
		}

		p.SetState(446)
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
		p.SetState(447)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(448)
		p.Match(nevaParserT__22)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(449)
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
		p.SetState(451)
		p.Match(nevaParserT__23)
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
		p.CompDefList()
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
		p.SetState(457)
		p.CompDef()
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
			p.CompDef()
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
	p.SetState(466)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 49, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(465)
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
		p.SetState(468)
		p.InterfaceDef()
	}
	{
		p.SetState(469)
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
	p.SetState(476)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(471)
			p.Match(nevaParserT__2)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(472)
			p.CompNodesDef()
		}

	case nevaParserT__26:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(473)
			p.CompNetDef()
		}
		{
			p.SetState(474)
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
		p.SetState(478)
		p.Match(nevaParserT__24)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(479)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(480)
		p.CompNodeDefList()
	}
	{
		p.SetState(481)
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
	p.SetState(485)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(483)
			p.AbsNodeDef()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(484)
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
		p.SetState(487)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(488)
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
		p.SetState(490)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(491)
		p.Match(nevaParserT__18)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(492)
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
		p.SetState(494)
		p.NodeRef()
	}
	{
		p.SetState(495)
		p.NodeArgs()
	}
	{
		p.SetState(496)
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
		p.SetState(498)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(503)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__25 {
		{
			p.SetState(499)
			p.Match(nevaParserT__25)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(500)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(505)
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
		p.SetState(506)
		p.Match(nevaParserT__15)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(507)
		p.NodeArgList()
	}
	{
		p.SetState(508)
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
		p.SetState(510)
		p.NodeArg()
	}

	{
		p.SetState(511)
		p.Match(nevaParserT__9)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(513)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserNEWLINE {
		{
			p.SetState(512)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(515)
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
		p.SetState(517)
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
		p.SetState(519)
		p.Match(nevaParserT__26)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(520)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(521)
		p.ConnDefList()
	}
	{
		p.SetState(522)
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
		p.SetState(524)
		p.ConnDef()
	}
	p.SetState(529)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(525)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(526)
			p.ConnDef()
		}

		p.SetState(531)
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
		p.SetState(532)
		p.PortAddr()
	}
	{
		p.SetState(533)
		p.Match(nevaParserT__27)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(534)
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

	p.SetState(546)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 57, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(537)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserIDENTIFIER {
			{
				p.SetState(536)
				p.Match(nevaParserIDENTIFIER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(539)
			p.PortDirection()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(540)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(544)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserT__11 {
			{
				p.SetState(541)
				p.Match(nevaParserT__11)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(542)
				p.Match(nevaParserINT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(543)
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
		p.SetState(548)
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
	p.SetState(552)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__28, nevaParserT__29, nevaParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(550)
			p.PortAddr()
		}

	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(551)
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
		p.SetState(554)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(555)
		p.PortAddr()
	}
	p.SetState(560)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(556)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(557)
			p.PortAddr()
		}

		p.SetState(562)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(563)
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

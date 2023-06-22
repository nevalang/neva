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
		"constDef", "constVal", "arrLit", "vecItems", "recLit", "recValueFields",
		"recValueField", "compStmt", "compDefList", "compDef", "compBody", "compNodesDef",
		"compNodeDefList", "absNodeDef", "concreteNodeDef", "concreteNodeInst",
		"nodeRef", "nodeArgs", "nodeArgList", "nodeArg", "compNetDef", "connDefList",
		"connDef", "portAddr", "portDirection", "connReceiverSide", "connReceivers",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 36, 672, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 2, 53, 7, 53, 2, 54, 7, 54, 1, 0, 1, 0, 1, 0, 5, 0, 114, 8, 0, 10,
		0, 12, 0, 117, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 5, 1, 123, 8, 1, 10, 1, 12,
		1, 126, 9, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 133, 8, 2, 1, 3, 1, 3,
		1, 3, 5, 3, 138, 8, 3, 10, 3, 12, 3, 141, 9, 3, 1, 3, 5, 3, 144, 8, 3,
		10, 3, 12, 3, 147, 9, 3, 1, 3, 1, 3, 1, 4, 3, 4, 152, 8, 4, 1, 4, 1, 4,
		5, 4, 156, 8, 4, 10, 4, 12, 4, 159, 9, 4, 1, 5, 3, 5, 162, 8, 5, 1, 5,
		1, 5, 1, 5, 5, 5, 167, 8, 5, 10, 5, 12, 5, 170, 9, 5, 1, 6, 1, 6, 1, 6,
		5, 6, 175, 8, 6, 10, 6, 12, 6, 178, 9, 6, 1, 6, 5, 6, 181, 8, 6, 10, 6,
		12, 6, 184, 9, 6, 1, 6, 1, 6, 1, 7, 3, 7, 189, 8, 7, 1, 7, 1, 7, 3, 7,
		193, 8, 7, 1, 7, 1, 7, 5, 7, 197, 8, 7, 10, 7, 12, 7, 200, 9, 7, 1, 8,
		1, 8, 5, 8, 204, 8, 8, 10, 8, 12, 8, 207, 9, 8, 1, 8, 1, 8, 1, 8, 5, 8,
		212, 8, 8, 10, 8, 12, 8, 215, 9, 8, 1, 8, 5, 8, 218, 8, 8, 10, 8, 12, 8,
		221, 9, 8, 1, 8, 5, 8, 224, 8, 8, 10, 8, 12, 8, 227, 9, 8, 1, 8, 1, 8,
		1, 9, 1, 9, 3, 9, 233, 8, 9, 1, 10, 1, 10, 1, 10, 3, 10, 238, 8, 10, 1,
		11, 1, 11, 3, 11, 242, 8, 11, 1, 12, 1, 12, 5, 12, 246, 8, 12, 10, 12,
		12, 12, 249, 9, 12, 1, 12, 1, 12, 1, 12, 5, 12, 254, 8, 12, 10, 12, 12,
		12, 257, 9, 12, 1, 12, 5, 12, 260, 8, 12, 10, 12, 12, 12, 263, 9, 12, 1,
		12, 5, 12, 266, 8, 12, 10, 12, 12, 12, 269, 9, 12, 1, 12, 1, 12, 1, 13,
		1, 13, 1, 13, 3, 13, 276, 8, 13, 1, 14, 1, 14, 5, 14, 280, 8, 14, 10, 14,
		12, 14, 283, 9, 14, 1, 14, 1, 14, 1, 14, 5, 14, 288, 8, 14, 10, 14, 12,
		14, 291, 9, 14, 1, 14, 5, 14, 294, 8, 14, 10, 14, 12, 14, 297, 9, 14, 1,
		14, 5, 14, 300, 8, 14, 10, 14, 12, 14, 303, 9, 14, 1, 14, 1, 14, 1, 15,
		1, 15, 5, 15, 309, 8, 15, 10, 15, 12, 15, 312, 9, 15, 1, 15, 1, 15, 5,
		15, 316, 8, 15, 10, 15, 12, 15, 319, 9, 15, 1, 15, 1, 15, 1, 15, 1, 16,
		1, 16, 5, 16, 326, 8, 16, 10, 16, 12, 16, 329, 9, 16, 1, 16, 3, 16, 332,
		8, 16, 1, 16, 1, 16, 1, 17, 1, 17, 4, 17, 338, 8, 17, 11, 17, 12, 17, 339,
		1, 17, 5, 17, 343, 8, 17, 10, 17, 12, 17, 346, 9, 17, 1, 18, 1, 18, 1,
		18, 5, 18, 351, 8, 18, 10, 18, 12, 18, 354, 9, 18, 1, 19, 1, 19, 5, 19,
		358, 8, 19, 10, 19, 12, 19, 361, 9, 19, 1, 19, 1, 19, 5, 19, 365, 8, 19,
		10, 19, 12, 19, 368, 9, 19, 1, 19, 4, 19, 371, 8, 19, 11, 19, 12, 19, 372,
		1, 20, 1, 20, 3, 20, 377, 8, 20, 1, 21, 1, 21, 1, 21, 5, 21, 382, 8, 21,
		10, 21, 12, 21, 385, 9, 21, 1, 21, 5, 21, 388, 8, 21, 10, 21, 12, 21, 391,
		9, 21, 1, 21, 1, 21, 1, 22, 3, 22, 396, 8, 22, 1, 22, 1, 22, 3, 22, 400,
		8, 22, 1, 22, 1, 22, 1, 22, 5, 22, 405, 8, 22, 10, 22, 12, 22, 408, 9,
		22, 1, 23, 1, 23, 1, 24, 1, 24, 1, 25, 1, 25, 5, 25, 416, 8, 25, 10, 25,
		12, 25, 419, 9, 25, 1, 25, 3, 25, 422, 8, 25, 1, 25, 1, 25, 1, 25, 5, 25,
		427, 8, 25, 10, 25, 12, 25, 430, 9, 25, 3, 25, 432, 8, 25, 1, 25, 1, 25,
		1, 26, 5, 26, 437, 8, 26, 10, 26, 12, 26, 440, 9, 26, 1, 26, 1, 26, 1,
		26, 5, 26, 445, 8, 26, 10, 26, 12, 26, 448, 9, 26, 1, 27, 1, 27, 1, 27,
		5, 27, 453, 8, 27, 10, 27, 12, 27, 456, 9, 27, 1, 27, 5, 27, 459, 8, 27,
		10, 27, 12, 27, 462, 9, 27, 1, 27, 1, 27, 1, 28, 3, 28, 467, 8, 28, 1,
		28, 1, 28, 1, 28, 1, 28, 1, 28, 5, 28, 474, 8, 28, 10, 28, 12, 28, 477,
		9, 28, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 3, 29, 487,
		8, 29, 1, 30, 1, 30, 5, 30, 491, 8, 30, 10, 30, 12, 30, 494, 9, 30, 1,
		30, 3, 30, 497, 8, 30, 1, 30, 1, 30, 1, 31, 1, 31, 1, 31, 1, 31, 5, 31,
		505, 8, 31, 10, 31, 12, 31, 508, 9, 31, 1, 31, 1, 31, 5, 31, 512, 8, 31,
		10, 31, 12, 31, 515, 9, 31, 5, 31, 517, 8, 31, 10, 31, 12, 31, 520, 9,
		31, 3, 31, 522, 8, 31, 1, 32, 1, 32, 5, 32, 526, 8, 32, 10, 32, 12, 32,
		529, 9, 32, 1, 32, 3, 32, 532, 8, 32, 1, 32, 1, 32, 1, 33, 1, 33, 5, 33,
		538, 8, 33, 10, 33, 12, 33, 541, 9, 33, 1, 33, 5, 33, 544, 8, 33, 10, 33,
		12, 33, 547, 9, 33, 1, 34, 1, 34, 1, 34, 1, 34, 5, 34, 553, 8, 34, 10,
		34, 12, 34, 556, 9, 34, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 36,
		1, 36, 1, 36, 5, 36, 567, 8, 36, 10, 36, 12, 36, 570, 9, 36, 1, 37, 3,
		37, 573, 8, 37, 1, 37, 1, 37, 1, 37, 1, 38, 1, 38, 1, 38, 1, 38, 1, 38,
		3, 38, 583, 8, 38, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 40, 1, 40, 3,
		40, 592, 8, 40, 1, 41, 1, 41, 1, 41, 1, 42, 1, 42, 1, 42, 1, 42, 1, 43,
		1, 43, 1, 43, 1, 43, 1, 44, 1, 44, 1, 44, 5, 44, 608, 8, 44, 10, 44, 12,
		44, 611, 9, 44, 1, 45, 1, 45, 1, 45, 1, 45, 1, 46, 1, 46, 1, 46, 3, 46,
		620, 8, 46, 1, 46, 1, 46, 1, 47, 1, 47, 1, 48, 1, 48, 1, 48, 1, 48, 1,
		48, 1, 49, 1, 49, 1, 49, 5, 49, 634, 8, 49, 10, 49, 12, 49, 637, 9, 49,
		1, 50, 1, 50, 1, 50, 1, 50, 1, 51, 3, 51, 644, 8, 51, 1, 51, 1, 51, 1,
		51, 1, 51, 1, 51, 3, 51, 651, 8, 51, 3, 51, 653, 8, 51, 1, 52, 1, 52, 1,
		53, 1, 53, 3, 53, 659, 8, 53, 1, 54, 1, 54, 1, 54, 1, 54, 5, 54, 665, 8,
		54, 10, 54, 12, 54, 668, 9, 54, 1, 54, 1, 54, 1, 54, 0, 0, 55, 0, 2, 4,
		6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42,
		44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78,
		80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 0, 2,
		1, 0, 35, 35, 1, 0, 29, 30, 710, 0, 115, 1, 0, 0, 0, 2, 120, 1, 0, 0, 0,
		4, 132, 1, 0, 0, 0, 6, 134, 1, 0, 0, 0, 8, 151, 1, 0, 0, 0, 10, 161, 1,
		0, 0, 0, 12, 171, 1, 0, 0, 0, 14, 188, 1, 0, 0, 0, 16, 201, 1, 0, 0, 0,
		18, 230, 1, 0, 0, 0, 20, 237, 1, 0, 0, 0, 22, 239, 1, 0, 0, 0, 24, 243,
		1, 0, 0, 0, 26, 275, 1, 0, 0, 0, 28, 277, 1, 0, 0, 0, 30, 306, 1, 0, 0,
		0, 32, 323, 1, 0, 0, 0, 34, 335, 1, 0, 0, 0, 36, 347, 1, 0, 0, 0, 38, 355,
		1, 0, 0, 0, 40, 376, 1, 0, 0, 0, 42, 378, 1, 0, 0, 0, 44, 395, 1, 0, 0,
		0, 46, 409, 1, 0, 0, 0, 48, 411, 1, 0, 0, 0, 50, 413, 1, 0, 0, 0, 52, 438,
		1, 0, 0, 0, 54, 449, 1, 0, 0, 0, 56, 466, 1, 0, 0, 0, 58, 486, 1, 0, 0,
		0, 60, 488, 1, 0, 0, 0, 62, 521, 1, 0, 0, 0, 64, 523, 1, 0, 0, 0, 66, 535,
		1, 0, 0, 0, 68, 548, 1, 0, 0, 0, 70, 557, 1, 0, 0, 0, 72, 563, 1, 0, 0,
		0, 74, 572, 1, 0, 0, 0, 76, 582, 1, 0, 0, 0, 78, 584, 1, 0, 0, 0, 80, 591,
		1, 0, 0, 0, 82, 593, 1, 0, 0, 0, 84, 596, 1, 0, 0, 0, 86, 600, 1, 0, 0,
		0, 88, 604, 1, 0, 0, 0, 90, 612, 1, 0, 0, 0, 92, 616, 1, 0, 0, 0, 94, 623,
		1, 0, 0, 0, 96, 625, 1, 0, 0, 0, 98, 630, 1, 0, 0, 0, 100, 638, 1, 0, 0,
		0, 102, 652, 1, 0, 0, 0, 104, 654, 1, 0, 0, 0, 106, 658, 1, 0, 0, 0, 108,
		660, 1, 0, 0, 0, 110, 114, 5, 35, 0, 0, 111, 114, 3, 2, 1, 0, 112, 114,
		3, 4, 2, 0, 113, 110, 1, 0, 0, 0, 113, 111, 1, 0, 0, 0, 113, 112, 1, 0,
		0, 0, 114, 117, 1, 0, 0, 0, 115, 113, 1, 0, 0, 0, 115, 116, 1, 0, 0, 0,
		116, 118, 1, 0, 0, 0, 117, 115, 1, 0, 0, 0, 118, 119, 5, 0, 0, 1, 119,
		1, 1, 0, 0, 0, 120, 124, 5, 1, 0, 0, 121, 123, 8, 0, 0, 0, 122, 121, 1,
		0, 0, 0, 123, 126, 1, 0, 0, 0, 124, 122, 1, 0, 0, 0, 124, 125, 1, 0, 0,
		0, 125, 3, 1, 0, 0, 0, 126, 124, 1, 0, 0, 0, 127, 133, 3, 6, 3, 0, 128,
		133, 3, 12, 6, 0, 129, 133, 3, 42, 21, 0, 130, 133, 3, 54, 27, 0, 131,
		133, 3, 70, 35, 0, 132, 127, 1, 0, 0, 0, 132, 128, 1, 0, 0, 0, 132, 129,
		1, 0, 0, 0, 132, 130, 1, 0, 0, 0, 132, 131, 1, 0, 0, 0, 133, 5, 1, 0, 0,
		0, 134, 135, 5, 2, 0, 0, 135, 139, 5, 3, 0, 0, 136, 138, 5, 35, 0, 0, 137,
		136, 1, 0, 0, 0, 138, 141, 1, 0, 0, 0, 139, 137, 1, 0, 0, 0, 139, 140,
		1, 0, 0, 0, 140, 145, 1, 0, 0, 0, 141, 139, 1, 0, 0, 0, 142, 144, 3, 8,
		4, 0, 143, 142, 1, 0, 0, 0, 144, 147, 1, 0, 0, 0, 145, 143, 1, 0, 0, 0,
		145, 146, 1, 0, 0, 0, 146, 148, 1, 0, 0, 0, 147, 145, 1, 0, 0, 0, 148,
		149, 5, 4, 0, 0, 149, 7, 1, 0, 0, 0, 150, 152, 5, 31, 0, 0, 151, 150, 1,
		0, 0, 0, 151, 152, 1, 0, 0, 0, 152, 153, 1, 0, 0, 0, 153, 157, 3, 10, 5,
		0, 154, 156, 5, 35, 0, 0, 155, 154, 1, 0, 0, 0, 156, 159, 1, 0, 0, 0, 157,
		155, 1, 0, 0, 0, 157, 158, 1, 0, 0, 0, 158, 9, 1, 0, 0, 0, 159, 157, 1,
		0, 0, 0, 160, 162, 5, 5, 0, 0, 161, 160, 1, 0, 0, 0, 161, 162, 1, 0, 0,
		0, 162, 163, 1, 0, 0, 0, 163, 168, 5, 31, 0, 0, 164, 165, 5, 6, 0, 0, 165,
		167, 5, 31, 0, 0, 166, 164, 1, 0, 0, 0, 167, 170, 1, 0, 0, 0, 168, 166,
		1, 0, 0, 0, 168, 169, 1, 0, 0, 0, 169, 11, 1, 0, 0, 0, 170, 168, 1, 0,
		0, 0, 171, 172, 5, 7, 0, 0, 172, 176, 5, 3, 0, 0, 173, 175, 5, 35, 0, 0,
		174, 173, 1, 0, 0, 0, 175, 178, 1, 0, 0, 0, 176, 174, 1, 0, 0, 0, 176,
		177, 1, 0, 0, 0, 177, 182, 1, 0, 0, 0, 178, 176, 1, 0, 0, 0, 179, 181,
		3, 14, 7, 0, 180, 179, 1, 0, 0, 0, 181, 184, 1, 0, 0, 0, 182, 180, 1, 0,
		0, 0, 182, 183, 1, 0, 0, 0, 183, 185, 1, 0, 0, 0, 184, 182, 1, 0, 0, 0,
		185, 186, 5, 4, 0, 0, 186, 13, 1, 0, 0, 0, 187, 189, 5, 8, 0, 0, 188, 187,
		1, 0, 0, 0, 188, 189, 1, 0, 0, 0, 189, 190, 1, 0, 0, 0, 190, 192, 5, 31,
		0, 0, 191, 193, 3, 16, 8, 0, 192, 191, 1, 0, 0, 0, 192, 193, 1, 0, 0, 0,
		193, 194, 1, 0, 0, 0, 194, 198, 3, 20, 10, 0, 195, 197, 5, 35, 0, 0, 196,
		195, 1, 0, 0, 0, 197, 200, 1, 0, 0, 0, 198, 196, 1, 0, 0, 0, 198, 199,
		1, 0, 0, 0, 199, 15, 1, 0, 0, 0, 200, 198, 1, 0, 0, 0, 201, 205, 5, 9,
		0, 0, 202, 204, 5, 35, 0, 0, 203, 202, 1, 0, 0, 0, 204, 207, 1, 0, 0, 0,
		205, 203, 1, 0, 0, 0, 205, 206, 1, 0, 0, 0, 206, 208, 1, 0, 0, 0, 207,
		205, 1, 0, 0, 0, 208, 219, 3, 18, 9, 0, 209, 213, 5, 10, 0, 0, 210, 212,
		5, 35, 0, 0, 211, 210, 1, 0, 0, 0, 212, 215, 1, 0, 0, 0, 213, 211, 1, 0,
		0, 0, 213, 214, 1, 0, 0, 0, 214, 216, 1, 0, 0, 0, 215, 213, 1, 0, 0, 0,
		216, 218, 3, 18, 9, 0, 217, 209, 1, 0, 0, 0, 218, 221, 1, 0, 0, 0, 219,
		217, 1, 0, 0, 0, 219, 220, 1, 0, 0, 0, 220, 225, 1, 0, 0, 0, 221, 219,
		1, 0, 0, 0, 222, 224, 5, 35, 0, 0, 223, 222, 1, 0, 0, 0, 224, 227, 1, 0,
		0, 0, 225, 223, 1, 0, 0, 0, 225, 226, 1, 0, 0, 0, 226, 228, 1, 0, 0, 0,
		227, 225, 1, 0, 0, 0, 228, 229, 5, 11, 0, 0, 229, 17, 1, 0, 0, 0, 230,
		232, 5, 31, 0, 0, 231, 233, 3, 20, 10, 0, 232, 231, 1, 0, 0, 0, 232, 233,
		1, 0, 0, 0, 233, 19, 1, 0, 0, 0, 234, 238, 3, 22, 11, 0, 235, 238, 3, 26,
		13, 0, 236, 238, 3, 38, 19, 0, 237, 234, 1, 0, 0, 0, 237, 235, 1, 0, 0,
		0, 237, 236, 1, 0, 0, 0, 238, 21, 1, 0, 0, 0, 239, 241, 5, 31, 0, 0, 240,
		242, 3, 24, 12, 0, 241, 240, 1, 0, 0, 0, 241, 242, 1, 0, 0, 0, 242, 23,
		1, 0, 0, 0, 243, 247, 5, 9, 0, 0, 244, 246, 5, 35, 0, 0, 245, 244, 1, 0,
		0, 0, 246, 249, 1, 0, 0, 0, 247, 245, 1, 0, 0, 0, 247, 248, 1, 0, 0, 0,
		248, 250, 1, 0, 0, 0, 249, 247, 1, 0, 0, 0, 250, 261, 3, 20, 10, 0, 251,
		255, 5, 10, 0, 0, 252, 254, 5, 35, 0, 0, 253, 252, 1, 0, 0, 0, 254, 257,
		1, 0, 0, 0, 255, 253, 1, 0, 0, 0, 255, 256, 1, 0, 0, 0, 256, 258, 1, 0,
		0, 0, 257, 255, 1, 0, 0, 0, 258, 260, 3, 20, 10, 0, 259, 251, 1, 0, 0,
		0, 260, 263, 1, 0, 0, 0, 261, 259, 1, 0, 0, 0, 261, 262, 1, 0, 0, 0, 262,
		267, 1, 0, 0, 0, 263, 261, 1, 0, 0, 0, 264, 266, 5, 35, 0, 0, 265, 264,
		1, 0, 0, 0, 266, 269, 1, 0, 0, 0, 267, 265, 1, 0, 0, 0, 267, 268, 1, 0,
		0, 0, 268, 270, 1, 0, 0, 0, 269, 267, 1, 0, 0, 0, 270, 271, 5, 11, 0, 0,
		271, 25, 1, 0, 0, 0, 272, 276, 3, 28, 14, 0, 273, 276, 3, 30, 15, 0, 274,
		276, 3, 32, 16, 0, 275, 272, 1, 0, 0, 0, 275, 273, 1, 0, 0, 0, 275, 274,
		1, 0, 0, 0, 276, 27, 1, 0, 0, 0, 277, 281, 5, 3, 0, 0, 278, 280, 5, 35,
		0, 0, 279, 278, 1, 0, 0, 0, 280, 283, 1, 0, 0, 0, 281, 279, 1, 0, 0, 0,
		281, 282, 1, 0, 0, 0, 282, 284, 1, 0, 0, 0, 283, 281, 1, 0, 0, 0, 284,
		295, 5, 31, 0, 0, 285, 289, 5, 10, 0, 0, 286, 288, 5, 35, 0, 0, 287, 286,
		1, 0, 0, 0, 288, 291, 1, 0, 0, 0, 289, 287, 1, 0, 0, 0, 289, 290, 1, 0,
		0, 0, 290, 292, 1, 0, 0, 0, 291, 289, 1, 0, 0, 0, 292, 294, 5, 31, 0, 0,
		293, 285, 1, 0, 0, 0, 294, 297, 1, 0, 0, 0, 295, 293, 1, 0, 0, 0, 295,
		296, 1, 0, 0, 0, 296, 301, 1, 0, 0, 0, 297, 295, 1, 0, 0, 0, 298, 300,
		5, 35, 0, 0, 299, 298, 1, 0, 0, 0, 300, 303, 1, 0, 0, 0, 301, 299, 1, 0,
		0, 0, 301, 302, 1, 0, 0, 0, 302, 304, 1, 0, 0, 0, 303, 301, 1, 0, 0, 0,
		304, 305, 5, 4, 0, 0, 305, 29, 1, 0, 0, 0, 306, 310, 5, 12, 0, 0, 307,
		309, 5, 35, 0, 0, 308, 307, 1, 0, 0, 0, 309, 312, 1, 0, 0, 0, 310, 308,
		1, 0, 0, 0, 310, 311, 1, 0, 0, 0, 311, 313, 1, 0, 0, 0, 312, 310, 1, 0,
		0, 0, 313, 317, 5, 32, 0, 0, 314, 316, 5, 35, 0, 0, 315, 314, 1, 0, 0,
		0, 316, 319, 1, 0, 0, 0, 317, 315, 1, 0, 0, 0, 317, 318, 1, 0, 0, 0, 318,
		320, 1, 0, 0, 0, 319, 317, 1, 0, 0, 0, 320, 321, 5, 13, 0, 0, 321, 322,
		3, 20, 10, 0, 322, 31, 1, 0, 0, 0, 323, 327, 5, 3, 0, 0, 324, 326, 5, 35,
		0, 0, 325, 324, 1, 0, 0, 0, 326, 329, 1, 0, 0, 0, 327, 325, 1, 0, 0, 0,
		327, 328, 1, 0, 0, 0, 328, 331, 1, 0, 0, 0, 329, 327, 1, 0, 0, 0, 330,
		332, 3, 34, 17, 0, 331, 330, 1, 0, 0, 0, 331, 332, 1, 0, 0, 0, 332, 333,
		1, 0, 0, 0, 333, 334, 5, 4, 0, 0, 334, 33, 1, 0, 0, 0, 335, 344, 3, 36,
		18, 0, 336, 338, 5, 35, 0, 0, 337, 336, 1, 0, 0, 0, 338, 339, 1, 0, 0,
		0, 339, 337, 1, 0, 0, 0, 339, 340, 1, 0, 0, 0, 340, 341, 1, 0, 0, 0, 341,
		343, 3, 36, 18, 0, 342, 337, 1, 0, 0, 0, 343, 346, 1, 0, 0, 0, 344, 342,
		1, 0, 0, 0, 344, 345, 1, 0, 0, 0, 345, 35, 1, 0, 0, 0, 346, 344, 1, 0,
		0, 0, 347, 348, 5, 31, 0, 0, 348, 352, 3, 20, 10, 0, 349, 351, 5, 35, 0,
		0, 350, 349, 1, 0, 0, 0, 351, 354, 1, 0, 0, 0, 352, 350, 1, 0, 0, 0, 352,
		353, 1, 0, 0, 0, 353, 37, 1, 0, 0, 0, 354, 352, 1, 0, 0, 0, 355, 370, 3,
		40, 20, 0, 356, 358, 5, 35, 0, 0, 357, 356, 1, 0, 0, 0, 358, 361, 1, 0,
		0, 0, 359, 357, 1, 0, 0, 0, 359, 360, 1, 0, 0, 0, 360, 362, 1, 0, 0, 0,
		361, 359, 1, 0, 0, 0, 362, 366, 5, 14, 0, 0, 363, 365, 5, 35, 0, 0, 364,
		363, 1, 0, 0, 0, 365, 368, 1, 0, 0, 0, 366, 364, 1, 0, 0, 0, 366, 367,
		1, 0, 0, 0, 367, 369, 1, 0, 0, 0, 368, 366, 1, 0, 0, 0, 369, 371, 3, 40,
		20, 0, 370, 359, 1, 0, 0, 0, 371, 372, 1, 0, 0, 0, 372, 370, 1, 0, 0, 0,
		372, 373, 1, 0, 0, 0, 373, 39, 1, 0, 0, 0, 374, 377, 3, 22, 11, 0, 375,
		377, 3, 26, 13, 0, 376, 374, 1, 0, 0, 0, 376, 375, 1, 0, 0, 0, 377, 41,
		1, 0, 0, 0, 378, 379, 5, 15, 0, 0, 379, 383, 5, 3, 0, 0, 380, 382, 5, 35,
		0, 0, 381, 380, 1, 0, 0, 0, 382, 385, 1, 0, 0, 0, 383, 381, 1, 0, 0, 0,
		383, 384, 1, 0, 0, 0, 384, 389, 1, 0, 0, 0, 385, 383, 1, 0, 0, 0, 386,
		388, 3, 44, 22, 0, 387, 386, 1, 0, 0, 0, 388, 391, 1, 0, 0, 0, 389, 387,
		1, 0, 0, 0, 389, 390, 1, 0, 0, 0, 390, 392, 1, 0, 0, 0, 391, 389, 1, 0,
		0, 0, 392, 393, 5, 4, 0, 0, 393, 43, 1, 0, 0, 0, 394, 396, 5, 8, 0, 0,
		395, 394, 1, 0, 0, 0, 395, 396, 1, 0, 0, 0, 396, 397, 1, 0, 0, 0, 397,
		399, 5, 31, 0, 0, 398, 400, 3, 16, 8, 0, 399, 398, 1, 0, 0, 0, 399, 400,
		1, 0, 0, 0, 400, 401, 1, 0, 0, 0, 401, 402, 3, 46, 23, 0, 402, 406, 3,
		48, 24, 0, 403, 405, 5, 35, 0, 0, 404, 403, 1, 0, 0, 0, 405, 408, 1, 0,
		0, 0, 406, 404, 1, 0, 0, 0, 406, 407, 1, 0, 0, 0, 407, 45, 1, 0, 0, 0,
		408, 406, 1, 0, 0, 0, 409, 410, 3, 50, 25, 0, 410, 47, 1, 0, 0, 0, 411,
		412, 3, 50, 25, 0, 412, 49, 1, 0, 0, 0, 413, 431, 5, 16, 0, 0, 414, 416,
		5, 35, 0, 0, 415, 414, 1, 0, 0, 0, 416, 419, 1, 0, 0, 0, 417, 415, 1, 0,
		0, 0, 417, 418, 1, 0, 0, 0, 418, 432, 1, 0, 0, 0, 419, 417, 1, 0, 0, 0,
		420, 422, 3, 52, 26, 0, 421, 420, 1, 0, 0, 0, 421, 422, 1, 0, 0, 0, 422,
		432, 1, 0, 0, 0, 423, 428, 3, 52, 26, 0, 424, 425, 5, 10, 0, 0, 425, 427,
		3, 52, 26, 0, 426, 424, 1, 0, 0, 0, 427, 430, 1, 0, 0, 0, 428, 426, 1,
		0, 0, 0, 428, 429, 1, 0, 0, 0, 429, 432, 1, 0, 0, 0, 430, 428, 1, 0, 0,
		0, 431, 417, 1, 0, 0, 0, 431, 421, 1, 0, 0, 0, 431, 423, 1, 0, 0, 0, 432,
		433, 1, 0, 0, 0, 433, 434, 5, 17, 0, 0, 434, 51, 1, 0, 0, 0, 435, 437,
		5, 35, 0, 0, 436, 435, 1, 0, 0, 0, 437, 440, 1, 0, 0, 0, 438, 436, 1, 0,
		0, 0, 438, 439, 1, 0, 0, 0, 439, 441, 1, 0, 0, 0, 440, 438, 1, 0, 0, 0,
		441, 442, 5, 31, 0, 0, 442, 446, 3, 20, 10, 0, 443, 445, 5, 35, 0, 0, 444,
		443, 1, 0, 0, 0, 445, 448, 1, 0, 0, 0, 446, 444, 1, 0, 0, 0, 446, 447,
		1, 0, 0, 0, 447, 53, 1, 0, 0, 0, 448, 446, 1, 0, 0, 0, 449, 450, 5, 18,
		0, 0, 450, 454, 5, 3, 0, 0, 451, 453, 5, 35, 0, 0, 452, 451, 1, 0, 0, 0,
		453, 456, 1, 0, 0, 0, 454, 452, 1, 0, 0, 0, 454, 455, 1, 0, 0, 0, 455,
		460, 1, 0, 0, 0, 456, 454, 1, 0, 0, 0, 457, 459, 3, 56, 28, 0, 458, 457,
		1, 0, 0, 0, 459, 462, 1, 0, 0, 0, 460, 458, 1, 0, 0, 0, 460, 461, 1, 0,
		0, 0, 461, 463, 1, 0, 0, 0, 462, 460, 1, 0, 0, 0, 463, 464, 5, 4, 0, 0,
		464, 55, 1, 0, 0, 0, 465, 467, 5, 8, 0, 0, 466, 465, 1, 0, 0, 0, 466, 467,
		1, 0, 0, 0, 467, 468, 1, 0, 0, 0, 468, 469, 5, 31, 0, 0, 469, 470, 3, 20,
		10, 0, 470, 471, 5, 19, 0, 0, 471, 475, 3, 58, 29, 0, 472, 474, 5, 35,
		0, 0, 473, 472, 1, 0, 0, 0, 474, 477, 1, 0, 0, 0, 475, 473, 1, 0, 0, 0,
		475, 476, 1, 0, 0, 0, 476, 57, 1, 0, 0, 0, 477, 475, 1, 0, 0, 0, 478, 487,
		5, 20, 0, 0, 479, 487, 5, 21, 0, 0, 480, 487, 5, 32, 0, 0, 481, 487, 5,
		33, 0, 0, 482, 487, 5, 34, 0, 0, 483, 487, 3, 60, 30, 0, 484, 487, 3, 64,
		32, 0, 485, 487, 5, 22, 0, 0, 486, 478, 1, 0, 0, 0, 486, 479, 1, 0, 0,
		0, 486, 480, 1, 0, 0, 0, 486, 481, 1, 0, 0, 0, 486, 482, 1, 0, 0, 0, 486,
		483, 1, 0, 0, 0, 486, 484, 1, 0, 0, 0, 486, 485, 1, 0, 0, 0, 487, 59, 1,
		0, 0, 0, 488, 492, 5, 12, 0, 0, 489, 491, 5, 35, 0, 0, 490, 489, 1, 0,
		0, 0, 491, 494, 1, 0, 0, 0, 492, 490, 1, 0, 0, 0, 492, 493, 1, 0, 0, 0,
		493, 496, 1, 0, 0, 0, 494, 492, 1, 0, 0, 0, 495, 497, 3, 62, 31, 0, 496,
		495, 1, 0, 0, 0, 496, 497, 1, 0, 0, 0, 497, 498, 1, 0, 0, 0, 498, 499,
		5, 13, 0, 0, 499, 61, 1, 0, 0, 0, 500, 522, 3, 58, 29, 0, 501, 518, 3,
		58, 29, 0, 502, 506, 5, 10, 0, 0, 503, 505, 5, 35, 0, 0, 504, 503, 1, 0,
		0, 0, 505, 508, 1, 0, 0, 0, 506, 504, 1, 0, 0, 0, 506, 507, 1, 0, 0, 0,
		507, 509, 1, 0, 0, 0, 508, 506, 1, 0, 0, 0, 509, 513, 3, 58, 29, 0, 510,
		512, 5, 35, 0, 0, 511, 510, 1, 0, 0, 0, 512, 515, 1, 0, 0, 0, 513, 511,
		1, 0, 0, 0, 513, 514, 1, 0, 0, 0, 514, 517, 1, 0, 0, 0, 515, 513, 1, 0,
		0, 0, 516, 502, 1, 0, 0, 0, 517, 520, 1, 0, 0, 0, 518, 516, 1, 0, 0, 0,
		518, 519, 1, 0, 0, 0, 519, 522, 1, 0, 0, 0, 520, 518, 1, 0, 0, 0, 521,
		500, 1, 0, 0, 0, 521, 501, 1, 0, 0, 0, 522, 63, 1, 0, 0, 0, 523, 527, 5,
		3, 0, 0, 524, 526, 5, 35, 0, 0, 525, 524, 1, 0, 0, 0, 526, 529, 1, 0, 0,
		0, 527, 525, 1, 0, 0, 0, 527, 528, 1, 0, 0, 0, 528, 531, 1, 0, 0, 0, 529,
		527, 1, 0, 0, 0, 530, 532, 3, 66, 33, 0, 531, 530, 1, 0, 0, 0, 531, 532,
		1, 0, 0, 0, 532, 533, 1, 0, 0, 0, 533, 534, 5, 4, 0, 0, 534, 65, 1, 0,
		0, 0, 535, 545, 3, 68, 34, 0, 536, 538, 5, 35, 0, 0, 537, 536, 1, 0, 0,
		0, 538, 541, 1, 0, 0, 0, 539, 537, 1, 0, 0, 0, 539, 540, 1, 0, 0, 0, 540,
		542, 1, 0, 0, 0, 541, 539, 1, 0, 0, 0, 542, 544, 3, 68, 34, 0, 543, 539,
		1, 0, 0, 0, 544, 547, 1, 0, 0, 0, 545, 543, 1, 0, 0, 0, 545, 546, 1, 0,
		0, 0, 546, 67, 1, 0, 0, 0, 547, 545, 1, 0, 0, 0, 548, 549, 5, 31, 0, 0,
		549, 550, 5, 23, 0, 0, 550, 554, 3, 58, 29, 0, 551, 553, 5, 35, 0, 0, 552,
		551, 1, 0, 0, 0, 553, 556, 1, 0, 0, 0, 554, 552, 1, 0, 0, 0, 554, 555,
		1, 0, 0, 0, 555, 69, 1, 0, 0, 0, 556, 554, 1, 0, 0, 0, 557, 558, 5, 24,
		0, 0, 558, 559, 5, 3, 0, 0, 559, 560, 3, 72, 36, 0, 560, 561, 5, 4, 0,
		0, 561, 562, 5, 35, 0, 0, 562, 71, 1, 0, 0, 0, 563, 568, 3, 74, 37, 0,
		564, 565, 5, 35, 0, 0, 565, 567, 3, 74, 37, 0, 566, 564, 1, 0, 0, 0, 567,
		570, 1, 0, 0, 0, 568, 566, 1, 0, 0, 0, 568, 569, 1, 0, 0, 0, 569, 73, 1,
		0, 0, 0, 570, 568, 1, 0, 0, 0, 571, 573, 5, 8, 0, 0, 572, 571, 1, 0, 0,
		0, 572, 573, 1, 0, 0, 0, 573, 574, 1, 0, 0, 0, 574, 575, 3, 44, 22, 0,
		575, 576, 3, 76, 38, 0, 576, 75, 1, 0, 0, 0, 577, 578, 5, 3, 0, 0, 578,
		583, 3, 78, 39, 0, 579, 580, 3, 96, 48, 0, 580, 581, 5, 4, 0, 0, 581, 583,
		1, 0, 0, 0, 582, 577, 1, 0, 0, 0, 582, 579, 1, 0, 0, 0, 583, 77, 1, 0,
		0, 0, 584, 585, 5, 25, 0, 0, 585, 586, 5, 3, 0, 0, 586, 587, 3, 80, 40,
		0, 587, 588, 5, 4, 0, 0, 588, 79, 1, 0, 0, 0, 589, 592, 3, 82, 41, 0, 590,
		592, 3, 84, 42, 0, 591, 589, 1, 0, 0, 0, 591, 590, 1, 0, 0, 0, 592, 81,
		1, 0, 0, 0, 593, 594, 5, 31, 0, 0, 594, 595, 3, 22, 11, 0, 595, 83, 1,
		0, 0, 0, 596, 597, 5, 31, 0, 0, 597, 598, 5, 19, 0, 0, 598, 599, 3, 86,
		43, 0, 599, 85, 1, 0, 0, 0, 600, 601, 3, 88, 44, 0, 601, 602, 3, 90, 45,
		0, 602, 603, 3, 24, 12, 0, 603, 87, 1, 0, 0, 0, 604, 609, 5, 31, 0, 0,
		605, 606, 5, 26, 0, 0, 606, 608, 5, 31, 0, 0, 607, 605, 1, 0, 0, 0, 608,
		611, 1, 0, 0, 0, 609, 607, 1, 0, 0, 0, 609, 610, 1, 0, 0, 0, 610, 89, 1,
		0, 0, 0, 611, 609, 1, 0, 0, 0, 612, 613, 5, 16, 0, 0, 613, 614, 3, 92,
		46, 0, 614, 615, 5, 17, 0, 0, 615, 91, 1, 0, 0, 0, 616, 617, 3, 94, 47,
		0, 617, 619, 5, 10, 0, 0, 618, 620, 5, 35, 0, 0, 619, 618, 1, 0, 0, 0,
		619, 620, 1, 0, 0, 0, 620, 621, 1, 0, 0, 0, 621, 622, 3, 94, 47, 0, 622,
		93, 1, 0, 0, 0, 623, 624, 3, 86, 43, 0, 624, 95, 1, 0, 0, 0, 625, 626,
		5, 27, 0, 0, 626, 627, 5, 3, 0, 0, 627, 628, 3, 98, 49, 0, 628, 629, 5,
		4, 0, 0, 629, 97, 1, 0, 0, 0, 630, 635, 3, 100, 50, 0, 631, 632, 5, 35,
		0, 0, 632, 634, 3, 100, 50, 0, 633, 631, 1, 0, 0, 0, 634, 637, 1, 0, 0,
		0, 635, 633, 1, 0, 0, 0, 635, 636, 1, 0, 0, 0, 636, 99, 1, 0, 0, 0, 637,
		635, 1, 0, 0, 0, 638, 639, 3, 102, 51, 0, 639, 640, 5, 28, 0, 0, 640, 641,
		3, 106, 53, 0, 641, 101, 1, 0, 0, 0, 642, 644, 5, 31, 0, 0, 643, 642, 1,
		0, 0, 0, 643, 644, 1, 0, 0, 0, 644, 645, 1, 0, 0, 0, 645, 653, 3, 104,
		52, 0, 646, 650, 5, 31, 0, 0, 647, 648, 5, 12, 0, 0, 648, 649, 5, 32, 0,
		0, 649, 651, 5, 13, 0, 0, 650, 647, 1, 0, 0, 0, 650, 651, 1, 0, 0, 0, 651,
		653, 1, 0, 0, 0, 652, 643, 1, 0, 0, 0, 652, 646, 1, 0, 0, 0, 653, 103,
		1, 0, 0, 0, 654, 655, 7, 1, 0, 0, 655, 105, 1, 0, 0, 0, 656, 659, 3, 102,
		51, 0, 657, 659, 3, 108, 54, 0, 658, 656, 1, 0, 0, 0, 658, 657, 1, 0, 0,
		0, 659, 107, 1, 0, 0, 0, 660, 661, 5, 3, 0, 0, 661, 666, 3, 102, 51, 0,
		662, 663, 5, 35, 0, 0, 663, 665, 3, 102, 51, 0, 664, 662, 1, 0, 0, 0, 665,
		668, 1, 0, 0, 0, 666, 664, 1, 0, 0, 0, 666, 667, 1, 0, 0, 0, 667, 669,
		1, 0, 0, 0, 668, 666, 1, 0, 0, 0, 669, 670, 5, 4, 0, 0, 670, 109, 1, 0,
		0, 0, 81, 113, 115, 124, 132, 139, 145, 151, 157, 161, 168, 176, 182, 188,
		192, 198, 205, 213, 219, 225, 232, 237, 241, 247, 255, 261, 267, 275, 281,
		289, 295, 301, 310, 317, 327, 331, 339, 344, 352, 359, 366, 372, 376, 383,
		389, 395, 399, 406, 417, 421, 428, 431, 438, 446, 454, 460, 466, 475, 486,
		492, 496, 506, 513, 518, 521, 527, 531, 539, 545, 554, 568, 572, 582, 591,
		609, 619, 635, 643, 650, 652, 658, 666,
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
	nevaParserRULE_constDef         = 28
	nevaParserRULE_constVal         = 29
	nevaParserRULE_arrLit           = 30
	nevaParserRULE_vecItems         = 31
	nevaParserRULE_recLit           = 32
	nevaParserRULE_recValueFields   = 33
	nevaParserRULE_recValueField    = 34
	nevaParserRULE_compStmt         = 35
	nevaParserRULE_compDefList      = 36
	nevaParserRULE_compDef          = 37
	nevaParserRULE_compBody         = 38
	nevaParserRULE_compNodesDef     = 39
	nevaParserRULE_compNodeDefList  = 40
	nevaParserRULE_absNodeDef       = 41
	nevaParserRULE_concreteNodeDef  = 42
	nevaParserRULE_concreteNodeInst = 43
	nevaParserRULE_nodeRef          = 44
	nevaParserRULE_nodeArgs         = 45
	nevaParserRULE_nodeArgList      = 46
	nevaParserRULE_nodeArg          = 47
	nevaParserRULE_compNetDef       = 48
	nevaParserRULE_connDefList      = 49
	nevaParserRULE_connDef          = 50
	nevaParserRULE_portAddr         = 51
	nevaParserRULE_portDirection    = 52
	nevaParserRULE_connReceiverSide = 53
	nevaParserRULE_connReceivers    = 54
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
	p.SetState(115)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&34376810630) != 0 {
		p.SetState(113)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case nevaParserNEWLINE:
			{
				p.SetState(110)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case nevaParserT__0:
			{
				p.SetState(111)
				p.Comment()
			}

		case nevaParserT__1, nevaParserT__6, nevaParserT__14, nevaParserT__17, nevaParserT__23:
			{
				p.SetState(112)
				p.Stmt()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(117)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(118)
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
		p.SetState(120)
		p.Match(nevaParserT__0)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(124)
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
				p.SetState(121)
				_la = p.GetTokenStream().LA(1)

				if _la <= 0 || _la == nevaParserNEWLINE {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
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
	p.SetState(132)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(127)
			p.UseStmt()
		}

	case nevaParserT__6:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(128)
			p.TypeStmt()
		}

	case nevaParserT__14:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(129)
			p.IoStmt()
		}

	case nevaParserT__17:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(130)
			p.ConstStmt()
		}

	case nevaParserT__23:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(131)
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
		p.SetState(134)
		p.Match(nevaParserT__1)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(135)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(139)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(136)
			p.Match(nevaParserNEWLINE)
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
	}
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__4 || _la == nevaParserIDENTIFIER {
		{
			p.SetState(142)
			p.ImportDef()
		}

		p.SetState(147)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(148)
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
	p.SetState(151)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(150)
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
		p.SetState(153)
		p.ImportPath()
	}
	p.SetState(157)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(154)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(159)
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
	p.SetState(161)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__4 {
		{
			p.SetState(160)
			p.Match(nevaParserT__4)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(163)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(168)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__5 {
		{
			p.SetState(164)
			p.Match(nevaParserT__5)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
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
		p.SetState(171)
		p.Match(nevaParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(172)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(176)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(173)
			p.Match(nevaParserNEWLINE)
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
	}
	p.SetState(182)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 || _la == nevaParserIDENTIFIER {
		{
			p.SetState(179)
			p.TypeDef()
		}

		p.SetState(184)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(185)
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
	p.SetState(188)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__7 {
		{
			p.SetState(187)
			p.Match(nevaParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(190)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(192)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__8 {
		{
			p.SetState(191)
			p.TypeParams()
		}

	}
	{
		p.SetState(194)
		p.TypeExpr()
	}
	p.SetState(198)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(195)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(200)
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
		p.SetState(201)
		p.Match(nevaParserT__8)
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

	for _la == nevaParserNEWLINE {
		{
			p.SetState(202)
			p.Match(nevaParserNEWLINE)
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
	}
	{
		p.SetState(208)
		p.TypeParam()
	}
	p.SetState(219)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__9 {
		{
			p.SetState(209)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(213)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(210)
				p.Match(nevaParserNEWLINE)
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
		}
		{
			p.SetState(216)
			p.TypeParam()
		}

		p.SetState(221)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(225)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(222)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(227)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(228)
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
		p.SetState(230)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(232)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2147487752) != 0 {
		{
			p.SetState(231)
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
	p.SetState(237)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(234)
			p.TypeInstExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(235)
			p.TypeLitExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(236)
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
		p.SetState(239)
		p.Match(nevaParserIDENTIFIER)
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

	if _la == nevaParserT__8 {
		{
			p.SetState(240)
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
		p.SetState(243)
		p.Match(nevaParserT__8)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(247)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(244)
			p.Match(nevaParserNEWLINE)
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
	}
	{
		p.SetState(250)
		p.TypeExpr()
	}
	p.SetState(261)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__9 {
		{
			p.SetState(251)
			p.Match(nevaParserT__9)
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

		for _la == nevaParserNEWLINE {
			{
				p.SetState(252)
				p.Match(nevaParserNEWLINE)
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
	}
	p.SetState(267)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(264)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(269)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(270)
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
	p.SetState(275)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(272)
			p.EnumTypeExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(273)
			p.ArrTypeExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(274)
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
		p.SetState(277)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(281)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(278)
			p.Match(nevaParserNEWLINE)
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
	}
	{
		p.SetState(284)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(295)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__9 {
		{
			p.SetState(285)
			p.Match(nevaParserT__9)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(289)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(286)
				p.Match(nevaParserNEWLINE)
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
		}
		{
			p.SetState(292)
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
	}
	p.SetState(301)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(298)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(303)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(304)
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
		p.SetState(306)
		p.Match(nevaParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(310)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(307)
			p.Match(nevaParserNEWLINE)
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
	}
	{
		p.SetState(313)
		p.Match(nevaParserINT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(317)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(314)
			p.Match(nevaParserNEWLINE)
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
	}
	{
		p.SetState(320)
		p.Match(nevaParserT__12)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(321)
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
		p.SetState(323)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
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
	p.SetState(331)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserIDENTIFIER {
		{
			p.SetState(330)
			p.RecFields()
		}

	}
	{
		p.SetState(333)
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
		p.SetState(335)
		p.RecField()
	}
	p.SetState(344)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		p.SetState(337)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == nevaParserNEWLINE {
			{
				p.SetState(336)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(339)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(341)
			p.RecField()
		}

		p.SetState(346)
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
		p.SetState(347)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(348)
		p.TypeExpr()
	}
	p.SetState(352)
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
				p.SetState(349)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

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
		p.SetState(355)
		p.NonUnionTypeExpr()
	}
	p.SetState(370)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			p.SetState(359)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(356)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(361)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(362)
				p.Match(nevaParserT__13)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(366)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(363)
					p.Match(nevaParserNEWLINE)
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
			}
			{
				p.SetState(369)
				p.NonUnionTypeExpr()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(372)
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
	p.SetState(376)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(374)
			p.TypeInstExpr()
		}

	case nevaParserT__2, nevaParserT__11:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(375)
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
		p.SetState(378)
		p.Match(nevaParserT__14)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(379)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(383)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(380)
			p.Match(nevaParserNEWLINE)
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
	}
	p.SetState(389)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 || _la == nevaParserIDENTIFIER {
		{
			p.SetState(386)
			p.InterfaceDef()
		}

		p.SetState(391)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(392)
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
	p.SetState(395)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__7 {
		{
			p.SetState(394)
			p.Match(nevaParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(397)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(399)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__8 {
		{
			p.SetState(398)
			p.TypeParams()
		}

	}
	{
		p.SetState(401)
		p.InPortsDef()
	}
	{
		p.SetState(402)
		p.OutPortsDef()
	}
	p.SetState(406)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(403)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(408)
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
		p.SetState(409)
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
		p.SetState(413)
		p.Match(nevaParserT__15)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(431)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 50, p.GetParserRuleContext()) {
	case 1:
		p.SetState(417)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(414)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(419)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		p.SetState(421)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserIDENTIFIER || _la == nevaParserNEWLINE {
			{
				p.SetState(420)
				p.PortAndType()
			}

		}

	case 3:
		{
			p.SetState(423)
			p.PortAndType()
		}
		p.SetState(428)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserT__9 {
			{
				p.SetState(424)
				p.Match(nevaParserT__9)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
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
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	{
		p.SetState(433)
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
	p.SetState(438)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(435)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(440)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(441)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(442)
		p.TypeExpr()
	}
	p.SetState(446)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(443)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(448)
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
	AllConstDef() []IConstDefContext
	ConstDef(i int) IConstDefContext

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

func (s *ConstStmtContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ConstStmtContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *ConstStmtContext) AllConstDef() []IConstDefContext {
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

func (s *ConstStmtContext) ConstDef(i int) IConstDefContext {
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
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(449)
		p.Match(nevaParserT__17)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(450)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(454)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(451)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(456)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(460)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 || _la == nevaParserIDENTIFIER {
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
	}
	{
		p.SetState(463)
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

// IConstDefContext is an interface to support dynamic dispatch.
type IConstDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TypeExpr() ITypeExprContext
	ConstVal() IConstValContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *ConstDefContext) ConstVal() IConstValContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstValContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstValContext)
}

func (s *ConstDefContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ConstDefContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.EnterRule(localctx, 56, nevaParserRULE_constDef)
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
		p.ConstVal()
	}
	p.SetState(475)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(472)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(477)
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

// IConstValContext is an interface to support dynamic dispatch.
type IConstValContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INT() antlr.TerminalNode
	FLOAT() antlr.TerminalNode
	STRING() antlr.TerminalNode
	ArrLit() IArrLitContext
	RecLit() IRecLitContext

	// IsConstValContext differentiates from other interfaces.
	IsConstValContext()
}

type ConstValContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConstValContext() *ConstValContext {
	var p = new(ConstValContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_constVal
	return p
}

func InitEmptyConstValContext(p *ConstValContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_constVal
}

func (*ConstValContext) IsConstValContext() {}

func NewConstValContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstValContext {
	var p = new(ConstValContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_constVal

	return p
}

func (s *ConstValContext) GetParser() antlr.Parser { return s.parser }

func (s *ConstValContext) INT() antlr.TerminalNode {
	return s.GetToken(nevaParserINT, 0)
}

func (s *ConstValContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(nevaParserFLOAT, 0)
}

func (s *ConstValContext) STRING() antlr.TerminalNode {
	return s.GetToken(nevaParserSTRING, 0)
}

func (s *ConstValContext) ArrLit() IArrLitContext {
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

func (s *ConstValContext) RecLit() IRecLitContext {
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

func (s *ConstValContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConstValContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConstValContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterConstVal(s)
	}
}

func (s *ConstValContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitConstVal(s)
	}
}

func (p *nevaParser) ConstVal() (localctx IConstValContext) {
	localctx = NewConstValContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, nevaParserRULE_constVal)
	p.SetState(486)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__19:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(478)
			p.Match(nevaParserT__19)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserT__20:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(479)
			p.Match(nevaParserT__20)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserINT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(480)
			p.Match(nevaParserINT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserFLOAT:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(481)
			p.Match(nevaParserFLOAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserSTRING:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(482)
			p.Match(nevaParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserT__11:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(483)
			p.ArrLit()
		}

	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(484)
			p.RecLit()
		}

	case nevaParserT__21:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(485)
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
	VecItems() IVecItemsContext

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

func (s *ArrLitContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *ArrLitContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *ArrLitContext) VecItems() IVecItemsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVecItemsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVecItemsContext)
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
	p.EnterRule(localctx, 60, nevaParserRULE_arrLit)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(488)
		p.Match(nevaParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(492)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(489)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(494)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(496)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&30072115208) != 0 {
		{
			p.SetState(495)
			p.VecItems()
		}

	}
	{
		p.SetState(498)
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

// IVecItemsContext is an interface to support dynamic dispatch.
type IVecItemsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllConstVal() []IConstValContext
	ConstVal(i int) IConstValContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsVecItemsContext differentiates from other interfaces.
	IsVecItemsContext()
}

type VecItemsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVecItemsContext() *VecItemsContext {
	var p = new(VecItemsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_vecItems
	return p
}

func InitEmptyVecItemsContext(p *VecItemsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_vecItems
}

func (*VecItemsContext) IsVecItemsContext() {}

func NewVecItemsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VecItemsContext {
	var p = new(VecItemsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_vecItems

	return p
}

func (s *VecItemsContext) GetParser() antlr.Parser { return s.parser }

func (s *VecItemsContext) AllConstVal() []IConstValContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IConstValContext); ok {
			len++
		}
	}

	tst := make([]IConstValContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IConstValContext); ok {
			tst[i] = t.(IConstValContext)
			i++
		}
	}

	return tst
}

func (s *VecItemsContext) ConstVal(i int) IConstValContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstValContext); ok {
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

	return t.(IConstValContext)
}

func (s *VecItemsContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *VecItemsContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *VecItemsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VecItemsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VecItemsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterVecItems(s)
	}
}

func (s *VecItemsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitVecItems(s)
	}
}

func (p *nevaParser) VecItems() (localctx IVecItemsContext) {
	localctx = NewVecItemsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, nevaParserRULE_vecItems)
	var _la int

	p.SetState(521)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 63, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(500)
			p.ConstVal()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(501)
			p.ConstVal()
		}
		p.SetState(518)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserT__9 {
			{
				p.SetState(502)
				p.Match(nevaParserT__9)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(506)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(503)
					p.Match(nevaParserNEWLINE)
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
			}
			{
				p.SetState(509)
				p.ConstVal()
			}
			p.SetState(513)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(510)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(515)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

			p.SetState(520)
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
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

func (s *RecLitContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *RecLitContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

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
	p.EnterRule(localctx, 64, nevaParserRULE_recLit)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(523)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(527)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(524)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(529)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(531)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserIDENTIFIER {
		{
			p.SetState(530)
			p.RecValueFields()
		}

	}
	{
		p.SetState(533)
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
	p.EnterRule(localctx, 66, nevaParserRULE_recValueFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(535)
		p.RecValueField()
	}
	p.SetState(545)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserIDENTIFIER || _la == nevaParserNEWLINE {
		p.SetState(539)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(536)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(541)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(542)
			p.RecValueField()
		}

		p.SetState(547)
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
	ConstVal() IConstValContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *RecValueFieldContext) ConstVal() IConstValContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstValContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstValContext)
}

func (s *RecValueFieldContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *RecValueFieldContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	p.EnterRule(localctx, 68, nevaParserRULE_recValueField)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(548)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(549)
		p.Match(nevaParserT__22)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(550)
		p.ConstVal()
	}
	p.SetState(554)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 68, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(551)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		p.SetState(556)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 68, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 70, nevaParserRULE_compStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(557)
		p.Match(nevaParserT__23)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(558)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(559)
		p.CompDefList()
	}
	{
		p.SetState(560)
		p.Match(nevaParserT__3)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(561)
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
	p.EnterRule(localctx, 72, nevaParserRULE_compDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(563)
		p.CompDef()
	}
	p.SetState(568)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(564)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(565)
			p.CompDef()
		}

		p.SetState(570)
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
	p.EnterRule(localctx, 74, nevaParserRULE_compDef)
	p.EnterOuterAlt(localctx, 1)
	p.SetState(572)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 70, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(571)
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
		p.SetState(574)
		p.InterfaceDef()
	}
	{
		p.SetState(575)
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
	p.EnterRule(localctx, 76, nevaParserRULE_compBody)
	p.SetState(582)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(577)
			p.Match(nevaParserT__2)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(578)
			p.CompNodesDef()
		}

	case nevaParserT__26:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(579)
			p.CompNetDef()
		}
		{
			p.SetState(580)
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
	p.EnterRule(localctx, 78, nevaParserRULE_compNodesDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(584)
		p.Match(nevaParserT__24)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(585)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(586)
		p.CompNodeDefList()
	}
	{
		p.SetState(587)
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
	p.EnterRule(localctx, 80, nevaParserRULE_compNodeDefList)
	p.SetState(591)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 72, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(589)
			p.AbsNodeDef()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(590)
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
	p.EnterRule(localctx, 82, nevaParserRULE_absNodeDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(593)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(594)
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
	p.EnterRule(localctx, 84, nevaParserRULE_concreteNodeDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(596)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(597)
		p.Match(nevaParserT__18)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(598)
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
	p.EnterRule(localctx, 86, nevaParserRULE_concreteNodeInst)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(600)
		p.NodeRef()
	}
	{
		p.SetState(601)
		p.NodeArgs()
	}
	{
		p.SetState(602)
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
	p.EnterRule(localctx, 88, nevaParserRULE_nodeRef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(604)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(609)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__25 {
		{
			p.SetState(605)
			p.Match(nevaParserT__25)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(606)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(611)
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
	p.EnterRule(localctx, 90, nevaParserRULE_nodeArgs)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(612)
		p.Match(nevaParserT__15)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(613)
		p.NodeArgList()
	}
	{
		p.SetState(614)
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
	p.EnterRule(localctx, 92, nevaParserRULE_nodeArgList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(616)
		p.NodeArg()
	}

	{
		p.SetState(617)
		p.Match(nevaParserT__9)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(619)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserNEWLINE {
		{
			p.SetState(618)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(621)
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
	p.EnterRule(localctx, 94, nevaParserRULE_nodeArg)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(623)
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
	p.EnterRule(localctx, 96, nevaParserRULE_compNetDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(625)
		p.Match(nevaParserT__26)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(626)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(627)
		p.ConnDefList()
	}
	{
		p.SetState(628)
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
	p.EnterRule(localctx, 98, nevaParserRULE_connDefList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(630)
		p.ConnDef()
	}
	p.SetState(635)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(631)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(632)
			p.ConnDef()
		}

		p.SetState(637)
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
	p.EnterRule(localctx, 100, nevaParserRULE_connDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(638)
		p.PortAddr()
	}
	{
		p.SetState(639)
		p.Match(nevaParserT__27)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(640)
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
	p.EnterRule(localctx, 102, nevaParserRULE_portAddr)
	var _la int

	p.SetState(652)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 78, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(643)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserIDENTIFIER {
			{
				p.SetState(642)
				p.Match(nevaParserIDENTIFIER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(645)
			p.PortDirection()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(646)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(650)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserT__11 {
			{
				p.SetState(647)
				p.Match(nevaParserT__11)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(648)
				p.Match(nevaParserINT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(649)
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
	p.EnterRule(localctx, 104, nevaParserRULE_portDirection)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(654)
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
	p.EnterRule(localctx, 106, nevaParserRULE_connReceiverSide)
	p.SetState(658)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__28, nevaParserT__29, nevaParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(656)
			p.PortAddr()
		}

	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(657)
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
	p.EnterRule(localctx, 108, nevaParserRULE_connReceivers)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(660)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(661)
		p.PortAddr()
	}
	p.SetState(666)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(662)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(663)
			p.PortAddr()
		}

		p.SetState(668)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(669)
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

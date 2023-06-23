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
		"'<'", "'>'", "','", "'['", "']'", "'|'", "'io'", "'('", "')'", "'const'",
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
		"typeDef", "typeParams", "typeParamList", "typeParam", "typeExpr", "typeInstExpr",
		"typeArgs", "typeLitExpr", "enumTypeExpr", "arrTypeExpr", "recTypeExpr",
		"recFields", "recField", "unionTypeExpr", "nonUnionTypeExpr", "ioStmt",
		"interfaceDef", "inPortsDef", "outPortsDef", "portsDef", "portAndType",
		"constStmt", "constDef", "constVal", "arrLit", "vecItems", "recLit",
		"recValueFields", "recValueField", "compStmt", "compDef", "compBody",
		"compNodesDef", "compNodeDefList", "absNodeDef", "concreteNodeDef",
		"concreteNodeInst", "nodeRef", "nodeArgs", "nodeArgList", "nodeArg",
		"compNetDef", "connDefList", "connDef", "portAddr", "portDirection",
		"connReceiverSide", "connReceivers",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 36, 743, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
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
		5, 3, 137, 8, 3, 10, 3, 12, 3, 140, 9, 3, 1, 3, 1, 3, 5, 3, 144, 8, 3,
		10, 3, 12, 3, 147, 9, 3, 1, 3, 5, 3, 150, 8, 3, 10, 3, 12, 3, 153, 9, 3,
		1, 3, 1, 3, 1, 4, 3, 4, 158, 8, 4, 1, 4, 1, 4, 5, 4, 162, 8, 4, 10, 4,
		12, 4, 165, 9, 4, 1, 5, 3, 5, 168, 8, 5, 1, 5, 1, 5, 1, 5, 5, 5, 173, 8,
		5, 10, 5, 12, 5, 176, 9, 5, 1, 6, 1, 6, 5, 6, 180, 8, 6, 10, 6, 12, 6,
		183, 9, 6, 1, 6, 1, 6, 5, 6, 187, 8, 6, 10, 6, 12, 6, 190, 9, 6, 1, 6,
		3, 6, 193, 8, 6, 1, 6, 1, 6, 5, 6, 197, 8, 6, 10, 6, 12, 6, 200, 9, 6,
		5, 6, 202, 8, 6, 10, 6, 12, 6, 205, 9, 6, 1, 6, 1, 6, 1, 7, 1, 7, 3, 7,
		211, 8, 7, 1, 7, 1, 7, 1, 8, 1, 8, 5, 8, 217, 8, 8, 10, 8, 12, 8, 220,
		9, 8, 1, 8, 3, 8, 223, 8, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 5, 9, 230, 8,
		9, 10, 9, 12, 9, 233, 9, 9, 1, 9, 1, 9, 5, 9, 237, 8, 9, 10, 9, 12, 9,
		240, 9, 9, 5, 9, 242, 8, 9, 10, 9, 12, 9, 245, 9, 9, 1, 10, 1, 10, 3, 10,
		249, 8, 10, 1, 11, 1, 11, 1, 11, 3, 11, 254, 8, 11, 1, 12, 1, 12, 3, 12,
		258, 8, 12, 1, 13, 1, 13, 5, 13, 262, 8, 13, 10, 13, 12, 13, 265, 9, 13,
		1, 13, 1, 13, 1, 13, 5, 13, 270, 8, 13, 10, 13, 12, 13, 273, 9, 13, 1,
		13, 5, 13, 276, 8, 13, 10, 13, 12, 13, 279, 9, 13, 1, 13, 5, 13, 282, 8,
		13, 10, 13, 12, 13, 285, 9, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 3, 14,
		292, 8, 14, 1, 15, 1, 15, 5, 15, 296, 8, 15, 10, 15, 12, 15, 299, 9, 15,
		1, 15, 1, 15, 1, 15, 5, 15, 304, 8, 15, 10, 15, 12, 15, 307, 9, 15, 1,
		15, 5, 15, 310, 8, 15, 10, 15, 12, 15, 313, 9, 15, 1, 15, 5, 15, 316, 8,
		15, 10, 15, 12, 15, 319, 9, 15, 1, 15, 1, 15, 1, 16, 1, 16, 5, 16, 325,
		8, 16, 10, 16, 12, 16, 328, 9, 16, 1, 16, 1, 16, 5, 16, 332, 8, 16, 10,
		16, 12, 16, 335, 9, 16, 1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 5, 17, 342,
		8, 17, 10, 17, 12, 17, 345, 9, 17, 1, 17, 3, 17, 348, 8, 17, 1, 17, 1,
		17, 1, 18, 1, 18, 4, 18, 354, 8, 18, 11, 18, 12, 18, 355, 1, 18, 5, 18,
		359, 8, 18, 10, 18, 12, 18, 362, 9, 18, 1, 19, 1, 19, 1, 19, 5, 19, 367,
		8, 19, 10, 19, 12, 19, 370, 9, 19, 1, 20, 1, 20, 5, 20, 374, 8, 20, 10,
		20, 12, 20, 377, 9, 20, 1, 20, 1, 20, 5, 20, 381, 8, 20, 10, 20, 12, 20,
		384, 9, 20, 1, 20, 4, 20, 387, 8, 20, 11, 20, 12, 20, 388, 1, 21, 1, 21,
		3, 21, 393, 8, 21, 1, 22, 1, 22, 5, 22, 397, 8, 22, 10, 22, 12, 22, 400,
		9, 22, 1, 22, 1, 22, 5, 22, 404, 8, 22, 10, 22, 12, 22, 407, 9, 22, 1,
		22, 3, 22, 410, 8, 22, 1, 22, 5, 22, 413, 8, 22, 10, 22, 12, 22, 416, 9,
		22, 1, 22, 1, 22, 1, 23, 1, 23, 3, 23, 422, 8, 23, 1, 23, 1, 23, 1, 23,
		5, 23, 427, 8, 23, 10, 23, 12, 23, 430, 9, 23, 1, 24, 1, 24, 1, 25, 1,
		25, 1, 26, 1, 26, 5, 26, 438, 8, 26, 10, 26, 12, 26, 441, 9, 26, 1, 26,
		3, 26, 444, 8, 26, 1, 26, 1, 26, 1, 26, 5, 26, 449, 8, 26, 10, 26, 12,
		26, 452, 9, 26, 3, 26, 454, 8, 26, 1, 26, 1, 26, 1, 27, 5, 27, 459, 8,
		27, 10, 27, 12, 27, 462, 9, 27, 1, 27, 1, 27, 1, 27, 5, 27, 467, 8, 27,
		10, 27, 12, 27, 470, 9, 27, 1, 28, 1, 28, 5, 28, 474, 8, 28, 10, 28, 12,
		28, 477, 9, 28, 1, 28, 1, 28, 5, 28, 481, 8, 28, 10, 28, 12, 28, 484, 9,
		28, 1, 28, 3, 28, 487, 8, 28, 1, 28, 5, 28, 490, 8, 28, 10, 28, 12, 28,
		493, 9, 28, 1, 28, 1, 28, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 5, 29, 502,
		8, 29, 10, 29, 12, 29, 505, 9, 29, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1,
		30, 1, 30, 1, 30, 3, 30, 515, 8, 30, 1, 31, 1, 31, 5, 31, 519, 8, 31, 10,
		31, 12, 31, 522, 9, 31, 1, 31, 3, 31, 525, 8, 31, 1, 31, 1, 31, 1, 32,
		1, 32, 1, 32, 1, 32, 5, 32, 533, 8, 32, 10, 32, 12, 32, 536, 9, 32, 1,
		32, 1, 32, 5, 32, 540, 8, 32, 10, 32, 12, 32, 543, 9, 32, 5, 32, 545, 8,
		32, 10, 32, 12, 32, 548, 9, 32, 3, 32, 550, 8, 32, 1, 33, 1, 33, 5, 33,
		554, 8, 33, 10, 33, 12, 33, 557, 9, 33, 1, 33, 3, 33, 560, 8, 33, 1, 33,
		1, 33, 1, 34, 1, 34, 5, 34, 566, 8, 34, 10, 34, 12, 34, 569, 9, 34, 1,
		34, 5, 34, 572, 8, 34, 10, 34, 12, 34, 575, 9, 34, 1, 35, 1, 35, 1, 35,
		1, 35, 5, 35, 581, 8, 35, 10, 35, 12, 35, 584, 9, 35, 1, 36, 1, 36, 5,
		36, 588, 8, 36, 10, 36, 12, 36, 591, 9, 36, 1, 36, 1, 36, 5, 36, 595, 8,
		36, 10, 36, 12, 36, 598, 9, 36, 1, 36, 3, 36, 601, 8, 36, 1, 36, 5, 36,
		604, 8, 36, 10, 36, 12, 36, 607, 9, 36, 1, 36, 1, 36, 1, 37, 1, 37, 1,
		37, 5, 37, 614, 8, 37, 10, 37, 12, 37, 617, 9, 37, 1, 38, 1, 38, 5, 38,
		621, 8, 38, 10, 38, 12, 38, 624, 9, 38, 1, 38, 1, 38, 3, 38, 628, 8, 38,
		1, 38, 1, 38, 1, 39, 1, 39, 5, 39, 634, 8, 39, 10, 39, 12, 39, 637, 9,
		39, 1, 39, 1, 39, 5, 39, 641, 8, 39, 10, 39, 12, 39, 644, 9, 39, 1, 39,
		1, 39, 1, 39, 1, 40, 1, 40, 3, 40, 651, 8, 40, 1, 41, 1, 41, 1, 41, 1,
		42, 1, 42, 1, 42, 1, 42, 1, 43, 1, 43, 1, 43, 1, 43, 1, 44, 1, 44, 1, 44,
		5, 44, 667, 8, 44, 10, 44, 12, 44, 670, 9, 44, 1, 45, 1, 45, 1, 45, 1,
		45, 1, 46, 1, 46, 1, 46, 3, 46, 679, 8, 46, 1, 46, 1, 46, 1, 47, 1, 47,
		1, 48, 1, 48, 5, 48, 687, 8, 48, 10, 48, 12, 48, 690, 9, 48, 1, 48, 1,
		48, 5, 48, 694, 8, 48, 10, 48, 12, 48, 697, 9, 48, 1, 48, 1, 48, 1, 48,
		1, 49, 1, 49, 1, 49, 5, 49, 705, 8, 49, 10, 49, 12, 49, 708, 9, 49, 1,
		50, 1, 50, 1, 50, 1, 50, 1, 51, 3, 51, 715, 8, 51, 1, 51, 1, 51, 1, 51,
		1, 51, 1, 51, 3, 51, 722, 8, 51, 3, 51, 724, 8, 51, 1, 52, 1, 52, 1, 53,
		1, 53, 3, 53, 730, 8, 53, 1, 54, 1, 54, 1, 54, 1, 54, 5, 54, 736, 8, 54,
		10, 54, 12, 54, 739, 9, 54, 1, 54, 1, 54, 1, 54, 0, 0, 55, 0, 2, 4, 6,
		8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42,
		44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78,
		80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 0, 2,
		1, 0, 35, 35, 1, 0, 29, 30, 795, 0, 115, 1, 0, 0, 0, 2, 120, 1, 0, 0, 0,
		4, 132, 1, 0, 0, 0, 6, 134, 1, 0, 0, 0, 8, 157, 1, 0, 0, 0, 10, 167, 1,
		0, 0, 0, 12, 177, 1, 0, 0, 0, 14, 208, 1, 0, 0, 0, 16, 214, 1, 0, 0, 0,
		18, 226, 1, 0, 0, 0, 20, 246, 1, 0, 0, 0, 22, 253, 1, 0, 0, 0, 24, 255,
		1, 0, 0, 0, 26, 259, 1, 0, 0, 0, 28, 291, 1, 0, 0, 0, 30, 293, 1, 0, 0,
		0, 32, 322, 1, 0, 0, 0, 34, 339, 1, 0, 0, 0, 36, 351, 1, 0, 0, 0, 38, 363,
		1, 0, 0, 0, 40, 371, 1, 0, 0, 0, 42, 392, 1, 0, 0, 0, 44, 394, 1, 0, 0,
		0, 46, 419, 1, 0, 0, 0, 48, 431, 1, 0, 0, 0, 50, 433, 1, 0, 0, 0, 52, 435,
		1, 0, 0, 0, 54, 460, 1, 0, 0, 0, 56, 471, 1, 0, 0, 0, 58, 496, 1, 0, 0,
		0, 60, 514, 1, 0, 0, 0, 62, 516, 1, 0, 0, 0, 64, 549, 1, 0, 0, 0, 66, 551,
		1, 0, 0, 0, 68, 563, 1, 0, 0, 0, 70, 576, 1, 0, 0, 0, 72, 585, 1, 0, 0,
		0, 74, 610, 1, 0, 0, 0, 76, 618, 1, 0, 0, 0, 78, 631, 1, 0, 0, 0, 80, 650,
		1, 0, 0, 0, 82, 652, 1, 0, 0, 0, 84, 655, 1, 0, 0, 0, 86, 659, 1, 0, 0,
		0, 88, 663, 1, 0, 0, 0, 90, 671, 1, 0, 0, 0, 92, 675, 1, 0, 0, 0, 94, 682,
		1, 0, 0, 0, 96, 684, 1, 0, 0, 0, 98, 701, 1, 0, 0, 0, 100, 709, 1, 0, 0,
		0, 102, 723, 1, 0, 0, 0, 104, 725, 1, 0, 0, 0, 106, 729, 1, 0, 0, 0, 108,
		731, 1, 0, 0, 0, 110, 114, 5, 35, 0, 0, 111, 114, 3, 2, 1, 0, 112, 114,
		3, 4, 2, 0, 113, 110, 1, 0, 0, 0, 113, 111, 1, 0, 0, 0, 113, 112, 1, 0,
		0, 0, 114, 117, 1, 0, 0, 0, 115, 113, 1, 0, 0, 0, 115, 116, 1, 0, 0, 0,
		116, 118, 1, 0, 0, 0, 117, 115, 1, 0, 0, 0, 118, 119, 5, 0, 0, 1, 119,
		1, 1, 0, 0, 0, 120, 124, 5, 1, 0, 0, 121, 123, 8, 0, 0, 0, 122, 121, 1,
		0, 0, 0, 123, 126, 1, 0, 0, 0, 124, 122, 1, 0, 0, 0, 124, 125, 1, 0, 0,
		0, 125, 3, 1, 0, 0, 0, 126, 124, 1, 0, 0, 0, 127, 133, 3, 6, 3, 0, 128,
		133, 3, 12, 6, 0, 129, 133, 3, 44, 22, 0, 130, 133, 3, 56, 28, 0, 131,
		133, 3, 72, 36, 0, 132, 127, 1, 0, 0, 0, 132, 128, 1, 0, 0, 0, 132, 129,
		1, 0, 0, 0, 132, 130, 1, 0, 0, 0, 132, 131, 1, 0, 0, 0, 133, 5, 1, 0, 0,
		0, 134, 138, 5, 2, 0, 0, 135, 137, 5, 35, 0, 0, 136, 135, 1, 0, 0, 0, 137,
		140, 1, 0, 0, 0, 138, 136, 1, 0, 0, 0, 138, 139, 1, 0, 0, 0, 139, 141,
		1, 0, 0, 0, 140, 138, 1, 0, 0, 0, 141, 145, 5, 3, 0, 0, 142, 144, 5, 35,
		0, 0, 143, 142, 1, 0, 0, 0, 144, 147, 1, 0, 0, 0, 145, 143, 1, 0, 0, 0,
		145, 146, 1, 0, 0, 0, 146, 151, 1, 0, 0, 0, 147, 145, 1, 0, 0, 0, 148,
		150, 3, 8, 4, 0, 149, 148, 1, 0, 0, 0, 150, 153, 1, 0, 0, 0, 151, 149,
		1, 0, 0, 0, 151, 152, 1, 0, 0, 0, 152, 154, 1, 0, 0, 0, 153, 151, 1, 0,
		0, 0, 154, 155, 5, 4, 0, 0, 155, 7, 1, 0, 0, 0, 156, 158, 5, 31, 0, 0,
		157, 156, 1, 0, 0, 0, 157, 158, 1, 0, 0, 0, 158, 159, 1, 0, 0, 0, 159,
		163, 3, 10, 5, 0, 160, 162, 5, 35, 0, 0, 161, 160, 1, 0, 0, 0, 162, 165,
		1, 0, 0, 0, 163, 161, 1, 0, 0, 0, 163, 164, 1, 0, 0, 0, 164, 9, 1, 0, 0,
		0, 165, 163, 1, 0, 0, 0, 166, 168, 5, 5, 0, 0, 167, 166, 1, 0, 0, 0, 167,
		168, 1, 0, 0, 0, 168, 169, 1, 0, 0, 0, 169, 174, 5, 31, 0, 0, 170, 171,
		5, 6, 0, 0, 171, 173, 5, 31, 0, 0, 172, 170, 1, 0, 0, 0, 173, 176, 1, 0,
		0, 0, 174, 172, 1, 0, 0, 0, 174, 175, 1, 0, 0, 0, 175, 11, 1, 0, 0, 0,
		176, 174, 1, 0, 0, 0, 177, 181, 5, 7, 0, 0, 178, 180, 5, 35, 0, 0, 179,
		178, 1, 0, 0, 0, 180, 183, 1, 0, 0, 0, 181, 179, 1, 0, 0, 0, 181, 182,
		1, 0, 0, 0, 182, 184, 1, 0, 0, 0, 183, 181, 1, 0, 0, 0, 184, 188, 5, 3,
		0, 0, 185, 187, 5, 35, 0, 0, 186, 185, 1, 0, 0, 0, 187, 190, 1, 0, 0, 0,
		188, 186, 1, 0, 0, 0, 188, 189, 1, 0, 0, 0, 189, 203, 1, 0, 0, 0, 190,
		188, 1, 0, 0, 0, 191, 193, 5, 8, 0, 0, 192, 191, 1, 0, 0, 0, 192, 193,
		1, 0, 0, 0, 193, 194, 1, 0, 0, 0, 194, 198, 3, 14, 7, 0, 195, 197, 5, 35,
		0, 0, 196, 195, 1, 0, 0, 0, 197, 200, 1, 0, 0, 0, 198, 196, 1, 0, 0, 0,
		198, 199, 1, 0, 0, 0, 199, 202, 1, 0, 0, 0, 200, 198, 1, 0, 0, 0, 201,
		192, 1, 0, 0, 0, 202, 205, 1, 0, 0, 0, 203, 201, 1, 0, 0, 0, 203, 204,
		1, 0, 0, 0, 204, 206, 1, 0, 0, 0, 205, 203, 1, 0, 0, 0, 206, 207, 5, 4,
		0, 0, 207, 13, 1, 0, 0, 0, 208, 210, 5, 31, 0, 0, 209, 211, 3, 16, 8, 0,
		210, 209, 1, 0, 0, 0, 210, 211, 1, 0, 0, 0, 211, 212, 1, 0, 0, 0, 212,
		213, 3, 22, 11, 0, 213, 15, 1, 0, 0, 0, 214, 218, 5, 9, 0, 0, 215, 217,
		5, 35, 0, 0, 216, 215, 1, 0, 0, 0, 217, 220, 1, 0, 0, 0, 218, 216, 1, 0,
		0, 0, 218, 219, 1, 0, 0, 0, 219, 222, 1, 0, 0, 0, 220, 218, 1, 0, 0, 0,
		221, 223, 3, 18, 9, 0, 222, 221, 1, 0, 0, 0, 222, 223, 1, 0, 0, 0, 223,
		224, 1, 0, 0, 0, 224, 225, 5, 10, 0, 0, 225, 17, 1, 0, 0, 0, 226, 243,
		3, 20, 10, 0, 227, 231, 5, 11, 0, 0, 228, 230, 5, 35, 0, 0, 229, 228, 1,
		0, 0, 0, 230, 233, 1, 0, 0, 0, 231, 229, 1, 0, 0, 0, 231, 232, 1, 0, 0,
		0, 232, 234, 1, 0, 0, 0, 233, 231, 1, 0, 0, 0, 234, 238, 3, 20, 10, 0,
		235, 237, 5, 35, 0, 0, 236, 235, 1, 0, 0, 0, 237, 240, 1, 0, 0, 0, 238,
		236, 1, 0, 0, 0, 238, 239, 1, 0, 0, 0, 239, 242, 1, 0, 0, 0, 240, 238,
		1, 0, 0, 0, 241, 227, 1, 0, 0, 0, 242, 245, 1, 0, 0, 0, 243, 241, 1, 0,
		0, 0, 243, 244, 1, 0, 0, 0, 244, 19, 1, 0, 0, 0, 245, 243, 1, 0, 0, 0,
		246, 248, 5, 31, 0, 0, 247, 249, 3, 22, 11, 0, 248, 247, 1, 0, 0, 0, 248,
		249, 1, 0, 0, 0, 249, 21, 1, 0, 0, 0, 250, 254, 3, 24, 12, 0, 251, 254,
		3, 28, 14, 0, 252, 254, 3, 40, 20, 0, 253, 250, 1, 0, 0, 0, 253, 251, 1,
		0, 0, 0, 253, 252, 1, 0, 0, 0, 254, 23, 1, 0, 0, 0, 255, 257, 5, 31, 0,
		0, 256, 258, 3, 26, 13, 0, 257, 256, 1, 0, 0, 0, 257, 258, 1, 0, 0, 0,
		258, 25, 1, 0, 0, 0, 259, 263, 5, 9, 0, 0, 260, 262, 5, 35, 0, 0, 261,
		260, 1, 0, 0, 0, 262, 265, 1, 0, 0, 0, 263, 261, 1, 0, 0, 0, 263, 264,
		1, 0, 0, 0, 264, 266, 1, 0, 0, 0, 265, 263, 1, 0, 0, 0, 266, 277, 3, 22,
		11, 0, 267, 271, 5, 11, 0, 0, 268, 270, 5, 35, 0, 0, 269, 268, 1, 0, 0,
		0, 270, 273, 1, 0, 0, 0, 271, 269, 1, 0, 0, 0, 271, 272, 1, 0, 0, 0, 272,
		274, 1, 0, 0, 0, 273, 271, 1, 0, 0, 0, 274, 276, 3, 22, 11, 0, 275, 267,
		1, 0, 0, 0, 276, 279, 1, 0, 0, 0, 277, 275, 1, 0, 0, 0, 277, 278, 1, 0,
		0, 0, 278, 283, 1, 0, 0, 0, 279, 277, 1, 0, 0, 0, 280, 282, 5, 35, 0, 0,
		281, 280, 1, 0, 0, 0, 282, 285, 1, 0, 0, 0, 283, 281, 1, 0, 0, 0, 283,
		284, 1, 0, 0, 0, 284, 286, 1, 0, 0, 0, 285, 283, 1, 0, 0, 0, 286, 287,
		5, 10, 0, 0, 287, 27, 1, 0, 0, 0, 288, 292, 3, 30, 15, 0, 289, 292, 3,
		32, 16, 0, 290, 292, 3, 34, 17, 0, 291, 288, 1, 0, 0, 0, 291, 289, 1, 0,
		0, 0, 291, 290, 1, 0, 0, 0, 292, 29, 1, 0, 0, 0, 293, 297, 5, 3, 0, 0,
		294, 296, 5, 35, 0, 0, 295, 294, 1, 0, 0, 0, 296, 299, 1, 0, 0, 0, 297,
		295, 1, 0, 0, 0, 297, 298, 1, 0, 0, 0, 298, 300, 1, 0, 0, 0, 299, 297,
		1, 0, 0, 0, 300, 311, 5, 31, 0, 0, 301, 305, 5, 11, 0, 0, 302, 304, 5,
		35, 0, 0, 303, 302, 1, 0, 0, 0, 304, 307, 1, 0, 0, 0, 305, 303, 1, 0, 0,
		0, 305, 306, 1, 0, 0, 0, 306, 308, 1, 0, 0, 0, 307, 305, 1, 0, 0, 0, 308,
		310, 5, 31, 0, 0, 309, 301, 1, 0, 0, 0, 310, 313, 1, 0, 0, 0, 311, 309,
		1, 0, 0, 0, 311, 312, 1, 0, 0, 0, 312, 317, 1, 0, 0, 0, 313, 311, 1, 0,
		0, 0, 314, 316, 5, 35, 0, 0, 315, 314, 1, 0, 0, 0, 316, 319, 1, 0, 0, 0,
		317, 315, 1, 0, 0, 0, 317, 318, 1, 0, 0, 0, 318, 320, 1, 0, 0, 0, 319,
		317, 1, 0, 0, 0, 320, 321, 5, 4, 0, 0, 321, 31, 1, 0, 0, 0, 322, 326, 5,
		12, 0, 0, 323, 325, 5, 35, 0, 0, 324, 323, 1, 0, 0, 0, 325, 328, 1, 0,
		0, 0, 326, 324, 1, 0, 0, 0, 326, 327, 1, 0, 0, 0, 327, 329, 1, 0, 0, 0,
		328, 326, 1, 0, 0, 0, 329, 333, 5, 32, 0, 0, 330, 332, 5, 35, 0, 0, 331,
		330, 1, 0, 0, 0, 332, 335, 1, 0, 0, 0, 333, 331, 1, 0, 0, 0, 333, 334,
		1, 0, 0, 0, 334, 336, 1, 0, 0, 0, 335, 333, 1, 0, 0, 0, 336, 337, 5, 13,
		0, 0, 337, 338, 3, 22, 11, 0, 338, 33, 1, 0, 0, 0, 339, 343, 5, 3, 0, 0,
		340, 342, 5, 35, 0, 0, 341, 340, 1, 0, 0, 0, 342, 345, 1, 0, 0, 0, 343,
		341, 1, 0, 0, 0, 343, 344, 1, 0, 0, 0, 344, 347, 1, 0, 0, 0, 345, 343,
		1, 0, 0, 0, 346, 348, 3, 36, 18, 0, 347, 346, 1, 0, 0, 0, 347, 348, 1,
		0, 0, 0, 348, 349, 1, 0, 0, 0, 349, 350, 5, 4, 0, 0, 350, 35, 1, 0, 0,
		0, 351, 360, 3, 38, 19, 0, 352, 354, 5, 35, 0, 0, 353, 352, 1, 0, 0, 0,
		354, 355, 1, 0, 0, 0, 355, 353, 1, 0, 0, 0, 355, 356, 1, 0, 0, 0, 356,
		357, 1, 0, 0, 0, 357, 359, 3, 38, 19, 0, 358, 353, 1, 0, 0, 0, 359, 362,
		1, 0, 0, 0, 360, 358, 1, 0, 0, 0, 360, 361, 1, 0, 0, 0, 361, 37, 1, 0,
		0, 0, 362, 360, 1, 0, 0, 0, 363, 364, 5, 31, 0, 0, 364, 368, 3, 22, 11,
		0, 365, 367, 5, 35, 0, 0, 366, 365, 1, 0, 0, 0, 367, 370, 1, 0, 0, 0, 368,
		366, 1, 0, 0, 0, 368, 369, 1, 0, 0, 0, 369, 39, 1, 0, 0, 0, 370, 368, 1,
		0, 0, 0, 371, 386, 3, 42, 21, 0, 372, 374, 5, 35, 0, 0, 373, 372, 1, 0,
		0, 0, 374, 377, 1, 0, 0, 0, 375, 373, 1, 0, 0, 0, 375, 376, 1, 0, 0, 0,
		376, 378, 1, 0, 0, 0, 377, 375, 1, 0, 0, 0, 378, 382, 5, 14, 0, 0, 379,
		381, 5, 35, 0, 0, 380, 379, 1, 0, 0, 0, 381, 384, 1, 0, 0, 0, 382, 380,
		1, 0, 0, 0, 382, 383, 1, 0, 0, 0, 383, 385, 1, 0, 0, 0, 384, 382, 1, 0,
		0, 0, 385, 387, 3, 42, 21, 0, 386, 375, 1, 0, 0, 0, 387, 388, 1, 0, 0,
		0, 388, 386, 1, 0, 0, 0, 388, 389, 1, 0, 0, 0, 389, 41, 1, 0, 0, 0, 390,
		393, 3, 24, 12, 0, 391, 393, 3, 28, 14, 0, 392, 390, 1, 0, 0, 0, 392, 391,
		1, 0, 0, 0, 393, 43, 1, 0, 0, 0, 394, 398, 5, 15, 0, 0, 395, 397, 5, 35,
		0, 0, 396, 395, 1, 0, 0, 0, 397, 400, 1, 0, 0, 0, 398, 396, 1, 0, 0, 0,
		398, 399, 1, 0, 0, 0, 399, 401, 1, 0, 0, 0, 400, 398, 1, 0, 0, 0, 401,
		405, 5, 3, 0, 0, 402, 404, 5, 35, 0, 0, 403, 402, 1, 0, 0, 0, 404, 407,
		1, 0, 0, 0, 405, 403, 1, 0, 0, 0, 405, 406, 1, 0, 0, 0, 406, 414, 1, 0,
		0, 0, 407, 405, 1, 0, 0, 0, 408, 410, 5, 8, 0, 0, 409, 408, 1, 0, 0, 0,
		409, 410, 1, 0, 0, 0, 410, 411, 1, 0, 0, 0, 411, 413, 3, 46, 23, 0, 412,
		409, 1, 0, 0, 0, 413, 416, 1, 0, 0, 0, 414, 412, 1, 0, 0, 0, 414, 415,
		1, 0, 0, 0, 415, 417, 1, 0, 0, 0, 416, 414, 1, 0, 0, 0, 417, 418, 5, 4,
		0, 0, 418, 45, 1, 0, 0, 0, 419, 421, 5, 31, 0, 0, 420, 422, 3, 16, 8, 0,
		421, 420, 1, 0, 0, 0, 421, 422, 1, 0, 0, 0, 422, 423, 1, 0, 0, 0, 423,
		424, 3, 48, 24, 0, 424, 428, 3, 50, 25, 0, 425, 427, 5, 35, 0, 0, 426,
		425, 1, 0, 0, 0, 427, 430, 1, 0, 0, 0, 428, 426, 1, 0, 0, 0, 428, 429,
		1, 0, 0, 0, 429, 47, 1, 0, 0, 0, 430, 428, 1, 0, 0, 0, 431, 432, 3, 52,
		26, 0, 432, 49, 1, 0, 0, 0, 433, 434, 3, 52, 26, 0, 434, 51, 1, 0, 0, 0,
		435, 453, 5, 16, 0, 0, 436, 438, 5, 35, 0, 0, 437, 436, 1, 0, 0, 0, 438,
		441, 1, 0, 0, 0, 439, 437, 1, 0, 0, 0, 439, 440, 1, 0, 0, 0, 440, 454,
		1, 0, 0, 0, 441, 439, 1, 0, 0, 0, 442, 444, 3, 54, 27, 0, 443, 442, 1,
		0, 0, 0, 443, 444, 1, 0, 0, 0, 444, 454, 1, 0, 0, 0, 445, 450, 3, 54, 27,
		0, 446, 447, 5, 11, 0, 0, 447, 449, 3, 54, 27, 0, 448, 446, 1, 0, 0, 0,
		449, 452, 1, 0, 0, 0, 450, 448, 1, 0, 0, 0, 450, 451, 1, 0, 0, 0, 451,
		454, 1, 0, 0, 0, 452, 450, 1, 0, 0, 0, 453, 439, 1, 0, 0, 0, 453, 443,
		1, 0, 0, 0, 453, 445, 1, 0, 0, 0, 454, 455, 1, 0, 0, 0, 455, 456, 5, 17,
		0, 0, 456, 53, 1, 0, 0, 0, 457, 459, 5, 35, 0, 0, 458, 457, 1, 0, 0, 0,
		459, 462, 1, 0, 0, 0, 460, 458, 1, 0, 0, 0, 460, 461, 1, 0, 0, 0, 461,
		463, 1, 0, 0, 0, 462, 460, 1, 0, 0, 0, 463, 464, 5, 31, 0, 0, 464, 468,
		3, 22, 11, 0, 465, 467, 5, 35, 0, 0, 466, 465, 1, 0, 0, 0, 467, 470, 1,
		0, 0, 0, 468, 466, 1, 0, 0, 0, 468, 469, 1, 0, 0, 0, 469, 55, 1, 0, 0,
		0, 470, 468, 1, 0, 0, 0, 471, 475, 5, 18, 0, 0, 472, 474, 5, 35, 0, 0,
		473, 472, 1, 0, 0, 0, 474, 477, 1, 0, 0, 0, 475, 473, 1, 0, 0, 0, 475,
		476, 1, 0, 0, 0, 476, 478, 1, 0, 0, 0, 477, 475, 1, 0, 0, 0, 478, 482,
		5, 3, 0, 0, 479, 481, 5, 35, 0, 0, 480, 479, 1, 0, 0, 0, 481, 484, 1, 0,
		0, 0, 482, 480, 1, 0, 0, 0, 482, 483, 1, 0, 0, 0, 483, 491, 1, 0, 0, 0,
		484, 482, 1, 0, 0, 0, 485, 487, 5, 8, 0, 0, 486, 485, 1, 0, 0, 0, 486,
		487, 1, 0, 0, 0, 487, 488, 1, 0, 0, 0, 488, 490, 3, 58, 29, 0, 489, 486,
		1, 0, 0, 0, 490, 493, 1, 0, 0, 0, 491, 489, 1, 0, 0, 0, 491, 492, 1, 0,
		0, 0, 492, 494, 1, 0, 0, 0, 493, 491, 1, 0, 0, 0, 494, 495, 5, 4, 0, 0,
		495, 57, 1, 0, 0, 0, 496, 497, 5, 31, 0, 0, 497, 498, 3, 22, 11, 0, 498,
		499, 5, 19, 0, 0, 499, 503, 3, 60, 30, 0, 500, 502, 5, 35, 0, 0, 501, 500,
		1, 0, 0, 0, 502, 505, 1, 0, 0, 0, 503, 501, 1, 0, 0, 0, 503, 504, 1, 0,
		0, 0, 504, 59, 1, 0, 0, 0, 505, 503, 1, 0, 0, 0, 506, 515, 5, 20, 0, 0,
		507, 515, 5, 21, 0, 0, 508, 515, 5, 32, 0, 0, 509, 515, 5, 33, 0, 0, 510,
		515, 5, 34, 0, 0, 511, 515, 3, 62, 31, 0, 512, 515, 3, 66, 33, 0, 513,
		515, 5, 22, 0, 0, 514, 506, 1, 0, 0, 0, 514, 507, 1, 0, 0, 0, 514, 508,
		1, 0, 0, 0, 514, 509, 1, 0, 0, 0, 514, 510, 1, 0, 0, 0, 514, 511, 1, 0,
		0, 0, 514, 512, 1, 0, 0, 0, 514, 513, 1, 0, 0, 0, 515, 61, 1, 0, 0, 0,
		516, 520, 5, 12, 0, 0, 517, 519, 5, 35, 0, 0, 518, 517, 1, 0, 0, 0, 519,
		522, 1, 0, 0, 0, 520, 518, 1, 0, 0, 0, 520, 521, 1, 0, 0, 0, 521, 524,
		1, 0, 0, 0, 522, 520, 1, 0, 0, 0, 523, 525, 3, 64, 32, 0, 524, 523, 1,
		0, 0, 0, 524, 525, 1, 0, 0, 0, 525, 526, 1, 0, 0, 0, 526, 527, 5, 13, 0,
		0, 527, 63, 1, 0, 0, 0, 528, 550, 3, 60, 30, 0, 529, 546, 3, 60, 30, 0,
		530, 534, 5, 11, 0, 0, 531, 533, 5, 35, 0, 0, 532, 531, 1, 0, 0, 0, 533,
		536, 1, 0, 0, 0, 534, 532, 1, 0, 0, 0, 534, 535, 1, 0, 0, 0, 535, 537,
		1, 0, 0, 0, 536, 534, 1, 0, 0, 0, 537, 541, 3, 60, 30, 0, 538, 540, 5,
		35, 0, 0, 539, 538, 1, 0, 0, 0, 540, 543, 1, 0, 0, 0, 541, 539, 1, 0, 0,
		0, 541, 542, 1, 0, 0, 0, 542, 545, 1, 0, 0, 0, 543, 541, 1, 0, 0, 0, 544,
		530, 1, 0, 0, 0, 545, 548, 1, 0, 0, 0, 546, 544, 1, 0, 0, 0, 546, 547,
		1, 0, 0, 0, 547, 550, 1, 0, 0, 0, 548, 546, 1, 0, 0, 0, 549, 528, 1, 0,
		0, 0, 549, 529, 1, 0, 0, 0, 550, 65, 1, 0, 0, 0, 551, 555, 5, 3, 0, 0,
		552, 554, 5, 35, 0, 0, 553, 552, 1, 0, 0, 0, 554, 557, 1, 0, 0, 0, 555,
		553, 1, 0, 0, 0, 555, 556, 1, 0, 0, 0, 556, 559, 1, 0, 0, 0, 557, 555,
		1, 0, 0, 0, 558, 560, 3, 68, 34, 0, 559, 558, 1, 0, 0, 0, 559, 560, 1,
		0, 0, 0, 560, 561, 1, 0, 0, 0, 561, 562, 5, 4, 0, 0, 562, 67, 1, 0, 0,
		0, 563, 573, 3, 70, 35, 0, 564, 566, 5, 35, 0, 0, 565, 564, 1, 0, 0, 0,
		566, 569, 1, 0, 0, 0, 567, 565, 1, 0, 0, 0, 567, 568, 1, 0, 0, 0, 568,
		570, 1, 0, 0, 0, 569, 567, 1, 0, 0, 0, 570, 572, 3, 70, 35, 0, 571, 567,
		1, 0, 0, 0, 572, 575, 1, 0, 0, 0, 573, 571, 1, 0, 0, 0, 573, 574, 1, 0,
		0, 0, 574, 69, 1, 0, 0, 0, 575, 573, 1, 0, 0, 0, 576, 577, 5, 31, 0, 0,
		577, 578, 5, 23, 0, 0, 578, 582, 3, 60, 30, 0, 579, 581, 5, 35, 0, 0, 580,
		579, 1, 0, 0, 0, 581, 584, 1, 0, 0, 0, 582, 580, 1, 0, 0, 0, 582, 583,
		1, 0, 0, 0, 583, 71, 1, 0, 0, 0, 584, 582, 1, 0, 0, 0, 585, 589, 5, 24,
		0, 0, 586, 588, 5, 35, 0, 0, 587, 586, 1, 0, 0, 0, 588, 591, 1, 0, 0, 0,
		589, 587, 1, 0, 0, 0, 589, 590, 1, 0, 0, 0, 590, 592, 1, 0, 0, 0, 591,
		589, 1, 0, 0, 0, 592, 596, 5, 3, 0, 0, 593, 595, 5, 35, 0, 0, 594, 593,
		1, 0, 0, 0, 595, 598, 1, 0, 0, 0, 596, 594, 1, 0, 0, 0, 596, 597, 1, 0,
		0, 0, 597, 605, 1, 0, 0, 0, 598, 596, 1, 0, 0, 0, 599, 601, 5, 8, 0, 0,
		600, 599, 1, 0, 0, 0, 600, 601, 1, 0, 0, 0, 601, 602, 1, 0, 0, 0, 602,
		604, 3, 74, 37, 0, 603, 600, 1, 0, 0, 0, 604, 607, 1, 0, 0, 0, 605, 603,
		1, 0, 0, 0, 605, 606, 1, 0, 0, 0, 606, 608, 1, 0, 0, 0, 607, 605, 1, 0,
		0, 0, 608, 609, 5, 4, 0, 0, 609, 73, 1, 0, 0, 0, 610, 611, 3, 46, 23, 0,
		611, 615, 3, 76, 38, 0, 612, 614, 5, 35, 0, 0, 613, 612, 1, 0, 0, 0, 614,
		617, 1, 0, 0, 0, 615, 613, 1, 0, 0, 0, 615, 616, 1, 0, 0, 0, 616, 75, 1,
		0, 0, 0, 617, 615, 1, 0, 0, 0, 618, 622, 5, 3, 0, 0, 619, 621, 5, 35, 0,
		0, 620, 619, 1, 0, 0, 0, 621, 624, 1, 0, 0, 0, 622, 620, 1, 0, 0, 0, 622,
		623, 1, 0, 0, 0, 623, 627, 1, 0, 0, 0, 624, 622, 1, 0, 0, 0, 625, 628,
		3, 78, 39, 0, 626, 628, 3, 96, 48, 0, 627, 625, 1, 0, 0, 0, 627, 626, 1,
		0, 0, 0, 627, 628, 1, 0, 0, 0, 628, 629, 1, 0, 0, 0, 629, 630, 5, 4, 0,
		0, 630, 77, 1, 0, 0, 0, 631, 635, 5, 25, 0, 0, 632, 634, 5, 35, 0, 0, 633,
		632, 1, 0, 0, 0, 634, 637, 1, 0, 0, 0, 635, 633, 1, 0, 0, 0, 635, 636,
		1, 0, 0, 0, 636, 638, 1, 0, 0, 0, 637, 635, 1, 0, 0, 0, 638, 642, 5, 3,
		0, 0, 639, 641, 5, 35, 0, 0, 640, 639, 1, 0, 0, 0, 641, 644, 1, 0, 0, 0,
		642, 640, 1, 0, 0, 0, 642, 643, 1, 0, 0, 0, 643, 645, 1, 0, 0, 0, 644,
		642, 1, 0, 0, 0, 645, 646, 3, 80, 40, 0, 646, 647, 5, 4, 0, 0, 647, 79,
		1, 0, 0, 0, 648, 651, 3, 82, 41, 0, 649, 651, 3, 84, 42, 0, 650, 648, 1,
		0, 0, 0, 650, 649, 1, 0, 0, 0, 651, 81, 1, 0, 0, 0, 652, 653, 5, 31, 0,
		0, 653, 654, 3, 24, 12, 0, 654, 83, 1, 0, 0, 0, 655, 656, 5, 31, 0, 0,
		656, 657, 5, 19, 0, 0, 657, 658, 3, 86, 43, 0, 658, 85, 1, 0, 0, 0, 659,
		660, 3, 88, 44, 0, 660, 661, 3, 90, 45, 0, 661, 662, 3, 26, 13, 0, 662,
		87, 1, 0, 0, 0, 663, 668, 5, 31, 0, 0, 664, 665, 5, 26, 0, 0, 665, 667,
		5, 31, 0, 0, 666, 664, 1, 0, 0, 0, 667, 670, 1, 0, 0, 0, 668, 666, 1, 0,
		0, 0, 668, 669, 1, 0, 0, 0, 669, 89, 1, 0, 0, 0, 670, 668, 1, 0, 0, 0,
		671, 672, 5, 16, 0, 0, 672, 673, 3, 92, 46, 0, 673, 674, 5, 17, 0, 0, 674,
		91, 1, 0, 0, 0, 675, 676, 3, 94, 47, 0, 676, 678, 5, 11, 0, 0, 677, 679,
		5, 35, 0, 0, 678, 677, 1, 0, 0, 0, 678, 679, 1, 0, 0, 0, 679, 680, 1, 0,
		0, 0, 680, 681, 3, 94, 47, 0, 681, 93, 1, 0, 0, 0, 682, 683, 3, 86, 43,
		0, 683, 95, 1, 0, 0, 0, 684, 688, 5, 27, 0, 0, 685, 687, 5, 35, 0, 0, 686,
		685, 1, 0, 0, 0, 687, 690, 1, 0, 0, 0, 688, 686, 1, 0, 0, 0, 688, 689,
		1, 0, 0, 0, 689, 691, 1, 0, 0, 0, 690, 688, 1, 0, 0, 0, 691, 695, 5, 3,
		0, 0, 692, 694, 5, 35, 0, 0, 693, 692, 1, 0, 0, 0, 694, 697, 1, 0, 0, 0,
		695, 693, 1, 0, 0, 0, 695, 696, 1, 0, 0, 0, 696, 698, 1, 0, 0, 0, 697,
		695, 1, 0, 0, 0, 698, 699, 3, 98, 49, 0, 699, 700, 5, 4, 0, 0, 700, 97,
		1, 0, 0, 0, 701, 706, 3, 100, 50, 0, 702, 703, 5, 35, 0, 0, 703, 705, 3,
		100, 50, 0, 704, 702, 1, 0, 0, 0, 705, 708, 1, 0, 0, 0, 706, 704, 1, 0,
		0, 0, 706, 707, 1, 0, 0, 0, 707, 99, 1, 0, 0, 0, 708, 706, 1, 0, 0, 0,
		709, 710, 3, 102, 51, 0, 710, 711, 5, 28, 0, 0, 711, 712, 3, 106, 53, 0,
		712, 101, 1, 0, 0, 0, 713, 715, 5, 31, 0, 0, 714, 713, 1, 0, 0, 0, 714,
		715, 1, 0, 0, 0, 715, 716, 1, 0, 0, 0, 716, 724, 3, 104, 52, 0, 717, 721,
		5, 31, 0, 0, 718, 719, 5, 12, 0, 0, 719, 720, 5, 32, 0, 0, 720, 722, 5,
		13, 0, 0, 721, 718, 1, 0, 0, 0, 721, 722, 1, 0, 0, 0, 722, 724, 1, 0, 0,
		0, 723, 714, 1, 0, 0, 0, 723, 717, 1, 0, 0, 0, 724, 103, 1, 0, 0, 0, 725,
		726, 7, 1, 0, 0, 726, 105, 1, 0, 0, 0, 727, 730, 3, 102, 51, 0, 728, 730,
		3, 108, 54, 0, 729, 727, 1, 0, 0, 0, 729, 728, 1, 0, 0, 0, 730, 107, 1,
		0, 0, 0, 731, 732, 5, 3, 0, 0, 732, 737, 3, 102, 51, 0, 733, 734, 5, 35,
		0, 0, 734, 736, 3, 102, 51, 0, 735, 733, 1, 0, 0, 0, 736, 739, 1, 0, 0,
		0, 737, 735, 1, 0, 0, 0, 737, 738, 1, 0, 0, 0, 738, 740, 1, 0, 0, 0, 739,
		737, 1, 0, 0, 0, 740, 741, 5, 4, 0, 0, 741, 109, 1, 0, 0, 0, 94, 113, 115,
		124, 132, 138, 145, 151, 157, 163, 167, 174, 181, 188, 192, 198, 203, 210,
		218, 222, 231, 238, 243, 248, 253, 257, 263, 271, 277, 283, 291, 297, 305,
		311, 317, 326, 333, 343, 347, 355, 360, 368, 375, 382, 388, 392, 398, 405,
		409, 414, 421, 428, 439, 443, 450, 453, 460, 468, 475, 482, 486, 491, 503,
		514, 520, 524, 534, 541, 546, 549, 555, 559, 567, 573, 582, 589, 596, 600,
		605, 615, 622, 627, 635, 642, 650, 668, 678, 688, 695, 706, 714, 721, 723,
		729, 737,
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
	nevaParserRULE_typeParamList    = 9
	nevaParserRULE_typeParam        = 10
	nevaParserRULE_typeExpr         = 11
	nevaParserRULE_typeInstExpr     = 12
	nevaParserRULE_typeArgs         = 13
	nevaParserRULE_typeLitExpr      = 14
	nevaParserRULE_enumTypeExpr     = 15
	nevaParserRULE_arrTypeExpr      = 16
	nevaParserRULE_recTypeExpr      = 17
	nevaParserRULE_recFields        = 18
	nevaParserRULE_recField         = 19
	nevaParserRULE_unionTypeExpr    = 20
	nevaParserRULE_nonUnionTypeExpr = 21
	nevaParserRULE_ioStmt           = 22
	nevaParserRULE_interfaceDef     = 23
	nevaParserRULE_inPortsDef       = 24
	nevaParserRULE_outPortsDef      = 25
	nevaParserRULE_portsDef         = 26
	nevaParserRULE_portAndType      = 27
	nevaParserRULE_constStmt        = 28
	nevaParserRULE_constDef         = 29
	nevaParserRULE_constVal         = 30
	nevaParserRULE_arrLit           = 31
	nevaParserRULE_vecItems         = 32
	nevaParserRULE_recLit           = 33
	nevaParserRULE_recValueFields   = 34
	nevaParserRULE_recValueField    = 35
	nevaParserRULE_compStmt         = 36
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
	p.SetState(138)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(135)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(140)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(141)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(142)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(147)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(151)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__4 || _la == nevaParserIDENTIFIER {
		{
			p.SetState(148)
			p.ImportDef()
		}

		p.SetState(153)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(154)
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
	p.SetState(157)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(156)
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
		p.SetState(159)
		p.ImportPath()
	}
	p.SetState(163)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(160)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(165)
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
	p.SetState(167)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__4 {
		{
			p.SetState(166)
			p.Match(nevaParserT__4)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(169)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(174)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__5 {
		{
			p.SetState(170)
			p.Match(nevaParserT__5)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(171)
			p.Match(nevaParserIDENTIFIER)
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
		p.SetState(177)
		p.Match(nevaParserT__6)
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

	for _la == nevaParserNEWLINE {
		{
			p.SetState(178)
			p.Match(nevaParserNEWLINE)
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
	{
		p.SetState(184)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(188)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(185)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(190)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(203)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 || _la == nevaParserIDENTIFIER {
		p.SetState(192)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserT__7 {
			{
				p.SetState(191)
				p.Match(nevaParserT__7)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(194)
			p.TypeDef()
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

		p.SetState(205)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(206)
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
	{
		p.SetState(208)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(210)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserT__8 {
		{
			p.SetState(209)
			p.TypeParams()
		}

	}
	{
		p.SetState(212)
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
	TypeParamList() ITypeParamListContext

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

func (s *TypeParamsContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *TypeParamsContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *TypeParamsContext) TypeParamList() ITypeParamListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeParamListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeParamListContext)
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
		p.SetState(214)
		p.Match(nevaParserT__8)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(218)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(215)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(220)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(222)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserIDENTIFIER {
		{
			p.SetState(221)
			p.TypeParamList()
		}

	}
	{
		p.SetState(224)
		p.Match(nevaParserT__9)
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

// ITypeParamListContext is an interface to support dynamic dispatch.
type ITypeParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTypeParam() []ITypeParamContext
	TypeParam(i int) ITypeParamContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

	// IsTypeParamListContext differentiates from other interfaces.
	IsTypeParamListContext()
}

type TypeParamListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeParamListContext() *TypeParamListContext {
	var p = new(TypeParamListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeParamList
	return p
}

func InitEmptyTypeParamListContext(p *TypeParamListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = nevaParserRULE_typeParamList
}

func (*TypeParamListContext) IsTypeParamListContext() {}

func NewTypeParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeParamListContext {
	var p = new(TypeParamListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = nevaParserRULE_typeParamList

	return p
}

func (s *TypeParamListContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeParamListContext) AllTypeParam() []ITypeParamContext {
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

func (s *TypeParamListContext) TypeParam(i int) ITypeParamContext {
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

func (s *TypeParamListContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *TypeParamListContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *TypeParamListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeParamListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.EnterTypeParamList(s)
	}
}

func (s *TypeParamListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(nevaListener); ok {
		listenerT.ExitTypeParamList(s)
	}
}

func (p *nevaParser) TypeParamList() (localctx ITypeParamListContext) {
	localctx = NewTypeParamListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, nevaParserRULE_typeParamList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(226)
		p.TypeParam()
	}
	p.SetState(243)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__10 {
		{
			p.SetState(227)
			p.Match(nevaParserT__10)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(231)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(228)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(233)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(234)
			p.TypeParam()
		}
		p.SetState(238)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(235)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(240)
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
	p.EnterRule(localctx, 20, nevaParserRULE_typeParam)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(246)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(248)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2147487752) != 0 {
		{
			p.SetState(247)
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
	p.EnterRule(localctx, 22, nevaParserRULE_typeExpr)
	p.SetState(253)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 23, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(250)
			p.TypeInstExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(251)
			p.TypeLitExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(252)
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
	p.EnterRule(localctx, 24, nevaParserRULE_typeInstExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(255)
		p.Match(nevaParserIDENTIFIER)
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

	if _la == nevaParserT__8 {
		{
			p.SetState(256)
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
	p.EnterRule(localctx, 26, nevaParserRULE_typeArgs)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(259)
		p.Match(nevaParserT__8)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(263)
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

		p.SetState(265)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(266)
		p.TypeExpr()
	}
	p.SetState(277)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__10 {
		{
			p.SetState(267)
			p.Match(nevaParserT__10)
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

		for _la == nevaParserNEWLINE {
			{
				p.SetState(268)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(273)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(274)
			p.TypeExpr()
		}

		p.SetState(279)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
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
		p.Match(nevaParserT__9)
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
	p.EnterRule(localctx, 28, nevaParserRULE_typeLitExpr)
	p.SetState(291)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 29, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(288)
			p.EnumTypeExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(289)
			p.ArrTypeExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(290)
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
	p.EnterRule(localctx, 30, nevaParserRULE_enumTypeExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(293)
		p.Match(nevaParserT__2)
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

	for _la == nevaParserNEWLINE {
		{
			p.SetState(294)
			p.Match(nevaParserNEWLINE)
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
	{
		p.SetState(300)
		p.Match(nevaParserIDENTIFIER)
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

	for _la == nevaParserT__10 {
		{
			p.SetState(301)
			p.Match(nevaParserT__10)
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

		for _la == nevaParserNEWLINE {
			{
				p.SetState(302)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(307)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(308)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(313)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
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
	p.EnterRule(localctx, 32, nevaParserRULE_arrTypeExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(322)
		p.Match(nevaParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(326)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(323)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(328)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(329)
		p.Match(nevaParserINT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(333)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(330)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(335)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(336)
		p.Match(nevaParserT__12)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(337)
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
	p.EnterRule(localctx, 34, nevaParserRULE_recTypeExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(339)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(343)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(340)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(345)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(347)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserIDENTIFIER {
		{
			p.SetState(346)
			p.RecFields()
		}

	}
	{
		p.SetState(349)
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
	p.EnterRule(localctx, 36, nevaParserRULE_recFields)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(351)
		p.RecField()
	}
	p.SetState(360)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		p.SetState(353)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == nevaParserNEWLINE {
			{
				p.SetState(352)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(355)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(357)
			p.RecField()
		}

		p.SetState(362)
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
	p.EnterRule(localctx, 38, nevaParserRULE_recField)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(363)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(364)
		p.TypeExpr()
	}
	p.SetState(368)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(365)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		p.SetState(370)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
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
	p.EnterRule(localctx, 40, nevaParserRULE_unionTypeExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(371)
		p.NonUnionTypeExpr()
	}
	p.SetState(386)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			p.SetState(375)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(372)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(377)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(378)
				p.Match(nevaParserT__13)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(382)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(379)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(384)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(385)
				p.NonUnionTypeExpr()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(388)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 42, nevaParserRULE_nonUnionTypeExpr)
	p.SetState(392)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(390)
			p.TypeInstExpr()
		}

	case nevaParserT__2, nevaParserT__11:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(391)
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
	p.EnterRule(localctx, 44, nevaParserRULE_ioStmt)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(394)
		p.Match(nevaParserT__14)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(398)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(395)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(400)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(401)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(405)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(402)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(407)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(414)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 || _la == nevaParserIDENTIFIER {
		p.SetState(409)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserT__7 {
			{
				p.SetState(408)
				p.Match(nevaParserT__7)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(411)
			p.InterfaceDef()
		}

		p.SetState(416)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(417)
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
	p.EnterRule(localctx, 46, nevaParserRULE_interfaceDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(419)
		p.Match(nevaParserIDENTIFIER)
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

	if _la == nevaParserT__8 {
		{
			p.SetState(420)
			p.TypeParams()
		}

	}
	{
		p.SetState(423)
		p.InPortsDef()
	}
	{
		p.SetState(424)
		p.OutPortsDef()
	}
	p.SetState(428)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(425)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(430)
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
	p.EnterRule(localctx, 48, nevaParserRULE_inPortsDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(431)
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
	p.EnterRule(localctx, 50, nevaParserRULE_outPortsDef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(433)
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
	p.EnterRule(localctx, 52, nevaParserRULE_portsDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(435)
		p.Match(nevaParserT__15)
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

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 54, p.GetParserRuleContext()) {
	case 1:
		p.SetState(439)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserNEWLINE {
			{
				p.SetState(436)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(441)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		p.SetState(443)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserIDENTIFIER || _la == nevaParserNEWLINE {
			{
				p.SetState(442)
				p.PortAndType()
			}

		}

	case 3:
		{
			p.SetState(445)
			p.PortAndType()
		}
		p.SetState(450)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserT__10 {
			{
				p.SetState(446)
				p.Match(nevaParserT__10)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(447)
				p.PortAndType()
			}

			p.SetState(452)
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
		p.SetState(455)
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
	p.EnterRule(localctx, 54, nevaParserRULE_portAndType)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(460)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(457)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(464)
		p.TypeExpr()
	}
	p.SetState(468)
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

		p.SetState(470)
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
	p.EnterRule(localctx, 56, nevaParserRULE_constStmt)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(471)
		p.Match(nevaParserT__17)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
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
	{
		p.SetState(478)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(482)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(479)
			p.Match(nevaParserNEWLINE)
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
	}
	p.SetState(491)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 || _la == nevaParserIDENTIFIER {
		p.SetState(486)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserT__7 {
			{
				p.SetState(485)
				p.Match(nevaParserT__7)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(488)
			p.ConstDef()
		}

		p.SetState(493)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(494)
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
	p.EnterRule(localctx, 58, nevaParserRULE_constDef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(496)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(497)
		p.TypeExpr()
	}
	{
		p.SetState(498)
		p.Match(nevaParserT__18)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(499)
		p.ConstVal()
	}
	p.SetState(503)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(500)
			p.Match(nevaParserNEWLINE)
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
	p.EnterRule(localctx, 60, nevaParserRULE_constVal)
	p.SetState(514)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__19:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(506)
			p.Match(nevaParserT__19)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserT__20:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(507)
			p.Match(nevaParserT__20)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserINT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(508)
			p.Match(nevaParserINT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserFLOAT:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(509)
			p.Match(nevaParserFLOAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserSTRING:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(510)
			p.Match(nevaParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case nevaParserT__11:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(511)
			p.ArrLit()
		}

	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(512)
			p.RecLit()
		}

	case nevaParserT__21:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(513)
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
	p.EnterRule(localctx, 62, nevaParserRULE_arrLit)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(516)
		p.Match(nevaParserT__11)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(520)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(517)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(522)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(524)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&30072115208) != 0 {
		{
			p.SetState(523)
			p.VecItems()
		}

	}
	{
		p.SetState(526)
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
	p.EnterRule(localctx, 64, nevaParserRULE_vecItems)
	var _la int

	p.SetState(549)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 68, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(528)
			p.ConstVal()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(529)
			p.ConstVal()
		}
		p.SetState(546)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == nevaParserT__10 {
			{
				p.SetState(530)
				p.Match(nevaParserT__10)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(534)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(531)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(536)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(537)
				p.ConstVal()
			}
			p.SetState(541)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == nevaParserNEWLINE {
				{
					p.SetState(538)
					p.Match(nevaParserNEWLINE)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				p.SetState(543)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

			p.SetState(548)
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
	p.EnterRule(localctx, 66, nevaParserRULE_recLit)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(551)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(555)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(552)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(557)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(559)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserIDENTIFIER {
		{
			p.SetState(558)
			p.RecValueFields()
		}

	}
	{
		p.SetState(561)
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
		p.SetState(563)
		p.RecValueField()
	}
	p.SetState(573)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserIDENTIFIER || _la == nevaParserNEWLINE {
		p.SetState(567)
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

			p.SetState(569)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(570)
			p.RecValueField()
		}

		p.SetState(575)
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
	p.EnterRule(localctx, 70, nevaParserRULE_recValueField)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(576)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(577)
		p.Match(nevaParserT__22)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(578)
		p.ConstVal()
	}
	p.SetState(582)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 73, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(579)
				p.Match(nevaParserNEWLINE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		p.SetState(584)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 73, p.GetParserRuleContext())
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
	AllCompDef() []ICompDefContext
	CompDef(i int) ICompDefContext

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

func (s *CompStmtContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *CompStmtContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

func (s *CompStmtContext) AllCompDef() []ICompDefContext {
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

func (s *CompStmtContext) CompDef(i int) ICompDefContext {
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
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(585)
		p.Match(nevaParserT__23)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(589)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(586)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(591)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(592)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(596)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(593)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(598)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(605)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__7 || _la == nevaParserIDENTIFIER {
		p.SetState(600)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserT__7 {
			{
				p.SetState(599)
				p.Match(nevaParserT__7)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(602)
			p.CompDef()
		}

		p.SetState(607)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(608)
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

// ICompDefContext is an interface to support dynamic dispatch.
type ICompDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	InterfaceDef() IInterfaceDefContext
	CompBody() ICompBodyContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *CompDefContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *CompDefContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(610)
		p.InterfaceDef()
	}
	{
		p.SetState(611)
		p.CompBody()
	}
	p.SetState(615)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(612)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(617)
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

// ICompBodyContext is an interface to support dynamic dispatch.
type ICompBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode
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

func (s *CompBodyContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *CompBodyContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
}

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
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(618)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(622)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(619)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(624)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(627)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	switch p.GetTokenStream().LA(1) {
	case nevaParserT__24:
		{
			p.SetState(625)
			p.CompNodesDef()
		}

	case nevaParserT__26:
		{
			p.SetState(626)
			p.CompNetDef()
		}

	case nevaParserT__3:

	default:
	}
	{
		p.SetState(629)
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

// ICompNodesDefContext is an interface to support dynamic dispatch.
type ICompNodesDefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CompNodeDefList() ICompNodeDefListContext
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *CompNodesDefContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *CompNodesDefContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(631)
		p.Match(nevaParserT__24)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(635)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(632)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(637)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(638)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(642)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(639)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(644)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(645)
		p.CompNodeDefList()
	}
	{
		p.SetState(646)
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
	p.SetState(650)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 83, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(648)
			p.AbsNodeDef()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(649)
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
		p.SetState(652)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(653)
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
		p.SetState(655)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(656)
		p.Match(nevaParserT__18)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(657)
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
		p.SetState(659)
		p.NodeRef()
	}
	{
		p.SetState(660)
		p.NodeArgs()
	}
	{
		p.SetState(661)
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
		p.SetState(663)
		p.Match(nevaParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(668)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserT__25 {
		{
			p.SetState(664)
			p.Match(nevaParserT__25)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(665)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(670)
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
		p.SetState(671)
		p.Match(nevaParserT__15)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(672)
		p.NodeArgList()
	}
	{
		p.SetState(673)
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
		p.SetState(675)
		p.NodeArg()
	}

	{
		p.SetState(676)
		p.Match(nevaParserT__10)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(678)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == nevaParserNEWLINE {
		{
			p.SetState(677)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(680)
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
		p.SetState(682)
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
	AllNEWLINE() []antlr.TerminalNode
	NEWLINE(i int) antlr.TerminalNode

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

func (s *CompNetDefContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(nevaParserNEWLINE)
}

func (s *CompNetDefContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(nevaParserNEWLINE, i)
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
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(684)
		p.Match(nevaParserT__26)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(688)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(685)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(690)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(691)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(695)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(692)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(697)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(698)
		p.ConnDefList()
	}
	{
		p.SetState(699)
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
		p.SetState(701)
		p.ConnDef()
	}
	p.SetState(706)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(702)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(703)
			p.ConnDef()
		}

		p.SetState(708)
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
		p.SetState(709)
		p.PortAddr()
	}
	{
		p.SetState(710)
		p.Match(nevaParserT__27)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(711)
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

	p.SetState(723)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 91, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(714)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserIDENTIFIER {
			{
				p.SetState(713)
				p.Match(nevaParserIDENTIFIER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(716)
			p.PortDirection()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(717)
			p.Match(nevaParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(721)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == nevaParserT__11 {
			{
				p.SetState(718)
				p.Match(nevaParserT__11)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(719)
				p.Match(nevaParserINT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(720)
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
		p.SetState(725)
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
	p.SetState(729)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case nevaParserT__28, nevaParserT__29, nevaParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(727)
			p.PortAddr()
		}

	case nevaParserT__2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(728)
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
		p.SetState(731)
		p.Match(nevaParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(732)
		p.PortAddr()
	}
	p.SetState(737)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == nevaParserNEWLINE {
		{
			p.SetState(733)
			p.Match(nevaParserNEWLINE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(734)
			p.PortAddr()
		}

		p.SetState(739)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(740)
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

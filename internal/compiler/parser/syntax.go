package parser

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/nevalang/neva/internal/compiler"
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
)

// ParsedSource is the syntax-preserving result of parsing one Neva source file.
//
// Unlike the compiler AST, it retains the original source and every parser-visible
// token, including comments and newlines. Source tooling such as the formatter
// must use this representation rather than the semantic AST, whose maps do not
// retain source order or complete trivia.
type ParsedSource struct {
	Tree   generated.IProgContext
	Source []byte
	Tokens []antlr.Token
}

// ParseSource parses source without building the compiler's semantic AST.
// It returns the complete parser-visible token stream so source tooling can
// preserve comments and newlines while rendering the parsed syntax.
func ParseSource(source []byte) (ParsedSource, *compiler.Error) {
	input := antlr.NewInputStream(string(source))
	lexer := generated.NewnevaLexer(input)
	lexerErrors := &CustomErrorListener{}
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexerErrors)

	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parserErrors := &CustomErrorListener{}
	prsr := generated.NewnevaParser(tokenStream)
	prsr.RemoveErrorListeners()
	prsr.AddErrorListener(parserErrors)
	prsr.BuildParseTrees = true

	tree := prsr.Prog()
	tokenStream.Fill()

	if len(lexerErrors.Errors) > 0 {
		return ParsedSource{}, lexerErrors.Errors[0]
	}
	if len(parserErrors.Errors) > 0 {
		return ParsedSource{}, parserErrors.Errors[0]
	}

	return ParsedSource{
		Source: source,
		Tokens: tokenStream.GetAllTokens(),
		Tree:   tree,
	}, nil
}

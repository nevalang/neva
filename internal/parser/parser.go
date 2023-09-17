// Package parser implements source code parsing.
package parser

import (
	"context"

	"github.com/antlr4-go/antlr/v4"

	generated "github.com/nevalang/neva/internal/parser/generated"
	"github.com/nevalang/neva/internal/src"
)

type treeShapeListener struct {
	*generated.BasenevaListener
	file src.File
}

type Parser struct{}

// Currently only returns one "main" package
func (p Parser) Parse(ctx context.Context, bb []byte) (map[string]src.File, error) {
	input := antlr.NewInputStream(string(bb))
	lexer := generated.NewnevaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	prsr := generated.NewnevaParser(stream)
	prsr.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	prsr.BuildParseTrees = true
	tree := prsr.Prog()
	listener := &treeShapeListener{}

	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	return map[string]src.File{"main": listener.file}, nil
}

func New() Parser {
	return Parser{}
}

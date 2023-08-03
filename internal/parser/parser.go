package parser

import (
	"context"

	"github.com/antlr4-go/antlr/v4"

	generated "github.com/nevalang/neva/internal/parser/generated"
	"github.com/nevalang/neva/internal/shared"
)

type treeShapeListener struct {
	*generated.BasenevaListener
	file shared.File
}

type Parser struct{}

func (p Parser) Parse(ctx context.Context, bb []byte) (map[string]shared.File, error) {
	input := antlr.NewInputStream(string(bb))
	lexer := generated.NewnevaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	prsr := generated.NewnevaParser(stream)
	prsr.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	prsr.BuildParseTrees = true
	tree := prsr.Prog()
	listener := &treeShapeListener{}

	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	return map[string]shared.File{"main": listener.file}, nil
}

func New() Parser {
	return Parser{}
}

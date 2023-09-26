// Package parser implements source code parsing.
package parser

import (
	"context"

	"github.com/antlr4-go/antlr/v4"
	"golang.org/x/sync/errgroup"

	generated "github.com/nevalang/neva/internal/parser/generated"
	"github.com/nevalang/neva/internal/src"
)

type treeShapeListener struct {
	*generated.BasenevaListener
	file src.File
}

type Parser struct {
	isDebug bool
}

func (p Parser) ParseFiles(ctx context.Context, files map[string][]byte) (map[string]src.File, error) {
	result := make(map[string]src.File, len(files))
	g, gctx := errgroup.WithContext(ctx)
	for name, bb := range files {
		name := name
		bb := bb
		g.Go(func() error {
			v, err := p.ParseFile(gctx, bb)
			if err != nil {
				return err
			}
			result[name] = v
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}

func (p Parser) ParseFile(ctx context.Context, bb []byte) (src.File, error) {
	input := antlr.NewInputStream(string(bb))
	lexer := generated.NewnevaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	parse := generated.NewnevaParser(stream)
	if p.isDebug {
		parse.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	}
	parse.BuildParseTrees = true

	tree := parse.Prog()
	listener := &treeShapeListener{}

	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	return listener.file, nil
}

func New(isDebug bool) Parser {
	return Parser{
		isDebug: isDebug,
	}
}

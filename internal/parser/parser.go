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
	file src.Package
}

type Parser struct{}

func (p Parser) ParseFiles(ctx context.Context, files map[string][]byte) (map[string]src.Package, error) {
	result := make(map[string]src.Package, len(files))
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

func (p Parser) ParseFile(ctx context.Context, bb []byte) (src.Package, error) {
	input := antlr.NewInputStream(string(bb))
	lexer := generated.NewnevaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	prsr := generated.NewnevaParser(stream)
	prsr.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	prsr.BuildParseTrees = true
	tree := prsr.Prog()
	listener := &treeShapeListener{}

	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	return listener.file, nil
}

func New() Parser {
	return Parser{}
}

package main

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	parser "github.com/nevalang/neva/internal/compiler/frontend/generated"
)

type TreeShapeListener struct {
	*parser.BasenevaListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (this *TreeShapeListener) EnterEveryRule(actx antlr.ParserRuleContext) {
	fmt.Println(actx.GetText())
}

func main() {
	input, err := antlr.NewFileStream(os.Args[1])
	if err != nil {
		panic(err)
	}

	lexer := parser.NewnevaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewnevaParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.Prog()

	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
}

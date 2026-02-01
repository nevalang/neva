package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/ast/core"
)

type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []*compiler.Error
}

func (c *CustomErrorListener) SyntaxError(
	_ antlr.Recognizer,
	offendingSymbol any,
	line, column int,
	msg string,
	e antlr.RecognitionException,
) {
	c.Errors = append(c.Errors, &compiler.Error{
		Message: msg,
		Meta: &core.Meta{
			Start: core.Position{
				Line:   line,
				Column: column,
			},
		},
	})
}

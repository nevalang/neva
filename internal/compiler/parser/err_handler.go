package parser

import (
	"errors"

	"github.com/antlr4-go/antlr/v4"
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []error
}

func (c *CustomErrorListener) SyntaxError(
	_ antlr.Recognizer,
	offendingSymbol any,
	line, column int,
	msg string,
	e antlr.RecognitionException,
) {
	c.Errors = append(c.Errors, compiler.Error{
		Err: errors.New(msg),
		Meta: &src.Meta{
			Start: src.Position{
				Line:   line,
				Column: column,
			},
		},
	})
}

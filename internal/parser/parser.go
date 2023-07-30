package parser

import (
	"context"

	"github.com/nevalang/neva/internal/shared"
)

type Parser struct{}

func (p Parser) Parse(context.Context, []byte) (shared.HighLvlProgram, error) {
	return shared.HighLvlProgram{}, nil
}

func MustNew() Parser {
	return Parser{}
}

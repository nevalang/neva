package parser

import (
	"context"

	"github.com/nevalang/neva/internal/shared"
)

type Parser struct{}

func (p Parser) Parse(context.Context, []byte) (map[string]shared.HLPackage, error) {
	return map[string]shared.HLPackage{}, nil
}

func MustNew() Parser {
	return Parser{}
}

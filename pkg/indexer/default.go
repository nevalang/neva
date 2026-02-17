package indexer

import (
	"github.com/tliron/commonlog"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/compiler/typesystem"
)

// NewDefault creates an Indexer with default compiler frontend dependencies.
func NewDefault(logger commonlog.Logger) (Indexer, error) {
	p := parser.New()

	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)

	b, err := builder.New(p)
	if err != nil {
		return Indexer{}, err
	}

	return New(b, p, analyzer.MustNew(resolver), logger), nil
}

// MustNewDefault creates an Indexer with default dependencies and panics on error.
func MustNewDefault(logger commonlog.Logger) Indexer {
	idx, err := NewDefault(logger)
	if err != nil {
		panic(err)
	}
	return idx
}

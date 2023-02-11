package analyze

import (
	"context"

	"github.com/emil14/neva/internal/compiler/src"
)

type Analyzer struct{}

func (a Analyzer) Analyze(ctx context.Context, prog src.Prog) error {
	return nil
}

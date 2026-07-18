package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
)

func Test(t *testing.T) {
	// A tag-only destination may safely ignore a payload from the same tag.
	e2e.Run(t, []string{"run", "main"})
}

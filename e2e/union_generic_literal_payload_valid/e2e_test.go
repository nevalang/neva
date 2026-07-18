package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
)

func Test(t *testing.T) {
	// Valid generic union payloads should keep compiling after payload checks.
	e2e.Run(t, []string{"run", "main"})
}

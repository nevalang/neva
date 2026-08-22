package parser

import (
	"testing"

	"github.com/nevalang/neva/pkg/core"
)

// FuzzParserParseFiles verifies that arbitrary source returns a diagnostic
// instead of panicking while the parser builds the semantic tree.
func FuzzParserParseFiles(f *testing.F) {
	for _, source := range [][]byte{
		nil,
		[]byte("def Main(start any) (stop any) {\n\t:start -> :stop\n}\n"),
		[]byte("def Main(start any) (stop any) {\n"),
		[]byte("// comment\n\xff\n"),
	} {
		f.Add(source)
	}

	f.Fuzz(func(t *testing.T, source []byte) {
		if _, err := New().ParseFiles(
			core.ModuleRef{},
			"fuzz",
			map[string][]byte{"main.neva": source},
		); err != nil && err.Message == "" {
			t.Fatal("ParseFiles returned an error without a diagnostic")
		}
	})
}

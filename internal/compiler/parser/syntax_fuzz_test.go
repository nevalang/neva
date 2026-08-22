package parser

import (
	"testing"

	"github.com/nevalang/neva/pkg/core"
)

// FuzzParserParseFiles checks that arbitrary source either parses or returns a
// diagnostic with its input location, never panicking during tree construction.
func FuzzParserParseFiles(f *testing.F) {
	const (
		packageName = "fuzz"
		fileName    = "main.neva"
	)
	modRef := core.ModuleRef{Path: "fuzz"}
	wantLocation := core.Location{
		ModRef:   modRef,
		Package:  packageName,
		Filename: fileName,
	}

	// These seeds run in ordinary go test. With -fuzz, Go mutates their []byte
	// values and retains variants that cover new parser or listener paths.
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
			modRef,
			packageName,
			map[string][]byte{fileName: source},
		); err != nil {
			if err.Message == "" {
				t.Fatal("ParseFiles returned an error without a message")
			}
			if err.Meta == nil {
				t.Fatal("ParseFiles returned an error without metadata")
			}
			if err.Meta.Location != wantLocation {
				t.Fatalf("error location = %#v, want %#v", err.Meta.Location, wantLocation)
			}
		}
	})
}

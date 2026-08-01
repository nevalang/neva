package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSourcePreservesSourceTokensAndTree(t *testing.T) {
	t.Parallel()

	source := []byte("// Greets the caller\n" +
		"def Main(start any) (stop any) {\n" +
		"\t:start -> :stop\n" +
		"}\n")

	parsed, err := ParseSource(source)
	require.Nil(t, err)
	require.Equal(t, source, parsed.Source)
	require.NotNil(t, parsed.Tree)
	require.NotEmpty(t, parsed.Tokens)

	var comments, newlines int
	for _, token := range parsed.Tokens {
		switch {
		case strings.HasPrefix(token.GetText(), "//"):
			comments++
		case token.GetText() == "\n":
			newlines++
		}
	}
	require.Equal(t, 1, comments)
	require.Equal(t, 4, newlines)
}

func TestParseSourceReturnsSyntaxError(t *testing.T) {
	t.Parallel()

	_, err := ParseSource([]byte("def Main(start any) (stop any) {\n"))
	require.Error(t, err)
	require.Equal(t, 2, err.Meta.Start.Line)
}

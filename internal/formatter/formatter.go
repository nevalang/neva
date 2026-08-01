// Package formatter renders syntax-valid Neva source into its canonical layout.
package formatter

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/parser"
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
)

// Format parses and formats one Neva source file. Invalid source is never
// formatted; the returned compiler error contains the parser diagnostic.
func Format(source []byte) ([]byte, *compiler.Error) {
	parsed, err := parser.ParseSource(source)
	if err != nil {
		return nil, err
	}

	return FormatParsed(parsed), nil
}

// FormatParsed renders syntax-valid parsed source. It has no dependency on
// module resolution, analysis, desugaring, or code generation.
func FormatParsed(parsed parser.ParsedSource) []byte {
	if len(parsed.Tokens) == 1 {
		return nil
	}

	layout := newLayoutAnnotations(parsed.Tokens)
	antlr.ParseTreeWalkerDefault.Walk(layout, parsed.Tree)

	lines := strings.Split(strings.ReplaceAll(string(parsed.Source), "\r\n", "\n"), "\n")
	byLine := tokensByLine(parsed.Tokens)

	var output strings.Builder
	depth := 0
	for line := 1; line <= len(lines); line++ {
		tokens := byLine[line]
		if len(tokens) == 0 {
			if line < len(lines) {
				output.WriteByte('\n')
			}
			continue
		}

		lineDepth := max(depth-leadingClosers(tokens), 0)
		output.WriteString(strings.Repeat("\t", lineDepth))
		output.WriteString(renderLine(tokens, layout))
		output.WriteByte('\n')

		depth += braceDelta(tokens)
	}

	return []byte(output.String())
}

type indexedToken struct {
	text  string
	index int
}

func tokensByLine(tokens []antlr.Token) map[int][]indexedToken {
	result := make(map[int][]indexedToken)
	for index, token := range tokens {
		if token.GetTokenType() == antlr.TokenEOF || token.GetText() == "\n" || token.GetText() == "\r\n" {
			continue
		}
		result[token.GetLine()] = append(result[token.GetLine()], indexedToken{
			index: index,
			text:  token.GetText(),
		})
	}
	return result
}

func leadingClosers(tokens []indexedToken) int {
	closers := 0
	for _, token := range tokens {
		if token.text != "}" {
			break
		}
		closers++
	}
	return closers
}

func braceDelta(tokens []indexedToken) int {
	delta := 0
	for _, token := range tokens {
		switch token.text {
		case "{":
			delta++
		case "}":
			delta--
		}
	}
	return delta
}

func renderLine(tokens []indexedToken, layout *layoutAnnotations) string {
	var output strings.Builder
	for index, token := range tokens {
		if index > 0 && needsSpace(tokens[index-1], token, layout) {
			output.WriteByte(' ')
		}
		output.WriteString(token.text)
	}
	return output.String()
}

func needsSpace(previous, current indexedToken, layout *layoutAnnotations) bool {
	if space, handled := spaceBefore(current, previous, layout); handled {
		return space
	}

	return spaceAfter(previous, layout)
}

func spaceBefore(current, previous indexedToken, layout *layoutAnnotations) (bool, bool) {
	if strings.HasPrefix(current.text, "//") {
		return true, true
	}

	switch current.text {
	case ")", "]", ",", ".", "?", "::":
		return false, true
	case ":":
		return previous.text == "->", true
	case "}":
		return !layout.nodeDIArgBraces[current.index] && previous.text != "{", true
	case "(":
		return previous.text == ")", true
	case "[", "<", ">":
		return false, true
	case "{":
		return !layout.nodeDIArgBraces[current.index], true
	case "-":
		return previous.text != "(" && previous.text != "[" && previous.text != "=", true
	}
	return false, false
}

func spaceAfter(previous indexedToken, layout *layoutAnnotations) bool {
	switch previous.text {
	case "(", "[", "<", ">", ".", "$", "@", "#", "::":
		return false
	case "{":
		return !layout.nodeDIArgBraces[previous.index]
	case ":":
		return layout.structValueFieldColons[previous.index]
	case ",", "=", "->":
		return true
	case "-":
		return false
	}

	return true
}

type layoutAnnotations struct {
	generated.BasenevaListener
	structValueFieldColons map[int]bool
	nodeDIArgBraces        map[int]bool
	tokens                 []antlr.Token
}

func newLayoutAnnotations(tokens []antlr.Token) *layoutAnnotations {
	return &layoutAnnotations{
		tokens:                 tokens,
		structValueFieldColons: make(map[int]bool),
		nodeDIArgBraces:        make(map[int]bool),
	}
}

func (l *layoutAnnotations) EnterStructValueField(ctx *generated.StructValueFieldContext) {
	l.markFirst(ctx, ":", l.structValueFieldColons)
}

func (l *layoutAnnotations) EnterNodeDIArgs(ctx *generated.NodeDIArgsContext) {
	l.markFirst(ctx, "{", l.nodeDIArgBraces)
	l.markLast(ctx, "}", l.nodeDIArgBraces)
}

func (l *layoutAnnotations) markFirst(ctx antlr.ParserRuleContext, text string, marks map[int]bool) {
	for index := ctx.GetStart().GetTokenIndex(); index <= ctx.GetStop().GetTokenIndex(); index++ {
		if l.tokens[index].GetText() == text {
			marks[index] = true
			return
		}
	}
}

func (l *layoutAnnotations) markLast(ctx antlr.ParserRuleContext, text string, marks map[int]bool) {
	for index := ctx.GetStop().GetTokenIndex(); index >= ctx.GetStart().GetTokenIndex(); index-- {
		if l.tokens[index].GetText() == text {
			marks[index] = true
			return
		}
	}
}

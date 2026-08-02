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

	// Mark syntax contexts whose delimiters require exceptional spacing.
	layout := newLayoutAnnotations(parsed.Tokens)
	antlr.ParseTreeWalkerDefault.Walk(layout, parsed.Tree)

	return renderTokens(parsed.Tokens, layout)
}

// renderTokens preserves source line breaks except inside non-empty composite
// literals, whose parser contexts define their canonical multiline layout.
//
//nolint:gocognit,gocyclo,cyclop // Token emission must combine source breaks and parse-tree layout marks.
func renderTokens(tokens []antlr.Token, layout *layoutAnnotations) []byte {
	var output strings.Builder
	var line []indexedToken

	depth, previousWasNewline := 0, false
	flush := func() {
		if len(line) == 0 {
			return
		}
		lineDepth := max(depth-leadingClosers(line, layout), 0)
		output.WriteString(strings.Repeat("\t", lineDepth))
		output.WriteString(renderLine(line, layout))
		output.WriteByte('\n')
		depth += delimiterDelta(line, layout)
		line = nil
	}

	for index, token := range tokens {
		text := token.GetText()
		if token.GetTokenType() == antlr.TokenEOF {
			break
		}
		//nolint:nestif // A source newline is either preserved, collapsed, or suppressed by literal layout.
		if text == "\n" || text == "\r\n" {
			if !layout.compositeNewlines[index] {
				if len(line) > 0 {
					flush()
				} else if previousWasNewline {
					output.WriteByte('\n')
				}
				previousWasNewline = true
			}
			continue
		}
		previousWasNewline = false

		if layout.compositeCloses[index] {
			if len(line) > 0 && line[len(line)-1].text != "," {
				line = append(line, indexedToken{text: ",", index: -1})
			}
			flush()
		}

		line = append(line, indexedToken{text: text, index: index})
		if layout.compositeOpens[index] || layout.compositeCommas[index] {
			flush()
		}
	}

	flush()
	return []byte(output.String())
}

// indexedToken retains a token's source-stream index with its printable text.
type indexedToken struct {
	text  string
	index int
}

// leadingClosers reports how many layout delimiters close at a source line's start.
func leadingClosers(tokens []indexedToken, layout *layoutAnnotations) int {
	closers := 0
	for _, token := range tokens {
		if token.text != "}" && !layout.compositeCloses[token.index] {
			break
		}
		closers++
	}
	return closers
}

// delimiterDelta reports the net indentation change produced by one source line.
func delimiterDelta(tokens []indexedToken, layout *layoutAnnotations) int {
	delta := 0
	for _, token := range tokens {
		switch token.text {
		case "{":
			delta++
		case "}":
			delta--
		case "[":
			if layout.compositeOpens[token.index] {
				delta++
			}
		case "]":
			if layout.compositeCloses[token.index] {
				delta--
			}
		}
	}
	return delta
}

// renderLine joins a line's tokens with their canonical horizontal spacing.
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

// needsSpace reports whether a space belongs between two adjacent tokens.
func needsSpace(previous, current indexedToken, layout *layoutAnnotations) bool {
	if space, handled := spaceBefore(current, previous, layout); handled {
		return space
	}

	return spaceAfter(previous, layout)
}

// spaceBefore handles spacing rules determined by the current token.
//
//nolint:gocyclo,cyclop // Token categories map directly to the language punctuation table.
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
	case "[":
		return spaceBeforeOpenBracket(previous), true
	case "<", ">":
		return false, true
	case "{":
		return !layout.nodeDIArgBraces[current.index], true
	case "-":
		return previous.text != "(" && previous.text != "[" && previous.text != "=", true
	}
	return false, false
}

// spaceBeforeOpenBracket keeps square brackets attached except after an operator.
func spaceBeforeOpenBracket(previous indexedToken) bool {
	return previous.text == "=" || previous.text == "->"
}

// spaceAfter handles spacing rules determined by the previous token.
func spaceAfter(previous indexedToken, layout *layoutAnnotations) bool {
	switch previous.text {
	case "(", "[", "<", ".", "$", "@", "#", "::":
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

// layoutAnnotations records syntax-context exceptions to generic token spacing.
type layoutAnnotations struct {
	generated.BasenevaListener
	structValueFieldColons map[int]bool
	nodeDIArgBraces        map[int]bool
	compositeOpens         map[int]bool
	compositeCloses        map[int]bool
	compositeCommas        map[int]bool
	compositeNewlines      map[int]bool
	tokens                 []antlr.Token
}

// newLayoutAnnotations creates a listener backed by one complete token stream.
func newLayoutAnnotations(tokens []antlr.Token) *layoutAnnotations {
	return &layoutAnnotations{
		tokens:                 tokens,
		structValueFieldColons: make(map[int]bool),
		nodeDIArgBraces:        make(map[int]bool),
		compositeOpens:         make(map[int]bool),
		compositeCloses:        make(map[int]bool),
		compositeCommas:        make(map[int]bool),
		compositeNewlines:      make(map[int]bool),
	}
}

// EnterListLit marks non-empty list literals for canonical multiline rendering.
func (l *layoutAnnotations) EnterListLit(ctx *generated.ListLitContext) {
	items, ok := ctx.ListItems().(antlr.ParserRuleContext)
	if !ok {
		return
	}
	l.markComposite(ctx, items)
}

// EnterStructLit marks non-empty struct literals for canonical multiline rendering.
func (l *layoutAnnotations) EnterStructLit(ctx *generated.StructLitContext) {
	fields, ok := ctx.StructValueFields().(antlr.ParserRuleContext)
	if !ok {
		return
	}
	l.markComposite(ctx, fields)
}

// markComposite records one literal's delimiter, separator, and newline layout.
func (l *layoutAnnotations) markComposite(ctx antlr.ParserRuleContext, items antlr.ParserRuleContext) {
	start := ctx.GetStart().GetTokenIndex()
	stop := ctx.GetStop().GetTokenIndex()
	l.compositeOpens[start] = true
	l.compositeCloses[stop] = true
	for index := start; index <= stop; index++ {
		if text := l.tokens[index].GetText(); text == "\n" || text == "\r\n" {
			l.compositeNewlines[index] = true
		}
	}
	for _, child := range items.GetChildren() {
		terminal, ok := child.(antlr.TerminalNode)
		if ok && terminal.GetText() == "," {
			l.compositeCommas[terminal.GetSymbol().GetTokenIndex()] = true
		}
	}
}

// EnterStructValueField marks a struct-literal colon, which requires a trailing space.
func (l *layoutAnnotations) EnterStructValueField(ctx *generated.StructValueFieldContext) {
	l.markFirst(ctx, ":", l.structValueFieldColons)
}

// EnterNodeDIArgs marks dependency-injection braces, which bind to their contents.
func (l *layoutAnnotations) EnterNodeDIArgs(ctx *generated.NodeDIArgsContext) {
	l.markFirst(ctx, "{", l.nodeDIArgBraces)
	l.markLast(ctx, "}", l.nodeDIArgBraces)
}

// markFirst records the first matching token inside one parser-rule context.
func (l *layoutAnnotations) markFirst(ctx antlr.ParserRuleContext, text string, marks map[int]bool) {
	for index := ctx.GetStart().GetTokenIndex(); index <= ctx.GetStop().GetTokenIndex(); index++ {
		if l.tokens[index].GetText() == text {
			marks[index] = true
			return
		}
	}
}

// markLast records the last matching token inside one parser-rule context.
func (l *layoutAnnotations) markLast(ctx antlr.ParserRuleContext, text string, marks map[int]bool) {
	for index := ctx.GetStop().GetTokenIndex(); index >= ctx.GetStart().GetTokenIndex(); index-- {
		if l.tokens[index].GetText() == text {
			marks[index] = true
			return
		}
	}
}

// Package formatter renders syntax-valid Neva source into its canonical layout.
package formatter

import (
	"sort"
	"strings"
	"unicode/utf8"

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
//nolint:funlen,gocognit,gocyclo,cyclop // Token emission must combine source breaks and parse-tree layout marks.
func renderTokens(tokens []antlr.Token, layout *layoutAnnotations) []byte {
	var output strings.Builder
	var line []indexedToken

	depth, previousWasNewline := 0, false
	flush := func() {
		if len(line) == 0 {
			return
		}
		lineDepth := max(depth-leadingClosers(line, layout), 0)
		renderWrappedLine(&output, line, layout, lineDepth)
		depth += delimiterDelta(line, layout)
		line = nil
	}

	for index := 0; index < len(tokens); index++ {
		if importBlock, ok := layout.importBlocks[index]; ok {
			flush()
			output.Write(renderImportBlock(importBlock, layout))
			previousWasNewline = false
			index = importBlock.end
			continue
		}

		token := tokens[index]
		text := token.GetText()
		if token.GetTokenType() == antlr.TokenEOF {
			break
		}
		//nolint:nestif // A source newline is either preserved, collapsed, or suppressed by literal layout.
		if text == "\n" || text == "\r\n" {
			if layout.declarationNewlines[index] || (!layout.compositeNewlines[index] && !layout.sequenceNewlines[index]) {
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

		if composite, ok := layout.composites[index]; ok {
			if !compositeExceedsLineWidth(line, composite, layout, depth) {
				line = append(line, layout.compactTokens(composite.start, composite.stop)...)
				index = composite.stop
				continue
			}
			layout.activateComposite(composite)
		}
		if sequence, ok := layout.sequences[index]; ok {
			if !layout.sequenceAlwaysVertical[sequence.start] &&
				!compositeExceedsLineWidth(line, sequence, layout, depth) {
				line = append(line, layout.compactTokens(sequence.start, sequence.stop)...)
				index = sequence.stop
				continue
			}
			layout.activateSequence(sequence)
		}

		if layout.compositeCloses[index] {
			if len(line) > 0 && line[len(line)-1].text != "," {
				line = append(line, indexedToken{text: ",", index: -1})
			}
			flush()
		}
		if layout.sequenceCloses[index] {
			if layout.sequenceTrailing[index] && len(line) > 0 && line[len(line)-1].text != "," {
				line = append(line, indexedToken{text: ",", index: -1})
			}
			flush()
		}
		if layout.declarationCloses[index] {
			flush()
		}

		line = append(line, indexedToken{text: text, index: index})
		if layout.compositeOpens[index] || layout.compositeCommas[index] ||
			layout.sequenceOpens[index] || layout.sequenceCommas[index] || layout.declarationOpens[index] {
			flush()
		}
	}

	flush()
	return []byte(output.String())
}

// compositeExceedsLineWidth reports whether a literal can share the current
// logical line in its compact form. Multiline literals are deliberately
// decided as a unit: once one does not fit, each item is rendered vertically.
func compositeExceedsLineWidth(line []indexedToken, composite compositeLayout, layout *layoutAnnotations, depth int) bool {
	for index := composite.start; index <= composite.stop; index++ {
		if layout.declarationOpens[index] {
			return true
		}
	}

	flat := append(append([]indexedToken(nil), line...), layout.compactTokens(composite.start, composite.stop)...)
	indent := strings.Repeat("\t", depth-leadingClosers(line, layout))
	return utf8.RuneCountInString(indent)+utf8.RuneCountInString(renderLine(flat, layout)) > maxLineLength
}

const maxLineLength = 80

// renderWrappedLine emits one logical line, inserting safe syntax-level
// line breaks only when its canonical rendering exceeds maxLineLength.
func renderWrappedLine(output *strings.Builder, tokens []indexedToken, layout *layoutAnnotations, depth int) {
	continuationDepth := depth + 1
	for start := 0; start < len(tokens); {
		indent := strings.Repeat("\t", depth)
		end, breakAfter := len(tokens), -1
		for candidate := start + 1; candidate <= len(tokens); candidate++ {
			line := renderLine(tokens[start:candidate], layout)
			if utf8.RuneCountInString(indent)+utf8.RuneCountInString(line) > maxLineLength {
				if breakAfter >= 0 {
					end = breakAfter
				}
				break
			}
		}

		output.WriteString(indent)
		output.WriteString(renderLine(tokens[start:end], layout))
		output.WriteByte('\n')
		if end == len(tokens) {
			return
		}
		start = end
		depth = continuationDepth
	}
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
		if token.text != "}" && !layout.compositeCloses[token.index] && !layout.sequenceCloses[token.index] {
			break
		}
		closers++
	}
	return closers
}

// delimiterDelta reports the net indentation change produced by one source line.
//
//nolint:gocognit,gocyclo,cyclop // Delimiter cases directly encode the formatter's indentation rules.
func delimiterDelta(tokens []indexedToken, layout *layoutAnnotations) int {
	delta := 0
	for _, token := range tokens {
		if layout.sequenceOpens[token.index] && token.text != "[" && token.text != "{" {
			delta++
		}
		if layout.sequenceCloses[token.index] && token.text != "]" && token.text != "}" {
			delta--
		}
		switch token.text {
		case "{":
			delta++
		case "}":
			delta--
		case "[":
			if layout.compositeOpens[token.index] || layout.sequenceOpens[token.index] {
				delta++
			}
		case "]":
			if layout.compositeCloses[token.index] || layout.sequenceCloses[token.index] {
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
	case ")", "]", ",", "::":
		return false, true
	case ".", "?":
		return previous.text == "->", true
	case ":":
		return previous.text == "->" || previous.text == "," || previous.text == "{", true
	case "/":
		return !layout.importPathSlashes[current.index], true
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
	return previous.text == "=" || previous.text == "->" || previous.text == ","
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
	case "/":
		return !layout.importPathSlashes[previous.index]
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
	importPathSlashes      map[int]bool
	compositeOpens         map[int]bool
	compositeCloses        map[int]bool
	compositeCommas        map[int]bool
	compositeNewlines      map[int]bool
	composites             map[int]compositeLayout
	trailingCommas         map[int]bool
	sequenceOpens          map[int]bool
	sequenceCloses         map[int]bool
	sequenceCommas         map[int]bool
	sequenceNewlines       map[int]bool
	sequenceTrailing       map[int]bool
	sequenceAlwaysVertical map[int]bool
	sequences              map[int]compositeLayout
	declarationOpens       map[int]bool
	declarationCloses      map[int]bool
	declarationNewlines    map[int]bool
	importBlocks           map[int]importBlock
	tokens                 []antlr.Token
}

// newLayoutAnnotations creates a listener backed by one complete token stream.
func newLayoutAnnotations(tokens []antlr.Token) *layoutAnnotations {
	return &layoutAnnotations{
		tokens:                 tokens,
		structValueFieldColons: make(map[int]bool),
		nodeDIArgBraces:        make(map[int]bool),
		importPathSlashes:      make(map[int]bool),
		compositeOpens:         make(map[int]bool),
		compositeCloses:        make(map[int]bool),
		compositeCommas:        make(map[int]bool),
		compositeNewlines:      make(map[int]bool),
		composites:             make(map[int]compositeLayout),
		trailingCommas:         make(map[int]bool),
		sequenceOpens:          make(map[int]bool),
		sequenceCloses:         make(map[int]bool),
		sequenceCommas:         make(map[int]bool),
		sequenceNewlines:       make(map[int]bool),
		sequenceTrailing:       make(map[int]bool),
		sequenceAlwaysVertical: make(map[int]bool),
		sequences:              make(map[int]compositeLayout),
		declarationOpens:       make(map[int]bool),
		declarationCloses:      make(map[int]bool),
		declarationNewlines:    make(map[int]bool),
		importBlocks:           make(map[int]importBlock),
	}
}

// EnterImportPath marks package path separators, which have no surrounding spaces.
func (l *layoutAnnotations) EnterImportPath(ctx *generated.ImportPathContext) {
	for index := ctx.GetStart().GetTokenIndex(); index <= ctx.GetStop().GetTokenIndex(); index++ {
		if l.tokens[index].GetText() == "/" {
			l.importPathSlashes[index] = true
		}
	}
}

// EnterPortsDef marks grammar-approved line-break positions in a port list.
func (l *layoutAnnotations) EnterPortsDef(ctx *generated.PortsDefContext) {
	l.markSequence(ctx, ctx, "(", ")")
}

// EnterTypeParams marks grammar-approved line-break positions in type parameters.
func (l *layoutAnnotations) EnterTypeParams(ctx *generated.TypeParamsContext) {
	params, ok := ctx.TypeParamList().(antlr.ParserRuleContext)
	if !ok {
		return
	}
	l.markSequence(ctx, params, "<", ">")
}

// EnterTypeArgs marks grammar-approved line-break positions in type arguments.
func (l *layoutAnnotations) EnterTypeArgs(ctx *generated.TypeArgsContext) {
	l.markSequence(ctx, ctx, "<", ">")
}

// EnterStructTypeExpr makes non-empty struct declarations into blocks.
func (l *layoutAnnotations) EnterStructTypeExpr(ctx *generated.StructTypeExprContext) {
	if ctx.StructFields() != nil {
		l.markDeclarationBlock(ctx)
	}
}

// EnterUnionTypeExpr makes non-empty union declarations into blocks.
func (l *layoutAnnotations) EnterUnionTypeExpr(ctx *generated.UnionTypeExprContext) {
	if ctx.UnionFields() != nil {
		l.markDeclarationBlock(ctx)
	}
}

// EnterMultipleSenderSide marks grammar-approved line-break positions in fan-in.
func (l *layoutAnnotations) EnterMultipleSenderSide(ctx *generated.MultipleSenderSideContext) {
	l.markSequence(ctx, ctx, "[", "]")
	l.markAlwaysVerticalSequence(ctx, "[")
}

// EnterMultipleReceiverSide marks grammar-approved line-break positions in fan-out.
func (l *layoutAnnotations) EnterMultipleReceiverSide(ctx *generated.MultipleReceiverSideContext) {
	l.markSequence(ctx, ctx, "[", "]")
	l.markAlwaysVerticalSequence(ctx, "[")
}

// EnterCompBody renders every non-empty component body as a block.
func (l *layoutAnnotations) EnterCompBody(ctx *generated.CompBodyContext) {
	open := l.firstToken(ctx, "{")
	closing := l.lastToken(ctx, "}")
	if open < 0 || closing < 0 || l.hasOnlyNewlines(open+1, closing) {
		return
	}
	l.markDeclarationBlock(ctx)
}

// EnterImportStmt records one import block for canonical ordering and grouping.
func (l *layoutAnnotations) EnterImportStmt(ctx *generated.ImportStmtContext) {
	entries := make([]importEntry, 0, len(ctx.AllImportBlockItem()))
	var leadingComments []string

	for _, item := range ctx.AllImportBlockItem() {
		if definition := item.ImportDef(); definition != nil {
			def, ok := definition.(antlr.ParserRuleContext)
			if !ok {
				continue
			}
			path := definition.ImportPath().GetText()
			entries = append(entries, importEntry{
				comments:  leadingComments,
				group:     classifyImport(path),
				start:     def.GetStart().GetTokenIndex(),
				stop:      def.GetStop().GetTokenIndex(),
				pathStart: definition.ImportPath().GetStart().GetTokenIndex(),
				pathStop:  definition.ImportPath().GetStop().GetTokenIndex(),
			})
			leadingComments = nil
			continue
		}

		if comment := item.COMMENT(); comment != nil {
			leadingComments = append(leadingComments, comment.GetText())
		}
	}

	if len(entries) == 0 {
		return
	}
	entries[len(entries)-1].comments = append(entries[len(entries)-1].comments, leadingComments...)
	start := ctx.GetStart().GetTokenIndex()
	l.importBlocks[start] = importBlock{end: ctx.GetStop().GetTokenIndex(), entries: entries}
}

// tokenRange returns printable tokens in one inclusive source-token range.
func tokenRange(tokens []antlr.Token, start, stop int) []indexedToken {
	line := make([]indexedToken, 0, stop-start+1)
	for index := start; index <= stop; index++ {
		text := tokens[index].GetText()
		if text != "\n" && text != "\r\n" {
			line = append(line, indexedToken{text: text, index: index})
		}
	}
	return line
}

type importGroup uint8

const (
	stdlibImportGroup importGroup = iota
	thirdPartyImportGroup
	localImportGroup
)

// importEntry keeps one import and the comments that travel with it.
type importEntry struct {
	comments  []string
	start     int
	stop      int
	pathStart int
	pathStop  int
	group     importGroup
}

// importBlock is a complete import statement marked by its token range.
type importBlock struct {
	entries []importEntry
	end     int
}

// classifyImport assigns the source-level import group defined by the style guide.
func classifyImport(path string) importGroup {
	switch {
	case strings.HasPrefix(path, "@:"):
		return localImportGroup
	case strings.Contains(path, ":"):
		return thirdPartyImportGroup
	default:
		return stdlibImportGroup
	}
}

// renderImportBlock sorts imports within their source-level groups and emits
// group spacing defined by the style guide.
func renderImportBlock(block importBlock, layout *layoutAnnotations) []byte {
	groups := [3][]importEntry{}
	for _, entry := range block.entries {
		groups[entry.group] = append(groups[entry.group], entry)
	}

	separateGroups := false
	for index := range groups {
		sort.SliceStable(groups[index], func(left, right int) bool {
			return tokenText(layout.tokens, groups[index][left].pathStart, groups[index][left].pathStop) <
				tokenText(layout.tokens, groups[index][right].pathStart, groups[index][right].pathStop)
		})
		separateGroups = separateGroups || len(groups[index]) > 2
	}

	var output strings.Builder
	output.WriteString("import {\n")
	for groupIndex, group := range groups {
		if len(group) == 0 {
			continue
		}
		if groupIndex > 0 && separateGroups && output.Len() > len("import {\n") {
			output.WriteByte('\n')
		}
		for _, entry := range group {
			for _, comment := range entry.comments {
				output.WriteByte('\t')
				output.WriteString(comment)
				output.WriteByte('\n')
			}
			output.WriteByte('\t')
			output.WriteString(renderLine(tokenRange(layout.tokens, entry.start, entry.stop), layout))
			output.WriteByte('\n')
		}
	}
	output.WriteString("}\n")
	return []byte(output.String())
}

// tokenText joins printable tokens in one inclusive range.
func tokenText(tokens []antlr.Token, start, stop int) string {
	var output strings.Builder
	for index := start; index <= stop; index++ {
		text := tokens[index].GetText()
		if text != "\n" && text != "\r\n" {
			output.WriteString(text)
		}
	}
	return output.String()
}

// EnterListLit records a non-empty list literal for width-directed rendering.
func (l *layoutAnnotations) EnterListLit(ctx *generated.ListLitContext) {
	items, ok := ctx.ListItems().(antlr.ParserRuleContext)
	if !ok {
		return
	}
	l.markComposite(ctx, items)
}

// EnterStructLit records a non-empty struct literal for width-directed rendering.
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
	composite := compositeLayout{start: start, stop: stop}
	for index := start; index <= stop; index++ {
		if text := l.tokens[index].GetText(); text == "\n" || text == "\r\n" {
			l.compositeNewlines[index] = true
		}
	}
	for _, child := range items.GetChildren() {
		terminal, ok := child.(antlr.TerminalNode)
		if ok && terminal.GetText() == "," {
			composite.commas = append(composite.commas, terminal.GetSymbol().GetTokenIndex())
		}
	}
	if len(composite.commas) > 0 {
		last := composite.commas[len(composite.commas)-1]
		if l.hasOnlyNewlines(last+1, stop) {
			l.trailingCommas[last] = true
		}
	}
	l.composites[start] = composite
}

// compositeLayout identifies one non-empty list or struct literal.
type compositeLayout struct {
	commas   []int
	start    int
	stop     int
	trailing bool
}

// activateComposite selects the vertical form for one literal.
func (l *layoutAnnotations) activateComposite(composite compositeLayout) {
	l.compositeOpens[composite.start] = true
	l.compositeCloses[composite.stop] = true
	for _, comma := range composite.commas {
		l.compositeCommas[comma] = true
	}
}

// markSequence records a comma-separated sequence for width-directed rendering.
func (l *layoutAnnotations) markSequence(ctx, items antlr.ParserRuleContext, open, closing string) {
	sequence := compositeLayout{
		start:    l.firstToken(ctx, open),
		stop:     l.lastToken(ctx, closing),
		trailing: true,
	}
	if sequence.start < 0 || sequence.stop < 0 {
		return
	}
	for index := sequence.start; index <= sequence.stop; index++ {
		text := l.tokens[index].GetText()
		if text == "\n" || text == "\r\n" {
			l.sequenceNewlines[index] = true
		}
	}
	for _, child := range items.GetChildren() {
		terminal, ok := child.(antlr.TerminalNode)
		if ok && terminal.GetText() == "," {
			sequence.commas = append(sequence.commas, terminal.GetSymbol().GetTokenIndex())
		}
	}
	if len(sequence.commas) > 0 {
		last := sequence.commas[len(sequence.commas)-1]
		if l.hasOnlyNewlines(last+1, sequence.stop) {
			l.trailingCommas[last] = true
		}
	}
	l.sequences[sequence.start] = sequence
}

// markAlwaysVerticalSequence records a graph branch sequence that never uses compact layout.
func (l *layoutAnnotations) markAlwaysVerticalSequence(ctx antlr.ParserRuleContext, open string) {
	start := l.firstToken(ctx, open)
	sequence, ok := l.sequences[start]
	if ok && len(sequence.commas) > 0 {
		l.sequenceAlwaysVertical[start] = true
	}
}

// markDeclarationBlock records the braces that delimit a multi-item type declaration.
func (l *layoutAnnotations) markDeclarationBlock(ctx antlr.ParserRuleContext) {
	open := l.firstToken(ctx, "{")
	closing := l.lastToken(ctx, "}")
	if open < 0 || closing < 0 {
		return
	}
	l.declarationOpens[open] = true
	l.declarationCloses[closing] = true
	for index := open; index <= closing; index++ {
		text := l.tokens[index].GetText()
		if text == "\n" || text == "\r\n" {
			l.declarationNewlines[index] = true
		}
	}
}

// firstToken returns the first matching token index in one parser-rule context.
func (l *layoutAnnotations) firstToken(ctx antlr.ParserRuleContext, text string) int {
	for index := ctx.GetStart().GetTokenIndex(); index <= ctx.GetStop().GetTokenIndex(); index++ {
		if l.tokens[index].GetText() == text {
			return index
		}
	}
	return -1
}

// lastToken returns the last matching token index in one parser-rule context.
func (l *layoutAnnotations) lastToken(ctx antlr.ParserRuleContext, text string) int {
	for index := ctx.GetStop().GetTokenIndex(); index >= ctx.GetStart().GetTokenIndex(); index-- {
		if l.tokens[index].GetText() == text {
			return index
		}
	}
	return -1
}

// activateSequence selects the vertical form for one graph or union sequence.
func (l *layoutAnnotations) activateSequence(sequence compositeLayout) {
	l.sequenceOpens[sequence.start] = true
	l.sequenceCloses[sequence.stop] = true
	for _, comma := range sequence.commas {
		l.sequenceCommas[comma] = true
	}
	if sequence.trailing {
		l.sequenceTrailing[sequence.stop] = true
	}
}

// compactTokens returns one literal without source newlines or trailing commas.
func (l *layoutAnnotations) compactTokens(start, stop int) []indexedToken {
	tokens := make([]indexedToken, 0, stop-start+1)
	for index := start; index <= stop; index++ {
		text := l.tokens[index].GetText()
		if text == "\n" || text == "\r\n" || l.trailingCommas[index] {
			continue
		}
		tokens = append(tokens, indexedToken{text: text, index: index})
	}
	return tokens
}

// hasOnlyNewlines reports whether the range before a closing delimiter is empty.
func (l *layoutAnnotations) hasOnlyNewlines(start, stop int) bool {
	for index := start; index < stop; index++ {
		text := l.tokens[index].GetText()
		if text != "\n" && text != "\r\n" {
			return false
		}
	}
	return true
}

// EnterStructValueField marks a struct-literal colon, which requires a trailing space.
func (l *layoutAnnotations) EnterStructValueField(ctx *generated.StructValueFieldContext) {
	l.markFirst(ctx, ":", l.structValueFieldColons)
}

// EnterNodeDIArgs renders every non-empty dependency-injection block vertically.
func (l *layoutAnnotations) EnterNodeDIArgs(ctx *generated.NodeDIArgsContext) {
	l.markFirst(ctx, "{", l.nodeDIArgBraces)
	l.markLast(ctx, "}", l.nodeDIArgBraces)
	l.markDeclarationBlock(ctx)
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

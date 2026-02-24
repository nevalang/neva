package smoke_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	"github.com/stretchr/testify/require"
)

type TreeShapeListener struct {
	*generated.BasenevaListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}

// MyErrorListener is a copy of antlr.ErrorListener just to generate mock
type MyErrorListener interface {
	antlr.ErrorListener
}

// FileAwareErrorListener provides better error reporting with file context
//
//nolint:govet // fieldalignment: small helper struct.
type FileAwareErrorListener struct {
	filename string
	t        *testing.T
}

func NewFileAwareErrorListener(t *testing.T, filename string) *FileAwareErrorListener {
	t.Helper()
	return &FileAwareErrorListener{
		filename: filename,
		t:        t,
	}
}

func (f *FileAwareErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol any, line, column int, msg string, e antlr.RecognitionException) {
	tokenText := "<unknown>"
	if token, ok := offendingSymbol.(antlr.Token); ok {
		tokenText = token.GetText()
	}
	f.t.Errorf("PARSER ERROR in %s at line %d:%d - %s\n  Token: '%s'",
		f.filename, line, column, msg, tokenText)
}

func (f *FileAwareErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	// Ignore ambiguity reports for now
}

func (f *FileAwareErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	// Ignore full context reports for now
}

func (f *FileAwareErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs *antlr.ATNConfigSet) {
	// Ignore context sensitivity reports for now
}

// TestSmoke reads all the ".neva" files in current directory and tries to parse them expecting zero errors.
func TestSmoke(t *testing.T) {
	err := os.Chdir("./happypath")
	require.NoError(t, err)

	nevaTestFiles, err := os.ReadDir(".")
	require.NoError(t, err)

	for _, file := range nevaTestFiles {
		fileName := file.Name()

		// skip current and mock files
		if !strings.HasSuffix(fileName, ".neva") {
			continue
		}

		fmt.Println("parsing was started for: ", fileName)

		// read file and create input
		input, err := antlr.NewFileStream(fileName)
		require.NoError(t, err)

		// create lexer and parser
		lexer := generated.NewnevaLexer(input)
		parser := generated.NewnevaParser(
			antlr.NewCommonTokenStream(lexer, 0),
		)

		// create file-aware error listener for better error reporting
		fileErrorListener := NewFileAwareErrorListener(t, fileName)
		parser.AddErrorListener(fileErrorListener)

		// create tree to walk
		parser.BuildParseTrees = true
		tree := parser.Prog()

		// walk the tree to catch potential errors
		antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
	}
}

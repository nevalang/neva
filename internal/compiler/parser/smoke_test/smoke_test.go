package smoke_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/golang/mock/gomock"

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

		fmt.Printf("Processing file: %s\n", fileName)

		// read file and create input
		input, err := antlr.NewFileStream(fileName)
		require.NoError(t, err)

		// create lexer and parser
		lexer := generated.NewnevaLexer(input)
		parser := generated.NewnevaParser(
			antlr.NewCommonTokenStream(lexer, 0),
		)

		// create mock and configure it to expect zero errors
		ctrl := gomock.NewController(t)
		mock := NewMockMyErrorListener(ctrl)
		initMock(mock.EXPECT())
		parser.AddErrorListener(mock)

		// create tree to walk
		parser.BuildParseTrees = true
		tree := parser.Prog()

		// walk the tree to catch potential errors
		antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
	}
}

// initMock configures the mock to expect zero calls
func initMock(recorder *MockMyErrorListenerMockRecorder) {
	// we don't care for now
	recorder.ReportContextSensitivity(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).AnyTimes()

	// we don't care for now
	recorder.ReportAmbiguity(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).AnyTimes()

	// we don't care for now
	recorder.ReportAttemptingFullContext(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).AnyTimes()

	// this is what we care about
	recorder.SyntaxError(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Times(0)
}

package frontend_test

import (
	"os"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	parsing "github.com/nevalang/neva/internal/compiler/frontend/generated"
)

type TreeShapeListener struct {
	*parsing.BasenevaListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}
type MyErrorListener interface {
	antlr.ErrorListener
}

// TestSmoke reads all the ".neva" files in current directory and tries to parse them expecting zero errors.
func TestSmoke(t *testing.T) {
	nevaTestFiles, err := os.ReadDir(".")
	require.NoError(t, err)

	for _, file := range nevaTestFiles {
		// skip current and mock files
		if !strings.HasSuffix(file.Name(), ".neva") {
			continue
		}

		// read file and create input
		input, err := antlr.NewFileStream(file.Name())
		require.NoError(t, err)

		// create lexer and parser
		lexer := parsing.NewnevaLexer(input)
		parser := parsing.NewnevaParser(
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
	recorder.SyntaxError(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Times(0)
}

package lsp

import (
	"sync"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/tliron/commonlog"
)

type Server struct {
	name, version string

	handler *Handler
	logger  commonlog.Logger
	indexer Indexer

	mu    *sync.Mutex
	state *State
}

type State struct {
	mod     src.Module
	problem string
}

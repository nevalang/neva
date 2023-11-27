package main

import (
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/tliron/glsp"
)

type FooBarParams struct {
	Document struct {
		Uri struct {
			Path   string
			fsPath string
		} `json:"uri"`
		FileName string `json:"fileName"`
	} `json:"document"`
}

type FooBarResp struct {
	File  src.File `json:"file"`
	Extra Extra    `json:"extra"` // info that is not presented in the file but needed for rendering
}

type Extra struct {
	NodesPorts map[string]src.Port `json:"nodesPorts"`
}

func (s *Server) FooBar(glspCtx *glsp.Context, params FooBarParams) (any, error) {
	s.logger.Info("FooBar")

	if s.indexedProgramState == nil {
		return nil, nil
	}

	// 1 get opened file
	// get nodes ports for every component in file
	// -- create scope with cur mod (entry), pkg and filename
	// -- use that scope to resolve all nodes'refs (interfaces and components)
	// -- gather their io and build response

	// TODO figure this out
	// var (
	// 	modName  = ""
	// 	pkgName  = ""
	// 	fileName = ""
	// )

	// scope := src.Scope{
	// 	Loc: src.ScopeLocation{
	// 		ModuleName: modName,
	// 		PkgName:    pkgName,
	// 		FileName:   fileName,
	// 	},
	// 	Module: *s.indexedProgramState,
	// }

	// TODO use scope to resolve all references in current file

	return FooBarResp{
		File: src.File{},
		Extra: Extra{
			NodesPorts: map[string]src.Port{},
		},
	}, nil
}

func (s *Server) setState(mod *src.Module, problem string) {
	s.mu.Lock()
	s.indexedProgramState = mod
	s.problem = problem
	s.mu.Unlock()
}

package main

import (
	"strings"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/tliron/glsp"
)

type FooBarRequest struct {
	WorkspaceUri URI `json:"workspaceUri"`
	Document     struct {
		Uri      URI    `json:"uri"`
		FileName string `json:"fileName"`
	} `json:"document"`
}

type URI struct {
	Path   string `json:"path"`
	FSPath string `json:"fsPath"`
}

type FooBarResp struct {
	File  src.File `json:"file"`
	Extra Extra    `json:"extra"` // info that is not presented in the file but needed for rendering
}

type Extra struct {
	NodesPorts map[string]map[string]map[string]src.Port `json:"nodesPorts"` // component -> node -> port
}

func (s *Server) FooBar(glspCtx *glsp.Context, req FooBarRequest) (any, error) {
	s.logger.Info("FooBar")

	if s.indexedProgramState == nil {
		return nil, nil
	}

	relPathToFile := strings.TrimPrefix(req.Document.FileName, req.WorkspaceUri.Path)
	fullFilePathParts := strings.Split(relPathToFile, "/")
	fileNameWithExt := fullFilePathParts[len(fullFilePathParts)-1]

	var (
		modName  = "entry"
		pkgName  = strings.Join(fullFilePathParts[0:len(fullFilePathParts)-1], "/")
		fileName = strings.TrimSuffix(fileNameWithExt, ".neva")
	)

	scope := src.Scope{
		Loc: src.ScopeLocation{
			ModuleName: modName,
			PkgName:    pkgName,
			FileName:   fileName,
		},
		Module: *s.indexedProgramState,
	}

	pkg := s.indexedProgramState.Packages[pkgName]
	file := pkg[fileName]

	extra := map[string]map[string]map[string]src.Port{} // c -> n -> p
	for entityName, entity := range file.Entities {
		if entity.Kind != src.ComponentEntity {
			continue
		}
		componentNodes := map[string]map[string]src.Port{}
		for _, node := range entity.Component.Nodes {
			// for every node find an entity
			// find out whether its component or interface
			// get ports and set to extra
			_, _, err := scope.Entity(node.EntityRef)
			if err != nil {
				panic(err)
			}
		}
		extra[entityName] = componentNodes
	}

	return FooBarResp{
		File: file,
		Extra: Extra{
			NodesPorts: extra,
		},
	}, nil
}

func (s *Server) setState(mod *src.Module, problem string) {
	s.mu.Lock()
	s.indexedProgramState = mod
	s.problem = problem
	s.mu.Unlock()
}

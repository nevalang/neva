package lsp

import (
	"errors"
	"strings"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/tliron/glsp"
)

type ResolveFileRequest struct {
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

type ResolveFileResponce struct {
	File  src.File `json:"file"`
	Extra Extra    `json:"extra"` // info that is not presented in the file but needed for rendering
}

type Extra struct {
	NodesPorts map[string]map[string]src.Interface `json:"nodesPorts"` // components -> nodes -> interface
}

func (s *Server) ResolveFile(glspCtx *glsp.Context, req ResolveFileRequest) (ResolveFileResponce, error) {
	if s.state == nil {
		return ResolveFileResponce{}, nil
	}

	relFilePath := strings.TrimPrefix(req.Document.FileName, req.WorkspaceUri.Path)
	relFilePath = strings.TrimPrefix(relFilePath, "/")

	relPathParts := strings.Split(relFilePath, "/")           // relative path to file in slice
	relPathLastPart := relPathParts[len(relPathParts)-1]      // file name with extension
	relPartsWithoutFile := relPathParts[:len(relPathParts)-1] // relative path to package

	pkgName := strings.Join(relPartsWithoutFile, "/")
	fileName := strings.TrimSuffix(relPathLastPart, ".neva")

	scope := src.Scope{
		Loc: src.ScopeLocation{
			ModuleName: "entry",
			PkgName:    pkgName,
			FileName:   fileName,
		},
		Module: s.state.mod,
	}

	pkg, ok := s.state.mod.Packages[pkgName]
	if !ok {
		return ResolveFileResponce{}, errors.New("no such package: " + pkgName)
	}

	file, ok := pkg[fileName]
	if !ok {
		return ResolveFileResponce{}, errors.New("no such file: " + fileName + "." + pkgName)
	}

	extra, err := getExtraForFile(file, scope)
	if err != nil {
		return ResolveFileResponce{}, err
	}

	return ResolveFileResponce{
		File: file,
		Extra: Extra{
			NodesPorts: extra,
		},
	}, nil
}

func getExtraForFile(file src.File, scope src.Scope) (map[string]map[string]src.Interface, error) {
	extra := map[string]map[string]src.Interface{}
	for entityName, entity := range file.Entities {
		if entity.Kind != src.ComponentEntity {
			continue
		}

		nodesIfaces := map[string]src.Interface{}
		for nodeName, node := range entity.Component.Nodes {
			nodeEntity, _, err := scope.Entity(node.EntityRef)
			if err != nil {
				return nil, err
			}

			var iface src.Interface
			if nodeEntity.Kind == src.ComponentEntity {
				iface = nodeEntity.Component.Interface
			} else if nodeEntity.Kind == src.InterfaceEntity {
				iface = nodeEntity.Interface
			}

			nodesIfaces[nodeName] = iface
		}

		extra[entityName] = nodesIfaces
	}

	return extra, nil
}

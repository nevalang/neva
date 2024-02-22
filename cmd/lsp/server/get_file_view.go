package server

import (
	"errors"
	"strings"

	src "github.com/nevalang/neva/pkg/sourcecode"
	"github.com/tliron/glsp"
)

type GetFileViewRequest struct {
	WorkspaceURI URI `json:"workspaceUri"`
	Document     struct {
		URI      URI    `json:"uri"`
		FileName string `json:"fileName"`
	} `json:"document"`
}

type URI struct {
	Path   string `json:"path"`
	FSPath string `json:"fsPath"`
}

type GetFileViewResponce struct {
	File  src.File `json:"file"`
	Extra Extra    `json:"extra"` // info that is not presented in the file but needed for rendering
}

type Extra struct {
	NodesPorts map[string]map[string]src.Interface `json:"nodesPorts"` // components -> nodes -> interface
}

func (s *Server) GetFileView(glspCtx *glsp.Context, req GetFileViewRequest) (GetFileViewResponce, error) {
	if s.index == nil {
		return GetFileViewResponce{}, nil
	}

	relFilePath := strings.TrimPrefix(req.Document.FileName, req.WorkspaceURI.Path)
	relFilePath = strings.TrimPrefix(relFilePath, "/")

	relPathParts := strings.Split(relFilePath, "/")           // relative path to file in slice
	relPathLastPart := relPathParts[len(relPathParts)-1]      // file name with extension
	relPartsWithoutFile := relPathParts[:len(relPathParts)-1] // relative path to package

	pkgName := strings.Join(relPartsWithoutFile, "/")
	fileName := strings.TrimSuffix(relPathLastPart, ".neva")

	scope := src.Scope{
		Location: src.Location{
			ModRef:   s.index.EntryModRef,
			PkgName:  pkgName,
			FileName: fileName,
		},
		Build: *s.index,
	}

	pkg, ok := s.index.Modules[s.index.EntryModRef].Packages[pkgName]
	if !ok {
		return GetFileViewResponce{}, errors.New("no such package: " + pkgName)
	}

	file, ok := pkg[fileName]
	if !ok {
		return GetFileViewResponce{}, errors.New("no such file: " + fileName + "." + pkgName)
	}

	extra, err := getExtraForFile(file, scope)
	if err != nil {
		return GetFileViewResponce{}, err
	}

	return GetFileViewResponce{
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

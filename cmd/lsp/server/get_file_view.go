package server

import (
	src "github.com/nevalang/neva/pkg/ast"
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
	NodesPorts map[string]map[string]src.Interface `json:"nodesPorts"` // flows -> nodes -> interface
}

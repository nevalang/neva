package server

import (
	"encoding/json"
	"fmt"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
)

type codeLensData struct {
	// URI points to the source file where the lens was created.
	URI string `json:"uri"`
	// Name is the entity name lens resolution should target.
	Name string `json:"name"`
	// Kind selects which query to run for the target entity.
	Kind codeLensKind `json:"kind"`
}

type codeLensKind string

const (
	codeLensKindReferences      codeLensKind = "references"
	codeLensKindImplementations codeLensKind = "implementations"
)

// TextDocumentCodeLens emits per-entity code lenses for references and implementations.
func (s *Server) TextDocumentCodeLens(
	glspCtx *glsp.Context,
	params *protocol.CodeLensParams,
) ([]protocol.CodeLens, error) {
	build, ok := s.getBuild()
	if !ok {
		return []protocol.CodeLens{}, nil
	}

	fileCtx, err := s.findFile(build, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	lenses := make([]protocol.CodeLens, 0, len(fileCtx.file.Entities)*2)
	missingMetaCount := 0
	for name, entity := range fileCtx.file.Entities {
		meta := entity.Meta()
		if meta == nil {
			missingMetaCount++
			continue
		}
		nameRange := rangeForName(*meta, name, 0)
		lenses = append(lenses, protocol.CodeLens{
			Range: nameRange,
			Data: codeLensData{
				URI:  pathToURI(fileCtx.filePath),
				Name: name,
				Kind: codeLensKindReferences,
			},
		})

		if entity.Kind == src.ComponentEntity || entity.Kind == src.InterfaceEntity {
			lenses = append(lenses, protocol.CodeLens{
				Range: nameRange,
				Data: codeLensData{
					URI:  pathToURI(fileCtx.filePath),
					Name: name,
					Kind: codeLensKindImplementations,
				},
			})
		}
	}
	if missingMetaCount > 0 {
		s.logger.Info("skipped code lenses for entities without metadata", "count", missingMetaCount, "file", fileCtx.filePath)
	}

	return lenses, nil
}

// CodeLensResolve computes locations and attaches the show-references command payload.
func (s *Server) CodeLensResolve(
	glspCtx *glsp.Context,
	lens *protocol.CodeLens,
) (*protocol.CodeLens, error) {
	parsedCodeLensData, ok := parseCodeLensData(lens.Data)
	if !ok {
		return lens, nil
	}

	build, ok := s.getBuild()
	if !ok {
		return lens, nil
	}

	fileCtx, err := s.findFile(build, parsedCodeLensData.URI)
	if err != nil {
		return nil, err
	}

	target, ok := s.resolveEntityRef(build, fileCtx, core.EntityRef{Name: parsedCodeLensData.Name})
	if !ok {
		return lens, nil
	}

	switch parsedCodeLensData.Kind {
	case codeLensKindImplementations:
		locations := s.appendDeclarationLocation(nil, target)
		title := fmt.Sprintf("%d implementations", len(locations))
		lens.Command = buildShowReferencesCommand(parsedCodeLensData.URI, lens.Range.Start, locations, title)
	case codeLensKindReferences:
		locations := s.referencesForEntity(build, target)
		title := fmt.Sprintf("%d references", len(locations))
		lens.Command = buildShowReferencesCommand(parsedCodeLensData.URI, lens.Range.Start, locations, title)
	default:
		return lens, nil
	}

	return lens, nil
}

// parseCodeLensData decodes strongly-typed lens metadata from LSP's generic data field.
func parseCodeLensData(raw any) (codeLensData, bool) {
	parsedCodeLensData := codeLensData{}
	rawJSON, err := json.Marshal(raw)
	if err != nil {
		return parsedCodeLensData, false
	}
	if err := json.Unmarshal(rawJSON, &parsedCodeLensData); err != nil {
		return parsedCodeLensData, false
	}
	if parsedCodeLensData.URI == "" || parsedCodeLensData.Name == "" {
		return parsedCodeLensData, false
	}
	switch parsedCodeLensData.Kind {
	case codeLensKindReferences, codeLensKindImplementations:
		return parsedCodeLensData, true
	default:
		return parsedCodeLensData, false
	}
}

// buildShowReferencesCommand creates the editor command expected by VS Code's references UI.
func buildShowReferencesCommand(
	uri string,
	position protocol.Position,
	locations []protocol.Location,
	title string,
) *protocol.Command {
	return &protocol.Command{
		Title:   title,
		Command: "editor.action.showReferences",
		Arguments: []any{
			uri,
			position,
			locations,
		},
	}
}

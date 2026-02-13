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
	URI  string `json:"uri"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// TextDocumentCodeLens emits per-entity code lenses for references and implementations.
func (s *Server) TextDocumentCodeLens(
	glspCtx *glsp.Context,
	params *protocol.CodeLensParams,
) ([]protocol.CodeLens, error) {
	build, ok := s.getBuild()
	if !ok {
		return []protocol.CodeLens{}, nil
	}

	ctx, err := s.findFile(build, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	lenses := make([]protocol.CodeLens, 0, len(ctx.file.Entities)*2)
	for name, entity := range ctx.file.Entities {
		meta := entity.Meta()
		if meta == nil {
			continue
		}
		rng := rangeForName(*meta, name, 0)
		lenses = append(lenses, protocol.CodeLens{
			Range: rng,
			Data: codeLensData{
				URI:  pathToURI(ctx.filePath),
				Name: name,
				Kind: "references",
			},
		})

		if entity.Kind == src.ComponentEntity || entity.Kind == src.InterfaceEntity {
			lenses = append(lenses, protocol.CodeLens{
				Range: rng,
				Data: codeLensData{
					URI:  pathToURI(ctx.filePath),
					Name: name,
					Kind: "implementations",
				},
			})
		}
	}

	return lenses, nil
}

// CodeLensResolve computes locations and attaches the show-references command payload.
func (s *Server) CodeLensResolve(
	glspCtx *glsp.Context,
	lens *protocol.CodeLens,
) (*protocol.CodeLens, error) {
	data, ok := parseCodeLensData(lens.Data)
	if !ok {
		return lens, nil
	}

	build, ok := s.getBuild()
	if !ok {
		return lens, nil
	}

	ctx, err := s.findFile(build, data.URI)
	if err != nil {
		return nil, err
	}

	target, ok := s.resolveEntityRef(build, ctx, core.EntityRef{Name: data.Name})
	if !ok {
		return lens, nil
	}

	switch data.Kind {
	case "implementations":
		locations := s.appendDeclarationLocation(nil, target)
		title := fmt.Sprintf("%d implementations", len(locations))
		lens.Command = buildShowReferencesCommand(data.URI, lens.Range.Start, locations, title)
	default:
		locations := s.referencesForEntity(build, target)
		title := fmt.Sprintf("%d references", len(locations))
		lens.Command = buildShowReferencesCommand(data.URI, lens.Range.Start, locations, title)
	}

	return lens, nil
}

// parseCodeLensData decodes strongly-typed lens metadata from LSP's generic data field.
func parseCodeLensData(raw any) (codeLensData, bool) {
	data := codeLensData{}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return data, false
	}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return data, false
	}
	if data.URI == "" || data.Name == "" {
		return data, false
	}
	return data, true
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

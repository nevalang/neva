// CodeLens handlers add actionable inline annotations above declarations.
// In Neva LSP we currently expose entity-level lenses for references and interface implementations,
// then resolve them into VS Code's show-references command payload.
package server

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
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

// TextDocumentCodeLens emits per-entity code lenses for references and, for interfaces, implementations.
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

		if entity.Kind == src.InterfaceEntity {
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
		locations := s.implementationLocationsForEntity(build, target)
		title := fmt.Sprintf("%d implementations", len(locations))
		lens.Command = buildShowReferencesCommand(parsedCodeLensData.URI, lens.Range.Start, locations, title)
	case codeLensKindReferences:
		locations := s.referenceLocationsForEntity(build, target)
		title := fmt.Sprintf("%d references", len(locations))
		lens.Command = buildShowReferencesCommand(parsedCodeLensData.URI, lens.Range.Start, locations, title)
	default:
		s.logger.Info("skipped code lens resolve for unknown kind", "kind", parsedCodeLensData.Kind, "name", parsedCodeLensData.Name)
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
	return parsedCodeLensData, isKnownCodeLensKind(parsedCodeLensData.Kind)
}

func isKnownCodeLensKind(kind codeLensKind) bool {
	return kind == codeLensKindReferences || kind == codeLensKindImplementations
}

// referenceLocationsForEntity includes explicit references and interface/component relation references.
func (s *Server) referenceLocationsForEntity(build *src.Build, target *resolvedEntity) []protocol.Location {
	locations := s.referencesForEntity(build, target)
	seen := make(map[protocol.Location]struct{}, len(locations))
	for _, location := range locations {
		seen[location] = struct{}{}
	}
	appendUnique := func(extra []protocol.Location) {
		for _, location := range extra {
			if _, ok := seen[location]; ok {
				continue
			}
			seen[location] = struct{}{}
			locations = append(locations, location)
		}
	}

	switch target.entity.Kind {
	case src.ComponentEntity:
		for _, implementedInterface := range s.implementedInterfacesForComponent(build, target) {
			appendUnique(s.referencesForEntity(build, implementedInterface))
		}
	case src.InterfaceEntity:
		appendUnique(s.implementationLocationsForEntity(build, target))
	case src.ConstEntity, src.TypeEntity:
		// No extra relationship references for const/type entities.
	}
	return locations
}

// implementationLocationsForEntity returns implementation-related locations for interfaces.
func (s *Server) implementationLocationsForEntity(build *src.Build, target *resolvedEntity) []protocol.Location {
	if target.entity.Kind != src.InterfaceEntity {
		return []protocol.Location{}
	}
	return s.implementationLocationsForInterface(build, target)
}

// implementationLocationsForInterface finds all components that structurally implement the interface.
func (s *Server) implementationLocationsForInterface(build *src.Build, ifaceTarget *resolvedEntity) []protocol.Location {
	interfaceDef := ifaceTarget.entity.Interface
	locations := []protocol.Location{}
	seen := map[protocol.Location]struct{}{}
	s.forEachWorkspaceFile(build, func(fileCtx fileContext) bool {
		for componentName, entity := range fileCtx.file.Entities {
			if entity.Kind != src.ComponentEntity {
				continue
			}
			for _, component := range entity.Component {
				if !componentImplementsInterface(component.IO, interfaceDef.IO) {
					continue
				}
				componentMeta := entity.Meta()
				if componentMeta == nil {
					continue
				}
				location := protocol.Location{
					URI:   pathToURI(fileCtx.filePath),
					Range: rangeForName(*componentMeta, componentName, 0),
				}
				if _, ok := seen[location]; ok {
					continue
				}
				seen[location] = struct{}{}
				locations = append(locations, location)
			}
		}
		return true
	})
	return locations
}

// implementedInterfacesForComponent returns interfaces that are structurally implemented by a component.
func (s *Server) implementedInterfacesForComponent(build *src.Build, componentTarget *resolvedEntity) []*resolvedEntity {
	if componentTarget.entity.Kind != src.ComponentEntity || len(componentTarget.entity.Component) == 0 {
		return nil
	}
	componentIO := componentTarget.entity.Component[0].IO
	implementedInterfaces := []*resolvedEntity{}
	s.forEachWorkspaceFile(build, func(fileCtx fileContext) bool {
		for interfaceName, entity := range fileCtx.file.Entities {
			if entity.Kind != src.InterfaceEntity {
				continue
			}
			if !componentImplementsInterface(componentIO, entity.Interface.IO) {
				continue
			}
			implementedInterfaces = append(implementedInterfaces, &resolvedEntity{
				moduleRef:   fileCtx.moduleRef,
				packageName: fileCtx.packageName,
				name:        interfaceName,
				filePath:    fileCtx.filePath,
				entity:      entity,
			})
		}
		return true
	})
	return implementedInterfaces
}

// componentImplementsInterface performs a structural port-level check for MVP interface implementation.
func componentImplementsInterface(componentIO src.IO, interfaceIO src.IO) bool {
	return ioContainsPorts(componentIO.In, interfaceIO.In) && ioContainsPorts(componentIO.Out, interfaceIO.Out)
}

func ioContainsPorts(componentPorts map[string]src.Port, interfacePorts map[string]src.Port) bool {
	interfacePortNames := make([]string, 0, len(interfacePorts))
	for name := range interfacePorts {
		interfacePortNames = append(interfacePortNames, name)
	}
	slices.Sort(interfacePortNames)
	for _, name := range interfacePortNames {
		interfacePort := interfacePorts[name]
		componentPort, ok := componentPorts[name]
		if !ok {
			return false
		}
		if componentPort.IsArray != interfacePort.IsArray {
			return false
		}
		if componentPort.TypeExpr.String() != interfacePort.TypeExpr.String() {
			return false
		}
	}
	return true
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

// Semantic tokens provide semantic highlighting categories beyond lexical tokenization.
// We emit declaration/reference tokens for Neva entities plus node/port address segments
// so editors can color symbols consistently across the current document.
package server

import (
	"math"
	"sort"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
)

// semanticTokenTypes returns the token names declared in the semantic token legend.
func semanticTokenTypes() []string {
	return []string{
		"namespace",
		"type",
		"class",
		"interface",
		"function",
		"variable",
		"property",
		"keyword",
		"constant",
	}
}

// semanticTokensLegend returns the static token schema advertised during initialize.
func semanticTokensLegend() protocol.SemanticTokensLegend {
	tokenTypes := semanticTokenTypes()
	return protocol.SemanticTokensLegend{
		TokenTypes:     tokenTypes,
		TokenModifiers: []string{},
	}
}

type semanticToken struct {
	line      int
	start     int
	length    int
	tokenType int
	modifiers int
}

// TextDocumentSemanticTokensFull returns full semantic tokens for a single Neva document.
func (s *Server) TextDocumentSemanticTokensFull(
	glspCtx *glsp.Context,
	params *protocol.SemanticTokensParams,
) (*protocol.SemanticTokens, error) {
	build, ok := s.getBuild()
	if !ok {
		return &protocol.SemanticTokens{Data: []uint32{}}, nil
	}

	fileCtx, err := s.findFile(build, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	tokens := s.collectSemanticTokens(build, fileCtx)
	encodedTokenData := encodeSemanticTokens(tokens)
	return &protocol.SemanticTokens{Data: encodedTokenData}, nil
}

// collectSemanticTokens gathers declaration, reference, and port-address tokens from a file.
func (s *Server) collectSemanticTokens(build *src.Build, fileCtx *fileContext) []semanticToken {
	typeIndex := tokenTypeIndex()
	tokens := make([]semanticToken, 0, len(fileCtx.file.Entities))

	for name, entity := range fileCtx.file.Entities {
		meta := entity.Meta()
		if meta == nil {
			continue
		}
		if tokenType, ok := entityTokenType(entity.Kind, typeIndex); ok {
			tokens = append(tokens, makeToken(*meta, 0, len(name), tokenType))
		}

		if entity.Kind == src.ComponentEntity {
			for _, comp := range entity.Component {
				for _, conn := range comp.Net {
					tokens = append(tokens, collectPortTokens(conn, typeIndex)...)
				}
			}
		}
	}

	refs := collectRefsInFile(fileCtx.file)
	for _, ref := range refs {
		resolved, ok := s.resolveEntityRef(build, fileCtx, ref.ref)
		if !ok {
			continue
		}
		tokenType, ok := entityTokenType(resolved.entity.Kind, typeIndex)
		if !ok {
			continue
		}
		offset := nameOffsetForRef(ref.meta)
		tokens = append(tokens, makeToken(ref.meta, offset, len(resolved.name), tokenType))
	}

	// LSP expects monotonically ordered tokens before delta encoding.
	sort.Slice(tokens, func(i, j int) bool {
		if tokens[i].line == tokens[j].line {
			return tokens[i].start < tokens[j].start
		}
		return tokens[i].line < tokens[j].line
	})

	return tokens
}

// tokenTypeIndex maps legend token names to numeric indices.
func tokenTypeIndex() map[string]int {
	tokenTypes := semanticTokenTypes()
	index := make(map[string]int, len(tokenTypes))
	for i, t := range tokenTypes {
		index[t] = i
	}
	return index
}

// entityTokenType maps Neva entity kinds to semantic token categories.
func entityTokenType(kind src.EntityKind, index map[string]int) (int, bool) {
	switch kind {
	case src.TypeEntity:
		return index["type"], true
	case src.InterfaceEntity:
		return index["interface"], true
	case src.ComponentEntity:
		return index["function"], true
	case src.ConstEntity:
		return index["constant"], true
	default:
		return 0, false
	}
}

// makeToken converts Neva metadata into absolute token coordinates.
func makeToken(meta core.Meta, offset int, length int, tokenType int) semanticToken {
	line := meta.Start.Line - 1
	start := meta.Start.Column + offset
	if line < 0 {
		line = 0
	}
	if start < 0 {
		start = 0
	}
	return semanticToken{
		line:      line,
		start:     start,
		length:    length,
		tokenType: tokenType,
		modifiers: 0,
	}
}

// collectPortTokens walks a connection tree and tokenizes node/port address segments.
func collectPortTokens(conn src.Connection, index map[string]int) []semanticToken {
	var tokens []semanticToken
	for _, sender := range conn.Senders {
		if sender.PortAddr != nil {
			tokens = append(tokens, portAddrTokens(*sender.PortAddr, index)...)
		}
	}
	for _, receiver := range conn.Receivers {
		if receiver.PortAddr != nil {
			tokens = append(tokens, portAddrTokens(*receiver.PortAddr, index)...)
		}
		if receiver.ChainedConnection != nil {
			tokens = append(tokens, collectPortTokens(*receiver.ChainedConnection, index)...)
		}
	}
	return tokens
}

// portAddrTokens emits separate tokens for node names and port names inside an address.
func portAddrTokens(addr src.PortAddr, index map[string]int) []semanticToken {
	if addr.Meta.Text == "" {
		return nil
	}
	text := addr.Meta.Text
	var tokens []semanticToken

	if addr.Node != "" {
		nodeLen := len(addr.Node)
		tokens = append(tokens, makeToken(addr.Meta, 0, nodeLen, index["variable"]))
	}

	if addr.Port != "" {
		idx := strings.Index(text, ":"+addr.Port)
		if idx >= 0 {
			portOffset := idx + 1
			tokens = append(tokens, makeToken(addr.Meta, portOffset, len(addr.Port), index["property"]))
		}
	}

	return tokens
}

// encodeSemanticTokens converts absolute tokens to LSP delta-encoded payload format.
func encodeSemanticTokens(tokens []semanticToken) []uint32 {
	encodedData := make([]uint32, 0, len(tokens)*5)
	lastLine := 0
	lastStart := 0

	for _, token := range tokens {
		deltaLine := token.line - lastLine
		deltaStart := token.start
		if deltaLine == 0 {
			deltaStart = token.start - lastStart
		}

		deltaLine32, ok := intToUint32(deltaLine)
		if !ok {
			continue
		}
		deltaStart32, ok := intToUint32(deltaStart)
		if !ok {
			continue
		}
		length32, ok := intToUint32(token.length)
		if !ok {
			continue
		}
		tokenType32, ok := intToUint32(token.tokenType)
		if !ok {
			continue
		}
		modifiers32, ok := intToUint32(token.modifiers)
		if !ok {
			continue
		}

		encodedData = append(encodedData, deltaLine32, deltaStart32, length32, tokenType32, modifiers32)

		lastLine = token.line
		lastStart = token.start
	}

	return encodedData
}

// intToUint32 safely converts a non-negative int that fits into uint32.
func intToUint32(value int) (uint32, bool) {
	if value < 0 || value > math.MaxUint32 {
		return 0, false
	}
	// #nosec G115 -- range is checked above before conversion.
	return uint32(value), true
}

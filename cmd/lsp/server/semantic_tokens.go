package server

import (
	"sort"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
)

var semanticTokenTypes = []string{
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

func semanticTokensLegend() protocol.SemanticTokensLegend {
	return protocol.SemanticTokensLegend{
		TokenTypes:     semanticTokenTypes,
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

func (s *Server) TextDocumentSemanticTokensFull(
	glspCtx *glsp.Context,
	params *protocol.SemanticTokensParams,
) (*protocol.SemanticTokens, error) {
	build, ok := s.getBuild()
	if !ok {
		return nil, nil
	}

	ctx, err := s.findFile(build, params.TextDocument.URI)
	if err != nil {
		return nil, nil
	}

	tokens := s.collectSemanticTokens(build, ctx)
	data := encodeSemanticTokens(tokens)
	return &protocol.SemanticTokens{Data: data}, nil
}

func (s *Server) collectSemanticTokens(build *src.Build, ctx *fileContext) []semanticToken {
	typeIndex := tokenTypeIndex()
	var tokens []semanticToken

	for name, entity := range ctx.file.Entities {
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

	refs := collectRefsInFile(ctx.file)
	for _, ref := range refs {
		resolved, ok := s.resolveEntityRef(build, ctx, ref.ref)
		if !ok {
			continue
		}
		tokenType, ok := entityTokenType(resolved.entity.Kind, typeIndex)
		if !ok {
			continue
		}
		offset := nameOffsetForRef(ref.meta, resolved.name)
		tokens = append(tokens, makeToken(ref.meta, offset, len(resolved.name), tokenType))
	}

	sort.Slice(tokens, func(i, j int) bool {
		if tokens[i].line == tokens[j].line {
			return tokens[i].start < tokens[j].start
		}
		return tokens[i].line < tokens[j].line
	})

	return tokens
}

func tokenTypeIndex() map[string]int {
	index := make(map[string]int, len(semanticTokenTypes))
	for i, t := range semanticTokenTypes {
		index[t] = i
	}
	return index
}

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

func collectPortTokens(conn src.Connection, index map[string]int) []semanticToken {
	var tokens []semanticToken
	if conn.Normal != nil {
		for _, sender := range conn.Normal.Senders {
			if sender.PortAddr != nil {
				tokens = append(tokens, portAddrTokens(*sender.PortAddr, index)...)
			}
		}
		for _, receiver := range conn.Normal.Receivers {
			if receiver.PortAddr != nil {
				tokens = append(tokens, portAddrTokens(*receiver.PortAddr, index)...)
			}
			if receiver.ChainedConnection != nil {
				tokens = append(tokens, collectPortTokens(*receiver.ChainedConnection, index)...)
			}
			if receiver.DeferredConnection != nil {
				tokens = append(tokens, collectPortTokens(*receiver.DeferredConnection, index)...)
			}
		}
	}
	if conn.ArrayBypass != nil {
		tokens = append(tokens, portAddrTokens(conn.ArrayBypass.SenderOutport, index)...)
		tokens = append(tokens, portAddrTokens(conn.ArrayBypass.ReceiverInport, index)...)
	}
	return tokens
}

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

func encodeSemanticTokens(tokens []semanticToken) []uint32 {
	var data []uint32
	lastLine := 0
	lastStart := 0

	for _, token := range tokens {
		deltaLine := token.line - lastLine
		deltaStart := token.start
		if deltaLine == 0 {
			deltaStart = token.start - lastStart
		}

		data = append(data,
			uint32(deltaLine),
			uint32(deltaStart),
			uint32(token.length),
			uint32(token.tokenType),
			uint32(token.modifiers),
		)

		lastLine = token.line
		lastStart = token.start
	}

	return data
}

package parser

import (
	"fmt"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
)

func (s *treeShapeListener) parseLeadingComments(
	entityLine int,
	io src.IO,
) (*src.Comments, *compiler.Error) {
	lines, startLine, stopLine := s.leadingCommentLines(entityLine)
	if len(lines) == 0 {
		return nil, nil
	}

	comments := &src.Comments{
		Inports:  map[string]string{},
		Outports: map[string]string{},
		Meta: core.Meta{
			Start: core.Position{Line: startLine, Column: 0},
			Stop:  core.Position{Line: stopLine, Column: 0},
			Location: core.Location{
				ModRef:   s.loc.ModRef,
				Package:  s.loc.Package,
				Filename: s.loc.Filename,
			},
		},
	}

	textBuf := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			s.flushTextBlock(comments, textBuf)
			textBuf = textBuf[:0]
			continue
		}

		if !strings.HasPrefix(trimmed, "@") {
			textBuf = append(textBuf, line)
			continue
		}

		s.flushTextBlock(comments, textBuf)
		textBuf = textBuf[:0]

		tagName, tagValue := splitTag(trimmed[1:])
		switch tagName {
		case "inport":
			portName, desc := splitFirstWord(tagValue)
			if portName == "" {
				return nil, s.commentParseError("invalid @inport tag: port name is required", startLine)
			}
			if _, ok := io.In[portName]; !ok {
				return nil, s.commentParseError(fmt.Sprintf("unknown @inport reference: %s", portName), startLine)
			}
			comments.Inports[portName] = desc
		case "outport":
			portName, desc := splitFirstWord(tagValue)
			if portName == "" {
				return nil, s.commentParseError("invalid @outport tag: port name is required", startLine)
			}
			if _, ok := io.Out[portName]; !ok {
				return nil, s.commentParseError(fmt.Sprintf("unknown @outport reference: %s", portName), startLine)
			}
			comments.Outports[portName] = desc
		case "example":
			comments.Examples = append(comments.Examples, tagValue)
		default:
			// Parser intentionally keeps unknown tags for CLI/AI and third-party tooling.
			comments.Tags = append(comments.Tags, src.CommentTag{
				Name:  tagName,
				Value: tagValue,
			})
		}
	}
	s.flushTextBlock(comments, textBuf)

	if len(comments.Inports) == 0 {
		comments.Inports = nil
	}
	if len(comments.Outports) == 0 {
		comments.Outports = nil
	}
	if len(comments.TextBlocks) == 0 && len(comments.Inports) == 0 &&
		len(comments.Outports) == 0 && len(comments.Examples) == 0 && len(comments.Tags) == 0 {
		return nil, nil
	}

	return comments, nil
}

func (s *treeShapeListener) flushTextBlock(comments *src.Comments, textLines []string) {
	if len(textLines) == 0 {
		return
	}
	comments.TextBlocks = append(comments.TextBlocks, strings.Join(textLines, "\n"))
}

func (s *treeShapeListener) leadingCommentLines(entityLine int) ([]string, int, int) {
	if entityLine <= 1 || entityLine-2 >= len(s.sourceLines) {
		return nil, 0, 0
	}

	i := entityLine - 2 // 0-based line before entity declaration.
	line := strings.TrimSpace(s.sourceLines[i])
	if !strings.HasPrefix(line, "//") {
		return nil, 0, 0
	}

	raw := make([]string, 0, 8)
	start := i + 1
	stop := i + 1
	for ; i >= 0; i-- {
		trimmed := strings.TrimSpace(s.sourceLines[i])
		if !strings.HasPrefix(trimmed, "//") {
			break
		}
		content := strings.TrimPrefix(trimmed, "//")
		content = strings.TrimPrefix(content, " ")
		raw = append(raw, content)
		start = i + 1
	}

	// reverse to keep source order.
	for left, right := 0, len(raw)-1; left < right; left, right = left+1, right-1 {
		raw[left], raw[right] = raw[right], raw[left]
	}

	return raw, start, stop
}

func splitTag(s string) (string, string) {
	name, value := splitFirstWord(s)
	return strings.ToLower(name), value
}

func splitFirstWord(s string) (string, string) {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return "", ""
	}
	name := parts[0]
	idx := strings.Index(s, name)
	rest := strings.TrimSpace(s[idx+len(name):])
	return name, rest
}

func (s *treeShapeListener) commentParseError(message string, line int) *compiler.Error {
	return &compiler.Error{
		Message: message,
		Meta: &core.Meta{
			Start: core.Position{
				Line:   line,
				Column: 0,
			},
			Location: core.Location{
				ModRef:   s.loc.ModRef,
				Package:  s.loc.Package,
				Filename: s.loc.Filename,
			},
		},
	}
}

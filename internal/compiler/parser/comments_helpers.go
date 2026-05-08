package parser

import (
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
)

//nolint:gocyclo,cyclop // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (s *treeShapeListener) parseLeadingComments(
	entityLine int,
	ports *src.IO,
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
		if err := s.consumeTag(comments, ports, startLine, tagName, tagValue); err != nil {
			return nil, err
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

func (s *treeShapeListener) consumeTag(
	comments *src.Comments,
	ports *src.IO,
	startLine int,
	tagName string,
	tagValue string,
) *compiler.Error {
	switch tagName {
	case "inport":
		portName, desc := splitFirstWord(tagValue)
		if portName == "" {
			return s.commentParseError("invalid @inport tag: port name is required", startLine)
		}
		if _, ok := ports.In[portName]; !ok {
			return s.commentParseError("unknown @inport reference: "+portName, startLine)
		}
		comments.Inports[portName] = desc
	case "outport":
		portName, desc := splitFirstWord(tagValue)
		if portName == "" {
			return s.commentParseError("invalid @outport tag: port name is required", startLine)
		}
		if _, ok := ports.Out[portName]; !ok {
			return s.commentParseError("unknown @outport reference: "+portName, startLine)
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

	return nil
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

	lineIdx := entityLine - 2 // 0-based line before entity declaration.
	line := strings.TrimSpace(s.sourceLines[lineIdx])
	if !strings.HasPrefix(line, "//") {
		return nil, 0, 0
	}

	raw := make([]string, 0, 8)
	start := lineIdx + 1
	stop := lineIdx + 1
	for ; lineIdx >= 0; lineIdx-- {
		trimmed := strings.TrimSpace(s.sourceLines[lineIdx])
		if !strings.HasPrefix(trimmed, "//") {
			break
		}
		content := strings.TrimPrefix(trimmed, "//")
		content = strings.TrimPrefix(content, " ")
		raw = append(raw, content)
		start = lineIdx + 1
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

func splitFirstWord(text string) (string, string) {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return "", ""
	}
	name := parts[0]
	idx := strings.Index(text, name)
	rest := strings.TrimSpace(text[idx+len(name):])
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

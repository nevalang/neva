package funcs

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/nevalang/neva/internal/runtime"
)

type dotenvLoadFrom struct {
	override bool
}

type dotenvLoad struct {
	override bool
}

func (d dotenvLoadFrom) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	filenameIn, err := rio.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			filenameMsg, ok := filenameIn.Receive(ctx)
			if !ok {
				return
			}

			filename := strings.TrimSpace(filenameMsg.Str())
			if filename == "" {
				filename = ".env"
			}

			if err := loadDotenvFile(filename, d.override); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, emptyStruct()) {
				return
			}
		}
	}, nil
}

func (d dotenvLoad) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			if err := loadDotenvFile(".env", d.override); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, emptyStruct()) {
				return
			}
		}
	}, nil
}

func loadDotenvFile(path string, override bool) error {
	// #nosec G304,G703 -- path comes from explicit component input (caller-controlled by design).
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	parsedDotenvValues, err := parseDotenv(file)
	if err != nil {
		return err
	}

	for key, value := range parsedDotenvValues {
		if !override {
			if _, exists := os.LookupEnv(key); exists {
				continue
			}
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set %q: %w", key, err)
		}
	}

	return nil
}

func parseDotenv(r io.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 4096), 1024*1024)

	result := make(map[string]string)
	lineNum := 0

	for scanner.Scan() {
		lineNum++

		raw := scanner.Text()
		if lineNum == 1 {
			raw = strings.TrimPrefix(raw, "\ufeff")
		}

		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if after, ok := strings.CutPrefix(line, "export "); ok {
			line = strings.TrimSpace(after)
		}

		key, value, err := parseDotenvEntry(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func parseDotenvEntry(line string) (string, string, error) {
	idx := strings.IndexRune(line, '=')
	if idx == -1 {
		return "", "", errors.New("missing '='")
	}

	key := strings.TrimSpace(line[:idx])
	if key == "" {
		return "", "", errors.New("missing key")
	}

	valuePart := stripInlineComment(line[idx+1:])
	value, err := parseDotenvValue(valuePart)
	if err != nil {
		return "", "", err
	}

	return key, value, nil
}

func parseDotenvValue(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}

	if raw[0] == '"' {
		if len(raw) == 1 || raw[len(raw)-1] != '"' {
			return "", errors.New("unterminated double-quoted value")
		}
		return unescapeDoubleQuoted(raw[1 : len(raw)-1])
	}

	if raw[0] == '\'' {
		if len(raw) == 1 || raw[len(raw)-1] != '\'' {
			return "", errors.New("unterminated single-quoted value")
		}
		inner := raw[1 : len(raw)-1]
		return strings.ReplaceAll(inner, "\\'", "'"), nil
	}

	return raw, nil
}

func unescapeDoubleQuoted(s string) (string, error) {
	var builder strings.Builder
	builder.Grow(len(s))

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch != '\\' {
			builder.WriteByte(ch)
			continue
		}

		i++
		if i >= len(s) {
			return "", errors.New("unterminated escape sequence")
		}

		switch s[i] {
		case 'n':
			builder.WriteByte('\n')
		case 'r':
			builder.WriteByte('\r')
		case 't':
			builder.WriteByte('\t')
		case 'b':
			builder.WriteByte('\b')
		case 'f':
			builder.WriteByte('\f')
		case 'v':
			builder.WriteByte('\v')
		case '\\':
			builder.WriteByte('\\')
		case '"':
			builder.WriteByte('"')
		case '$':
			builder.WriteByte('$')
		default:
			builder.WriteByte(s[i])
		}
	}

	return builder.String(), nil
}

func stripInlineComment(value string) string {
	value = strings.TrimRightFunc(value, unicode.IsSpace)

	inSingle := false
	inDouble := false

	for i := 0; i < len(value); i++ {
		switch value[i] {
		case '\'':
			if inDouble || isEscaped(value, i) {
				continue
			}
			inSingle = !inSingle
		case '"':
			if inSingle || isEscaped(value, i) {
				continue
			}
			inDouble = !inDouble
		case '#':
			if inSingle || inDouble {
				continue
			}
			if i == 0 || unicode.IsSpace(rune(value[i-1])) {
				trimmed := strings.TrimSpace(value[:i])
				if trimmed == "" {
					return ""
				}
				return trimmed
			}
		}
	}

	return strings.TrimSpace(value)
}

func isEscaped(s string, idx int) bool {
	if idx == 0 {
		return false
	}

	backslashes := 0
	for i := idx - 1; i >= 0 && s[i] == '\\'; i-- {
		backslashes++
	}

	return backslashes%2 == 1
}

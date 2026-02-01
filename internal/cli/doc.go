package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"

	cli "github.com/urfave/cli/v2"

	"github.com/nevalang/neva/std"
)

func newDocCmd() *cli.Command {
	return &cli.Command{
		Name:      "doc",
		Usage:     "Search the standard library for symbols using regular expressions",
		ArgsUsage: "[package/path] [.] <pattern>",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "context",
				Aliases: []string{"C"},
				Value:   2,
				Usage:   "number of context lines to display around each match",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			pkgPath, pattern, err := parseDocArgs(cliCtx.Args())
			if err != nil {
				return err
			}

			re, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("invalid pattern: %w", err)
			}

			matches, err := searchStdlib(pkgPath, re, cliCtx.Int("context"))
			if err != nil {
				return err
			}

			if len(matches) == 0 {
				fmt.Println("No matches found")
				return nil
			}

			sort.Slice(matches, func(i, j int) bool {
				if matches[i].file == matches[j].file {
					return matches[i].line < matches[j].line
				}
				return matches[i].file < matches[j].file
			})

			for _, m := range matches {
				fmt.Printf("std/%s:%d\n", m.file, m.line)
				for _, ctxLine := range m.context {
					marker := " "
					if ctxLine.line == m.line {
						marker = ">"
					}
					fmt.Printf("%s %6d | %s\n", marker, ctxLine.line, ctxLine.text)
				}
				fmt.Println()
			}

			return nil
		},
	}
}

//nolint:govet // fieldalignment: small helper struct.
type docMatch struct {
	file    string
	line    int
	context []docContextLine
}

//nolint:govet // fieldalignment: small helper struct.
type docContextLine struct {
	line int
	text string
}

func parseDocArgs(args cli.Args) (string, string, error) {
	switch args.Len() {
	case 0:
		return "", "", errors.New("expected at least 1 argument")
	case 1:
		return "", args.Get(0), nil
	case 2:
		pkgPath := normalizePkgArg(args.Get(0))
		if pkgPath == "" {
			return "", args.Get(1), nil
		}
		return pkgPath, args.Get(1), nil
	case 3:
		if args.Get(1) != "." {
			return "", "", errors.New("second argument must be '.' when three arguments are provided")
		}
		pkgPath := normalizePkgArg(args.Get(0))
		return pkgPath, args.Get(2), nil
	default:
		return "", "", errors.New("expected at most 3 arguments")
	}
}

func normalizePkgArg(arg string) string {
	trimmed := strings.TrimSpace(arg)
	if trimmed == "." {
		return ""
	}
	trimmed = strings.Trim(trimmed, "/")
	return strings.ReplaceAll(trimmed, "\\", "/")
}

func searchStdlib(pkgPath string, re *regexp.Regexp, context int) ([]docMatch, error) {
	if context < 0 {
		context = 0
	}

	var fsys fs.FS = std.FS
	basePath := ""

	if pkgPath != "" {
		sub, err := fs.Sub(std.FS, pkgPath)
		if err != nil {
			return nil, fmt.Errorf("package %q not found in stdlib: %w", pkgPath, err)
		}
		fsys = sub
		basePath = pkgPath
	}

	matches := make([]docMatch, 0)
	err := fs.WalkDir(fsys, ".", func(entryPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(entryPath, ".neva") {
			return nil
		}

		data, err := fs.ReadFile(fsys, entryPath)
		if err != nil {
			return err
		}

		lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
		for i, line := range lines {
			if !re.MatchString(line) {
				continue
			}

			start := i - context
			if start < 0 {
				start = 0
			}
			end := i + context
			if end >= len(lines) {
				end = len(lines) - 1
			}

			contextLines := make([]docContextLine, 0, end-start+1)
			for idx := start; idx <= end; idx++ {
				contextLines = append(contextLines, docContextLine{
					line: idx + 1,
					text: lines[idx],
				})
			}

			relativePath := entryPath
			if basePath != "" {
				relativePath = path.Join(basePath, entryPath)
			}

			matches = append(matches, docMatch{
				file:    relativePath,
				line:    i + 1,
				context: contextLines,
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return matches, nil
}

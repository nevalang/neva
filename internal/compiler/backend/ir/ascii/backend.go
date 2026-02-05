package ascii

import (
	"fmt"
	"io"
	"strings"

	ma "github.com/AlexanderGrooff/mermaid-ascii/cmd"

	"github.com/nevalang/neva/internal/compiler/backend/ir/mermaid"
	"github.com/nevalang/neva/internal/compiler/backend/ir/report"
	"github.com/nevalang/neva/internal/compiler/ir"
)

type Encoder struct{}

func (e Encoder) Encode(w io.Writer, prog *ir.Program) error {
	rep := report.Build(prog)
	if err := report.WriteHeader(w, rep); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "## 1. Visual Flow"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, ""); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "```"); err != nil {
		return err
	}

	flow, err := mermaid.EncodeFlowchart(prog)
	if err != nil {
		return err
	}

	rendered, err := renderASCII(flow)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, rendered); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "```"); err != nil {
		return err
	}

	if err := report.WriteComponentsTable(w, rep); err != nil {
		return err
	}
	return report.WriteMetrics(w, rep)
}

func renderASCII(flowchart string) (string, error) {
	diagramText := strings.Join([]string{
		"flowchart TD",
		strings.TrimPrefix(flowchart, "flowchart TD"),
	}, "\n")

	out, err := ma.RenderDiagram(diagramText, nil)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(toASCII(out), "\n"), nil
}

func toASCII(input string) string {
	replacer := strings.NewReplacer(
		"─", "-",
		"━", "-",
		"│", "|",
		"┃", "|",
		"┌", "+",
		"┐", "+",
		"└", "+",
		"┘", "+",
		"├", "+",
		"┤", "+",
		"┬", "+",
		"┴", "+",
		"┼", "+",
		"╭", "+",
		"╮", "+",
		"╰", "+",
		"╯", "+",
		"╱", "/",
		"╲", "\\",
		"→", ">",
		"←", "<",
		"↑", "^",
		"↓", "v",
	)
	return replacer.Replace(input)
}

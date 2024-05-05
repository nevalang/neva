package funcs

import (
	"context"
	"fmt"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type trace struct{}

func (t trace) Create(
	io runtime.FuncIO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	msgIn, err := io.In.Port("msg")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var msg runtime.Msg
			select {
			case <-ctx.Done():
				return
			case msg = <-msgIn:
			}

			s := formatTrace(msg.Trace())
			select {
			case <-ctx.Done():
				return
			case dataOut <- runtime.NewStrMsg(s):
			}
		}
	}, nil
}

func formatTrace(trace []runtime.PortAddr) string {
	var b strings.Builder
	depths := make([]int, len(trace)) // to store depths of all items

	for i, portAddr := range trace {
		depths[i] = strings.Count(portAddr.Path, "/")
	}

	for i, portAddr := range trace {
		currentDepth := depths[i]

		// Indent according to depth, managing spaces and vertical bars accurately
		if i > 0 {
			for j := 0; j < currentDepth; j++ {
				if shouldDisplayVerticalBar(depths, j, i) {
					b.WriteString("│  ")
				} else {
					b.WriteString("   ")
					// b.WriteString("└─ ")
				}
			}
		}

		// Determine if it's the last element at its level
		if i == len(trace)-1 || (i < len(trace)-1 && depths[i+1] <= currentDepth) {
			b.WriteString("└─ ")
		} else {
			b.WriteString("├─ ")
		}

		// Append the port address
		b.WriteString(fmt.Sprintf("%v:%v[%v]\n", portAddr.Path, portAddr.Port, portAddr.Idx))
	}

	return b.String()
}

// Helper function to determine if a vertical bar should be displayed at the given depth
func shouldDisplayVerticalBar(depths []int, depth int, index int) bool {
	for k := index + 1; k < len(depths); k++ {
		if depths[k] > depth {
			return true // Keep the vertical bar as there are deeper entries
		}
		if depths[k] < depth {
			return false // Stop the vertical bar as it returns to shallower entries
		}
		// Continue the bar if same depth is found
		if depths[k] == depth {
			return true
		}
	}
	return false
}

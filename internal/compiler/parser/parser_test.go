// We need unit tests for parser because it contains not only antlr grammar but also mapping logic.

package parser_test

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/stretchr/testify/require"
)

// TestParser_ParseFile_Directives checks only
// how compiler directives are parsed.
func TestParser_ParseFile_Directives(t *testing.T) {
	text := []byte(`
		component {
			#extern(d1)
			C1() ()

			#extern(d2)
			C2() () {
				nodes {
					#bind(d3)
					n1 C1

					#bind(d4)
					n2 C1
				}
			}

			#autoports
			C3() ()

			#extern(d5)
			#autoports
			C4() ()
		}
	`)

	p := parser.New(false)

	got, err := p.ParseFile(text)
	require.True(t, err == nil)

	d1 := got.Entities["C1"].Component.Directives[compiler.ExternDirective][0]
	require.Equal(t, "d1", d1)

	c2 := got.Entities["C2"].Component

	d2 := c2.Directives[compiler.ExternDirective][0]
	require.Equal(t, "d2", d2)

	d3 := c2.Nodes["n1"].Directives[compiler.BindDirective][0]
	require.Equal(t, "d3", d3)

	d4 := c2.Nodes["n2"].Directives[compiler.BindDirective][0]
	require.Equal(t, "d4", d4)

	c3 := got.Entities["C3"].Component
	_, ok := c3.Directives[compiler.AutoportsDirective]
	require.Equal(t, true, ok)

	c4 := got.Entities["C4"].Component
	d5, ok := c4.Directives[compiler.ExternDirective]
	require.Equal(t, true, ok)
	require.Equal(t, "d5", d5[0])
	_, ok = c4.Directives[compiler.AutoportsDirective]
	require.Equal(t, true, ok)
}

// Check that parser correctly parses port addresses without
// explicitly specified nodes.
func TestParser_ParseFile_IONodes(t *testing.T) {
	text := []byte(`
		component C1(start any) (stop any) {
			net { :start -> :stop }
		}
	`)

	p := parser.New(false)

	got, err := p.ParseFile(text)
	require.True(t, err == nil)

	conn := got.Entities["C1"].Component.Net[0]
	sender := conn.Normal.SenderSide.PortAddr.Node
	receiver := conn.Normal.ReceiverSide.Receivers[0].PortAddr.Node
	require.Equal(t, "in", sender)
	require.Equal(t, "out", receiver)
}

// We need unit tests for parser because it contains not only antlr grammar but also mapping logic.

package parser_test

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/stretchr/testify/require"
)

// TestParser_ParseFile_Directives checks only how compiler directives are parsed.
func TestParser_ParseFile_Directives(t *testing.T) {
	text := []byte(`
		components {
			#runtime_func(d1)
			C1() ()

			#runtime_func(d2)
			C2() () {
				nodes {
					#runtime_func_msg(d3)
					n1 C1

					#runtime_func_msg(d4)
					n2 C1
				}
			}

			#struct_inports
			C3() ()

			#runtime_func(d5)
			#struct_inports
			C4() ()
		}
	`)

	p := parser.New(false)

	got, err := p.ParseFile(text)
	require.True(t, err == nil)

	d1 := got.Entities["C1"].Component.Directives[compiler.RuntimeFuncDirective][0]
	require.Equal(t, "d1", d1)

	c2 := got.Entities["C2"].Component

	d2 := c2.Directives[compiler.RuntimeFuncDirective][0]
	require.Equal(t, "d2", d2)

	d3 := c2.Nodes["n1"].Directives[compiler.RuntimeFuncMsgDirective][0]
	require.Equal(t, "d3", d3)

	d4 := c2.Nodes["n2"].Directives[compiler.RuntimeFuncMsgDirective][0]
	require.Equal(t, "d4", d4)

	c3 := got.Entities["C3"].Component
	_, ok := c3.Directives[compiler.StructInports]
	require.Equal(t, true, ok)

	c4 := got.Entities["C4"].Component
	d5, ok := c4.Directives[compiler.RuntimeFuncDirective]
	require.Equal(t, true, ok)
	require.Equal(t, "d5", d5[0])
	_, ok = c4.Directives[compiler.StructInports]
	require.Equal(t, true, ok)
}

package parser

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	"github.com/stretchr/testify/require"
)

var location = core.Location{
	ModRef: core.ModuleRef{
		Path:    "test",
		Version: "1",
	},
	Package:  "test",
	Filename: "test",
}

func TestParser_ParseFile_StructSelectorsWithLonelyChain(t *testing.T) {
	text := []byte(`
		def C1() () {
			userSender -> .pet.name -> println -> :stop
		}`,
	)
	p := New()
	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component[0].Net
	require.Equal(t, 1, len(net))

	conn := net[0].Normal
	require.Equal(t, "userSender", conn.Senders[0].PortAddr.Node)
	require.Equal(t, "", conn.Senders[0].PortAddr.Port)

	chain := conn.Receivers[0].ChainedConnection.Normal
	require.Equal(t, "pet", chain.Senders[0].StructSelector[0])
	require.Equal(t, "name", chain.Senders[0].StructSelector[1])

	secondChain := chain.Receivers[0].ChainedConnection.Normal
	require.Equal(t, "println", secondChain.Senders[0].PortAddr.Node)
	require.Equal(t, "", secondChain.Senders[0].PortAddr.Port)

	chainEnd := secondChain.Receivers[0].PortAddr
	require.Equal(t, "stop", chainEnd.Port)
}

func TestParser_ParseFile_PortlessArrPortAddr(t *testing.T) {
	text := []byte(`
		def C1() () {
			foo[0] -> bar[255]
		}
	`)

	p := New()

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.Equal(t, true, err == nil)

	net := got.Entities["C1"].Component[0].Net
	conn := net[0].Normal

	// foo[0]->
	require.Equal(t, "foo", conn.Senders[0].PortAddr.Node)
	require.Equal(t, "", conn.Senders[0].PortAddr.Port)
	require.Equal(t, compiler.Pointer(uint8(0)), conn.Senders[0].PortAddr.Idx)

	// ->bar[255]
	require.Equal(t, "bar", conn.Receivers[0].PortAddr.Node)
	require.Equal(t, "", conn.Receivers[0].PortAddr.Port)
	require.Equal(t, compiler.Pointer(uint8(255)), conn.Receivers[0].PortAddr.Idx)
}

func TestParser_ParseFile_LonelyPorts(t *testing.T) {
	text := []byte(`
		def C1() () {
			:port -> lonely
			lonely -> :port
		}
	`)

	p := New()

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.True(t, err == nil)

	// 1) :port -> lonely
	// 2) lonely -> :port
	net := got.Entities["C1"].Component[0].Net
	require.Equal(t, 2, len(net))

	// 1) :port -> lonely
	receiverPortAddr := net[0].Normal.Receivers[0].PortAddr
	require.Equal(t, "lonely", receiverPortAddr.Node)
	require.Equal(t, "", receiverPortAddr.Port)

	// 2) lonely -> :port
	senderPortAddr := net[1].Normal.Senders[0].PortAddr
	require.Equal(t, "lonely", senderPortAddr.Node)
	require.Equal(t, "", senderPortAddr.Port)
}

func TestParser_ParseFile_ChainedConnections(t *testing.T) {
	text := []byte(`
		def C1() () { :foo -> n1:p1 -> :bar }
	`)

	p := New()

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component[0].Net
	require.Equal(t, 1, len(net))
	conn := net[0].Normal

	sender := conn.Senders[0].PortAddr
	require.Equal(t, "in", sender.Node)
	require.Equal(t, "foo", sender.Port)

	chain := conn.Receivers[0].ChainedConnection.Normal
	chainSender := chain.Senders[0].PortAddr
	require.Equal(t, "n1", chainSender.Node)
	require.Equal(t, "p1", chainSender.Port)

	chainReceiver := chain.Receivers[0].PortAddr
	require.Equal(t, "out", chainReceiver.Node)
	require.Equal(t, "bar", chainReceiver.Port)

	require.Greater(t, chain.Meta.Start.Line, 0)
	require.Greater(t, chain.Meta.Stop.Line, 0)
}

func TestParser_ParseFile_ChainedConnectionsWithConstants(t *testing.T) {
	tests := []struct { //nolint:govet // fieldalignment
		name  string
		text  string
		check func(t *testing.T, net []src.Connection)
	}{
		{
			name: "const ref in chain",
			text: `
				const greeting string = 'hello'
				def C1() () {
					:start -> $greeting -> :stop
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				require.Equal(t, "in", conn.Senders[0].PortAddr.Node)
				require.Equal(t, "start", conn.Senders[0].PortAddr.Port)

				chain := conn.Receivers[0].ChainedConnection.Normal
				require.NotNil(t, chain.Senders[0].Const)
				require.Equal(t, "greeting", chain.Senders[0].Const.Value.Ref.Name)
				require.Equal(t, "out", chain.Receivers[0].PortAddr.Node)
				require.Equal(t, "stop", chain.Receivers[0].PortAddr.Port)
			},
		},
		{
			name: "message literal in chain",
			text: `
				def C1() () {
					:start -> 'hello' -> :stop
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				require.Equal(t, "in", conn.Senders[0].PortAddr.Node)
				require.Equal(t, "start", conn.Senders[0].PortAddr.Port)

				chain := conn.Receivers[0].ChainedConnection.Normal
				require.NotNil(t, chain.Senders[0].Const)
				require.Equal(t, "hello", *chain.Senders[0].Const.Value.Message.Str)
				require.Equal(t, "out", chain.Receivers[0].PortAddr.Node)
				require.Equal(t, "stop", chain.Receivers[0].PortAddr.Port)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()

			got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(tt.text))
			require.Nil(t, err)

			net := got.Entities["C1"].Component[0].Net
			tt.check(t, net)
		})
	}
}

func TestParser_ParseFile_Comments(t *testing.T) {
	text := []byte(`
	// comment
	`)

	p := New()

	_, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.True(t, err == nil)
}

func TestParser_ParseFile_Directives(t *testing.T) {
	text := []byte(`
		#extern(d1)
		def C1() ()

		#extern(d2)
		def C2() () {
			#bind(d3)
			n1 C1

			#bind(d4)
			n2 C1
			---
		}

		#autoports
		def C3() ()

		#extern(d5)
		#autoports
		def C4() ()
	`)

	p := New()

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.True(t, err == nil)

	d1 := got.Entities["C1"].Component[0].Directives[compiler.ExternDirective]
	require.Equal(t, "d1", d1)

	c2 := got.Entities["C2"].Component

	d2 := c2[0].Directives[compiler.ExternDirective]
	require.Equal(t, "d2", d2)

	d3 := c2[0].Nodes["n1"].Directives[compiler.BindDirective]
	require.Equal(t, "d3", d3)

	d4 := c2[0].Nodes["n2"].Directives[compiler.BindDirective]
	require.Equal(t, "d4", d4)

	c3 := got.Entities["C3"].Component
	_, ok := c3[0].Directives[compiler.AutoportsDirective]
	require.Equal(t, true, ok)

	c4 := got.Entities["C4"].Component
	d5, ok := c4[0].Directives[compiler.ExternDirective]
	require.Equal(t, true, ok)
	require.Equal(t, "d5", d5)
	_, ok = c4[0].Directives[compiler.AutoportsDirective]
	require.Equal(t, true, ok)
}

func TestParser_ParseFile_IONodes(t *testing.T) {
	text := []byte(`
		def C1(start any) (stop any) {
			:start -> :stop
		}
	`)

	p := New()

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.True(t, err == nil)

	conn := got.Entities["C1"].Component[0].Net[0]

	sender := conn.Normal.Senders[0].PortAddr.Node
	require.Equal(t, "in", sender)

	receiver := conn.Normal.Receivers[0].PortAddr.Node
	require.Equal(t, "out", receiver)
}

func TestParser_ParseFile_AnonymousNodes(t *testing.T) {
	text := []byte(`
		def C1(start any) (stop any) {
			Scanner
			Printer<int>
			---
		}
	`)

	p := New()

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.True(t, err == nil)

	nodes := got.Entities["C1"].Component[0].Nodes

	_, ok := nodes["scanner"]
	require.Equal(t, true, ok)

	_, ok = nodes["printer"]
	require.Equal(t, true, ok)
}

func TestParser_ParseFile_TaggedUnionTypeExpr(t *testing.T) {
	tests := []struct { //nolint:govet // fieldalignment
		name  string
		text  string
		check func(t *testing.T, got src.File)
	}{
		{
			name: "simple union",
			text: `
				type Input union {
					Int int
					None
				}
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				unionType := got.Entities["Input"].Type.BodyExpr.Lit.Union
				require.NotNil(t, unionType)

				// Check Int tag with value
				intTag := unionType["Int"]
				require.NotNil(t, intTag)
				require.Equal(t, "int", intTag.Inst.Ref.Name)

				// Check None tag without value
				noneTag, ok := unionType["None"]
				require.True(t, ok)
				require.Nil(t, noneTag)
			},
		},
		{
			name: "one-line union",
			text: `
				type Result union { Ok int, Err string }
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				unionType := got.Entities["Result"].Type.BodyExpr.Lit.Union
				require.NotNil(t, unionType)

				// Check Ok tag with int value
				okTag := unionType["Ok"]
				require.NotNil(t, okTag)
				require.Equal(t, "int", okTag.Inst.Ref.Name)

				// Check Err tag with string value
				errTag := unionType["Err"]
				require.NotNil(t, errTag)
				require.Equal(t, "string", errTag.Inst.Ref.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(tt.text))
			require.Nil(t, err)
			tt.check(t, got)
		})
	}
}

func TestParser_ParseFile_TaggedUnionConstLiteral(t *testing.T) {
	tests := []struct { //nolint:govet // fieldalignment
		name  string
		text  string
		check func(t *testing.T, got src.File)
	}{
		{
			name: "union with value",
			text: `
				const c0 Input = Input::Int(42)
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c0"].Const.Value.Message.Union
				require.Equal(t, "", union.EntityRef.Pkg)
				require.Equal(t, "Input", union.EntityRef.Name)
				require.Equal(t, "Int", union.Tag)
				require.Equal(t, 42, *union.Data.Message.Int)
			},
		},
		{
			name: "union without value",
			text: `
				const c1 pkg.Input = pkg.Input::None
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c1"].Const.Value.Message.Union
				require.Equal(t, "pkg", union.EntityRef.Pkg)
				require.Equal(t, "Input", union.EntityRef.Name)
				require.Equal(t, "None", union.Tag)
				require.Nil(t, got.Entities["c1"].Const.Value.Message.Int)
			},
		},
		{
			name: "imported union with value",
			text: `
				const c2 pkg.Input = pkg.Input::Int(100)
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c2"].Const.Value.Message.Union
				require.Equal(t, "pkg", union.EntityRef.Pkg)
				require.Equal(t, "Input", union.EntityRef.Name)
				require.Equal(t, "Int", union.Tag)
				require.Equal(t, 100, *union.Data.Message.Int)
			},
		},
		{
			name: "nested union const value",
			text: "const c3 Input = Input::Nested(Input::Int(7))",
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c3"].Const.Value.Message.Union
				require.Equal(t, "", union.EntityRef.Pkg)
				require.Equal(t, "Input", union.EntityRef.Name)
				require.Equal(t, "Nested", union.Tag)

				nestedUnion := union.Data.Message.Union
				require.Equal(t, "", nestedUnion.EntityRef.Pkg)
				require.Equal(t, "Input", nestedUnion.EntityRef.Name)
				require.Equal(t, "Int", nestedUnion.Tag)
				require.Equal(t, 7, *nestedUnion.Data.Message.Int)
			},
		},
		{
			name: "tag-only",
			text: `
				type U union { Empty }
				const c4 U = U::Empty
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c4"].Const.Value.Message.Union
				require.Equal(t, "U", union.EntityRef.Name)
				require.Equal(t, "Empty", union.Tag)
				require.Nil(t, union.Data)
			},
		},
		{
			name: "wraps int positive",
			text: `
				type U union { Int int }
				const c5 U = U::Int(42)
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c5"].Const.Value.Message.Union
				require.Equal(t, "U", union.EntityRef.Name)
				require.Equal(t, "Int", union.Tag)
				require.Equal(t, 42, *union.Data.Message.Int)
			},
		},
		{
			name: "wraps int negative",
			text: `
				type U union { Int int }
				const c6 U = U::Int(-7)
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c6"].Const.Value.Message.Union
				require.Equal(t, "U", union.EntityRef.Name)
				require.Equal(t, "Int", union.Tag)
				require.Equal(t, -7, *union.Data.Message.Int)
			},
		},
		{
			name: "wraps float positive",
			text: `
				type U union { Float float }
				const c7 U = U::Float(1.5)
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c7"].Const.Value.Message.Union
				require.Equal(t, "U", union.EntityRef.Name)
				require.Equal(t, "Float", union.Tag)
				require.InEpsilon(t, 1.5, *union.Data.Message.Float, 0.0001)
			},
		},
		{
			name: "wraps float negative",
			text: `
				type U union { Float float }
				const c8 U = U::Float(-2.25)
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c8"].Const.Value.Message.Union
				require.Equal(t, "U", union.EntityRef.Name)
				require.Equal(t, "Float", union.Tag)
				require.InEpsilon(t, -2.25, *union.Data.Message.Float, 0.0001)
			},
		},
		{
			name: "wraps string",
			text: `
				type U union { Str string }
				const c9 U = U::Str('hello')
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c9"].Const.Value.Message.Union
				require.Equal(t, "U", union.EntityRef.Name)
				require.Equal(t, "Str", union.Tag)
				require.Equal(t, "hello", *union.Data.Message.Str)
			},
		},
		{
			name: "wraps bool true",
			text: `
				type U union { Flag bool }
				const c10 U = U::Flag(true)
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c10"].Const.Value.Message.Union
				require.Equal(t, "U", union.EntityRef.Name)
				require.Equal(t, "Flag", union.Tag)
				require.Equal(t, true, *union.Data.Message.Bool)
			},
		},
		{
			name: "wraps bool false",
			text: `
				type U union { Flag bool }
				const c11 U = U::Flag(false)
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c11"].Const.Value.Message.Union
				require.Equal(t, "U", union.EntityRef.Name)
				require.Equal(t, "Flag", union.Tag)
				require.Equal(t, false, *union.Data.Message.Bool)
			},
		},
		{
			name: "const ref to union const",
			text: `
				type U union { A }
				const c0 U = U::A
				const c1 U = c0
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				ref := got.Entities["c1"].Const.Value.Ref
				require.NotNil(t, ref)
				require.Equal(t, "c0", ref.Name)
			},
		},
		{
			name: "sender const tag-only",
			text: `
				type U union { Empty }
				const c12 U = U::Empty
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c12"].Const.Value.Message.Union
				require.Equal(t, "U", union.EntityRef.Name)
				require.Equal(t, "Empty", union.Tag)
				require.Nil(t, union.Data)
			},
		},
		{
			name: "sender const with int",
			text: `
				type U union { Int int }
				const c13 U = U::Int(5)
			`,
			check: func(t *testing.T, got src.File) {
			    t.Helper()
				union := got.Entities["c13"].Const.Value.Message.Union
				require.Equal(t, "U", union.EntityRef.Name)
				require.Equal(t, "Int", union.Tag)
				require.Equal(t, 5, *union.Data.Message.Int)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(tt.text))
			require.Nil(t, err)
			tt.check(t, got)
		})
	}
}

func TestParser_ParseFile_UnionLiteralConstSenders(t *testing.T) {
	tests := []struct { //nolint:govet // fieldalignment
		name  string
		text  string
		check func(t *testing.T, net []src.Connection)
	}{
		{
			name: "direct tag-only",
			text: `
				type U union { A }
				def C1() () {
					U::A -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				senderUnion := conn.Senders[0].Const.Value.Message.Union
				require.NotNil(t, senderUnion)
				require.Equal(t, "U", senderUnion.EntityRef.Name)
				require.Equal(t, "A", senderUnion.Tag)
				require.Nil(t, senderUnion.Data)
			},
		},
		{
			name: "direct int positive",
			text: `
				type U union { I int }
				def C1() () {
					U::I(7) -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				senderUnion := conn.Senders[0].Const.Value.Message.Union
				require.Equal(t, "U", senderUnion.EntityRef.Name)
				require.Equal(t, "I", senderUnion.Tag)
				require.Equal(t, 7, *senderUnion.Data.Message.Int)
			},
		},
		{
			name: "direct int negative",
			text: `
				type U union { I int }
				def C1() () {
					U::I(-3) -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				senderUnion := conn.Senders[0].Const.Value.Message.Union
				require.Equal(t, "U", senderUnion.EntityRef.Name)
				require.Equal(t, "I", senderUnion.Tag)
				require.Equal(t, -3, *senderUnion.Data.Message.Int)
			},
		},
		{
			name: "direct float positive",
			text: `
				type U union { F float }
				def C1() () {
					U::F(1.25) -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				senderUnion := conn.Senders[0].Const.Value.Message.Union
				require.Equal(t, "U", senderUnion.EntityRef.Name)
				require.Equal(t, "F", senderUnion.Tag)
				require.InEpsilon(t, 1.25, *senderUnion.Data.Message.Float, 0.0001)
			},
		},
		{
			name: "direct float negative",
			text: `
				type U union { F float }
				def C1() () {
					U::F(-2.5) -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				senderUnion := conn.Senders[0].Const.Value.Message.Union
				require.Equal(t, "U", senderUnion.EntityRef.Name)
				require.Equal(t, "F", senderUnion.Tag)
				require.InEpsilon(t, -2.5, *senderUnion.Data.Message.Float, 0.0001)
			},
		},
		{
			name: "direct string",
			text: `
				type U union { S string }
				def C1() () {
					U::S('hi') -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				senderUnion := conn.Senders[0].Const.Value.Message.Union
				require.Equal(t, "U", senderUnion.EntityRef.Name)
				require.Equal(t, "S", senderUnion.Tag)
				require.Equal(t, "hi", *senderUnion.Data.Message.Str)
			},
		},
		{
			name: "direct bool true",
			text: `
				type U union { B bool }
				def C1() () {
					U::B(true) -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				senderUnion := conn.Senders[0].Const.Value.Message.Union
				require.Equal(t, "U", senderUnion.EntityRef.Name)
				require.Equal(t, "B", senderUnion.Tag)
				require.Equal(t, true, *senderUnion.Data.Message.Bool)
			},
		},
		{
			name: "direct bool false",
			text: `
				type U union { B bool }
				def C1() () {
					U::B(false) -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				senderUnion := conn.Senders[0].Const.Value.Message.Union
				require.Equal(t, "U", senderUnion.EntityRef.Name)
				require.Equal(t, "B", senderUnion.Tag)
				require.Equal(t, false, *senderUnion.Data.Message.Bool)
			},
		},
		{
			name: "chained tag-only",
			text: `
				type U union { A }
				def C1() () {
					:start -> U::A -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				chain := conn.Receivers[0].ChainedConnection.Normal
				senderUnion := chain.Senders[0].Const.Value.Message.Union
				require.Equal(t, "U", senderUnion.EntityRef.Name)
				require.Equal(t, "A", senderUnion.Tag)
			},
		},
		{
			name: "chained with string",
			text: `
				type U union { S string }
				def C1() () {
					:start -> U::S('ok') -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
			    t.Helper()
				conn := net[0].Normal
				chain := conn.Receivers[0].ChainedConnection.Normal
				senderUnion := chain.Senders[0].Const.Value.Message.Union
				require.Equal(t, "U", senderUnion.EntityRef.Name)
				require.Equal(t, "S", senderUnion.Tag)
				require.Equal(t, "ok", *senderUnion.Data.Message.Str)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(tt.text))
			require.Nil(t, err)
			net := got.Entities["C1"].Component[0].Net
			tt.check(t, net)
		})
	}
}

func TestParser_ParseFile_ConnectionSendersConstRefAndPortAddr(t *testing.T) {
	text := `
		type U union { A }
		const c0 U = U::A
		def C1() (out any) {
			a Pass<any>
			---
			$c0 -> a
			a -> :out
			a:res -> :out
		}
	`

	p := New()
	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(text))
	require.Nil(t, err)

	net := got.Entities["C1"].Component[0].Net
	require.Len(t, net, 3)

	constSender := net[0].Normal.Senders[0].Const
	require.NotNil(t, constSender)
	require.NotNil(t, constSender.Value.Ref)
	require.Equal(t, "c0", constSender.Value.Ref.Name)

	implicitPortSender := net[1].Normal.Senders[0].PortAddr
	require.NotNil(t, implicitPortSender)
	require.Equal(t, "a", implicitPortSender.Node)
	require.Equal(t, "", implicitPortSender.Port)

	explicitPortSender := net[2].Normal.Senders[0].PortAddr
	require.NotNil(t, explicitPortSender)
	require.Equal(t, "a", explicitPortSender.Node)
	require.Equal(t, "res", explicitPortSender.Port)
}

func TestParser_ParseFile_StructLiteralTrailingComma(t *testing.T) {
	text := `
		const user User = { name: 'Ada', }
	`

	p := New()
	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(text))
	require.Nil(t, err)

	fields := got.Entities["user"].Const.Value.Message.DictOrStruct
	require.Equal(t, "Ada", *fields["name"].Message.Str)
}

func TestParser_ParseFile_OverloadedComponentDefinitions(t *testing.T) {
	text := []byte(`
		#extern(v1)
		pub def Add(left int, right int) (res int)
		#extern(v2)
		pub def Add(left float, right float) (res float)
		#extern(v3)
		pub def Add(left string, right string) (res string)
	`)

	p := New()

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.Nil(t, err)

	entity, ok := got.Entities["Add"]
	require.True(t, ok)
	require.Equal(t, src.ComponentEntity, entity.Kind)
	require.Equal(t, 3, len(entity.Component))

	// check extern directives preserved and inputs match
	require.Equal(t, "v1", entity.Component[0].Directives[compiler.ExternDirective])
	require.Equal(t, "int", entity.Component[0].IO.In["left"].TypeExpr.Inst.Ref.Name)
	require.Equal(t, "int", entity.Component[0].IO.In["right"].TypeExpr.Inst.Ref.Name)

	require.Equal(t, "v2", entity.Component[1].Directives[compiler.ExternDirective])
	require.Equal(t, "float", entity.Component[1].IO.In["left"].TypeExpr.Inst.Ref.Name)
	require.Equal(t, "float", entity.Component[1].IO.In["right"].TypeExpr.Inst.Ref.Name)

	require.Equal(t, "v3", entity.Component[2].Directives[compiler.ExternDirective])
	require.Equal(t, "string", entity.Component[2].IO.In["left"].TypeExpr.Inst.Ref.Name)
	require.Equal(t, "string", entity.Component[2].IO.In["right"].TypeExpr.Inst.Ref.Name)
}

package parser

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseFile_LonelyPorts(t *testing.T) {
	text := []byte(`
		component C1() () {
			:port -> lonely
			lonely -> :port
			:port -> lonely -> :port
		}
	`)

	p := New(false)

	got, err := p.parseFile(src.Location{}, text)
	require.True(t, err == nil)

	// 1) :port -> lonely
	// 2) lonely -> :port
	// 3) :port -> lonely
	// 4) lonely -> :port
	net := got.Entities["C1"].Component.Net
	require.Equal(t, 4, len(net))

	// 1) :port -> lonely
	receiverPortAddr := net[0].Normal.ReceiverSide.Receivers[0].PortAddr
	require.Equal(t, "lonely", receiverPortAddr.Node)
	require.Equal(t, "", receiverPortAddr.Port)

	// 2) lonely -> :port
	senderPortAddr := net[1].Normal.SenderSide.PortAddr
	require.Equal(t, "lonely", senderPortAddr.Node)
	require.Equal(t, "", senderPortAddr.Port)

	// 3) :port -> lonely
	receiverPortAddr = net[2].Normal.ReceiverSide.Receivers[0].PortAddr
	require.Equal(t, "lonely", receiverPortAddr.Node)
	require.Equal(t, "", receiverPortAddr.Port)

	// 4) lonely -> :port
	senderPortAddr = net[3].Normal.SenderSide.PortAddr
	require.Equal(t, "lonely", senderPortAddr.Node)
	require.Equal(t, "", senderPortAddr.Port)
}

func TestParser_ParseFile_ChainedConnections(t *testing.T) {
	text := []byte(`
		component C1() () { :in -> n1:p1 -> :out }
	`)

	p := New(false)

	got, err := p.parseFile(src.Location{}, text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 2, len(net))
}

func TestParser_ParseFile_Comments(t *testing.T) {
	text := []byte(`
	// comment
	`)

	p := New(false)

	_, err := p.parseFile(src.Location{}, text)
	require.True(t, err == nil)
}

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

	p := New(false)

	got, err := p.parseFile(src.Location{}, text)
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

func TestParser_ParseFile_IONodes(t *testing.T) {
	text := []byte(`
		component C1(start any) (stop any) {
			net { :start -> :stop }
		}
	`)

	p := New(false)

	got, err := p.parseFile(src.Location{}, text)
	require.True(t, err == nil)

	conn := got.Entities["C1"].Component.Net[0]

	sender := conn.Normal.SenderSide.PortAddr.Node
	require.Equal(t, "in", sender)

	receiver := conn.Normal.ReceiverSide.Receivers[0].PortAddr.Node
	require.Equal(t, "out", receiver)
}

func TestParser_ParseFile_AnonymousNodes(t *testing.T) {
	text := []byte(`
		component C1(start any) (stop any) {
			nodes {
				Scanner
				Printer<int>
			}
		}
	`)

	p := New(false)

	got, err := p.parseFile(src.Location{}, text)
	require.True(t, err == nil)

	nodes := got.Entities["C1"].Component.Nodes

	_, ok := nodes["scanner"]
	require.Equal(t, true, ok)

	_, ok = nodes["printer"]
	require.Equal(t, true, ok)
}

func TestParser_ParseFile_EnumLiterals(t *testing.T) {
	text := []byte(`
		const c0 Enum = Enum1::Foo
		const c1 pkg.Enum = pkg.Enum2::Bar
	`)

	p := New(false)

	got, err := p.parseFile(src.Location{}, text)
	require.True(t, err == nil)

	enum := got.Entities["c0"].Const.Message.Enum
	require.Equal(t, "", enum.EnumRef.Pkg)
	require.Equal(t, "Enum1", enum.EnumRef.Name)
	require.Equal(t, "Foo", enum.MemberName)

	enum = got.Entities["c1"].Const.Message.Enum
	require.Equal(t, "pkg", enum.EnumRef.Pkg)
	require.Equal(t, "Enum2", enum.EnumRef.Name)
	require.Equal(t, "Bar", enum.MemberName)
}

func TestParser_ParseFile_EnumLiteralSenders(t *testing.T) {
	text := []byte(`
		component C1() () {
			net {
				Foo::Bar -> :out
				foo.Bar::Baz -> :out
			}
		}
	`)

	p := New(false)

	got, err := p.parseFile(src.Location{}, text)
	require.True(t, err == nil)

	conn := got.Entities["C1"].Component.Net[0]

	senderEnum := conn.Normal.SenderSide.Const.Message.Enum
	require.Equal(t, "", senderEnum.EnumRef.Pkg)
	require.Equal(t, "Foo", senderEnum.EnumRef.Name)
	require.Equal(t, "Bar", senderEnum.MemberName)

	conn = got.Entities["C1"].Component.Net[1]
	senderEnum = conn.Normal.SenderSide.Const.Message.Enum
	require.Equal(t, "foo", senderEnum.EnumRef.Pkg)
	require.Equal(t, "Bar", senderEnum.EnumRef.Name)
	require.Equal(t, "Baz", senderEnum.MemberName)
}

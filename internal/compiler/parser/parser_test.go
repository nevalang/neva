package parser

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseFile_TernaryExpression(t *testing.T) {
	text := []byte(`
		def C1() () {
			(condition ? trueValue : falseValue) -> receiver
		}
	`)

	p := New()

	got, err := p.parseFile(text)
	require.NoError(t, err)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 1, len(net))

	conn := net[0].Normal
	require.Equal(t, 1, len(conn.SenderSide))

	ternary := conn.SenderSide[0].TernaryExpr
	require.NotNil(t, ternary)

	require.Equal(t, "condition", ternary.Condition.PortAddr.Node)
	require.Equal(t, "trueValue", ternary.Left.PortAddr.Node)
	require.Equal(t, "falseValue", ternary.Right.PortAddr.Node)

	require.Equal(t, "receiver", conn.ReceiverSide[0].PortAddr.Node)
}

func TestParser_ParseFile_StructSelectorsWithLonelyChain(t *testing.T) {
	text := []byte(`
		def C1() () {
			userSender -> .pet.name -> println -> :stop
		}`,
	)
	p := New()
	got, err := p.parseFile(text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 1, len(net))

	conn := net[0].Normal
	require.Equal(t, "userSender", conn.SenderSide[0].PortAddr.Node)
	require.Equal(t, "", conn.SenderSide[0].PortAddr.Port)

	chain := conn.ReceiverSide[0].ChainedConnection.Normal
	require.Equal(t, "pet", chain.SenderSide[0].StructSelector[0])
	require.Equal(t, "name", chain.SenderSide[0].StructSelector[1])

	secondChain := chain.ReceiverSide[0].ChainedConnection.Normal
	require.Equal(t, "println", secondChain.SenderSide[0].PortAddr.Node)
	require.Equal(t, "", secondChain.SenderSide[0].PortAddr.Port)

	chainEnd := secondChain.ReceiverSide[0].PortAddr
	require.Equal(t, "stop", chainEnd.Port)
}

func TestParser_ParseFile_PortlessArrPortAddr(t *testing.T) {
	text := []byte(`
		def C1() () {
			foo[0] -> bar[255]
		}
	`)

	p := New()

	got, err := p.parseFile(text)
	require.Equal(t, true, err == nil)

	net := got.Entities["C1"].Component.Net
	conn := net[0].Normal

	// foo[0]->
	require.Equal(t, "foo", conn.SenderSide[0].PortAddr.Node)
	require.Equal(t, "", conn.SenderSide[0].PortAddr.Port)
	require.Equal(t, compiler.Pointer(uint8(0)), conn.SenderSide[0].PortAddr.Idx)

	// ->bar[255]
	require.Equal(t, "bar", conn.ReceiverSide[0].PortAddr.Node)
	require.Equal(t, "", conn.ReceiverSide[0].PortAddr.Port)
	require.Equal(t, compiler.Pointer(uint8(255)), conn.ReceiverSide[0].PortAddr.Idx)
}

func TestParser_ParseFile_ChainedConnectionsWithDefer(t *testing.T) {
	text := []byte(`
		def C1() () {
			:start -> { foo -> bar -> :stop }
		}
	`)

	p := New()

	got, err := p.parseFile(text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 1, len(net))

	conn := net[0].Normal
	require.Equal(t, "in", conn.SenderSide[0].PortAddr.Node)
	require.Equal(t, "start", conn.SenderSide[0].PortAddr.Port)

	deferred := conn.ReceiverSide[0].DeferredConnection

	deferSender := deferred.Normal.SenderSide[0].PortAddr
	require.Equal(t, "foo", deferSender.Node)
	require.Equal(t, "", deferSender.Port)

	chainHead := deferred.Normal.ReceiverSide[0].ChainedConnection.Normal
	require.Equal(t, "bar", chainHead.SenderSide[0].PortAddr.Node)
	require.Equal(t, "", chainHead.SenderSide[0].PortAddr.Port)

	chainTail := chainHead.ReceiverSide[0].PortAddr
	require.Equal(t, "out", chainTail.Node)
	require.Equal(t, "stop", chainTail.Port)
}

func TestParser_ParseFile_LonelyPorts(t *testing.T) {
	text := []byte(`
		def C1() () {
			:port -> lonely
			lonely -> :port
		}
	`)

	p := New()

	got, err := p.parseFile(text)
	require.True(t, err == nil)

	// 1) :port -> lonely
	// 2) lonely -> :port
	net := got.Entities["C1"].Component.Net
	require.Equal(t, 2, len(net))

	// 1) :port -> lonely
	receiverPortAddr := net[0].Normal.ReceiverSide[0].PortAddr
	require.Equal(t, "lonely", receiverPortAddr.Node)
	require.Equal(t, "", receiverPortAddr.Port)

	// 2) lonely -> :port
	senderPortAddr := net[1].Normal.SenderSide[0].PortAddr
	require.Equal(t, "lonely", senderPortAddr.Node)
	require.Equal(t, "", senderPortAddr.Port)
}

func TestParser_ParseFile_ChainedConnections(t *testing.T) {
	text := []byte(`
		def C1() () { :foo -> n1:p1 -> :bar }
	`)

	p := New()

	got, err := p.parseFile(text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 1, len(net))
	conn := net[0].Normal

	sender := conn.SenderSide[0].PortAddr
	require.Equal(t, "in", sender.Node)
	require.Equal(t, "foo", sender.Port)

	chain := conn.ReceiverSide[0].ChainedConnection.Normal
	chainSender := chain.SenderSide[0].PortAddr
	require.Equal(t, "n1", chainSender.Node)
	require.Equal(t, "p1", chainSender.Port)

	chainReceiver := chain.ReceiverSide[0].PortAddr
	require.Equal(t, "out", chainReceiver.Node)
	require.Equal(t, "bar", chainReceiver.Port)
}

func TestParser_ParseFile_Comments(t *testing.T) {
	text := []byte(`
	// comment
	`)

	p := New()

	_, err := p.parseFile(text)
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

	got, err := p.parseFile(text)
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
		def C1(start any) (stop any) {
			:start -> :stop
		}
	`)

	p := New()

	got, err := p.parseFile(text)
	require.True(t, err == nil)

	conn := got.Entities["C1"].Component.Net[0]

	sender := conn.Normal.SenderSide[0].PortAddr.Node
	require.Equal(t, "in", sender)

	receiver := conn.Normal.ReceiverSide[0].PortAddr.Node
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

	got, err := p.parseFile(text)
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

	p := New()

	got, err := p.parseFile(text)
	require.True(t, err == nil)

	enum := got.Entities["c0"].Const.Value.Message.Enum
	require.Equal(t, "", enum.EnumRef.Pkg)
	require.Equal(t, "Enum1", enum.EnumRef.Name)
	require.Equal(t, "Foo", enum.MemberName)

	enum = got.Entities["c1"].Const.Value.Message.Enum
	require.Equal(t, "pkg", enum.EnumRef.Pkg)
	require.Equal(t, "Enum2", enum.EnumRef.Name)
	require.Equal(t, "Bar", enum.MemberName)
}

func TestParser_ParseFile_EnumLiteralSenders(t *testing.T) {
	text := []byte(`
		def C1() () {
			Foo::Bar -> :out
			foo.Bar::Baz -> :out
		}
	`)

	p := New()

	got, err := p.parseFile(text)
	require.True(t, err == nil)

	conn := got.Entities["C1"].Component.Net[0]

	senderEnum := conn.Normal.SenderSide[0].Const.Value.Message.Enum
	require.Equal(t, "", senderEnum.EnumRef.Pkg)
	require.Equal(t, "Foo", senderEnum.EnumRef.Name)
	require.Equal(t, "Bar", senderEnum.MemberName)

	conn = got.Entities["C1"].Component.Net[1]
	senderEnum = conn.Normal.SenderSide[0].Const.Value.Message.Enum
	require.Equal(t, "foo", senderEnum.EnumRef.Pkg)
	require.Equal(t, "Bar", senderEnum.EnumRef.Name)
	require.Equal(t, "Baz", senderEnum.MemberName)
}

func TestParser_ParseFile_RangeExpression(t *testing.T) {
	text := []byte(`
		def C1() () {
			1..10 -> :out
		}
	`)

	p := New()

	got, err := p.parseFile(text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 1, len(net))

	conn := net[0].Normal
	require.NotNil(t, conn.SenderSide[0].Range)
	require.Equal(t, int64(1), conn.SenderSide[0].Range.From)
	require.Equal(t, int64(10), conn.SenderSide[0].Range.To)
	require.Equal(t, "out", conn.ReceiverSide[0].PortAddr.Port)
}

func TestParser_ParseFile_MultipleRangeExpressions(t *testing.T) {
	text := []byte(`
		def C1() () {
			1..5 -> :out1
			10..20 -> :out2
		}
	`)

	p := New()

	got, err := p.parseFile(text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 2, len(net))

	conn1 := net[0].Normal
	require.NotNil(t, conn1.SenderSide[0].Range)
	require.Equal(t, int64(1), conn1.SenderSide[0].Range.From)
	require.Equal(t, int64(5), conn1.SenderSide[0].Range.To)
	require.Equal(t, "out1", conn1.ReceiverSide[0].PortAddr.Port)

	conn2 := net[1].Normal
	require.NotNil(t, conn2.SenderSide[0].Range)
	require.Equal(t, int64(10), conn2.SenderSide[0].Range.From)
	require.Equal(t, int64(20), conn2.SenderSide[0].Range.To)
	require.Equal(t, "out2", conn2.ReceiverSide[0].PortAddr.Port)
}

func TestParser_ParseFile_RangeExpressionWithNegativeNumbers(t *testing.T) {
	text := []byte(`
		def C1() () {
			-5..5 -> :out
		}
	`)

	p := New()

	got, err := p.parseFile(text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 1, len(net))

	conn := net[0].Normal
	require.NotNil(t, conn.SenderSide[0].Range)
	require.Equal(t, int64(-5), conn.SenderSide[0].Range.From)
	require.Equal(t, int64(5), conn.SenderSide[0].Range.To)
	require.Equal(t, "out", conn.ReceiverSide[0].PortAddr.Port)
}

func TestParser_ParseFile_RangeExpressionMixedWithOtherConnections(t *testing.T) {
	text := []byte(`
		def C1() () {
			1..10 -> :out1
			:in -> :out2
			20..30 -> :out3
		}
	`)

	p := New()

	got, err := p.parseFile(text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 3, len(net))

	conn1 := net[0].Normal
	require.NotNil(t, conn1.SenderSide[0].Range)
	require.Equal(t, int64(1), conn1.SenderSide[0].Range.From)
	require.Equal(t, int64(10), conn1.SenderSide[0].Range.To)
	require.Equal(t, "out1", conn1.ReceiverSide[0].PortAddr.Port)

	conn2 := net[1].Normal
	require.Nil(t, conn2.SenderSide[0].Range)
	require.Equal(t, "in", conn2.SenderSide[0].PortAddr.Node)
	require.Equal(t, "out2", conn2.ReceiverSide[0].PortAddr.Port)

	conn3 := net[2].Normal
	require.NotNil(t, conn3.SenderSide[0].Range)
	require.Equal(t, int64(20), conn3.SenderSide[0].Range.From)
	require.Equal(t, int64(30), conn3.SenderSide[0].Range.To)
	require.Equal(t, "out3", conn3.ReceiverSide[0].PortAddr.Port)
}

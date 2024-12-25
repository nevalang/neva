package parser

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
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

func TestParser_ParseFile_TernaryExpression(t *testing.T) {
	text := []byte(`
		def C1() () {
			(condition ? trueValue : falseValue) -> receiver
		}
	`)

	p := New()

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.Nil(t, err)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 1, len(net))

	conn := net[0].Normal
	require.Equal(t, 1, len(conn.Senders))

	ternary := conn.Senders[0].Ternary
	require.NotNil(t, ternary)

	require.Equal(t, "condition", ternary.Condition.PortAddr.Node)
	require.Equal(t, "trueValue", ternary.Left.PortAddr.Node)
	require.Equal(t, "falseValue", ternary.Right.PortAddr.Node)

	require.Equal(t, "receiver", conn.Receivers[0].PortAddr.Node)
}

func TestParser_ParseFile_NestedTernaryExpression(t *testing.T) {
	text := []byte(`
		def C1() () {
			(cond1 ? (cond2 ? val1 : val2) : val3) -> receiver
		}
	`)

	p := New()

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.Nil(t, err)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 1, len(net))

	conn := net[0].Normal
	require.Equal(t, 1, len(conn.Senders))

	outerTernary := conn.Senders[0].Ternary
	require.NotNil(t, outerTernary)

	require.Equal(t, "cond1", outerTernary.Condition.PortAddr.Node)
	require.Equal(t, "val3", outerTernary.Right.PortAddr.Node)

	innerTernary := outerTernary.Left.Ternary
	require.NotNil(t, innerTernary)
	require.Equal(t, "cond2", innerTernary.Condition.PortAddr.Node)
	require.Equal(t, "val1", innerTernary.Left.PortAddr.Node)
	require.Equal(t, "val2", innerTernary.Right.PortAddr.Node)

	require.Equal(t, "receiver", conn.Receivers[0].PortAddr.Node)
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

	net := got.Entities["C1"].Component.Net
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

	net := got.Entities["C1"].Component.Net
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

func TestParser_ParseFile_ChainedConnectionsWithDefer(t *testing.T) {
	text := []byte(`
		def C1() () {
			:start -> { foo -> bar -> :stop }
		}
	`)

	p := New()

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.True(t, err == nil)

	net := got.Entities["C1"].Component.Net
	require.Equal(t, 1, len(net))

	conn := net[0].Normal
	require.Equal(t, "in", conn.Senders[0].PortAddr.Node)
	require.Equal(t, "start", conn.Senders[0].PortAddr.Port)

	deferred := conn.Receivers[0].DeferredConnection

	deferSender := deferred.Normal.Senders[0].PortAddr
	require.Equal(t, "foo", deferSender.Node)
	require.Equal(t, "", deferSender.Port)

	chainHead := deferred.Normal.Receivers[0].ChainedConnection.Normal
	require.Equal(t, "bar", chainHead.Senders[0].PortAddr.Node)
	require.Equal(t, "", chainHead.Senders[0].PortAddr.Port)

	chainTail := chainHead.Receivers[0].PortAddr
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

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.True(t, err == nil)

	// 1) :port -> lonely
	// 2) lonely -> :port
	net := got.Entities["C1"].Component.Net
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

	net := got.Entities["C1"].Component.Net
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
}

func TestParser_ParseFile_ChainedConnectionsWithConstants(t *testing.T) {
	tests := []struct {
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

			net := got.Entities["C1"].Component.Net
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

	got, err := p.parseFile(location.ModRef, location.Package, location.Filename, text)
	require.True(t, err == nil)

	conn := got.Entities["C1"].Component.Net[0]

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

	nodes := got.Entities["C1"].Component.Nodes

	_, ok := nodes["scanner"]
	require.Equal(t, true, ok)

	_, ok = nodes["printer"]
	require.Equal(t, true, ok)
}

func TestParser_ParseFile_Range(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		check func(t *testing.T, net []src.Connection)
	}{
		{
			name: "simple range",
			text: `
				def C1() () {
					1..10 -> :out
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				conn := net[0].Normal
				require.NotNil(t, conn.Senders[0].Range)
				require.Equal(t, int64(1), conn.Senders[0].Range.From)
				require.Equal(t, int64(10), conn.Senders[0].Range.To)
				require.Equal(t, "out", conn.Receivers[0].PortAddr.Port)
			},
		},
		{
			name: "multiple ranges",
			text: `
				def C1() () {
					1..5 -> :out1
					10..20 -> :out2
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				require.Equal(t, 2, len(net))

				conn1 := net[0].Normal
				require.NotNil(t, conn1.Senders[0].Range)
				require.Equal(t, int64(1), conn1.Senders[0].Range.From)
				require.Equal(t, int64(5), conn1.Senders[0].Range.To)
				require.Equal(t, "out1", conn1.Receivers[0].PortAddr.Port)

				conn2 := net[1].Normal
				require.NotNil(t, conn2.Senders[0].Range)
				require.Equal(t, int64(10), conn2.Senders[0].Range.From)
				require.Equal(t, int64(20), conn2.Senders[0].Range.To)
				require.Equal(t, "out2", conn2.Receivers[0].PortAddr.Port)
			},
		},
		{
			name: "negative_from",
			text: `
				def C1() () {
					-5..5 -> :out
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				require.Equal(t, 1, len(net))

				conn := net[0].Normal
				require.NotNil(t, conn.Senders[0].Range)
				require.Equal(t, int64(-5), conn.Senders[0].Range.From)
				require.Equal(t, int64(5), conn.Senders[0].Range.To)
				require.Equal(t, "out", conn.Receivers[0].PortAddr.Port)
			},
		},
		{
			name: "negative_to",
			text: `
				def C1() () {
					1..-5 -> :out
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				require.Equal(t, 1, len(net))

				conn := net[0].Normal
				require.NotNil(t, conn.Senders[0].Range)
				require.Equal(t, int64(1), conn.Senders[0].Range.From)
				require.Equal(t, int64(-5), conn.Senders[0].Range.To)
				require.Equal(t, "out", conn.Receivers[0].PortAddr.Port)
			},
		},
		{
			name: "mixed range expressions",
			text: `
				def C1() () {
					1..10 -> :out1
					:in -> :out2
					20..30 -> :out3
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				require.Equal(t, 3, len(net))

				conn1 := net[0].Normal
				require.NotNil(t, conn1.Senders[0].Range)
				require.Equal(t, int64(1), conn1.Senders[0].Range.From)
				require.Equal(t, int64(10), conn1.Senders[0].Range.To)
				require.Equal(t, "out1", conn1.Receivers[0].PortAddr.Port)

				conn2 := net[1].Normal
				require.Nil(t, conn2.Senders[0].Range)
				require.Equal(t, "in", conn2.Senders[0].PortAddr.Node)
				require.Equal(t, "out2", conn2.Receivers[0].PortAddr.Port)

				conn3 := net[2].Normal
				require.NotNil(t, conn3.Senders[0].Range)
				require.Equal(t, int64(20), conn3.Senders[0].Range.From)
				require.Equal(t, int64(30), conn3.Senders[0].Range.To)
				require.Equal(t, "out3", conn3.Receivers[0].PortAddr.Port)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()

			got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(tt.text))
			require.Nil(t, err)

			net := got.Entities["C1"].Component.Net
			tt.check(t, net)
		})
	}
}

func TestParser_ParseFile_Binary(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		operator string
	}{
		// Arithmetic
		{
			name: "addition",
			text: `
				def C1() () {
					(a + b) -> receiver
				}
			`,
			operator: "+",
		},
		{
			name: "subtraction",
			text: `
				def C1() () {
					(a - b) -> receiver
				}
			`,
			operator: "-",
		},
		{
			name: "multiplication",
			text: `
				def C1() () {
					(a * b) -> receiver
				}
			`,
			operator: "*",
		},
		{
			name: "division",
			text: `
				def C1() () {
					(a / b) -> receiver
				}
			`,
			operator: "/",
		},
		{
			name: "modulo",
			text: `
				def C1() () {
					(a % b) -> receiver
				}
			`,
			operator: "%",
		},
		{
			name: "power",
			text: `
				def C1() () {
					(a ** b) -> receiver
				}
			`,
			operator: "**",
		},
		// Comparison
		{
			name: "equality",
			text: `
				def C1() () {
					(a == b) -> receiver
				}
			`,
			operator: "==",
		},
		{
			name: "inequality",
			text: `
				def C1() () {
					(a != b) -> receiver
				}
			`,
			operator: "!=",
		},
		{
			name: "greater than",
			text: `
				def C1() () {
					(a > b) -> receiver
				}
			`,
			operator: ">",
		},
		{
			name: "less than",
			text: `
				def C1() () {
					(a < b) -> receiver
				}
			`,
			operator: "<",
		},
		{
			name: "greater than or equal",
			text: `
				def C1() () {
					(a >= b) -> receiver
				}
			`,
			operator: ">=",
		},
		{
			name: "less than or equal",
			text: `
				def C1() () {
					(a <= b) -> receiver
				}
			`,
			operator: "<=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()

			got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(tt.text))
			require.Nil(t, err)

			net := got.Entities["C1"].Component.Net
			require.Equal(t, 1, len(net))

			conn := net[0].Normal
			require.Equal(t, 1, len(conn.Senders))

			binary := conn.Senders[0].Binary
			require.NotNil(t, binary)

			require.Equal(t, "a", binary.Left.PortAddr.Node)
			require.Equal(t, "b", binary.Right.PortAddr.Node)
			require.Equal(t, src.BinaryOperator(tt.operator), binary.Operator)
			require.Equal(t, "receiver", conn.Receivers[0].PortAddr.Node)
		})
	}
}

func TestParser_ParseFile_ComplexBinaryAndTernary(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		check func(t *testing.T, conn *src.NormalConnection)
	}{
		{
			name: "nested binary expressions",
			text: `
				def C1() () {
					((a + b) * (c - d)) -> receiver
				}
			`,
			check: func(t *testing.T, conn *src.NormalConnection) {
				binary := conn.Senders[0].Binary
				require.NotNil(t, binary)
				require.Equal(t, src.MulOp, binary.Operator)

				leftBinary := binary.Left.Binary
				require.NotNil(t, leftBinary)
				require.Equal(t, src.AddOp, leftBinary.Operator)
				require.Equal(t, "a", leftBinary.Left.PortAddr.Node)
				require.Equal(t, "b", leftBinary.Right.PortAddr.Node)

				rightBinary := binary.Right.Binary
				require.NotNil(t, rightBinary)
				require.Equal(t, src.SubOp, rightBinary.Operator)
				require.Equal(t, "c", rightBinary.Left.PortAddr.Node)
				require.Equal(t, "d", rightBinary.Right.PortAddr.Node)
			},
		},
		{
			name: "binary with ternary",
			text: `
				def C1() () {
					(a + (b ? c : d)) -> receiver
				}
			`,
			check: func(t *testing.T, conn *src.NormalConnection) {
				binary := conn.Senders[0].Binary
				require.NotNil(t, binary)
				require.Equal(t, src.AddOp, binary.Operator)
				require.Equal(t, "a", binary.Left.PortAddr.Node)

				ternary := binary.Right.Ternary
				require.NotNil(t, ternary)
				require.Equal(t, "b", ternary.Condition.PortAddr.Node)
				require.Equal(t, "c", ternary.Left.PortAddr.Node)
				require.Equal(t, "d", ternary.Right.PortAddr.Node)
			},
		},
		{
			name: "ternary with binary branches",
			text: `
				def C1() () {
					(cond ? (a + b) : (c * d)) -> receiver
				}
			`,
			check: func(t *testing.T, conn *src.NormalConnection) {
				ternary := conn.Senders[0].Ternary
				require.NotNil(t, ternary)
				require.Equal(t, "cond", ternary.Condition.PortAddr.Node)

				leftBinary := ternary.Left.Binary
				require.NotNil(t, leftBinary)
				require.Equal(t, src.AddOp, leftBinary.Operator)
				require.Equal(t, "a", leftBinary.Left.PortAddr.Node)
				require.Equal(t, "b", leftBinary.Right.PortAddr.Node)

				rightBinary := ternary.Right.Binary
				require.NotNil(t, rightBinary)
				require.Equal(t, src.MulOp, rightBinary.Operator)
				require.Equal(t, "c", rightBinary.Left.PortAddr.Node)
				require.Equal(t, "d", rightBinary.Right.PortAddr.Node)
			},
		},
		{
			name: "ternary with binary condition",
			text: `
				def C1() () {
					((a == b) ? c : d) -> receiver
				}
			`,
			check: func(t *testing.T, conn *src.NormalConnection) {
				ternary := conn.Senders[0].Ternary
				require.NotNil(t, ternary)

				condBinary := ternary.Condition.Binary
				require.NotNil(t, condBinary)
				require.Equal(t, src.EqOp, condBinary.Operator)
				require.Equal(t, "a", condBinary.Left.PortAddr.Node)
				require.Equal(t, "b", condBinary.Right.PortAddr.Node)

				require.Equal(t, "c", ternary.Left.PortAddr.Node)
				require.Equal(t, "d", ternary.Right.PortAddr.Node)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()

			got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(tt.text))
			require.Nil(t, err)

			net := got.Entities["C1"].Component.Net
			require.Equal(t, 1, len(net))

			conn := net[0].Normal
			require.Equal(t, 1, len(conn.Senders))
			require.Equal(t, "receiver", conn.Receivers[0].PortAddr.Node)

			tt.check(t, conn)
		})
	}
}

func TestParser_ParseFile_Switch(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		check func(t *testing.T, net []src.Connection)
	}{
		{
			name: "simple switch",
			text: `
				def C1() () {
					sender -> switch {
						true -> receiver1
						false -> receiver2
						_ -> receiver3
					}
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				conn := net[0].Normal
				require.Equal(t, "sender", conn.Senders[0].PortAddr.Node)

				switchStmt := conn.Receivers[0].Switch
				require.Equal(t, 2, len(switchStmt.Cases))

				// true -> receiver1
				require.Equal(t, true, *switchStmt.Cases[0].Senders[0].Const.Value.Message.Bool)
				require.Equal(t, "receiver1", switchStmt.Cases[0].Receivers[0].PortAddr.Node)

				// false -> receiver2
				require.Equal(t, false, *switchStmt.Cases[1].Senders[0].Const.Value.Message.Bool)
				require.Equal(t, "receiver2", switchStmt.Cases[1].Receivers[0].PortAddr.Node)

				// default -> receiver3
				require.Equal(t, "receiver3", switchStmt.Default[0].PortAddr.Node)
			},
		},
		{
			name: "switch with multiple senders and receivers",
			text: `
				def C1() () {
					sender -> switch {
						[a, b] -> [receiver1, receiver2]
						c -> [receiver3, receiver4]
					}
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				conn := net[0].Normal
				require.Equal(t, "sender", conn.Senders[0].PortAddr.Node)

				switchStmt := conn.Receivers[0].Switch
				require.Equal(t, 2, len(switchStmt.Cases))

				// [a, b] -> [receiver1, receiver2]
				require.Equal(t, 2, len(switchStmt.Cases[0].Senders))
				require.Equal(t, "a", switchStmt.Cases[0].Senders[0].PortAddr.Node)
				require.Equal(t, "b", switchStmt.Cases[0].Senders[1].PortAddr.Node)
				require.Equal(t, 2, len(switchStmt.Cases[0].Receivers))
				require.Equal(t, "receiver1", switchStmt.Cases[0].Receivers[0].PortAddr.Node)
				require.Equal(t, "receiver2", switchStmt.Cases[0].Receivers[1].PortAddr.Node)

				// c -> [receiver3, receiver4]
				require.Equal(t, 1, len(switchStmt.Cases[1].Senders))
				require.Equal(t, "c", switchStmt.Cases[1].Senders[0].PortAddr.Node)
				require.Equal(t, 2, len(switchStmt.Cases[1].Receivers))
				require.Equal(t, "receiver3", switchStmt.Cases[1].Receivers[0].PortAddr.Node)
				require.Equal(t, "receiver4", switchStmt.Cases[1].Receivers[1].PortAddr.Node)
			},
		},
		{
			name: "switch with binary expressions",
			text: `
				def C1() () {
					sender -> switch {
						(a + b) -> receiver1
						(c * d) -> receiver2
					}
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				conn := net[0].Normal
				require.Equal(t, "sender", conn.Senders[0].PortAddr.Node)

				switchStmt := conn.Receivers[0].Switch
				require.Equal(t, 2, len(switchStmt.Cases))

				// (a + b) -> receiver1
				binary1 := switchStmt.Cases[0].Senders[0].Binary
				require.Equal(t, src.AddOp, binary1.Operator)
				require.Equal(t, "a", binary1.Left.PortAddr.Node)
				require.Equal(t, "b", binary1.Right.PortAddr.Node)
				require.Equal(t, "receiver1", switchStmt.Cases[0].Receivers[0].PortAddr.Node)

				// (c * d) -> receiver2
				binary2 := switchStmt.Cases[1].Senders[0].Binary
				require.Equal(t, src.MulOp, binary2.Operator)
				require.Equal(t, "c", binary2.Left.PortAddr.Node)
				require.Equal(t, "d", binary2.Right.PortAddr.Node)
				require.Equal(t, "receiver2", switchStmt.Cases[1].Receivers[0].PortAddr.Node)
			},
		},
		{
			name: "nested switch",
			text: `
				def C1() () {
					sender -> switch {
						true -> switch {
							1 -> receiver1
							2 -> receiver2
						}
						false -> receiver3
					}
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				conn := net[0].Normal
				require.Equal(t, "sender", conn.Senders[0].PortAddr.Node)

				switchStmt := conn.Receivers[0].Switch
				require.Equal(t, 2, len(switchStmt.Cases))

				// true -> switch {...}
				require.Equal(t, true, *switchStmt.Cases[0].Senders[0].Const.Value.Message.Bool)
				nestedSwitch := switchStmt.Cases[0].Receivers[0].Switch
				require.Equal(t, 2, len(nestedSwitch.Cases))

				// 1 -> receiver1
				require.Equal(t, int(1), *nestedSwitch.Cases[0].Senders[0].Const.Value.Message.Int)
				require.Equal(t, "receiver1", nestedSwitch.Cases[0].Receivers[0].PortAddr.Node)

				// 2 -> receiver2
				require.Equal(t, int(2), *nestedSwitch.Cases[1].Senders[0].Const.Value.Message.Int)
				require.Equal(t, "receiver2", nestedSwitch.Cases[1].Receivers[0].PortAddr.Node)

				// false -> receiver3
				require.Equal(t, false, *switchStmt.Cases[1].Senders[0].Const.Value.Message.Bool)
				require.Equal(t, "receiver3", switchStmt.Cases[1].Receivers[0].PortAddr.Node)
			},
		},
		{
			name: "switch in chained connection",
			text: `
				def C1() () {
					sender -> .field -> switch {
						true -> receiver1
						false -> receiver2
					}
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				// sender ->
				conn := net[0].Normal
				require.Equal(t, "sender", conn.Senders[0].PortAddr.Node)

				// -> .field
				chain := conn.Receivers[0].ChainedConnection.Normal
				require.Equal(t, "field", chain.Senders[0].StructSelector[0])

				// -> switch {...}
				switchStmt := chain.Receivers[0].Switch
				require.Equal(t, 2, len(switchStmt.Cases))

				// { true -> receiver1
				require.Equal(t, true, *switchStmt.Cases[0].Senders[0].Const.Value.Message.Bool)
				require.Equal(t, "receiver1", switchStmt.Cases[0].Receivers[0].PortAddr.Node)

				// false -> receiver2 }
				require.Equal(t, false, *switchStmt.Cases[1].Senders[0].Const.Value.Message.Bool)
				require.Equal(t, "receiver2", switchStmt.Cases[1].Receivers[0].PortAddr.Node)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()

			got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(tt.text))
			require.Nil(t, err)

			net := got.Entities["C1"].Component.Net
			require.Equal(t, 1, len(net))

			tt.check(t, net)
		})
	}
}

func TestParser_ParseFile_TaggedUnionTypeExpr(t *testing.T) {
	tests := []struct {
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

func TestParser_ParseFile_TaggedUnionSender(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		check func(t *testing.T, net []src.Connection)
	}{
		{
			name: "non-chained tag-only",
			text: `
				def C1() () {
					Input::Int -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				conn := net[0].Normal
				senderUnion := conn.Senders[0].Const.Value.Message.Union
				require.Equal(t, "Input", senderUnion.EntityRef.Name)
				require.Equal(t, "Int", senderUnion.Tag)
				require.Equal(t, "receiver", conn.Receivers[0].PortAddr.Node)
			},
		},
		{
			name: "chained tag-only",
			text: `
				def C1() () {
					:start -> Input::Int -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				conn := net[0].Normal
				chain := conn.Receivers[0].ChainedConnection.Normal
				senderUnion := chain.Senders[0].Const.Value.Message.Union
				require.Equal(t, "Input", senderUnion.EntityRef.Name)
				require.Equal(t, "Int", senderUnion.Tag)
				require.Equal(t, "receiver", chain.Receivers[0].PortAddr.Node)
			},
		},
		{
			name: "non-chained with value",
			text: `
				def C1() () {
					Input::Int(42) -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				conn := net[0].Normal
				senderUnion := conn.Senders[0].Const.Value.Message.Union
				require.Equal(t, "Input", senderUnion.EntityRef.Name)
				require.Equal(t, "Int", senderUnion.Tag)
				require.Equal(t, 42, *conn.Senders[0].Const.Value.Message.Int)
				require.Equal(t, "receiver", conn.Receivers[0].PortAddr.Node)
			},
		},
		{
			name: "chained with value",
			text: `
				def C1() () {
					:start -> Input::Int(42) -> receiver
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				conn := net[0].Normal
				chain := conn.Receivers[0].ChainedConnection.Normal
				senderUnion := chain.Senders[0].Const.Value.Message.Union
				require.Equal(t, "Input", senderUnion.EntityRef.Name)
				require.Equal(t, "Int", senderUnion.Tag)
				require.Equal(t, 42, *chain.Senders[0].Const.Value.Message.Int)
				require.Equal(t, "receiver", chain.Receivers[0].PortAddr.Node)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(tt.text))
			require.Nil(t, err)
			net := got.Entities["C1"].Component.Net
			tt.check(t, net)
		})
	}
}

func TestParser_ParseFile_TaggedUnionPatternMatching(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		check func(t *testing.T, net []src.Connection)
	}{
		{
			name: "pattern matching without value",
			text: `
				def C1() () {
					sender -> switch {
						Input::Int -> receiver1
						Input::None -> receiver2
						_ -> receiver3
					}
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				require.Equal(t, 1, len(net))

				conn := net[0].Normal
				require.Equal(t, "sender", conn.Senders[0].PortAddr.Node)

				switchStmt := conn.Receivers[0].Switch
				require.Equal(t, 2, len(switchStmt.Cases))

				// Input::Int -> receiver1
				intCase := switchStmt.Cases[0].Senders[0].Const.Value.Message.Union
				require.Equal(t, "Input", intCase.EntityRef.Name)
				require.Equal(t, "Int", intCase.Tag)
				require.Equal(t, "receiver1", switchStmt.Cases[0].Receivers[0].PortAddr.Node)

				// Input::None -> receiver2
				noneCase := switchStmt.Cases[1].Senders[0].Const.Value.Message.Union
				require.Equal(t, "Input", noneCase.EntityRef.Name)
				require.Equal(t, "None", noneCase.Tag)
				require.Equal(t, "receiver2", switchStmt.Cases[1].Receivers[0].PortAddr.Node)

				// _ -> receiver3
				require.Equal(t, "receiver3", switchStmt.Default[0].PortAddr.Node)
			},
		},
		{
			name: "pattern matching with value",
			text: `
				def C1() () {
					sender -> switch {
						Input::Int(42) -> receiver1
						Input::None -> receiver2
						_ -> receiver3
					}
				}
			`,
			check: func(t *testing.T, net []src.Connection) {
				require.Equal(t, 1, len(net))

				conn := net[0].Normal
				require.Equal(t, "sender", conn.Senders[0].PortAddr.Node)

				switchStmt := conn.Receivers[0].Switch
				require.Equal(t, 2, len(switchStmt.Cases))

				// Input::Int(42) -> receiver1
				intCase := switchStmt.Cases[0].Senders[0].Const.Value.Message.Union
				require.Equal(t, "Input", intCase.EntityRef.Name)
				require.Equal(t, "Int", intCase.Tag)
				require.Equal(t, 42, *switchStmt.Cases[0].Senders[0].Const.Value.Message.Int)
				require.Equal(t, "receiver1", switchStmt.Cases[0].Receivers[0].PortAddr.Node)

				// Input::None -> receiver2
				noneCase := switchStmt.Cases[1].Senders[0].Const.Value.Message.Union
				require.Equal(t, "Input", noneCase.EntityRef.Name)
				require.Equal(t, "None", noneCase.Tag)
				require.Equal(t, "receiver2", switchStmt.Cases[1].Receivers[0].PortAddr.Node)

				// _ -> receiver3
				require.Equal(t, "receiver3", switchStmt.Default[0].PortAddr.Node)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			got, err := p.parseFile(location.ModRef, location.Package, location.Filename, []byte(tt.text))
			require.Nil(t, err)
			net := got.Entities["C1"].Component.Net
			tt.check(t, net)
		})
	}
}

func TestParser_ParseFile_TaggedUnionConstLiteral(t *testing.T) {
	tests := []struct {
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

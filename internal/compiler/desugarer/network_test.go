package desugarer

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: some cases are hard to test this way because desugarer depends on Scope object
// which is normally passed from top-level functions in this package.
func TestDesugarNetwork(t *testing.T) {
	tests := []struct {
		name           string
		iface          src.Interface
		mockScope      func(scope *MockScopeMockRecorder)
		net            []src.Connection
		nodes          map[string]src.Node
		expectedResult handleNetworkResult
	}{
		{
			// node1:out -> node2:in
			name: "simple_1_to_1",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{
								PortAddr: &src.PortAddr{Node: "node1", Port: "out"},
							},
						},
						Receivers: []src.ConnectionReceiver{
							{
								PortAddr: &src.PortAddr{Node: "node2", Port: "in"},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"node1": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node1"}},
				"node2": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node2"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "node1", Port: "out"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "node2", Port: "in"},
								},
							},
						},
					},
				},
				constsToInsert: map[string]src.Const{},
				nodesToInsert:  map[string]src.Node{},
			},
		},
		// [node1:out, node2:out] -> node3:in
		{
			name: "fan_in",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{PortAddr: &src.PortAddr{Node: "node1", Port: "out"}},
							{PortAddr: &src.PortAddr{Node: "node2", Port: "out"}},
						},
						Receivers: []src.ConnectionReceiver{
							{
								PortAddr: &src.PortAddr{Node: "node3", Port: "in"},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"node1": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node1"}},
				"node2": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node2"}},
				"node3": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node3"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "node1", Port: "out"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "__fan_in__1", Port: "data", Idx: compiler.Pointer(uint8(0))},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "node2", Port: "out"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "__fan_in__1", Port: "data", Idx: compiler.Pointer(uint8(1))},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "__fan_in__1", Port: "res"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "node3", Port: "in"},
								},
							},
						},
					},
				},
				constsToInsert: map[string]src.Const{},
				nodesToInsert: map[string]src.Node{
					"__fan_in__1": {
						EntityRef: core.EntityRef{Pkg: "builtin", Name: "FanIn"},
					},
				},
			},
		},
		// node1:foo -> node2:bar -> node3:baz
		{
			name: "chained",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{
								PortAddr: &src.PortAddr{Node: "node1", Port: "foo"},
							},
						},
						Receivers: []src.ConnectionReceiver{
							{
								ChainedConnection: &src.Connection{
									Normal: &src.NormalConnection{
										Senders: []src.ConnectionSender{
											{
												PortAddr: &src.PortAddr{Node: "node2", Port: "bar"},
											},
										},
										Receivers: []src.ConnectionReceiver{
											{
												PortAddr: &src.PortAddr{Node: "node3", Port: "baz"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"node1": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node1"}},
				"node2": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node2"}},
				"node3": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node3"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "node1", Port: "foo"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "node2", Port: "bar"},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "node2", Port: "bar"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "node3", Port: "baz"},
								},
							},
						},
					},
				},
				constsToInsert: map[string]src.Const{},
				nodesToInsert:  map[string]src.Node{},
			},
		},
		// foo:bar -> .a.b.c -> baz:bax
		{
			name: "struct_selector_chain",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{
								PortAddr: &src.PortAddr{Node: "foo", Port: "bar"},
							},
						},
						Receivers: []src.ConnectionReceiver{
							{
								ChainedConnection: &src.Connection{
									Normal: &src.NormalConnection{
										Senders: []src.ConnectionSender{
											{StructSelector: []string{"a", "b", "c"}},
										},
										Receivers: []src.ConnectionReceiver{
											{PortAddr: &src.PortAddr{Node: "baz", Port: "bax"}},
										},
									},
								},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"foo": {EntityRef: core.EntityRef{Pkg: "test", Name: "Foo"}},
				"baz": {EntityRef: core.EntityRef{Pkg: "test", Name: "Baz"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "foo", Port: "bar"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "__field__1", Port: "data"},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "__field__1", Port: "res"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "baz", Port: "bax"},
								},
							},
						},
					},
				},
				nodesToInsert: map[string]src.Node{
					"__field__1": {
						EntityRef: core.EntityRef{Pkg: "builtin", Name: "Field"},
						Directives: map[src.Directive]string{
							compiler.BindDirective: "__const__1",
						},
					},
				},
				constsToInsert: map[string]src.Const{
					"__const__1": {
						TypeExpr: ts.Expr{
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Pkg: "builtin", Name: "list"},
								Args: []ts.Expr{
									{
										Inst: &ts.InstExpr{
											Ref: core.EntityRef{Pkg: "builtin", Name: "string"},
										},
									},
								},
							},
						},
						Value: src.ConstValue{
							Message: &src.MsgLiteral{
								List: []src.ConstValue{
									{Message: &src.MsgLiteral{Str: compiler.Pointer("a")}},
									{Message: &src.MsgLiteral{Str: compiler.Pointer("b")}},
									{Message: &src.MsgLiteral{Str: compiler.Pointer("c")}},
								},
							},
						},
					},
				},
			},
		},
		// (:a + :b) -> :c
		{
			name: "binary_expressions",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{
								Binary: &src.Binary{
									Operator: src.AddOp,
									Left: src.ConnectionSender{
										PortAddr: &src.PortAddr{Port: "a"},
									},
									Right: src.ConnectionSender{
										PortAddr: &src.PortAddr{Port: "b"},
									},
									AnalyzedType: ts.Expr{
										Inst: &ts.InstExpr{
											Ref: core.EntityRef{Name: "int"},
										},
									},
								},
							},
						},
						Receivers: []src.ConnectionReceiver{
							{PortAddr: &src.PortAddr{Port: "c"}},
						},
					},
				},
			},
			mockScope: func(mock *MockScopeMockRecorder) {
				mock.
					Entity(core.EntityRef{Name: "Add"}).
					// In reality Add is overloaded, but for for now we only test simplified case with one version.
					Return(
						src.Entity{
							Kind: src.ComponentEntity,
							Component: []src.Component{
								{
									Interface: src.Interface{
										IO: src.IO{
											In: map[string]src.Port{
												"left":  {TypeExpr: ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}}},
												"right": {TypeExpr: ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}}},
											},
											Out: map[string]src.Port{
												"res": {TypeExpr: ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}}},
											},
										},
									},
								},
							},
						},
						core.Location{},
						nil,
					)
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						// __add__1:res -> :c
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Node: "__add__1", Port: "res"}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Port: "c"}},
							},
						},
					},
					{
						// :a -> __add__1:left
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Port: "a"}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "__add__1", Port: "left"}},
							},
						},
					},
					{
						// :b -> __add__1:right
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Port: "b"}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "__add__1", Port: "right"}},
							},
						},
					},
				},
				nodesToInsert: map[string]src.Node{
					"__add__1": {
						EntityRef: core.EntityRef{
							Pkg:  "builtin",
							Name: "Add",
						},
						TypeArgs:      nil,
						OverloadIndex: compiler.Pointer(0),
					},
				},
				constsToInsert: map[string]src.Const{},
			},
		},
		{
			// node1:x -> switch {
			//     node2:y -> node3:z
			//     node4:y -> node5:z
			//     _ -> node6:z
			// }
			name: "switch_receiver",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{PortAddr: &src.PortAddr{Node: "node1", Port: "x"}},
						},
						Receivers: []src.ConnectionReceiver{
							{
								Switch: &src.Switch{
									Cases: []src.NormalConnection{
										{
											Senders: []src.ConnectionSender{
												{PortAddr: &src.PortAddr{Node: "node2", Port: "y"}},
											},
											Receivers: []src.ConnectionReceiver{
												{PortAddr: &src.PortAddr{Node: "node3", Port: "z"}},
											},
										},
										{
											Senders: []src.ConnectionSender{
												{PortAddr: &src.PortAddr{Node: "node4", Port: "y"}},
											},
											Receivers: []src.ConnectionReceiver{
												{PortAddr: &src.PortAddr{Node: "node5", Port: "z"}},
											},
										},
									},
									Default: []src.ConnectionReceiver{
										{PortAddr: &src.PortAddr{Node: "node6", Port: "z"}},
									},
								},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"node1": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node1"}},
				"node2": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node2"}},
				"node3": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node3"}},
				"node4": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node4"}},
				"node5": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node5"}},
				"node6": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node6"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Node: "node1", Port: "x"}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "__switch__1", Port: "data"}},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Node: "node2", Port: "y"}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "__switch__1", Port: "case", Idx: compiler.Pointer(uint8(0))}},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Node: "__switch__1", Port: "case", Idx: compiler.Pointer(uint8(0))}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "node3", Port: "z"}},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Node: "node4", Port: "y"}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "__switch__1", Port: "case", Idx: compiler.Pointer(uint8(1))}},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Node: "__switch__1", Port: "case", Idx: compiler.Pointer(uint8(1))}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "node5", Port: "z"}},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Node: "__switch__1", Port: "else"}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "node6", Port: "z"}},
							},
						},
					},
				},
				nodesToInsert: map[string]src.Node{
					"__switch__1": {
						EntityRef: core.EntityRef{
							Pkg:  "builtin",
							Name: "Switch",
						},
					},
				},
				constsToInsert: map[string]src.Const{},
			},
		},
		// $foo -> bar:baz
		{
			name: "const_ref_sender",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{
								Const: &src.Const{
									Value: src.ConstValue{
										Ref: &core.EntityRef{Name: "foo"},
									},
								},
							},
						},
						Receivers: []src.ConnectionReceiver{
							{
								PortAddr: &src.PortAddr{Node: "bar", Port: "baz"},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"bar": {EntityRef: core.EntityRef{Name: "Bar"}},
			},
			mockScope: func(mock *MockScopeMockRecorder) {
				constEntity := src.Const{
					TypeExpr: ts.Expr{
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
					},
					Value: src.ConstValue{
						Message: &src.MsgLiteral{Int: compiler.Pointer(42)},
					},
				}
				mock.
					Entity(core.EntityRef{Name: "foo"}).
					Return(src.Entity{Kind: src.ConstEntity, Const: constEntity}, core.Location{}, nil)
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "__new__1", Port: "res"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "bar", Port: "baz"},
								},
							},
						},
					},
				},
				nodesToInsert: map[string]src.Node{
					"__new__1": {
						EntityRef: core.EntityRef{Pkg: "builtin", Name: "New"},
						TypeArgs: src.TypeArgs{
							{
								Inst: &ts.InstExpr{
									Ref: core.EntityRef{Name: "int"},
								},
							},
						},
						Directives: map[src.Directive]string{
							compiler.BindDirective: "foo",
						},
					},
				},
				constsToInsert: map[string]src.Const{},
			},
		},
		// a:b -> $c -> d:e
		{
			name: "const_ref_in_chain",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{
								PortAddr: &src.PortAddr{Node: "a", Port: "b"},
							},
						},
						Receivers: []src.ConnectionReceiver{
							{
								ChainedConnection: &src.Connection{
									Normal: &src.NormalConnection{
										Senders: []src.ConnectionSender{
											{
												Const: &src.Const{
													Value: src.ConstValue{
														Ref: &core.EntityRef{Name: "c"},
													},
												},
											},
										},
										Receivers: []src.ConnectionReceiver{
											{
												PortAddr: &src.PortAddr{Node: "d", Port: "e"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"a": {EntityRef: core.EntityRef{Name: "A"}},
				"d": {EntityRef: core.EntityRef{Name: "D"}},
			},
			mockScope: func(mock *MockScopeMockRecorder) {
				constEntity := src.Const{
					TypeExpr: ts.Expr{
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
					},
					Value: src.ConstValue{
						Message: &src.MsgLiteral{Int: compiler.Pointer(42)},
					},
				}
				mock.
					Entity(core.EntityRef{Name: "c"}).
					Return(src.Entity{Kind: src.ConstEntity, Const: constEntity}, core.Location{}, nil)
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "a", Port: "b"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "__newv2__1", Port: "sig"},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "__newv2__1", Port: "res"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "d", Port: "e"},
								},
							},
						},
					},
				},
				nodesToInsert: map[string]src.Node{
					"__newv2__1": {
						EntityRef: core.EntityRef{Pkg: "builtin", Name: "NewV2"},
						TypeArgs: src.TypeArgs{
							{
								Inst: &ts.InstExpr{
									Ref: core.EntityRef{Name: "int"},
								},
							},
						},
						Directives: map[src.Directive]string{
							compiler.BindDirective: "c",
						},
					},
				},
				constsToInsert: map[string]src.Const{},
			},
		},
		// a:b -> 42 -> c:d
		{
			name: "const_literal_in_chain",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{
								PortAddr: &src.PortAddr{Node: "a", Port: "b"},
							},
						},
						Receivers: []src.ConnectionReceiver{
							{
								ChainedConnection: &src.Connection{
									Normal: &src.NormalConnection{
										Senders: []src.ConnectionSender{
											{
												Const: &src.Const{
													TypeExpr: ts.Expr{
														Inst: &ts.InstExpr{
															Ref: core.EntityRef{Name: "int"},
														},
													},
													Value: src.ConstValue{
														Message: &src.MsgLiteral{Int: compiler.Pointer(42)},
													},
												},
											},
										},
										Receivers: []src.ConnectionReceiver{
											{
												PortAddr: &src.PortAddr{Node: "c", Port: "d"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"a": {EntityRef: core.EntityRef{Name: "A"}},
				"c": {EntityRef: core.EntityRef{Name: "C"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "a", Port: "b"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "__newv2__1", Port: "sig"},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "__newv2__1", Port: "res"},
								},
							},
							Receivers: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "c", Port: "d"},
								},
							},
						},
					},
				},
				nodesToInsert: map[string]src.Node{
					"__newv2__1": {
						EntityRef: core.EntityRef{Pkg: "builtin", Name: "NewV2"},
						TypeArgs: src.TypeArgs{
							{
								Inst: &ts.InstExpr{
									Ref: core.EntityRef{Name: "int"},
								},
							},
						},
						Directives: map[src.Directive]string{
							compiler.BindDirective: "__const__1",
						},
					},
				},
				constsToInsert: map[string]src.Const{
					"__const__1": {
						TypeExpr: ts.Expr{
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Name: "int"},
							},
						},
						Value: src.ConstValue{
							Message: &src.MsgLiteral{Int: compiler.Pointer(42)},
						},
					},
				},
			},
		},
		// (5 + 3) -> foo:bar
		{
			name: "binary_expression_with_literals",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{
								Binary: &src.Binary{
									Operator: src.AddOp,
									Left: src.ConnectionSender{
										Const: &src.Const{
											TypeExpr: ts.Expr{
												Inst: &ts.InstExpr{
													Ref: core.EntityRef{Name: "int"},
												},
											},
											Value: src.ConstValue{
												Message: &src.MsgLiteral{Int: compiler.Pointer(5)},
											},
										},
									},
									Right: src.ConnectionSender{
										Const: &src.Const{
											TypeExpr: ts.Expr{
												Inst: &ts.InstExpr{
													Ref: core.EntityRef{Name: "int"},
												},
											},
											Value: src.ConstValue{
												Message: &src.MsgLiteral{Int: compiler.Pointer(3)},
											},
										},
									},
									AnalyzedType: ts.Expr{
										Inst: &ts.InstExpr{
											Ref: core.EntityRef{Name: "int"},
										},
									},
								},
							},
						},
						Receivers: []src.ConnectionReceiver{
							{PortAddr: &src.PortAddr{Node: "foo", Port: "bar"}},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"foo": {EntityRef: core.EntityRef{Name: "Foo"}},
			},
			mockScope: func(mock *MockScopeMockRecorder) {
				mock.
					Entity(core.EntityRef{Name: "Add"}).
					// In reality Add is overloaded, but for for now we only test simplified case with one version.
					Return(
						src.Entity{
							Kind: src.ComponentEntity,
							Component: []src.Component{
								{
									Interface: src.Interface{
										IO: src.IO{
											In: map[string]src.Port{
												"left":  {TypeExpr: ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}}},
												"right": {TypeExpr: ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}}},
											},
											Out: map[string]src.Port{
												"res": {TypeExpr: ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}}},
											},
										},
									},
								},
							},
						},
						core.Location{},
						nil,
					)
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						// __add__1:res -> foo:bar
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Node: "__add__1", Port: "res"}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "foo", Port: "bar"}},
							},
						},
					},
					{
						// __new__1:res -> __add__1:left
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Node: "__new__1", Port: "res"}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "__add__1", Port: "left"}},
							},
						},
					},
					{
						// __new__2:res -> __add__1:right
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{
								{PortAddr: &src.PortAddr{Node: "__new__2", Port: "res"}},
							},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "__add__1", Port: "right"}},
							},
						},
					},
				},
				nodesToInsert: map[string]src.Node{
					"__add__1": {
						EntityRef:     core.EntityRef{Pkg: "builtin", Name: "Add"},
						TypeArgs:      nil,
						OverloadIndex: compiler.Pointer(0),
					},
					"__new__1": {
						EntityRef: core.EntityRef{Pkg: "builtin", Name: "New"},
						TypeArgs: []ts.Expr{
							{
								Inst: &ts.InstExpr{
									Ref: core.EntityRef{Name: "int"},
								},
							},
						},
						Directives: map[src.Directive]string{
							compiler.BindDirective: "__const__1",
						},
					},
					"__new__2": {
						EntityRef: core.EntityRef{Pkg: "builtin", Name: "New"},
						TypeArgs: []ts.Expr{
							{
								Inst: &ts.InstExpr{
									Ref: core.EntityRef{Name: "int"},
								},
							},
						},
						Directives: map[src.Directive]string{
							compiler.BindDirective: "__const__2",
						},
					},
				},
				constsToInsert: map[string]src.Const{
					"__const__1": {
						TypeExpr: ts.Expr{
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Name: "int"},
							},
						},
						Value: src.ConstValue{
							Message: &src.MsgLiteral{Int: compiler.Pointer(5)},
						},
					},
					"__const__2": {
						TypeExpr: ts.Expr{
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Name: "int"},
							},
						},
						Value: src.ConstValue{
							Message: &src.MsgLiteral{Int: compiler.Pointer(3)},
						},
					},
				},
			},
		},
		{
			name: "union_sender_tag_only",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{
								Union: &src.UnionSender{
									EntityRef: core.EntityRef{Name: "Input"},
									Tag:       "Int",
								},
							},
						},
						Receivers: []src.ConnectionReceiver{
							{PortAddr: &src.PortAddr{Node: "foo", Port: "bar"}},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"foo": {EntityRef: core.EntityRef{Name: "Foo"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{{
								PortAddr: &src.PortAddr{
									Node: "__new__1",
									Port: "res",
								},
							}},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "foo", Port: "bar"}},
							},
						},
					},
				},
				nodesToInsert: map[string]src.Node{
					"__new__1": {
						EntityRef: core.EntityRef{
							Pkg:  "builtin",
							Name: "New",
						},
						TypeArgs: src.TypeArgs{
							{
								Inst: &ts.InstExpr{
									Ref: core.EntityRef{Name: "Input"},
								},
							},
						},
						Directives: map[src.Directive]string{
							compiler.BindDirective: "__union_const__1",
						},
					},
				},
				constsToInsert: map[string]src.Const{
					"__union_const__1": {
						Value: src.ConstValue{
							Message: &src.MsgLiteral{
								Union: &src.UnionLiteral{
									EntityRef: core.EntityRef{Name: "Input"},
									Tag:       "Int",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "union_sender_with_value",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						Senders: []src.ConnectionSender{
							{
								Union: &src.UnionSender{
									EntityRef: core.EntityRef{Name: "Input"},
									Tag:       "Int",
									Data: &src.ConnectionSender{
										Const: &src.Const{
											TypeExpr: ts.Expr{
												Inst: &ts.InstExpr{
													Ref: core.EntityRef{Name: "int"},
												},
											},
											Value: src.ConstValue{
												Message: &src.MsgLiteral{Int: compiler.Pointer(42)},
											},
										},
									},
								},
							},
						},
						Receivers: []src.ConnectionReceiver{
							{PortAddr: &src.PortAddr{Node: "foo", Port: "bar"}},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"foo": {EntityRef: core.EntityRef{Name: "Foo"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{{
								PortAddr: &src.PortAddr{
									Node: "__union__1",
									Port: "res",
								},
							}},
							Receivers: []src.ConnectionReceiver{
								{PortAddr: &src.PortAddr{Node: "foo", Port: "bar"}},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							Senders: []src.ConnectionSender{{
								PortAddr: &src.PortAddr{
									Node: "__new__1",
									Port: "res",
								},
							}},
							Receivers: []src.ConnectionReceiver{{
								PortAddr: &src.PortAddr{
									Node: "__union__1",
									Port: "data",
								},
							}},
						},
					},
				},
				nodesToInsert: map[string]src.Node{
					"__new__1": {
						EntityRef: core.EntityRef{
							Pkg:  "builtin",
							Name: "New",
						},
						TypeArgs: src.TypeArgs{
							{
								Inst: &ts.InstExpr{
									Ref: core.EntityRef{Name: "int"},
								},
							},
						},
						Directives: map[src.Directive]string{
							compiler.BindDirective: "__const__2",
						},
					},
					"__union__1": {
						Directives: map[src.Directive]string{
							compiler.BindDirective: "__union_tag__1",
						},
						EntityRef: core.EntityRef{
							Pkg:  "builtin",
							Name: "UnionWrapV1",
						},
					},
				},
				constsToInsert: map[string]src.Const{
					"__union_tag__1": {
						Value: src.ConstValue{
							Message: &src.MsgLiteral{
								Str: compiler.Pointer("Int"),
							},
						},
					},
					"__const__2": {
						TypeExpr: ts.Expr{
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Name: "int"},
							},
						},
						Value: src.ConstValue{
							Message: &src.MsgLiteral{Int: compiler.Pointer(42)},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Desugarer{}

			scope := NewMockScope(gomock.NewController(t))
			if tt.mockScope != nil {
				tt.mockScope(scope.EXPECT())
			}

			result, err := d.desugarNetwork(tt.iface, tt.net, tt.nodes, scope)
			require.Nil(t, err)

			assert.Equal(t, tt.expectedResult.desugaredConnections, result.desugaredConnections)
			assert.Equal(t, tt.expectedResult.constsToInsert, result.constsToInsert)
			assert.Equal(t, tt.expectedResult.nodesToInsert, result.nodesToInsert)
		})
	}
}

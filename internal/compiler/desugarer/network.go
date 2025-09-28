package desugarer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	"github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type handleNetworkResult struct {
	desugaredConnections []src.Connection
	constsToInsert       map[string]src.Const
	nodesToInsert        map[string]src.Node
	nodesPortsUsed       nodeOutportsUsed
}

func (d *Desugarer) desugarNetwork(
	iface src.Interface,
	net []src.Connection,
	nodes map[string]src.Node,
	scope Scope,
) (handleNetworkResult, error) {
	nodesToInsert := map[string]src.Node{}
	constsToInsert := map[string]src.Const{}
	nodesPortsUsed := newNodePortsMap()

	desugaredConnections, err := d.desugarConnections(
		iface,
		net,
		nodesPortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return handleNetworkResult{}, err
	}

	return handleNetworkResult{
		desugaredConnections: desugaredConnections,
		nodesPortsUsed:       nodesPortsUsed,
		constsToInsert:       constsToInsert,
		nodesToInsert:        nodesToInsert,
	}, nil
}

func (d *Desugarer) desugarConnections(
	iface src.Interface,
	net []src.Connection,
	nodePortsUsed nodeOutportsUsed,
	scope Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) ([]src.Connection, error) {
	desugaredConnections := make([]src.Connection, 0, len(net))

	for _, conn := range net {
		result, err := d.desugarConnection(
			iface,
			conn,
			nodePortsUsed,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return nil, err
		}

		if result.replace != nil {
			desugaredConnections = append(desugaredConnections, *result.replace)
		}

		desugaredConnections = append(desugaredConnections, result.insert...)
	}

	return desugaredConnections, nil
}

type desugarConnectionResult struct {
	replace *src.Connection
	insert  []src.Connection
}

func (d *Desugarer) desugarConnection(
	iface src.Interface,
	conn src.Connection,
	nodePortsUsed nodeOutportsUsed,
	scope Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (desugarConnectionResult, error) {
	if conn.ArrayBypass != nil { // nothing to desugar, just mark ports as used
		nodePortsUsed.set(
			conn.ArrayBypass.SenderOutport.Node,
			conn.ArrayBypass.SenderOutport.Port,
		)
		nodePortsUsed.set(
			conn.ArrayBypass.ReceiverInport.Node,
			conn.ArrayBypass.ReceiverInport.Port,
		)
		return desugarConnectionResult{replace: &conn}, nil
	}

	return d.desugarNormalConnection(
		iface,
		*conn.Normal,
		nodePortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
	)
}

func (d *Desugarer) desugarNormalConnection(
	iface src.Interface,
	normConn src.NormalConnection,
	nodePortsUsed nodeOutportsUsed,
	scope Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (desugarConnectionResult, error) {
	if fanIn := len(normConn.Senders) > 1; fanIn {
		result, err := d.desugarFanIn(
			iface,
			normConn,
			nodesToInsert,
			constsToInsert,
			nodePortsUsed,
			scope,
			nodes,
		)
		if err != nil {
			return desugarConnectionResult{}, fmt.Errorf("desugar fan in: %w", err)
		}
		// one fan-in connection was replaced by multiple desugared ones
		return desugarConnectionResult{insert: result}, nil
	}

	desugarSenderResult, err := d.desugarSingleSender(
		iface,
		normConn,
		scope,
		nodes,
		nodePortsUsed,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return desugarConnectionResult{}, fmt.Errorf("desugar single sender: %w", err)
	}

	normConn = *desugarSenderResult.replace.Normal
	insert := desugarSenderResult.insert

	// We need to first degugar fan-out and convert it to a single receiver in this connection
	if len(normConn.Receivers) > 1 {
		result, err := d.desugarFanOut(
			iface,
			normConn,
			nodesToInsert,
			constsToInsert,
			nodePortsUsed,
			scope,
			nodes,
		)
		if err != nil {
			return desugarConnectionResult{}, err
		}

		return desugarConnectionResult{
			replace: &result.replace,
			insert:  append(insert, result.insert...),
		}, nil
	}

	desugarReceiverResult, err := d.desugarSingleReceiver(
		iface,
		normConn,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
		nodePortsUsed,
	)
	if err != nil {
		return desugarConnectionResult{}, err
	}

	return desugarConnectionResult{
		replace: &desugarReceiverResult.replace,
		insert:  append(insert, desugarReceiverResult.insert...),
	}, nil
}

type desugarReceiverResult struct {
	replace src.Connection
	insert  []src.Connection
}

// desugarSingleReceiver expects connection without fan-out (it must be desugared before).
func (d *Desugarer) desugarSingleReceiver(
	iface src.Interface,
	normConn src.NormalConnection,
	scope Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	nodePortsUsed nodeOutportsUsed,
) (desugarReceiverResult, error) {
	// FIXME for some reason normConn.Meta sometimes empty
	locOnlyMeta := core.Meta{
		Location: normConn.Senders[0].Meta.Location,
	}

	receiver := normConn.Receivers[0]

	if receiver.PortAddr != nil {
		if receiver.PortAddr.Port != "" {
			return desugarReceiverResult{
				replace: src.Connection{
					Normal: &src.NormalConnection{
						Senders:   normConn.Senders,
						Receivers: []src.ConnectionReceiver{receiver},
					},
					Meta: locOnlyMeta,
				},
				insert: []src.Connection{},
			}, nil
		}

		firstInportName, err := d.getFirstInportName(scope, nodes, *receiver.PortAddr)
		if err != nil {
			return desugarReceiverResult{}, fmt.Errorf("get first inport name: %w", err)
		}

		// if node is interface with anonymous port, port-addr will remain empty string
		// to be later desugared at irgen step, because it's not possible to do here
		return desugarReceiverResult{
			replace: src.Connection{
				Normal: &src.NormalConnection{
					Senders: normConn.Senders,
					Receivers: []src.ConnectionReceiver{
						{
							PortAddr: &src.PortAddr{
								Port: firstInportName,
								Node: receiver.PortAddr.Node,
								Idx:  receiver.PortAddr.Idx,
								Meta: receiver.PortAddr.Meta,
							},
						},
					},
					Meta: locOnlyMeta,
				},
				Meta: locOnlyMeta,
			},
		}, nil
	}

	if receiver.DeferredConnection != nil {
		result, err := d.desugarDeferredConnection(
			iface,
			normConn,
			scope,
			constsToInsert,
			nodesToInsert,
			nodePortsUsed,
			nodes,
		)
		if err != nil {
			return desugarReceiverResult{}, err
		}

		return desugarReceiverResult(result), nil
	}

	if receiver.Switch != nil {
		d.switchCounter++
		switchNodeName := fmt.Sprintf("__switch__%d", d.switchCounter)

		nodesToInsert[switchNodeName] = src.Node{
			EntityRef: core.EntityRef{
				Pkg:  "builtin",
				Name: "Switch",
				Meta: locOnlyMeta,
			},
			Meta: locOnlyMeta,
		}

		// Connect original sender to switch:data
		replace := src.Connection{
			Normal: &src.NormalConnection{
				Senders: normConn.Senders,
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: switchNodeName,
							Port: "data",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
			},
			Meta: locOnlyMeta,
		}

		insert := []src.Connection{}

		// For each case in the switch
		for i, caseConn := range receiver.Switch.Cases {
			// Connect case-sender to switch:case[i]
			insert = append(insert, src.Connection{
				Normal: &src.NormalConnection{
					Senders: caseConn.Senders,
					Receivers: []src.ConnectionReceiver{
						{
							PortAddr: &src.PortAddr{
								Node: switchNodeName,
								Port: "case",
								Idx:  compiler.Pointer(uint8(i)),
								Meta: locOnlyMeta,
							},
							Meta: locOnlyMeta,
						},
					},
				},
				Meta: locOnlyMeta,
			})

			// Connect switch:case[i] to case receiver
			insert = append(insert, src.Connection{
				Normal: &src.NormalConnection{
					Senders: []src.ConnectionSender{
						{
							PortAddr: &src.PortAddr{
								Node: switchNodeName,
								Port: "case",
								Idx:  compiler.Pointer(uint8(i)),
								Meta: locOnlyMeta,
							},
							Meta: locOnlyMeta,
						},
					},
					Receivers: caseConn.Receivers,
					Meta:      locOnlyMeta,
				},
				Meta: locOnlyMeta,
			})
		}

		// Connect switch:default to its receiver
		insert = append(insert, src.Connection{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{
							Node: switchNodeName,
							Port: "else",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Receivers: receiver.Switch.Default,
				Meta:      locOnlyMeta,
			},
		})

		desugaredInsert, err := d.desugarConnections(
			iface,
			insert,
			nodePortsUsed,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarReceiverResult{}, err
		}

		return desugarReceiverResult{
			replace: replace,
			insert:  desugaredInsert,
		}, nil
	}

	desugarChainResult, err := d.desugarChainedConnection(
		iface,
		receiver,
		scope,
		nodes,
		nodePortsUsed,
		nodesToInsert,
		constsToInsert,
		normConn,
	)
	if err != nil {
		return desugarReceiverResult{}, err
	}

	return desugarReceiverResult{
		replace: *desugarChainResult.replace,
		insert:  desugarChainResult.insert,
	}, nil
}

func (d *Desugarer) desugarChainedConnection(
	iface src.Interface,
	receiver src.ConnectionReceiver,
	scope Scope,
	nodes map[string]src.Node,
	nodePortsUsed nodeOutportsUsed,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	normConn src.NormalConnection,
) (desugarConnectionResult, error) {
	chainedConn := *receiver.ChainedConnection
	chainHead := chainedConn.Normal.Senders[0] // chain head is always single sender

	// it's only possible to find receiver port before desugaring of chained connection
	var chainHeadPort string
	switch {
	case chainHead.Range != nil:
		chainHeadPort = "sig" // Range has sig inport
	case chainHead.Const != nil:
		chainHeadPort = "sig" // NewV2 has sig inport
	case len(chainHead.StructSelector) != 0:
		chainHeadPort = "data"
	case chainHead.PortAddr != nil:
		chainHeadPort = chainHead.PortAddr.Port
		if chainHeadPort == "" {
			var err error
			chainHeadPort, err = d.getFirstInportName(scope, nodes, *chainHead.PortAddr)
			if err != nil {
				return desugarConnectionResult{}, fmt.Errorf("get first inport name: %w", err)
			}
		}
	default:
		panic("unexpected chain head type")
	}

	locOnlyMeta := core.Meta{Location: chainHead.Meta.Location}

	if chainHead.Const != nil {
		d.virtualTriggersCount++
		triggerNodeName := fmt.Sprintf("__newv2__%d", d.virtualTriggersCount)

		if chainHead.Const.Value.Ref != nil {
			constTypeExpr, err := d.getConstTypeByRef(*chainHead.Const.Value.Ref, scope)
			if err != nil {
				return desugarConnectionResult{}, fmt.Errorf("get const type by ref: %w", err)
			}

			nodesToInsert[triggerNodeName] = src.Node{
				EntityRef: core.EntityRef{
					Pkg:  "builtin",
					Name: "NewV2",
					Meta: locOnlyMeta,
				},
				TypeArgs: []ts.Expr{constTypeExpr},
				Directives: map[src.Directive]string{
					compiler.BindDirective: chainHead.Const.Value.Ref.String(),
				},
				Meta: locOnlyMeta,
			}
		} else {
			d.virtualConstCount++
			virtualConstName := fmt.Sprintf("__const__%d", d.virtualConstCount)
			constsToInsert[virtualConstName] = *chainHead.Const

			nodesToInsert[triggerNodeName] = src.Node{
				EntityRef: core.EntityRef{
					Pkg:  "builtin",
					Name: "NewV2",
					Meta: locOnlyMeta,
				},
				TypeArgs: []ts.Expr{chainHead.Const.TypeExpr},
				Directives: map[src.Directive]string{
					compiler.BindDirective: virtualConstName,
				},
				Meta: locOnlyMeta,
			}
		}

		chainedConn.Normal.Senders = []src.ConnectionSender{
			{
				PortAddr: &src.PortAddr{
					Node: triggerNodeName,
					Port: "res",
					Meta: locOnlyMeta,
				},
				Meta: locOnlyMeta,
			},
		}
	}

	desugarChainResult, err := d.desugarConnection(
		iface,
		chainedConn,
		nodePortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return desugarConnectionResult{}, err
	}

	desugaredHead := desugarChainResult.replace.Normal.Senders[0]

	replace := src.Connection{
		Normal: &src.NormalConnection{
			Senders: normConn.Senders,
			Receivers: []src.ConnectionReceiver{
				{
					PortAddr: &src.PortAddr{
						Node: desugaredHead.PortAddr.Node,
						Port: chainHeadPort,
						Idx:  desugaredHead.PortAddr.Idx,
						Meta: chainHead.Meta,
					},
					Meta: chainHead.Meta,
				},
			},
			Meta: locOnlyMeta,
		},
	}

	// we need to insert both: replace and insert, example:
	// input = a -> b -> c -> d
	// sender = a
	// chain = b -> c -> d
	// ---
	// replace = b -> c
	// insert = c -> d
	// (and replace existing one with `a -> b`)
	insert := append([]src.Connection{}, desugarChainResult.insert...)
	insert = append(insert, *desugarChainResult.replace)

	return desugarConnectionResult{
		replace: &replace,
		insert:  insert,
	}, nil
}

type desugarDeferredConnectionsResult struct {
	replace src.Connection
	insert  []src.Connection
}

func (d *Desugarer) desugarDeferredConnection(
	iface src.Interface,
	normConn src.NormalConnection,
	scope Scope,
	constsToInsert map[string]src.Const,
	nodesToInsert map[string]src.Node,
	nodesPortsUsed nodeOutportsUsed,
	nodes map[string]src.Node,
) (desugarDeferredConnectionsResult, error) {
	deferredConnection := *normConn.Receivers[0].DeferredConnection

	desugarDeferredConnResult, err := d.desugarConnection(
		iface,
		deferredConnection,
		nodesPortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return desugarDeferredConnectionsResult{}, err
	}

	locOnlyMeta := core.Meta{Location: deferredConnection.Meta.Location}

	deferredConnection = *desugarDeferredConnResult.replace
	connsToInsert := desugarDeferredConnResult.insert

	// 1) create lock node
	d.virtualLocksCounter++
	lockNodeName := fmt.Sprintf("__lock__%d", d.virtualLocksCounter)
	nodesToInsert[lockNodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "Lock",
			Meta: locOnlyMeta,
		},
		TypeArgs: []typesystem.Expr{
			ts.Expr{
				Inst: &typesystem.InstExpr{
					Ref: core.EntityRef{Pkg: "builtin", Name: "any"},
				},
				Meta: locOnlyMeta,
			},
		},
		Meta: locOnlyMeta,
	}

	// 2) connect original sender to lock receiver
	replace := src.Connection{
		Normal: &src.NormalConnection{
			Senders: normConn.Senders,
			Receivers: []src.ConnectionReceiver{
				{
					PortAddr: &src.PortAddr{
						Node: lockNodeName,
						Port: "sig",
						Meta: locOnlyMeta,
					},
				},
			},
		},
	}

	connsToInsert = append(
		// 3) connect deferred sender to lock data
		connsToInsert,
		src.Connection{
			Normal: &src.NormalConnection{
				Senders: deferredConnection.Normal.Senders,
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: lockNodeName,
							Port: "data",
							Meta: locOnlyMeta,
						},
					},
				},
				Meta: locOnlyMeta,
			},
		},
		// 4) create connection from lock:data to receiver-side of deferred connection
		src.Connection{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{
							Node: lockNodeName,
							Port: "data",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Receivers: deferredConnection.Normal.Receivers,
				Meta:      locOnlyMeta,
			},
		},
	)

	return desugarDeferredConnectionsResult{
		replace: replace,
		insert:  connsToInsert,
	}, nil
}

type desugarSenderResult struct {
	replace src.Connection   // receiver side might need desugaring
	insert  []src.Connection // already desugared
}

// desugarSingleSender keeps receiver side untouched so it must be desugared by caller (except for selectors).
func (d *Desugarer) desugarSingleSender(
	iface src.Interface,
	normConn src.NormalConnection,
	scope Scope,
	nodes map[string]src.Node,
	usedNodeOutports nodeOutportsUsed,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (desugarSenderResult, error) {
	sender := normConn.Senders[0]

	if sender.PortAddr != nil {
		portName := sender.PortAddr.Port
		if sender.PortAddr.Port == "" {
			firstOutportName, err := d.getFirstOutportName(scope, nodes, *sender.PortAddr)
			if err != nil {
				return desugarSenderResult{}, fmt.Errorf("get first outport name: %w", err)
			}
			portName = firstOutportName
			normConn.Senders = []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Port: portName, // <- desugaring
						Node: sender.PortAddr.Node,
						Idx:  sender.PortAddr.Idx,
						Meta: sender.PortAddr.Meta,
					},
					Meta: sender.Meta,
				},
			}
		}
		usedNodeOutports.set(
			sender.PortAddr.Node,
			portName,
		)
		// if node is interface with anonymous port, port-addr will remain empty string
		// to be later desugared at irgen step, because it's not possible to do here
		return desugarSenderResult{
			replace: src.Connection{Normal: &normConn},
			insert:  nil,
		}, nil
	}

	if len(sender.StructSelector) != 0 {
		result, err := d.desugarStructSelectors(
			normConn,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarSenderResult{}, fmt.Errorf("desugar struct selectors: %w", err)
		}

		// connection that replaces original one might need desugaring itself
		replacedConnDesugarRes, err := d.desugarConnection(
			iface,
			result.replace,
			usedNodeOutports,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarSenderResult{}, err
		}

		return desugarSenderResult{
			replace: src.Connection{
				Normal: replacedConnDesugarRes.replace.Normal,
				Meta:   replacedConnDesugarRes.replace.Meta,
			},
			insert: replacedConnDesugarRes.insert,
		}, nil
	}

	if sender.Const != nil {
		if sender.Const.Value.Ref != nil {
			portAddr, err := d.handleConstRefSender(*sender.Const.Value.Ref, nodesToInsert, scope)
			if err != nil {
				return desugarSenderResult{}, err
			}

			normConn = src.NormalConnection{
				Senders: []src.ConnectionSender{
					{
						PortAddr: &portAddr,
						Meta:     sender.Meta,
					},
				},
				Receivers: normConn.Receivers,
				Meta:      sender.Meta,
			}
		} else if sender.Const.Value.Message != nil {
			portAddr, err := d.handleLiteralSender(*sender.Const, nodesToInsert, constsToInsert)
			if err != nil {
				return desugarSenderResult{}, err
			}

			normConn = src.NormalConnection{
				Senders: []src.ConnectionSender{
					{
						PortAddr: &portAddr,
						Meta:     sender.Meta,
					},
				},
				Receivers: normConn.Receivers,
				Meta:      sender.Meta,
			}
		}

		return desugarSenderResult{
			replace: src.Connection{Normal: &normConn},
			insert:  nil,
		}, nil
	}

	if sender.Union != nil {
		result, err := d.desugarUnionSender(
			*sender.Union,
			normConn,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarSenderResult{}, fmt.Errorf("desugar union sender: %w", err)
		}
		return desugarSenderResult{
			replace: result.replace,
			insert:  result.insert,
		}, nil
	}

	if sender.Ternary != nil {
		result, err := d.desugarTernarySender(
			iface,
			*sender.Ternary,
			normConn,
			nodesToInsert,
			constsToInsert,
			usedNodeOutports,
			scope,
			nodes,
		)
		if err != nil {
			return desugarSenderResult{}, fmt.Errorf("desugar ternary sender: %w", err)
		}

		return desugarSenderResult(result), nil
	}

	if sender.Binary != nil {
		result, err := d.desugarBinarySender(
			iface,
			*sender.Binary,
			normConn,
			nodesToInsert,
			constsToInsert,
			usedNodeOutports,
			scope,
			nodes,
		)
		if err != nil {
			return desugarSenderResult{}, err
		}
		return desugarSenderResult(result), nil
	}

	result, err := d.desugarRangeSender(
		*sender.Range,
		normConn,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return desugarSenderResult{}, err
	}

	return desugarSenderResult{
		replace: src.Connection{Normal: &result.replace},
		insert:  result.insert,
	}, nil
}

func (d *Desugarer) getFirstInportName(
	scope Scope,
	nodes map[string]src.Node,
	portAddr src.PortAddr,
) (string, error) {
	io, err := scope.GetNodeIOByPortAddr(nodes, portAddr)
	if err != nil {
		return "", err
	}
	for inport := range io.In {
		return inport, nil
	}
	return "", errors.New("first inport not found")
}

func (d *Desugarer) getFirstOutportName(
	scope Scope,
	nodes map[string]src.Node,
	portAddr src.PortAddr,
) (string, error) {
	io, err := scope.GetNodeIOByPortAddr(nodes, portAddr)
	if err != nil {
		return "", err
	}

	// important: skip `err` outport if node has err guard
	for outport := range io.Out {
		if outport == "err" && nodes[portAddr.Node].ErrGuard {
			continue
		}
		return outport, nil
	}

	return "", errors.New("first outport not found")
}

var newComponentRef = core.EntityRef{
	Pkg:  "builtin",
	Name: "New",
}

func (d *Desugarer) handleLiteralSender(
	constant src.Const,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (src.PortAddr, error) {
	d.virtualConstCount++
	constName := fmt.Sprintf("__const__%d", d.virtualConstCount)

	// we can't call d.handleConstRefSender()
	// because our virtual const isn't in the scope

	d.virtualEmittersCount++
	emitterNodeName := fmt.Sprintf("__new__%d", d.virtualEmittersCount)

	locOnlyMeta := core.Meta{Location: constant.Meta.Location}

	emitterNode := src.Node{
		Directives: map[src.Directive]string{
			compiler.BindDirective: constName,
		},
		EntityRef: core.EntityRef{
			Pkg:  newComponentRef.Pkg,
			Name: newComponentRef.Name,
			Meta: locOnlyMeta,
		},
		TypeArgs: []ts.Expr{constant.TypeExpr},
		Meta:     locOnlyMeta,
	}

	nodesToInsert[emitterNodeName] = emitterNode
	constsToInsert[constName] = constant

	emitterNodeOutportAddr := src.PortAddr{
		Node: emitterNodeName,
		Port: "res",
	}

	return emitterNodeOutportAddr, nil
}

func (d *Desugarer) handleConstRefSender(
	ref core.EntityRef,
	nodesToInsert map[string]src.Node,
	scope Scope,
) (src.PortAddr, error) {
	constTypeExpr, err := d.getConstTypeByRef(ref, scope)
	if err != nil {
		return src.PortAddr{}, fmt.Errorf("get const type by ref: %w", err)
	}

	locOnlyMeta := core.Meta{Location: ref.Meta.Location}

	d.virtualEmittersCount++
	virtualEmitterName := fmt.Sprintf("__new__%d", d.virtualEmittersCount)

	emitterNode := src.Node{
		// don't forget to bind
		Directives: map[src.Directive]string{
			compiler.BindDirective: ref.String(),
		},
		EntityRef: core.EntityRef{
			Pkg:  newComponentRef.Pkg,
			Name: newComponentRef.Name,
			Meta: locOnlyMeta,
		},
		TypeArgs: []ts.Expr{constTypeExpr},
		Meta:     locOnlyMeta,
	}

	emitterNodeOutportAddr := src.PortAddr{
		Node: virtualEmitterName,
		Port: "res",
		Meta: locOnlyMeta,
	}

	nodesToInsert[virtualEmitterName] = emitterNode

	return emitterNodeOutportAddr, nil
}

// getConstTypeByRef is needed to figure out type parameters for Const node
func (d *Desugarer) getConstTypeByRef(ref core.EntityRef, scope Scope) (ts.Expr, error) {
	entity, _, err := scope.Entity(ref)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("get entity: %w", err)
	}

	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, fmt.Errorf("entity is not a constant: %v", entity.Kind)
	}

	if entity.Const.Value.Ref != nil {
		expr, err := d.getConstTypeByRef(*entity.Const.Value.Ref, scope)
		if err != nil {
			return ts.Expr{}, fmt.Errorf("get const type by ref: %w", err)
		}
		return expr, nil
	}

	return entity.Const.TypeExpr, nil
}

type desugarFanOutResult struct {
	replace src.Connection   // original sender -> fanOut receiver
	insert  []src.Connection // fanOut sender -> original receivers
}

func (d *Desugarer) desugarFanOut(
	iface src.Interface,
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	nodePortsUsed nodeOutportsUsed,
	scope Scope,
	nodes map[string]src.Node,
) (desugarFanOutResult, error) {
	d.fanOutCounter++
	nodeName := fmt.Sprintf("__fan_out__%d", d.fanOutCounter)

	locOnlyMeta := core.Meta{Location: normConn.Senders[0].Meta.Location} // FIXME for some reason norm-conn sometimes doesn't have meta

	nodesToInsert[nodeName] = src.Node{
		EntityRef: core.EntityRef{
			Name: "FanOut",
			Pkg:  "builtin",
			Meta: locOnlyMeta,
		},
		Meta: locOnlyMeta,
	}

	receiverToReplace := src.ConnectionReceiver{
		PortAddr: &src.PortAddr{
			Node: nodeName,
			Port: "data",
			Meta: locOnlyMeta,
		},
		Meta: locOnlyMeta,
	}

	insert := make([]src.Connection, 0, len(normConn.Receivers))
	for i, receiver := range normConn.Receivers {
		conn := src.Connection{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{
							Node: nodeName,
							Port: "data",
							Idx:  compiler.Pointer(uint8(i)),
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Receivers: []src.ConnectionReceiver{receiver},
				Meta:      locOnlyMeta,
			},
		}

		desugarConnRes, err := d.desugarConnection(
			iface,
			conn,
			nodePortsUsed,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarFanOutResult{}, err
		}

		insert = append(insert, *desugarConnRes.replace)
		insert = append(insert, desugarConnRes.insert...)
	}

	return desugarFanOutResult{
		replace: src.Connection{
			Normal: &src.NormalConnection{
				Senders:   normConn.Senders, // senders must be desugared
				Receivers: []src.ConnectionReceiver{receiverToReplace},
				Meta:      locOnlyMeta,
			},
			Meta: locOnlyMeta,
		},
		insert: insert,
	}, nil
}

// Add a new function to handle range senders
type handleRangeSenderResult struct {
	replace src.NormalConnection
	insert  []src.Connection
}

// desugarRangeSender desugars `from..to -> XXX` part.
// It does not create connection to range:sig,
// it's done in chained connection desugaring.
func (d *Desugarer) desugarRangeSender(
	rangeExpr src.Range,
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (handleRangeSenderResult, error) {
	locOnlyMeta := core.Meta{Location: rangeExpr.Meta.Location}

	d.rangeCounter++

	rangeNodeName := fmt.Sprintf("__range%d__", d.rangeCounter)
	fromConstName := fmt.Sprintf("__range%d_from__", d.rangeCounter)
	toConstName := fmt.Sprintf("__range%d_to__", d.rangeCounter)

	constsToInsert[fromConstName] = src.Const{
		TypeExpr: ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Pkg: "builtin", Name: "int"}}},
		Value:    src.ConstValue{Message: &src.MsgLiteral{Int: compiler.Pointer(int(rangeExpr.From))}},
		Meta:     locOnlyMeta,
	}
	constsToInsert[toConstName] = src.Const{
		TypeExpr: ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Pkg: "builtin", Name: "int"}}},
		Value:    src.ConstValue{Message: &src.MsgLiteral{Int: compiler.Pointer(int(rangeExpr.To))}},
		Meta:     locOnlyMeta,
	}

	nodesToInsert[rangeNodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "Range",
			Meta: locOnlyMeta,
		},
		Meta: locOnlyMeta,
	}
	nodesToInsert[fromConstName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "New",
			Meta: locOnlyMeta,
		},
		Directives: map[src.Directive]string{
			compiler.BindDirective: fromConstName,
		},
		Meta: locOnlyMeta,
	}
	nodesToInsert[toConstName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "New",
			Meta: locOnlyMeta,
		},
		Directives: map[src.Directive]string{
			compiler.BindDirective: toConstName,
		},
		Meta: locOnlyMeta,
	}

	replace := src.NormalConnection{
		Senders: []src.ConnectionSender{
			{
				PortAddr: &src.PortAddr{
					Node: rangeNodeName,
					Port: "res",
					Meta: locOnlyMeta,
				},
				Meta: locOnlyMeta,
			},
		},
		Receivers: normConn.Receivers,
		Meta:      locOnlyMeta,
	}

	insert := []src.Connection{
		// $from -> range:from
		{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{
							Node: fromConstName,
							Port: "res",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: rangeNodeName,
							Port: "from",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Meta: locOnlyMeta,
			},
			Meta: locOnlyMeta,
		},
		// $to -> range:to
		{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{
							Node: toConstName,
							Port: "res",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: rangeNodeName,
							Port: "to",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Meta: locOnlyMeta,
			},
			Meta: locOnlyMeta,
		},
	}

	return handleRangeSenderResult{
		insert:  insert,
		replace: replace,
	}, nil
}

// desugarFanIn returns connections that must be used instead of given one.
// It recursevely desugars each connection before return so result is final.
func (d *Desugarer) desugarFanIn(
	iface src.Interface,
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	nodePortsUsed nodeOutportsUsed,
	scope Scope,
	nodes map[string]src.Node,
) ([]src.Connection, error) {
	locOnlyMeta := normConn.Senders[0].Meta // FIXME for some reason normConn.Meta sometimes empty

	// 1. create fan-in node with unique name and insert into nodes
	d.fanInCounter++
	fanInNodeName := fmt.Sprintf("__fan_in__%d", d.fanInCounter)
	nodesToInsert[fanInNodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "FanIn",
			Meta: locOnlyMeta,
		},
		Meta: locOnlyMeta,
	}

	// 2. connect each sender of this connection with fan-in node
	desugaredFanIn := make([]src.Connection, 0, len(normConn.Senders))
	for i, originalSender := range normConn.Senders {
		desugaredFanIn = append(desugaredFanIn, src.Connection{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{originalSender},
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: fanInNodeName,
							Port: "data",
							Idx:  compiler.Pointer(uint8(i)),
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Meta: locOnlyMeta,
			},
			Meta: locOnlyMeta,
		})
	}

	// 3. insert new connection: fan_in -> original receivers
	desugaredFanIn = append(desugaredFanIn, src.Connection{
		Normal: &src.NormalConnection{
			Senders: []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Node: fanInNodeName,
						Port: "res",
						Meta: locOnlyMeta,
					},
					Meta: locOnlyMeta,
				},
			},
			Receivers: normConn.Receivers,
			Meta:      locOnlyMeta,
		},
		Meta: locOnlyMeta,
	})

	// 4. desugar each connection (original senders and receivers might need it)
	desugaredConnections, err := d.desugarConnections(
		iface,
		desugaredFanIn,
		nodePortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return nil, fmt.Errorf("desugar fan in: %w", err)
	}

	return desugaredConnections, nil
}

type handleTernarySenderResult struct {
	replace src.Connection
	insert  []src.Connection
}

// (cond ? left : right) -> XXX;
// =>
// 1) cond -> ternary:if;
// 2) left -> ternary:then;
// 3) right -> ternary:else;
// 4) ternary:res -> XXX;
func (d *Desugarer) desugarTernarySender(
	iface src.Interface,
	ternary src.Ternary,
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	usedNodeOutports nodeOutportsUsed,
	scope Scope,
	nodes map[string]src.Node,
) (handleTernarySenderResult, error) {
	locOnlyMeta := core.Meta{Location: ternary.Meta.Location}

	d.ternaryCounter++
	ternaryNodeName := fmt.Sprintf("__ternary__%d", d.ternaryCounter)

	nodesToInsert[ternaryNodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "Ternary",
			Meta: locOnlyMeta,
		},
		Meta: locOnlyMeta,
	}

	sugaredInsert := []src.Connection{
		// 1) cond -> ternary:if
		{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{ternary.Condition},
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: ternaryNodeName,
							Port: "if",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
			},
			Meta: locOnlyMeta,
		},
		// 2) left -> ternary:then
		{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{ternary.Left},
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: ternaryNodeName,
							Port: "then",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Meta: locOnlyMeta,
			},
		},
		// right -> ternary:else
		{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{ternary.Right},
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: ternaryNodeName,
							Port: "else",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Meta: locOnlyMeta,
			},
			Meta: locOnlyMeta,
		},
	}

	desugaredInsert := make([]src.Connection, 0, len(sugaredInsert))
	for _, conn := range sugaredInsert {
		desugarConnRes, err := d.desugarConnection(
			iface,
			conn,
			usedNodeOutports,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return handleTernarySenderResult{}, err
		}
		desugaredInsert = append(desugaredInsert, *desugarConnRes.replace)
		desugaredInsert = append(desugaredInsert, desugarConnRes.insert...)
	}

	// 4) ternary:res -> XXX;
	sugaredReplace := src.Connection{
		Normal: &src.NormalConnection{
			Senders: []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Node: ternaryNodeName,
						Port: "res",
						Meta: locOnlyMeta,
					},
				},
			},
			Receivers: normConn.Receivers,
			Meta:      locOnlyMeta,
		},
	}

	return handleTernarySenderResult{
		replace: sugaredReplace,
		insert:  desugaredInsert,
	}, nil
}

type handleBinarySenderResult struct {
	replace src.Connection
	insert  []src.Connection
}

func (d *Desugarer) desugarBinarySender(
	iface src.Interface,
	binary src.Binary,
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	usedNodeOutports nodeOutportsUsed,
	scope Scope,
	nodes map[string]src.Node,
) (handleBinarySenderResult, error) {
	locOnlyMeta := core.Meta{
		Location: binary.Meta.Location,
	}

	var (
		opNode      string
		opComponent string
	)

	switch binary.Operator {
	// Arithmetic
	case src.AddOp:
		d.addCounter++
		opNode = fmt.Sprintf("__add__%d", d.addCounter)
		opComponent = "Add"
	case src.SubOp:
		d.subCounter++
		opNode = fmt.Sprintf("__sub__%d", d.subCounter)
		opComponent = "Sub"
	case src.MulOp:
		d.mulCounter++
		opNode = fmt.Sprintf("__mul__%d", d.mulCounter)
		opComponent = "Mul"
	case src.DivOp:
		d.divCounter++
		opNode = fmt.Sprintf("__div__%d", d.divCounter)
		opComponent = "Div"
	case src.ModOp:
		d.modCounter++
		opNode = fmt.Sprintf("__mod__%d", d.modCounter)
		opComponent = "Mod"
	case src.PowOp:
		d.powCounter++
		opNode = fmt.Sprintf("__pow__%d", d.powCounter)
		opComponent = "Pow"
	// Comparison
	case src.EqOp:
		d.eqCounter++
		opNode = fmt.Sprintf("__eq__%d", d.eqCounter)
		opComponent = "Eq"
	case src.NeOp:
		d.neCounter++
		opNode = fmt.Sprintf("__ne__%d", d.neCounter)
		opComponent = "Ne"
	case src.GtOp:
		d.gtCounter++
		opNode = fmt.Sprintf("__gt__%d", d.gtCounter)
		opComponent = "Gt"
	case src.LtOp:
		d.ltCounter++
		opNode = fmt.Sprintf("__lt__%d", d.ltCounter)
		opComponent = "Lt"
	case src.GeOp:
		d.geCounter++
		opNode = fmt.Sprintf("__ge__%d", d.geCounter)
		opComponent = "Ge"
	case src.LeOp:
		d.leCounter++
		opNode = fmt.Sprintf("__le__%d", d.leCounter)
		opComponent = "Le"
	// Logical
	case src.AndOp:
		d.andCounter++
		opNode = fmt.Sprintf("__and__%d", d.andCounter)
		opComponent = "And"
	case src.OrOp:
		d.orCounter++
		opNode = fmt.Sprintf("__or__%d", d.orCounter)
		opComponent = "Or"
	// Bitwise
	case src.BitAndOp:
		d.bitAndCounter++
		opNode = fmt.Sprintf("__bitAnd__%d", d.bitAndCounter)
		opComponent = "BitAnd"
	case src.BitOrOp:
		d.bitOrCounter++
		opNode = fmt.Sprintf("__bitOr__%d", d.bitOrCounter)
		opComponent = "BitOr"
	case src.BitXorOp:
		d.bitXorCounter++
		opNode = fmt.Sprintf("__bitXor__%d", d.bitXorCounter)
		opComponent = "BitXor"
	case src.BitLshOp:
		d.bitLshCounter++
		opNode = fmt.Sprintf("__bitLsh__%d", d.bitLshCounter)
		opComponent = "BitLsh"
	case src.BitRshOp:
		d.bitRshCounter++
		opNode = fmt.Sprintf("__bitRsh__%d", d.bitRshCounter)
		opComponent = "BitRsh"
	default:
		return handleBinarySenderResult{}, fmt.Errorf(
			"unsupported binary operator: %s",
			binary.Operator,
		)
	}

	operatorNodeToInsert := src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: opComponent,
			Meta: locOnlyMeta,
		},
		OverloadIndex: nil, // To be set later
		Meta:          locOnlyMeta,
	}

	// set overload index for operator node based on analyzed operand type
	// use analyzer-provided binary.AnalyzedType to choose overload version
	if binary.AnalyzedType.Inst == nil {
		panic("binary analyzed type must be instantiation at desugaring stage")
	}

	opEntity, _, err := scope.Entity(core.EntityRef{Pkg: "builtin", Name: opComponent})
	if err != nil {
		panic(fmt.Sprintf("resolve operator entity: %v", err))
	}
	if opEntity.Kind != src.ComponentEntity {
		panic("operator entity must be a component")
	}
	if len(opEntity.Component) == 0 {
		panic("operator entity must have at least one component")
	}

	// scope = scope.Relocate(opEntityLoc) // mainly for completeness (not sure if needed)

	if len(opEntity.Component) == 1 {
		operatorNodeToInsert.OverloadIndex = compiler.Pointer(int(0))
	} else {
		wantName := binary.AnalyzedType.Inst.Ref.Name
		for i, version := range opEntity.Component {
			leftPort, okL := version.Interface.IO.In["left"]
			rightPort, okR := version.Interface.IO.In["right"]
			if !okL || !okR {
				continue
			}
			if leftPort.TypeExpr.Inst == nil || rightPort.TypeExpr.Inst == nil {
				continue
			}
			if leftPort.TypeExpr.Inst.Ref.Name == wantName && rightPort.TypeExpr.Inst.Ref.Name == wantName {
				operatorNodeToInsert.OverloadIndex = compiler.Pointer(i)
				break
			}
		}
	}

	if operatorNodeToInsert.OverloadIndex == nil {
		panic(fmt.Sprintf("no matching overload for operator %s and type %s", opComponent, binary.AnalyzedType.String()))
	}

	nodesToInsert[opNode] = operatorNodeToInsert

	// left -> op:left
	// right -> op:right
	sugaredInsert := []src.Connection{
		{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{binary.Left},
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: opNode,
							Port: "left",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Meta: locOnlyMeta,
			},
		},
		{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{binary.Right},
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: opNode,
							Port: "right",
							Meta: locOnlyMeta,
						},
						Meta: locOnlyMeta,
					},
				},
				Meta: locOnlyMeta,
			},
		},
	}

	// operand-senders might be sugared, so we need to desugar them
	desugaredInsert := make([]src.Connection, 0, len(sugaredInsert))
	for _, conn := range sugaredInsert {
		desugarConnRes, err := d.desugarConnection(
			iface,
			conn,
			usedNodeOutports,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return handleBinarySenderResult{}, err
		}
		desugaredInsert = append(desugaredInsert, *desugarConnRes.replace)
		desugaredInsert = append(desugaredInsert, desugarConnRes.insert...)
	}

	// op:res -> XXX
	replace := src.Connection{
		Normal: &src.NormalConnection{
			Senders: []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Node: opNode,
						Port: "res",
						Meta: locOnlyMeta,
					},
					Meta: locOnlyMeta,
				},
			},
			Receivers: normConn.Receivers, // desugaring of original receivers is job of caller
			Meta:      locOnlyMeta,
		},
		Meta: locOnlyMeta,
	}

	return handleBinarySenderResult{
		replace: replace,
		insert:  desugaredInsert,
	}, nil
}

package desugarer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

const maxUint8 = int(^uint8(0))

//nolint:govet // fieldalignment: keep semantic grouping.
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

	// Implicit fan-in is generally forbidden by analyzer, but err-guard (`?`)
	// desugaring can introduce it by wiring multiple errors to `:err`.
	// Merge those into a single multi-sender connection so FanIn is inserted.
	net = d.mergeImplicitFanIn(net)

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

// mergeImplicitFanIn merges multiple connections targeting the same outport receiver.
// Implicit fan-in is normally rejected by analyzer; currently `?` is the only
// desugaring that can create it (multiple `err` senders to `:err`).
// By merging into a single multi-sender connection, we reuse explicit fan-in
// desugaring to insert FanIn nodes.
func (d *Desugarer) mergeImplicitFanIn(
	net []src.Connection,
) []src.Connection {
	type receiverKey struct {
		node   string
		port   string
		hasIdx bool
		idx    uint8
	}
	//nolint:govet // fieldalignment: local helper layout.
	type group struct {
		receiver src.ConnectionReceiver
		senders  []src.ConnectionSender
		meta     core.Meta
	}

	groups := map[receiverKey]group{}
	order := make([]receiverKey, 0, len(net))
	positions := map[receiverKey]int{}
	kept := make([]src.Connection, 0, len(net))

	for _, conn := range net {
		if conn.ArrayBypass != nil || conn.Normal == nil {
			kept = append(kept, conn)
			continue
		}

		norm := conn.Normal
		if len(norm.Receivers) != 1 {
			kept = append(kept, conn)
			continue
		}

		receiver := norm.Receivers[0]
		if receiver.PortAddr == nil || receiver.ChainedConnection != nil {
			kept = append(kept, conn)
			continue
		}

		if receiver.PortAddr.Node != "out" {
			kept = append(kept, conn)
			continue
		}

		key := receiverKey{
			node: receiver.PortAddr.Node,
			port: receiver.PortAddr.Port,
		}
		if receiver.PortAddr.Idx != nil {
			key.hasIdx = true
			key.idx = *receiver.PortAddr.Idx
		}

		current, ok := groups[key]
		if !ok {
			current = group{
				receiver: receiver,
				meta:     core.Meta{Location: receiver.PortAddr.Meta.Location},
			}
			order = append(order, key)
			positions[key] = len(kept)
			kept = append(kept, src.Connection{})
		}
		current.senders = append(current.senders, norm.Senders...)
		groups[key] = current
	}

	for _, key := range order {
		g := groups[key]
		kept[positions[key]] = src.Connection{
			Normal: &src.NormalConnection{
				Senders:   g.senders,
				Receivers: []src.ConnectionReceiver{g.receiver},
				Meta:      g.meta,
			},
			Meta: g.meta,
		}
	}

	return kept
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
	case chainHead.Const != nil:
		chainHeadPort = "sig" // New has sig inport
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
					Name: "New",
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
					Name: "New",
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
		TypeArgs: []ts.Expr{
			{
				Inst: &ts.InstExpr{
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
		result := d.desugarStructSelectors(
			normConn,
			nodesToInsert,
			constsToInsert,
		)

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
			portAddr := d.handleLiteralSender(*sender.Const, nodesToInsert, constsToInsert)

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

	return desugarSenderResult{}, fmt.Errorf("unexpected sender type: %v", sender.Meta.Location)
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

func newComponentRef() core.EntityRef {
	return core.EntityRef{
		Pkg:  "builtin",
		Name: "New",
	}
}

func (d *Desugarer) handleLiteralSender(
	constant src.Const,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) src.PortAddr {
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
			Pkg:  newComponentRef().Pkg,
			Name: newComponentRef().Name,
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

	return emitterNodeOutportAddr
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
			Pkg:  newComponentRef().Pkg,
			Name: newComponentRef().Name,
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
		if i > maxUint8 {
			return desugarFanOutResult{}, fmt.Errorf("fan-out index %d overflows uint8", i)
		}
		conn := src.Connection{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{
							Node: nodeName,
							Port: "data",
							Idx:  compiler.Pointer(uint8(i)), // #nosec G115 -- bounds checked above
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
		if i > maxUint8 {
			return nil, fmt.Errorf("fan-in index %d overflows uint8", i)
		}
		desugaredFanIn = append(desugaredFanIn, src.Connection{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{originalSender},
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: fanInNodeName,
							Port: "data",
							Idx:  compiler.Pointer(uint8(i)), // #nosec G115 -- bounds checked above
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

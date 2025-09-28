package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	"github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type foundInterface struct {
	iface    src.Interface
	location core.Location
}

func (a Analyzer) analyzeNodes(
	iface src.Interface, // resolved interface of the component that contains the nodes
	nodes map[string]src.Node, // nodes to analyze
	net []src.Connection, // network of the component that contains the nodes
	scope src.Scope, // scope of the component
) (
	map[string]src.Node, // resolved nodes
	map[string]foundInterface, // resolved nodes interfaces with locations
	bool, // one of the nodes has error guard
	*compiler.Error, // err
) {
	analyzedNodes := make(map[string]src.Node, len(nodes))
	nodesInterfaces := make(map[string]foundInterface, len(nodes))
	hasErrGuard := false

	for nodeName, node := range nodes {
		if node.ErrGuard {
			hasErrGuard = true
		}

		analyzedNode, nodeInterface, err := a.analyzeNode(
			nodeName,
			node,
			scope,
			iface,
			nodes,
			net,
		)
		if err != nil {
			return nil, nil, false, compiler.Error{
				Meta: &node.Meta,
			}.Wrap(err)
		}

		nodesInterfaces[nodeName] = nodeInterface
		analyzedNodes[nodeName] = analyzedNode
	}

	return analyzedNodes, nodesInterfaces, hasErrGuard, nil
}

func (a Analyzer) analyzeNode(
	name string, // name of the node
	node src.Node, // node to analyze
	scope src.Scope, // scope of the component that contains the node
	iface src.Interface, // interface of the component that contains the node
	nodes map[string]src.Node, // nodes of the component that contains the node
	net []src.Connection, // network of the component that contains the node
) (src.Node, foundInterface, *compiler.Error) {
	parentTypeParams := iface.TypeParams

	nodeEntity, location, err := scope.Entity(node.EntityRef)
	if err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &node.Meta,
		}
	}

	if nodeEntity.Kind != src.ComponentEntity &&
		nodeEntity.Kind != src.InterfaceEntity {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: fmt.Sprintf("Node can only refer to flows or interfaces: %v", nodeEntity.Kind),
			Meta:    nodeEntity.Meta(),
		}
	}

	bindArg, hasBind := node.Directives[compiler.BindDirective]
	if hasBind && len(bindArg) != 1 {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: "Node with #bind directive must provide exactly one argument",
			Meta:    nodeEntity.Meta(),
		}
	}

	if hasBind && nodeEntity.Kind == src.InterfaceEntity {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: "Interface node cannot use #bind directive",
			Meta:    nodeEntity.Meta(),
		}
	}

	if nodeEntity.Kind == src.InterfaceEntity && node.DIArgs != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: "Only component node can have dependency injection",
			Meta:    nodeEntity.Meta(),
		}
	}

	// We need to get resolved frame from parent type parameters
	// in order to be able to resolve node's args
	// since they can refer to type parameter of the parent (interface)
	_, resolvedParentParamsFrame, err := a.resolver.ResolveParams(
		parentTypeParams.Params,
		scope,
	)
	if err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &node.Meta,
		}
	}

	// Now when we have frame made of parent type parameters constraints
	// we can resolve cases like `subnode SubFlow<T>`
	// where `T` refers to type parameter of the component/interface we're in.
	resolvedNodeArgs, err := a.resolver.ResolveExprsWithFrame(
		node.TypeArgs,
		resolvedParentParamsFrame,
		scope,
	)
	if err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &node.Meta,
		}
	}

	var (
		nodeIface     src.Interface
		overloadIndex *int // only for component nodes
	)
	if nodeEntity.Kind == src.InterfaceEntity {
		nodeIface = nodeEntity.Interface
	} else {
		var err *compiler.Error
		nodeIface, overloadIndex, err = a.getInterfaceAndOverloadingIndexForNode(
			name,
			nodeEntity,
			hasBind,
			node,
			scope,
			resolvedNodeArgs,
			iface,
			nodes,
			net,
		)
		if err != nil {
			return src.Node{}, foundInterface{}, err
		}
	}

	if node.ErrGuard {
		if _, ok := iface.IO.Out["err"]; !ok {
			return src.Node{}, foundInterface{}, &compiler.Error{
				Message: "Error-guard operator '?' can only be used in components with ':err' outport to propagate errors",
				Meta:    &node.Meta,
			}
		}
		if _, ok := nodeIface.IO.Out["err"]; !ok {
			return src.Node{}, foundInterface{}, &compiler.Error{
				Message: "Error-guard operator '?' requires node to have ':err' outport to propagate errors",
				Meta:    &node.Meta,
			}
		}
	}

	// default any
	if len(resolvedNodeArgs) == 0 && len(nodeIface.TypeParams.Params) == 1 {
		resolvedNodeArgs = []typesystem.Expr{
			{
				Inst: &typesystem.InstExpr{
					Ref: core.EntityRef{Name: "any"},
				},
			},
		}
	}

	// Finally check that every argument is compatible
	// with corresponding parameter of the node's interface.
	if err = a.resolver.CheckArgsCompatibility(
		resolvedNodeArgs,
		nodeIface.TypeParams.Params,
		scope,
	); err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: fmt.Sprintf("%s: %s", node.EntityRef.Name, err.Error()),
			Meta:    &node.Meta,
		}
	}

	if node.DIArgs == nil {
		return src.Node{
				Directives:    node.Directives,
				EntityRef:     node.EntityRef,
				TypeArgs:      resolvedNodeArgs,
				Meta:          node.Meta,
				OverloadIndex: overloadIndex,
				ErrGuard:      node.ErrGuard,
			}, foundInterface{
				iface:    nodeIface,
				location: location,
			}, nil
	}

	resolvedFlowDI := make(map[string]src.Node, len(node.DIArgs))
	for depName, depNode := range node.DIArgs {
		resolvedDep, _, err := a.analyzeNode(
			name, // TODO make sure DI works with overloading (example: Reduce{Add})
			depNode,
			scope,
			iface,
			nodes,
			net,
		)
		if err != nil {
			return src.Node{}, foundInterface{}, compiler.Error{
				Meta: &depNode.Meta,
			}.Wrap(err)
		}
		resolvedFlowDI[depName] = resolvedDep
	}

	return src.Node{
			Directives:    node.Directives,
			EntityRef:     node.EntityRef,
			TypeArgs:      resolvedNodeArgs,
			DIArgs:        resolvedFlowDI,
			ErrGuard:      node.ErrGuard,
			OverloadIndex: overloadIndex,
			Meta:          node.Meta,
		}, foundInterface{
			iface:    nodeIface,
			location: location,
		}, nil
}

// getInterfaceAndOverloadingIndexForNode returns interface and overload index for the given node.
// Overloading at the level of sourcecode is implemented here.
func (a Analyzer) getInterfaceAndOverloadingIndexForNode(
	nodeName string,
	entity src.Entity,
	hasBind bool,
	node src.Node,
	scope src.Scope,
	resolvedNodeArgs []typesystem.Expr,
	resolvedParentIface src.Interface, // resolved interface of the component that contains the node
	allParentNodes map[string]src.Node, // nodes of the component that contains the node
	net []src.Connection, // network of the component that contains the node
) (src.Interface, *int, *compiler.Error) {
	var (
		overloadIndex *int
		version       src.Component
	)
	if len(entity.Component) == 1 {
		version = entity.Component[0]
	} else {
		var err *compiler.Error
		version, overloadIndex, err = a.getNodeOverloadVersionAndIndex(
			nodeName,
			resolvedParentIface,
			scope,
			resolvedNodeArgs,
			allParentNodes,
			net,
			entity.Component,
		)
		if err != nil {
			return src.Interface{}, nil, &compiler.Error{
				Message: "Node can't use #bind if it isn't instantiated with the component that use #extern",
				Meta:    &node.Meta,
			}
		}
	}

	_, hasExtern := version.Directives[compiler.ExternDirective]
	if hasBind && !hasExtern {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Node can't use #bind if it isn't instantiated with the component that use #extern",
			Meta:    entity.Meta(),
		}
	}

	versionIface := version.Interface

	_, hasAutoPortsDirective := version.Directives[compiler.AutoportsDirective]
	if !hasAutoPortsDirective {
		return versionIface, overloadIndex, nil
	}

	// if we here then we have #autoports (only for structs)

	if len(versionIface.IO.In) != 0 {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Component that uses struct inports directive must have no defined inports",
			Meta:    entity.Meta(),
		}
	}

	if len(versionIface.TypeParams.Params) != 1 {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Exactly one type parameter expected",
			Meta:    entity.Meta(),
		}
	}

	resolvedTypeParamConstr, err := a.resolver.ResolveExpr(versionIface.TypeParams.Params[0].Constr, scope)
	if err != nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: err.Error(),
			Meta:    entity.Meta(),
		}
	}

	if resolvedTypeParamConstr.Lit == nil || resolvedTypeParamConstr.Lit.Struct == nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Struct type expected",
			Meta:    entity.Meta(),
		}
	}

	if len(resolvedNodeArgs) != 1 {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Exactly one type argument expected",
			Meta:    entity.Meta(),
		}
	}

	resolvedNodeArg, err := a.resolver.ResolveExpr(resolvedNodeArgs[0], scope)
	if err != nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: err.Error(),
			Meta:    entity.Meta(),
		}
	}

	if resolvedNodeArg.Lit == nil || resolvedNodeArg.Lit.Struct == nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Struct argument expected",
			Meta:    entity.Meta(),
		}
	}

	structFields := resolvedNodeArg.Lit.Struct
	inports := make(map[string]src.Port, len(structFields))
	for fieldName, fieldTypeExpr := range structFields {
		inports[fieldName] = src.Port{
			TypeExpr: fieldTypeExpr,
		}
	}

	return src.Interface{
		TypeParams: versionIface.TypeParams,
		IO: src.IO{
			In: inports,
			Out: map[string]src.Port{
				// struct builder has exactly one outport - created structure
				"res": {
					TypeExpr: resolvedNodeArg,
					IsArray:  false,
					Meta:     versionIface.IO.Out["res"].Meta,
				},
			},
		},
		Meta: versionIface.Meta,
	}, overloadIndex, nil
}

// getNodeOverloadVersionAndIndex determines which overload of a component to use for a node.
// This is called when we know the node references an overloaded component with multiple implementations.
// It analyzes how the node is used in connections to determine the appropriate implementation.
func (a Analyzer) getNodeOverloadVersionAndIndex(
	nodeName string,
	resolvedParentIface src.Interface,
	scope src.Scope,
	resolvedNodeArgs []typesystem.Expr,
	allParentNodes map[string]src.Node,
	net []src.Connection,
	allNodeComponentVersions []src.Component,
) (src.Component, *int, *compiler.Error) {
	nodeRefs := findNodeRefsInNet(nodeName, net)
	if len(nodeRefs) == 0 {
		return src.Component{}, nil, &compiler.Error{
			Message: fmt.Sprintf("no usages found for node %s", nodeName),
			Meta:    &resolvedParentIface.Meta,
		}
	}

	// nodeConstraints := a.collectUsageDerivedTypeConstraintsForNode(
	// 	nodeName,
	// 	resolvedParentIface,
	// 	scope,
	// 	allParentNodes,
	// 	net,
	// )

	var (
		remainingIdx   []int
		remainingComps []src.Component
	)
	for i, component := range allNodeComponentVersions {
		if !isCandidateCompatibleWithAllNodeRefs(component, nodeRefs) {
			continue
		}
		// TODO: uncomment after the related code will be carefully reviewed!
		// if !a.doesCandidateSatisfyTypeConstraints(
		// 	component.Interface,
		// 	resolvedNodeArgs,
		// 	nodeConstraints,
		// 	scope,
		// ) {
		// 	continue
		// }
		remainingIdx = append(remainingIdx, i)
		remainingComps = append(remainingComps, component)
	}

	if len(remainingComps) == 0 {
		return src.Component{}, nil, &compiler.Error{
			Message: fmt.Sprintf("no compatible overload found for node %s", nodeName),
			Meta:    &resolvedParentIface.Meta,
		}
	}

	if len(remainingComps) > 1 {
		return src.Component{}, nil, &compiler.Error{
			Message: fmt.Sprintf("ambiguous overload for node %s: multiple candidates satisfy usage", nodeName),
			Meta:    &resolvedParentIface.Meta,
		}
	}

	return remainingComps[0], &remainingIdx[0], nil
}

// nodeUsageConstraints captures incoming produced types and outgoing expected types per port.
// type nodeUsageConstraints struct {
// 	incoming map[string][]typesystem.Expr // for inports: types produced by connected senders
// 	outgoing map[string][]typesystem.Expr // for outports: types expected by connected receivers
// }

// collectUsageDerivedTypeConstraintsForNode inspects the network and extracts type constraints for the given node.
// it only uses available information (parent iface, literals, neighbor node interfaces). it does not depend on nodesIfaces.
// func (a Analyzer) collectUsageDerivedTypeConstraintsForNode(
// 	nodeName string,
// 	resolvedParentIface src.Interface,
// 	scope src.Scope,
// 	nodes map[string]src.Node,
// 	net []src.Connection,
// ) (nodeUsageConstraints) {
// 	c := nodeUsageConstraints{
// 		incoming: map[string][]typesystem.Expr{},
// 		outgoing: map[string][]typesystem.Expr{},
// 	}

// 	// helper to append unique types (by String) to a slice
// 	appendUnique := func(dst *[]typesystem.Expr, t typesystem.Expr) {
// 		s := t.String()
// 		for _, e := range *dst {
// 			if e.String() == s {
// 				return
// 			}
// 		}
// 		*dst = append(*dst, t)
// 	}

// 	// resolve parent param frame once
// 	_, parentFrame, err := a.resolver.ResolveParams(
// 		resolvedParentIface.TypeParams.Params, scope,
// 	)
// 	if err != nil {
// 		return c
// 	}

// 	// walk all connections
// 	for _, conn := range net {
// 		if conn.Normal == nil {
// 			continue
// 		}

// 		// check if our node is a sender in this connection
// 		for _, sender := range conn.Normal.Senders {
// 			if sender.PortAddr == nil || sender.PortAddr.Node != nodeName {
// 				continue
// 			}
// 			port := sender.PortAddr.Port
// 			// derive expected types from all receivers
// 			recvPortAddrs := a.flattenReceiversPortAddrs(conn.Normal.Receivers)
// 			for _, rpa := range recvPortAddrs {
// 				// parent out
// 				if rpa.Node == "out" {
// 					if p, ok := resolvedParentIface.IO.Out[rpa.Port]; ok {
// 						if resolved, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, parentFrame, scope); err == nil {
// 							list := c.outgoing[port]
// 							appendUnique(&list, resolved)
// 							c.outgoing[port] = list
// 						}
// 					}
// 					continue
// 				}
// 				// other node inport
// 				recvNode, ok := nodes[rpa.Node]
// 				if !ok {
// 					continue
// 				}
// 				// get possible inport types across overloads
// 				for _, t := range a.getPossibleNodePortTypes(scope, parentFrame, recvNode, true, rpa.Port) {
// 					list := c.outgoing[port]
// 					appendUnique(&list, t)
// 					c.outgoing[port] = list
// 				}
// 			}
// 		}

// 		// check if our node is a receiver in this connection
// 		recvPortAddrs := a.flattenReceiversPortAddrs(conn.Normal.Receivers)
// 		for _, rpa := range recvPortAddrs {
// 			if rpa.Node != nodeName {
// 				continue
// 			}
// 			port := rpa.Port
// 			// derive produced types from all senders
// 			for _, sender := range conn.Normal.Senders {
// 				for _, t := range a.getPossibleSenderTypes(scope, parentFrame, resolvedParentIface, nodes, sender) {
// 					list := c.incoming[port]
// 					appendUnique(&list, t)
// 					c.incoming[port] = list
// 				}
// 			}
// 		}
// 	}

// 	return c
// }

// flattenReceiversPortAddrs extracts all direct receiver port addresses from potentially nested receivers.
// func (a Analyzer) flattenReceiversPortAddrs(receivers []src.ConnectionReceiver) []src.PortAddr {
// 	var res []src.PortAddr
// 	var visit func(recs []src.ConnectionReceiver)
// 	visit = func(recs []src.ConnectionReceiver) {
// 		for _, r := range recs {
// 			if r.PortAddr != nil {
// 				res = append(res, *r.PortAddr)
// 			}
// 			if r.ChainedConnection != nil && r.ChainedConnection.Normal != nil {
// 				// only the head receiver is relevant as a consumer of our sender
// 				visit(r.ChainedConnection.Normal.Receivers)
// 			}
// 			if r.DeferredConnection != nil && r.DeferredConnection.Normal != nil {
// 				visit(r.DeferredConnection.Normal.Receivers)
// 			}
// 			if r.Switch != nil {
// 				for _, cse := range r.Switch.Cases {
// 					visit(cse.Receivers)
// 				}
// 				if r.Switch.Default != nil {
// 					visit(r.Switch.Default)
// 				}
// 			}
// 		}
// 	}
// 	visit(receivers)
// 	return res
// }

// getPossibleSenderTypes returns a set of possible types produced by a sender without requiring resolved node interfaces.
// func (a Analyzer) getPossibleSenderTypes(
// 	scope src.Scope,
// 	parentFrame map[string]typesystem.Def,
// 	parentIface src.Interface,
// 	nodes map[string]src.Node,
// 	sender src.ConnectionSender,
// ) []typesystem.Expr {
// 	// const sender
// 	if sender.Const != nil {
// 		if _, t, err := a.getConstSenderType(*sender.Const, scope); err == nil {
// 			return []typesystem.Expr{t}
// 		}
// 	}
// 	// range sender: stream<int>
// 	if sender.Range != nil {
// 		return []typesystem.Expr{
// 			{
// 				Inst: &typesystem.InstExpr{
// 					Ref:  core.EntityRef{Name: "stream"},
// 					Args: []typesystem.Expr{{Inst: &typesystem.InstExpr{Ref: core.EntityRef{Name: "int"}}}},
// 				},
// 			},
// 		}
// 	}
// 	// union sender produces the union type itself
// 	if sender.Union != nil {
// 		if entity, _, err := scope.GetType(sender.Union.EntityRef); err == nil {
// 			if t, e := a.analyzeTypeExpr(*entity.BodyExpr, scope); e == nil {
// 				return []typesystem.Expr{t}
// 			}
// 		}
// 	}
// 	// binary sender: approximate using left operand type and operator rule
// 	if sender.Binary != nil {
// 		lefts := a.getPossibleSenderTypes(scope, parentFrame, parentIface, nodes, sender.Binary.Left)
// 		if len(lefts) > 0 {
// 			return []typesystem.Expr{a.getBinaryExprType(sender.Binary.Operator, lefts[0])}
// 		}
// 	}
// 	// ternary: approximate with left branch
// 	if sender.Ternary != nil {
// 		lefts := a.getPossibleSenderTypes(scope, parentFrame, parentIface, nodes, sender.Ternary.Left)
// 		if len(lefts) > 0 {
// 			return []typesystem.Expr{lefts[0]}
// 		}
// 	}
// 	// port-addr
// 	if sender.PortAddr != nil {
// 		pa := *sender.PortAddr
// 		if pa.Node == "in" {
// 			if p, ok := parentIface.IO.In[pa.Port]; ok {
// 				if resolved, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, parentFrame, scope); err == nil {
// 					return []typesystem.Expr{resolved}
// 				}
// 			}
// 			return nil
// 		}
// 		// outport of another node
// 		other, ok := nodes[pa.Node]
// 		if !ok {
// 			return nil
// 		}
// 		return a.getPossibleNodePortTypes(scope, parentFrame, other, false, pa.Port)
// 	}
// 	return nil
// }

// getPossibleNodePortTypes collects possible port types across all overloads of a node's component.
// isInput determines whether we look at inports (true) or outports (false).
// func (a Analyzer) getPossibleNodePortTypes(
// 	scope src.Scope,
// 	parentFrame map[string]typesystem.Def,
// 	node src.Node,
// 	isInput bool,
// 	portName string,
// ) []typesystem.Expr {
// 	var out []typesystem.Expr
// 	// resolve node's entity
// 	entity, _, err := scope.Entity(node.EntityRef)
// 	if err != nil || entity.Kind != src.ComponentEntity {
// 		return out
// 	}
// 	// resolve node type args against parent frame
// 	resolvedArgs, err2 := a.resolver.ResolveExprsWithFrame(node.TypeArgs, parentFrame, scope)
// 	if err2 != nil {
// 		return out
// 	}
// 	for _, comp := range entity.Component {
// 		iface := comp.Interface
// 		// choose port
// 		var p src.Port
// 		var ok bool
// 		if isInput {
// 			if portName == "" {
// 				if len(iface.IO.In) == 1 {
// 					for _, v := range iface.IO.In {
// 						p = v
// 						ok = true
// 						break
// 					}
// 				}
// 			} else {
// 				p, ok = iface.IO.In[portName]
// 			}
// 		} else {
// 			if portName == "" {
// 				// if multiple, prefer first non-err
// 				if len(iface.IO.Out) == 1 {
// 					for _, v := range iface.IO.Out {
// 						p = v
// 						ok = true
// 						break
// 					}
// 				} else {
// 					for name, v := range iface.IO.Out {
// 						if name != "err" {
// 							p = v
// 							ok = true
// 							break
// 						}
// 					}
// 				}
// 			} else {
// 				p, ok = iface.IO.Out[portName]
// 			}
// 		}
// 		if !ok {
// 			continue
// 		}
// 		// substitute node args into port type
// 		frame := make(map[string]typesystem.Def, len(iface.TypeParams.Params))
// 		for i, param := range iface.TypeParams.Params {
// 			if i < len(resolvedArgs) {
// 				frame[param.Name] = typesystem.Def{BodyExpr: &resolvedArgs[i]}
// 			}
// 		}
// 		if resolved, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, frame, scope); err == nil {
// 			out = append(out, resolved)
// 		}
// 	}
// 	return out
// }

// // doesCandidateSatisfyTypeConstraints checks that a candidate interface matches all collected constraints.
// func (a Analyzer) doesCandidateSatisfyTypeConstraints(
// 	candidate src.Interface,
// 	resolvedNodeArgs []typesystem.Expr,
// 	nodeUsageConstr nodeUsageConstraints,
// 	scope src.Scope,
// ) bool {
// 	// build frame for candidate from node's resolved args
// 	frame := make(map[string]typesystem.Def, len(candidate.TypeParams.Params))
// 	for i, param := range candidate.TypeParams.Params {
// 		if i < len(resolvedNodeArgs) {
// 			frame[param.Name] = typesystem.Def{BodyExpr: &resolvedNodeArgs[i]}
// 		}
// 	}

// 	// check incoming constraints
// 	for port, types := range nodeUsageConstr.incoming {
// 		p, ok := candidate.IO.In[port]
// 		if !ok {
// 			return false
// 		}
// 		candType, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, frame, scope)
// 		if err != nil {
// 			return false
// 		}
// 		for _, t := range types {
// 			if err := a.resolver.IsSubtypeOf(t, candType, scope); err != nil {
// 				return false
// 			}
// 		}
// 	}

// 	// check outgoing constraints
// 	for port, types := range nodeUsageConstr.outgoing {
// 		p, ok := candidate.IO.Out[port]
// 		if !ok {
// 			return false
// 		}
// 		candType, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, frame, scope)
// 		if err != nil {
// 			return false
// 		}
// 		for _, exp := range types {
// 			if err := a.resolver.IsSubtypeOf(candType, exp, scope); err != nil {
// 				return false
// 			}
// 		}
// 	}

// 	return true
// }

// findNodeRefsInNet identifies all places where the specified node is referenced in a connection.
func findNodeRefsInNet(nodeName string, connections []src.Connection) []nodeRefInNet {
	var refs []nodeRefInNet

	for _, conn := range connections {
		if conn.ArrayBypass != nil {
			if conn.ArrayBypass.SenderOutport.Node == nodeName {
				refs = append(refs, nodeRefInNet{
					isOutgoing: true,
					port:       conn.ArrayBypass.SenderOutport.Port,
					arrayIdx:   conn.ArrayBypass.SenderOutport.Idx,
				})
			}
			if conn.ArrayBypass.ReceiverInport.Node == nodeName {
				refs = append(refs, nodeRefInNet{
					isOutgoing: false,
					port:       conn.ArrayBypass.ReceiverInport.Port,
					arrayIdx:   conn.ArrayBypass.ReceiverInport.Idx,
				})
			}
			continue
		}

		if conn.Normal != nil {
			for _, sender := range conn.Normal.Senders {
				if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
					refs = append(refs, nodeRefInNet{
						isOutgoing: true,
						port:       sender.PortAddr.Port,
						arrayIdx:   sender.PortAddr.Idx,
					})
				}
			}
			refs = append(refs, findNodeUsagesInReceivers(nodeName, conn.Normal.Receivers)...)
		}
	}

	return refs
}

type nodeRefInNet struct {
	isOutgoing bool
	port       string
	arrayIdx   *uint8
}

// findNodeUsagesInReceivers recursively checks receivers for node usages
func findNodeUsagesInReceivers(nodeName string, receivers []src.ConnectionReceiver) []nodeRefInNet {
	var nodeRefs []nodeRefInNet

	for _, receiver := range receivers {
		// Check direct port address
		if receiver.PortAddr != nil && receiver.PortAddr.Node == nodeName {
			nodeRefs = append(nodeRefs, nodeRefInNet{
				isOutgoing: false,
				port:       receiver.PortAddr.Port,
				arrayIdx:   receiver.PortAddr.Idx,
			})
		}

		// Check chained connection
		if receiver.ChainedConnection != nil {
			if receiver.ChainedConnection.Normal != nil {
				// Check senders in the chain
				for _, sender := range receiver.ChainedConnection.Normal.Senders {
					if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
						nodeRefs = append(nodeRefs, nodeRefInNet{
							isOutgoing: true,
							port:       sender.PortAddr.Port,
							arrayIdx:   sender.PortAddr.Idx,
						})
					}
				}

				// Recursively check receivers in the chain
				nodeRefs = append(nodeRefs, findNodeUsagesInReceivers(nodeName, receiver.ChainedConnection.Normal.Receivers)...)
			}
		}

		// Check deferred connection
		if receiver.DeferredConnection != nil {
			// Similar logic to what we do with normal connections
			if receiver.DeferredConnection.Normal != nil {
				for _, sender := range receiver.DeferredConnection.Normal.Senders {
					if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
						nodeRefs = append(nodeRefs, nodeRefInNet{
							isOutgoing: true,
							port:       sender.PortAddr.Port,
							arrayIdx:   sender.PortAddr.Idx,
						})
					}
				}

				nodeRefs = append(nodeRefs, findNodeUsagesInReceivers(nodeName, receiver.DeferredConnection.Normal.Receivers)...)
			}
		}

		// Check switch cases
		if receiver.Switch != nil {
			// Check each case in the switch
			for _, caseConn := range receiver.Switch.Cases {
				for _, sender := range caseConn.Senders {
					if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
						nodeRefs = append(nodeRefs, nodeRefInNet{
							isOutgoing: true,
							port:       sender.PortAddr.Port,
							arrayIdx:   sender.PortAddr.Idx,
						})
					}
				}

				nodeRefs = append(nodeRefs, findNodeUsagesInReceivers(nodeName, caseConn.Receivers)...)
			}

			// Check default case
			if receiver.Switch.Default != nil {
				nodeRefs = append(nodeRefs, findNodeUsagesInReceivers(nodeName, receiver.Switch.Default)...)
			}
		}
	}

	return nodeRefs
}

// isCandidateCompatibleWithAllNodeRefs checks if the given component is compatible
// with all the node references by checking that all the node references
// are compatible with the component's interface.
// Compatibility is checked by comparing the port types and array usage.
// The type of the port ignored for now, to be checked later.
func isCandidateCompatibleWithAllNodeRefs(
	component src.Component,
	nodeRefs []nodeRefInNet,
) bool {
	for _, nodeRef := range nodeRefs {
		var ports map[string]src.Port
		if !nodeRef.isOutgoing {
			ports = component.Interface.IO.In
		} else {
			ports = component.Interface.IO.Out
		}

		port, portExists := ports[nodeRef.port]
		if !portExists {
			return false
		}

		isArrExpected := nodeRef.arrayIdx != nil
		if port.IsArray != isArrExpected {
			return false
		}
	}

	return true
}

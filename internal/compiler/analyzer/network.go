// This file contains the logic for analyzing network connections.
// Some methods here might look like they are related to senders or receivers specifically,
// but they are actually related to both, so they are placed here.
package analyzer

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

var ErrComplexLiteralSender = errors.New("literal network sender must have primitive type")

// analyzeNetwork must be called after analyzeNodes so we sure nodes are resolved.
func (a Analyzer) analyzeNetwork(
	net []src.Connection, // network to analyze
	iface src.Interface, // resolved interface of the component that contains the network
	hasGuard bool, // whether `?` is used by at least one node in the network
	nodes map[string]src.Node, // nodes of the component that contains the network
	nodesIfaces map[string]foundInterface, // resolved interfaces of the nodes
	scope src.Scope,
) ([]src.Connection, *compiler.Error) {
	nodesUsage := make(map[string]netNodeUsage, len(nodes))

	// Read `Union:tag` wiring to bind each `Union<T>` node to a concrete tag.
	// Example: `Union<MyU>`; `MyU::Int -> union:tag`; `42 -> union:data`
	// We bind `tag=Int` so `union:data` must be compatible with int later.
	unionActiveTags, err := a.buildUnionActiveTagBindings(net, nodes, scope)
	if err != nil {
		return nil, err
	}

	analyzedConnections, err := a.analyzeConnections(
		net,
		iface,
		nodes,
		nodesIfaces,
		nodesUsage,
		scope,
		unionActiveTags,
	)
	if err != nil {
		return nil, err
	}

	if err := a.analyzeNetPortsUsage(
		iface,
		nodesIfaces,
		hasGuard,
		nodesUsage,
		nodes,
	); err != nil {
		return nil, err
	}

	return analyzedConnections, nil
}

// analyzeConnections does two things:
// 1. Analyzes every connection and terminates with non-nil error if any of them is invalid.
// 2. Updates nodesUsage (we mutate it in-place instead of returning to avoid merging across recursive calls).
func (a Analyzer) analyzeConnections(
	net []src.Connection,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
	scope src.Scope,
	unionActiveTags map[string]unionActiveTagInfo,
) ([]src.Connection, *compiler.Error) {
	analyzedConnections := make([]src.Connection, 0, len(net))

	for _, conn := range net {
		resolvedConn, err := a.analyzeConnection(
			conn,
			iface,
			nodes,
			nodesIfaces,
			scope,
			nodesUsage,
			nil,
			net,
			unionActiveTags,
		)
		if err != nil {
			return nil, err
		}
		analyzedConnections = append(analyzedConnections, resolvedConn)
	}

	return analyzedConnections, nil
}

func (a Analyzer) analyzeConnection(
	conn src.Connection,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
	nodesUsage map[string]netNodeUsage,
	prevChainLink []src.ConnectionSender,
	net []src.Connection,
	unionActiveTags map[string]unionActiveTagInfo,
) (src.Connection, *compiler.Error) {
	if conn.ArrayBypass != nil {
		if err := a.analyzeArrayBypassConnection(
			conn,
			scope,
			iface,
			nodes,
			nodesIfaces,
			nodesUsage,
		); err != nil {
			return src.Connection{}, err
		}
		return conn, nil
	}

	analyzedNormalConn, err := a.analyzeNormalConnection(
		conn.Normal,
		iface,
		nodes,
		nodesIfaces,
		scope,
		nodesUsage,
		prevChainLink,
		net,
		unionActiveTags,
	)
	if err != nil {
		return src.Connection{}, err
	}

	return src.Connection{
		Normal: analyzedNormalConn,
		Meta:   conn.Meta,
	}, nil
}

func (a Analyzer) analyzeNormalConnection(
	normConn *src.NormalConnection,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
	nodesUsage map[string]netNodeUsage,
	prevChainLink []src.ConnectionSender,
	net []src.Connection,
	unionActiveTags map[string]unionActiveTagInfo,
) (*src.NormalConnection, *compiler.Error) {
	// Check if any receiver is a Switch.case port - if so, senders are pattern senders
	isPatternMatchingContext := hasSwitchCaseReceiver(normConn.Receivers, nodes)

	analyzedSenders, resolvedSenderTypes, err := a.analyzeSenders(
		normConn.Senders,
		scope,
		iface,
		nodes,
		nodesIfaces,
		nodesUsage,
		prevChainLink,
		isPatternMatchingContext,
	)
	if err != nil {
		return nil, err
	}

	// Switch needs special typing: for union pattern matching, each case output may
	// have a different payload type, which the type system can't express directly.
	// We infer it from the wired case tag and must do this after senders are resolved.
	analyzedSenders, resolvedSenderTypes, err = a.patchSwitchSenders(
		analyzedSenders,
		resolvedSenderTypes,
		nodes,
		scope,
		net,
	)
	if err != nil {
		return nil, err
	}

	analyzedReceivers, err := a.analyzeReceivers(
		normConn.Receivers,
		scope,
		iface,
		nodes,
		nodesIfaces,
		nodesUsage,
		resolvedSenderTypes,
		analyzedSenders,
		net,
		unionActiveTags,
	)
	if err != nil {
		return nil, err
	}

	return &src.NormalConnection{
		Senders:   analyzedSenders,
		Receivers: analyzedReceivers,
		Meta:      normConn.Meta,
	}, nil
}

// patchSwitchSenders patches sender types for Switch:case[i] outports.
// Switch is special: with union pattern matching, each case output may have a
// different payload type, which the type system can't express directly, so we
// infer it from the wired case tag.
func (a Analyzer) patchSwitchSenders(
	analyzedSenders []src.ConnectionSender,
	resolvedSenderTypes []*ts.Expr,
	nodes map[string]src.Node,
	scope src.Scope,
	net []src.Connection,
) ([]src.ConnectionSender, []*ts.Expr, *compiler.Error) {
	// Must run after sender normalization so switch:case[i] senders are resolved.
	for i, sender := range analyzedSenders {
		if sender.PortAddr == nil {
			continue
		}
		if !isSwitchCasePort(*sender.PortAddr, nodes) {
			continue
		}
		// We found a Switch:case[i] output usage; infer its concrete type.
		resolvedType, err := a.getSwitchCaseOutportType(
			*sender.PortAddr,
			nodes,
			scope,
			net,
		)
		if err != nil {
			return nil, nil, err
		}
		resolvedSenderTypes[i] = resolvedType
	}

	return analyzedSenders, resolvedSenderTypes, nil
}

func (a Analyzer) analyzeArrayBypassConnection(
	conn src.Connection,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	nodesUsage map[string]netNodeUsage,
) *compiler.Error {
	arrBypassConn := conn.ArrayBypass

	_, senderType, isArray, err := a.getPortSenderType(
		arrBypassConn.SenderOutport,
		scope,
		iface,
		nodes,
		nodesIfaces,
	)
	if err != nil {
		return compiler.Error{
			Meta: &conn.Meta,
		}.Wrap(err)
	}
	if !isArray {
		return &compiler.Error{
			Message: "Non-array outport in array-bypass connection",
			Meta:    &arrBypassConn.SenderOutport.Meta,
		}
	}

	_, receiverType, isArray, err := a.getReceiverPortType(
		arrBypassConn.ReceiverInport,
		iface,
		nodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return compiler.Error{
			Meta: &conn.Meta,
		}.Wrap(err)
	}
	if !isArray {
		return &compiler.Error{
			Message: "Non-array outport in array-bypass connection",
			Meta:    &arrBypassConn.SenderOutport.Meta,
		}
	}

	if err := a.resolver.IsSubtypeOf(
		senderType,
		receiverType,
		scope,
	); err != nil {
		return &compiler.Error{
			Message: fmt.Sprintf(
				"Incompatible types: %v -> %v: %v",
				arrBypassConn.SenderOutport, arrBypassConn.ReceiverInport, err.Error(),
			),
			Meta: &conn.Meta,
		}
	}

	if err := netNodesUsage(nodesUsage).trackOutportUsage(arrBypassConn.SenderOutport); err != nil {
		return &compiler.Error{
			Message: err.Error(),
			Meta:    &conn.Meta,
		}
	}

	if err := netNodesUsage(nodesUsage).trackInportUsage(arrBypassConn.ReceiverInport); err != nil {
		return &compiler.Error{
			Message: err.Error(),
			Meta:    &conn.Meta,
		}
	}

	return nil
}

func (a Analyzer) analyzeNetPortsUsage(
	iface src.Interface, // resolved interface of the component that contains the network
	nodesIfaces map[string]foundInterface, // resolved interfaces of the nodes in the network
	hasGuard bool, // whether `?` is used by at least one node in the network
	nodesUsage map[string]netNodeUsage,
	nodes map[string]src.Node,
) *compiler.Error {
	// 1. every self inport must be used
	inportsUsage, ok := nodesUsage["in"]
	if !ok {
		allInports := make([]string, 0, len(iface.IO.In))
		for inportName := range iface.IO.In {
			allInports = append(allInports, inportName)
		}

		return &compiler.Error{
			Message: unusedPortsMessage("inport", allInports),
			Meta:    &iface.Meta,
		}
	}

	unusedInports := make([]string, 0, len(iface.IO.In))
	for inportName := range iface.IO.In {
		if _, ok := inportsUsage.Out[inportName]; !ok { // note that self inports are outports for the network
			unusedInports = append(unusedInports, inportName)
		}
	}
	if len(unusedInports) > 0 {
		return &compiler.Error{
			Message: unusedPortsMessage("inport", unusedInports),
			Meta:    &iface.Meta,
		}
	}

	// 2. every self-outport must be used
	outportsUsage, ok := nodesUsage["out"]
	if !ok {
		return &compiler.Error{
			Message: "Component must use its outports",
			Meta:    &iface.Meta,
		}
	}

	for outportName := range iface.IO.Out {
		if _, ok := outportsUsage.In[outportName]; ok { // self outports are inports in network
			continue
		}

		// err outport is allowed to be unused if parent uses guard
		if outportName == "err" && hasGuard {
			continue
		}

		return &compiler.Error{
			Message: fmt.Sprintf("Unused outport: %v", outportName),
		}
	}

	// 3. check sub-nodes usage in network
	for nodeName, nodeIface := range nodesIfaces {
		nodeMeta := nodes[nodeName].Meta

		// every sub-node must be used
		nodeUsage, ok := nodesUsage[nodeName]
		if !ok {
			return &compiler.Error{
				Message: fmt.Sprintf("Unused node found: %v", nodeName),
				Meta:    &nodeMeta,
			}
		}

		// every sub-node's inport must be used
		for inportName := range nodeIface.iface.IO.In {
			if _, ok := nodeUsage.In[inportName]; ok {
				continue
			}

			_, portless := nodeUsage.In[""]
			if portless && len(nodeIface.iface.IO.In) == 1 {
				continue
			}

			return &compiler.Error{
				Message: fmt.Sprintf(
					"Unused node inport: %v:%v",
					nodeName,
					inportName,
				),
				Meta: &nodeMeta,
			}
		}

		if len(nodeIface.iface.IO.Out) == 0 { // e.g. Del
			continue
		}

		// :err outport must always be used + at least one outport must be used in general
		atLeastOneOutportIsUsed := false
		for outportName := range nodeIface.iface.IO.Out {
			if _, ok := nodeUsage.Out[outportName]; ok {
				atLeastOneOutportIsUsed = true
				continue
			}

			if outportName == "err" && !nodes[nodeName].ErrGuard {
				return &compiler.Error{
					Message: fmt.Sprintf("unhandled error: %v:err", nodeName),
					Meta:    &nodeMeta,
				}
			}
		}

		if !atLeastOneOutportIsUsed {
			if _, ok := nodeUsage.Out[""]; ok && len(nodeIface.iface.IO.Out) == 1 {
				continue
			}
			return &compiler.Error{
				Message: fmt.Sprintf("All node's outports are unused: %v", nodeName),
				Meta:    &nodeMeta,
			}
		}
	}

	// 4. check that array ports are used correctly (from 0 and without holes)
	for nodeName, nodeUsage := range nodesUsage {
		nodeMeta := nodes[nodeName].Meta

		for portName, usedSlots := range nodeUsage.In {
			if usedSlots == nil {
				continue // skip non-array ports
			}

			maxSlot := uint8(0)
			for slot := range usedSlots {
				if slot > maxSlot {
					maxSlot = slot
				}
			}

			for i := uint8(0); i <= maxSlot; i++ {
				if _, ok := usedSlots[i]; !ok {
					return &compiler.Error{
						Message: fmt.Sprintf(
							"array inport '%s:%s' is used incorrectly: slot %d is missing",
							nodeName,
							portName,
							i,
						),
						Meta: &nodeMeta,
					}
				}
			}
		}

		for portName, usedSlots := range nodeUsage.Out {
			if usedSlots == nil {
				continue // skip non-array ports
			}

			maxSlot := uint8(0)
			for slot := range usedSlots {
				if slot > maxSlot {
					maxSlot = slot
				}
			}

			for i := uint8(0); i <= maxSlot; i++ {
				if _, ok := usedSlots[i]; !ok {
					return &compiler.Error{
						Message: fmt.Sprintf(
							"array outport '%s:%s' is used incorrectly: slot %d is missing",
							nodeName,
							portName,
							i,
						),
						Meta: &nodeMeta,
					}
				}
			}
		}
	}

	return nil
}

func unusedPortsMessage(portType string, ports []string) string {
	sort.Strings(ports)
	if len(ports) == 1 {
		return fmt.Sprintf("Unused %s: %s", portType, ports[0])
	}
	return fmt.Sprintf("Unused %ss: %s", portType, strings.Join(ports, ", "))
}

// getResolvedPortType returns resolved port-addr, type expr and isArray bool.
// Resolved port is equal to the given one unless it was an "" empty string.
func (a Analyzer) getResolvedPortType(
	ports map[string]src.Port,
	nodeIfaceParams []ts.Param,
	portAddr src.PortAddr,
	node src.Node,
	scope src.Scope,
	isInput bool,
) (src.PortAddr, ts.Expr, bool, *compiler.Error) {
	if portAddr.Port == "" {
		resolvedPortAddr, err := a.resolveUnnamedPort(ports, portAddr, node, isInput)
		if err != nil {
			return src.PortAddr{}, ts.Expr{}, false, err
		}
		portAddr = resolvedPortAddr
	}

	port, ok := ports[portAddr.Port]
	if !ok {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: fmt.Sprintf(
				"Port not found '%v'",
				portAddr,
			),
			Meta: &portAddr.Meta,
		}
	}

	resolvedPortType, err := a.resolvePortTypeWithFrame(
		port,
		nodeIfaceParams,
		node.TypeArgs,
		scope,
	)
	if err != nil {
		return src.PortAddr{}, ts.Expr{}, false, err
	}

	return portAddr, resolvedPortType, port.IsArray, nil
}

func (a Analyzer) resolveUnnamedPort(
	ports map[string]src.Port,
	portAddr src.PortAddr,
	node src.Node,
	isInput bool,
) (src.PortAddr, *compiler.Error) {
	allowUnnamed := len(ports) == 1 || (!isInput && len(ports) == 2 && node.ErrGuard)
	if !allowUnnamed {
		kind := "outports"
		if isInput {
			kind = "inports"
		}
		return src.PortAddr{}, &compiler.Error{
			Message: fmt.Sprintf(
				"node '%v' has multiple %s - port name must be specified",
				portAddr.Node,
				kind,
			),
			Meta: &portAddr.Meta,
		}
	}

	// for output ports with error guard, skip the 'err' port and select the first non-error port
	if !isInput && node.ErrGuard {
		for name := range ports {
			if name != "err" {
				portAddr.Port = name
				break
			}
		}
	} else {
		for name := range ports {
			portAddr.Port = name
			break
		}
	}
	return portAddr, nil
}

func (a Analyzer) resolvePortTypeWithFrame(
	port src.Port,
	nodeIfaceParams []ts.Param,
	resolvedNodeArgs []ts.Expr,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	// we don't resolve node's args assuming they resolved already
	frame := make(map[string]ts.Def, len(nodeIfaceParams))
	for i, param := range nodeIfaceParams {
		arg := resolvedNodeArgs[i]
		frame[param.Name] = ts.Def{
			BodyExpr: &arg,
			Meta:     arg.Meta,
		}
	}

	resolvedPortType, err := a.resolver.ResolveExprWithFrame(
		port.TypeExpr,
		frame,
		scope,
	)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &port.Meta,
		}
	}

	return resolvedPortType, nil
}

func (a Analyzer) getResolvedSenderType(
	sender src.ConnectionSender,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
	prevChainLink []src.ConnectionSender,
	nodesUsage map[string]netNodeUsage,
	isPatternSender bool,
) (src.ConnectionSender, ts.Expr, bool, *compiler.Error) {
	if sender.Const != nil {
		resolvedConst, resolvedExpr, err := a.getConstSenderType(*sender.Const, scope)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, err
		}

		return src.ConnectionSender{
			Const: &resolvedConst,
			Meta:  sender.Meta,
		}, resolvedExpr, false, nil
	}

	if len(sender.StructSelector) > 0 {
		_, chainLinkType, _, err := a.getResolvedSenderType(
			prevChainLink[0],
			iface,
			nodes,
			nodesIfaces,
			scope,
			prevChainLink,
			nodesUsage,
			isPatternSender,
		)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, err
		}

		lastFieldType, err := a.getSelectorsSenderType(
			chainLinkType,
			sender.StructSelector,
			scope,
		)
		if err != nil {
			return src.ConnectionSender{}, ts.Expr{}, false, compiler.Error{
				Meta: &sender.Meta,
			}.Wrap(err)
		}

		return sender, lastFieldType, false, nil
	}

	// handle port-address sender
	resolvedPort, resolvedExpr, isArr, err := a.getPortSenderType(
		*sender.PortAddr,
		scope,
		iface,
		nodes,
		nodesIfaces,
	)
	if err != nil {
		return src.ConnectionSender{}, ts.Expr{}, false, err
	}

	return src.ConnectionSender{
		PortAddr: &resolvedPort,
		Meta:     sender.Meta,
	}, resolvedExpr, isArr, nil
}

// getPortSenderType returns resolved port-addr, type expr and isArray bool.
// Resolved port is equal to the given one unless it was an "" empty string.
func (a Analyzer) getPortSenderType(
	senderSidePortAddr src.PortAddr,
	scope src.Scope,
	iface src.Interface,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
) (src.PortAddr, ts.Expr, bool, *compiler.Error) {
	if senderSidePortAddr.Node == "out" {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: "Component cannot read from self outport",
			Meta:    &senderSidePortAddr.Meta,
		}
	}

	if senderSidePortAddr.Node == "in" {
		inports := iface.IO.In

		inport, ok := inports[senderSidePortAddr.Port]
		if !ok {
			return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
				Message: fmt.Sprintf("Referenced inport not found in component's interface: %v", senderSidePortAddr.Port),
				Meta:    &senderSidePortAddr.Meta,
			}
		}

		resolvedInportType, err := a.resolver.ResolveExprWithFrame(
			inport.TypeExpr,
			iface.TypeParams.ToFrame(),
			scope,
		)
		if err != nil {
			return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
				Message: err.Error(),
				Meta:    &senderSidePortAddr.Meta,
			}
		}

		return senderSidePortAddr, resolvedInportType, inport.IsArray, nil
	}

	return a.getNodeOutportType(
		senderSidePortAddr, nodes, nodesIfaces, scope,
	)
}

func (a Analyzer) getConstSenderType(
	constSender src.Const,
	scope src.Scope,
) (src.Const, ts.Expr, *compiler.Error) {
	if constSender.Value.Ref != nil {
		expr, err := a.getResolvedConstTypeByRef(*constSender.Value.Ref, scope)
		if err != nil {
			return src.Const{}, ts.Expr{}, compiler.Error{
				Meta: &constSender.Value.Ref.Meta,
			}.Wrap(err)
		}
		return constSender, expr, nil
	}

	if constSender.Value.Message == nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Message: "Literal sender type is empty",
			Meta:    &constSender.Meta,
		}
	}

	resolvedExpr, err := a.resolver.ResolveExpr(
		constSender.TypeExpr,
		scope,
	)
	if err != nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &constSender.Value.Message.Meta,
		}
	}

	if err := a.validateLiteralSender(resolvedExpr); err != nil {
		return src.Const{}, ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &constSender.Value.Message.Meta,
		}
	}

	return src.Const{
		TypeExpr: resolvedExpr,
		Value: src.ConstValue{
			Message: &src.MsgLiteral{
				Bool:         constSender.Value.Message.Bool,
				Int:          constSender.Value.Message.Int,
				Float:        constSender.Value.Message.Float,
				Str:          constSender.Value.Message.Str,
				List:         constSender.Value.Message.List,
				DictOrStruct: constSender.Value.Message.DictOrStruct,
				Union:        constSender.Value.Message.Union,
				Meta:         constSender.Value.Message.Meta,
			},
		},
		Meta: constSender.Meta,
	}, resolvedExpr, nil
}

func (a Analyzer) validateLiteralSender(resolvedExpr ts.Expr) error {
	if resolvedExpr.Inst != nil {
		switch resolvedExpr.Inst.Ref.String() {
		case "bool", "int", "float", "string":
			return nil
		}
		return ErrComplexLiteralSender
	}

	if resolvedExpr.Lit == nil ||
		resolvedExpr.Lit.Union == nil {
		return ErrComplexLiteralSender
	}

	return nil
}

// getNodeOutportType returns resolved port-addr, type expr and isArray bool.
// Resolved port is equal to the given one unless it was an "" empty string.
func (a Analyzer) getNodeOutportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
) (src.PortAddr, ts.Expr, bool, *compiler.Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: fmt.Sprintf("Referenced node not found: %v", portAddr.Node),
			Meta:    &portAddr.Meta,
		}
	}

	nodeIface, ok := nodesIfaces[portAddr.Node]
	if !ok {
		return src.PortAddr{}, ts.Expr{}, false, &compiler.Error{
			Message: fmt.Sprintf("Referenced node not found: %v", portAddr.Node),
			Meta:    &portAddr.Meta,
		}
	}

	resolvedPort, resolvedType, isArray, err := a.getResolvedPortType(
		nodeIface.iface.IO.Out,
		nodeIface.iface.TypeParams.Params,
		portAddr,
		node,
		scope.Relocate(nodeIface.location),
		false,
	)
	if err != nil {
		return src.PortAddr{}, ts.Expr{}, false, err
	}

	return resolvedPort, resolvedType, isArray, nil
}

func (a Analyzer) getResolvedConstTypeByRef(
	ref core.EntityRef,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	constant, loc, err := scope.GetConst(ref)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &ref.Meta,
		}
	}

	if constant.Value.Ref != nil {
		expr, err := a.getResolvedConstTypeByRef(*constant.Value.Ref, scope)
		if err != nil {
			return ts.Expr{}, compiler.Error{
				Meta: &constant.Meta,
			}.Wrap(err)
		}
		return expr, nil
	}

	scope = scope.Relocate(loc)

	resolvedExpr, err := a.resolver.ResolveExpr(constant.TypeExpr, scope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &constant.Value.Message.Meta,
		}
	}

	return resolvedExpr, nil
}

func (a Analyzer) getSelectorsSenderType(
	senderType ts.Expr,
	selectors []string,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	if len(selectors) == 0 {
		return senderType, nil
	}

	// if it's an instance type, resolve it to get the actual struct definition
	if senderType.Inst != nil {
		resolvedType, err := a.resolver.ResolveExpr(senderType, scope)
		if err != nil {
			return ts.Expr{}, &compiler.Error{
				Message: fmt.Sprintf("Failed to resolve type: %v", err),
			}
		}
		senderType = resolvedType
	}

	if senderType.Lit == nil || senderType.Lit.Struct == nil {
		return ts.Expr{}, &compiler.Error{
			Message: fmt.Sprintf("Type not struct: %v", senderType.String()),
		}
	}

	curField := selectors[0]
	fieldType, ok := senderType.Lit.Struct[curField]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Message: fmt.Sprintf("struct field '%v' not found", curField),
		}
	}

	return a.getSelectorsSenderType(fieldType, selectors[1:], scope)
}

package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

type switchAnalysis struct {
	byNode map[string]*switchCaseAnalysis
}

type switchCaseAnalysis struct {
	dataType     ts.Expr
	unionTags    map[string]*ts.Expr
	casePatterns map[uint8]switchCasePattern
	caseUsed     map[uint8]struct{}
	hasElse      bool
}

type switchCasePattern struct {
	tag         string
	payloadType *ts.Expr
}

func (s switchAnalysis) caseOutportType(nodeName string, idx uint8) (ts.Expr, bool) {
	info, ok := s.byNode[nodeName]
	if !ok {
		return ts.Expr{}, false
	}
	pattern, ok := info.casePatterns[idx]
	if !ok {
		return info.dataType, true
	}
	if pattern.payloadType == nil {
		return info.dataType, true
	}
	return *pattern.payloadType, true
}

func (a Analyzer) collectSwitchAnalysis(
	net []src.Connection,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
) (switchAnalysis, *compiler.Error) {
	analysis := switchAnalysis{byNode: make(map[string]*switchCaseAnalysis)}

	getOrInitSwitch := func(nodeName string) (*switchCaseAnalysis, *compiler.Error) {
		if existing, ok := analysis.byNode[nodeName]; ok {
			return existing, nil
		}

		dataType, err := a.getSwitchDataType(nodeName, nodes, nodesIfaces, scope)
		if err != nil {
			return nil, err
		}

		resolvedDataType, err := a.resolveTypeLiteral(dataType, nodesIfaces, nodeName, scope)
		if err != nil {
			return nil, err
		}

		info := &switchCaseAnalysis{
			dataType:     resolvedDataType,
			casePatterns: make(map[uint8]switchCasePattern),
			caseUsed:     make(map[uint8]struct{}),
		}
		if resolvedDataType.Lit != nil && resolvedDataType.Lit.Union != nil {
			info.unionTags = resolvedDataType.Lit.Union
		}

		analysis.byNode[nodeName] = info
		return info, nil
	}

	var visitConnection func(conn src.Connection) *compiler.Error
	visitConnection = func(conn src.Connection) *compiler.Error {
		if conn.Normal == nil {
			return nil
		}

		for _, sender := range conn.Normal.Senders {
			if sender.PortAddr != nil && isSwitchElsePort(*sender.PortAddr, nodes) {
				info, err := getOrInitSwitch(sender.PortAddr.Node)
				if err != nil {
					return err
				}
				info.hasElse = true
			}
		}

		switchReceivers := collectSwitchReceiverPorts(conn.Normal.Receivers, nodes)
		if len(switchReceivers) > 0 {
			unionSenders, unionSenderMeta := collectUnionSenders(conn.Normal.Senders)
			for _, receiver := range switchReceivers {
				info, err := getOrInitSwitch(receiver.PortAddr.Node)
				if err != nil {
					return err
				}

				if receiver.PortAddr.Port != "case" {
					continue
				}
				if receiver.PortAddr.Idx != nil {
					info.caseUsed[*receiver.PortAddr.Idx] = struct{}{}
				}

				if len(unionSenders) == 0 {
					continue
				}

				if info.unionTags == nil {
					return &compiler.Error{
						Message: "Switch pattern matching requires a tagged union input",
						Meta:    &receiver.PortAddr.Meta,
					}
				}

				tag := unionSenders[0].Tag
				for _, sender := range unionSenders[1:] {
					if sender.Tag != tag {
						return &compiler.Error{
							Message: "Switch case port cannot match multiple union tags",
							Meta:    &unionSenderMeta,
						}
					}
				}

				tagType, ok := info.unionTags[tag]
				if !ok {
					return &compiler.Error{
						Message: fmt.Sprintf("tag %q is not part of union input", tag),
						Meta:    &receiver.PortAddr.Meta,
					}
				}

				if receiver.PortAddr.Idx == nil {
					continue
				}

				existing, hasExisting := info.casePatterns[*receiver.PortAddr.Idx]
				if hasExisting && existing.tag != tag {
					return &compiler.Error{
						Message: "Switch case port cannot match multiple union tags",
						Meta:    &receiver.PortAddr.Meta,
					}
				}

				info.casePatterns[*receiver.PortAddr.Idx] = switchCasePattern{
					tag:         tag,
					payloadType: tagType,
				}
			}
		}

		for _, receiver := range conn.Normal.Receivers {
			if receiver.ChainedConnection != nil {
				if err := visitConnection(*receiver.ChainedConnection); err != nil {
					return err
				}
			}
			if receiver.DeferredConnection != nil {
				if err := visitConnection(*receiver.DeferredConnection); err != nil {
					return err
				}
			}
		}

		return nil
	}

	for _, conn := range net {
		if err := visitConnection(conn); err != nil {
			return switchAnalysis{}, err
		}
	}

	for nodeName, info := range analysis.byNode {
		node, ok := nodes[nodeName]
		if !ok {
			continue
		}
		if info.unionTags == nil {
			continue
		}

		if len(info.caseUsed) == 0 {
			return switchAnalysis{}, &compiler.Error{
				Message: "Switch with union input must define at least one case",
				Meta:    &node.Meta,
			}
		}

		covered := make(map[string]struct{}, len(info.casePatterns))
		for _, pattern := range info.casePatterns {
			covered[pattern.tag] = struct{}{}
		}

		allCovered := true
		for tag := range info.unionTags {
			if _, ok := covered[tag]; !ok {
				allCovered = false
				break
			}
		}

		if allCovered {
			if info.hasElse {
				return switchAnalysis{}, &compiler.Error{
					Message: "Switch :else is not allowed when all union tags are covered",
					Meta:    &node.Meta,
				}
			}
			continue
		}

		if !info.hasElse {
			return switchAnalysis{}, &compiler.Error{
				Message: "Switch with union input must connect :else when not all tags are covered",
				Meta:    &node.Meta,
			}
		}
	}

	return analysis, nil
}

func (a Analyzer) collectSwitchOutportTypesForUsage(
	net []src.Connection,
	nodes map[string]src.Node,
	parentFrame map[string]ts.Def,
	scope src.Scope,
) map[string]map[uint8]ts.Expr {
	type usageInfo struct {
		dataType     ts.Expr
		unionTags    map[string]*ts.Expr
		casePatterns map[uint8]switchCasePattern
	}

	switches := make(map[string]*usageInfo)
	getOrInitSwitch := func(nodeName string) *usageInfo {
		if existing, ok := switches[nodeName]; ok {
			return existing
		}

		node, ok := nodes[nodeName]
		if !ok || len(node.TypeArgs) == 0 {
			return nil
		}

		resolvedDataType, err := a.resolver.ResolveExprWithFrame(node.TypeArgs[0], parentFrame, scope)
		if err != nil {
			return nil
		}

		info := &usageInfo{
			dataType:     resolvedDataType,
			casePatterns: make(map[uint8]switchCasePattern),
		}
		if resolvedDataType.Lit != nil && resolvedDataType.Lit.Union != nil {
			info.unionTags = resolvedDataType.Lit.Union
		}
		switches[nodeName] = info
		return info
	}

	var visitConnection func(conn src.Connection)
	visitConnection = func(conn src.Connection) {
		if conn.Normal == nil {
			return
		}

		switchReceivers := collectSwitchReceiverPorts(conn.Normal.Receivers, nodes)
		if len(switchReceivers) > 0 {
			unionSenders, _ := collectUnionSenders(conn.Normal.Senders)
			for _, receiver := range switchReceivers {
				if receiver.PortAddr.Port != "case" || receiver.PortAddr.Idx == nil {
					continue
				}

				info := getOrInitSwitch(receiver.PortAddr.Node)
				if info == nil || info.unionTags == nil || len(unionSenders) == 0 {
					continue
				}

				tag := unionSenders[0].Tag
				for _, sender := range unionSenders[1:] {
					if sender.Tag != tag {
						return
					}
				}

				tagType, ok := info.unionTags[tag]
				if !ok {
					return
				}

				info.casePatterns[*receiver.PortAddr.Idx] = switchCasePattern{
					tag:         tag,
					payloadType: tagType,
				}
			}
		}

		for _, receiver := range conn.Normal.Receivers {
			if receiver.ChainedConnection != nil {
				visitConnection(*receiver.ChainedConnection)
			}
			if receiver.DeferredConnection != nil {
				visitConnection(*receiver.DeferredConnection)
			}
		}
	}

	for _, conn := range net {
		visitConnection(conn)
	}

	result := make(map[string]map[uint8]ts.Expr)
	for nodeName, info := range switches {
		for idx, pattern := range info.casePatterns {
			outType := info.dataType
			if pattern.payloadType != nil {
				outType = *pattern.payloadType
			}
			if result[nodeName] == nil {
				result[nodeName] = make(map[uint8]ts.Expr)
			}
			result[nodeName][idx] = outType
		}
	}

	return result
}

func (a Analyzer) getSwitchDataType(
	nodeName string,
	nodes map[string]src.Node,
	nodesIfaces map[string]foundInterface,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	portAddr := src.PortAddr{Node: nodeName, Port: "data"}
	_, resolvedType, _, err := a.getNodeInportType(portAddr, nodes, nodesIfaces, scope)
	if err != nil {
		return ts.Expr{}, err
	}
	return resolvedType, nil
}

func (a Analyzer) resolveTypeLiteral(
	expr ts.Expr,
	nodesIfaces map[string]foundInterface,
	nodeName string,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	resolvedScope := scope
	if nodeIface, ok := nodesIfaces[nodeName]; ok {
		resolvedScope = scope.Relocate(nodeIface.location)
	}
	resolvedExpr, err := a.resolver.ResolveExpr(expr, resolvedScope)
	if err != nil {
		return ts.Expr{}, &compiler.Error{Message: err.Error()}
	}
	return resolvedExpr, nil
}

type switchReceiverPort struct {
	PortAddr src.PortAddr
}

func collectSwitchReceiverPorts(receivers []src.ConnectionReceiver, nodes map[string]src.Node) []switchReceiverPort {
	var res []switchReceiverPort
	for _, receiver := range receivers {
		if receiver.PortAddr != nil {
			if isSwitchCasePort(*receiver.PortAddr, nodes) {
				res = append(res, switchReceiverPort{PortAddr: *receiver.PortAddr})
			}
		}

		if receiver.ChainedConnection != nil && receiver.ChainedConnection.Normal != nil {
			chainHead := receiver.ChainedConnection.Normal.Senders
			if len(chainHead) > 0 && chainHead[0].PortAddr != nil {
				portAddr := *chainHead[0].PortAddr
				if isSwitchCasePort(portAddr, nodes) {
					res = append(res, switchReceiverPort{PortAddr: portAddr})
				}
			}
		}
	}
	return res
}

func collectUnionSenders(senders []src.ConnectionSender) ([]src.UnionSender, core.Meta) {
	var res []src.UnionSender
	var meta core.Meta
	for _, sender := range senders {
		if sender.Union != nil {
			res = append(res, *sender.Union)
			meta = sender.Union.Meta
		}
	}
	return res, meta
}

func isSwitchElsePort(portAddr src.PortAddr, nodes map[string]src.Node) bool {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return false
	}
	return node.EntityRef.Name == "Switch" && portAddr.Port == "else"
}

package irgen

import (
	"fmt"
	"sort"
	"strings"

	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ir"
)

// processNetwork inserts connections to result and returns metadata about the network.
func (g Generator) processNetwork(
	conns []src.Connection,
	scope *src.Scope,
	nodeCtx nodeContext,
	result *ir.Program,
) map[string]portsUsage {
	nodesPortsUsage := map[string]portsUsage{}

	for _, conn := range conns {
		if sender, receiver, ok := src.ArrayBypassPorts(conn.Normal); ok {
			g.processArrayBypassConnection(
				*sender,
				*receiver,
				nodesPortsUsage,
				nodeCtx,
				result,
			)
			continue
		}

		if len(conn.Normal.Senders) != 1 || len(conn.Normal.Receivers) != 1 {
			panic("not 1-1 connection found after desugaring")
		}

		g.processNormalConnection(
			nodeCtx,
			scope,
			conn,
			nodesPortsUsage,
			result,
		)
	}

	return nodesPortsUsage
}

func (Generator) processArrayBypassConnection(
	senderPort src.PortAddr,
	receiverPort src.PortAddr,
	nodesPortsUsage map[string]portsUsage,
	nodeCtx nodeContext,
	result *ir.Program,
) {
	// here's how we handle array-bypass connections:
	// sender is always component's inport
	// based on that, we can set receiver's inport slots
	// equal to slots of our own inport

	arrBypassSender := senderPort
	arrBypassReceiver := receiverPort

	if _, ok := nodesPortsUsage[arrBypassReceiver.Node]; !ok {
		nodesPortsUsage[arrBypassReceiver.Node] = portsUsage{
			in:  map[relPortAddr]struct{}{},
			out: map[relPortAddr]struct{}{},
		}
	}

	var slotIdx uint8 = 0
	for usageAddr := range nodeCtx.portsUsage.in {
		if usageAddr.Port != arrBypassSender.Port {
			continue
		}

		slotIdxCopy := slotIdx

		addr := relPortAddr{Port: arrBypassReceiver.Port, Idx: &slotIdxCopy}
		nodesPortsUsage[arrBypassReceiver.Node].in[addr] = struct{}{}

		irSenderSlot := ir.PortAddr{
			Path:    joinNodePath(nodeCtx.path, arrBypassSender.Node),
			Port:    arrBypassSender.Port,
			Idx:     slotIdxCopy,
			IsArray: true,
		}

		irReceiverSlot := ir.PortAddr{
			Path:    joinNodePath(nodeCtx.path, arrBypassReceiver.Node) + "/in",
			Port:    arrBypassReceiver.Port,
			Idx:     slotIdxCopy,
			IsArray: true,
		}

		result.Connections[irSenderSlot] = irReceiverSlot

		slotIdx++
	}
}

func (g Generator) processNormalConnection(
	nodeCtx nodeContext,
	scope *src.Scope,
	conn src.Connection,
	nodesPortsUsage map[string]portsUsage,
	result *ir.Program,
) {
	irSenderSidePortAddr := g.processSender(
		nodeCtx,
		scope,
		conn.Normal.Senders[0],
		nodesPortsUsage,
	)
	irReceiverPortAddr := g.processReceiver(
		nodeCtx,
		scope,
		conn.Normal.Receivers[0],
		nodesPortsUsage,
	)
	result.Connections[irSenderSidePortAddr] = irReceiverPortAddr
}

func (g Generator) processSender(
	nodeCtx nodeContext,
	scope *src.Scope,
	sender src.ConnectionSender,
	nodesUsage map[string]portsUsage,
) ir.PortAddr {
	// other special senders should also have been desugared
	if sender.PortAddr == nil {
		panic(fmt.Sprintf(
			"INTERNAL ERROR: sender with nil PortAddr was not desugared (const=%v, location: %v)",
			sender.Const != nil,
			sender.Meta.Location,
		))
	}

	// there could be many connections with the same sender but we must only add it once
	if _, ok := nodesUsage[sender.PortAddr.Node]; !ok {
		nodesUsage[sender.PortAddr.Node] = portsUsage{
			in:  map[relPortAddr]struct{}{},
			out: map[relPortAddr]struct{}{},
		}
	}

	// if sender node is dependency from DI and if port we are referring to is an empty string
	// we need to find dependency component and use its outport name
	// this is techically desugaring at irgen level but it's impossible to desugare before
	// because only irgen really builds nodes and passes DI args to them
	depNode, isNodeDep := nodeCtx.node.DIArgs[sender.PortAddr.Node]
	if isNodeDep && sender.PortAddr.Port == "" {
		versions, err := scope.
			Relocate(depNode.Meta.Location).
			GetComponent(depNode.EntityRef)
		if err != nil {
			panic(err)
		}

		var version src.Component
		if len(versions) == 1 {
			version = versions[0]
		} else {
			version = versions[*depNode.OverloadIndex]
		}

		for outport := range version.IO.Out {
			sender.PortAddr.Port = outport
			break
		}
	}

	// insert outport usage
	nodesUsage[sender.PortAddr.Node].out[relPortAddr{
		Port: sender.PortAddr.Port,
		Idx:  sender.PortAddr.Idx,
	}] = struct{}{}

	// create ir port addr
	irSenderPort := ir.PortAddr{
		Path: joinNodePath(nodeCtx.path, sender.PortAddr.Node),
		Port: sender.PortAddr.Port,
	}
	if sender.PortAddr.Idx != nil {
		irSenderPort.IsArray = true
		irSenderPort.Idx = *sender.PortAddr.Idx
	}
	if sender.PortAddr.Node != "in" {
		irSenderPort.Path += "/out"
	}

	return irSenderPort
}

// processReceiver maps src connection side to ir connection side 1-1 just making the port addr's path absolute
func (g Generator) processReceiver(
	nodeCtx nodeContext,
	scope *src.Scope,
	receiver src.ConnectionReceiver,
	nodesUsage map[string]portsUsage,
) ir.PortAddr {
	// same receiver can be used by multiple senders so we only add it once
	if _, ok := nodesUsage[receiver.PortAddr.Node]; !ok {
		nodesUsage[receiver.PortAddr.Node] = portsUsage{
			in:  map[relPortAddr]struct{}{},
			out: map[relPortAddr]struct{}{},
		}
	}

	// if receiver node DI
	diArgNode, isDI := nodeCtx.node.DIArgs[receiver.PortAddr.Node]
	// and if port we are referring to is an empty string
	if isDI && receiver.PortAddr.Port == "" {
		// we need to find dependency component and use its inport name
		// this is techically desugaring at irgen level
		// but it's impossible to desugare before, because only irgen really builds nodes

		versions, err := scope.
			Relocate(diArgNode.Meta.Location).
			GetComponent(diArgNode.EntityRef)
		if err != nil {
			panic(err)
		}

		var version src.Component
		if len(versions) == 1 {
			version = versions[0]
		} else {
			version = versions[*diArgNode.OverloadIndex]
		}

		for inport := range version.IO.In {
			receiver.PortAddr.Port = inport
			break
		}
	}

	receiverNode := receiver.PortAddr.Node
	receiverPortAddr := relPortAddr{
		Port: receiver.PortAddr.Port,
		Idx:  receiver.PortAddr.Idx,
	}

	nodesUsage[receiverNode].in[receiverPortAddr] = struct{}{}

	result := ir.PortAddr{
		Path: joinNodePath(nodeCtx.path, receiver.PortAddr.Node),
		Port: receiver.PortAddr.Port,
	}
	if receiver.PortAddr.Idx != nil {
		result.IsArray = true
		result.Idx = *receiver.PortAddr.Idx
	}

	// 'out' node is actually receiver but we don't want to have 'out/in' path
	if receiver.PortAddr.Node != "out" {
		result.Path += "/in"
	}

	return result
}

func joinNodePath(nodePath []string, nodeName string) string {
	newPath := make([]string, len(nodePath))
	copy(newPath, nodePath)
	newPath = append(newPath, nodeName)
	return strings.Join(newPath, "/")
}

// this is very important because runtime function calls depends on this order.
func sortPortAddrs(addrs []ir.PortAddr) {
	sort.Slice(addrs, func(i, j int) bool {
		if addrs[i].Path != addrs[j].Path {
			return addrs[i].Path < addrs[j].Path
		}
		if addrs[i].Port != addrs[j].Port {
			return addrs[i].Port < addrs[j].Port
		}
		if !addrs[i].IsArray {
			return true
		}
		return addrs[i].Idx < addrs[j].Idx
	})
}

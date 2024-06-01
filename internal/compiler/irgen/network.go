package irgen

import (
	"fmt"
	"sort"
	"strings"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/runtime/ir"
)

// processNetwork
// 1) inserts network connections
// 2) returns metadata about how subnodes are used by this network
func (g Generator) processNetwork(
	conns []src.Connection,
	nodeCtx nodeContext,
	result *ir.Program,
) (map[string]portsUsage, error) {
	nodesPortsUsage := map[string]portsUsage{}

	for _, conn := range conns {
		// here's how we handle array-bypass connections:
		// sender is always flow's inport
		// based on that, we can set receiver's inport slots
		// equal to slots of our inport
		if conn.ArrayBypass != nil {
			arrBypassSender := conn.ArrayBypass.SenderOutport
			arrBypassReceiver := conn.ArrayBypass.ReceiverInport

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

				addr := relPortAddr{Port: arrBypassReceiver.Port, Idx: slotIdx}
				nodesPortsUsage[arrBypassReceiver.Node].in[addr] = struct{}{}

				irSenderSlot := ir.PortAddr{
					Path: joinNodePath(nodeCtx.path, arrBypassSender.Node),
					Port: arrBypassSender.Port,
					Idx:  uint32(slotIdx),
				}

				irReceiverSlot := ir.PortAddr{
					Path: joinNodePath(nodeCtx.path, arrBypassReceiver.Node) + "/in",
					Port: arrBypassReceiver.Port,
					Idx:  uint32(slotIdx),
				}

				result.Connections[irSenderSlot] = map[ir.PortAddr]struct{}{
					irReceiverSlot: {},
				}

				slotIdx++
			}

			continue
		}

		senderSide := conn.Normal.SenderSide

		irSenderSidePortAddr, err := g.processSenderSide(
			nodeCtx,
			senderSide,
			nodesPortsUsage,
		)
		if err != nil {
			return nil, fmt.Errorf("process sender side: %w", err)
		}

		receiverPortsIR := make(map[ir.PortAddr]struct{}, len(conn.Normal.ReceiverSide.Receivers))
		for _, receiverSide := range conn.Normal.ReceiverSide.Receivers {
			receiverSideIR := g.mapReceiverSide(nodeCtx.path, receiverSide)
			receiverPortsIR[*receiverSideIR] = struct{}{}

			// same receiver can be used by multiple senders so we only add it once
			if _, ok := nodesPortsUsage[receiverSide.PortAddr.Node]; !ok {
				nodesPortsUsage[receiverSide.PortAddr.Node] = portsUsage{
					in:  map[relPortAddr]struct{}{},
					out: map[relPortAddr]struct{}{},
				}
			}

			var idx uint8
			if receiverSide.PortAddr.Idx != nil {
				idx = *receiverSide.PortAddr.Idx
			}

			receiverNode := receiverSide.PortAddr.Node
			receiverPortAddr := relPortAddr{
				Port: receiverSide.PortAddr.Port,
				Idx:  idx,
			}

			nodesPortsUsage[receiverNode].in[receiverPortAddr] = struct{}{}
		}

		result.Connections[*irSenderSidePortAddr] = receiverPortsIR
	}

	return nodesPortsUsage, nil
}

func (g Generator) processSenderSide(
	nodeCtx nodeContext,
	senderSide src.ConnectionSenderSide,
	result map[string]portsUsage,
) (*ir.PortAddr, error) {
	// there could be many connections with the same sender but we must only add it once
	if _, ok := result[senderSide.PortAddr.Node]; !ok {
		result[senderSide.PortAddr.Node] = portsUsage{
			in:  map[relPortAddr]struct{}{},
			out: map[relPortAddr]struct{}{},
		}
	}

	var idx uint8
	if senderSide.PortAddr.Idx != nil {
		idx = *senderSide.PortAddr.Idx
	}

	// insert outport usage
	result[senderSide.PortAddr.Node].out[relPortAddr{
		Port: senderSide.PortAddr.Port,
		Idx:  idx,
	}] = struct{}{}

	irSenderSide := &ir.PortAddr{
		Path: joinNodePath(nodeCtx.path, senderSide.PortAddr.Node),
		Port: senderSide.PortAddr.Port,
		Idx:  uint32(idx),
	}

	if senderSide.PortAddr.Node == "in" {
		return irSenderSide, nil
	}
	irSenderSide.Path += "/out"

	return irSenderSide, nil
}

// for in we only use parent ctx cuz all inports are used
func (Generator) getFuncInports(nodeCtx nodeContext) []ir.PortAddr {
	inports := make([]ir.PortAddr, 0, len(nodeCtx.portsUsage.in))

	// in valid program all inports are used, so it's safe to depend on nodeCtx and not use flow's IO
	// actually we can't use IO because we need to know how many slots are used
	for addr := range nodeCtx.portsUsage.in {
		addr := &ir.PortAddr{
			Path: joinNodePath(nodeCtx.path, "in"),
			Port: addr.Port,
			Idx:  uint32(addr.Idx),
		}
		inports = append(inports, *addr)
	}

	sortPortAddrs(inports)

	return inports
}

// sortPortAddrs sorts port addresses by path, port and idx,
// this is very important because runtime function calls depends on this order.
func sortPortAddrs(addrs []ir.PortAddr) {
	sort.Slice(addrs, func(i, j int) bool {
		if addrs[i].Path != addrs[j].Path {
			return addrs[i].Path < addrs[j].Path
		}
		if addrs[i].Port != addrs[j].Port {
			return addrs[i].Port < addrs[j].Port
		}
		return addrs[i].Idx < addrs[j].Idx
	})
}

// for out we use both parent ctx and interface
func (Generator) getFuncOutports(nodeCtx nodeContext) []ir.PortAddr {
	outports := make([]ir.PortAddr, 0, len(nodeCtx.portsUsage.out))

	// In a valid (desugared) program all outports are used so it's safe to depend on nodeCtx and not use flow's IO.
	// Actually we can't use IO because we need to know how many slots are used.
	for addr := range nodeCtx.portsUsage.out {
		irAddr := &ir.PortAddr{
			Path: joinNodePath(nodeCtx.path, "out"),
			Port: addr.Port,
			Idx:  uint32(addr.Idx),
		}

		outports = append(outports, *irAddr)
	}

	sortPortAddrs(outports)

	return outports
}

// mapReceiverSide maps src connection side to ir connection side 1-1 just making the port addr's path absolute
func (g Generator) mapReceiverSide(nodeCtxPath []string, side src.ConnectionReceiver) *ir.PortAddr {
	var idx uint8
	if side.PortAddr.Idx != nil {
		idx = *side.PortAddr.Idx
	}

	result := &ir.PortAddr{
		Path: joinNodePath(nodeCtxPath, side.PortAddr.Node),
		Port: side.PortAddr.Port,
		Idx:  uint32(idx),
	}
	if side.PortAddr.Node == "out" { // 'out' node is actually receiver but we don't want to have 'out.in' addresses
		return result
	}
	result.Path += "/in"

	return result
}

func joinNodePath(nodePath []string, nodeName string) string {
	newPath := make([]string, len(nodePath))
	copy(newPath, nodePath)
	newPath = append(newPath, nodeName)
	return strings.Join(newPath, "/")
}

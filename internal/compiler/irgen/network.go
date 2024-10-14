package irgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nevalang/neva/internal/compiler/ir"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
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
		// sender is always component's inport
		// based on that, we can set receiver's inport slots
		// equal to slots of our own inport
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
				if usageAddr.Port != arrBypassSender.Port { // we only care about the port we're bypassing
					continue
				}

				slotIdxCopy := slotIdx // fix of "Using the variable on range scope" issue

				addr := relPortAddr{Port: arrBypassReceiver.Port, Idx: &slotIdxCopy} // TODO check if pointer is ok to use here
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

			continue
		}

		if len(conn.Normal.SenderSide) != 1 {
			return nil, fmt.Errorf(
				"expected exactly one sender side in desugared network, got %v",
				len(conn.Normal.SenderSide),
			)
		}

		sender := conn.Normal.SenderSide[0]

		irSenderSidePortAddr, err := g.processSenderSide(
			nodeCtx,
			sender,
			nodesPortsUsage,
		)
		if err != nil {
			return nil, fmt.Errorf("process sender side: %w", err)
		}

		for _, receiverSide := range conn.Normal.ReceiverSide.Receivers {
			receiverSideIR := g.mapReceiverSide(nodeCtx.path, receiverSide)

			result.Connections[irSenderSidePortAddr] = receiverSideIR

			// same receiver can be used by multiple senders so we only add it once
			if _, ok := nodesPortsUsage[receiverSide.PortAddr.Node]; !ok {
				nodesPortsUsage[receiverSide.PortAddr.Node] = portsUsage{
					in:  map[relPortAddr]struct{}{},
					out: map[relPortAddr]struct{}{},
				}
			}

			receiverNode := receiverSide.PortAddr.Node
			receiverPortAddr := relPortAddr{
				Port: receiverSide.PortAddr.Port,
				Idx:  receiverSide.PortAddr.Idx,
			}

			nodesPortsUsage[receiverNode].in[receiverPortAddr] = struct{}{}
		}
	}

	return nodesPortsUsage, nil
}

func (g Generator) processSenderSide(
	nodeCtx nodeContext,
	senderSide src.ConnectionSender,
	result map[string]portsUsage,
) (ir.PortAddr, error) {
	// there could be many connections with the same sender but we must only add it once
	if _, ok := result[senderSide.PortAddr.Node]; !ok {
		result[senderSide.PortAddr.Node] = portsUsage{
			in:  map[relPortAddr]struct{}{},
			out: map[relPortAddr]struct{}{},
		}
	}

	// insert outport usage
	result[senderSide.PortAddr.Node].out[relPortAddr{
		Port: senderSide.PortAddr.Port,
		Idx:  senderSide.PortAddr.Idx,
	}] = struct{}{}

	irSenderPort := ir.PortAddr{
		Path: joinNodePath(nodeCtx.path, senderSide.PortAddr.Node),
		Port: senderSide.PortAddr.Port,
	}
	if senderSide.PortAddr.Idx != nil {
		irSenderPort.IsArray = true
		irSenderPort.Idx = *senderSide.PortAddr.Idx
	}

	if senderSide.PortAddr.Node != "in" {
		irSenderPort.Path += "/out"
	}

	return irSenderPort, nil
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

// mapReceiverSide maps src connection side to ir connection side 1-1 just making the port addr's path absolute
func (g Generator) mapReceiverSide(nodeCtxPath []string, side src.ConnectionPortReceiver) ir.PortAddr {
	result := ir.PortAddr{
		Path: joinNodePath(nodeCtxPath, side.PortAddr.Node),
		Port: side.PortAddr.Port,
	}
	if side.PortAddr.Idx != nil {
		result.IsArray = true
		result.Idx = *side.PortAddr.Idx
	}

	// 'out' node is actually receiver but we don't want to have 'out/in' path
	if side.PortAddr.Node != "out" {
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

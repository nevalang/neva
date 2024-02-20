package irgen

import (
	"fmt"
	"strings"

	"github.com/nevalang/neva/pkg/ir"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

// processNet
// 1) inserts network connections
// 2) returns metadata about how subnodes are used by this network
func (g Generator) processNet(
	scope src.Scope,
	conns []src.Connection,
	nodeCtx nodeContext,
	irResult *ir.Program,
) (map[string]portsUsage, error) {
	result := map[string]portsUsage{}

	for _, conn := range conns {
		irSenderSidePortAddr, err := g.processSenderSide(scope, nodeCtx, conn.SenderSide, result)
		if err != nil {
			return nil, fmt.Errorf("process sender side: %w", err)
		}

		receiverSidesIR := make([]ir.ReceiverConnectionSide, 0, len(conn.ReceiverSide.Receivers))
		for _, receiverSide := range conn.ReceiverSide.Receivers {
			receiverSideIR := g.mapReceiverSide(nodeCtx.path, receiverSide)
			receiverSidesIR = append(receiverSidesIR, *receiverSideIR)

			// same receiver can be used by multiple senders so we only add it once
			if _, ok := result[receiverSide.PortAddr.Node]; !ok {
				result[receiverSide.PortAddr.Node] = portsUsage{
					in:  map[relPortAddr]struct{}{},
					out: map[relPortAddr]struct{}{},
				}
			}

			var idx uint8
			if receiverSide.PortAddr.Idx != nil {
				idx = *receiverSide.PortAddr.Idx
			}

			result[receiverSide.PortAddr.Node].in[relPortAddr{
				Port: receiverSide.PortAddr.Port,
				Idx:  idx,
			}] = struct{}{}
		}

		irResult.Connections = append(irResult.Connections, ir.Connection{
			SenderSide:    *irSenderSidePortAddr,
			ReceiverSides: receiverSidesIR,
		})
	}

	return result, nil
}

func (g Generator) processSenderSide(
	scope src.Scope,
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

func (Generator) insertAndReturnInports(
	nodeCtx nodeContext,
	result *ir.Program,
) []ir.PortAddr {
	inports := make([]ir.PortAddr, 0, len(nodeCtx.portsUsage.in))

	// in valid program all inports are used, so it's safe to depend on nodeCtx and not use component's IO
	// actually we can't use IO because we need to know how many slots are used
	for addr := range nodeCtx.portsUsage.in {
		addr := &ir.PortAddr{
			Path: joinNodePath(nodeCtx.path, "in"),
			Port: addr.Port,
			Idx:  uint32(addr.Idx),
		}
		result.Ports = append(result.Ports, ir.PortInfo{
			PortAddr: *addr,
			BufSize:  0,
		})
		inports = append(inports, *addr)
	}

	return inports
}

func (Generator) insertAndReturnOutports(
	outports map[string]src.Port,
	nodeCtx nodeContext,
	result *ir.Program,
) []ir.PortAddr {
	runtimeFuncOutportAddrs := make([]ir.PortAddr, 0, len(nodeCtx.portsUsage.out))

	// In a valid (desugared) program all outports are used so it's safe to depend on nodeCtx and not use component's IO.
	// Actually we can't use IO because we need to know how many slots are used.
	for addr := range nodeCtx.portsUsage.out {
		irAddr := &ir.PortAddr{
			Path: joinNodePath(nodeCtx.path, "out"),
			Port: addr.Port,
			Idx:  uint32(addr.Idx),
		}
		result.Ports = append(result.Ports, ir.PortInfo{
			PortAddr: *irAddr,
			BufSize:  0,
		})
		runtimeFuncOutportAddrs = append(runtimeFuncOutportAddrs, *irAddr)
	}

	return runtimeFuncOutportAddrs
}

// mapReceiverSide maps compiler connection side to ir connection side 1-1 just making the port addr's path absolute
func (g Generator) mapReceiverSide(nodeCtxPath []string, side src.ConnectionReceiver) *ir.ReceiverConnectionSide {
	var idx uint8
	if side.PortAddr.Idx != nil {
		idx = *side.PortAddr.Idx
	}

	result := &ir.ReceiverConnectionSide{
		PortAddr: ir.PortAddr{
			Path: joinNodePath(nodeCtxPath, side.PortAddr.Node),
			Port: side.PortAddr.Port,
			Idx:  uint32(idx),
		},
	}
	if side.PortAddr.Node == "out" { // 'out' node is actually receiver but we don't want to have 'out.in' addresses
		return result
	}
	result.PortAddr.Path += "/in"
	return result
}

func joinNodePath(nodePath []string, nodeName string) string {
	newPath := append(nodePath, nodeName)
	return strings.Join(newPath, "/")
}

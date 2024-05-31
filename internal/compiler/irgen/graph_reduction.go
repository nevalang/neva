package irgen

import "github.com/nevalang/neva/internal/runtime/ir"

// reduceGraph transforms program to a state where it doesn't have intermediate connections.
// It's not optimization, it's a functional requirement of runtime.
func reduceGraph(prog *ir.Program) ([]ir.PortInfo, []ir.Connection) {
	intermediatePorts := map[ir.PortAddr]struct{}{}
	netWithoutIntermediateReceivers := make([]ir.Connection, 0, len(prog.Connections))

	for _, conn := range prog.Connections {
		// it's possible that we already saw this sender as a receiver in previous iterations
		if _, ok := intermediatePorts[conn.SenderSide]; ok {
			continue
		}

		// find final receivers for every intermediate one, also remember all intermediates
		receivers := make([]ir.ReceiverConnectionSide, 0, len(conn.ReceiverSides))
		for _, receiver := range conn.ReceiverSides {
			finalReceivers, wasIntermediate := getFinalReceivers(receiver)
			receivers = append(receivers, finalReceivers...)
			if wasIntermediate {
				intermediatePorts[receiver.PortAddr] = struct{}{}
			}
		}

		// every connection in resultNet has only final receivers
		// (it still might have intermediate ports as senders though)
		netWithoutIntermediateReceivers = append(netWithoutIntermediateReceivers, ir.Connection{
			SenderSide:    conn.SenderSide,
			ReceiverSides: receivers,
		})
	}

	// resultNet only contains connections with final receivers
	// but some of them has senders that are intermediate resultPorts
	// we need to remove these connections and ports for those nodes
	netWithoutIntermediatePorts := make([]ir.Connection, 0, len(netWithoutIntermediateReceivers))
	finalPorts := make([]ir.PortInfo, 0, len(prog.Ports))

	for _, conn := range netWithoutIntermediateReceivers {
		// intermediate receiver is always also a sender in some other connection
		if _, ok := intermediatePorts[conn.SenderSide]; ok {
			continue // skip this connection and don't add its ports
		}

		// basically just add ports for this connection
		finalPorts = append(finalPorts, ir.PortInfo{
			PortAddr: conn.SenderSide,
			BufSize:  0, // TODO https://github.com/nevalang/neva/issues/665
		})

		for _, receiver := range conn.ReceiverSides {
			finalPorts = append(finalPorts, ir.PortInfo{
				PortAddr: receiver.PortAddr,
				BufSize:  0, // TODO https://github.com/nevalang/neva/issues/665
			})
		}

		netWithoutIntermediatePorts = append(netWithoutIntermediatePorts, conn)
	}

	return finalPorts, netWithoutIntermediatePorts
}

// getFinalReceivers returns all final receivers that are behind the given one.
func getFinalReceivers(ir.ReceiverConnectionSide) ([]ir.ReceiverConnectionSide, bool) {
	return nil, false
}

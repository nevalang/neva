package irgen

import "github.com/nevalang/neva/internal/runtime/ir"

// reduceGraph transforms program to a state where it doesn't have intermediate connections.
// It's not optimization, it's a functional requirement of runtime.
func reduceGraph(prog *ir.Program) ([]ir.PortAddr, []ir.Connection) {
	intermediatePorts := map[ir.PortAddr]struct{}{}
	netWithoutIntermediateReceivers := make([]ir.Connection, 0, len(prog.Connections))

	for sender, receivers := range prog.Connections {
		// it's possible that we already saw this sender as a receiver in previous iterations
		if _, ok := intermediatePorts[sender]; ok {
			continue
		}

		// find final receivers for every intermediate one, also remember all intermediates
		finalReceivers := make([]ir.PortAddr, 0, len(receivers))
		for receiver := range receivers {
			curFinalReceivers, wasIntermediate := getFinalReceivers(receiver, prog.Connections)
			finalReceivers = append(finalReceivers, curFinalReceivers...)
			if wasIntermediate {
				intermediatePorts[receiver] = struct{}{}
			}
		}

		// every connection in resultNet has only final receivers
		// (it still might have intermediate ports as senders though)
		netWithoutIntermediateReceivers = append(netWithoutIntermediateReceivers, ir.Connection{
			SenderSide:    sender,
			ReceiverSides: finalReceivers,
		})
	}

	// resultNet only contains connections with final receivers
	// but some of them has senders that are intermediate resultPorts
	// we need to remove these connections and ports for those nodes
	netWithoutIntermediatePorts := make([]ir.Connection, 0, len(netWithoutIntermediateReceivers))
	finalPorts := make([]ir.PortAddr, 0, len(prog.Ports))

	for _, conn := range netWithoutIntermediateReceivers {
		// intermediate receiver is always also a sender in some other connection
		if _, ok := intermediatePorts[conn.SenderSide]; ok {
			continue // skip this connection and don't add its ports
		}

		// basically just add ports for this connection
		finalPorts = append(finalPorts, conn.SenderSide)
		for _, receiver := range conn.ReceiverSides {
			finalPorts = append(finalPorts, receiver)
		}

		netWithoutIntermediatePorts = append(netWithoutIntermediatePorts, conn)
	}

	return finalPorts, netWithoutIntermediatePorts
}

// getFinalReceivers returns all final receivers that are behind the given one.
// It also returns true if given port address was intermediate and false otherwise.
// If given port was already final then a slice with one original port is returned.
func getFinalReceivers(
	receiver ir.PortAddr,
	net map[ir.PortAddr]map[ir.PortAddr]struct{},
) (final []ir.PortAddr, intermediate bool) {
	next, ok := net[receiver]
	if !ok { // given receiver is not a sender for anyone, so it's NOT intermediate port
		return []ir.PortAddr{receiver}, false
	}

	final = make([]ir.PortAddr, 0, len(next))
	for nextReceiver := range next {
		// we don't care about the flag, it's intermediate
		nextNext, _ := getFinalReceivers(nextReceiver, net) // <- recursion
		final = append(final, nextNext...)
	}

	return final, true // yep it was intermediate and here's the finals
}

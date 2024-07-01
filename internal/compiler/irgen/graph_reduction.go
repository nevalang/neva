package irgen

import "github.com/nevalang/neva/internal/compiler/ir"

// reduceGraph transforms program to a state where it doesn't have intermediate connections.
func reduceGraph(prog *ir.Program) (map[ir.PortAddr]struct{}, map[ir.PortAddr]map[ir.PortAddr]struct{}) {
	intermediatePorts := map[ir.PortAddr]struct{}{}

	netWithoutIntermediateReceivers := make(
		map[ir.PortAddr]map[ir.PortAddr]struct{},
		len(prog.Connections),
	)

	for sender, receivers := range prog.Connections {
		// it's possible that we already saw this sender as a receiver in previous iterations
		if _, ok := intermediatePorts[sender]; ok {
			continue
		}

		// find final receivers for every intermediate one, also remember all intermediates
		finalReceivers := make(map[ir.PortAddr]struct{}, len(receivers))
		for receiver := range receivers {
			curFinalReceivers, wasIntermediate := getFinalReceivers(receiver, prog.Connections)
			for _, curFinalReceiver := range curFinalReceivers {
				finalReceivers[curFinalReceiver] = struct{}{}
			}
			if wasIntermediate {
				intermediatePorts[receiver] = struct{}{}
			}
		}

		// every connection in resultNet has only final receivers
		// (it still might have intermediate ports as senders though)
		netWithoutIntermediateReceivers[sender] = finalReceivers
	}

	// resultNet only contains connections with final receivers
	// but some of them has senders that are intermediate resultPorts
	// we need to remove these connections and ports for those nodes
	finalPorts := make(map[ir.PortAddr]struct{}, len(prog.Ports))
	netWithoutIntermediatePorts := make(
		map[ir.PortAddr]map[ir.PortAddr]struct{},
		len(netWithoutIntermediateReceivers),
	)

	for sender, receivers := range netWithoutIntermediateReceivers {
		// intermediate receiver is always also a sender in some other connection
		if _, ok := intermediatePorts[sender]; ok {
			continue // skip this connection and don't add its ports
		}

		netWithoutIntermediatePorts[sender] = receivers

		// basically just add ports for every nonskipped connection
		finalPorts[sender] = struct{}{}
		for receiver := range receivers {
			finalPorts[receiver] = struct{}{}
		}
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

package irgen

import "github.com/nevalang/neva/internal/compiler/ir"

// reduceFinalGraph transforms program to a state where it doesn't have intermediate connections.
func (Generator) reduceFinalGraph(connections map[ir.PortAddr]ir.PortAddr) map[ir.PortAddr]ir.PortAddr {
	intermediatePorts := map[ir.PortAddr]struct{}{}

	netWithoutIntermediateReceivers := make(
		map[ir.PortAddr]ir.PortAddr,
		len(connections),
	)

	for sender, receiver := range connections {
		// it's possible that we already saw this sender as a receiver in previous iterations
		if _, ok := intermediatePorts[sender]; ok {
			continue
		}

		curFinalReceiver, wasIntermediate := getFinalReceiver(receiver, connections)
		if wasIntermediate {
			intermediatePorts[receiver] = struct{}{}
		}

		// every connection in resultNet has final receiver
		// (it still might have intermediate ports as senders though)
		netWithoutIntermediateReceivers[sender] = curFinalReceiver
	}

	// resultNet only contains connections with final receivers
	// but some of them has senders that are intermediate resultPorts
	// we need to remove these connections and ports for those nodes
	netWithoutIntermediatePorts := make(
		map[ir.PortAddr]ir.PortAddr,
		len(netWithoutIntermediateReceivers),
	)

	for sender, receiver := range netWithoutIntermediateReceivers {
		// intermediate receiver is always also a sender in some other connection
		if _, ok := intermediatePorts[sender]; ok {
			continue // skip this connection and don't add its ports
		}
		netWithoutIntermediatePorts[sender] = receiver
	}

	return netWithoutIntermediatePorts
}

// getFinalReceiver returns all final receivers that are behind the given one.
// It also returns true if given port address was intermediate and false otherwise.
// If given port was already final then a slice with one original port is returned.
func getFinalReceiver(
	receiver ir.PortAddr,
	net map[ir.PortAddr]ir.PortAddr,
) (final ir.PortAddr, intermediate bool) {
	next, isSender := net[receiver]
	if !isSender {
		return receiver, false
	}
	v, _ := getFinalReceiver(next, net)
	return v, true
}

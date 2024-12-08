package ir

// GraphReduction removes intermediate connections and returns reduced graph.
func GraphReduction(connections map[PortAddr]PortAddr) map[PortAddr]PortAddr {
	intermediatePorts := map[PortAddr]struct{}{}

	withoutIntermediateReceivers := make(
		map[PortAddr]PortAddr,
		len(connections),
	)

	// after this loop we'll get net where all senders have final receivers
	// but senders themselves may still be intermediate
	for sender, receiver := range connections {
		curFinalReceiver, wasIntermediate := getFinalReceiver(receiver, connections)
		if wasIntermediate {
			intermediatePorts[receiver] = struct{}{}
		}
		withoutIntermediateReceivers[sender] = curFinalReceiver
	}

	// second pass: remove connections with intermediate senders
	result := make(
		map[PortAddr]PortAddr,
		len(withoutIntermediateReceivers),
	)
	for sender, receiver := range withoutIntermediateReceivers {
		if _, isIntermediate := intermediatePorts[sender]; !isIntermediate {
			result[sender] = receiver
		}
	}

	return result
}

// getFinalReceiver returns the final receiver for a given port address.
// It also returns true if the given port address was intermediate, false otherwise.
func getFinalReceiver(
	receiver PortAddr,
	connections map[PortAddr]PortAddr,
) (final PortAddr, intermediate bool) {
	visited := make(map[PortAddr]struct{})
	current := receiver
	for {
		visited[current] = struct{}{}
		next, exists := connections[current]
		if !exists {
			return current, len(visited) > 1
		}
		if _, alreadyVisited := visited[next]; alreadyVisited {
			return current, true // we've detected a cycle, return the current node
		}
		current = next
	}
}

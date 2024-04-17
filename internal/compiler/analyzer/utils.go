package analyzer

type nodesNetUsage map[string]nodeNetUsage

func (n nodesNetUsage) AddOutport(node, port string) {
	if _, ok := n[node]; !ok {
		defaultValue := nodeNetUsage{
			In:  map[string]struct{}{},
			Out: map[string]struct{}{},
		}
		n[node] = defaultValue
	}
	n[node].Out[port] = struct{}{}
}

func (n nodesNetUsage) AddInport(node, port string) {
	if _, ok := n[node]; !ok {
		defaultValue := nodeNetUsage{
			In:  map[string]struct{}{},
			Out: map[string]struct{}{},
		}
		n[node] = defaultValue
	}
	n[node].In[port] = struct{}{}
}

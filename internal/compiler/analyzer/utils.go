package analyzer

import src "github.com/nevalang/neva/internal/compiler/sourcecode"

type netNodesUsage map[string]netNodeUsage

type netNodeUsage struct {
	In, Out map[string]struct{}
}

func (n netNodesUsage) AddOutport(addr src.PortAddr) {
	if _, ok := n[addr.Node]; !ok {
		n[addr.Node] = netNodeUsage{
			In:  map[string]struct{}{},
			Out: map[string]struct{}{},
		}
	}
	n[addr.Node].Out[addr.Port] = struct{}{}
}

func (n netNodesUsage) AddInport(addr src.PortAddr) {
	if _, ok := n[addr.Node]; !ok {
		n[addr.Node] = netNodeUsage{
			In:  map[string]struct{}{},
			Out: map[string]struct{}{},
		}
	}
	n[addr.Node].In[addr.Port] = struct{}{}
}

package portgen

import (
	"fmt"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type PortGen struct{}

func (p PortGen) Ports(nodes map[string]runtime.Node) (map[runtime.PortAddr]chan core.Msg, error) {
	nodesPorts := make(map[runtime.PortAddr]chan core.Msg, len(nodes)*2)

	for name, node := range nodes {
		nodePorts, err := p.nodePorts(name, node.IO)
		if err != nil {
			return nil, fmt.Errorf("node ports: %w", err)
		}

		nodesPorts, err = p.mergePorts(nodesPorts, nodePorts)
		if err != nil {
			return nil, fmt.Errorf("merge ports: %w", err)
		}
	}

	return nodesPorts, nil
}

func (p PortGen) nodePorts(string, runtime.NodeIO) (map[runtime.PortAddr]chan core.Msg, error) {
	return nil, nil
}

func (p PortGen) mergePorts(x, y map[runtime.PortAddr]chan core.Msg) (map[runtime.PortAddr]chan core.Msg, error) {
	return nil, nil
}
